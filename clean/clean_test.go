package clean

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
	"github.com/vanng822/css"

	"github.com/asavoy/reprint/book"
	cleanCSS "github.com/asavoy/reprint/clean/css"
)

func TestDecomposePage(t *testing.T) {
	page := book.Resource{
		Path: "text/page.xhtml",
		Contents: []byte(`<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en-US">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<link href="../stylesheet.css" rel="stylesheet" type="text/css"/>
	<style type="text/css">h1{ color: green; }</style>
	<style type="text/css">h2{ color: blue; }</style>
</head>
<body>
	<h1>Chapter 1</h1>
	<p>The quick brown fox jumps over the lazy dog</p>
</body>
</html>`),
	}
	ssResource := book.Resource{
		Path:     "stylesheet.css",
		Contents: []byte(`h3{ color: purple; }`),
	}
	b := book.Book{
		Resources: []book.Resource{
			ssResource,
			page,
		},
	}
	doc, ss, gotSSResources, err := decomposePage(page, b)
	if err != nil {
		t.Fatal(err)
	}
	gotHTML, err := doc.Html()
	if err != nil {
		t.Fatal(err)
	}
	gotCSS := cleanCSS.Render(ss)
	wantHTML := `<!--?xml version='1.0' encoding='utf-8'?--><html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en-US"><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	
	
	
</head>
<body>
	<h1>Chapter 1</h1>
	<p>The quick brown fox jumps over the lazy dog</p>

</body></html>`
	wantCSS := `h1 {
    color: green;
}
h2 {
    color: blue;
}
h3 {
    color: purple;
}
`
	wantSSResources := []book.Resource{ssResource}
	if diff := cmp.Diff(wantHTML, gotHTML); diff != "" {
		t.Error("html got != want:\n", diff)
	}
	if diff := cmp.Diff(wantCSS, gotCSS); diff != "" {
		t.Error("css got != want:\n", diff)
	}
	if diff := cmp.Diff(wantSSResources, gotSSResources); diff != "" {
		t.Error("gotSSResources got != want:\n", diff)
	}
}

func TestExtractInlineStyles(t *testing.T) {
	h := `
<html>
<body>
<h1 style="text-align: center;">Heading</h1>
</body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	ss := &css.CSSStyleSheet{}
	extractInlineStyles(doc, ss)
	gotHTML, _ := doc.Html()
	gotCSS := cleanCSS.Render(ss)
	wantHTML := `<html><head></head><body>
<h1 class="reprint_h1_0">Heading</h1>


</body></html>`
	wantCSS := `.reprint_h1_0 {
    text-align: center !important;
}
`

	if diff := cmp.Diff(wantHTML, gotHTML); diff != "" {
		t.Error("gotHTML != wantHTML:\n", diff)
	}
	if diff := cmp.Diff(wantCSS, gotCSS); diff != "" {
		t.Error("gotCSS != wantCSS:\n", diff)
	}
}

func TestExtractImageStyles(t *testing.T) {
	html := []byte(`<?xml version='1.0' encoding='utf-8'?>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en-US">
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	
</head>
<body>
	<figure>
		<img src="flourish-ultra-high-definition.png" />
	</figure>
	<p>Paragraph</p>
</body>
</html>`)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))

	ss := css.Parse(`
figure {
	width: 2em;
    line-height: 1.5;
}
img {
	width: 100%;
}
body {
	width: 600px;
}
`)
	gotSS := extractImageStyles(doc, ss)
	if err != nil {
		t.Fatal(err)
	}
	gotHTML, err := doc.Html()
	if err != nil {
		t.Fatal(err)
	}
	gotCSS := cleanCSS.Render(gotSS)
	wantHTML := `<!--?xml version='1.0' encoding='utf-8'?--><html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en-US"><head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	
</head>
<body>
	<figure class="reprint_images_0_0">
		<img src="flourish-ultra-high-definition.png" class="reprint_images_0"/>
	</figure>
	<p>Paragraph</p>

</body></html>`
	wantCSS := `.reprint_images_0 {
    width: 100%;
}
.reprint_images_0_0 {
    width: 2em;
}
`
	if diff := cmp.Diff(wantHTML, gotHTML); diff != "" {
		t.Error("html got != want:\n", diff)
	}
	if diff := cmp.Diff(wantCSS, gotCSS); diff != "" {
		t.Error("css got != want:\n", diff)
	}
}
