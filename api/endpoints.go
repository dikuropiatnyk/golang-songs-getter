package api

import (
	"encoding/json"
	u "github.com/dikuropiatnyk/golang-songs-getter/utils"
	"log"
	"net/http"
	"os"
)

var Up = func(w http.ResponseWriter, r *http.Request) {
	log.Println("We came in town!")
	name := r.URL.Query().Get("name")

	resp := u.Message(true, "success")
	resp["data"] = "It's alive!" + name
	u.Respond(w, resp, 200)
}

var GetSongs = func(w http.ResponseWriter, r *http.Request) {
	log.Println("We are in the songs getter!")
	//expectedEndpoint := u.DefaultEndpoint
	//if r.URL.Query().Has("endpoint") {
	//	expectedEndpoint = r.URL.Query().Get("endpoint")
	//}
	//fullUrl := u.EndpointGetter(expectedEndpoint)

	//resp, err := http.Get(fullUrl)
	//if err != nil {
	//	log.Println(err)
	//	u.Respond(w, u.Message(false, "Request to the music service wasn't succeeded"), 400)
	//}
	//
	//if resp != nil {
	//	defer resp.Body.Close()
	//}
	//
	//// Create an empty
	//var data map[string]interface{}
	////Convert the body to JSON and get preview
	//err = json.NewDecoder(resp.Body).Decode(&data)
	//if err != nil {
	//	log.Println(err)
	//	u.Respond(w, u.Message(false, "Music service response isn't parseable!"), 400)
	//}

	songLinks := make([]string, 0, u.DefaultSongsCount)
	// Create array of string channels with slice capacity
	var songChannels [u.DefaultSongsCount]chan string
	for i := range songChannels {
		songChannels[i] = make(chan string)
	}

	// Create goroutines to parse tracks list concurrently
	// Open our jsonFile
	jsonFile, err := os.Open("tracks.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
		u.Respond(w, u.Message(false, "Request to the music service wasn't succeeded"), 400)
		return
	}
	log.Println("Successfully Opened tracks.json")
	defer jsonFile.Close()

	// Create an empty interface to read data
	var data u.Tracks
	//Convert the body to JSON and get preview
	err = json.NewDecoder(jsonFile).Decode(&data)
	if err != nil {
		log.Println(err)
		u.Respond(w, u.Message(false, "Music service response isn't parseable!"), 400)
		return
	}

	log.Println(data)

	// Iterate through received tracks
	for _, track := range data.Data {
		if track.Preview == "" {
			log.Printf("Song '%s' doesn't have preview. Skipping...", track.Title)
			continue
		}
		if len(songLinks) < u.DefaultSongsCount {
			songLinks = append(songLinks, track.Preview)
		}
	}

	if len(songLinks) < u.DefaultSongsCount {
		u.Respond(w, u.Message(false, "Music service didn't provide enough songs!"), 400)
		return
	}

	// Finally, count all songs concurrently
	result := make([]string, 0, len(songLinks))
	log.Println("Start concurrent songs download")
	// Run all downloading goroutines
	for songIndex, songLink := range songLinks {
		go u.DownloadSong(songLink, songChannels[songIndex])
	}
	// Gather all results into one slice
	for i := range songChannels {
		result = append(result, <-songChannels[i])
	}
	log.Println("Finish concurrent songs download")
	log.Println(result)

	response := u.Message(true, "All good!")
	response["songs"] = result

	u.Respond(w, response, 200)
}
