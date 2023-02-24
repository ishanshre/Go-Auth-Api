package v1

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	/*
		A struct that stores the listening address and access to Storage interface
	*/
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	/*
		Initializing variable for starting the http server
	*/
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	/*
		This method is initiates the url controller after database connection and ping is successfull
		We use gorialla mux library to handle the urls
	*/
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/auth/sign-up", makeHttpHandler(s.handleUserSignup))
	router.HandleFunc("/api/v1/auth/login", makeHttpHandler(s.handleUserLogin))
	router.HandleFunc("/api/v1/admin/user/{id}", jwtAuthAdminHandler(makeHttpHandler(s.handleAdminUserById), s.store))
	router.HandleFunc("/api/v1/user", jwtAuthAdminHandler(makeHttpHandler(s.handleGetUser), s.store))
	router.HandleFunc("/api/v1/user/{id}", jwtAuthHandler(makeHttpHandler(s.handleUsersById), s.store))
	log.Println("Starting server at ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error // signature of our handler

type ApiError struct {
	Error string `json:"error"`
}

type ApiSuccess struct {
	Success string `json:"success"`
}
