// Copyright 2022 by lolorenzo77. All rights reserved.
// Use of this source code is governed by MIT licence that can be found in the LICENSE file.

// this main package contains the web assembly source code.
// It's compiled into a '.wasm' file with "GOOS=js GOARCH=wasm go build -o ../webapp/main.wasm"
package main

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/lolorenzo777/linkerpod/pkg/yamlpod"

	"github.com/icecake-framework/icecake/pkg/console"
	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ick/ickui"
	"github.com/icecake-framework/icecake/pkg/ickcore"
	"github.com/icecake-framework/icecake/pkg/namingpattern"
)

const (
	ICON_TILES string = "bi bi-columns-gap"
	ICON_LIST  string = "bi bi-view-list"
)

const (
	_setupfile        string = "linkerpod.yaml"
	_setupdefaultfile string = "linkerpod_default.yaml"
)

var (
	_btnLayout = ickui.Button("Tiles", "").
			SetId("btnlayout").
			SetColor(ick.COLOR_PRIMARY).
			SetOutlined(true).
			SetDisabled(true).
			SetIcon(*ick.Icon(ICON_TILES), false).
			SetSize(ick.SIZE_SMALL)

	_lp LinkerPod
)

// The main func is required by the wasm GO builder.
// outputs will appears in the console of the browser
func main() {
	c := make(chan struct{})
	fmt.Println("Go/WASM loaded and running...")
	// verbose.IsOn = true
	// verbose.IsDebugging = true

	var err error

	start := time.Now()

	// title
	u := dom.Doc().Body().BaseURI()
	dom.Id("title").InsertText(dom.INSERT_BODY, path.Dir(u.String()))

	yaml := _setupfile

	// extract query if any
	// TODO: encode the query
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
	dom.Id("subtitle").InsertText(dom.INSERT_BODY, yaml[:len(yaml)-5])

	if err == nil {
		_lp, err = DownloadData(yaml)
		if yaml == _setupfile && errors.Is(err, yamlpod.ErrGetYamlFile) {
			_lp, err = DownloadData(_setupdefaultfile)
		}
	}

	if err == nil {
		dom.Id("webapp").InsertSnippet(dom.INSERT_BODY, &_lp)
		_lp.Mount(u.Fragment)

		_btnLayout.OnClick = OnToggleLayout
		dom.Id("webapp").InsertSnippet(dom.INSERT_FIRST_CHILD,
			ick.Elem("div", `class="level is-mobile"`,
				ick.Elem("div", `class="level-left"`,
					ick.Elem("div", `class="level-item"`, _btnLayout))))

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

func OnToggleLayout() {
	if _lp.IsTiles {
		_btnLayout.Title = "Tiles"
		_btnLayout.OpeningIcon.Key = ICON_TILES
		_lp.SetTiles(false)
	} else {
		_btnLayout.Title = "List"
		_btnLayout.OpeningIcon.Key = ICON_LIST
		_lp.SetTiles(true)
	}
	_btnLayout.RefreshContent(_btnLayout)
	_btnLayout.DOM.Blur()
}

/******************************************************************************/

func DownloadData(yaml string) (LinkerPod, error) {

	lp := NewLinkerPod()

	// download yaml file
	ys, err := yamlpod.DownloadYaml(yaml)
	if err != nil {
		return *lp, fmt.Errorf("DownloadData: %w", err)
	}

	// parse lp.MiniPodMap
	for ympk, ymp := range ys.MiniPods {
		if ymp.Name == "" {
			ymp.Name = ympk
		}
		ympk = "mp-" + namingpattern.MakeValidName(strings.ToLower(ympk))
		if _, found := lp.MiniPodMap[ympk]; found {
			console.Warnf("DownloadData.minipods: duplicate minipod id %q", ympk)
			continue
		}
		mp := MiniPod(ympk, ymp.Name, ymp.Icon, strings.ToLower(ymp.ABC))
		lp.MiniPodMap[ympk] = mp
	}

	// parse lp.LinksMap
	for ylnkk, ylnk := range ys.Links {
		if ylnk.Link == "" {
			console.Warnf("DownloadData.links: missing link %q", ylnkk)
			continue
		}

		if ylnk.Name == "" {
			ylnk.Name = ylnkk
		}
		lnkkey := "lnk-" + namingpattern.MakeValidName(strings.ToLower(ylnkk))
		if _, found := lp.LinksMap[lnkkey]; found {
			console.Warnf("DownloadData.links: duplicate %q", lnkkey)
			continue
		}

		lnk := Card(lnkkey, ylnk.Name).ParseHRef(ylnk.Link)
		lnk.SetIcon(ylnk.Icon)
		lp.LinksMap[lnkkey] = lnk

		// insert card in Minipods
		for _, mpinlnk := range ylnk.Minipods {
			mpkey := "mp-" + namingpattern.MakeValidName(strings.ToLower(mpinlnk.MinipodKey))
			mp, found := lp.MiniPodMap[mpkey]
			if !found {
				console.Warnf("DownloadData.links: minipod %q not found", mpinlnk.MinipodKey)
				continue
			}

			lp.LinksMap[lnkkey].InMiniPods++
			inlnkkey := mp.Tag().SubId(lnkkey)
			lnkin := *lnk
			lnkin.Tag().SetId(inlnkkey)
			mp.InsertCard(lnkin, mpinlnk.ABC)
		}
	}

	if len(lp.LinksMap) == 0 {
		return *lp, fmt.Errorf("empty pod")
	}

	return *lp, nil
}
