package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type image struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type images []image

type album struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	List images `json:"list"`
}

type albums []album

var gallery = albums{}

func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var a album
	a.Name = mux.Vars(r)["name"]
	a.ID = len(gallery)

	for _,b := range gallery{
		if b.Name == a.Name{
			fmt.Fprintf(w, "Album %v already exists", a.Name)
			json.NewEncoder(w).Encode(gallery)
			w.WriteHeader(http.StatusAlreadyReported)
			return
		}
	}

	gallery = append(gallery, a)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(gallery)

	fmt.Fprintf(w, "Album %v Created", a)
}

func AddImage(w http.ResponseWriter, r *http.Request) {
	//name := mux.Vars(r)["name"]
	name := r.FormValue("name")
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



	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	tempFile, err := ioutil.TempFile(fmt.Sprintf("%v/gallery/", path), "image-*"+handler.Filename)
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
	for i, a := range gallery {
		if a.Name == name {
			n := image{
				Name: tempFile.Name(),
				ID:   len(a.List),
			}
			gallery[i].List = append(gallery[i].List, n)
		}
	}

	fmt.Fprintf(w, "Successfully Uploaded Image\n")


}

func Deleteimage(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for _, a := range gallery {
		if a.Name == name {
			for _, b := range a.List {
				if b.ID == id {
					path, err := os.Getwd()
					if err != nil {
						log.Println(err)
					}
					err = os.Remove(path + "/" + b.Name)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Fprintf(w, "Image %v Deleted", id)
				}
			}

		}
	}
}

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	for i, a := range gallery {
		if a.Name == name {
			path, err := os.Getwd()
			if err != nil {
				log.Println(err)
			}
			for _, b := range a.List {
				err = os.Remove(path + "/" + b.Name)
				if err != nil {
					log.Fatal(err)
				}
			}
			gallery = append(gallery[:i], gallery[i+1:]...)
		}
	}
	fmt.Fprintf(w, "Album %v Deleted", name)
}

func GetallAlbum(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(gallery)
}

func Listimagesinalbum(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	for _, a := range gallery {
		if a.Name == name {
			json.NewEncoder(w).Encode(a)
		}
	}
}

func Listallimages(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(gallery)
}
