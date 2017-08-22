package mdbstitch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const DEFAULT_STITCH_SERVER_URL = "https://stitch.mongodb.com/api/client/v1.0/app/"

type StitchClient struct {
	AppId       string
	AccessToken string
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserId       string `json:"userId"`
	DeviceId     string `json:"deviceId"`
}

type StitchResult struct {
	Result []interface{} `json:"result"`
}

type APIKeyForm struct {
	Key     string `json:"key"`
	Options struct {
		Device struct {
			AppId           string `json:"appId"`
			AppVersion      string `json:"appVersion"`
			SdkVersion      string `json:"sdkVersion"`
			DeviceId        string `json:"deviceId"`
			Platform        string `json:"platform"`
			PlatformVersion string `json:"platformVersion"`
		} `json:"device"`
	} `json:"options"`
}

func (sc StitchClient) APIKeyAuth(appKey string) (AuthResponse, error) {
	var uri bytes.Buffer
	uri.WriteString(DEFAULT_STITCH_SERVER_URL)
	uri.WriteString(sc.AppId)
	uri.WriteString("/auth/api/key")
	var response AuthResponse

	var data APIKeyForm
	data.Key = appKey
	data.Options.Device.AppId = sc.AppId
	data.Options.Device.AppVersion = ""
	data.Options.Device.SdkVersion = "0.0.24"
	data.Options.Device.DeviceId = ""
	data.Options.Device.Platform = "chrome"
	data.Options.Device.PlatformVersion = "60.0.3112"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return response, err
	}
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{},
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var client = &http.Client{Timeout: 30 * time.Second, Transport: tr}
	req, err := http.NewRequest("POST", uri.String(), bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	if err != nil {
		panic(err)
	}
	return response, nil
}

func (sc StitchClient) Query(action string, db string, coll string, query *bson.M, project *bson.M, limit int) (StitchResult, error) {
	var stitchResult StitchResult

	var uri bytes.Buffer
	uri.WriteString(DEFAULT_STITCH_SERVER_URL)
	uri.WriteString(sc.AppId)
	uri.WriteString("/pipeline")

	data := []bson.M{bson.M{
		"service": "mongodb-atlas",
		"action":  action,
		"args": bson.M{
			"database":   db,
			"collection": coll,
			"query":      query,
			"project":    project,
			"limit":      limit,
		}}}

	jsonData, _ := json.Marshal(data)
	var client = &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("POST", uri.String(), bytes.NewBuffer(jsonData))
	var token bytes.Buffer
	token.WriteString("Bearer ")
	token.WriteString(sc.AccessToken)
	req.Header.Add("Authorization", token.String())
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return stitchResult, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&stitchResult)
	if err != nil {
		panic(err)
	}
	return stitchResult, nil
}
