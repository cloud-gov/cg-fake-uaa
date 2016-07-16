package main

import (
	"encoding/json"
)

type GoxcConfig struct {
	PackageVersion string
}

func GetVersion() string {
	config := &GoxcConfig{}
	err := json.Unmarshal(GetAsset(".goxc.json"), config)
	if err != nil {
		panic("Unable to parse .goxc.json!")
	}
	return config.PackageVersion
}
