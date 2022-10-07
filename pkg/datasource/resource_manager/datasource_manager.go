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
	"github.com/prometheus/client_golang/prometheus"

	"github.com/seata/seata-go/pkg/constant"
	"github.com/seata/seata-go/pkg/protocol/message"
	"github.com/seata/seata-go/pkg/rm"
)

type DataSourceManager interface {
	RegisterResource(resource rm.Resource) error
	UnregisterResource(resource rm.Resource) error
	GetManagedResources() map[string]rm.Resource
	BranchRollback(ctx context.Context, req message.BranchRollbackRequest) (int, error)
	BranchCommit(ctx context.Context, req message.BranchCommitRequest) (int, error)
	LockQuery(ctx context.Context, req message.GlobalLockQueryRequest) (bool, error)
	BranchRegister(ctx context.Context, clientId string, req message.BranchRegisterRequest) (int64, error)
	BranchReport(ctx context.Context, req message.BranchReportRequest) error
}

// ResourceManager all the branch type resource manager
type ResourceManager struct {
	managers map[int]DataSourceManager

	prom prometheus.Registerer
}

func NewResourceManagers() *ResourceManager {
	return &ResourceManager{
		managers: make(map[int]DataSourceManager),
	}
}

// RegisterResourceManager register a resource manager
func (r *ResourceManager) RegisterResourceManager(prom prometheus.Registerer, branchTypes ...int) error {
	var err error
	for _, branchType := range branchTypes {
		if _, ok := r.managers[branchType]; ok {
			continue
		}

		switch branchType {
		case constant.BranchTypeAT:
			r.managers[constant.BranchTypeAT] = NewATSourceManager(prom)
		case constant.BranchTypeXA:
		default:
			err = errors.New("not support resource manager type")
		}
	}
	return err
}

func (r *ResourceManager) GetDataSourceManager(branchType int) (DataSourceManager, bool) {
	datasourceManager, exist := r.managers[branchType]
	return datasourceManager, exist
}
