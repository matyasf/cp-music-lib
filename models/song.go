package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Song struct {
	SongID int `json:"song_id" db:"id"`
	SongName string `json:"song_name" db:"song_name"`
	SongLqURL string `json:"song_lq_url" db:"song_lq_url"`
	SongHqURL string `json:"song_hq_url" db:"song_hq_url"`
	SongInstrumentalLqURL string `json:"song_instrumental_lq_url" db:"instrumental_lq_url"`
	SongInstrumentalHqURL string `json:"song_instrumental_hq_url" db:"instrumental_hq_url"`
	SongAlbum int `json:"song_album" db:"song_album"`
	SongTags []*Tag `json:"song_tags" db:"song_tag"`
	boxID int
}

func NewSong(name string, album int, boxID int) Song {
	return Song{
		SongName:name,
		SongAlbum:album,
		boxID:boxID,
	}
}

func (song *Song) CreateSong(db *sqlx.DB) (int, bool) {
	isInstrumental := false
	createdNew := false
	var id int
	if strings.Contains(song.SongName, "nstrumental") {
		if strings.Index(song.SongName, "(") > 0 {
		song.SongName = strings.TrimSpace(song.SongName[:strings.Index(song.SongName,"(")])
		isInstrumental = true}
	}
	db.QueryRow(
		`SELECT id from song where song_name = $1`, song.SongName,
	).Scan(&id)
	if id != 0 {
		if isInstrumental {
			db.Query(`UPDATE song SET instrumental_lq_url = $1 WHERE id=$2`, song.boxID, id)
			fmt.Printf("\nfound %v and inserted instrumental version", song.SongName)
		} else {
			db.Query(`UPDATE song SET song_lq_url = $1 WHERE id=$2`, song.boxID, id)
			fmt.Printf("\nfound %v and inserted vocal version", song.SongName)

		}
	} else {
		if isInstrumental {
			db.QueryRow(
				`INSERT INTO song (song_name, song_album,instrumental_lq_url) VALUES ($1, $2, $3) RETURNING id`,
				song.SongName,  song.SongAlbum, song.boxID,
			).Scan(&id)
		} else {
			db.QueryRow(
				`INSERT INTO song (song_name, song_album,song_lq_url) VALUES ($1, $2, $3) RETURNING id`,
				song.SongName,  song.SongAlbum, song.boxID,
			).Scan(&id)
		}
		createdNew=true
	}

	return id, createdNew
}

func GetSongByID(id int, db *sqlx.DB) (*Song, error) {
	var song Song
	err := db.QueryRowx(`SELECT * FROM song WHERE id = $q`, id,
	).StructScan(&song)
	if err!=nil{
		return nil, err
	}
	return &song, nil
}

func GetSongByArtist(id int, db *sqlx.DB) ([]Song, error) {
	rows, err := db.Queryx(
		`SELECT * from song join album on song.song_album = album.id where album_artist = $1` , id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := make([]Song, 0)
	for rows.Next() {
		var song Song
		rows.StructScan(&song)
		songs = append(songs, song)
	}
	return songs, nil
}

func GetSongByAlbum(id int, db *sqlx.DB) ([]Song, error) {
	rows, err := db.Queryx(
		`SELECT * FROM album WHERE song_album = $1` , id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := make([]Song, 0)
	for rows.Next() {
		var song Song
		rows.StructScan(&song)
		songs = append(songs, song)
	}
	return songs, nil
}

func GetSongByTag(id int, db *sqlx.DB) ([]Song, error) {
	rows, err := db.Queryx(
		`select * from tag_song join tag on tag.id = tag_song.map_tag join song on song.id = tag_song.map_song where tag.id = $1` , id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	songs := make([]Song, 0)
	for rows.Next() {
		var song Song
		rows.StructScan(&song)
		songs = append(songs, song)
	}
	return songs, nil
}

//func GetSongByTitleAndAlbum(title string, album string, db *sqlx.DB) ([]Song, error) {
//
//	//Get song object
//	rows, err := db.Queryx(
//		`SELECT id, song_name, song_artist, song_album, song_lq_url, song_hq_url, instrumental_lq_url, instrumental_hq_url FROM song WHERE song.song_name = $1` , title,
//	)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//	songs := make([]Song, 0, 5)
//
//	for rows.Next() {
//		var id int
//		var song_name string
//		var song_artist int
//		var song_album int
//		var song_lq_url string
//		var song_hq_url string
//		var instrumental_lq_url string
//		var instrumental_hq_url string
//
//		err := rows.Scan(&id, &song_name, &song_artist, &song_album, &song_hq_url, &song_lq_url, &instrumental_hq_url, &instrumental_lq_url)
//		if err != nil {
//			fmt.Print(err.Error())
//			return nil, err
//		}
//		//artist, err := GetArtistByID(song_artist, db)
//		album, err := GetAlbumByID(song_album, db)
//		tags, err := GetTagsOfSong(id, db)
//		if err != nil {
//			fmt.Print(err.Error())
//			return nil, err
//		}
//
//		songs = append(songs, Song{
//			SongID:                id,
//			SongName:              song_name,
//			SongLqURL:             song_lq_url,
//			SongHqURL:             song_hq_url,
//			SongInstrumentalLqURL: instrumental_lq_url,
//			SongInstrumentalHqURL: instrumental_hq_url,
//			SongAlbum:             &album,
//			SongArtist:            album.AlbumArtist,
//			SongTags:              tags,
//		})
//	}
//
//	return songs, nil
//}