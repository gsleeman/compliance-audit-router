package splunk

import "time"

// AlertDetails is a structured Splunk alert details
type AlertDetails struct {
	SearchResult
	AlertName  string
	User       string
	Group      string
	Timestamp  time.Time
	ClusterIDs []string
	Reasons    []string
}

// NewAlertDetails creates a new AlertDetails from a SearchResult
func NewAlertDetails(result SearchResult) AlertDetails {
	return AlertDetails{
		SearchResult: result,
		AlertName:    result.string("alertname"),
		User:         result.string("username"),
		Group:        result.string("group"),
		Timestamp:    result.time("timestamp"),
		ClusterIDs:   result.slice("clusterid"),
		Reasons:      result.slice("reason"),
	}
}

// Webhook.Details returns the AlertDetails for a webhook
func (w Webhook) Details() AlertDetails {
	return NewAlertDetails(w.Result)
}

// Alert.Details returns a slice of AlertDetails from the alert search results
func (w Alert) Details() []AlertDetails {
	alerts := []AlertDetails{}
	for _, result := range w.SearchResults.Results {
		alerts = append(alerts, NewAlertDetails(result))
	}
	return alerts
}
