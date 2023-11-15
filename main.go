package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)
type Quote struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
			log.Printf("defaulting to port %s", port)
	}

	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/index.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/quote", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/fragments/quote.html"))
		resp, err := http.Get("https://api.quotable.io/random")
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)

		data := struct {
			Content string `json:"content"`
			Author  string `json:"author"`
		}{}
		json.Unmarshal(body, &data)
		tmpl.Execute(w, data)
	})

	log.Println("Bauhaus running on localhost:"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}