package v2

import "testing"

func TestHTTPServer_Start(t *testing.T) {
	s := &HTTPServer{}
	s.Start(":8081")
}
