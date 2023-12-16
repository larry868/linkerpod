package main

import (
	"io"
	"net/url"
	"strings"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/event"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
	"github.com/lolorenzo777/loadfavicon/v2/pkg/svg"
)

type MiniPodSnippet struct {
	ickcore.BareSnippet
	dom.UI

	Separator string // could be "blank" or "splitter" used by the father snipper rendering process
	Name      string // mini pod name

	iconSrc      *url.URL // URL of the icon
	iconSVG      string
	iconCssClass string
	iconChar     string // single char like a favicon

	Body []*CardSnippet // rendered as <div class="card-content">

	//IsExpanded bool
	IsOpen bool
	//HasMore    int
}

// Ensuring MiniPodSnippet implements the right interface
var _ dom.UIComposer = (*MiniPodSnippet)(nil)

func MiniPod(key string, name string, iconkey string, seperator string) *MiniPodSnippet {
	n := new(MiniPodSnippet)
	n.Body = make([]*CardSnippet, 0)
	n.Tag().SetId(key)
	n.Name = name
	if iconkey != "" {
		n.SetIcon(iconkey)
	} else {
		n.SetIcon("css=bi bi-dot")
	}
	n.Separator = seperator
	return n
}

func (mp MiniPodSnippet) IsSeparatorOnly() bool {
	return mp.Name == "" && mp.Separator != ""
}

// InsertCard inserts c card to the minipod. Return false if not inserted because of duplicate key
func (mp *MiniPodSnippet) InsertCard(c CardSnippet) bool {
	// cmore := ""
	// if len(c.SortIndex) >= 5 {
	// 	cmore = c.SortIndex[:5]
	// }
	// if cmore == "_more" {
	// 	c.Tag().AddClass("more is-hidden")
	// 	//mp.HasMore++
	// }

	for i, cinbody := range mp.Body {
		// check if not a duplicate id
		if c.Tag().Id() == cinbody.Tag().Id() {
			return false
		}

		// verbose.Printf(verbose.INFO, "%v newc=%s scanc=%s\n", i, c.Tag().Id(), cinbody.Tag().Id())
		if c.Tag().Id() < cinbody.Tag().Id() {
			mp.Body = append(mp.Body[:i+1], mp.Body[i:]...)
			mp.Body[i] = &c
			return true
		}
	}
	// verbose.Printf(verbose.INFO, "newc=%s append\n", c.Tag().Id())
	mp.Body = append(mp.Body, &c)
	return true
}

// SetIcon
func (mp *MiniPodSnippet) SetIcon(icon string) *MiniPodSnippet {
	switch {
	case svg.IsValidSVG([]byte(icon)):
		mp.iconSVG = icon
		// TODO: CardSnippet- change internal svg properties such as weight and height
	case len(icon) > 4 && icon[:4] == "chr=":
		mp.iconChar = strings.Trim(icon[4:], `"`)
		for index, runeValue := range mp.iconChar {
			if index == 0 {
				mp.iconChar = string(runeValue)
				break
			}
		}
	case len(icon) > 4 && icon[:4] == "css=":
		mp.iconCssClass = strings.Trim(icon[4:], `"`)
	case len(icon) > 0:
		mp.iconSrc, _ = url.Parse(icon)
	}
	return mp
}

/******************************************************************************/

func (mp *MiniPodSnippet) NeedRendering() bool {
	return true
	//XXX:
	return len(mp.Body) > 0
}

// BuildTag
func (mp *MiniPodSnippet) BuildTag() ickcore.Tag {
	mp.Tag().
		SetTagName("div").
		AddClass("card mb-1")
	return *mp.Tag()
}

// RenderContent
func (mp *MiniPodSnippet) RenderContent(out io.Writer) error {

	ickcore.RenderString(out, `<header class="card-header">`, `<p class="card-header-title pl-3">`)

	imgc := ick.Elem("span", `class="cardlink-img"`)
	var img *ickcore.HTMLString
	switch {
	case mp.iconSVG != "":
		img = ickcore.ToHTML(mp.iconSVG)
		imgc.Append(img)
	case mp.iconChar != "":
		img = ickcore.ToHTML(`<span>` + mp.iconChar + `</span>`)
		imgc.Append(img)
	case mp.iconCssClass != "":
		img = ickcore.ToHTML(`<span class="icon"><i class="` + mp.iconCssClass + `"></i></span>`)
		imgc.Append(img)
	case mp.iconSrc != nil && mp.iconSrc.Path != "":
		img = ickcore.ToHTML(`<img role="img" src="` + mp.iconSrc.String() + `">`)
		imgc.Append(img)
	default:
		img = ickcore.ToHTML(_defaultIcon)
		imgc.Append(img)
	}
	ickcore.RenderChild(out, mp, imgc)

	ickcore.RenderString(out, `<span class="ml-2">`+mp.Name+`</span>`)
	ickcore.RenderString(out, `</p>`, `</header>`)

	ishidden := "is-hidden"
	if mp.IsOpen {
		ishidden = ""
	}
	ickcore.RenderString(out, `<div class="card-content pt-2 px-2 pb-1 `+ishidden+`">`)
	mp.RenderLinkList(out)
	ickcore.RenderString(out, `</div>`)

	// if mp.HasMore > 0 {
	// 	btnmore := ickui.Button("More...").SetId(mp.Tag().SubId("btnmore")).SetColor(ick.COLOR_PRIMARY).SetOutlined(true).SetSize(ick.SIZE_SMALL)
	// 	btnmore.OnClick = mp.OnShowMeMore
	// 	ickcore.RenderString(out, `<div class="card-footer is-hidden">`)
	// 	ickcore.RenderString(out, `<span class="card-footer-item is-justify-content-flex-start">`)
	// 	ickcore.RenderChild(out, mp, btnmore)
	// 	ickcore.RenderString(out, `</span>`)
	// 	ickcore.RenderString(out, `</div>`)
	// }
	return nil
}

// RenderLinkList
func (mp *MiniPodSnippet) RenderLinkList(out io.Writer) error {
	for _, cinbody := range mp.Body {
		//		if lastabc != 0 && len(cinbody.ABC) > 0 && cinbody.ABC[0] != lastabc {
		//hidesplit := ""
		// if cinbody.ABC[0] >= 'c' {
		// 	hidesplit = "more is-hidden"
		// }
		//ickcore.RenderString(out, `<span class="hsplitter `+hidesplit+`"></span>`)
		// }
		err := ickcore.RenderChild(out, mp, cinbody)
		if err != nil {
			return err
		}
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
