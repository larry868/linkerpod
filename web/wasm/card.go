package main

import (
	"io"
	"net/url"

	"github.com/gosimple/slug"
	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
)

var (
	// _defaultIcon = ick.Icon("bi bi-x-square-fill").SetColor(ick.TXTCOLOR_GREY_LIGHT)
	_defaultIcon = `<svg xmlns="http://www.w3.org/2000/svg" fill="#e8e9ea" class="bi bi-x-square-fill" viewBox="0 0 16 16">
	<path d="M2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2zm3.354 4.646L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 1 1 .708-.708z"/>
  </svg>`
)

type CardSnippet struct {
	ickcore.BareSnippet
	dom.UI

	Name       string   // link card name, must be unique
	IconSrc    *url.URL // URL of the icon
	HRef       *url.URL // URL link card
	IsExpanded bool
	InMiniPods int
	ABC        string

	dftfaviconprefix string
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
	var err error
	card.HRef, err = url.Parse(rawUrl)
	if err == nil {
		card.dftfaviconprefix = slug.Make(card.HRef.Host)
	}
	return card
}

// SetHRef sets the href url
func (card *CardSnippet) SetHRef(href *url.URL) *CardSnippet {
	if href == nil {
		card.HRef = nil
	} else {
		h := *href
		card.HRef = &h
		card.dftfaviconprefix = slug.Make(card.HRef.Host)
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

	imgc := ick.Elem("div", `class="cardlink-img"`)
	if card.IconSrc == nil || card.IconSrc.Path == "" {
		img := ickcore.ToHTML(_defaultIcon)
		imgc.Append(img)
	} else {
		img := ickcore.ToHTML(`<img role="img" src="` + card.IconSrc.String() + `">`)
		imgc.Append(img)
	}

	c := ick.Elem("div", `class="cardlink"`,
		ick.Elem("div", `class="cardlink-left"`, imgc),
		ick.Elem("div", `class="cardlink-content"`,
			ick.Link(ickcore.ToHTML(card.Name)).SetHRef(card.HRef)))

	ickcore.RenderChild(out, card, c)
	return nil
}
