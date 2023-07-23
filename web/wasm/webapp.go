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
	"path"
	"strings"
	"time"

	"github.com/icecake-framework/icecake/pkg/console"
	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ick/ickui"
	"github.com/icecake-framework/icecake/pkg/ickcore"
	"github.com/lolorenzo777/verbose"
	"gopkg.in/yaml.v3"
)

var (
	_btnShrink = ickui.Button("Shrink", "").
			SetId("btnshrink").
			SetColor(ick.COLOR_PRIMARY).
			SetOutlined(true).
			SetDisabled(true).
			SetSize(ick.SIZE_SMALL)
	_btnLayout = ickui.Button("Tiles", "").
			SetId("btnlayout").
			SetColor(ick.COLOR_PRIMARY).
			SetOutlined(true).
			SetDisabled(true).
			SetSize(ick.SIZE_SMALL)

	_lp LinkerPod
)

// The main func is required by the wasm GO builder.
// outputs will appears in the console of the browser
func main() {
	c := make(chan struct{})
	fmt.Println("Go/WASM loaded and running...")
	verbose.IsOn = true
	verbose.IsDebugging = true

	var err error

	start := time.Now()

	// extract query if any
	// TODO: encode the query
	yaml := "linkerpod.yaml"
	u := dom.Doc().Body().BaseURI()
	q := u.Query()
	if len(q) > 0 {
		pod := q["pod"]
		if len(pod) > 0 && pod[0] != "" {
			qbase := path.Base(pod[0])
			qext := path.Ext(pod[0])
			if (qext != "" && qext != ".yaml") || (qbase != pod[0]) {
				err = errors.New("malformed query")
			} else if qext == ".yaml" {
				yaml = qbase
			} else {
				yaml = qbase + ".yaml"
			}
		}
	}

	if err == nil {
		_lp, err = DownloadData(yaml)
	}

	if err == nil {
		dom.Id("webapp").InsertSnippet(dom.INSERT_BODY, &_lp)
		_lp.Mount()

		_btnShrink.OnClick = OnToggleShrink
		_btnLayout.OnClick = OnToggleLayout
		dom.Id("webapp").InsertSnippet(dom.INSERT_FIRST_CHILD,
			ick.Elem("div", `class="level"`,
				ick.Elem("div", `class="level-left"`,
					ick.Elem("div", `class="level-item"`, _btnShrink),
					ick.Elem("div", `class="level-item"`, _btnLayout))))

		_btnShrink.SetDisabled(false)
		_btnLayout.SetDisabled(false)
		fmt.Printf("Linkerpod loaded and displayed in %v\n", time.Since(start).Round(time.Millisecond))
	}

	if err != nil {
		console.Errorf(err.Error())
		dom.Id("webapp").InsertSnippet(dom.INSERT_BODY, ick.Message(ickcore.ToHTML("Unable to load this linkerpod.")).SetColor(ick.COLOR_DANGER))
	}

	// let's go
	fmt.Println("Go/WASM ready and listening browser events")
	<-c
}

func OnToggleShrink() {
	if _lp.IsShrunk {
		_btnShrink.SetTitle("Shrink")
		_lp.SetShrunk(false)
	} else {
		_btnShrink.SetTitle("Expand")
		_lp.SetShrunk(true)
	}
	_btnShrink.DOM.Blur()
}

func OnToggleLayout() {
	if _lp.IsTiles {
		_btnLayout.SetTitle("Tiles")
		_lp.SetTiles(false)
	} else {
		_btnLayout.SetTitle("List")
		_lp.SetTiles(true)
	}
	_btnLayout.DOM.Blur()
}

/******************************************************************************/

type YamlMiniPod struct {
	Icon  string     `yaml:"icon"`
	Links []YamlLink `yaml:"links"`
}

type YamlLink struct {
	Id  string `yaml:"id"`
	ABC string `yaml:"abc"`
}

type YamlLinkEntry struct {
	Link string `yaml:"link"`
	ABC  string `yaml:"abc"`
}

type YamlStruct struct {
	Links    map[string]YamlLinkEntry `yaml:"links"`
	MiniPods map[string]YamlMiniPod   `yaml:"minipods"`
}

func DownloadData(yaml string) (LinkerPod, error) {

	lp := NewLinkerPod()

	// download yaml file
	ys, err := DownloadYaml(yaml)
	if err != nil {
		return *lp, fmt.Errorf("DownloadData: %w", err)
	}

	// parse lp.LinksMap
	for k, v := range ys.Links {
		if v.Link == "" {
			console.Warnf("DownloadData.links: [%s] missing link id", k)
		} else {
			c := Card(k).ParseHRef(v.Link)
			if abc := strings.ToLower(strings.Trim(v.ABC, " ")); abc != "" {
				c.ABC = abc
			}
			lp.LinksMap[k] = c
		}
	}

	if len(lp.LinksMap) == 0 {
		return *lp, fmt.Errorf("empty pod")
	}

	// parse lp.MiniPodMap
	for k, mp := range ys.MiniPods {
		lp.MiniPodMap[k] = MiniPod(k, mp.Icon)
		for _, mpl := range mp.Links {
			if c, found := lp.LinksMap[mpl.Id]; found {
				lp.MiniPodMap[k].InsertCard(*c, strings.ToLower(strings.Trim(mpl.ABC, " ")))
				c.InMiniPods += 1
			} else {
				console.Warnf("DownloadData.minipods: link %q not referenced", mpl)
			}
		}
	}

	return *lp, nil
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
