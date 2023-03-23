package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Sequence   string `json:"sequence"`
	DeviceId   string `json:"deviceid"`
	SelfAPIKey string `json:"selfApikey"`
	Data       string `json:"data"`
	Encrypt    bool   `json:"encrypt"`
	IV         string `json:"iv"`
}

type Data struct {
	Switch string `json:"switch"`
}

type Response struct {
	Seq      int    `json:"seq"`
	Sequence string `json:"sequence"`
	Error    int    `json:"error"`
	Info     string `json:"info"`
}

func turnOn(deviceIp string, key string) {
	var data Data
	data.Switch = "on"
	body := generateBody(data, key)
	_, err := sendRequest(deviceIp, body)
	if err != nil {
		log.Print(err)
	}
}

func turnOff(deviceIp string, key string) {
	var data Data
	data.Switch = "off"
	body := generateBody(data, key)
	_, err := sendRequest(deviceIp, body)
	if err != nil {
		log.Print(err)
	}
}

func sendRequest(deviceIp string, body Request) (Response, error) {
	var result Response
	// convert body to bytes
	bodyJSON, _ := json.Marshal(body)

	if config.Debug {
		return result, nil
	}

	// create and send request
	req, err := http.NewRequest("POST", "http://"+deviceIp+":8081/zeroconf/switch", bytes.NewBuffer(bodyJSON))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "close")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// convert response to Response type
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println(err)
	}
	result.Info = GetInfo(result.Error)

	// if non-200 status code or debug enabled then dump out request/response info
	if config.Debug || resp.StatusCode != 200 || result.Error != 0 {

		fmt.Printf("Request:\n %v\n\n", string(bodyJSON))

		fmt.Printf("Response:\n %v\n\n", string(respBody))
	}

	return result, nil
}

func GetInfo(ErrorCode int) string {
	switch ErrorCode {
	case 0:
		return "Success."
	case 400:
		return "The operation failed and the request was formatted incorrectly. The request body is not a valid JSON format."
	case 401:
		return "The operation failed and the request was unauthorized. Device information encryption is enabled on the device, but the request is not encrypted."
	case 404:
		return "The operation failed and the device does not exist. The device does not support the requested deviceid."
	case 422:
		return "The operation failed and the request parameters are invalid. For example, the device does not support setting specific device information."
	}
	return "Unknown error code."
}

// random hex code string, length 16
func generateIV() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:16], nil
}

func generateBody(data Data, key string) Request {
	var body Request
	var encryptionKey = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	iv, _ := generateIV()
	body.Sequence = "0000000000000"
	body.DeviceId = "0000000000"
	body.SelfAPIKey = "123"
	body.Encrypt = true
	body.IV = base64.StdEncoding.EncodeToString([]byte(iv))

	// do the encoding
	json, _ := json.Marshal(data)
	encoded := opensslCommand(string(json), encryptionKey, iv)
	body.Data = encoded
	return body
}
