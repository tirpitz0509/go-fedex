package common

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	FEDEX_API_TEST_URL = "https://wsbeta.fedex.com:443/web-services"
	FEDEX_API_URL      = "https://ws.fedex.com:443/web-services"
	XSI                = "http://www.w3.org/2001/XMLSchema-instance"
	XSD                = "http://www.w3.org/2001/XMLSchema"
	ENV                = "http://schemas.xmlsoap.org/soap/envelope/"
	ENC                = "http://schemas.xmlsoap.org/soap/encoding/"
)

type Fedex struct {
	TestMode bool //
}

func (c Fedex) PostRequest(xml string, path string) (content []byte, err error) {
	var url string
	if c.TestMode {
		url = FEDEX_API_TEST_URL + path
	} else {
		url = FEDEX_API_URL + path
	}
	xml = `<?xml version="1.0" encoding="UTF-8"?>` + xml
	//log.Println(xml)
	resp, err := http.Post(url, "text/xml", strings.NewReader(xml))
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()

	//log.Println(resp.StatusCode)
	return ioutil.ReadAll(resp.Body)
}
