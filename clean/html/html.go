package html

import (
	"strings"

	"golang.org/x/net/html/atom"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func RemoveEmptySpans(doc *goquery.Document) {
	doc.Find("span").Each(func(_ int, s *goquery.Selection) {
		elementCount := s.Children().Length()
		text := strings.TrimSpace(s.Text())
		if elementCount == 0 && text == "" {
			s.Remove()
		}
	})
}

func RemoveEmptyDivs(doc *goquery.Document) {
	doc.Find("div").Each(func(_ int, s *goquery.Selection) {
		elementCount := s.Children().Length()
		text := strings.TrimSpace(s.Text())
		if elementCount == 0 && text == "" {
			// IDs may be referenced by links, preserve them
			ID := s.AttrOr("id", "")
			if ID == "" {
				s.Remove()
			} else {
				s.ReplaceWithNodes(anchorNode(ID))
			}
		}
	})
}

func RemoveLineBreaks(doc *goquery.Document) {
	doc.Find("p > br, li > br").Each(func(_ int, s *goquery.Selection) {
		parent := s.Parent()
		text := strings.TrimSpace(parent.Text())

		siblingElements := s.Siblings()
		styleElements := siblingElements.Filter("b, em, i, span, strong, u")

		if text == "" && siblingElements.Length() == styleElements.Length() {
			parent.Remove()
		}
	})
}

func RemoveContainers(doc *goquery.Document) {
	body := doc.Find("body")

	for {
		childElements := body.Children().Filter("*:not(a[id])")
		// Detect elements that container the body contents
		if childElements.Length() == 1 && childElements.Filter("div, blockquote").Length() == 1 {
			// Replace with container's child nodes
			childElements.Each(func(_ int, s *goquery.Selection) {
				node := s.Nodes[0]
				var grandChildNodes []*html.Node
				// IDs may be referenced by links, preserve them
				if ID := s.AttrOr("id", ""); ID != "" {
					grandChildNodes = append(grandChildNodes, anchorNode(ID))
				}
				// To capture element and text nodes
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					grandChildNodes = append(grandChildNodes, c)
				}
				s.ReplaceWithNodes(grandChildNodes...)
			})
		} else {
			break
		}
	}
}

func RemoveBoldInHeadings(doc *goquery.Document) {
	doc.Find("h1, h2, h3, h4, h5, h6").Each(func(_ int, s *goquery.Selection) {
		s.Find("b, strong").Each(func(_ int, s2 *goquery.Selection) {
			node := s2.Nodes[0]
			// To capture element and text nodes
			var childNodes []*html.Node
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				childNodes = append(childNodes, c)
			}
			s2.ReplaceWithNodes(childNodes...)
		})
	})
}

func RemoveExcessBlockquotes(doc *goquery.Document) {
	doc.Find("h1, h2, h3, h4, h5, h6, li").Each(func(_ int, s *goquery.Selection) {
		s.Find("blockquote").Each(func(_ int, s2 *goquery.Selection) {
			node := s2.Nodes[0]
			// To capture element and text nodes
			var childNodes []*html.Node
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				childNodes = append(childNodes, c)
			}
			s2.ReplaceWithNodes(childNodes...)
		})
	})
}

func anchorNode(ID string) *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.A,
		Data:     "a",
		Attr: []html.Attribute{
			{Key: "id", Val: ID},
		},
	}
}
