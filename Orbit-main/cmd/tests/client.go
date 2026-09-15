package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

const BaseURL = "http://localhost:9132"

type Client struct {
	client *http.Client
}

func NewClient() *Client {

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return &Client{
		client: &http.Client{
			Jar: jar,
			Timeout: 20 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        6000,
				MaxIdleConnsPerHost: 6000,
				MaxConnsPerHost:     6000,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (c *Client) do(req *http.Request) ([]byte, int, error) {

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func (c *Client) Get(path string) ([]byte, int, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		BaseURL+path,
		nil,
	)
	if err != nil {
		return nil, 0, err
	}

	return c.do(req)
}

func (c *Client) Delete(path string) ([]byte, int, error) {

	req, err := http.NewRequest(
		http.MethodDelete,
		BaseURL+path,
		nil,
	)
	if err != nil {
		return nil, 0, err
	}

	return c.do(req)
}

func (c *Client) Post(path string, payload any) ([]byte, int, error) {

	var body io.Reader

	if payload != nil {

		b, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, err
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		BaseURL+path,
		body,
	)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.do(req)
}

func (c *Client) Put(path string, payload any) ([]byte, int, error) {

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		BaseURL+path,
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.do(req)
}

func (c *Client) UploadMultipart(
	method string,
	path string,
	fields map[string]string,
	fileField string,
	filePath string,
) ([]byte, int, error) {

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, 0, err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fileField, file.Name())
	if err != nil {
		return nil, 0, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, 0, err
	}

	if err := writer.Close(); err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(
		method,
		BaseURL+path,
		&body,
	)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.do(req)
}

func Decode[T any](body []byte) T {

	var v T

	if err := json.Unmarshal(body, &v); err != nil {
		panic(err)
	}

	return v
}
