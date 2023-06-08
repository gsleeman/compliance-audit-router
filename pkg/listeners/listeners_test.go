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

package listeners

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestListenerURIs(t *testing.T) {
	for _, l := range Listeners {
		_, err := url.ParseRequestURI(l.Path)
		if err != nil {
			t.Errorf("%s is not a valid url", l.Path)
		}
	}
}

func TestInitRoutes(t *testing.T) {
	r := chi.NewRouter()
	InitRoutes(r)

	expectedRouteLen := 4
	if routeLen := len(r.Routes()); routeLen != expectedRouteLen {
		t.Errorf("Error initializing routes. Expected %v but got %v.", expectedRouteLen, routeLen)
	}

	paths := []string{"/readyz", "/healthz", "/api/v1/alert", "/api/v1/jira_webhook"}

	for _, route := range r.Routes() {
		found := false
		for i, path := range paths {
			if route.Pattern == path {
				found = true
				paths = append(paths[:i], paths[i+1:]...)
				break
			}
		}
		if !found {
			t.Errorf("Unexpected route pattern found: %v", route.Pattern)
		}
	}
}

func TestRespondOKHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	RespondOKHandler(recorder, nil)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if body := recorder.Body.String(); body != "OK" {
		t.Errorf("handler returned wrong body: got %v want %v", body, "OK")
	}
}

func TestProcessAlertHandler(t *testing.T) {
	// Example webhook payloads that might be received from the
	// alerting system (ie: Splunk)
	tests := []struct {
		name                string
		incomingWebhookBody string
		status              int
		contentType         string
		body                string
	}{
		{
			name:                "empty webhook should fail",
			incomingWebhookBody: "",
			status:              http.StatusBadRequest,
			contentType:         "text/plain; charset=utf-8",
			body:                "Request body must not be empty",
		},
	}

	for _, tt := range tests {
		// Create a test incoming webhook request to the http server
		req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(tt.incomingWebhookBody))
		if err != nil {
			t.Fatal(err)
		}

		// ResponseRecorder is used to store the server response
		recorder := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {
			ProcessAlertHandler(recorder, req)
			// Test the returned http status code
			if status := recorder.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v, want %v",
					status, tt.status)
			}

			// Test the returned header Content-Type
			if contentType := recorder.Header().Get("Content-Type"); contentType != tt.contentType {
				t.Errorf("handler returned wrong Content-Type: got %v, want %v",
					contentType, tt.contentType)
			}

			// Test the returned body
			if body := strings.TrimSpace(recorder.Body.String()); body != tt.body {
				t.Errorf("handler returned unexpected body: got %v, want %v",
					body, tt.body)
			}
		})
	}
}

func TestProcessJiraWebhook(t *testing.T) {
	tests := []struct {
		name                string
		incomingWebhookBody string
		status              int
		contentType         string
		expectedBody        string
	}{
		{
			name:                "empty webhook should fail",
			incomingWebhookBody: "",
			status:              http.StatusBadRequest,
			contentType:         "text/plain",
			expectedBody:        "failed to parse webhook",
		},
	}

	for _, tt := range tests {
		// Create a test incoming webhook request to the http server
		req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(tt.incomingWebhookBody))
		if err != nil {
			t.Fatal(err)
		}

		// ResponseRecorder is used to store the server response
		recorder := httptest.NewRecorder()

		t.Run(tt.name, func(t *testing.T) {
			ProcessJiraWebhook(recorder, req)
			// Test the returned http status code
			if status := recorder.Code; status != tt.status {
				t.Errorf("handler returned wrong status code: got %v, want %v",
					status, tt.status)
			}

			// Test the returned header Content-Type
			if contentType := recorder.Header().Get("Content-Type"); contentType != tt.contentType {
				t.Errorf("handler returned wrong Content-Type: got %v, want %v",
					contentType, tt.contentType)
			}

			// Test the returned body
			if body := strings.TrimSpace(recorder.Body.String()); body != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v, want %v",
					body, tt.expectedBody)
			}
		})
	}
}
