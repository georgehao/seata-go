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

package sql

import (
	"context"
	"database/sql/driver"
	"errors"
	"io"
	"sync"

	"github.com/go-sql-driver/mysql"

	"github.com/seata/seata-go/pkg/datasource/sql/types"
	"github.com/seata/seata-go/pkg/util/log"
)

type seataATConnector struct {
	*seataConnector
	transType types.TransactionMode
}

func (c *seataATConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.seataConnector.Connect(ctx)
	if err != nil {
		return nil, err
	}

	_conn, _ := conn.(*Conn)

	return &ATConn{
		Conn: _conn,
	}, nil
}

func (c *seataATConnector) Driver() driver.Driver {
	return &seataATDriver{
		seataDriver: c.seataConnector.Driver().(*seataDriver),
	}
}

type seataXAConnector struct {
	*seataConnector
	transType types.TransactionMode
}

func (c *seataXAConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.seataConnector.Connect(ctx)
	if err != nil {
		return nil, err
	}

	_conn, _ := conn.(*Conn)

	return &XAConn{
		Conn: _conn,
	}, nil
}

func (c *seataXAConnector) Driver() driver.Driver {
	return &seataXADriver{
		seataDriver: c.seataConnector.Driver().(*seataDriver),
	}
}

// A Connector represents a driver in a fixed configuration
// and can create any number of equivalent Conns for use
// by multiple goroutines.
//
// A Connector can be passed to sql.OpenDB, to allow drivers
// to implement their own sql.DB constructors, or returned by
// DriverContext's OpenConnector method, to allow drivers
// access to context and to avoid repeated parsing of driver
// configuration.
//
// If a Connector implements io.Closer, the sql package's DB.Close
// method will call Close and return error (if any).
type seataConnector struct {
	transType types.TransactionMode
	res       *DBResource
	once      sync.Once
	driver    driver.Driver
	target    driver.Connector
	cfg       *mysql.Config
}

// Connect returns a connection to the database.
// Connect may return a cached connection (one previously
// closed), but doing so is unnecessary; the sql package
// maintains a pool of idle connections for efficient re-use.
//
// The provided context.Context is for dialing purposes only
// (see net.DialContext) and should not be stored or used for
// other purposes. A default timeout should still be used
// when dialing as a connection pool may call Connect
// asynchronously to any query.
//
// The returned connection is only used by one goroutine at a
// time.
func (c *seataConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.target.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// get the version of mysql for xa.
	if c.transType == types.XAMode {
		version, err := c.dbVersion(ctx, conn)
		if err != nil {
			return nil, err
		}
		c.res.SetDbVersion(version)
	}

	return &Conn{
		targetConn: conn,
		res:        c.res,
		txCtx:      types.NewTxCtx(),
		autoCommit: true,
		dbName:     c.cfg.DBName,
		dbType:     types.DBTypeMySQL,
	}, nil
}

func (c *seataConnector) dbVersion(ctx context.Context, conn driver.Conn) (string, error) {
	queryConn, isQueryContext := conn.(driver.QueryerContext)
	if !isQueryContext {
		return "", errors.New("get db version error for unexpected driver conn")
	}

	res, err := queryConn.QueryContext(ctx, "SELECT VERSION()", nil)
	if err != nil {
		log.Errorf("seata connector get the xa mysql version err:%v", err)
		return "", err
	}

	dest := make([]driver.Value, 1)
	var version string
	for true {
		if err = res.Next(dest); err != nil {
			if err == io.EOF {
				return version, nil
			}
			return "", err
		}
		if len(dest) != 1 {
			return "", errors.New("get the mysql version is not column 1")
		}

		var isVersionOk bool
		version, isVersionOk = dest[0].(string)
		if !isVersionOk {
			return "", errors.New("get the mysql version is not a string")
		}
	}
	return "", errors.New("get the mysql version is error")
}

// Driver returns the underlying Driver of the Connector,
// mainly to maintain compatibility with the Driver method
// on sql.DB.
func (c *seataConnector) Driver() driver.Driver {
	c.once.Do(func() {
		d := c.target.Driver()
		c.driver = d
	})

	return &seataDriver{
		target: c.driver,
	}
}
