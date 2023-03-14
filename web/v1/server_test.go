package v1

import "testing"

func TestServer(t *testing.T) {
	server := &HTTPServer{}
	server.Start(":8080")
}
