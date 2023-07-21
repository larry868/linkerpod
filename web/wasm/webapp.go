// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

// this main package contains the web assembly source code.
// It's compiled into a '.wasm' file with "GOOS=js GOARCH=wasm go build -o ../webapp/main.wasm"
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/icecake-framework/icecake/pkg/console"
	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
	"gopkg.in/yaml.v3"
)

// The main func is required by the wasm GO builder.
// outputs will appears in the console of the browser
func main() {
	c := make(chan struct{})
	fmt.Println("Go/WASM loaded and running...")
	start := time.Now()

	lmap, err := DownloadLinks()
	if err != nil {
		console.Errorf(err.Error())
		dom.Id("webapp").InsertSnippet(dom.INSERT_LAST_CHILD, ick.Message(ickcore.ToHTML("Unable to load this linkerpod.")).SetColor(ick.COLOR_DANGER))
	} else {
		app := dom.Id("webapp")
		for _, l := range lmap {
			app.InsertSnippet(dom.INSERT_LAST_CHILD, l)
		}
		fmt.Printf("Linkerpod loaded in %v\n", time.Since(start).Round(time.Millisecond))
	}

	// let's go
	fmt.Println("Go/WASM ready and listening browser events")
	<-c
}

type YamlLinkEntry struct {
	Link string `yaml:"link"`
}

type YamlStruct struct {
	Links map[string]YamlLinkEntry `yaml:"links"`
}

func DownloadLinks() (map[string]*LinkCardSnippet, error) {
	lmap := make(map[string]*LinkCardSnippet)
	ys, err := DownloadYaml("linkerpod.yaml")
	if err != nil {
		return nil, fmt.Errorf("DownloadYaml: %w", err)
	}

	for k, l := range ys.Links {
		if l.Link == "" {
			console.Warnf("[%s] missing link", k)
		} else {
			lmap[k] = LinkCard(k).ParseHRef(l.Link)
		}
	}

	if len(lmap) == 0 {
		return nil, fmt.Errorf("empty pod")
	}

	return lmap, nil
}

func DownloadYaml(url string) (*YamlStruct, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Permissions-Policy", "interest-cohort=()")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	// Write the body to the writer
	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, resp.Body)
	if err == nil && n == 0 {
		err = errors.New("empty yaml file")
	}
	if err != nil {
		return nil, err
	}

	ys := &YamlStruct{}
	err = yaml.Unmarshal(buf.Bytes(), ys)
	if err != nil {
		return nil, err
	}

	return ys, nil
}
