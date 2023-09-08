package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/tools/go/analysis/passes/defers"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token string
	httpClient http.Client
	RemainingTimes int
}

func NewClient(token string) *Client {
	client := http.Client{}
	return &Client{
		Token:token,
		httpClient: client,
	}
}

type SearchResult struct {
	Page int `json:"page"`
	PerPage int `json:"per_page"`
	TotalResults int `json:"total_Results"`
	NextPage string `json:"next_page"`
	Photos []Photo `json:"photos"`
}

type Photo struct {
	Id		int			`json:"id"`
	Width	int			`json:"width"`
	Height	int			`json:"height"`
	Src		PhotoSource	`json:"src"`
}

type PhotoSource struct {
	Original 	string	`json:"original"`
	Large 		string	`json:"large"`
	Large2x 	string	`json:"large2x"`
	Medium 		string	`json:"medium"`
	Small 		string	`json:"small"`
	Portrait 	string	`json:"portrait"`
	Square 		string	`json:"square"`
	Landscape 	string	`json:"landscape"`
	Tiny 		string	`json:"tiny"`
}

func (client *Client) SearchPhotos(search string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d",search,perPage,page)
	res, err := client.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult

	err = json.Unmarshal(data, &result)
	return &result, err
}

func (client *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authentication", client.Token)
	res, err := client.httpClient.Do(req)
	if err != nil {
		return res, err
	}

	times, err := strconv.Atoi(res.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return res, nil
	} else {
		client.RemainingTimes = int(times)
	}

	return res, nil
}

func main() {
	os.Setenv("PexelsToken", "pUK6JmqnsLWLulPDdjOOwU8nrz0efX1qV9twL5JQIEaFZJCB09Kcwu3y")
	Token := os.Getenv("PexelsToken")
	
	client := NewClient(Token)

	result, err := client.SearchPhotos("waves",15,1)
	if err != nil {
		fmt.Errorf("search error: %v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("search result wrong")
	}

	fmt.Println(result)
}