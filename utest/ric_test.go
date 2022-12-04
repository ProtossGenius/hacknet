package main

import (
	"fmt"
	"testing"
)

// TestClient test clientf.
type TestClient struct {
}

func (t *TestClient) SendForward(email string, data []byte) {
	fmt.Println("recv email .. ", email)
}

func (t *TestClient) OnForward(data []byte) {
	panic("not implemented") // TODO: Implement
}

func TestHelloWorld(t *testing.T) {
}
