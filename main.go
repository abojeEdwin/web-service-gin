package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	Genre  string  `json:"genre"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99, Genre: "Pop"},
	{ID: "2", Title: "Odumodublavk", Artist: "Gerry Mulligan", Price: 17.99, Genre: "RnB"},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99, Genre: "Country"},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/add", addAlbum)
	router.GET("/get/:id", getAlbumByID)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(context *gin.Context) {

	var newAlbum album
	if err := context.BindJSON(&newAlbum); err != nil || newAlbum.ID == "" {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//if newAlbum.ID == "" || newAlbum.Title == "" || newAlbum.Artist == "" {
	//	context.IndentedJSON(http.StatusNotFound, gin.H{"error": "Album ID and Title/Artist are required"})
	//}
	albums = append(albums, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
