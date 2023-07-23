package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

type LinkerPod struct {
	ickcore.BareSnippet
	dom.UI

	IsTiles bool

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

func (mp *LinkerPod) SetTiles(f bool) {
	mp.IsTiles = f

	for _, l := range mp.LinksMap {
		l.DOM.SetClassIf(mp.IsTiles, "mr-4")
	}

	for _, p := range mp.MiniPodMap {
		p.DOM.SetClassIf(mp.IsTiles, "mr-4")
		p.OnOpenClose(false)
	}

	dom.Id("minipods").SetClassIf(mp.IsTiles, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
	dom.Id("nominipod").SetClassIf(mp.IsTiles, "is-flex is-flex-direction-row is-flex-wrap-wrap is-justify-content-flex-start is-align-content-flex-start is-align-items-flex-start")
}

/******************************************************************************/

func (mp *LinkerPod) BuildTag() ickcore.Tag {
	mp.Tag().SetTagName("layout")
	return *mp.Tag()
}

func (mp *LinkerPod) RenderContent(out io.Writer) error {
	ickcore.RenderString(out, `<div id="minipods"></div>`)
	ickcore.RenderString(out, `<div id="nominipod" class="pt-2"></div>`)
	return nil
}

/******************************************************************************/

// Mount inserts minipods and single cards. If at correspond to the id of a minipod then opens it
func (mp *LinkerPod) Mount(at string) error {
	at = strings.ToLower(at)

	// Insert Sorted MiniPodMap
	kmp := make([]string, 0, len(_lp.MiniPodMap))
	for k := range _lp.MiniPodMap {
		kmp = append(kmp, k)
	}
	sort.Slice(kmp, func(i, j int) bool {
		if _lp.MiniPodMap[kmp[i]].ABC == _lp.MiniPodMap[kmp[j]].ABC {
			return _lp.MiniPodMap[kmp[i]].Name < _lp.MiniPodMap[kmp[j]].Name
		}
		return _lp.MiniPodMap[kmp[i]].ABC < _lp.MiniPodMap[kmp[j]].ABC
	})
	emp := dom.Id("minipods")
	for _, k := range kmp {
		fmt.Printf("%+v --- %+v\n", _lp.MiniPodMap[k].Tag(), at)
		if strings.ToLower(_lp.MiniPodMap[k].Tag().Id()) == at {
			_lp.MiniPodMap[k].IsOpen = true
		}
		emp.InsertSnippet(dom.INSERT_LAST_CHILD, _lp.MiniPodMap[k])
	}

	// Insert Sorted LinksMap
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
