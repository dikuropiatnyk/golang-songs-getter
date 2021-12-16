package utils

// Get songs from radio `Hits`
const DefaultEndpoint = "radio/37151/tracks"
const DefaultSongsCount = 5

type Tracks struct {
	Data []Track `json:"data"`
}

type Track struct {
	Preview string `json:"preview"`
	Title   string `json:"title"`
	Artist  struct {
		Name string `json:"name"`
	} `json:"artist"`
}
