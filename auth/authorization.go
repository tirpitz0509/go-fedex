package auth

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

const (
	apiURLTest  = "https://apis-sandbox.fedex.com"
	apiURLLive  = "https://apis.fedex.com"
	apiTokenUrl = "https://developer.fedex.com/api/en-ae/catalog/authorization/v1"
)

type FedEXAuth struct {
	GrantType    string `json:"grant_type"`    //
	ClientId     string `json:"client_id"`     //
	ClientSecret string `json:"client_secret"` //
	TestMode     bool   `json:"-"`             //
}

type FedexAuthResponse struct {
	AccessToken string `json:"access_token,omitempty"` //
	TokenType   string `json:"token_type,omitempty"`   //
	ExpiresIn   int    `json:"expires_in,omitempty"`   //
	Scope       string `json:"scope,omitempty"`        //
	Url         string `json:"url,omitempty"`          //
	Errors      []struct {
		Code    string `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"errors,omitempty"`
}

func (c FedEXAuth) Authorization() (*FedexAuthResponse, error) {
	var _response FedexAuthResponse
	var reqUrl string
	client := &http.Client{}

	if c.TestMode {
		reqUrl = apiURLTest + "/oauth/token"
	} else {
		reqUrl = apiURLLive + "/oauth/token"
	}

	request := strings.NewReader(
		`grant_type=` + c.GrantType + `&client_id=` + c.ClientId + `&client_secret=` + c.ClientSecret,
	)

	log.Printf("%s", request)

	req, err := http.NewRequest("POST", reqUrl, request)
	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return &_response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		data := json.NewDecoder(resp.Body)
		errjson := data.Decode(&_response)
		if errjson != nil {
			return &_response, errjson
		}
		if c.TestMode {
			_response.Url = apiURLTest
		} else {
			_response.Url = apiURLLive
		}
		return &_response, nil
	} else {
		dataError := json.NewDecoder(resp.Body)
		dataError.Decode(&_response)
		return &_response, errors.New(_response.Errors[0].Message)
	}

	return &_response, nil

}
