// Package exporter
// Time    : 2021/7/22 2:25 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import (
	"bytes"
	"fmt"
	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type ReqClient struct {
	BaseUrl string
	hc      *http.Client
}

// ResponseStruct is a struct who returns after requests
type ResponseStruct struct {
	Status        string
	StatusCode    int
	Header        http.Header
	ContentLength int64
	Body          []byte
}

// NewReqClient get a request client use default RoundTripper
func NewReqClient(baseUrl string) *ReqClient {
	return &ReqClient{
		baseUrl,
		&http.Client{},
	}
}

// Get get a GET Request
func (rc *ReqClient) Get(path string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", rc.BaseUrl, path), bytes.NewBuffer([]byte{}))
}

// GetWith func returns a request
func (rc *ReqClient) GetWith(path string, params interface{}) (*http.Request, error) {
	queryString, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s?%s", rc.BaseUrl, path, queryString.Encode()), bytes.NewBuffer([]byte{}))
}

// Do func returns a response with your data
func (rc *ReqClient) Do(request *http.Request) (*ResponseStruct, error) {
	response, err := rc.hc.Do(request)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("%s req <%s> close body failed, err: <%#v>", request.Method, request.RequestURI, err)
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &ResponseStruct{
		Status:        response.Status,
		StatusCode:    response.StatusCode,
		Header:        response.Header,
		ContentLength: response.ContentLength,
		Body:          body,
	}, nil
}
