package preferences

import (
	"newtofile/internal/appui"
	"os"
	"path/filepath"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings/preferences.json")
	want := appui.Preferences{Directory: "/tmp/中文 & files", OpenAfter: true, Appearance: 2}
	if err := Save(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path, appui.Preferences{})
	if err != nil || got != want {
		t.Fatalf("got=%+v err=%v", got, err)
	}
	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("permissions=%v", info.Mode())
	}
}

func TestInvalidPreferenceValuesFallBack(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preferences.json")
	if err := os.WriteFile(path, []byte(`{"directory":"relative","appearance":100}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path, appui.Preferences{Directory: "/tmp"})
	if err != nil || got.Directory != "/tmp" || got.Appearance != 0 {
		t.Fatalf("got=%+v err=%v", got, err)
	}
}
