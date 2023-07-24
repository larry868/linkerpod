package main

import (
	"io"
	"net/url"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

var (
	_defaultImage = ick.Icon("bi bi-x-square-fill").SetColor(ick.TXTCOLOR_GREY_LIGHT)
)

type CardSnippet struct {
	ickcore.BareSnippet
	dom.UI

	Name       string // link card name, must be unique
	Image      ick.ICKImage
	HRef       *url.URL // URL link card
	IsExpanded bool
	InMiniPods int
	ABC        string
}

// Ensuring LinkCardSnippet implements the right interface
var _ dom.UIComposer = (*CardSnippet)(nil)

func Card(key string, name string) *CardSnippet {
	n := new(CardSnippet)
	n.Tag().SetId(key)
	n.Name = name
	n.ABC = "b"
	return n
}

// ParseHRef tries to parse rawUrl to HRef ignoring error.
func (card *CardSnippet) ParseHRef(rawUrl string) *CardSnippet {
	card.HRef, _ = url.Parse(rawUrl)
	return card
}

// SetHRef sets the href url
func (card *CardSnippet) SetHRef(href *url.URL) *CardSnippet {
	if href == nil {
		card.HRef = nil
	} else {
		h := *href
		card.HRef = &h
	}
	return card
}

func (card *CardSnippet) Expand(f bool) *CardSnippet {
	card.IsExpanded = f
	card.DOM.SetClassIf(!f, "py-1 px-3", "py-3 px-5")
	return card
}

/******************************************************************************/

func (card *CardSnippet) BuildTag() ickcore.Tag {
	card.Tag().
		SetTagName("div").
		AddClass("card mb-1 px-3").
		SetClassIf(!card.IsExpanded, "py-1", "py-3").
		SetAttributeIf(card.ABC != "", "data-abc", card.ABC)
	return *card.Tag()
}

func (card *CardSnippet) RenderContent(out io.Writer) error {

	var img ickcore.ContentComposer
	if !card.Image.NeedRendering() {
		img = _defaultImage
	} else {
		img = &card.Image
	}

	c := ick.Elem("div", `class="cardlink"`,
		ick.Elem("div", `class="cardlink-left"`, img),
		ick.Elem("div", `class="cardlink-content"`,
			ick.Link(ickcore.ToHTML(card.Name)).SetHRef(card.HRef)))

	ickcore.RenderChild(out, card, c)
	return nil
}
