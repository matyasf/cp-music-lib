package controller

import (
	"fmt"
	box_lib "github.com/annakertesz/cp-music-lib/box-lib"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Server struct {
	db           *sqlx.DB
	Token        string
	musicFolder  int
	coverFolder  int
	clientID     string
	clientSecret string
	privateKey   string
}

type e map[string]string

func NewServer(db *sqlx.DB, clientID string, clientSecret string, privateKey string, musicFolder, coverFolder int) *Server {
	token := box_lib.AuthOfBox(clientID, clientSecret, privateKey)
	return &Server{db: db, Token: token, musicFolder: musicFolder, coverFolder: coverFolder, clientID: clientID, clientSecret: clientSecret, privateKey: privateKey}
}

func (s *Server) GetBoxToken() {
	s.Token = box_lib.AuthOfBox(s.clientID, s.clientSecret, s.privateKey)
}

func (server *Server) Routes() chi.Router {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getIndex")
		fmt.Fprint(w, "cp")
	})

	//Songs
	r.Get("/song", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getAllSong")
		getAllSongs(server.db, w, r)
	})
	r.Get("/song/{songID}", func(w http.ResponseWriter, r *http.Request) {
		getSongByID(server.db, w, r)
	})
	r.Get("/song/findByAlbum", func(w http.ResponseWriter, r *http.Request) {
		getSongByAlbum(server.db, w, r)
	})
	r.Get("/song/findByArtist", func(w http.ResponseWriter, r *http.Request) {
		getSongByArtist(server.db, w, r)
	})
	r.Get("/song/findByTag", func(w http.ResponseWriter, r *http.Request) {
		getSongByTag(server.db, w, r)
	})
	//TODO
	r.Get("/song/findByFreeSearch", func(w http.ResponseWriter, r *http.Request) {
		searchSong(server.db, w, r)
	})

	//Albums
	r.Get("/album", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("album request")
		getAllAlbum(server.db, w, r)

	})
	r.Get("/album/findByArtist", func(w http.ResponseWriter, r *http.Request) {
		getAlbumsByArtist(server.db, w, r)
	})
	r.Get("/album/{albumID}", func(w http.ResponseWriter, r *http.Request) {
		getAlbumsById(server.db, w, r)
	})

	//Artists
	r.Get("/artist", func(w http.ResponseWriter, r *http.Request) {
		getAllArtist(server.db, w, r)
	})
	r.Get("/artist/{artistID}", func(w http.ResponseWriter, r *http.Request) {
		getArtistById(server.db, w, r)
	})

	//Tags
	r.Get("/tag", func(w http.ResponseWriter, r *http.Request) {
		getAllTag(server.db, w, r)
	})

	r.Get("/update", func(w http.ResponseWriter, r *http.Request) {
		update(server.db, w, r, server.Token, server.coverFolder, server.musicFolder)
	})

	r.Get("/download/{boxID}", func(w http.ResponseWriter, r *http.Request) {
		err := download(server.db, server.Token, w, r)
		if err != nil {
			server.GetBoxToken()
			err := download(server.db, server.Token, w, r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})

	//Playlist
	//TODO
	r.Post("/playlist", func(w http.ResponseWriter, r *http.Request) {
		createPlaylist(server.db, w, r)
	})
	//TODO
	r.Get("/playlist/{playlistID}", func(w http.ResponseWriter, r *http.Request) {
		getPlaylistById(server.db, w, r)
	})

	//User
	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		createUser(server.db, w, r)
	})
	r.Options("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Accept, Content-Length, Accept-Encoding, X-CSRF-TOKEN, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/user/{userID}/validate", func(w http.ResponseWriter, r *http.Request) {
		validateUser(server.db, w, r)
	})
	r.Post("/user/login", func(w http.ResponseWriter, r *http.Request) {
		loginUser(server.db, w, r)
	})
	return r
}
