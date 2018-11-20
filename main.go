package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Synonym struct {
	Word  string `json:"word"`
	Score int    `json:"score"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		resp, err := http.Get(fmt.Sprintf("https://api.datamuse.com/words?rel_syn=%v", message))

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(message)
		fmt.Println(resp)
		defer resp.Body.Close()

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var synonymArray []Synonym
		err = json.Unmarshal(responseData, &synonymArray)

		var synonyms []string
		for i := 0; i < len(synonymArray); i++ {
			synonyms = append(synonyms, synonymArray[i].Word)
		}

		synonymList := strings.Join(synonyms, ",")
		synonymList = strings.Title(synonymList)
		fmt.Fprint(w, fmt.Sprintf("%v", synonymList))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
