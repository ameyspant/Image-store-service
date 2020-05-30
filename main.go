package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type album struct{
	Id int `json:"ID"`
	Name string `json:"Name"`
}

type allAlbum []album

var albums = allAlbum{
	{
		Id: 1,
		Name: "photo",
	},
}
var i = 2

func createAlbum(w http.ResponseWriter, r *http.Request) {
	newalbum := album{
		Id: i,
		Name: mux.Vars(r)["name"],
	}
	i = i+1

	albums= append(albums,newalbum)
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(path+"/"+newalbum.Name); os.IsNotExist(err) {
		os.Mkdir(path+"/"+newalbum.Name, os.ModePerm)
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newalbum)

	fmt.Fprintf(w, "Album %v Created", newalbum.Name)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	albumname := mux.Vars(r)["name"]

	for i, singlealbum := range albums{
		if singlealbum.Name == albumname {
			albums = append(albums[:i], albums[i+1:]...)
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			err =os.RemoveAll(path+"/"+albumname+"/")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "Album %v Deleted", albumname)
		}
		fmt.Fprintf(w,"id %v== %v \n", singlealbum.Name, albumname)
	}

}

func getAlbum(w http.ResponseWriter, r *http.Request) {
	albumid := mux.Vars(r)["name"]

	for _, singlealbum := range albums {
		if singlealbum.Name == albumid{
			json.NewEncoder(w).Encode(singlealbum)
		}
	}
	fmt.Fprintf(w, "Album Listed")
}

func getallAlbum(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(albums)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {



	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	newalbum := r.FormValue("name")

	fmt.Println("File Upload Endpoint Hit")

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("%v/%s",path,newalbum)
	tempFile, err := ioutil.TempFile(fmt.Sprintf("%v/%s",path,newalbum), "image-*"+handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func main()  {
	//fmt.Println(albums)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/create/{name}", createAlbum).Methods("POST")
	router.HandleFunc("/delete/{name}", deleteAlbum).Methods("DELETE")
	router.HandleFunc("/get/{name}", getAlbum).Methods("GET")
	router.HandleFunc("/getall", getallAlbum).Methods("GET")
	router.HandleFunc("/uploadfile", uploadFile).Methods("POST")
	router.PathPrefix("/upload").Handler(http.StripPrefix("/upload",http.FileServer(http.Dir("./html/")))).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000",router))
}


