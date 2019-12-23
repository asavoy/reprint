package ncx

import (
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
	tocNCX := []byte(
		`<?xml version="1.0" encoding="utf-8" ?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
    <head>
        <meta name="dtb:uid" content="ID:ISBN:123456789"/>
        <meta name="dtb:depth" content="2"/>
        <meta name="dtb:totalPageCount" content="0"/>
        <meta name="dtb:maxPageNumber" content="0"/>
    </head>
    <docTitle>
        <text>A Sample Book</text>
    </docTitle>
    <docAuthor>
        <text>John Smith, Jane Doe</text>
    </docAuthor>
    <navMap>
        <navPoint id="p1" playOrder="1">
            <navLabel>
                <text>&amp;#160;Foreword</text>
            </navLabel>
            <content src="1.xhtml#start"/>
        </navPoint>
        <navPoint id="p2" playOrder="2">
            <navLabel>
                <text>Main</text>
            </navLabel>
            <content src="2.xhtml#start"/>
            <navPoint id="p3" playOrder="3">
                <navLabel>
                    <text>Chapter 1</text>
                </navLabel>
                <content src="3.xhtml#start"/>
            </navPoint>
            <navPoint id="p4" playOrder="4">
                <navLabel>
                    <text>Chapter 2</text>
                </navLabel>
                <content src="4.xhtml#start"/>
            </navPoint>
        </navPoint>
    </navMap>
</ncx>`)
	got, err := Read(tocNCX)
	if err != nil {
		t.Fatal(err)
	}
	want := NCX{
		XMLName: xml.Name{
			Space: "http://www.daisy.org/z3986/2005/ncx/",
			Local: "ncx",
		},
		Version: "2005-1",
		Metas: []Meta{
			{Name: "dtb:uid", Content: "ID:ISBN:123456789"},
			{Name: "dtb:depth", Content: "2"},
			{Name: "dtb:totalPageCount", Content: "0"},
			{Name: "dtb:maxPageNumber", Content: "0"},
		},
		Title:  "A Sample Book",
		Author: "John Smith, Jane Doe",
		NavPoints: []NavPoint{
			{
				ID:        "p1",
				PlayOrder: "1",
				Label:     "&#160;Foreword",
				Content:   Content{Src: "1.xhtml#start"},
			},
			{
				ID:        "p2",
				PlayOrder: "2",
				Label:     "Main",
				Content:   Content{Src: "2.xhtml#start"},
				NavPoints: []NavPoint{
					{
						ID:        "p3",
						PlayOrder: "3",
						Label:     "Chapter 1",
						Content:   Content{Src: "3.xhtml#start"},
					},
					{
						ID:        "p4",
						PlayOrder: "4",
						Label:     "Chapter 2",
						Content:   Content{Src: "4.xhtml#start"},
					},
				},
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestWrite(t *testing.T) {
	ncx := NCX{
		Version: "2005-1",
		Metas: []Meta{
			{Name: "dtb:uid", Content: "ID:ISBN:123456789"},
			{Name: "dtb:depth", Content: "2"},
			{Name: "dtb:totalPageCount", Content: "0"},
			{Name: "dtb:maxPageNumber", Content: "0"},
		},
		Title:  "A Sample Book",
		Author: "John Smith, Jane Doe",
		NavPoints: []NavPoint{
			{
				ID:        "p1",
				PlayOrder: "1",
				Label:     "&#160;Foreword",
				Content:   Content{Src: "1.xhtml#start"},
			},
			{
				ID:        "p2",
				PlayOrder: "2",
				Label:     "Main",
				Content:   Content{Src: "2.xhtml#start"},
				NavPoints: []NavPoint{
					{
						ID:        "p3",
						PlayOrder: "3",
						Label:     "Chapter 1",
						Content:   Content{Src: "3.xhtml#start"},
					},
					{
						ID:        "p4",
						PlayOrder: "4",
						Label:     "Chapter 2",
						Content:   Content{Src: "4.xhtml#start"},
					},
				},
			},
		},
	}
	got, err := Write(ncx)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(
		`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
    <head>
        <meta name="dtb:uid" content="ID:ISBN:123456789"></meta>
        <meta name="dtb:depth" content="2"></meta>
        <meta name="dtb:totalPageCount" content="0"></meta>
        <meta name="dtb:maxPageNumber" content="0"></meta>
    </head>
    <docTitle>
        <text>A Sample Book</text>
    </docTitle>
    <docAuthor>
        <text>John Smith, Jane Doe</text>
    </docAuthor>
    <navMap>
        <navPoint id="p1" playOrder="1">
            <navLabel>
                <text>&amp;#160;Foreword</text>
            </navLabel>
            <content src="1.xhtml#start"></content>
        </navPoint>
        <navPoint id="p2" playOrder="2">
            <navLabel>
                <text>Main</text>
            </navLabel>
            <content src="2.xhtml#start"></content>
            <navPoint id="p3" playOrder="3">
                <navLabel>
                    <text>Chapter 1</text>
                </navLabel>
                <content src="3.xhtml#start"></content>
            </navPoint>
            <navPoint id="p4" playOrder="4">
                <navLabel>
                    <text>Chapter 2</text>
                </navLabel>
                <content src="4.xhtml#start"></content>
            </navPoint>
        </navPoint>
    </navMap>
</ncx>`)
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Error("got != want:\n", diff)
	}
}
