package secrets

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/docker/go-connections/sockets"
)

func TestHandler(t *testing.T) {
	p := &testPlugin{}
	h := NewHandler(p)
	l := sockets.NewInmemSocket("test", 0)
	go h.Serve(l)
	defer l.Close()

	client := &http.Client{Transport: &http.Transport{
		Dial: l.Dial,
	}}

	resp, err := pluginRequest(client, getPath, Request{SecretName: "my-secret"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatalf("error while getting secret: %v", resp.Err)
	}
	if p.get != 1 {
		t.Fatalf("expected get 1, got %d", p.get)
	}
	if !bytes.EqualFold(secret, resp.Value) {
		t.Fatalf("expecting secret value %s, got %s", secret, resp.Value)
	}
	resp, err = pluginRequest(client, getPath, Request{SecretName: ""})
	if err != nil {
		t.Fatal(err)
	}
	if p.get != 2 {
		t.Fatalf("expected get 2, got %d", p.get)
	}
	if resp.Err == "" {
		t.Fatalf("expected missing secret")
	}
	resp, err = pluginRequest(client, getPath, Request{SecretName: "another-secret", SecretLabels: map[string]string{"prefix": "p-"}})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err != "" {
		t.Fatalf("error while getting secret: %v", resp.Err)
	}
	if !bytes.EqualFold(append([]byte("p-"), secret...), resp.Value) {
		t.Fatalf("expecting secret value %s, got %s", secret, resp.Value)
	}
}

func pluginRequest(client *http.Client, method string, req Request) (*Response, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp Response
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

type testPlugin struct {
	get int
}

var secret = []byte("secret")

func (p *testPlugin) Get(req Request) Response {
	p.get++
	if req.SecretName == "" {
		return Response{Err: "missing secret name"}
	}
	if prefix, exists := req.SecretLabels["prefix"]; exists {
		return Response{Value: append([]byte(prefix), secret...)}
	}
	return Response{Value: secret}
}
