package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"os"
	"io/ioutil"

	"github.com/leoferlopes/desafio-stone/router"
	"github.com/leoferlopes/desafio-stone/config"
)

func main() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var config *config.Config = &config.Settings

	err := json.Unmarshal(file, config)
	if err != nil{
		fmt.Printf("File sintax error: %v\n", e)
		os.Exit(1)
	}

	router := router.NewRouter()

	port := ":" + config.Port
	fmt.Printf("listening port: %s\n", config.Port)

	log.Fatal(http.ListenAndServe(port, router))
}