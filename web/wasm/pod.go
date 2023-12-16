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

// LinkerPod snippet embedding all the page content
type LinkerPod struct {
	ickcore.BareSnippet
	dom.UI

	Layout LAYOUT

	MiniPodMap    map[string]*MiniPodSnippet // mini pods containing links
	SingleLinkMap map[string]*CardSnippet    // standalone links
}

// Ensuring LinkerPod implements the right interface
var _ dom.UIComposer = (*LinkerPod)(nil)

// LinkerPod factory
func NewLinkerPod() *LinkerPod {
	n := new(LinkerPod)
	n.MiniPodMap = make(map[string]*MiniPodSnippet, 0)
	n.SingleLinkMap = make(map[string]*CardSnippet, 0)
	return n
}

func (lp *LinkerPod) SetLayout(layout LAYOUT) {
	lp.Layout = layout

	for _, p := range lp.MiniPodMap {
		p.DOM.SetClassIf(lp.Layout == LAYOUT_TILES, "mr-4")
		p.OnOpenClose(false)
	}
	dom.Id("minipods").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")

	for _, l := range lp.SingleLinkMap {
		l.DOM.SetClassIf(lp.Layout == LAYOUT_TILES, "mr-4")
	}
	dom.Id("singlelinks").SetClassIf(lp.Layout == LAYOUT_TILES, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
}

/******************************************************************************/

func (lp *LinkerPod) BuildTag() ickcore.Tag {
	lp.Tag().SetTagName("layout")
	return *lp.Tag()
}

func (lp *LinkerPod) RenderContent(out io.Writer) error {
	ickcore.RenderString(out, `<div id="minipods"></div>`)
	ickcore.RenderString(out, `<div id="singlelinks"></div>`)
	return nil
}

/******************************************************************************/

// Mount inserts minipods and single cards. If at correspond to the id of a minipod then opens it
func (lp *LinkerPod) Mount(at string) error {
	at = strings.ToLower(at)

	// sort list of Minipods' ABCGroups
	mpkeys := make([]string, 0)
	for mpkey := range lp.MiniPodMap {
		mpkeys = append(mpkeys, mpkey)
	}
	sort.Strings(mpkeys)

	// Insert sorted MiniPods
	for _, mpkey := range mpkeys {

		if lp.MiniPodMap[mpkey].Separator == "splitter" {
			dom.Id("minipods").InsertRawHTML(dom.INSERT_LAST_CHILD, `<hr class="my-4"/>`)
		}

		if !lp.MiniPodMap[mpkey].IsSeparatorOnly() {
			if mpkey == at {
				lp.MiniPodMap[mpkey].IsOpen = true
			}

			mpg := ick.Elem("div", ``)
			if lp.MiniPodMap[mpkey].Separator == "splitter" {
				mpg.Tag().AddClass("pt-1")
			} else if lp.MiniPodMap[mpkey].Separator != "" {
				mpg.Tag().AddClass("pt-3")
			}
			mpg.Tag().SetId("minipod." + mpkey)
			mpg.Append(lp.MiniPodMap[mpkey])
			dom.Id("minipods").InsertSnippet(dom.INSERT_LAST_CHILD, mpg)
		}
	}

	// Insert single links
	slkeys := make([]string, 0, len(_glp.SingleLinkMap))
	for slkey, _ := range _glp.SingleLinkMap {
		slkeys = append(slkeys, slkey)
	}
	if len(slkeys) > 0 {
		eno := dom.Id("singlelinks")
		eno.InsertRawHTML(dom.INSERT_LAST_CHILD, `<hr class="my-4"/>`)

		sort.Strings(slkeys)
		for _, k := range slkeys {
			_glp.SingleLinkMap[k].Expand(true)
			eno.InsertSnippet(dom.INSERT_LAST_CHILD, _glp.SingleLinkMap[k])
		}
	}

	return nil
}
