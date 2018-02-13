package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type languageDetail struct {
	Boilerplate   string `json:"boilerplate"`
	CommentPrefix string `json:"commentPrefix"`
}

var languages map[string]languageDetail
var cbAddress string

func main() {
	// Read settings. Compilebox address/port
	port, portOk := os.LookupEnv("COMPILEBOX_PORT")
	address, addressOk := os.LookupEnv("COMPILEBOX_ADDRESS")
	if !portOk || !addressOk {
		log.Fatal("Missing compilebox environment variables, please make sure service is available")
	}
	cbAddress = address + ":" + port

	// Check to ensure the compilebox is up by trying to fill langs variable
	fmt.Printf("Requesting language list from compilebox (%s)...\n", cbAddress)
	populateLanguages()

	/* Serve REST API endpoints for

	- various challenge searches

	- simple code submission (just pass along to testbox)

	- code submission checked against challenge

	*/

	// http.HandleFunc("/get_challenge/", getChallenge)
	// http.HandleFunc("/submit/", submitTest)
	// http.HandleFunc("/stdout/", getStdout)
	// http.HandleFunc("/languages/", getLangs)
	// http.HandleFunc("/", frontPage)

	port = getEnv("TESTBOX_PORT", "31336")
	fmt.Println("testbox listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	/* Serve administrative interface for challenge collection

	 */
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	fmt.Printf("Environment variable %s not found, setting to %s\n", key, fallback)
	os.Setenv(key, fallback)
	return fallback
}

func populateLanguages() {
	r, err := http.Get(cbAddress + "/languages/")

	if err != nil {
		log.Fatal("Unable to contact compilebox, please ensure service is available")
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	// b := make([]byte, 256)
	// _, _ = r.Body.Read(b)
	// fmt.Printf("response: %s", b)

	err = decoder.Decode(&languages)
	if err != nil {
		panic(err)
	}

	supportedLangs := make([]string, 0, len(languages))
	for k := range languages {
		supportedLangs = append(supportedLangs, fmt.Sprintf("%s", k))
	}
	fmt.Printf("Supporting: %s\n", strings.Join(supportedLangs, ", "))
}
