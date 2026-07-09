package utils

import (
	"net/http"
	"testing"
)

func TestGetIPv4fromRequestRemoteAddr(t *testing.T) {
	header := http.Header{}
	header.Add(FORWARD, "192.168.0.1")
	httpRequest := &http.Request{
		RemoteAddr: "192.168.0.1:9090",
		Header:     header,
	}
	type args struct {
		request *http.Request
	}
	input := args{
		request: httpRequest,
	}

	tests := []struct {
		name    string
		input   args
		want    string
		wantErr bool
	}{{
		name:    "valid",
		input:   input,
		want:    "192.168.0.1",
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIPv4fromRequestRemoteAddr(tt.input.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPv4fromRequestRemoteAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIPv4fromRequestRemoteAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIPv4FromRequestHeader(t *testing.T) {
	header := http.Header{}
	header.Add(FORWARD, "192.168.0.1")
	httpRequest := &http.Request{
		RemoteAddr: "192.168.0.1:9090",
		Header:     header,
	}
	type args struct {
		request *http.Request
	}
	input := args{
		request: httpRequest,
	}

	tests := []struct {
		name    string
		input   args
		want    string
		wantErr bool
	}{{
		name:    "valid",
		input:   input,
		want:    "192.168.0.1",
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIPv4FromRequestHeader(tt.input.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPv4FromRequestHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIPv4FromRequestHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIP(t *testing.T) {
	header := http.Header{}
	header.Add(FORWARD, "192.168.0.1")
	httpRequest := &http.Request{
		RemoteAddr: "192.168.0.1:9090",
		Header:     header,
	}
	type args struct {
		request *http.Request
	}
	input := args{
		request: httpRequest,
	}

	tests := []struct {
		name    string
		input   args
		want    string
		wantErr bool
	}{{
		name:    "valid",
		input:   input,
		want:    "192.168.0.1",
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIPAddress(tt.input.request)
			if got != tt.want {
				t.Errorf("GetIPv4FromRequestHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
