package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"example.com/database"
	"example.com/types"
	"github.com/gin-gonic/gin" //Gin - библиотека для создания серверов: для получения сообщений с серверов и возрата ответов, через функции-обрабочтики программы
)

// album represents data about a record album.
/*type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}*/

// albums slice to seed record album data.
/*var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}*/

func main() {

	fmt.Println(32 << (^uint(0) >> 63))

	database.WorkForDB()

	router := gin.Default()
	router.GET("/albums", getAlbums) //GET Метод добавления endpoint в router.
	router.GET("/albums/:ID", getAlbumById)
	router.PUT("/albums/:ID", putAlbumById) // Нужно написать в DB, вызов альбома по ID
	//router.PUT("album/:Title, putAlbumByTitle") // Нужно напистать в DB, вызов альбома по ID
	//router.PUT("album/:Artist, putAlbumByArtist") // Нужно напистать в DB, вызов альбома по ID
	//router.PUT("album/:Complex, putAlbumComplex") // Нужно напистать в DB, вызов альбома по ID
	router.POST("/albums", postAlbums)
	router.DELETE("/albums/:ID", deleteAlbum)
	router.Run("localhost:8080") //Run - команда запуска router
	//8080 - порт данного компьютера, не служебный. Служебные от 1 до 1023. Можем использовать порты с 1024 до 49000
	//localhost синоним адреса IP-адреса 127.0.0.1 - личный адрес внутренней сети, который не доступен из внешних сетей. 127.x.x.x есть у всех пользователей.
	//В итоге получаем адрес localhost:8080/albums, по которому лежат срезы

}

//var M map [string]func(*gin.Context)

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	allalbums, err := database.GetAllAlbums()
	if err != nil {
		c.String(http.StatusInternalServerError, "{\"Error\": \"%v\"}", err)
		return
		// c.String(http.StatusInternalServerError, err.Error())
	}
	c.IndentedJSON(http.StatusOK, allalbums)
	//f:= M["abc"]
	//f(c)
}

func getAlbumById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("ID"), 10, 64)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	album, err := database.AlbumByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum types.Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Add the new album to the slice.
	if _, err := database.AddAlbum(newAlbum); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func putAlbumById(c *gin.Context) {
	var changealb types.AlbumForChanges

	if err := c.BindJSON(&changealb); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseInt(c.Param("ID"), 10, 64)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = database.PutAlbumById(id, changealb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func deleteAlbum(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("ID"), 10, 64)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteData(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)

}

//curl http://localhost:8080/albums --header "Content-Type: application/json" --request "POST" --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
//curl -d '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}' -H "Content-Type: application/json" -X POST http://localhost:8080/albums
