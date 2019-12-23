package container

import (
	"encoding/xml"
)

var Path = "META-INF/container.xml"

// XML structure of META-INF/container.xml
type Container struct {
	XMLName   xml.Name   `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Version   string     `xml:"version,attr"`
	RootFiles []RootFile `xml:"rootfiles>rootfile"`
}

type RootFile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

func Read(xmlBytes []byte) (Container, error) {
	var container Container
	err := xml.Unmarshal(xmlBytes, &container)
	if err != nil {
		return Container{}, err
	}
	return container, nil
}

func Write(container Container) ([]byte, error) {
	xmlBytes, err := xml.MarshalIndent(container, "", "    ")
	if err != nil {
		return nil, err
	}
	return []byte(xml.Header + string(xmlBytes)), nil
}
