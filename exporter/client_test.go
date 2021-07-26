// Package exporter
// Time    : 2021/7/22 5:22 下午
// Author  : xushiyin
// contact : yuqingxushiyin@gmail.com
package exporter

import (
	"fmt"
	"testing"
)

func TestGetLogstashRoot(t *testing.T) {
	tests := []struct {
		baseUrl string
		path    string
		version string
	}{
		{"http://192.168.210.250:9601", "/", "6.8.12"},
	}
	for _, ts := range tests {
		nrc := NewReqClient(ts.baseUrl)

		ri, err := GetLogstashRootInfo(nrc, ts.path)
		if err != nil {
			t.Errorf("err %#v", err)
		}
		if ri.Version == ts.version {
			fmt.Println(ri)
		} else {
			t.Errorf("incorrect version %s", ri.Version)
		}
	}
}

func TestGetLogstashNodeStats(t *testing.T) {
	tests := []struct {
		baseUrl string
		path    string
		version string
	}{
		{"http://192.168.210.250:9601", "/_node/stats", "5.0.0"},
		{"http://192.168.210.250:9600", "/_node/stats", "6.8.12"},
		{"http://192.168.210.250:9601", "/_node/stats", "7.3.0"},
	}
	for _, ts := range tests {
		nrc := NewReqClient(ts.baseUrl)
		nsi, err := GetLogstashNodeStats(nrc, ts.path)
		if err != nil {
			t.Errorf("err %#v", err)
		}
		if nsi.Version == ts.version {
			fmt.Println(nsi)
		} else {
			t.Errorf("incorrect version %s", nsi.Version)
		}
	}
}
