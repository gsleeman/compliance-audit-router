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
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/openshift/compliance-audit-router/pkg/config"
	"github.com/openshift/compliance-audit-router/pkg/helpers"
)

// Webhook is the JSON structure for a Splunk webhook
type Webhook struct {
	Sid         string       `json:"sid"`
	SearchName  string       `json:"search_name"`
	App         string       `json:"app"`
	Owner       string       `json:"owner"`
	ResultsLink string       `json:"results_link"`
	Result      SearchResult `json:"result"`
}

// Alert describes a Splunk alert
type Alert struct {
	SearchID      string
	SearchResults SearchResults
}

// searchResult represents an actual SPLUNK search result
type SearchResult map[string]interface{}

// searchResults represents the results of a Splunk API */results call
type SearchResults struct {
	InitOffset  int                 `json:"init_offset"`
	Messages    []map[string]string `json:"messages"`
	Preview     bool                `json:"preview"`
	Results     []SearchResult      `json:"results"`
	Highlighted map[string]string   `json:"highlighted"`
}

type SplunkServer config.SplunkConfig

// NOTE: The webhook itself contains the search result. So this may not be necessary

// RetrieveSearchFromAlert parses the received webhook, and looks up the data for the alert in Splunk,
// and returns the information in an Alert struct
func (s SplunkServer) RetrieveSearchFromAlert(sid string) (Alert, error) {

	splunkHttpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: s.AllowInsecure,
			},
		},
	}

	url := fmt.Sprintf("%s/services/search/jobs/%s/results?output_mode=json", s.Host, sid)

	var alert = Alert{
		SearchID:      sid,
		SearchResults: SearchResults{},
	}

	// Create a new HTTP client; don't modify the default client
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return alert, err
	}

	bearerToken := "Bearer" + s.Token
	req.Header.Add("Authorization", bearerToken)

	resp, err := splunkHttpClient.Do(req)
	if err != nil {
		return alert, err
	}

	// Process the response
	err = helpers.DecodeJSONResponseBody(resp, &alert.SearchResults)
	if err != nil {
		return alert, err
	}

	log.Printf("Received alert from Splunk: %s, %v", alert.SearchID, alert.Details())
	return alert, err

}
