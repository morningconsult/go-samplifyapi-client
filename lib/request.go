package samplify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

// APIResponse ...
type APIResponse struct {
	Body      json.RawMessage
	RequestID string
}

// SendRequest exposing sendrequest to enable custom requests
func SendRequest(host, method, url, accessToken string, body interface{}, timeout int) (*APIResponse, error) {
	dur := time.Duration(timeout)
	httpClient := newDefaultHTTPClient()
	httpClient.SetTimeout(dur)
	c := NewClient("", "", "", httpClient, nil)
	return c.sendRequest(host, method, url, accessToken, body)
}

func (c *Client) sendRequestWithContext(
	ctx context.Context,
	host string,
	method string,
	url string,
	accessToken string,
	body interface{},
) (*APIResponse, error) {
	path := fmt.Sprintf("%s%s", host, url)

	newRequest := func(body io.Reader) (*http.Request, error) {
		return http.NewRequestWithContext(ctx, method, path, body)
	}

	return c.makeRequest(path, newRequest, accessToken, body)
}

func (c *Client) sendRequest(host, method, url, accessToken string, body interface{}) (*APIResponse, error) {
	path := fmt.Sprintf("%s%s", host, url)

	newRequest := func(body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, path, body)
	}

	return c.makeRequest(path, newRequest, accessToken, body)
}

func (c *Client) makeRequest(
	path string,
	newRequest func(io.Reader) (*http.Request, error),
	accessToken string,
	body interface{},
) (*APIResponse, error) {
	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	var buf *bytes.Buffer
	if reqData != nil {
		buf = bytes.NewBuffer(reqData)
	}

	req, err := newRequest(buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")

	if accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ar := &APIResponse{RequestID: resp.Header.Get("x-request-id")}

	if resp.StatusCode >= http.StatusBadRequest {
		t := time.Now()
		err := &ErrorResponse{
			Timestamp:  &t,
			RequestID:  ar.RequestID,
			HTTPCode:   resp.StatusCode,
			HTTPPhrase: resp.Status,
			Path:       path,
			Errors:     []*Error{{Path: path, Message: resp.Status}},
		}

		ar.Body = json.RawMessage(respData)

		return ar, err
	}

	ar.Body = json.RawMessage(respData)

	return ar, nil
}

func (c *Client) sendFormDataWithContext(
	ctx context.Context,
	host string,
	method string,
	path string,
	accessToken string,
	file multipart.File,
	fileName string,
	message string,
) (*APIResponse, error) {
	path = fmt.Sprintf("%s%s", host, path)

	newRequest := func(body io.Reader) (*http.Request, error) {
		return http.NewRequestWithContext(ctx, method, path, body)
	}

	return c.makeFormDataRequest(path, newRequest, accessToken, file, fileName, message)
}

func (c *Client) sendFormData(
	host string,
	method string,
	path string,
	accessToken string,
	file multipart.File,
	fileName string,
	message string,
) (*APIResponse, error) {
	path = fmt.Sprintf("%s%s", host, path)

	newRequest := func(body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, path, body)
	}

	return c.makeFormDataRequest(path, newRequest, accessToken, file, fileName, message)
}

func (c *Client) makeFormDataRequest(
	path string,
	newRequest func(io.Reader) (*http.Request, error),
	accessToken string,
	file multipart.File,
	fileName string,
	message string,
) (*APIResponse, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(fileWriter, file); err != nil {
		return nil, err
	}

	bodyWriter.WriteField("message", message)
	bodyWriter.Close()

	req, err := newRequest(bodyBuf)
	if err != nil {
		return nil, err
	}

	if accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ar := &APIResponse{RequestID: resp.Header.Get("x-request-id")}

	if resp.StatusCode >= http.StatusBadRequest {
		t := time.Now()
		err := &ErrorResponse{
			Timestamp:  &t,
			RequestID:  ar.RequestID,
			HTTPCode:   resp.StatusCode,
			HTTPPhrase: resp.Status,
			Path:       path,
			Errors:     []*Error{{Path: path, Message: resp.Status}},
		}

		ar.Body = json.RawMessage(respData)

		return ar, err
	}

	ar.Body = json.RawMessage(respData)

	return ar, nil
}
