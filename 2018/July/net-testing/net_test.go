package net

import (
	"bytes"
	"log"
	"net"
	"testing"
)

var srv Server

func init() {
	// Start the new server
	srv, err := NewServer("tcp", ":1123")
	if err != nil {
		log.Println("error starting TCP server")
		return
	}

	// Run the server in goroutine to stop blocking
	go func() {
		srv.Run()
	}()
}

func TestNETServer_Running(t *testing.T) {
	// Simply check that the server is up and can
	// accept connections.
	conn, err := net.Dial("tcp", ":1123")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()
}

func TestNETServer_Request(t *testing.T) {
	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		{"Sending a simple request returns result", []byte("hello world\n"), []byte("Request received: hello world")},
		{"Sending another simple request works", []byte("goodbye world\n"), []byte("Request received: goodbye world")},
	}

	for _, tc := range tt {
		t.Run(tc.test, func(t *testing.T) {
			conn, err := net.Dial("tcp", ":1123")
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer conn.Close()

			if _, err := conn.Write(tc.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
			}

			out := make([]byte, 1024)
			if _, err := conn.Read(out); err == nil {
				if bytes.Compare(out, tc.want) == 0 {
					t.Error("response did match expected output")
				}
			} else {
				t.Error("could not read from connection")
			}
		})
	}
}
