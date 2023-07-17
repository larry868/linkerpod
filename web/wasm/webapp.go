// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

// this main package contains the web assembly source code.
// It's compiled into a '.wasm' file with "GOOS=js GOARCH=wasm go build -o ../webapp/main.wasm"
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/icecake-framework/icecake/pkg/dom"
	"gopkg.in/yaml.v3"
)

// The main func is required by the wasm GO builder.
// outputs will appears in the console of the browser
func main() {
	// TODO: handle id="wasm-status" element in the html in case wasm not loading (see. wasm_loader.js)

	c := make(chan struct{})
	fmt.Println("Go/WASM loaded.")

	u, errp := url.Parse(dom.Doc().Body().BaseURI())
	if errp != nil {
		fmt.Println("parsing error:", errp.Error())
	}
	u = u.JoinPath("./linkerpod.yaml")

	lmap := make(map[string]*LinkCardSnippet)
	ys, err := DownloadYaml(u.String())
	if err != nil {
		// TODO: hendle error with ick.console and display a message to the user
		fmt.Println(err)
	} else {
		for k, l := range ys.Links {
			lmap[k] = LinkCard(k).ParseHRef(l.Link)
		}
	}

	app := dom.Id("app")
	for _, l := range lmap {
		app.InsertSnippet(dom.INSERT_LAST_CHILD, l)
	}

	// let's go
	fmt.Println("Go/WASM listening browser events")
	<-c
}

type YamlLinkEntry struct {
	Link string `yaml:"link"`
}

type YamlStruct struct {
	Links map[string]YamlLinkEntry `yaml:"links"`
}

func DownloadYaml(url string) (*YamlStruct, error) {
	fmt.Println("DownloadYaml:", url)

	buf := &bytes.Buffer{}
	err := DownloadFile(buf, url)
	if err != nil {
		return nil, fmt.Errorf("DownloadYaml: %+w", err)
	}

	ys := &YamlStruct{}
	err = yaml.Unmarshal(buf.Bytes(), ys)
	if err != nil {
		return nil, fmt.Errorf("DownloadYaml: %+w", err)
	}

	//fmt.Println("YAML=> %+v", ys)

	return ys, nil
}

func DownloadFile(w io.Writer, url string) error {
	// Get the data
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Permissions-Policy", "interest-cohort=()")

	resp, err := client.Do(req)
	// resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to the writer
	_, err = io.Copy(w, resp.Body)
	return err
}
