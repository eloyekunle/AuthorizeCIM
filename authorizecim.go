package AuthorizeCIM

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var api_endpoint string = "https://apitest.authorize.net/xml/v1/request.api"
var apiName *string
var apiKey *string
var testMode string
var showLogs bool = true
var connected bool = false

func SetAPIInfo(name string, key string, mode string) {
	apiKey = &key
	apiName = &name
	if mode == "live" {
		showLogs = false
		testMode = "liveMode"
		api_endpoint = "https://api.authorize.net/xml/v1/request.api"
	} else {
		showLogs = true
		testMode = "testMode"
		api_endpoint = "https://apitest.authorize.net/xml/v1/request.api"
	}
}

func GetAuthentication() MerchantAuthentication {
	auth := MerchantAuthentication{
		Name:           apiName,
		TransactionKey: apiKey,
	}
	return auth
}

func SendRequest(input []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", api_endpoint, bytes.NewBuffer(input))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	if showLogs {
		fmt.Println(string(body))
	}
	return body, err
}
