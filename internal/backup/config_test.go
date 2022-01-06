package backup

import (
	"testing"
)

func TestGetClientIP(t *testing.T) {
	URL := "http://127.0.0.1:8080/ip"
	IP := GetClientIP(URL)
	t.Log(IP)
}
