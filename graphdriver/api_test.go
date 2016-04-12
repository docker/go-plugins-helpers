package graphdriver

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

	// Init
	init, err := pluginInitRequest(client, InitRequest{Home: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if init.Err != "" {
		t.Fatalf("got error initialising graph drivers: %v", init.Err)
	}
	if p.init != 1 {
		t.Fatalf("expected init 1, got %d", p.init)
	}

	// Create
	create, err := pluginCreateRequest(client, CreateRequest{ID: "foo", Parent: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if create.Err != "" {
		t.Fatalf("got error creating graph drivers: %v", create.Err)
	}
	if p.create != 1 {
		t.Fatalf("expected create 1, got %d", p.create)
	}

	// Remove
	remove, err := pluginRemoveRequest(client, RemoveRequest{ID: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if remove.Err != "" {
		t.Fatalf("got error removing graph drivers: %s", remove.Err)
	}
	if p.remove != 1 {
		t.Fatalf("expected remove 1, got %d", p.remove)
	}

	// Get
	get, err := pluginGetRequest(client, GetRequest{ID: "foo", MountLabel: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if get.Err != "" {
		t.Fatalf("got error getting graph drivers: %s", get.Err)
	}
	if p.get != 1 {
		t.Fatalf("expected get 1, got %d", p.get)
	}

	// Put
	put, err := pluginPutRequest(client, PutRequest{ID: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if put.Err != "" {
		t.Fatalf("got error putting graph drivers: %s", put.Err)
	}
	if p.put != 1 {
		t.Fatalf("expected put 1, got %d", p.put)
	}

	// Exists
	exists, err := pluginExistsRequest(client, ExistsRequest{ID: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if !exists.Exists {
		t.Fatalf("got error testing for existence of graph drivers: %v", exists.Exists)
	}
	if p.exists != 1 {
		t.Fatalf("expected exists 1, got %d", p.exists)
	}

	// Status
	status, err := pluginStatusRequest(client, StatusRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if status.Status != nil {
		t.Fatalf("got error putting graph drivers: %s", status.Status)
	}
	if p.status != 1 {
		t.Fatalf("expected get 1, got %d", p.get)
	}

	// GetMetadata
	getMetadata, err := pluginGetMetadataRequest(client, GetMetadataRequest{ID: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if getMetadata.Err != "" {
		t.Fatalf("got error getting metadata for graph drivers: %s", getMetadata.Err)
	}
	if p.getMetadata != 1 {
		t.Fatalf("expected getMetadata 1, got %d", p.getMetadata)
	}

	// Cleanup
	cleanup, err := pluginCleanupRequest(client, CleanupRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if cleanup.Err != "" {
		t.Fatalf("got error cleaning graph drivers up: %s", cleanup.Err)
	}
	if p.cleanup != 1 {
		t.Fatalf("expected cleanup 1, got %d", p.cleanup)
	}

	// Diff
	_, err = pluginDiffRequest(client, DiffRequest{ID: "foo", Parent: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if p.diff != 1 {
		t.Fatalf("expected diff 1, got %d", p.diff)
	}

	// Changes
	changes, err := pluginChangesRequest(client, ChangesRequest{ID: "foo", Parent: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if changes.Err != "" {
		t.Fatalf("got error putting graph drivers: %s", changes.Err)
	}
	if p.status != 1 {
		t.Fatalf("expected get 1, got %d", p.get)
	}

	// ApplyDiff
	b := new(bytes.Buffer)
	stream := bytes.NewReader(b.Bytes())
	applyDiff, err := pluginApplyDiffRequest(client,
		ApplyDiffRequest{ID: "foo", Parent: "bar", Stream: stream})
	if err != nil {
		t.Fatal(err)
	}
	if applyDiff.Err != "" {
		t.Fatalf("got error applying diff on graph drivers: %s", applyDiff.Err)
	}
	if p.status != 1 {
		t.Fatalf("expected applyDiff 1, got %d", p.applyDiff)
	}

	// DiffSize
	diffSize, err := pluginDiffSizeRequest(client, DiffSizeRequest{ID: "foo", Parent: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if diffSize.Err != "" {
		t.Fatalf("got error getting diff size for graph drivers: %s", diffSize.Err)
	}
	if p.diffSize != 1 {
		t.Fatalf("expected diffSize 1, got %d", p.diffSize)
	}
}

func pluginInitRequest(client *http.Client, req InitRequest) (*InitResponse, error) {
	method := initPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp InitResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginCreateRequest(client *http.Client, req CreateRequest) (*CreateResponse, error) {
	method := createPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginRemoveRequest(client *http.Client, req RemoveRequest) (*RemoveResponse, error) {
	method := removePath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp RemoveResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginGetRequest(client *http.Client, req GetRequest) (*GetResponse, error) {
	method := getPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp GetResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginPutRequest(client *http.Client, req PutRequest) (*PutResponse, error) {
	method := putPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp PutResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginExistsRequest(client *http.Client, req ExistsRequest) (*ExistsResponse, error) {
	method := existsPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp ExistsResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginStatusRequest(client *http.Client, req StatusRequest) (*StatusResponse, error) {
	method := statusPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginGetMetadataRequest(client *http.Client, req GetMetadataRequest) (*GetMetadataResponse, error) {
	method := getMetadataPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp GetMetadataResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginCleanupRequest(client *http.Client, req CleanupRequest) (*CleanupResponse, error) {
	method := cleanupPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp CleanupResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginDiffRequest(client *http.Client, req DiffRequest) (*DiffResponse, error) {
	method := diffPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return &DiffResponse{Stream: resp.Body}, nil
}

func pluginChangesRequest(client *http.Client, req ChangesRequest) (*ChangesResponse, error) {
	method := changesPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp ChangesResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginApplyDiffRequest(client *http.Client, req ApplyDiffRequest) (*ApplyDiffResponse, error) {
	method := applyDiffPath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp ApplyDiffResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

func pluginDiffSizeRequest(client *http.Client, req DiffSizeRequest) (*DiffSizeResponse, error) {
	method := diffSizePath
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := client.Post("http://localhost"+method, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var vResp DiffSizeResponse
	err = json.NewDecoder(resp.Body).Decode(&vResp)
	if err != nil {
		return nil, err
	}

	return &vResp, nil
}

type testPlugin struct {
	init        int
	create      int
	remove      int
	get         int
	put         int
	exists      int
	status      int
	getMetadata int
	cleanup     int
	diff        int
	changes     int
	applyDiff   int
	diffSize    int
}

var _ Driver = &testPlugin{}

func (p *testPlugin) Init(string, map[string]string) error {
	p.init++
	return nil
}

func (p *testPlugin) Create(string, string) error {
	p.create++
	return nil
}

func (p *testPlugin) Remove(string) error {
	p.remove++
	return nil
}

func (p *testPlugin) Get(string, string) (string, error) {
	p.get++
	return "", nil
}

func (p *testPlugin) Put(string) error {
	p.put++
	return nil
}

func (p *testPlugin) Exists(string) bool {
	p.exists++
	return true
}

func (p *testPlugin) Status() [][2]string {
	p.status++
	return nil
}

func (p *testPlugin) GetMetadata(string) (map[string]string, error) {
	p.getMetadata++
	return nil, nil
}

func (p *testPlugin) Cleanup() error {
	p.cleanup++
	return nil
}

func (p *testPlugin) Diff(string, string) io.ReadCloser {
	p.diff++
	b := new(bytes.Buffer)
	x := ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	return x
}

func (p *testPlugin) Changes(string, string) ([]Change, error) {
	p.changes++
	return nil, nil
}

func (p *testPlugin) ApplyDiff(string, string, io.Reader) (int64, error) {
	p.applyDiff++
	return 42, nil
}

func (p *testPlugin) DiffSize(string, string) (int64, error) {
	p.diffSize++
	return 42, nil
}
