package repotrack

import (
	"fmt"
	"testing"
)

func TestPopulate(t *testing.T) {

	protocol := "https"
	repoName := "myrepo"

	rt := NewRepoTrack()
	rt.URL = fmt.Sprintf("%s://github.com/chrishiestand/%s.git", protocol, repoName)
	rt = Populate(rt)

	if rt.Protocol != protocol {
		t.Fatalf("repotrack.Populate() did not return correct protocol, expected '%s' got '%s'", protocol, rt.Protocol)
	}

	if rt.Name != repoName {
		t.Fatalf("repotrack.Populate() did not return correct name, expected '%s' got '%s'", repoName, rt.Name)
	}
}
