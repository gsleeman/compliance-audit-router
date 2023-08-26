/*
Copyright © 2021 Red Hat, Inc

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package splunk

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/openshift/compliance-audit-router/pkg/config"
)

const TEST_SEARCH_API_RESPONSE string = `{"preview":false,"init_offset":0,"messages":[],"fields":[{"name":"_time"},{"name":"alertname","type":"str"},{"name":"clusterid","type":"str"},{"name":"group","type":"str"},{"name":"timestamp","type":"str"},{"name":"username","type":"str"}],"results":[{"_time":"2023-08-27T02:25:00.000+10:00","alertname":"TestAlert","clusterid":"testcluster","group":"testgroup","timestamp":"2023-08-27T02:25:01.GMT","username":"testuser"}], "highlighted":{}}`

var EXPECTED_ALERT_DETAILS = []AlertDetails{{
	AlertName:  string("TestAlert"),
	User:       "testuser",
	Group:      "testgroup",
	Timestamp:  time.Date(2023, 8, 27, 02, 25, 1, 0, time.UTC),
	ClusterIDs: []string{"testcluster"},
	Reasons:    []string{},
}}

func NewTestServer() (*httptest.Server, *SplunkServer) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(TEST_SEARCH_API_RESPONSE))
	}))
	splunkserver := SplunkServer(config.SplunkConfig{
		Token:         "test",
		Host:          server.URL,
		AllowInsecure: true,
	})
	return server, &splunkserver
}

func TestSplunkServer_RetrieveSearchFromAlert(t *testing.T) {
	type args struct {
		sid string
	}
	server, splunkserver := NewTestServer()
	tests := []struct {
		name         string
		splunkserver SplunkServer
		args         args
		want         []AlertDetails
		wantErr      bool
	}{{
		"test",
		*splunkserver,
		args{"test"},
		EXPECTED_ALERT_DETAILS,
		false,
	},

	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer server.Close()
			got, err := tt.splunkserver.RetrieveSearchFromAlert(tt.args.sid)
			log.Println(got.Details())
			if (err != nil) != tt.wantErr {
				t.Errorf("SplunkServer.RetrieveSearchFromAlert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Details(), tt.want) {
				t.Errorf("SplunkServer.RetrieveSearchFromAlert() = %v, want %v", got.Details(), tt.want)
			}
		})
	}
}
