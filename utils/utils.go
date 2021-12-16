package utils

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func EndpointGetter(endpoint string) (fullUrl string) {
	e := godotenv.Load()
	if e != nil {
		log.Panicln(e)
	}
	rootUrl, urlExist := os.LookupEnv("music_service_url")
	token, tokenExist := os.LookupEnv("access_token")

	if !urlExist || !tokenExist {
		log.Panicln(".env file doesn't include needed variables!")
	}

	accessToken := "access_token=" + token
	fullUrl = rootUrl + endpoint + accessToken
	return
}

func DownloadSong(songLink string, songChannel chan string) {
	songChannel <- songLink + "|||Danichka"
}
