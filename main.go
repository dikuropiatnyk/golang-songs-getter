package main

import (
	"encoding/json"
	"fmt"
	"github.com/dikuropiatnyk/golang-songs-getter/api"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func songGetter() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	rootUrl := os.Getenv("music_service_url")
	accessToken := "access_token=" + os.Getenv("access_token")

	resEndpoint := rootUrl + "track/1393481682?" + accessToken

	resp, err := http.Get(resEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	var data map[string]interface{}
	//Convert the body to JSON and get preview
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("The response:", data["preview"])
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/up", api.Up).Methods("GET")
	router.HandleFunc("/get_songs", api.GetSongs).Methods("GET")
	err := http.ListenAndServe(":8000", router)

	if err != nil {
		fmt.Print(err)
	}

}
