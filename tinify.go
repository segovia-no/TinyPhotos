package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
)

const TINIFY_API_URL = "https://api.tinify.com"

type TinifyClient struct {
	ApiKey string
}

type TinifyResponse struct {
	Headers TinifyHeadersResponse
	Input   TinifyImageResponse `json:"input,omitempty"`
	Output  TinifyImageResponse `json:"output,omitempty"`
}

type TinifyImageResponse struct {
	Size   int     `json:"size,omitempty"`
	Type   string  `json:"type,omitempty"`
	Width  int     `json:"width,omitempty"`
	Height int     `json:"height,omitempty"`
	Ratio  float32 `json:"ratio,omitempty"`
	Url    string  `json:"url,omitempty"`
}

type TinifyHeadersResponse struct {
	CompressionCount int
	Location         string
	ContentType      string
}

type TinifyPreserveBody struct {
	Preserve []string `json:"preserve,omitempty"`
}

func (t *TinifyClient) SetAPIKey(apiKey string) error {
	if apiKey == "" {
		return errors.New("cannot set an empty API key")
	}

	t.ApiKey = apiKey
	return nil
}

func (t *TinifyClient) setAuthHeader(req *http.Request) error {
	if t.ApiKey == "" {
		return errors.New("cannot set Auth header without an API key")
	}

	authString := "api:" + t.ApiKey
	b64Auth := base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Set("Authorization", "Basic "+b64Auth)
	return nil
}

func (t *TinifyClient) MakeRequest(path string, inputFilename string) (TinifyResponse, error) {
	if t.ApiKey == "" {
		return TinifyResponse{}, errors.New("cannot make request without an API key")
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		return TinifyResponse{}, errors.New("couldnt open the input file")
	}
	defer inputFile.Close()

	reqBody := bufio.NewReader(inputFile)

	req, err := http.NewRequest(http.MethodPost, TINIFY_API_URL+path, reqBody)
	if err != nil {
		return TinifyResponse{}, errors.New("couldnt create a new request: " + err.Error())
	}

	t.setAuthHeader(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TinifyResponse{}, errors.New("Error making a POST request to tinify API: " + err.Error())
	}
	defer res.Body.Close()

	stringBody, err := io.ReadAll(res.Body)
	if err != nil {
		return TinifyResponse{}, errors.New("Error reading response body: " + err.Error())
	}

	var jsonRes TinifyResponse
	json.Unmarshal(stringBody, &jsonRes)

	headerCompressionCount, _ := strconv.Atoi(res.Header.Get("Compression-Count"))
	jsonRes.Headers.CompressionCount = headerCompressionCount
	jsonRes.Headers.Location = res.Header.Get("Location")
	jsonRes.Headers.ContentType = res.Header.Get("Content-Type")

	return jsonRes, nil
}

func (t *TinifyClient) DownloadWithMetadata(locationPath string, outputFilepath string) error {
	if locationPath == "" {
		return errors.New("path to location of resulting image cannot be empty")
	}

	if outputFilepath == "" {
		return errors.New("output file path cannot be empty")
	}

	reqStruct := &TinifyPreserveBody{
		Preserve: []string{"copyright", "creation", "location"},
	}

	jsonData, err := json.Marshal(reqStruct)
	if err != nil {
		return errors.New("couldnt marshal JSON post data")
	}

	req, err := http.NewRequest(http.MethodPost, locationPath, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("couldnt create a new request: " + err.Error())
	}

	t.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Error making a POST request to tinify API: " + err.Error())
	}
	defer res.Body.Close()

	file, err := os.Create(outputFilepath)
	if err != nil {
		return errors.New("Error creating output file: " + err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return errors.New("Error writing to output file: " + err.Error())
	}

	return nil
}
