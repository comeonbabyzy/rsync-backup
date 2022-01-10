package backup

import (
	"testing"
)

func TestGetClientIP(t *testing.T) {
	URL := "http://192.168.191.143:8080/ip"
	IP := GetClientIP(URL)
	t.Log(IP)
}
