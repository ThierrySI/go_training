package album

import (
	"go_training/web-service/pkg/config"
	"go_training/web-service/pkg/database"
	"testing"
)

var albumsReference = []Album{
	{Code: "1", Title: "Abbey Road", Artist: "The Beatles", Price: 56.99},
	{Code: "2", Title: "What's Going On", Artist: "Marvin Gaye", Price: 17.99},
	{Code: "3", Title: "Purple Rain", Artist: "Prince", Price: 39.99},
	// {Code: "4", Title: "Bohemian Rhapsod", Artist: "Queen", Price: 12.99},
}

func setup() *database.ClickhouseDB {
	config := config.LoadConfigFile()
	// 2. DB Pool
	db, _ := database.ConnectCH(config.DBHostname, config.DBPort, config.DBUsername, config.DBPassword, config.DBDatabase, config.DBMaxOpenConn, config.DBMaxIdleConn)
	return db

}

func teardown(d *database.ClickhouseDB) {
	d.Connection.Close()
}

func TestCreateTableAndData(t *testing.T) {
	d := setup()
	h := &handler{DB: d}

	albumsReturned := h.createTableAndData()
	compareAlbums(albumsReturned, albumsReference, t)

	teardown(d)
}

func TestSelectAlbum(t *testing.T) {
	d := setup()
	h := &handler{DB: d}

	code := "1"
	albumsReturned := h.selectAlbum(code)
	compareAlbums(albumsReturned, albumsReference[:1], t)

	teardown(d)
}

func TestSelectAlbums(t *testing.T) {
	d := setup()
	h := &handler{DB: d}

	albumsReturned := h.selectAlbums()
	compareAlbums(albumsReturned, albumsReference, t)

	teardown(d)
}

func TestInsertNewAlbum(t *testing.T) {
	d := setup()
	h := &handler{DB: d}

	newItem := Album{Code: "4", Title: "Bohemian Rhapsod", Artist: "Queen", Price: 12.99}
	albumsAdded := h.insertNewAlbum(newItem)

	code := "4"
	albumsReturned := h.selectAlbum(code)
	compareAlbums(albumsReturned, albumsAdded, t)

	newAlbumsReference := append(albumsReference, newItem)
	newAlbumsReturned := h.selectAlbums()
	compareAlbums(newAlbumsReturned, newAlbumsReference, t)

	teardown(d)
}

func TestDeleteAlbum(t *testing.T) {
	d := setup()
	h := &handler{DB: d}

	code := "4"
	albumsReturned := h.selectAlbum(code)
	albumsDeleted := h.deleteAlbum(code)
	compareAlbums(albumsDeleted, albumsReturned, t)

	albumsReturned = h.selectAlbums()
	compareAlbums(albumsReturned, albumsReference, t)

	teardown(d)
}

func compareAlbums(albumsReturned []Album, albumsReference []Album, t *testing.T) {

	t.Logf("%s | LOG | Number of Album | received:%d, wanted:%d", t.Name(), len(albumsReturned), len(albumsReference))

	if len(albumsReference) != len(albumsReturned) {
		t.Errorf("%s | ERROR | Number of Album | got:%d, wanted:%d", t.Name(), len(albumsReturned), len(albumsReference))
	}

	for i, v := range albumsReturned {

		if v.Code != albumsReference[i].Code {
			t.Errorf("%s | ERROR | Album Code in raw | got:%s, wanted:%s", t.Name(), v.Code, albumsReference[i].Code)
		}

		if v.Title != albumsReference[i].Title {
			t.Errorf("%s | ERROR | Album Title in raw | got:%s, wanted:%s", t.Name(), v.Title, albumsReference[i].Title)
		}

		if v.Artist != albumsReference[i].Artist {
			t.Errorf("%s | ERROR | Album Artist in raw | got:%s, wanted:%s", t.Name(), v.Artist, albumsReference[i].Artist)
		}

		if v.Price != albumsReference[i].Price {
			t.Errorf(" %s | ERROR | Album Price in raw |got:%v, wanted:%v", t.Name(), v.Price, albumsReference[i].Price)
		}

	}
}
