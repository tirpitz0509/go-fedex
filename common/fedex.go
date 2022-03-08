package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	FEDEX_API_URL      = "https://wsbeta.fedex.com:443/web-services"
	FEDEX_API_TEST_URL = "https://ws.fedex.com:443/web-services"
)

type Fedex struct {
	Key, Password, Account, Meter string
	FedexUrl                      string
	TestMode                      bool
}

func (c Fedex) PostRequest(xml string, path string) (content []byte, err error) {
	var url string
	if c.TestMode {
		url = FEDEX_API_TEST_URL + path
	} else {
		url = FEDEX_API_URL + path
	}
	log.Printf("%s", xml)
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	return ioutil.ReadAll(resp.Body)
}
