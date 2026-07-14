package utils

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

const FORWARD string = "X-Forwarded-For"
const REAL_IP string = "X-Real-IP"

func GetIPv4fromRequestRemoteAddr(request *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("error to process request %s", err)
	}
	ipFromAddress := net.ParseIP(host)
	if ipFromAddress == nil {
		return "", fmt.Errorf("error to parsing hostname %s", err)
	}
	return ipFromAddress.String(), nil
}

func GetIPv4FromRequestHeader(request *http.Request) (string, error) {
	ipFromHeader := request.Header.Get(FORWARD)
	if ipFromHeader == "" {
		return "", fmt.Errorf("error to parsing hostname")
	}
	return ipFromHeader, nil
}

func GetIPAddress(request *http.Request) string {
	if ip := request.Header.Get(FORWARD); ip != "" {
		return strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	if ip := request.Header.Get(REAL_IP); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return request.RemoteAddr
	}

	return host
}
