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

package datasource

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/seata/seata-go/pkg/datasource/proxy"
	"github.com/seata/seata-go/pkg/datasource/types"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/seata/seata-go/pkg/config"
	"github.com/seata/seata-go/pkg/datasource/resource_manager"
	"github.com/seata/seata-go/pkg/util/log"
	"github.com/seata/seata-go/pkg/util/reflectx"
)

const (
	// SeataATMySQLDriver MySQL driver for AT mode
	SeataATMySQLDriver = "seata-at-mysql"
	// SeataXAMySQLDriver MySQL driver for XA mode
	SeataXAMySQLDriver = "seata-xa-mysql"
)

func init() {
	sql.Register(SeataATMySQLDriver, &seataATDriver{
		seataDriver: &seataDriver{
			transType: types.ATMode,
			target:    mysql.MySQLDriver{},
		},
	})
	sql.Register(SeataXAMySQLDriver, &seataXADriver{
		seataDriver: &seataDriver{
			transType: types.XAMode,
			target:    mysql.MySQLDriver{},
		},
	})
}

type dsnConnector struct {
	dsn    string
	driver driver.Driver
}

func (t *dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return t.driver.Open(t.dsn)
}

func (t *dsnConnector) Driver() driver.Driver {
	return t.driver
}

type seataATDriver struct {
	*seataDriver
}

func (d *seataATDriver) OpenConnector(name string) (c driver.Connector, err error) {
	connector, err := d.seataDriver.OpenConnector(name)
	if err != nil {
		return nil, err
	}

	tmpConnector, _ := connector.(*seataConnector)
	tmpConnector.transType = types.ATMode

	return &seataATConnector{
		seataConnector: tmpConnector,
	}, nil
}

type seataXADriver struct {
	*seataDriver
}

func (d *seataXADriver) OpenConnector(name string) (c driver.Connector, err error) {
	connector, err := d.seataDriver.OpenConnector(name)
	if err != nil {
		return nil, err
	}

	tmpConnector, _ := connector.(*seataConnector)
	tmpConnector.transType = types.XAMode

	return &seataXAConnector{
		seataConnector: tmpConnector,
	}, nil
}

type seataDriver struct {
	conf      config.Config
	transType types.TransactionType
	target    driver.Driver

	sourceManager *resource_manager.ResourceManager
}

func (d *seataDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.target.Open(name)
	if err != nil {
		log.Errorf("open target connection: %w", err)
		return nil, err
	}

	v := reflect.ValueOf(conn)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByName("connector")
	proxy, err := d.OpenConnector(name)
	if err != nil {
		log.Errorf("open connector: %w", err)
		return nil, err
	}

	reflectx.SetUnexportedField(field, proxy)
	return conn, nil
}

func (d *seataDriver) OpenConnector(name string) (c driver.Connector, err error) {
	c = &dsnConnector{dsn: name, driver: d.target}

	if driverCtx, ok := d.target.(driver.DriverContext); ok {
		c, err = driverCtx.OpenConnector(name)
		if err != nil {
			log.Errorf("open connector: %w", err)
			return nil, err
		}
	}

	var dbType types.DBType
	switch strings.ToLower(d.getTargetDriverName()) {
	case "mysql":
		dbType = types.DbTypeMySQL
	default:
		return nil, fmt.Errorf("unsupport conn type %s", d.getTargetDriverName())
	}

	db := sql.OpenDB(c)
	resourceID := parseResourceID(name)
	dbResource, err := proxy.NewResource(prometheus.DefaultRegisterer,
		&d.conf, proxy.WithTarget(db), proxy.WithResourceID(resourceID), proxy.WithDBType(dbType))
	if err != nil {
		log.Errorf("create new resource: %w", err)
		return nil, err
	}

	dataSourceManager, exist := d.sourceManager.GetDataSourceManager(d.conf.BranchType)
	if !exist {
		log.Error("get datasource manager branch type:%d err:%v", d.conf.BranchType, err)
		return nil, err
	}
	dataSourceManager.RegisterResource(dbResource)

	return &seataConnector{
		res:    dbResource,
		target: c,
		conf:   &d.conf,
	}, nil
}

func (d *seataDriver) getTargetDriverName() string {
	return "mysql"
}

func parseResourceID(dsn string) string {
	i := strings.Index(dsn, "?")

	res := dsn

	if i > 0 {
		res = dsn[:i]
	}

	return strings.ReplaceAll(res, ",", "|")
}
