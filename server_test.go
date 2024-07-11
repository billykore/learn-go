package learning

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	server := http.Server{
		Addr: "localhost:8080",
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
