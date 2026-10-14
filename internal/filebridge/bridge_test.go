package filebridge

import (
	"net/url"
	"testing"
)

func TestFinderURLPreservesPath(t *testing.T) {
	path := "/Users/test/新文件 & 100% # + ?"
	q := url.Values{"type": {"txt"}, "directory": {path}}
	r, err := ParseURL("newtofile://create?" + q.Encode())
	if err != nil || r.Directory != path || r.Kind != "txt" {
		t.Fatalf("request=%+v err=%v", r, err)
	}
}

func TestInvalidRequestsRejected(t *testing.T) {
	for _, raw := range []string{"https://create?type=txt&directory=/tmp", "newtofile://create?type=exe&directory=/tmp", "newtofile://create?type=txt&directory=relative", "newtofile://create?type=txt&directory=/tmp%00a", "newtofile://create?type=txt&directory=%zz", "newtofile://wrong"} {
		if _, err := ParseURL(raw); err == nil {
			t.Fatal(raw)
		}
	}
}

func TestSettingsURL(t *testing.T) {
	r, err := ParseURL("newtofile://settings")
	if err != nil || r.Operation != "settings" {
		t.Fatalf("request=%+v err=%v", r, err)
	}
}
