package graphdriver

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/go-connections/sockets"
)

const url = "http://localhost"

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
	init, err := CallInit(url, client, InitRequest{Home: "foo"})
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
	create, err := CallCreate(url, client, CreateRequest{ID: "foo", Parent: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if create.Err != "" {
		t.Fatalf("got error creating graph drivers: %v", create.Err)
	}
	if p.create != 1 {
		t.Fatalf("expected create 1, got %d", p.create)
	}

	// CreateReadWrite
	createReadWrite, err := CallCreateReadWrite(url, client,
		CreateRequest{ID: "foo", Parent: "bar", MountLabel: "toto"})
	if err != nil {
		t.Fatal(err)
	}
	if create.Err != "" {
		t.Fatalf("got error creating read-write graph drivers: %v", createReadWrite.Err)
	}
	if p.createRW != 1 {
		t.Fatalf("expected createReadWrite 1, got %d", p.createRW)
	}

	// Remove
	remove, err := CallRemove(url, client, RemoveRequest{ID: "foo"})
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
	get, err := CallGet(url, client, GetRequest{ID: "foo", MountLabel: "bar"})
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
	put, err := CallPut(url, client, PutRequest{ID: "foo"})
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
	exists, err := CallExists(url, client, ExistsRequest{ID: "foo"})
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
	status, err := CallStatus(url, client, StatusRequest{})
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
	getMetadata, err := CallGetMetadata(url, client, GetMetadataRequest{ID: "foo"})
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
	cleanup, err := CallCleanup(url, client, CleanupRequest{})
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
	_, err = CallDiff(url, client, DiffRequest{ID: "foo", Parent: "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if p.diff != 1 {
		t.Fatalf("expected diff 1, got %d", p.diff)
	}

	// Changes
	changes, err := CallChanges(url, client, ChangesRequest{ID: "foo", Parent: "bar"})
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
	applyDiff, err := CallApplyDiff(url, client,
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
	diffSize, err := CallDiffSize(url, client, DiffSizeRequest{ID: "foo", Parent: "bar"})
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

type testPlugin struct {
	init        int
	create      int
	createRW    int
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

func (p *testPlugin) Init(string, []string) error {
	p.init++
	return nil
}

func (p *testPlugin) Create(string, string) error {
	p.create++
	return nil
}

func (p *testPlugin) CreateReadWrite(string, string) error {
	p.createRW++
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
