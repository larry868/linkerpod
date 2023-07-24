package main

import (
	"io"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/event"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ick/ickui"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

type MiniPodSnippet struct {
	ickcore.BareSnippet
	dom.UI

	Name string // mini pod name
	Icon ick.ICKIcon
	Body []*CardSnippet // rendered as <div class="card-content">

	IsExpanded bool
	IsOpen     bool
	HasMore    int
	ABC        string
}

// Ensuring MiniPodSnippet implements the right interface
var _ dom.UIComposer = (*MiniPodSnippet)(nil)

func MiniPod(key string, name string, iconkey string, abc string) *MiniPodSnippet {
	n := new(MiniPodSnippet)
	n.Body = make([]*CardSnippet, 0)
	n.Tag().SetId(key)
	n.Name = name
	if iconkey != "" {
		n.Icon.Key = iconkey
	} else {
		n.Icon.Key = "bi bi-dot"
	}
	n.ABC = abc
	return n
}

// InsertCard inserts c card to the minipod according to abc position.
// If abc parameter is empty then the abc position of the card itself is used.
// If abc is greather or equal than c then the class "more" is added
func (mp *MiniPodSnippet) InsertCard(c CardSnippet, abc string) {
	if abc != "" {
		c.ABC = abc
	}
	c.Tag().SetClassIf(c.ABC >= "c", "more is-hidden")
	if c.ABC >= "c" {
		mp.HasMore++
	}

	for i, cinbody := range mp.Body {
		if c.ABC < cinbody.ABC || (c.ABC == cinbody.ABC && c.Name < cinbody.Name) {
			mp.Body = append(mp.Body[:i+1], mp.Body[i:]...)
			mp.Body[i] = &c
			return
		}
	}
	mp.Body = append(mp.Body, &c)
}

func (mp MiniPodSnippet) ABCGroup() string {
	if len(mp.ABC) > 0 {
		return string(mp.ABC[0])
	}
	return ""
}

/******************************************************************************/

func (mp *MiniPodSnippet) NeedRendering() bool {
	return len(mp.Body) > 0
}

// BuildTag
func (mp *MiniPodSnippet) BuildTag() ickcore.Tag {
	mp.Tag().
		SetTagName("div").
		AddClass("card mb-2").
		SetAttributeIf(mp.ABC != "", "data-abc", mp.ABC)
	return *mp.Tag()
}

// RenderContent
func (mp *MiniPodSnippet) RenderContent(out io.Writer) error {

	ickcore.RenderString(out, `<header class="card-header">`, `<p class="card-header-title">`)
	ickcore.RenderChild(out, mp, &mp.Icon)
	ickcore.RenderString(out, `<span class="ml-2">`+mp.Name+`</span>`)
	ickcore.RenderString(out, `</p>`, `</header>`)

	ishidden := "is-hidden"
	if mp.IsOpen {
		ishidden = ""
	}
	ickcore.RenderString(out, `<div class="card-content pt-2 px-2 pb-1 `+ishidden+`">`)
	var lastabc byte
	for _, cinbody := range mp.Body {
		if lastabc != 0 && len(cinbody.ABC) > 0 && cinbody.ABC[0] != lastabc {
			hidesplit := ""
			if cinbody.ABC[0] >= 'c' {
				hidesplit = "more is-hidden"
			}
			ickcore.RenderString(out, `<span class="hsplitter `+hidesplit+`"></span>`)
		}
		err := ickcore.RenderChild(out, mp, cinbody)
		if err != nil {
			return err
		}
		lastabc = cinbody.ABC[0]
	}
	ickcore.RenderString(out, `</div>`)

	if mp.HasMore > 0 {
		btnmore := ickui.Button("more").SetId(mp.Tag().SubId("btnmore")).SetColor(ick.COLOR_PRIMARY).SetOutlined(true).SetSize(ick.SIZE_SMALL)
		btnmore.OnClick = mp.OnShowMeMore
		ickcore.RenderString(out, `<div class="card-footer is-hidden">`)
		ickcore.RenderString(out, `<span class="card-footer-item is-justify-content-flex-start">`)
		ickcore.RenderChild(out, mp, btnmore)
		ickcore.RenderString(out, `</span>`)
		ickcore.RenderString(out, `</div>`)
	}
	return nil
}

// AddListeners
func (mp *MiniPodSnippet) AddListeners() {
	if !mp.DOM.IsInDOM() {
		return
	}

	chs := mp.DOM.ChildrenByClassName("card-header")
	chs[0].AddClass("is-clickable").Blur()

	chs[0].AddMouseEvent(event.MOUSE_ONCLICK, func(me *event.MouseEvent, e *dom.Element) {
		mp.OnOpenClose(!mp.IsOpen)
	})
}

/******************************************************************************/

func (mp *MiniPodSnippet) OnOpenClose(open bool) {
	if !open {
		mp.IsOpen = false
		cmores := mp.DOM.ChildrenByClassName("more")
		for _, cmore := range cmores {
			cmore.AddClass("is-hidden")
		}
	} else {
		mp.IsOpen = true
	}

	ccs := mp.DOM.ChildrenByClassName("card-content")
	ccs[0].SetClassIf(!mp.IsOpen, "is-hidden")

	cfs := mp.DOM.ChildrenByClassName("card-footer")
	if len(cfs) > 0 {
		cfs[0].SetClassIf(!mp.IsOpen, "is-hidden")
	}
}

func (mp *MiniPodSnippet) OnShowMeMore() {
	cmores := mp.DOM.ChildrenByClassName("more")
	for _, cmore := range cmores {
		cmore.RemoveClass("is-hidden")
	}
	cfs := mp.DOM.ChildrenByClassName("card-footer")
	cfs[0].AddClass("is-hidden")
}
