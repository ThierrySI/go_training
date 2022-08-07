package database

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"time"
)

type Conn = driver.Conn
type Ctx = context.Context

// Clickhouse clickhouse connection and basic info holder
type ClickhouseDB struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DatabaseName string
	Connection   Conn
	Context      Ctx
}

// return CH object with connection to local CH
func ConnectCH(host string, port int, user string, password string, database string, maxOpenCnx int, maxIdleCnx int) (*ClickhouseDB, error) {

	var clusterCH = &ClickhouseDB{Host: host, Port: port, UserName: user, Password: password, DatabaseName: database}
	addressCH := fmt.Sprintf("%v:%v", clusterCH.Host, clusterCH.Port)

	cnx, err := clickhouse.Open(
		&clickhouse.Options{
			Addr: []string{addressCH},
			Auth: clickhouse.Auth{
				Database: clusterCH.DatabaseName,
				Username: clusterCH.UserName,
				Password: clusterCH.Password,
			},
			// Debug:           true,
			DialTimeout:     time.Second,
			MaxOpenConns:    maxOpenCnx,
			MaxIdleConns:    maxIdleCnx,
			ConnMaxLifetime: time.Hour,
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
		})
	if err != nil {
		fmt.Printf("[Connect] Error: %v\n", err)
		return &ClickhouseDB{}, err
	}

	ctx := clickhouse.Context(context.Background(),
		clickhouse.WithSettings(clickhouse.Settings{
			"max_block_size": 10,
		}))

	if err := cnx.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return &ClickhouseDB{}, err
	}

	clusterCH.Context = ctx
	clusterCH.Connection = cnx

	return clusterCH, err
}

func CloseDB(c *ClickhouseDB) {
	c.Connection.Close()
}

//for any DDL command execution (Data Definition Language as CREATE, DROP, ALTER, etc ...)
func (c *ClickhouseDB) ExecuteDDL(query string) error {
	return c.Connection.Exec(c.Context, query)
}
