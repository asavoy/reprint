package ncx

import (
	"encoding/xml"
)

var (
	DocType = `<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
`
	MediaType = "application/x-dtbncx+xml"
)

// XML structure of toc.ncx
type NCX struct {
	XMLName   xml.Name   `xml:"http://www.daisy.org/z3986/2005/ncx/ ncx"`
	Version   string     `xml:"version,attr"`
	Metas     []Meta     `xml:"head>meta"`
	Title     string     `xml:"docTitle>text"`
	Author    string     `xml:"docAuthor>text"`
	NavPoints []NavPoint `xml:"navMap>navPoint"`
}

type Meta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}
type NavPoint struct {
	ID        string     `xml:"id,attr"`
	PlayOrder string     `xml:"playOrder,attr"`
	Label     string     `xml:"navLabel>text"`
	Content   Content    `xml:"content"`
	NavPoints []NavPoint `xml:"navPoint"`
}

type Content struct {
	Src string `xml:"src,attr"`
}

func Read(xmlBytes []byte) (NCX, error) {
	var ncx NCX
	err := xml.Unmarshal(xmlBytes, &ncx)
	if err != nil {
		return NCX{}, err
	}
	return ncx, nil
}

func Write(ncx NCX) ([]byte, error) {
	xmlBytes, err := xml.MarshalIndent(ncx, "", "    ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + DocType + string(xmlBytes)), nil
}
