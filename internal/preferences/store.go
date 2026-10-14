package preferences

import (
	"encoding/json"
	"os"
	"path/filepath"

	"newtofile/internal/appui"
)

func Load(path string, defaults appui.Preferences) (appui.Preferences, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return defaults, nil
	}
	if err != nil {
		return defaults, err
	}
	result := defaults
	if err := json.Unmarshal(data, &result); err != nil {
		return defaults, err
	}
	if !filepath.IsAbs(result.Directory) {
		result.Directory = defaults.Directory
	}
	if result.Appearance < 0 || result.Appearance > 2 {
		result.Appearance = 0
	}
	return result, nil
}

func Save(path string, preferences appui.Preferences) error {
	data, err := json.MarshalIndent(preferences, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".preferences-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
