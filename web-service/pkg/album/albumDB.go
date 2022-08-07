package album

import (
	"fmt"
	"go_training/web-service/pkg/database"
	"log"
)

func GetAlbums() []Album {
	albumsReturned, err :=


				selectI    selectItem(Album{})
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
func (c *database.Clickhouse) selectItem(row Album) ([]Album, error) {
	var whereClause string
	if (Album{}) == row {
		whereClause = `WHERE 1=1 ORDER BY code`
	} else {
		whereClause = fmt.Sprintf("WHERE code = '%s' ORDER BY code", row.Code)
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
