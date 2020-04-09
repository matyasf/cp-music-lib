package controller

import (
	"encoding/json"
	"github.com/annakertesz/cp-music-lib/models"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

func getAllAlbum(db *sqlx.DB, w http.ResponseWriter, r *http.Request){
	albums, err := models.GetAlbum(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	AlbumROs, err := albumROListFromAlbums(albums, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	albumsJSON, err := json.Marshal(AlbumROs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(albumsJSON)
	w.WriteHeader(http.StatusOK)
}


func getAlbumsById(db *sqlx.DB, w http.ResponseWriter, r *http.Request){
	param := chi.URLParam(r, "albumId")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)  //TODO: bad request to swagger
		return
	}
	album, err := models.GetAlbumByID(id, db)
	if album == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	artist, err := models.GetArtistByID(album.AlbumArtist, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	albumJSON, err := json.Marshal(albumROFromAlbum(*album, *artist))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(albumJSON)
	w.WriteHeader(http.StatusOK)
}


func getAlbumsByArtist(db *sqlx.DB, w http.ResponseWriter, r *http.Request){
	param := r.URL.Query().Get("artistID")
	id, err := strconv.Atoi(param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)  //TODO: bad request to swagger
		return
	}
	albums, err := models.GetAlbumByArtist(id, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	AlbumROs, err := albumROListFromAlbums(albums, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	albumsJSON, err := json.Marshal(AlbumROs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(albumsJSON)
	w.WriteHeader(http.StatusOK)
}

func albumROFromAlbum(album models.Album, artist models.Artist) AlbumRO {
	return AlbumRO{
		ID:        album.AlbumID,
		Title:     album.AlbumName,
		Artist:    artistROFromArtist(artist),
	}
}

func albumROListFromAlbums(albums []models.Album, db *sqlx.DB) ([]AlbumRO, error) {
	albumROs := make([]AlbumRO, 0)
	for _, album := range albums {
		artist, err := models.GetArtistByID(album.AlbumArtist, db)
		if err != nil {
			return nil, err
		}
		albumROs = append(albumROs, albumROFromAlbum(album, *artist))
	}
	return albumROs, nil
}