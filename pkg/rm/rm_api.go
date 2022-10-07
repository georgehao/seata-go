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

package rm

import (
	"context"
	"sync"
)

// Resource that can be managed by Resource Manager and involved into global transaction
type Resource interface {
	GetResourceGroupId() string
	GetResourceId() string
	GetBranchType() int
}

// BranchResource contains branch to commit or rollback
type BranchResource struct {
	BranchType      int
	Xid             string
	BranchId        int64
	ResourceId      string
	ApplicationData []byte
}

// ResourceManagerInbound Control a branch transaction commit or rollback
type ResourceManagerInbound interface {
	// BranchCommit Commit a branch transaction
	BranchCommit(ctx context.Context, resource BranchResource) (int, error)
	// BranchRollback Rollback a branch transaction
	BranchRollback(ctx context.Context, resource BranchResource) (int, error)
}

// BranchRegisterParam Branch register function param for ResourceManager
type BranchRegisterParam struct {
	BranchType      int
	ResourceId      string
	ClientId        string
	Xid             string
	ApplicationData string
	LockKeys        string
}

// BranchReportParam Branch report function param for ResourceManager
type BranchReportParam struct {
	BranchType      int
	Xid             string
	BranchId        int64
	Status          int
	ApplicationData string
}

// LockQueryParam Lock query function param for ResourceManager
type LockQueryParam struct {
	BranchType int
	ResourceId string
	Xid        string
	LockKeys   string
}

// ResourceManagerOutbound send outbound request to TC
type ResourceManagerOutbound interface {
	// BranchRegister register long
	BranchRegister(ctx context.Context, param BranchRegisterParam) (int64, error)
	// BranchReport Branch report
	BranchReport(ctx context.Context, param BranchReportParam) error
	// LockQuery Lock query boolean
	LockQuery(ctx context.Context, param LockQueryParam) (bool, error)
}

// ResourceManager Resource Manager: common behaviors
type ResourceManager interface {
	ResourceManagerInbound
	ResourceManagerOutbound

	// RegisterResource Register a Resource to be managed by Resource Manager
	RegisterResource(resource Resource) error
	// UnregisterResource Unregister a Resource from the Resource Manager
	UnregisterResource(resource Resource) error
	// GetCachedResources Get all resources managed by this manager
	GetCachedResources() *sync.Map
	// GetBranchType Get the BranchType
	GetBranchType() int
}

type ResourceManagerGetter interface {
	GetResourceManager(branchType int) ResourceManager
}
