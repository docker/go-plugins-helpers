package authorization

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/docker/go-plugins-helpers/sdk"
)

type TestPlugin struct {
	Plugin
}

func (p *TestPlugin) AuthZReq(r Request) Response {
	return Response{
		Allow: false,
		Msg:   "You are not authorized",
		Err:   "",
	}
}

func (p *TestPlugin) AuthZRes(r Request) Response {
	return Response{
		Allow: false,
		Msg:   "You are not authorized",
		Err:   "",
	}
}

func TestActivate(t *testing.T) {
	response, err := http.Get("http://localhost:32456/Plugin.Activate")

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	if string(body) != manifest+"\n" {
		t.Fatalf("Expected %s, got %s\n", manifest+"\n", string(body))
	}
}

func TestAuthZReq(t *testing.T) {
	request := `{"User":"bob","UserAuthNMethod":"","RequestMethod":"POST","RequestURI":"http://127.0.0.1/v.1.23/containers/json","RequestBody":"","RequestHeader":"","RequestStatusCode":""}`

	response, err := http.Post(
		"http://localhost:32456/AuthZPlugin.AuthZReq",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatal(err)
	}

	if r.Msg != "You are not authorized" {
		t.Fatal("Authorization message does not match")
	}

	if r.Allow {
		t.Fatal("The request has been allowed while should not be")
	}

	if r.Err != "" {
		t.Fatal("Authorization Error should be empty")
	}
}

func TestAuthZRes(t *testing.T) {
	request := `{"User":"bob","UserAuthNMethod":"","RequestMethod":"POST","RequestURI":"http://127.0.0.1/v.1.23/containers/json","RequestBody":"","RequestHeader":"","RequestStatusCode":"", "ResponseBody":"","ResponseHeader":"","ResponseStatusCode":200}`

	response, err := http.Post(
		"http://localhost:32456/AuthZPlugin.AuthZRes",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request),
	)

	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatal(err)
	}

	if r.Msg != "You are not authorized" {
		t.Fatal("Authorization message does not match")
	}

	if r.Allow {
		t.Fatal("The request has been allowed while should not be")
	}

	if r.Err != "" {
		t.Fatal("Authorization Error should be empty")
	}
}

func TestMain(m *testing.M) {
	d := &TestPlugin{}
	h := NewHandler(d)
	go h.ServeTCP("test", ":32456", nil)

	os.Exit(m.Run())
}
