package profile

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Profile, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, fmt.Errorf("read candidate profile %q: %w", path, err)
	}

	var candidate Profile
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	if err := decoder.Decode(&candidate); err != nil {
		return Profile{}, fmt.Errorf("decode candidate profile %q: %w", path, err)
	}
	if err := candidate.Validate(); err != nil {
		return Profile{}, fmt.Errorf("validate candidate profile %q: %w", path, err)
	}
	return candidate, nil
}
