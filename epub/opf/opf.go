package opf

import (
	"encoding/xml"
)

// XML structure of content.opf
type Package struct {
	XMLName          xml.Name `xml:"http://www.idpf.org/2007/opf package"`
	Version          string   `xml:"version,attr"`
	UniqueIdentifier string   `xml:"unique-identifier,attr"`
	Metadata         Metadata `xml:"http://www.idpf.org/2007/opf metadata"`
	Manifest         Manifest `xml:"http://www.idpf.org/2007/opf manifest"`
	Spine            Spine    `xml:"http://www.idpf.org/2007/opf spine"`
	Guide            Guide    `xml:"http://www.idpf.org/2007/opf guide"`
}

type Metadata struct {
	XMLName     xml.Name     `xml:"http://www.idpf.org/2007/opf metadata"`
	Title       string       `xml:"http://purl.org/dc/elements/1.1/ title"`
	Identifiers []Identifier `xml:"http://purl.org/dc/elements/1.1/ identifier"`
	Creators    []string     `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Publisher   string       `xml:"http://purl.org/dc/elements/1.1/ publisher"`
	Language    string       `xml:"http://purl.org/dc/elements/1.1/ language"`
	Subjects    []string     `xml:"http://purl.org/dc/elements/1.1/ subject"`
	Rights      string       `xml:"http://purl.org/dc/elements/1.1/ rights"`
	Source      string       `xml:"http://purl.org/dc/elements/1.1/ source"`
	Dates       []Date       `xml:"http://purl.org/dc/elements/1.1/ date"`
	Metas       []Meta       `xml:"http://www.idpf.org/2007/opf meta"`
}

type Meta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}

type Identifier struct {
	ID    string `xml:"id,attr"`
	Value string `xml:",innerxml"`
}

type Date struct {
	Event string `xml:"event,attr"`
	Value string `xml:",innerxml"`
}

type Manifest struct {
	XMLName xml.Name       `xml:"http://www.idpf.org/2007/opf manifest"`
	Items   []ManifestItem `xml:"http://www.idpf.org/2007/opf item"`
}

type ManifestItem struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	XMLName  xml.Name       `xml:"http://www.idpf.org/2007/opf spine"`
	Toc      string         `xml:"toc,attr"`
	ItemRefs []SpineItemRef `xml:"http://www.idpf.org/2007/opf itemref"`
}

type SpineItemRef struct {
	IDRef  string `xml:"idref,attr"`
	Linear string `xml:"linear,attr"`
}

type Guide struct {
	XMLName xml.Name   `xml:"http://www.idpf.org/2007/opf guide"`
	Refs    []GuideRef `xml:"http://www.idpf.org/2007/opf reference"`
}

type GuideRef struct {
	Type  string `xml:"type,attr"`
	Title string `xml:"title,attr"`
	Href  string `xml:"href,attr"`
}

var MediaType = "application/oebps-package+xml"

func Read(xmlBytes []byte) (Package, error) {
	var pack Package
	err := xml.Unmarshal(xmlBytes, &pack)
	if err != nil {
		return Package{}, err
	}
	return pack, nil
}

func Write(pack Package) ([]byte, error) {
	xmlBytes, err := xml.MarshalIndent(pack, "", "    ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(xmlBytes)), nil
}
