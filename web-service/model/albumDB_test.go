package model

import (
	"testing"
)

var albumsReference = []Album{
	{Code: "1", Title: "Abbey Road", Artist: "The Beatles", Price: 56.99},
	{Code: "2", Title: "What's Going On", Artist: "Marvin Gaye", Price: 17.99},
	{Code: "3", Title: "Purple Rain", Artist: "Prince", Price: 39.99},
	// {Code: "4", Title: "Bohemian Rhapsod", Artist: "Queen", Price: 12.99},
}

func TestCreateTableAndData(t *testing.T) {
	albumsReturned := CreateTableAndData()
	compareAlbums(albumsReturned, albumsReference, t)
}

func TestGetAlbumByCode(t *testing.T) {
	code := "1"
	albumsReturned := GetAlbumByCode(code)
	compareAlbums(albumsReturned, albumsReference[:1], t)
}

func TestGetAlbums(t *testing.T) {
	albumsReturned := GetAlbums()
	compareAlbums(albumsReturned, albumsReference, t)
}

func TestAddNewAlbum(t *testing.T) {
	newItem := Album{Code: "4", Title: "Bohemian Rhapsod", Artist: "Queen", Price: 12.99}
	albumsAdded := AddNewAlbum(newItem)

	code := "4"
	albumsReturned := GetAlbumByCode(code)
	compareAlbums(albumsReturned, albumsAdded, t)

	newAlbumsReference := append(albumsReference, newItem)
	newAlbumsReturned := GetAlbums()
	compareAlbums(newAlbumsReturned, newAlbumsReference, t)

}

func TestDeleteAlbum(t *testing.T) {
	code := "4"
	albumsReturned := GetAlbumByCode(code)
	albumsDeleted := DeleteAlbum(code)
	compareAlbums(albumsDeleted, albumsReturned, t)

	albumsReturned = GetAlbums()
	compareAlbums(albumsReturned, albumsReference, t)
}

func compareAlbums(albumsReturned []Album, albumsReference []Album, t *testing.T) {

	t.Logf("[Log - Number of Album] received:%d, wanted:%d", len(albumsReturned), len(albumsReference))

	if len(albumsReference) != len(albumsReturned) {
		t.Errorf("[Error - Number of Album] got:%d, wanted:%d", len(albumsReturned), len(albumsReference))
	}

	for i, v := range albumsReturned {

		if v.Code != albumsReference[i].Code {
			t.Errorf("[Error - Album Code in raw] got:%s, wanted:%s", v.Code, albumsReference[i].Code)
		}

		if v.Title != albumsReference[i].Title {
			t.Errorf("[Error - Album Title in raw] got:%s, wanted:%s", v.Title, albumsReference[i].Title)
		}

		if v.Artist != albumsReference[i].Artist {
			t.Errorf("[Error - Album Artist in raw] got:%s, wanted:%s", v.Artist, albumsReference[i].Artist)
		}

		if v.Price != albumsReference[i].Price {
			t.Errorf("[Error - Album Price in raw] got:%v, wanted:%v", v.Price, albumsReference[i].Price)
		}

	}
}
