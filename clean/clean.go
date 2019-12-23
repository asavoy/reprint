package clean

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/vanng822/css"
	"golang.org/x/net/html"

	"github.com/asavoy/reprint/book"
	cleanCSS "github.com/asavoy/reprint/clean/css"
	cleanHTML "github.com/asavoy/reprint/clean/html"
)

func Clean(b *book.Book) error {
	var newResources []book.Resource
	deleteResourceByPath := make(map[string]bool)

	for _, resource := range b.Resources {
		if resource.MediaType == "application/xhtml+xml" {
			doc, ss, ssResources, err := decomposePage(resource, *b)
			if err != nil {
				return err
			}

			cleanPage(doc, ss)
			docHTML, err := doc.Html()
			if err != nil {
				return err
			}
			newResources = append(newResources, book.Resource{
				ID:        resource.ID,
				Path:      resource.Path,
				MediaType: resource.MediaType,
				Contents:  []byte(docHTML),
			})

			deleteResourceByPath[resource.Path] = true
			for _, ssResource := range ssResources {
				deleteResourceByPath[ssResource.Path] = true
			}
		}
	}

	var resources []book.Resource
	for _, r := range newResources {
		resources = append(resources, r)
	}
	for _, r := range b.Resources {
		if _, ok := deleteResourceByPath[r.Path]; !ok {
			resources = append(resources, r)
		}
	}

	b.Resources = resources
	return nil
}

func cleanPage(doc *goquery.Document, ss *css.CSSStyleSheet) {
	extractInlineStyles(doc, ss)

	cleanCSS.RemoveMediaRules(ss)
	cleanCSS.RemoveKeyframeRules(ss)
	cleanCSS.RemoveFontFaceRules(ss)
	// TODO: consider making these optional
	cleanCSS.RemoveColors(ss)
	cleanCSS.RemoveTextAlignJustify(ss)

	cleanCSS.KeepSimpleStyles(ss)
	cleanCSS.AddHeadingStyles(ss)
	cleanCSS.AddFigureStyles(ss)
	cleanCSS.AddAsideStyles(ss)
	cleanCSS.AddTableStyles(ss)

	cleanHTML.RemoveEmptySpans(doc)
	cleanHTML.RemoveEmptyDivs(doc)
	cleanHTML.RemoveLineBreaks(doc)
	cleanHTML.RemoveContainers(doc)
	cleanHTML.RemoveBoldInHeadings(doc)
	cleanHTML.RemoveExcessBlockquotes(doc)

	// Add styles directly into document
	renderedStyles := cleanCSS.Render(ss)
	styleNode := &html.Node{
		Type: html.ElementNode,
		Data: "style",
		Attr: []html.Attribute{
			{Key: "type", Val: "text/css"},
		},
	}
	styleNode.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: renderedStyles,
	})
	doc.Find("head").AppendNodes(styleNode)
}

func decomposePage(page book.Resource, b book.Book) (*goquery.Document, *css.CSSStyleSheet, []book.Resource, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page.Contents))
	if err != nil {
		return nil, nil, nil, err
	}

	styleSelection := doc.Find("style")
	linkSelection := doc.Find("link[rel=stylesheet]")

	mergedStyles := styleSelection.Text()
	var linkedResources []book.Resource
	linkStyles := linkSelection.Map(func(i int, s *goquery.Selection) string {
		relPath, exists := s.Attr("href")
		if !exists {
			err = errors.New("link missing href attribute")
			return ""
		}
		absPath := path.Clean(path.Join(path.Dir(page.Path), relPath))
		r, e := b.GetResource(absPath)
		if e != nil {
			err = e
			return ""
		}
		linkedResources = append(linkedResources, r)
		linkedSS := css.Parse(string(r.Contents))
		cleanCSS.RebaseURLs(linkedSS, path.Dir(absPath), path.Dir(page.Path))
		return cleanCSS.Render(linkedSS)
	})
	if err != nil {
		return nil, nil, nil, err
	}
	for _, linkStyle := range linkStyles {
		mergedStyles += linkStyle
	}

	styleSelection.Remove()
	linkSelection.Remove()

	ss := css.Parse(mergedStyles)

	return doc, ss, linkedResources, nil
}

func extractInlineStyles(doc *goquery.Document, ss *css.CSSStyleSheet) {
	doc.Find("[style]").Each(func(i int, s *goquery.Selection) {
		className := fmt.Sprintf("reprint_%s_%d", s.Nodes[0].Data, i)
		s.AddClass(className)
		cssText, _ := s.Attr("style")
		s.RemoveAttr("style")

		var styles []*css.CSSStyleDeclaration
		for _, style := range css.ParseBlock(cssText) {
			styles = append(styles, &css.CSSStyleDeclaration{
				Property:  style.Property,
				Value:     style.Value,
				Important: 1,
			})
		}

		ss.CssRuleList = append(ss.CssRuleList, &css.CSSRule{
			Type: css.STYLE_RULE,
			Style: css.CSSStyleRule{
				SelectorText: "." + className,
				Styles:       styles,
			},
			Rules: nil,
		})
	})
}
