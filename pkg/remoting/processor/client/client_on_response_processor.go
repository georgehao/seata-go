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

package client

import (
	"context"

	"github.com/seata/seata-go/pkg/protocol/message"
	"github.com/seata/seata-go/pkg/remoting/getty"
	"github.com/seata/seata-go/pkg/util/log"
)

func initOnResponse(eventLister *getty.GettyClientHandler, gettyClient *getty.GettyRemotingClient) {
	clientOnResponseProcessor := &clientOnResponseProcessor{
		gettyClient: gettyClient,
	}

	eventLister.RegisterProcessor(message.MessageTypeSeataMergeResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeBranchRegisterResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeBranchStatusReportResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeGlobalLockQueryResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeRegRmResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeGlobalBeginResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeGlobalCommitResult, clientOnResponseProcessor)

	eventLister.RegisterProcessor(message.MessageTypeGlobalReportResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeGlobalRollbackResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeGlobalStatusResult, clientOnResponseProcessor)
	eventLister.RegisterProcessor(message.MessageTypeRegCltResult, clientOnResponseProcessor)
}

type clientOnResponseProcessor struct {
	gettyClient *getty.GettyRemotingClient
}

func (f *clientOnResponseProcessor) Process(ctx context.Context, rpcMessage message.RpcMessage) error {
	log.Infof("the rm client received  clientOnResponse msg %#v from tc server.", rpcMessage)
	if mergedResult, ok := rpcMessage.Body.(message.MergeResultMessage); ok {
		mergedMessage := f.gettyClient.GetMergedMessage(rpcMessage.ID)
		if mergedMessage != nil {
			for i := 0; i < len(mergedMessage.Msgs); i++ {
				msgID := mergedMessage.MsgIds[i]
				response := f.gettyClient.GetMessageFuture(msgID)
				if response != nil {
					response.Response = mergedResult.Msgs[i]
					response.Done <- struct{}{}
					f.gettyClient.RemoveMessageFuture(msgID)
				}
			}
			f.gettyClient.RemoveMergedMessageFuture(rpcMessage.ID)
		}
		return nil
	} else {
		// 如果是请求消息，做处理逻辑
		msgFuture := f.gettyClient.GetMessageFuture(rpcMessage.ID)
		if msgFuture != nil {
			f.gettyClient.NotifyRpcMessageResponse(rpcMessage)
			f.gettyClient.RemoveMessageFuture(rpcMessage.ID)
		} else {
			if _, ok := rpcMessage.Body.(message.AbstractResultMessage); ok {
				log.Infof("the rm client received response msg [{}] from tc server.", msgFuture)
			}
		}
	}
	return nil
}
