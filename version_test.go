package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGetVersionWorks(t *testing.T) {
	buf, err := ioutil.ReadFile(".goxc.json")
	if err != nil {
		panic("Error reading .goxc.json")
	}
	config := &GoxcConfig{}
	err = json.Unmarshal(buf, config)
	if err != nil {
		panic("Error parsing .goxc.json")
	}
	if len(config.PackageVersion) == 0 {
		t.Errorf("Expected PackageVersion to have nonzero length")
	}
	version := GetVersion()
	if version != config.PackageVersion {
		t.Errorf(
			"Expected GetVersion() to return '%s', not '%s'. Consider re-running 'go generate'.",
			config.PackageVersion,
			version,
		)
	}
}
