package album

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_training/web-service/pkg/database"
	"net/http"
)

type handler struct {
	DB *database.ClickhouseDB
}

func RegisterRoutes(router *gin.Engine, db *database.ClickhouseDB) {
	h := &handler{
		DB: db,
	}

	router.GET("/albums", h.GetAlbums)
	router.GET("/albums/:code", h.GetAlbum)
	router.POST("/albums", h.AddAlbum)
	router.DELETE("/albums/:code", h.RemoveAlbum)
	router.PUT("/albums/:code", h.ModifyAlbum)
	router.GET("/albums/initdb", h.InitDB)

}

// handler to return all items
func (h handler) GetAlbums(c *gin.Context) {
	albumsReturned := h.selectAlbums()

	if albumsReturned == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, albumsReturned)
	}
}

// handler to return a specific item
func (h handler) GetAlbum(c *gin.Context) {

	code := c.Param("code")
	albumReturned := h.selectAlbum(code)

	if albumReturned == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		c.IndentedJSON(http.StatusOK, albumReturned)
	}
}

// handler to add a new item - adds an album from JSON received in the request body.
func (h handler) AddAlbum(c *gin.Context) {

	var newAlbum Album
	err := c.BindJSON(&newAlbum)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	} else {

		albumAdded := h.insertNewAlbum(newAlbum)
		if albumAdded == nil {
			c.AbortWithStatus(http.StatusNotModified)
		} else {
			// stream.Producer(newAlbum)
			c.IndentedJSON(http.StatusCreated, albumAdded)
		}
	}
}

// handle to delete item
func (h handler) RemoveAlbum(c *gin.Context) {

	code := c.Param("code")
	albumDeleted := h.deleteAlbum(code)

	if albumDeleted == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album [%s] not found", code)})
	} else {
		c.IndentedJSON(http.StatusOK, albumDeleted)
	}

}

// handle to update item (=delete + insert)
func (h handler) ModifyAlbum(c *gin.Context) {

	code := c.Param("code")
	albumDeleted := h.deleteAlbum(code)

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
			albumAdded := h.insertNewAlbum(newAlbum)
			if albumAdded == nil {
				c.AbortWithStatus(http.StatusNotModified)
			} else {
				c.IndentedJSON(http.StatusCreated, albumAdded)
			}
		}
	}
}

func (h handler) InitDB(c *gin.Context) {

	albumsCreated := h.createTableAndData()

	if albumsCreated == nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, albumsCreated)
	}

}
