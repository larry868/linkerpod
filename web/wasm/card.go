package main

import (
	"io"
	"net/url"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

type CardSnippet struct {
	ickcore.BareSnippet
	dom.UI

	Name       string   // link card name, must be unique
	HRef       *url.URL // URL link card
	IsShrunk   bool
	InMiniPods int
	ABC        string
}

// Ensuring LinkCardSnippet implements the right interface
var _ dom.UIComposer = (*CardSnippet)(nil)

func Card(name string) *CardSnippet {
	n := new(CardSnippet)
	n.Name = name
	n.ABC = "b"
	n.Tag().SetId(n.Name)
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

func (card *CardSnippet) SetShrunk(shrunk bool) *CardSnippet {
	card.IsShrunk = shrunk
	card.DOM.SetClassIf(!shrunk, "p-2 mb-2")
	return card
}

/******************************************************************************/

func (card *CardSnippet) BuildTag() ickcore.Tag {
	card.Tag().
		SetTagName("div").
		AddClass("card").
		SetClassIf(!card.IsShrunk, "p-2 mb-2").
		SetAttributeIf(card.ABC != "", "data-abc", card.ABC)
	return *card.Tag()
}

func (card *CardSnippet) RenderContent(out io.Writer) error {
	l := ick.Link(ickcore.ToHTML(card.Name)).SetHRef(card.HRef)
	l.Tag().NoName = true
	ickcore.RenderChild(out, card, l)
	return nil
}
