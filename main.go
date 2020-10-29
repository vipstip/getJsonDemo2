package main

import (
	"fmt"
	"getJsonDemo2/GetData"
	"log"
	"net/http"
)

func jsonArr(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/json" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, string(GetData.GetJson()))
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		address := r.FormValue("address")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/json", jsonArr)
	fmt.Printf("Starting server for testing HTTP POST 8000...\n")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
