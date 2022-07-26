package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// handler to return all items
// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// handler to add a new item
// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {

	var newAlbum album

	err := c.BindJSON(&newAlbum)
	if err != nil {
		return
	}

	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// handler to return a specific item
// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {

	id := c.Param("id")

	// loop range by looking for ID
	for _, v := range albums {
		// fmt.Printf("[%s] -> Artiste=%s, Title=%s, Price=%s \n", v.ID, v.Artist, v.Title, v.Price)
		if v.ID == id {
			c.IndentedJSON(http.StatusOK, v)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// handle to delete item
func deleteAlbumByID(c *gin.Context) {
	var isDeleted = false

	id := c.Param("id")

	// loop range by looking for ID
	for idx, v := range albums {
		// fmt.Printf("[%s] -> Artiste=%s, Title=%s, Price=%s \n", v.ID, v.Artist, v.Title, v.Price)
		if v.ID == id {
			albums = append(albums[0:idx], albums[idx+1:]...)
			isDeleted = true
		}
	}

	if isDeleted {
		c.IndentedJSON(http.StatusOK, albums)
	} else {
		errorMessage := fmt.Sprintf("album [%s] not found", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": errorMessage})
	}
}

// handle to update item (=delete + insert)
func updateAlbumByID(c *gin.Context) {
	var isDeleted = false
	albumDeleted := album{}
	id := c.Param("id")

	// loop range by looking for ID
	for idx, v := range albums {
		// fmt.Printf("[%s] -> Artiste=%s, Title=%s, Price=%s \n", v.ID, v.Artist, v.Title, v.Price)
		if v.ID == id {
			albums = append(albums[0:idx], albums[idx+1:]...)
			isDeleted = true
			albumDeleted = v
			//			fmt.Printf("[Album Deleted %s] -> Artiste=%s, Title=%s, Price=%f \n", albumDeleted.ID, albumDeleted.Artist, albumDeleted.Title, albumDeleted.Price)
			//			break
		}
	}

	if isDeleted {
		newAlbum := album{}
		err := c.BindJSON(&newAlbum)
		if err != nil {
			return
		}
		//fmt.Printf("New Album [BEFORE] = %v | %d | %d | %v\n", newAlbum, len(newAlbum.Title), len(newAlbum.Title), newAlbum.Price)

		newAlbum.ID = id
		if len(newAlbum.Title) == 0 {
			newAlbum.Title = albumDeleted.Title
		}
		if newAlbum.Artist == "" {
			newAlbum.Artist = albumDeleted.Artist
		}
		if newAlbum.Price == 0 {
			newAlbum.Price = albumDeleted.Price
		}
		//fmt.Printf("New Album [AFTER] = %v \n", newAlbum)

		albums = append(albums, newAlbum)
		// c.IndentedJSON(http.StatusCreated, newAlbum)
		c.IndentedJSON(http.StatusOK, albums)
	} else {
		errorMessage := fmt.Sprintf("album [%s] not found", id)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": errorMessage})
	}

}

func main() {
	// gin.SetMode(gin.ReleaseMode)

	// Creates a gin router
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.DELETE("/albums/:id", deleteAlbumByID)
	router.PUT("/albums/:id", updateAlbumByID)

	// Attached the router to httpServer and start it on :8080
	router.Run("localhost:8080")
}
