package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/xanzy/go-gitlab"
)

var baseAPI *url.URL

var gitlabClient *gitlab.Client
var singleClient *bool

func getClient() *gitlab.Client {
	var err error
	if singleClient == nil {
		gitlabClient, err = gitlab.NewClient(getToken())
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		b := true
		singleClient = &b
	}
	return gitlabClient
}

var singleToken *bool
var apiToken string

func getToken() string {
	if singleToken == nil {
		var isSet bool
		apiToken, isSet = os.LookupEnv("GITLAB_API_PRIVATE_TOKEN")
		if apiToken == "" || !isSet {
			panic("[ERROR] Must set $GITLAB_API_PRIVATE_TOKEN!")
		}
		b := true
		singleToken = &b
	}
	return apiToken
}

func newRequest(api string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", baseAPI, api)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Set("Accept", apiToken)
	headers.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header[k] = v
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil

	//	err = json.NewDecoder(resp.Body).Decode(&i)
	//	return i, err
}

func main() {
	//	user := getUser()
	//	file, _ := json.Marshal(user)
	//	ioutil.WriteFile("boo.json", file, 0644)
	var err error
	if baseAPI, err = url.Parse("https://gitlab.com/api/v4/"); err != nil {
		log.Fatalf("Failed to parse the given REST API URL: %v", err)
	}
	getProjects()
}
