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

package proxy

import (
	gosql "database/sql"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/seata/seata-go/pkg/config"
	"github.com/seata/seata-go/pkg/datasource/table_meta_cache"
	"github.com/seata/seata-go/pkg/datasource/types"
	"github.com/seata/seata-go/pkg/datasource/undo"
	"github.com/seata/seata-go/pkg/util/log"
)

type dbOption func(db *DBResource)

func WithResourceID(id string) dbOption {
	return func(db *DBResource) {
		db.ResourceID = id
	}
}

func WithDBType(dt types.DBType) dbOption {
	return func(db *DBResource) {
		db.DBType = dt
	}
}

func WithBranchType(branchType int) dbOption {
	return func(db *DBResource) {
		db.BranchType = branchType
	}
}

func WithTarget(source *gosql.DB) dbOption {
	return func(db *DBResource) {
		db.Target = source
	}
}

// DBResource proxy sql.DB, enchance database/sql.DB to add distribute transaction ability
type DBResource struct {
	Conf *config.Config

	GroupID    string
	ResourceID string
	Target     *gosql.DB
	DBType     types.DBType
	BranchType int

	undoLogMgr       undo.UndoLogManager
	metaCacheManager *table_meta_cache.TableMetaCacheManager
}

func NewResource(prom prometheus.Registerer, conf *config.Config, opts ...dbOption) (*DBResource, error) {
	db := &DBResource{
		Conf: conf,
	}

	for i := range opts {
		opts[i](db)
	}

	// init table meta cache
	metaCacheManager, err := table_meta_cache.NewTableMetaCacheManager()
	if err != nil {
		log.Fatalf("proxy db resource init meta cache manager err:%v", err)
		return nil, err
	}
	db.metaCacheManager = metaCacheManager

	return db, nil
}

func (resource *DBResource) GetResourceGroupId() string {
	return ""
}

func (resource *DBResource) GetResourceId() string {
	return ""
}

func (resource *DBResource) GetBranchType() int {
	return 0
}
