package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Token struct {
	Token string
}

type SeriesSearchResults struct {
	Data []Series
}

type Series struct {
	FirstAired string
	Id         int
	Network    string
	Overview   string
	SeriesName string
	Status     string
}

type SeriesEpisodesQuery struct {
	Data   []Episode
	Errors string
}

type Episode struct {
	AbsoluteNumber     int
	AiredEpisodeNumber int
	AiredSeason        int
	Directors          []string
	EpisodeName        string
	FirstAired         string
	Overview           string
	Writers            []string
}

const (
	ApiUrl = "https://api.thetvdb.com"
	ApiKey = ""
)

var client http.Client
var token string

func init() {
	client = http.Client{}

	req, err := http.NewRequest(http.MethodPost, ApiUrl+"/login", strings.NewReader("{\"apikey\":\""+ApiKey+"\"}"))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	body, err := do(req)
	if err != nil {
		log.Fatal(err)
	}
	var t Token
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Fatal(err)
	}
	token = t.Token
}

func Search(name string) ([]Series, error) {
	params := url.Values{}
	params.Add("name", name)
	body, err := request("/search/series?" + params.Encode())
	if err != nil {
		return []Series{}, err
	}
	var res SeriesSearchResults
	err = json.Unmarshal(body, &res)
	if err != nil {
		return []Series{}, err
	}
	return res.Data, nil
}

func GetEpisode(series Series, season, episode int) (Episode, error) {
	params := url.Values{}
	params.Add("airedSeason", fmt.Sprintf("%d", season))
	params.Add("airedEpisode", fmt.Sprintf("%d", episode))
	endpoint := fmt.Sprintf("/series/%d/episodes/query?%s", series.Id, params.Encode())
	body, err := request(endpoint)
	if err != nil {
		return Episode{}, err
	}
	var res SeriesEpisodesQuery
	err = json.Unmarshal(body, &res)
	if err != nil {
		return Episode{}, err
	}
	if res.Errors != "" {
		fmt.Println("don't")
		return Episode{}, errors.New(res.Errors)
	}
	if len(res.Data) < 1 {
		return Episode{}, errors.New("No episode found for this season and episode number")
	} else if len(res.Data) > 1 {
		return Episode{}, errors.New("More than one episodes found for this season and episode number")
	}
	return res.Data[0], nil
}

func request(endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, ApiUrl+endpoint, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	body, err := do(req)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func do(req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
