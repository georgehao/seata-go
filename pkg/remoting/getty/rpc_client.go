/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package getty

import (
	"crypto/tls"
	"fmt"
	"net"

	getty "github.com/apache/dubbo-getty"
	gxsync "github.com/dubbogo/gost/sync"
	"github.com/pkg/errors"

	"github.com/seata/seata-go/pkg/util/log"
	"github.com/seata/seata-go/pkg/config"
)

type RpcClient struct {
	conf config.Config
}

func InitRpcClient(cfg *config.Config) {
	var rpcClient RpcClient
	rpcClient.conf = *cfg

	addressList := getAvailServerList()
	if len(addressList) == 0 {
		log.Fatal("no have valid seata server list")
		return
	}

	for _, address := range addressList {
		gettyClient := getty.NewTCPClient(
			getty.WithServerAddress(address),
			getty.WithConnectionNumber(cfg.GettyConfig.ConnectionNum),
			getty.WithReconnectInterval(cfg.GettyConfig.ReconnectInterval),
			getty.WithClientTaskPool(gxsync.NewTaskPoolSimple(0)),
		)
		go gettyClient.RunEventLoop(rpcClient.newSession)
	}
}

// todo mock
func getAvailServerList() []string {
	return []string{"127.0.0.1:8091"}
}

func (client *RpcClient) newSession(session getty.Session) error {
	var (
		ok      bool
		tcpConn *net.TCPConn
		err     error
	)

	if client.conf.GettyConfig.GettySessionParam.CompressEncoding {
		session.SetCompressType(getty.CompressZip)
	}
	if _, ok = session.Conn().(*tls.Conn); ok {
		session.SetName(client.conf.GettyConfig.GettySessionParam.SessionName)
		session.SetMaxMsgLen(client.conf.GettyConfig.GettySessionParam.MaxMsgLen)
		session.SetPkgHandler(rpcPkgHandler)
		session.SetEventListener(GetGettyClientHandlerInstance(client.conf))
		session.SetReadTimeout(client.conf.GettyConfig.GettySessionParam.TCPReadTimeout)
		session.SetWriteTimeout(client.conf.GettyConfig.GettySessionParam.TCPWriteTimeout)
		session.SetCronPeriod((int)(client.conf.GettyConfig.GettySessionParam.CronPeriod))
		session.SetWaitTime(client.conf.GettyConfig.GettySessionParam.WaitTimeout)
		log.Debugf("server accepts new tls session:%s\n", session.Stat())
		return nil
	}
	if _, ok = session.Conn().(*net.TCPConn); !ok {
		panic(fmt.Sprintf("%s, session.conn{%#v} is not a tcp connection\n", session.Stat(), session.Conn()))
	}

	if _, ok = session.Conn().(*tls.Conn); !ok {
		if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
			return errors.New(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection", session.Stat(), session.Conn()))
		}

		if err = tcpConn.SetNoDelay(client.conf.GettyConfig.GettySessionParam.TCPNoDelay); err != nil {
			return err
		}
		if err = tcpConn.SetKeepAlive(client.conf.GettyConfig.GettySessionParam.TCPKeepAlive); err != nil {
			return err
		}
		if client.conf.GettyConfig.GettySessionParam.TCPKeepAlive {
			if err = tcpConn.SetKeepAlivePeriod(client.conf.GettyConfig.GettySessionParam.KeepAlivePeriod); err != nil {
				return err
			}
		}
		if err = tcpConn.SetReadBuffer(client.conf.GettyConfig.GettySessionParam.TCPRBufSize); err != nil {
			return err
		}
		if err = tcpConn.SetWriteBuffer(client.conf.GettyConfig.GettySessionParam.TCPWBufSize); err != nil {
			return err
		}
	}

	session.SetName(client.conf.GettyConfig.GettySessionParam.SessionName)
	session.SetMaxMsgLen(client.conf.GettyConfig.GettySessionParam.MaxMsgLen)
	session.SetPkgHandler(rpcPkgHandler)
	session.SetEventListener(GetGettyClientHandlerInstance(client.conf))
	session.SetReadTimeout(client.conf.GettyConfig.GettySessionParam.TCPReadTimeout)
	session.SetWriteTimeout(client.conf.GettyConfig.GettySessionParam.TCPWriteTimeout)
	session.SetCronPeriod((int)(client.conf.GettyConfig.GettySessionParam.CronPeriod.Nanoseconds() / 1e6))
	session.SetWaitTime(client.conf.GettyConfig.GettySessionParam.WaitTimeout)
	log.Debugf("rpc_client new session:%s\n", session.Stat())

	return nil
}
