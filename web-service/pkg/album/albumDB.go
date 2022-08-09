package album

import (
	"fmt"
	"log"
)

func (h *handler) selectAlbums() []Album {

	albumsReturned, err := h.selectItem(Album{})
	if err != nil {
		return nil
	} else {
		return albumsReturned
	}
}

func (h *handler) selectAlbum(code string) []Album {
	var albumSelected Album
	albumSelected.Code = code

	albumsReturned, err := h.selectItem(albumSelected)
	if err != nil {
		return nil
	} else {
		return albumsReturned
	}
}

// select item
func (h *handler) selectItem(row Album) ([]Album, error) {
	var whereClause string
	if (Album{}) == row {
		whereClause = `WHERE 1=1 ORDER BY code`
	} else {
		whereClause = fmt.Sprintf("WHERE code = '%s' ORDER BY code", row.Code)
	}

	rows, err := h.DB.Connection.Query(h.DB.Context, "SELECT code, title, artist, price FROM albums "+whereClause)
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

	rows.Close()
	return albumSelected, rows.Err()

}

func (h *handler) insertNewAlbum(row Album) []Album {

	batch, err := h.DB.Connection.PrepareBatch(h.DB.Context, "INSERT INTO albums (code, title, artist, price)")
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

func (h *handler) deleteAlbum(code string) []Album {

	var albumDeleted Album
	albumDeleted.Code = code

	displayAlbumDeleted, err := h.selectItem(albumDeleted)
	if err != nil {
		return nil
	} else {
		query := fmt.Sprintf("ALTER TABLE albums DELETE WHERE code = '%s'", code)
		err = h.DB.ExecuteDDL(query)
		if err != nil {
			return nil
		} else {
			return displayAlbumDeleted
		}
	}

}

func (h *handler) createTableAndData() []Album {

	query := `DROP TABLE IF EXISTS albums`
	err := h.DB.ExecuteDDL(query)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	query = `CREATE TABLE IF NOT EXISTS albums (code String, title String, artist String, price Float32) engine=Memory`
	err = h.DB.ExecuteDDL(query)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	query = `INSERT INTO albums (code, title, artist, price)`
	batch, err := h.DB.Connection.PrepareBatch(h.DB.Context, query)
	if err != nil {
		return nil
	}
	albums := sampleAlbums()
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
