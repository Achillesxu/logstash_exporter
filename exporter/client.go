// Package exporter
// Time    : 2021/7/22 2:25 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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

type NodeRootInfo struct {
	Host        string `json:"host"`
	Version     string `json:"version"`
	HttpAddress string `json:"http_address"`
}

// Pipeline type
type Pipeline struct {
	Events struct {
		DurationInMillis int `json:"duration_in_millis"`
		In               int `json:"in"`
		Filtered         int `json:"filtered"`
		Out              int `json:"out"`
	} `json:"events"`
	Plugins struct {
		Inputs []struct {
			ID     string `json:"id"`
			Events struct {
				In  int `json:"in"`
				Out int `json:"out"`
			} `json:"events"`
			Name string `json:"name"`
		} `json:"inputs,omitempty"`
		Filters []struct {
			ID     string `json:"id"`
			Events struct {
				DurationInMillis int `json:"duration_in_millis"`
				In               int `json:"in"`
				Out              int `json:"out"`
			} `json:"events,omitempty"`
			Name             string `json:"name"`
			Matches          int    `json:"matches,omitempty"`
			Failures         int    `json:"failures,omitempty"`
			PatternsPerField struct {
				CapturedRequestHeaders int `json:"captured_request_headers"`
			} `json:"patterns_per_field,omitempty"`
			Formats int `json:"formats,omitempty"`
		} `json:"filters"`
		Outputs []struct {
			ID     string `json:"id"`
			Events struct {
				In  int `json:"in"`
				Out int `json:"out"`
			} `json:"events"`
			Name string `json:"name"`
		} `json:"outputs"`
	} `json:"plugins"`
	Reloads struct {
		LastError            interface{} `json:"last_error"`
		Successes            int         `json:"successes"`
		LastSuccessTimestamp interface{} `json:"last_success_timestamp"`
		LastFailureTimestamp interface{} `json:"last_failure_timestamp"`
		Failures             int         `json:"failures"`
	} `json:"reloads"`
	Queue struct {
		Events   int    `json:"events"`
		Type     string `json:"type"`
		Capacity struct {
			PageCapacityInBytes int   `json:"page_capacity_in_bytes"`
			MaxQueueSizeInBytes int64 `json:"max_queue_size_in_bytes"`
			MaxUnreadEvents     int   `json:"max_unread_events"`
		} `json:"capacity"`
		Data struct {
			Path             string `json:"path"`
			FreeSpaceInBytes int64  `json:"free_space_in_bytes"`
			StorageType      string `json:"storage_type"`
		} `json:"data"`
	} `json:"queue"`
	DeadLetterQueue struct {
		QueueSizeInBytes int `json:"queue_size_in_bytes"`
	} `json:"dead_letter_queue"`
}

// NodeStatsInfo type
type NodeStatsInfo struct {
	Host        string `json:"host"`
	Version     string `json:"version"`
	HTTPAddress string `json:"http_address"`
	Jvm         struct {
		Threads struct {
			Count     int `json:"count"`
			PeakCount int `json:"peak_count"`
		} `json:"threads"`
		Mem struct {
			HeapUsedInBytes         int `json:"heap_used_in_bytes"`
			HeapUsedPercent         int `json:"heap_used_percent"`
			HeapCommittedInBytes    int `json:"heap_committed_in_bytes"`
			HeapMaxInBytes          int `json:"heap_max_in_bytes"`
			NonHeapUsedInBytes      int `json:"non_heap_used_in_bytes"`
			NonHeapCommittedInBytes int `json:"non_heap_committed_in_bytes"`
			Pools                   struct {
				Survivor struct {
					PeakUsedInBytes  int `json:"peak_used_in_bytes"`
					UsedInBytes      int `json:"used_in_bytes"`
					PeakMaxInBytes   int `json:"peak_max_in_bytes"`
					MaxInBytes       int `json:"max_in_bytes"`
					CommittedInBytes int `json:"committed_in_bytes"`
				} `json:"survivor"`
				Old struct {
					PeakUsedInBytes  int `json:"peak_used_in_bytes"`
					UsedInBytes      int `json:"used_in_bytes"`
					PeakMaxInBytes   int `json:"peak_max_in_bytes"`
					MaxInBytes       int `json:"max_in_bytes"`
					CommittedInBytes int `json:"committed_in_bytes"`
				} `json:"old"`
				Young struct {
					PeakUsedInBytes  int `json:"peak_used_in_bytes"`
					UsedInBytes      int `json:"used_in_bytes"`
					PeakMaxInBytes   int `json:"peak_max_in_bytes"`
					MaxInBytes       int `json:"max_in_bytes"`
					CommittedInBytes int `json:"committed_in_bytes"`
				} `json:"young"`
			} `json:"pools"`
		} `json:"mem"`
		Gc struct {
			Collectors struct {
				Old struct {
					CollectionTimeInMillis int `json:"collection_time_in_millis"`
					CollectionCount        int `json:"collection_count"`
				} `json:"old"`
				Young struct {
					CollectionTimeInMillis int `json:"collection_time_in_millis"`
					CollectionCount        int `json:"collection_count"`
				} `json:"young"`
			} `json:"collectors"`
		} `json:"gc"`
	} `json:"jvm"`
	Process struct {
		OpenFileDescriptors     int `json:"open_file_descriptors"`
		PeakOpenFileDescriptors int `json:"peak_open_file_descriptors"`
		MaxFileDescriptors      int `json:"max_file_descriptors"`
		Mem                     struct {
			TotalVirtualInBytes int64 `json:"total_virtual_in_bytes"`
		} `json:"mem"`
		CPU struct {
			TotalInMillis int64 `json:"total_in_millis"`
			Percent       int   `json:"percent"`
		} `json:"cpu"`
	} `json:"process"`
	Pipeline  Pipeline            `json:"pipeline"`  // Logstash 5
	Pipelines map[string]Pipeline `json:"pipelines"` // Logstash >=6
}

// NewReqClient get a request client use default RoundTripper
func NewReqClient(baseUrl string) *ReqClient {
	return &ReqClient{
		baseUrl,
		&http.Client{},
	}
}

// Get returns a GET request
func (rc *ReqClient) Get(path string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", rc.BaseUrl, path), bytes.NewBuffer([]byte{}))
}

// GetQuery returns a GET request with query params
func (rc *ReqClient) GetQuery(path string, params interface{}) (*http.Request, error) {
	queryString, err := query.Values(params)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("get queryString error"))
		return nil, err
	}
	return http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s?%s", rc.BaseUrl, path, queryString.Encode()), bytes.NewBuffer([]byte{}))
}

// Do func returns a response with your data
func (rc *ReqClient) Do(request *http.Request, duration time.Duration) (*ResponseStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	request = request.WithContext(ctx)

	response, reqErr := rc.hc.Do(request)
	if reqErr != nil {
		reqErr = errors.Wrap(reqErr, fmt.Sprintf("%s %s error", request.Method, request.RequestURI))
		return nil, reqErr
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			reqErr = errors.Wrap(err, fmt.Sprintf("%s %s close body error", request.Method, request.RequestURI))
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("%s %s error", request.Method, request.RequestURI))
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

// GetLogstashRootInfo get Logstash root info
func GetLogstashRootInfo(rc *ReqClient, path string, milliseconds int64) (*NodeRootInfo, error) {
	reqGet, err := rc.Get(path)
	if err != nil {
		return nil, err
	}
	resp, err := rc.Do(reqGet, time.Duration(milliseconds)*time.Millisecond)
	if err != nil {
		return nil, err
	}
	host := gjson.GetBytes(resp.Body, "host")
	version := gjson.GetBytes(resp.Body, "version")
	httpAddress := gjson.GetBytes(resp.Body, "http_address")

	rootInfo := NodeRootInfo{
		host.String(),
		version.String(),
		httpAddress.String(),
	}
	return &rootInfo, nil
}

func GetLogstashNodeStats(rc *ReqClient, path string, milliseconds int64) (*NodeStatsInfo, error) {
	reqGet, err := rc.Get(path)
	if err != nil {
		return nil, err
	}
	resp, err := rc.Do(reqGet, time.Duration(milliseconds)*time.Millisecond)
	if err != nil {
		return nil, err
	}
	nsi := &NodeStatsInfo{}

	err = json.Unmarshal(resp.Body, nsi)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Unmarshal body <%#v> from <%s>", resp.Body, reqGet.RequestURI))
		return nil, err
	}
	return nsi, nil
}
