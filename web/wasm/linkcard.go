package main

import (
	"io"
	"net/url"

	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

type LinkCardSnippet struct {
	ickcore.BareSnippet
	dom.UI

	Name     string   // link card name, must be unique
	HRef     *url.URL // URL link card
	IsShrunk bool
}

// Ensuring LinkCardSnippet implements the right interface
var _ dom.UIComposer = (*LinkCardSnippet)(nil)

func LinkCard(name string) *LinkCardSnippet {
	n := new(LinkCardSnippet)
	n.Name = name
	return n
}

// ParseHRef tries to parse rawUrl to HRef ignoring error.
func (lnk *LinkCardSnippet) ParseHRef(rawUrl string) *LinkCardSnippet {
	lnk.HRef, _ = url.Parse(rawUrl)
	return lnk
}

// SetHRef sets the href url
func (lnk *LinkCardSnippet) SetHRef(href *url.URL) *LinkCardSnippet {
	if href == nil {
		lnk.HRef = nil
	} else {
		h := *href
		lnk.HRef = &h
	}
	return lnk
}

func (lnk *LinkCardSnippet) SetShrunk(shrunk bool) *LinkCardSnippet {
	lnk.IsShrunk = shrunk
	lnk.DOM.SetClassIf(!shrunk, "box")
	return lnk
}

/******************************************************************************/

// BuildTag builds the tag used to render the html element.
func (lnk *LinkCardSnippet) BuildTag() ickcore.Tag {
	lnk.Tag().SetTagName("div").AddClass("card").AddClassIf(!lnk.IsShrunk, "box")

	return *lnk.Tag()
}

// RenderContent writes the HTML string corresponding to the content of the HTML element.
// The default implementation for an HTMLSnippet snippet is to render all the internal stack of composers inside an enclosed HTML tag.
func (card *LinkCardSnippet) RenderContent(out io.Writer) error {
	l := ick.Link(ickcore.ToHTML(card.Name)).SetHRef(card.HRef)
	l.Tag().NoName = true
	ickcore.RenderChild(out, card, l)
	return nil
}
