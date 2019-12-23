package html

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
)

func TestRemoveEmptySpans(t *testing.T) {
	h := `
<html>
<body>
<p><span class="Apple-converted-space">    </span></p>
<p>Paragraph</p>
</body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	RemoveEmptySpans(doc)
	got, _ := doc.Html()
	want := `<html><head></head><body>
<p></p>
<p>Paragraph</p>


</body></html>`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveEmptyDivs(t *testing.T) {
	h := `
<html>
<body>
<div>Text</div>
<div><p>Paragraph</p></div>
<div></div>
<div> </div></body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	RemoveEmptyDivs(doc)
	got, _ := doc.Html()
	want := `<html><head></head><body>
<div>Text</div>
<div><p>Paragraph</p></div>



</body></html>`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveLineBreaks(t *testing.T) {
	h := `
<html>
<body>
<h1>Heading</h1>
<p><br/></p>
<p><span></span><br/></p>
<p><i></i><br/></p>
<p><b></b><br/></p>
<p><img src="image.gif"><br/></p>
<p>Text</p>
<ul>
    <li>Item</li>
    <li><br/></li>
</ul>

</body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	RemoveLineBreaks(doc)
	got, _ := doc.Html()
	want := `<html><head></head><body>
<h1>Heading</h1>




<p><img src="image.gif"/><br/></p>
<p>Text</p>
<ul>
    <li>Item</li>
    
</ul>



</body></html>`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveContainers(t *testing.T) {
	h := `
<html>
<body>
<div><blockquote><p>This is a paragraph.</p></blockquote></div>
</body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	RemoveContainers(doc)
	got, _ := doc.Html()
	want := `<html><head></head><body>
<p>This is a paragraph.</p>


</body></html>`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveBoldnHeadings(t *testing.T) {
	h := `
<html>
<body>
<h1>One</h1>
<h2><b>Two</b></h2>
<h3><strong>Three</strong></h3>
</body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	RemoveBoldInHeadings(doc)
	got, _ := doc.Html()
	want := `<html><head></head><body>
<h1>One</h1>
<h2>Two</h2>
<h3>Three</h3>


</body></html>`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveExcessBlockquotes(t *testing.T) {
	h := `
<html>
<body>
<h1>One</h1>
<h2><blockquote>Two</blockquote></h2>
<ul>
    <li>Three</li>
    <li><blockquote>Four</blockquote></li>
</ul>
</body>
</html>
`
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(h)))
	RemoveExcessBlockquotes(doc)
	got, _ := doc.Html()
	want := `<html><head></head><body>
<h1>One</h1>
<h2>Two</h2>
<ul>
    <li>Three</li>
    <li>Four</li>
</ul>


</body></html>`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}
