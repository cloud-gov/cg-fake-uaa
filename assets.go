//go:generate go-bindata data/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func readFile(path string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("Error reading file %s: %v", path, err)
	}
	return buf, err
}

func GetAsset(path string) ([]byte, error) {
	debug := os.Getenv("FAKECLOUDGOV_DEBUG")
	if (debug == "") {
		return Asset(path)
	} else {
		return readFile(path)
	}
}
