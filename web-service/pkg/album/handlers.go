package album

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_training/web-service/pkg/stream"
	"net/http"
)

// handler to return all items - responds with the list of all albums as JSON.
func GetAlbumsHandler(c *gin.Context) {

	albumsReturned := GetAlbums()

	if albumsReturned == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, albumsReturned)
	}

}

// handler to return a specific item
// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func GetAlbumByCodeHandler(c *gin.Context) {

	code := c.Param("code")
	albumReturned := GetAlbumByCode(code)

	if albumReturned == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		c.IndentedJSON(http.StatusOK, albumReturned)
	}

}

// handler to add a new item - adds an album from JSON received in the request body.
func AddAlbumHandler(c *gin.Context) {

	var newAlbum Album
	err := c.BindJSON(&newAlbum)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	albumResult := AddNewAlbum(newAlbum)
	if albumResult == nil {
		c.AbortWithStatus(http.StatusNotModified)
	} else {
		stream.Producer(newAlbum)
		c.IndentedJSON(http.StatusCreated, albumResult)
	}
}

// handle to delete item
func DeleteAlbumByCodeHandler(c *gin.Context) {

	code := c.Param("code")
	albumDeleted := DeleteAlbum(code)

	if albumDeleted == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		c.IndentedJSON(http.StatusOK, albumDeleted)
	}

}

// handle to update item (=delete + insert)
func ModifyAlbumByCodeHandler(c *gin.Context) {

	code := c.Param("code")
	albumDeleted := DeleteAlbum(code)

	if albumDeleted == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		// deleted - ready to add new album
		var newAlbum Album
		newAlbum.Code = code
		err := c.BindJSON(&newAlbum)

		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			albumResult := AddNewAlbum(newAlbum)
			if albumResult == nil {
				c.AbortWithStatus(http.StatusNotModified)
			} else {
				c.IndentedJSON(http.StatusCreated, albumResult)
			}
		}
	}
}

func InitDBHandler(c *gin.Context) {

	albumsCreated := CreateTableAndData()

	if albumsCreated == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, albumsCreated)
	}

}
