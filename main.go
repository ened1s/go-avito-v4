package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var userURL = "/user"
	var segmURL = "/segm"

	router := mux.NewRouter()
	router.HandleFunc(segmURL, CreateSegment).Methods("POST")
	//пример запроса через POSTMAN: POST localhost:1234/segm?slug=add1

	router.HandleFunc(segmURL, DeleteSegment).Methods("DELETE")
	//DELETE localhost:1234/segm?slug=add1

	router.HandleFunc(userURL, ChangeUserSegments).Methods("PATCH")
	//PATCH localhost:1234/user?add_seg=add1,add2&del_seg=del1,del2&user_id=123

	router.HandleFunc(userURL, GetUserSegments).Methods("GET")
	//GET localhost:1234/user?user_id=123

	fmt.Println("Server at 1234")
	log.Fatal(http.ListenAndServe("localhost:1234", router))
}
