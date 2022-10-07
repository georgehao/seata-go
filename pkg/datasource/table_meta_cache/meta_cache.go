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

package table_meta_cache

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// TableMetaCache tables metadata cache
type TableMetaCache interface {
	GetCacheKey() string
	FetchSchema()
}

type TableMetaCacheConf struct {
	MetaCacheExpireTime      time.Duration `yaml:"meta_cache_expire_time" json:"meta_cache_expire_time"`
	MetaCacheCleanInterval   time.Duration `yaml:"meta_cache_clean_interval" json:"meta_cache_clean_interval"`
	MetaCacheRefreshInterval time.Duration `yaml:"meta_cache_refresh_interval" json:"meta_cache_refresh_interval"`
	EnableMetaCacheRefresh   bool          `yaml:"enable_meta_cache_refresh" json:"enable_meta_cache_refresh"`
}

func (cfg *TableMetaCacheConf) RegisterFlags(f *flag.FlagSet) {
	f.DurationVar(&cfg.MetaCacheExpireTime, "table-meta-cache.expire.time", 900*1000*time.Second, "meta cache expire time")
	f.DurationVar(&cfg.MetaCacheCleanInterval, "table-meta-cache.clean.interval", 5*time.Minute, "meta cache clean interval")
	f.DurationVar(&cfg.MetaCacheRefreshInterval, "meta-cache-refresh-checker.refresh.interval", 5*time.Minute, "meta cache refresh checker interval")
	f.BoolVar(&cfg.EnableMetaCacheRefresh, "meta-cache-refresh-checker.refresh.enable", true, "meta cache refresh enable")
}

// TableMetaCacheManager the base table meta cache, all meta cache need use it
type TableMetaCacheManager struct {
	conf  TableMetaCacheConf
	Cache *cache.Cache
}

// NewTableMetaCacheManager create table meta cache
func NewTableMetaCacheManager() (*TableMetaCacheManager, error) {
	metaCache := &TableMetaCacheManager{}
	metaCache.Cache = cache.New(metaCache.conf.MetaCacheExpireTime, metaCache.conf.MetaCacheCleanInterval)

	metaCache.refresh()

	return metaCache, nil
}

// refresh
func (manager *TableMetaCacheManager) refresh() {
	if !manager.conf.EnableMetaCacheRefresh {
		return
	}

	ticker := time.NewTicker(manager.conf.MetaCacheRefreshInterval)
	defer ticker.Stop()
	for range ticker.C {
		// TODO
	}
}

// GetTableMeta get the table cache
func (manager *TableMetaCacheManager) GetTableMeta(context context.Context, metaCache TableMetaCache, tableName, resourceId string) (*schema.TableMeta, error) {
	if len(tableName) == 0 {
		return nil, errors.New("TableMeta cannot be fetched without tableName")
	}

	key := metaCache.GetCacheKey()
	imeta, ok := manager.Cache.Get(key)
	if !ok || imeta == nil {
		// todo get xid
		err := fmt.Errorf("xid:%s get table meta failed, please check whether the table `%s` exists", "xid", tableName)
		return nil, err
	}

	meta, ok := imeta.(*schema.TableMeta)
	if !ok {
		// todo get xid
		err := fmt.Errorf("xid:%s tableName %s get table meta failed, check meta type err", "xid", tableName)
		return nil, err
	}
	return meta, nil
}
