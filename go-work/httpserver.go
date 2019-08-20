// wiki.go
package main

import (
	// "encoding/json"
	"bookcache"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book_id := vars["book_id"]

	fmt.Printf("Request GET book_id=%v\n", book_id)

	msg := bookcache.NewMessage(bookcache.READ)
	msg.Key = book_id
	resultMsg := bookcache.SendMessage(msg) //block
	contents := resultMsg.Value

	fmt.Printf("Response GET book_id=%v ,contents=%v\n", book_id, contents[:10])

	if resultMsg.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found "+book_id)
		return
	}

	fmt.Fprint(w, contents)

}
func updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	book_id := vars["book_id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	contents := string(reqBody)
	fmt.Printf("Request PUT book_id=%v\n", book_id)

	msg := bookcache.NewMessage(bookcache.WRITE)
	msg.Key = book_id
	msg.Value = contents
	resultMsg := bookcache.SendMessage(msg) //block untill return
	if resultMsg.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Response PUT book_id=%v ,contents=%v\n", book_id, contents[:10])

	fmt.Fprint(w, "Updated Successful!")

	// decoder := json.NewDecoder(r.Body)
	// var data myData
	// err := decoder.Decode(&data)
	// if err != nil {
	// 	panic(err)
	// }
	// owner := data.Owner
	// name := data.Name
}

func main() {
	bookcache.Init(10 * 10000) //size of cached items
	bookcache.StartLoop()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/books/{book_id}", getBook).Methods("GET")
	router.HandleFunc("/api/v1/books/{book_id}", updateBook).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}
