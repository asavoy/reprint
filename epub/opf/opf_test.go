package opf

import (
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
	contentOPF := []byte(
		`<?xml version='1.0' encoding='UTF-8'?>
<package xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:opf="http://www.idpf.org/2007/opf" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="id">
    <metadata>
        <dc:rights>Public domain</dc:rights>
        <dc:identifier opf:scheme="URI" id="id">http://www.gutenberg.org/ebooks/1184</dc:identifier>
        <dc:creator opf:file-as="Dumas, Alexandre">Alexandre Dumas</dc:creator>
        <dc:title>The Count of Monte Cristo, Illustrated</dc:title>
        <dc:language xsi:type="dcterms:RFC4646">en</dc:language>
        <dc:subject>Historical fiction</dc:subject>
        <dc:subject>Adventure stories</dc:subject>
        <dc:date opf:event="publication">1998-01-01</dc:date>
        <dc:date opf:event="conversion">2019-05-31T14:42:00.544171+00:00</dc:date>
        <dc:source>http://www.gutenberg.org/files/1184/1184-h/1184-h.htm</dc:source>
        <meta name="cover" content="item1"/>
        <meta name="calibre:timestamp" content="2017-10-10T20:22:18.989147+00:00"/>
    </metadata>
    <manifest>
        <item media-type="image/jpeg" href="cover.jpg" id="item1"/>
        <item media-type="text/css" href="pgepub.css" id="item2"/>
        <item media-type="text/css" href="0.css" id="item3"/>
        <item media-type="text/css" href="1.css" id="item4"/>
        <item media-type="application/xhtml+xml" href="page1.html" id="item5"/>
        <item media-type="application/xhtml+xml" href="page2.html" id="item6"/>
        <item media-type="application/x-dtbncx+xml" href="toc.ncx" id="ncx"/>
        <item media-type="application/xhtml+xml" href="cover.html" id="coverpage-wrapper"/>
    </manifest>
    <spine toc="ncx">
        <itemref idref="coverpage-wrapper" linear="no"/>
        <itemref idref="item5" linear="yes"/>
        <itemref idref="item6" linear="yes"/>
    </spine>
    <guide>
        <reference type="toc" title="Contents" href="page1.html#pgepubid00001"/>
        <reference type="cover" title="CoverImageID" href="cover.html"/>
    </guide>
</package>
`)
	got, err := Read(contentOPF)
	if err != nil {
		t.Fatal(err)
	}
	want := Package{
		XMLName:          xml.Name{Space: "http://www.idpf.org/2007/opf", Local: "package"},
		UniqueIdentifier: "id",
		Version:          "2.0",
		Metadata: Metadata{
			XMLName: xml.Name{Space: "http://www.idpf.org/2007/opf", Local: "metadata"},
			Title:   "The Count of Monte Cristo, Illustrated",
			Identifiers: []Identifier{
				{ID: "id", Value: "http://www.gutenberg.org/ebooks/1184"},
			},
			Creators:  []string{"Alexandre Dumas"},
			Publisher: "",
			Language:  "en",
			Subjects:  []string{"Historical fiction", "Adventure stories"},
			Rights:    "Public domain",
			Source:    "http://www.gutenberg.org/files/1184/1184-h/1184-h.htm",
			Dates: []Date{
				{Event: "publication", Value: "1998-01-01"},
				{Event: "conversion", Value: "2019-05-31T14:42:00.544171+00:00"},
			},
			Metas: []Meta{
				{Name: "cover", Content: "item1"},
				{Name: "calibre:timestamp", Content: "2017-10-10T20:22:18.989147+00:00"},
			},
		},
		Manifest: Manifest{
			XMLName: xml.Name{Space: "http://www.idpf.org/2007/opf", Local: "manifest"},
			Items: []ManifestItem{
				{ID: "item1", Href: "cover.jpg", MediaType: "image/jpeg"},
				{ID: "item2", Href: "pgepub.css", MediaType: "text/css"},
				{ID: "item3", Href: "0.css", MediaType: "text/css"},
				{ID: "item4", Href: "1.css", MediaType: "text/css"},
				{ID: "item5", Href: "page1.html", MediaType: "application/xhtml+xml"},
				{ID: "item6", Href: "page2.html", MediaType: "application/xhtml+xml"},
				{ID: "ncx", Href: "toc.ncx", MediaType: "application/x-dtbncx+xml"},
				{ID: "coverpage-wrapper", Href: "cover.html", MediaType: "application/xhtml+xml"},
			},
		},
		Spine: Spine{
			XMLName: xml.Name{Space: "http://www.idpf.org/2007/opf", Local: "spine"},
			Toc:     "ncx",
			ItemRefs: []SpineItemRef{
				{IDRef: "coverpage-wrapper", Linear: "no"},
				{IDRef: "item5", Linear: "yes"},
				{IDRef: "item6", Linear: "yes"},
			},
		},
		Guide: Guide{
			XMLName: xml.Name{Space: "http://www.idpf.org/2007/opf", Local: "guide"},
			Refs: []GuideRef{
				{Type: "toc", Title: "Contents", Href: "page1.html#pgepubid00001"},
				{Type: "cover", Title: "CoverImageID", Href: "cover.html"},
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestWrite(t *testing.T) {
	ncx := Package{
		UniqueIdentifier: "id",
		Version:          "2.0",
		Metadata: Metadata{
			Title: "The Count of Monte Cristo, Illustrated",
			Identifiers: []Identifier{
				{ID: "id", Value: "http://www.gutenberg.org/ebooks/1184"},
			},
			Creators:  []string{"Alexandre Dumas"},
			Publisher: "",
			Language:  "en",
			Subjects:  []string{"Historical fiction", "Adventure stories"},
			Rights:    "Public domain",
			Source:    "http://www.gutenberg.org/files/1184/1184-h/1184-h.htm",
			Dates: []Date{
				{Event: "publication", Value: "1998-01-01"},
				{Event: "conversion", Value: "2019-05-31T14:42:00.544171+00:00"},
			},
			Metas: []Meta{
				{Name: "cover", Content: "item1"},
				{Name: "calibre:timestamp", Content: "2017-10-10T20:22:18.989147+00:00"},
			},
		},
		Manifest: Manifest{
			Items: []ManifestItem{
				{ID: "item1", Href: "cover.jpg", MediaType: "image/jpeg"},
				{ID: "item2", Href: "pgepub.css", MediaType: "text/css"},
				{ID: "item3", Href: "0.css", MediaType: "text/css"},
				{ID: "item4", Href: "1.css", MediaType: "text/css"},
				{ID: "item5", Href: "page1.html", MediaType: "application/xhtml+xml"},
				{ID: "item6", Href: "page2.html", MediaType: "application/xhtml+xml"},
				{ID: "ncx", Href: "toc.ncx", MediaType: "application/x-dtbncx+xml"},
				{ID: "coverpage-wrapper", Href: "cover.html", MediaType: "application/xhtml+xml"},
			},
		},
		Spine: Spine{
			Toc: "ncx",
			ItemRefs: []SpineItemRef{
				{IDRef: "coverpage-wrapper", Linear: "no"},
				{IDRef: "item5", Linear: "yes"},
				{IDRef: "item6", Linear: "yes"},
			},
		},
		Guide: Guide{
			Refs: []GuideRef{
				{Type: "toc", Title: "Contents", Href: "page1.html#pgepubid00001"},
				{Type: "cover", Title: "CoverImageID", Href: "cover.html"},
			},
		},
	}
	got, err := Write(ncx)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(
		`<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="id">
    <metadata xmlns="http://www.idpf.org/2007/opf">
        <title xmlns="http://purl.org/dc/elements/1.1/">The Count of Monte Cristo, Illustrated</title>
        <identifier xmlns="http://purl.org/dc/elements/1.1/" id="id">http://www.gutenberg.org/ebooks/1184</identifier>
        <creator xmlns="http://purl.org/dc/elements/1.1/">Alexandre Dumas</creator>
        <publisher xmlns="http://purl.org/dc/elements/1.1/"></publisher>
        <language xmlns="http://purl.org/dc/elements/1.1/">en</language>
        <subject xmlns="http://purl.org/dc/elements/1.1/">Historical fiction</subject>
        <subject xmlns="http://purl.org/dc/elements/1.1/">Adventure stories</subject>
        <rights xmlns="http://purl.org/dc/elements/1.1/">Public domain</rights>
        <source xmlns="http://purl.org/dc/elements/1.1/">http://www.gutenberg.org/files/1184/1184-h/1184-h.htm</source>
        <date xmlns="http://purl.org/dc/elements/1.1/" event="publication">1998-01-01</date>
        <date xmlns="http://purl.org/dc/elements/1.1/" event="conversion">2019-05-31T14:42:00.544171+00:00</date>
        <meta xmlns="http://www.idpf.org/2007/opf" name="cover" content="item1"></meta>
        <meta xmlns="http://www.idpf.org/2007/opf" name="calibre:timestamp" content="2017-10-10T20:22:18.989147+00:00"></meta>
    </metadata>
    <manifest xmlns="http://www.idpf.org/2007/opf">
        <item xmlns="http://www.idpf.org/2007/opf" id="item1" href="cover.jpg" media-type="image/jpeg"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="item2" href="pgepub.css" media-type="text/css"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="item3" href="0.css" media-type="text/css"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="item4" href="1.css" media-type="text/css"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="item5" href="page1.html" media-type="application/xhtml+xml"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="item6" href="page2.html" media-type="application/xhtml+xml"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="ncx" href="toc.ncx" media-type="application/x-dtbncx+xml"></item>
        <item xmlns="http://www.idpf.org/2007/opf" id="coverpage-wrapper" href="cover.html" media-type="application/xhtml+xml"></item>
    </manifest>
    <spine xmlns="http://www.idpf.org/2007/opf" toc="ncx">
        <itemref xmlns="http://www.idpf.org/2007/opf" idref="coverpage-wrapper" linear="no"></itemref>
        <itemref xmlns="http://www.idpf.org/2007/opf" idref="item5" linear="yes"></itemref>
        <itemref xmlns="http://www.idpf.org/2007/opf" idref="item6" linear="yes"></itemref>
    </spine>
    <guide xmlns="http://www.idpf.org/2007/opf">
        <reference xmlns="http://www.idpf.org/2007/opf" type="toc" title="Contents" href="page1.html#pgepubid00001"></reference>
        <reference xmlns="http://www.idpf.org/2007/opf" type="cover" title="CoverImageID" href="cover.html"></reference>
    </guide>
</package>`)
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Error("got != want:\n", diff)
	}
}
