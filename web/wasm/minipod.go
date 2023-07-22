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
	ickui.ICKCard

	Name string // mini pod name, must be unique
	Icon ick.ICKIcon
	//cards    []CardSnippet
	IsShrunk bool
	IsOpen   bool
}

// Ensuring MiniPodSnippet implements the right interface
var _ dom.UIComposer = (*MiniPodSnippet)(nil)

func MiniPod(name string, iconkey string) *MiniPodSnippet {
	n := new(MiniPodSnippet)
	//n.cards = make([]CardSnippet, 0)
	n.Name = name
	if iconkey != "" {
		n.Icon.Key = iconkey
	} else {
		n.Icon.Key = "bi bi-dot"
	}
	return n
}

func (mp *MiniPodSnippet) AppendCard(c CardSnippet) {
	mp.Body.Append(&c)
}

func (mp *MiniPodSnippet) SetShrunk(shrunk bool) *MiniPodSnippet {
	mp.IsShrunk = shrunk
	//mp.DOM.SetClassIf(!shrunk, "mb-2")
	return mp
}

/******************************************************************************/

func (mp *MiniPodSnippet) BuildTag() ickcore.Tag {
	tag := mp.ICKCard.BuildTag()
	tag.SetId("minipod-" + mp.Name)
	tag.SetClassIf(!mp.IsShrunk, "mb-2")

	mp.Body.Tag().AddClass("pt-2 px-2 pb-1 is-hidden")

	return tag
}

func (mp *MiniPodSnippet) RenderContent(out io.Writer) error {

	title := ick.Elem("span", "").
		Append(&mp.Icon).
		Append(ickcore.ToHTML(`<span class="ml-2">` + mp.Name + `</span>`))

	ickcore.RenderChild(&mp.Title, mp, title)

	return mp.ICKCard.RenderContent(out)
}

func (mp *MiniPodSnippet) AddListeners() {
	if !mp.DOM.IsInDOM() {
		return
	}

	chs := mp.DOM.ChildrenByClassName("card-header")
	chs[0].AddClass("is-clickable")

	ccs := mp.DOM.ChildrenByClassName("card-content")

	// etitle := dom.Id(mp.Tag().SubId("title"))
	chs[0].AddMouseEvent(event.MOUSE_ONCLICK, func(me *event.MouseEvent, e *dom.Element) {
		if mp.IsOpen {
			mp.IsOpen = false
		} else {
			mp.IsOpen = true
		}
		ccs[0].SetClassIf(!mp.IsOpen, "is-hidden")
	})
}
