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

package resource_manager

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/seata/seata-go/pkg/constant"
	"github.com/seata/seata-go/pkg/datasource/proxy"
	"github.com/seata/seata-go/pkg/datasource/types"
	"github.com/seata/seata-go/pkg/datasource/undo"
	"github.com/seata/seata-go/pkg/protocol/message"
	"github.com/seata/seata-go/pkg/rm"
)

// ATSourceManager The type AT Data source manager.
type ATSourceManager struct {
	BasicSourceManager

	resourceCache sync.Map
	worker        *AsyncWorker
}

// NewATSourceManager create at source manager
func NewATSourceManager(prom prometheus.Registerer) *ATSourceManager {
	sourceManager := &ATSourceManager{
		resourceCache: sync.Map{},
	}

	sourceManager.worker = NewAsyncWorker(prom, sourceManager)
	return sourceManager
}

// RegisterResource register a Resource to be managed by Resource Manager
func (mgr *ATSourceManager) RegisterResource(res rm.Resource) error {
	mgr.resourceCache.Store(res.GetResourceId(), res)
	return mgr.BasicSourceManager.RegisterResource(res)
}

// UnregisterResource unregister a Resource from the Resource Manager
func (mgr *ATSourceManager) UnregisterResource(res rm.Resource) error {
	return errors.New("not support unregister resource")
}

// GetManagedResources get all resources managed by this manager
func (mgr *ATSourceManager) GetManagedResources() map[string]rm.Resource {
	ret := make(map[string]rm.Resource)
	mgr.resourceCache.Range(func(key, value interface{}) bool {
		ret[key.(string)] = value.(rm.Resource)
		return true
	})

	return ret
}

// BranchRollback Rollback the corresponding transactions according to the request
func (mgr *ATSourceManager) BranchRollback(ctx context.Context, req message.BranchRollbackRequest) (int, error) {
	val, ok := mgr.resourceCache.Load(req.ResourceId)
	if !ok {
		return constant.BranchStatusUnknown, fmt.Errorf("resource %s not found", req.ResourceId)
	}

	res := val.(*proxy.DBResource)
	undoMgr, err := undo.GetUndoLogManager(res.DBType)
	if err != nil {
		return constant.BranchStatusUnknown, err
	}

	conn, err := res.Target.Conn(ctx)
	if err != nil {
		return constant.BranchStatusUnknown, err
	}

	if err := undoMgr.RunUndo(req.Xid, req.BranchId, conn); err != nil {
		transErr, ok := err.(*types.TransactionError)
		if !ok {
			return constant.BranchStatusPhaseOneFailed, err
		}

		if transErr.Code() == types.ErrorCodeBranchRollbackFailedUnretriable {
			return constant.BranchStatusPhaseTwoRollbackFailedUnretryable, nil
		}

		return constant.BranchStatusPhaseTwoRollbackFailedRetryable, nil
	}

	return constant.BranchStatusPhaseTwoRollbacked, nil
}

// BranchCommit a branch transaction
func (mgr *ATSourceManager) BranchCommit(ctx context.Context, req message.BranchCommitRequest) (int, error) {
	mgr.worker.BranchCommit(ctx, req)
	return constant.BranchStatusPhaseOneDone, nil
}

// LockQuery query
func (mgr *ATSourceManager) LockQuery(ctx context.Context, req message.GlobalLockQueryRequest) (bool, error) {
	return false, nil
}

// BranchRegister a branch transaction
func (mgr *ATSourceManager) BranchRegister(ctx context.Context, clientId string, req message.BranchRegisterRequest) (int64, error) {
	return 0, nil
}

// BranchReport report the branch transaction
func (mgr *ATSourceManager) BranchReport(ctx context.Context, req message.BranchReportRequest) error {
	return nil
}
