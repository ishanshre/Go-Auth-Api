package v1

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/user/auth/sign-up", makeHttpHandler(s.handleUserSignup))
	router.HandleFunc("/api/v1/user/auth/login", makeHttpHandler(s.handleUserLogin))
	router.HandleFunc("/api/v1/user", jwtAuthAdminHandler(makeHttpHandler(s.handleGetUser), s.store))
	router.HandleFunc("/api/v1/user/{id}", jwtAuthHandler(makeHttpHandler(s.handleUsersById), s.store))
	log.Println("Starting server at ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

type ApiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type ApiSuccess struct {
	Success string `json:"success"`
}
