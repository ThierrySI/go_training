package model

// Album represents data about a record album.
type Album struct {
	Code   string  `json:"code"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func SampleAlbums() []Album {
	// albums slice to seed record Album data.
	var albums = []Album{
		{Code: "1", Title: "Abbey Road", Artist: "The Beatles", Price: 56.99},
		{Code: "2", Title: "What's Going On", Artist: "Marvin Gaye", Price: 17.99},
		{Code: "3", Title: "Purple Rain", Artist: "Prince", Price: 39.99},
	}
	return albums
}
