package epub

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/asavoy/reprint/book"
	"github.com/asavoy/reprint/epub/container"
	"github.com/asavoy/reprint/epub/ncx"
	"github.com/asavoy/reprint/epub/opf"
)

// Write an EPUBv2 format book
func Write(filepath string, b book.Book) error {
	zipFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	var fileWriter io.Writer

	// Build mimetype file
	mimetypeContents := []byte("application/epub+zip")

	// Build container.xml
	opfPath := "content.opf"
	containerContents, err := container.Write(buildContainer(opfPath))
	if err != nil {
		return err
	}

	// Build the EPUB data
	resources := b.Resources

	// Build toc.ncx
	ncxContents, err := ncx.Write(buildNCX(b))
	if err != nil {
		return err
	}
	resources = append(resources, book.Resource{
		ID:        "ncx",
		Path:      "toc.ncx",
		MediaType: "application/x-dtbncx+xml",
		Contents:  ncxContents,
	})

	// Build content.opf
	var metas []opf.Meta
	if b.CoverImageID != "" {
		metas = append(metas, opf.Meta{
			Name:    "cover",
			Content: b.CoverImageID,
		})
	}
	var dates []opf.Date
	for _, bookDate := range b.Dates {
		dates = append(dates, opf.Date{
			Event: bookDate.Event,
			Value: bookDate.Value,
		})
	}
	identifierName := "PrimaryIdentifier"
	opfContents, err := opf.Write(opf.Package{
		Version:          "2.0",
		UniqueIdentifier: identifierName,
		Metadata: opf.Metadata{
			Title: b.Title,
			Identifiers: []opf.Identifier{
				{ID: identifierName, Value: b.Identifier},
			},
			Creators:  b.Creators,
			Publisher: b.Publisher,
			Language:  b.Language,
			Subjects:  b.Subjects,
			Rights:    b.Rights,
			Source:    b.Source,
			Metas:     metas,
			Dates:     dates,
		},
		Manifest: opf.Manifest{
			Items: buildManifestItems(resources),
		},
		Spine: opf.Spine{
			Toc:      "ncx",
			ItemRefs: buildSpineItemRefs(b.SpineItems),
		},
		// Don't bother with guide because Apple Books doesn't support it
		Guide: opf.Guide{},
	})
	if err != nil {
		return err
	}

	// Write mimetype file
	fileWriter, err = zipWriter.CreateHeader(&zip.FileHeader{
		Name:     "mimetype",
		Method:   zip.Store,
		Modified: time.Now(),
	})
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(mimetypeContents)
	if err != nil {
		return err
	}

	// Write container.xml
	fileWriter, err = zipWriter.Create(container.Path)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(containerContents)
	if err != nil {
		return err
	}

	// Write content.opf
	fileWriter, err = zipWriter.Create(opfPath)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(opfContents)
	if err != nil {
		return err
	}

	// Write all resources, including toc.ncx
	for _, resource := range resources {
		fileWriter, err = zipWriter.Create(resource.Path)
		if err != nil {
			return err
		}
		_, err = fileWriter.Write(resource.Contents)
		if err != nil {
			return err
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func buildContainer(opfPath string) container.Container {
	return container.Container{
		Version: "1.0",
		RootFiles: []container.RootFile{
			{
				FullPath:  opfPath,
				MediaType: "application/oebps-package+xml",
			},
		},
	}
}

func buildNCX(b book.Book) ncx.NCX {
	maxDepth := 0
	for _, tocItem := range b.TOCItems {
		depth := getMaxDepth(tocItem)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return ncx.NCX{
		Version: "2005-1",
		Metas: []ncx.Meta{
			{Name: "dtb:uid", Content: b.Identifier},
			{Name: "dtb:depth", Content: fmt.Sprintf("%d", maxDepth)},
			{Name: "dtb:totalPageCount", Content: "0"},
			{Name: "dtb:maxPageNumber", Content: "0"},
		},
		Title:     b.Title,
		Author:    strings.Join(b.Creators, ", "),
		NavPoints: buildNavPoints(b.TOCItems),
	}
}

func getMaxDepth(tocItem book.TOCItem) int {
	maxDepth := 1
	for _, child := range tocItem.Children {
		depth := 1 + getMaxDepth(child)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

func buildManifestItems(resources []book.Resource) []opf.ManifestItem {
	var manifestItems []opf.ManifestItem
	for _, resource := range resources {
		manifestItems = append(manifestItems, opf.ManifestItem{
			ID:        resource.ID,
			Href:      resource.Path,
			MediaType: resource.MediaType,
		})
	}
	return manifestItems
}

func buildSpineItemRefs(spineItems []book.SpineItem) []opf.SpineItemRef {
	var spineItemRefs []opf.SpineItemRef
	for _, spineItem := range spineItems {
		var linear string
		if spineItem.Linear {
			linear = "yes"
		} else {
			linear = "no"
		}
		spineItemRefs = append(spineItemRefs, opf.SpineItemRef{
			IDRef:  spineItem.ID,
			Linear: linear,
		})
	}
	return spineItemRefs
}

func buildNavPoints(tocItems []book.TOCItem) []ncx.NavPoint {
	var navPoints []ncx.NavPoint
	for _, item := range tocItems {
		children := buildNavPoints(item.Children)
		navPoints = append(navPoints, ncx.NavPoint{
			ID:        item.ID,
			PlayOrder: fmt.Sprintf("%d", item.PlayOrder),
			Label:     item.Label,
			Content:   ncx.Content{Src: item.Href},
			NavPoints: children,
		})
	}
	return navPoints
}
