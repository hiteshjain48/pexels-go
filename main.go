package main

import(
	"fmt"
	"net/http"
	"encoding/json"
	"os"
	"strconv"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

func main() {
	os.Setenv("PexelsToken", "")
	Token := os.Getenv("PexelsToken")
	
	client := NewClient(Token)

	result, err := client.SearchPhotos("waves")
	if err != nil {
		fmt.Errorf("search error: %v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("search result wrong")
	}

	fmt.Println(result)
}