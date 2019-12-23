package container

import (
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
	containerXML := []byte(
		`<?xml version="1.0"?>
		<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
			<rootfiles>
				<rootfile full-path="content.opf" media-type="application/oebps-package+xml"/>
			</rootfiles>
		</container>`)
	got, err := Read(containerXML)
	if err != nil {
		t.Fatal(err)
	}
	want := Container{
		Version: "1.0",
		XMLName: xml.Name{
			Space: "urn:oasis:names:tc:opendocument:xmlns:container",
			Local: "container",
		},
		RootFiles: []RootFile{
			{
				FullPath:  "content.opf",
				MediaType: "application/oebps-package+xml",
			},
		}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Error("got != want:\n", diff)
	}
}

func TestWrite(t *testing.T) {
	container := Container{
		Version: "1.0",
		XMLName: xml.Name{
			Space: "urn:oasis:names:tc:opendocument:xmlns:container",
			Local: "container",
		},
		RootFiles: []RootFile{
			{
				FullPath:  "content.opf",
				MediaType: "application/oebps-package+xml",
			},
		},
	}
	got, err := Write(container)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte(
		`<?xml version="1.0" encoding="UTF-8"?>
<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
    <rootfiles>
        <rootfile full-path="content.opf" media-type="application/oebps-package+xml"></rootfile>
    </rootfiles>
</container>`)
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Error("got != want:\n", diff)
	}
}
