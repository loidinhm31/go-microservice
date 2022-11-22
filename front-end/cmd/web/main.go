package main

import (
	"fmt"
	"github.com/loidinhm31/go-micro/common"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	log.Printf("Starting front end service on port %s\n", common.FrontEndPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", common.FrontEndPort), nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {
	wDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if !strings.HasSuffix(wDir, "/cmd/web") {
		err = os.Chdir("./cmd/web/")
		if err != nil {
			log.Println(err)
		}
	}

	partials := []string{
		"./templates/base.layout.gohtml",
		"./templates/header.partial.gohtml",
		"./templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
