package splunk

import (
	"reflect"
	"testing"
)

func TestNewAlertDetails(t *testing.T) {
	type args struct {
		result SearchResult
	}
	tests := []struct {
		name string
		args args
		want AlertDetails
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlertDetails(tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlertDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhook_Details(t *testing.T) {
	tests := []struct {
		name string
		w    Webhook
		want AlertDetails
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Details(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Webhook.Details() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlert_Details(t *testing.T) {
	tests := []struct {
		name string
		w    Alert
		want []AlertDetails
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Details(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Alert.Details() = %v, want %v", got, tt.want)
			}
		})
	}
}
