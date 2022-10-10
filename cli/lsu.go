package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dollarkillerx/analysis_japanese_schools/internal/server"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	lsu := server.LanguageSchoolUpdate{}
	update, err := lsu.Update()
	if err != nil {
		log.Fatalln(err)
	}

	marshal, err := json.Marshal(update)
	if err != nil {
		log.Fatalln(err)
	}

	os.WriteFile("lsu.json", marshal, 00666)
}
