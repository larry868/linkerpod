package main

import (
	"io"
	"sort"
	"strings"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

type LAYOUT int

const (
	LAYOUT_TILES LAYOUT = iota
	LAYOUT_LIST
)

type LinkerPod struct {
	ickcore.BareSnippet
	dom.UI

	Layout LAYOUT

	LinksMap   map[string]*CardSnippet
	MiniPodMap map[string]*MiniPodSnippet
}

// Ensuring MiniPodSnippet implements the right interface
var _ dom.UIComposer = (*LinkerPod)(nil)

func NewLinkerPod() *LinkerPod {
	n := new(LinkerPod)
	n.LinksMap = make(map[string]*CardSnippet, 0)
	n.MiniPodMap = make(map[string]*MiniPodSnippet, 0)
	return n
}

func (lp *LinkerPod) SetLayout(layout LAYOUT) {
	lp.Layout = layout

	for _, l := range lp.LinksMap {
		l.DOM.SetClassIf(lp.Layout == LAYOUT_TILES, "mr-4")
	}

	for _, p := range lp.MiniPodMap {
		p.DOM.SetClassIf(lp.Layout == LAYOUT_TILES, "mr-4")
		p.OnOpenClose(false)
	}

	dom.Id("podlinks-header").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	dom.Id("minipods").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	dom.Id("podlinks-footer").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	dom.Id("podlinks-orphan").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	// browser.LocalStorage().Set("layout", strconv.Itoa(int(lp.Layout)))
}

/******************************************************************************/

func (lp *LinkerPod) BuildTag() ickcore.Tag {
	lp.Tag().SetTagName("layout")
	return *lp.Tag()
}

func (lp *LinkerPod) RenderContent(out io.Writer) error {
	ickcore.RenderString(out, `<div id="podlinks-header"></div>`)
	ickcore.RenderString(out, `<div id="minipods"></div>`)
	ickcore.RenderString(out, `<div id="podlinks-footer"></div>`)
	ickcore.RenderString(out, `<div id="podlinks-orphan"></div>`)
	return nil
}

/******************************************************************************/

// Mount inserts minipods and single cards. If at correspond to the id of a minipod then opens it
func (lp *LinkerPod) Mount(at string) error {
	at = strings.ToLower(at)

	// sort list of Minipods' ABCGroups
	xabcgs := make([]string, 0)
nextmp:
	for _, mp := range _gpod.MiniPodMap {
		g := mp.ABCGroup()
		for _, k := range xabcgs {
			if g == k {
				continue nextmp
			}
		}
		xabcgs = append(xabcgs, g)
	}
	sort.Strings(xabcgs)

	// Insert Sorted MiniPodMap by abc
	for _, xabcg := range xabcgs {

		mpking := make([]string, 0)
		for mpk, mp := range _gpod.MiniPodMap {
			if mp.ABCGroup() == xabcg {
				mpking = append(mpking, mpk)
			}
		}
		sort.Slice(mpking, func(i, j int) bool {
			if _gpod.MiniPodMap[mpking[i]].ABC == _gpod.MiniPodMap[mpking[j]].ABC {
				return _gpod.MiniPodMap[mpking[i]].Name < _gpod.MiniPodMap[mpking[j]].Name
			}
			return _gpod.MiniPodMap[mpking[i]].ABC < _gpod.MiniPodMap[mpking[j]].ABC
		})

		mpg := ick.Elem("div", `class="pb-3"`)
		x := xabcg[0]
		switch x {
		case '1':
			mpg.Tag().SetId("minipod.header")
		case '3':
			mpg.Tag().SetId("minipod.footer")
		default:
			if len(xabcg) > 1 {
				mpg.Tag().SetId("minipod." + xabcg[1:len(xabcg)])
			} else {
				mpg.Tag().SetId("minipod._void")
			}
		}

		for _, k := range mpking {
			if strings.ToLower(_gpod.MiniPodMap[k].Tag().Id()) == at {
				_gpod.MiniPodMap[k].IsOpen = true
			}
			mpg.Append(_gpod.MiniPodMap[k])
		}

		switch x {
		case '1':
			dom.Id("podlinks-header").InsertSnippet(dom.INSERT_LAST_CHILD, mpg)
		case '3':
			dom.Id("podlinks-footer").InsertSnippet(dom.INSERT_LAST_CHILD, mpg)
		default:
			dom.Id("minipods").InsertSnippet(dom.INSERT_LAST_CHILD, mpg)
		}
	}

	// Insert orphan links (not in any minipod)
	kl := make([]string, 0, len(_gpod.LinksMap))
	for k, l := range _gpod.LinksMap {
		if l.InMiniPods == 0 {
			kl = append(kl, k)
		}
	}
	if len(kl) > 0 {
		eno := dom.Id("podlinks-orphan")
		eno.InsertRawHTML(dom.INSERT_LAST_CHILD, `<hr class="mt-2"/>`)

		sort.Strings(kl)
		for _, k := range kl {
			_gpod.LinksMap[k].Expand(true)
			eno.InsertSnippet(dom.INSERT_LAST_CHILD, _gpod.LinksMap[k])
		}
	}

	return nil
}
