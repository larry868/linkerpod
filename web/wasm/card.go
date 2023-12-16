package main

import (
	"io"
	"net/url"
	"strings"

	"github.com/gosimple/slug"
	"github.com/icecake-framework/icecake/pkg/dom"
	"github.com/icecake-framework/icecake/pkg/ick"
	"github.com/icecake-framework/icecake/pkg/ickcore"
	"github.com/lolorenzo777/loadfavicon/v2/pkg/svg"
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

	Name string // link card name, must be unique

	iconSrc      *url.URL // URL of the icon
	iconSVG      string
	iconCssClass string
	iconChar     string // single char like a favicon

	HRef       *url.URL // URL link card
	IsExpanded bool
	//InMiniPods int
	// SortIndex string

	dftfaviconprefix string
}

// Ensuring LinkCardSnippet implements the right interface
var _ dom.UIComposer = (*CardSnippet)(nil)

func Card(key string, name string) *CardSnippet {
	n := new(CardSnippet)
	n.Tag().SetId(key)
	n.Name = name
	// n.SortIndex = "b"
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

// SetIcon
func (card *CardSnippet) SetIcon(icon string) *CardSnippet {
	switch {
	case svg.IsValidSVG([]byte(icon)):
		card.iconSVG = icon
		// TODO: CardSnippet- change internal svg properties such as weight and height
	case len(icon) > 4 && icon[:4] == "chr=":
		card.iconChar = strings.Trim(icon[4:], `"`)
		for index, runeValue := range card.iconChar {
			if index == 0 {
				card.iconChar = string(runeValue)
				break
			}
		}
	case len(icon) > 4 && icon[:4] == "css=":
		card.iconCssClass = strings.Trim(icon[4:], `"`)
	case len(icon) > 0:
		card.iconSrc, _ = url.Parse(icon)
	}
	return card
}

/******************************************************************************/

func (card *CardSnippet) BuildTag() ickcore.Tag {
	card.Tag().
		SetTagName("div").
		AddClass("card mb-1 px-3").
		SetClassIf(!card.IsExpanded, "py-1", "py-3")
		// SetAttributeIf(card.SortIndex != "", "data-abc", card.SortIndex)
	return *card.Tag()
}

func (card *CardSnippet) RenderContent(out io.Writer) error {

	imgc := ick.Elem("div", `class="cardlink-img"`)
	switch {
	case card.iconSVG != "":
		img := ickcore.ToHTML(card.iconSVG)
		imgc.Append(img)
	case card.iconChar != "":
		img := ickcore.ToHTML(`<span>` + card.iconChar + `</span>`)
		imgc.Append(img)
	case card.iconCssClass != "":
		img := ickcore.ToHTML(`<span class="icon"><i class="` + card.iconCssClass + `"></i></span>`)
		imgc.Append(img)
	case card.iconSrc != nil && card.iconSrc.Path != "":
		img := ickcore.ToHTML(`<img role="img" src="` + card.iconSrc.String() + `">`)
		imgc.Append(img)
	default:
		img := ickcore.ToHTML(_defaultIcon)
		imgc.Append(img)
	}

	c := ick.Elem("div", `class="cardlink"`,
		ick.Elem("div", `class="cardlink-left"`, imgc),
		ick.Elem("div", `class="cardlink-content"`,
			ick.Link(ickcore.ToHTML(card.Name)).SetHRef(card.HRef)))

	ickcore.RenderChild(out, card, c)
	return nil
}
