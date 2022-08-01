package model

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
	"time"
)

type Conn = driver.Conn
type Ctx = context.Context

// Clickhouse clickhouse connection and basic info holder
type Clickhouse struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	DatabaseName string
	Connection   Conn
	Context      Ctx
}

const (
	DBHostname = "127.0.0.1"
	DBPort     = 9001
	DBUsername = "default"
	DBPassword = ""
	DBDatabase = "default"
)

func GetAlbums() []Album {
	db, err := connectCH(DBHostname, DBPort, DBUsername, DBPassword, DBDatabase)
	if err != nil {
		return nil
	}
	defer db.Connection.Close()

	albumsReturned, err := db.selectItem(Album{})
	if err != nil {
		return nil
	} else {
		return albumsReturned
	}

}

func GetAlbumByCode(code string) []Album {
	db, err := connectCH(DBHostname, DBPort, DBUsername, DBPassword, DBDatabase)
	if err != nil {
		return nil
	}
	defer db.Connection.Close()

	var albumAdded Album
	albumAdded.Code = code

	albumsReturned, err := db.selectItem(albumAdded)
	if err != nil {
		return nil
	} else {
		return albumsReturned
	}

}

func AddNewAlbum(row Album) []Album {
	db, err := connectCH(DBHostname, DBPort, DBUsername, DBPassword, DBDatabase)
	if err != nil {
		return nil
	}
	defer db.Connection.Close()

	batch, err := db.Connection.PrepareBatch(db.Context, "INSERT INTO albums (code, title, artist, price)")
	if err != nil {
		return nil
	}
	if err := batch.Append(row.Code, row.Title, row.Artist, row.Price); err != nil {
		return nil
	}
	if err := batch.Send(); err != nil {
		return nil
	}

	return append([]Album{}, row)
}

func DeleteAlbum(code string) []Album {
	db, err := connectCH(DBHostname, DBPort, DBUsername, DBPassword, DBDatabase)
	if err != nil {
		return nil
	}
	defer db.Connection.Close()

	var albumDeleted Album
	albumDeleted.Code = code

	albumsReturned, err := db.selectItem(albumDeleted)
	if err != nil {
		return nil
	} else {
		query := fmt.Sprintf("ALTER TABLE albums DELETE WHERE code = '%s'", code)
		err = db.Connection.Exec(db.Context, query)
		if err != nil {
			return nil
		} else {
			return albumsReturned
		}
	}

}

// return CH object with connection to local CH
func connectCH(host string, port int, user string, password string, database string) (*Clickhouse, error) {

	var clusterCH = &Clickhouse{Host: host, Port: port, UserName: user, Password: password, DatabaseName: database}
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
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
			Compression: &clickhouse.Compression{
				Method: clickhouse.CompressionLZ4,
			},
		})
	if err != nil {
		fmt.Printf("[Connect] Error: %v\n", err)
		return &Clickhouse{}, err
	}

	ctx := clickhouse.Context(context.Background(),
		clickhouse.WithSettings(clickhouse.Settings{
			"max_block_size": 10,
		}))
	/*		, clickhouse.WithProgress(func(p *clickhouse.Progress) {
				fmt.Println("progress: ", p)
			}), clickhouse.WithProfileInfo(func(p *clickhouse.ProfileInfo) {
			fmt.Println("profile info: ", p)
		}))*/

	if err := cnx.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return &Clickhouse{}, err
	}

	clusterCH.Context = ctx
	clusterCH.Connection = cnx

	return clusterCH, err

}

// for any DDL command execution (Data Definition Language as CREATE, DROP, ALTER, etc ...)
func (c *Clickhouse) executeDDL(query string) error {
	return c.Connection.Exec(c.Context, query)
}

func CreateTableAndData() []Album {
	db, err := connectCH(DBHostname, DBPort, DBUsername, DBPassword, DBDatabase)
	if err != nil {
		return nil
	}
	defer db.Connection.Close()

	query := `DROP TABLE IF EXISTS albums`
	err = db.executeDDL(query)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	query = `CREATE TABLE IF NOT EXISTS albums (code String, title String, artist String, price Float32) engine=Memory`
	err = db.executeDDL(query)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	query = `INSERT INTO albums (code, title, artist, price)`
	batch, err := db.Connection.PrepareBatch(db.Context, query)
	if err != nil {
		return nil
	}
	albums := SampleAlbums()
	for _, v := range albums {
		if err := batch.Append(v.Code, v.Title, v.Artist, v.Price); err != nil {
			return nil
		}
	}
	if err := batch.Send(); err != nil {
		return nil
	}

	return albums

}

// select item
func (c *Clickhouse) selectItem(row Album) ([]Album, error) {
	var whereClause string
	if (Album{}) == row {
		whereClause = `WHERE 1=1`
	} else {
		whereClause = fmt.Sprintf("WHERE code = '%s'", row.Code)
	}

	rows, err := c.Connection.Query(c.Context, "SELECT code, title, artist, price FROM albums "+whereClause)
	if err != nil {
		return nil, err
	}
	var albumSelected []Album
	for rows.Next() {
		var (
			albumInRow = Album{}
		)
		if err := rows.Scan(&albumInRow.Code, &albumInRow.Title, &albumInRow.Artist, &albumInRow.Price); err != nil {
			return nil, err
		}
		albumSelected = append(albumSelected, albumInRow)
	}
	// display(albumSelected)
	rows.Close()
	return albumSelected, rows.Err()
}
