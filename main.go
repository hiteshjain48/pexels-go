package main

import (
	"math/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

type CuratedPhotos struct {
	Page 		int 	`json:"page"`
	PerPage 	int 	`json:"per_page"`
	NextPage 	string 	`json:"next_page"`
	Photos 		[]Photo `json:"photos"`
}

type VideoSearchResult struct {
	Page 			int 	`json:"page"`
	PerPage 		int 	`json:"per_page"`
	TotalResults 	int 	`json:"total_results"`
	NextPage 		int 	`json:"next_page"`
	Videos 			[]Video `json:"videos"`
}

type Video struct {
	Id 				int 			`json:"id"`
	Width 			int 			`json:"width"`
	Height 			int 			`json:"Height"`
	Url 			string 			`json:"url"`
	Image  			string 			`json:"image"`
	FullRes 		interface{}		`json:"full_res"`
	Duration 		float64 		`json:"duration"`
	VideoFiles 		[]VideoFiles 	`json:"video_files"`
	VideoPictures 	[]VideoPictures `json:"video_pictures"`
}

type PopularVideos struct {
	Page 			int 	`json:"page"`
	PerPage 		int 	`json:"per_page"`
	TotalResults 	int 	`json:"total_results"`
	Url 			string 	`json:"url"`
	Videos 			[]Video `json:"videos"`
}

type VideoFiles struct {
	Id 			int 	`json:"id"`
	Quality 	string 	`json:"quality"`
	FileType 	string 	`json:"file_type"`
	Width 		int 	`json:"width"`
	Height 		int 	`json:"height"`
	Link 		string 	`json:"link:"`
}

type VideoPictures struct {
	Id 		int 	`json:"id"`
	Picture string 	`json:"picture"`
	Nr 		int 	`json:"nr"`
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

func (client *Client) CuratedPhotos(perPage, page int) (*CuratedPhotos, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d",perPage, page)
	res, err := client.requestDoWithAuth("GET",url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result CuratedPhotos
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (client *Client) GetPhoto(id int) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d",id)
	res, err := client.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result Photo
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (client *Client) GetRandomPhoto()(*Photo, error) {
	randNum := rand.Intn(1001)
	result, err := client.CuratedPhotos(1,randNum)
	if err == nil && len(result.Photos)==1 {
		return &result.Photos[0], nil
	}
	return nil , err
}

func (client *Client) SearchVideo(search string, perPage, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/query?=%sper_page=%d&page=%d",search, perPage, page)
	res, err := client.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result VideoSearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (client *Client) PopularVideo(perPage, page int) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d",perPage, page)
	res, err := client.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result PopularVideos
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (client *Client) GetRandomVideo()(*Video, error) {
	randNum := rand.Intn(1001)
	result, err := client.PopularVideo(1,randNum)
	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}
	return nil, err
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

func (client *Client) GetRemainingRequestsThisMonth() int {
	return client.RemainingTimes
}

func main() {
	os.Setenv("PexelsToken", "pUK6JmqnsLWLulPDdjOOwU8nrz0efX1qV9twL5JQIEaFZJCB09Kcwu3y")
	Token := os.Getenv("PexelsToken")
	
	client := NewClient(Token)

	result, err := client.SearchPhotos("planet",15,1)
	if err != nil {
		fmt.Printf("search error: %v", err)
	}

	if result.Page == 0 {
		fmt.Printf("search result wrong")
	}

	fmt.Println(result)
}