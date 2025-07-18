package types

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type AlbumForChanges struct {
	Title  *string  `json:"title"`
	Artist *string  `json:"artist"`
	Price  *float32 `json:"price"`
}
