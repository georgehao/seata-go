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
	"fmt"

	getty "github.com/apache/dubbo-getty"
	gxtime "github.com/dubbogo/gost/time"
	"go.uber.org/atomic"

	"github.com/seata/seata-go/pkg/protocol/codec"
	"github.com/seata/seata-go/pkg/protocol/message"
	"github.com/seata/seata-go/pkg/util/log"
)

type GettyRemotingClient struct {
	idGenerator *atomic.Uint32
	remoting    *GettyRemoting
}

func newGettyRemotingClient(loadBalanceType string, smgr *SessionManager) *GettyRemotingClient {
	return &GettyRemotingClient{
		idGenerator: &atomic.Uint32{},
		remoting:    newGettyRemotingInstance(loadBalanceType, smgr),
	}
}

func (client *GettyRemotingClient) SendAsyncRequest(msg interface{}, session getty.Session) error {
	var msgType message.GettyRequestType
	if _, ok := msg.(message.HeartBeatMessage); ok {
		msgType = message.GettyRequestTypeHeartbeatRequest
	} else {
		msgType = message.GettyRequestTypeRequestOneway
	}
	rpcMessage := message.RpcMessage{
		ID:         int32(client.idGenerator.Inc()),
		Type:       msgType,
		Codec:      byte(codec.CodecTypeSeata),
		Compressor: 0,
		Body:       msg,
	}
	return client.remoting.SendASync(rpcMessage, session, client.asyncCallback)
}

func (client *GettyRemotingClient) SendAsyncResponse(msgID int32, msg interface{}) error {
	rpcMessage := message.RpcMessage{
		ID:         msgID,
		Type:       message.GettyRequestTypeResponse,
		Codec:      byte(codec.CodecTypeSeata),
		Compressor: 0,
		Body:       msg,
	}
	return client.remoting.SendASync(rpcMessage, nil, nil)
}

func (client *GettyRemotingClient) SendSyncRequest(msg interface{}) (interface{}, error) {
	rpcMessage := message.RpcMessage{
		ID:         int32(client.idGenerator.Inc()),
		Type:       message.GettyRequestTypeRequestSync,
		Codec:      byte(codec.CodecTypeSeata),
		Compressor: 0,
		Body:       msg,
	}
	return client.remoting.SendSync(rpcMessage, nil, client.syncCallback)
}

func (client *GettyRemotingClient) asyncCallback(reqMsg message.RpcMessage, respMsg *message.MessageFuture) (interface{}, error) {
	go client.syncCallback(reqMsg, respMsg)
	return nil, nil
}

func (client *GettyRemotingClient) syncCallback(reqMsg message.RpcMessage, respMsg *message.MessageFuture) (interface{}, error) {
	select {
	case <-gxtime.GetDefaultTimerWheel().After(RpcRequestTimeout):
		client.remoting.RemoveMergedMessageFuture(reqMsg.ID)
		log.Errorf("wait resp timeout: %#v", reqMsg)
		return nil, fmt.Errorf("wait response timeout, request: %#v", reqMsg)
	case <-respMsg.Done:
		return respMsg.Response, respMsg.Err
	}
}

func (client *GettyRemotingClient) GetMergedMessage(msgID int32) *message.MergedWarpMessage {
	return client.remoting.GetMergedMessage(msgID)
}

func (client *GettyRemotingClient) GetMessageFuture(msgID int32) *message.MessageFuture {
	return client.remoting.GetMessageFuture(msgID)
}

func (client *GettyRemotingClient) RemoveMessageFuture(msgID int32) {
	client.remoting.RemoveMessageFuture(msgID)
}

func (client *GettyRemotingClient) RemoveMergedMessageFuture(msgID int32) {
	client.remoting.RemoveMergedMessageFuture(msgID)
}

func (client *GettyRemotingClient) NotifyRpcMessageResponse(rpcMessage message.RpcMessage) {
	client.remoting.NotifyRpcMessageResponse(rpcMessage)
}
