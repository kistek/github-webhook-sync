package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAliveHandler(t *testing.T) {

	req := httptest.NewRequest("GET", "/_alive", nil)
	w := httptest.NewRecorder()

	//http.ResponseWriter, *http.Request
	aliveHandler()(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Logf("got:\n%#v\n", resp.StatusCode)
		t.Logf("expected:\n%#v\n", http.StatusOK)
		t.Fatalf("/_alive wasn't returned as expected.\nGot:\n%d\nExpected:\n%d",
			resp.StatusCode, http.StatusOK)
	}

}

func TestReadyHandler200(t *testing.T) {

	req := httptest.NewRequest("GET", "/_ready", nil)
	w := httptest.NewRecorder()

	//http.ResponseWriter, *http.Request
	readyHandler("/tmp")(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("/_alive wasn't returned as expected.\nGot:\n%d\nExpected:\n%d",
			resp.StatusCode, http.StatusOK)
	}
}

func TestReadyHandler503(t *testing.T) {

	req := httptest.NewRequest("GET", "/_ready", nil)
	w := httptest.NewRecorder()

	//http.ResponseWriter, *http.Request
	readyHandler("/doesnotexist")(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("/_alive wasn't returned as expected.\nGot:\n%d\nExpected:\n%d",
			resp.StatusCode, http.StatusServiceUnavailable)
	}
}

func TestWebhookHandler(t *testing.T) {

	//TODO add headers - ping
	// X-GitHub-Delivery: 17d6eeb0-0ddf-11e8-8b63-2d493c6c1b09
	// X-GitHub-Event: ping
	// X-Hub-Signature: sha1=b58ed79d892ebfc3fe8e7dd733aa784dd67c09a4

	// push headers
	// X-GitHub-Delivery: c2690c8c-0ddf-11e8-83b2-75793281676e
	// X-GitHub-Event: push
	// X-Hub-Signature: sha1=24985b34e3bfe4c49a1ba846c090bddd30552905

	pushFixture := "testdata/example-firstpush.json"
	fixtureSourcePath := "testdata/repo_source"
	pushSha1 := "c18ad04687e0f9651d473212ab8fa8d6643f7c58"
	repoRootPath := "/tmp/"

	pushJSON, err := ioutil.ReadFile(pushFixture)
	if err != nil {
		t.Fatalf("could not read testdata file: %s", pushFixture)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(pushJSON, &raw); err != nil {
		t.Fatalf("could not unmarshal json: %s", string(pushJSON))
	}

	commitID, ok := raw["after"].(string)

	if ok == false {
		t.Fatalf("could not read commit string from testdata file: %s", pushFixture)
	}

	fmt.Print("commitID: " + commitID + "\n")

	jsonReader := bytes.NewReader(pushJSON)

	req := httptest.NewRequest("GET", "/webhook", jsonReader)
	w := httptest.NewRecorder()

	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-GitHub-Delivery", "bar")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature", "sha1="+pushSha1)

	rts, err2 := ConfigPathToRepoTracks(fixtureSourcePath)

	if err2 != nil {
		t.Fatalf("could not read repo metadata from path: %s\n%v", fixtureSourcePath, err2)
	}

	webhookHandler(rts, repoRootPath)(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("webhook wasn't handled as expected.\nGot:\n%d\nExpected:\n%d",
			resp.StatusCode, http.StatusServiceUnavailable)
	}

	//TODO check output?
	// repoName := "foo"
	// if err := unix.Access(repoRootPath+repoName, unix.R_OK); err != nil {
	// 	t.Fatalf("failed to clone git repo\n%s", repoRootPath+repoName)
	// }

}

func TestWebhookHandlerInvalidContentType(t *testing.T) {

	pushFixture := "testdata/example-push.json"
	fixtureSourcePath := "testdata/repo_source"

	pushJSON, err := ioutil.ReadFile(pushFixture)
	if err != nil {
		t.Fatalf("could not read testdata file: %s", pushFixture)
	}

	jsonReader := bytes.NewReader(pushJSON)

	req := httptest.NewRequest("GET", "/webhook", jsonReader)
	w := httptest.NewRecorder()

	req.Header.Set("content-type", "applicaton/x-www-form-urlencoded")
	req.Header.Set("X-GitHub-Delivery", "bar")
	req.Header.Set("X-GitHub-Event", "push")
	req.Header.Set("X-Hub-Signature", "sha1=24985b34e3bfe4c49a1ba846c090bddd30552905")

	repoRootPath := "/tmp/"

	rts, err := ConfigPathToRepoTracks(fixtureSourcePath)

	if err != nil {
		t.Fatalf("could not read repo metadata from path: %s", fixtureSourcePath)
	}

	webhookHandler(rts, repoRootPath)(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("/webhook wasn't returned as expected.\nGot:\n%d\nExpected:\n%d",
			resp.StatusCode, http.StatusBadRequest)
	}

}
