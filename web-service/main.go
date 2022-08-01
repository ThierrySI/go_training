package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"web-service/model"
)

// handler to return all items - responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {

	albumsReturned := model.GetAlbums()

	if albumsReturned == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, albumsReturned)
	}

}

// handler to return a specific item
// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByCode(c *gin.Context) {

	code := c.Param("code")
	albumReturned := model.GetAlbumByCode(code)

	if albumReturned == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		c.IndentedJSON(http.StatusOK, albumReturned)
	}

}

// handler to add a new item - adds an album from JSON received in the request body.
func addAlbum(c *gin.Context) {

	var newAlbum model.Album
	err := c.BindJSON(&newAlbum)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		albumResult := model.AddNewAlbum(newAlbum)
		if albumResult == nil {
			c.AbortWithStatus(http.StatusNotModified)
		} else {
			c.IndentedJSON(http.StatusCreated, albumResult)
		}
	}

}

// handle to delete item
func deleteAlbumByCode(c *gin.Context) {

	code := c.Param("code")
	albumDeleted := model.DeleteAlbum(code)

	if albumDeleted == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		c.IndentedJSON(http.StatusOK, albumDeleted)
	}

}

// handle to update item (=delete + insert)
func modifyAlbumByCode(c *gin.Context) {

	code := c.Param("code")
	albumDeleted := model.DeleteAlbum(code)

	if albumDeleted == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		// deleted - ready to add new album
		var newAlbum model.Album
		newAlbum.Code = code
		err := c.BindJSON(&newAlbum)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			albumResult := model.AddNewAlbum(newAlbum)
			if albumResult == nil {
				c.AbortWithStatus(http.StatusNotModified)
			} else {
				c.IndentedJSON(http.StatusCreated, albumResult)
			}
		}
	}
}

func initDB(c *gin.Context) {

	albumsCreated := model.CreateTableAndData()

	if albumsCreated == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, albumsCreated)
	}

}

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Creates a gin router
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/albums/:code", getAlbumByCode)
	router.GET("/albums/initdb", initDB)

	router.POST("/albums", addAlbum)
	router.DELETE("/albums/:code", deleteAlbumByCode)
	router.PUT("/albums/:code", modifyAlbumByCode)

	// Attached the router to httpServer and start it on :8080
	router.Run("localhost:8080")
}
