package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"testing"

	"example.com/types"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func WorkForDB() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "root",
		Passwd: "PoLiNa030516",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func TestMain(*testing.T) {

	WorkForDB()

	albums, err := AlbumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := AlbumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := AddAlbum(types.Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist queries for albums that have the specified artist name.
func AlbumsByArtist(name string) ([]types.Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []types.Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb types.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func AlbumByID(id int64) (types.Album, error) {
	// An album to hold data from the returned row.
	var alb types.Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, err
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func AddAlbum(alb types.Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func GetAllAlbums() ([]types.Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []types.Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb types.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %v", err)
	}
	return albums, nil
}

func ParseArgs(id int64, alb types.AlbumForChanges) (string, []any) {

	const_result_1 := "UPDATE album SET "
	const_result_2 := " WHERE ID = ?"
	args := []any{}

	if alb.Title != nil {
		const_result_1 += "title = ?," // Реализовать запись запятых
		args = append(args, alb.Title)
	}

	if alb.Artist != nil {
		const_result_1 += "artist = ?,"
		args = append(args, alb.Artist)
	}

	if alb.Price != nil {
		const_result_1 += "price = ?,"
		args = append(args, alb.Price)
	}

	const_result_1 = strings.TrimSuffix(const_result_1, ",")
	args = append(args, id)

	return const_result_1 + const_result_2, args
}

// необходимо получить id из запроса и внести исправлаение в Список альбомов: []Album
func PutAlbumById(id int64, alb types.AlbumForChanges) error { // обработать случаи с ошибкой, чтобы выводилась ошибка 404 и 500
	//1. Обратиться к конкертному альбому по Last ID
	//2. Изменить LastID на NewID

	query, args := ParseArgs(id, alb)

	_, err := db.Exec(query,
		args...)

	return err
}

func DeleteData(id int64) error {
	_, err := db.Exec("DELETE FROM album WHERE ID = ?", id)

	if err != nil {
		return err
	}
	return nil
}
