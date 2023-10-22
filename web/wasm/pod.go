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

	dom.Id("minipods").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	dom.Id("nominipod").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	// browser.LocalStorage().Set("layout", strconv.Itoa(int(lp.Layout)))
}

/******************************************************************************/

func (lp *LinkerPod) BuildTag() ickcore.Tag {
	lp.Tag().SetTagName("layout")
	return *lp.Tag()
}

func (lp *LinkerPod) RenderContent(out io.Writer) error {
	ickcore.RenderString(out, `<div id="minipods"></div>`)
	ickcore.RenderString(out, `<div id="nominipod"></div>`)
	return nil
}

/******************************************************************************/

// Mount inserts minipods and single cards. If at correspond to the id of a minipod then opens it
func (lp *LinkerPod) Mount(at string) error {
	at = strings.ToLower(at)

	// get sorted list of ABCGroups
	abcgs := make([]string, 0)
nextmp:
	for _, mp := range _lp.MiniPodMap {
		g := mp.ABCGroup()
		for _, k := range abcgs {
			if g == k {
				continue nextmp
			}
		}
		abcgs = append(abcgs, g)
	}
	sort.Strings(abcgs)

	// Insert Sorted MiniPodMap by abc
	for _, abcg := range abcgs {

		mpking := make([]string, 0)
		for mpk, mp := range _lp.MiniPodMap {
			if mp.ABCGroup() == abcg {
				mpking = append(mpking, mpk)
			}
		}
		sort.Slice(mpking, func(i, j int) bool {
			if _lp.MiniPodMap[mpking[i]].ABC == _lp.MiniPodMap[mpking[j]].ABC {
				return _lp.MiniPodMap[mpking[i]].Name < _lp.MiniPodMap[mpking[j]].Name
			}
			return _lp.MiniPodMap[mpking[i]].ABC < _lp.MiniPodMap[mpking[j]].ABC
		})

		mpg := ick.Elem("div", `class="pb-3"`)
		mpg.Tag().SetId("minipods." + abcg)
		for _, k := range mpking {
			if strings.ToLower(_lp.MiniPodMap[k].Tag().Id()) == at {
				_lp.MiniPodMap[k].IsOpen = true
			}
			mpg.Append(_lp.MiniPodMap[k])
		}
		dom.Id("minipods").InsertSnippet(dom.INSERT_LAST_CHILD, mpg)

	}
	// Insert Sorted single LinksMap
	kl := make([]string, 0, len(_lp.LinksMap))
	for k, l := range _lp.LinksMap {
		if l.InMiniPods == 0 {
			kl = append(kl, k)
		}
	}
	sort.Strings(kl)
	eno := dom.Id("nominipod")
	for _, k := range kl {
		_lp.LinksMap[k].Expand(true)
		eno.InsertSnippet(dom.INSERT_LAST_CHILD, _lp.LinksMap[k])
	}

	return nil
}
