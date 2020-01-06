package css

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vanng822/css"
)

func TestRebaseURLs(t *testing.T) {
	ss := css.Parse(`h1 {
    background: no-repeat url("../images/image.png");
}
h2 {
    background: no-repeat url(icons/icon.png);
}
h3 {
    background: url('../pages/page.png');
}`)
	RebaseURLs(ss, "styles", "pages")
	got := Render(ss)
	want := `h1 {
    background: no-repeat url("../images/image.png");
}
h2 {
    background: no-repeat url("../styles/icons/icon.png");
}
h3 {
    background: url("../pages/page.png");
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}
func TestRemoveMediaRules(t *testing.T) {
	ss := css.Parse(`h1 { color: green; }
@media all and (min-width: 48em) {
    h2 { color: red; }
}`)
	RemoveMediaRules(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveKeyframeRules(t *testing.T) {
	ss := css.Parse(`h1 { color: green; }
@keyframes slidein {
    from {
        margin-left: 100%;
        width: 300%; 
    }
    to {
        margin-left: 0%;
        width: 100%;
    }
}`)
	RemoveKeyframeRules(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveFontFaceRules(t *testing.T) {
	ss := css.Parse(`h1 { color: green; }
@font-face {
    font-family: "somefont";
    src: url(somefont.ttf);
}`)
	RemoveFontFaceRules(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestKeepSimpleStyles(t *testing.T) {
	ss := css.Parse(`h1 { 
    color: green;
    margin: 10px;
    padding: 10px;
    line-height: 1.5;
}
blockquote.code {
    white-space: pre;
}
.squeeze-amzn {
    display: none;
}
`)
	KeepSimpleStyles(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
blockquote.code {
    white-space: pre;
}
.squeeze-amzn {
    display: none;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRemoveTextAlignJustify(t *testing.T) {
	ss := css.Parse(`img { 
    text-align: center;
}
p {
    text-align: justify;
}
`)
	RemoveTextAlignJustify(ss)
	got := Render(ss)
	want := `img {
    text-align: center;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestAddHeadingStyles(t *testing.T) {
	ss := css.Parse(`h1 { 
    color: green;
}
`)
	AddHeadingStyles(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
h1, h2, h3, h4, h5, h6 {
    font-weight: bold;
    -webkit-hyphens: none !important;
    hyphens: none !important;
    page-break-inside: avoid;
    page-break-after: avoid;
}
h5, h6 {
    text-transform: uppercase;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestAddFigureStyles(t *testing.T) {
	ss := css.Parse(`h1 { 
    color: green;
}
`)
	AddFigureStyles(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
.figure, figure {
    page-break-inside: avoid;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestAddAsideStyles(t *testing.T) {
	ss := css.Parse(`h1 { 
    color: green;
}
`)
	AddAsideStyles(ss)
	got := Render(ss)
	want := `h1 {
    color: green;
}
aside, .aside, .box, .boxg, .note, .note1, sidebar, .sidebar1, [data-type="note"], [data-type="tip"], [data-type="warning"] {
    border: 1px dotted #ddd;
    padding: 0em 1em !important;
    margin-top: 1em !important;
    margin-bottom: 1em !important;
    page-break-inside: avoid;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestRender(t *testing.T) {
	got := Render(css.Parse(`h1 { color: green; }

@media all and (min-width: 48em) {
    h2 { color: green; }
}

@font-face {
    font-family: "somefont";
    src: url(somefont.ttf);
}

@page {
    margin-bottom: 5pt;
    margin-top: 5pt;
}
`))
	want := `h1 {
    color: green;
}
@media all and (min-width: 48em) {
    h2 {
        color: green;
    }
}
@font-face {
    font-family: "somefont";
    src: url(somefont.ttf);
}
@page {
    margin-bottom: 5pt;
    margin-top: 5pt;
}
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}
