package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//fmt.Println(albums)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/create/{name}", CreateAlbum).Methods("POST")
	router.HandleFunc("/delete/{name}", DeleteAlbum).Methods("DELETE")
	router.HandleFunc("/get/{name}", Listimagesinalbum).Methods("GET")
	router.HandleFunc("/getall", GetallAlbum).Methods("GET")
	router.HandleFunc("/uploadfile", AddImage).Methods("POST")
	router.HandleFunc("/delimage/{name,id}", Deleteimage).Methods("DELETE")
	router.HandleFunc("/listall", Listallimages).Methods("GET")
	router.PathPrefix("/upload").Handler(http.StripPrefix("/upload", http.FileServer(http.Dir("./html/")))).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
	
}
