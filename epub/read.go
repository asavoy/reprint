package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/asavoy/reprint/book"
	"github.com/asavoy/reprint/epub/container"
	"github.com/asavoy/reprint/epub/ncx"
	"github.com/asavoy/reprint/epub/opf"
)

// Read an EPUBv2 format book
func Read(filepath string) (book.Book, error) {
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return book.Book{}, err
	}
	defer r.Close()

	ctr, err := container.Read(readFile(r, container.Path))
	if err != nil {
		return book.Book{}, err
	}
	opfPath := ctr.RootFiles[0].FullPath
	pack, err := opf.Read(readFile(r, opfPath))
	if err != nil {
		return book.Book{}, err
	}
	fmt.Println(pack.Metadata.Title)

	tocResource, err := parseTOCResource(pack.Manifest.Items, r, opfPath)
	if err != nil {
		return book.Book{}, nil
	}
	tocNCX, err := ncx.Read(tocResource.Contents)
	if err != nil {
		return book.Book{}, err
	}
	resources := parseResources(pack.Manifest.Items, r, opfPath)
	tocItems, err := parseTOCItems(tocNCX.NavPoints, tocResource.Path)
	if err != nil {
		return book.Book{}, err
	}

	spineItems, err := parseSpineItems(pack.Spine.ItemRefs)
	if err != nil {
		return book.Book{}, err
	}

	var dates []book.Date
	for _, date := range pack.Metadata.Dates {
		dates = append(dates, book.Date{
			Event: date.Event,
			Value: date.Value,
		})
	}

	var metas []book.Meta
	for _, meta := range pack.Metadata.Metas {
		metas = append(metas, book.Meta{
			Name:    meta.Name,
			Content: meta.Content,
		})
	}

	uniqueID, err := parseUniqueID(pack)
	if err != nil {
		return book.Book{}, err
	}

	coverImageID, err := parseCoverImageID(pack.Metadata.Metas, pack.Manifest.Items)
	if err != nil {
		return book.Book{}, err
	}

	b := book.Book{
		Title:        pack.Metadata.Title,
		Identifier:   uniqueID,
		Creators:     pack.Metadata.Creators,
		Publisher:    pack.Metadata.Publisher,
		Language:     pack.Metadata.Language,
		Subjects:     pack.Metadata.Subjects,
		Rights:       pack.Metadata.Rights,
		Source:       pack.Metadata.Source,
		Dates:        dates,
		Metas:        metas,
		Resources:    resources,
		SpineItems:   spineItems,
		CoverImageID: coverImageID,
		TOCItems:     tocItems,
	}

	return b, nil
}

func absPath(docPath string, relPath string) string {
	return path.Clean(path.Join(path.Dir(docPath), relPath))
}

func readFile(r *zip.ReadCloser, path string) []byte {
	for _, f := range r.File {
		if f.Name == path {
			return readAll(f)
		}
	}
	panic(fmt.Sprintf("could not find in zip %s", path))
}

func readAll(file *zip.File) []byte {
	fc, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		panic(err)
	}

	return content
}

func parseTOCResource(items []opf.ManifestItem, r *zip.ReadCloser, opfPath string) (book.Resource, error) {
	for _, item := range items {
		itemPath := absPath(opfPath, item.Href)
		if item.MediaType == ncx.MediaType {
			// This is the toc.ncx file, which we treat separately as metadata
			return book.Resource{
				ID:        item.ID,
				Path:      itemPath,
				MediaType: item.MediaType,
				Contents:  readFile(r, itemPath),
			}, nil
		}
	}
	return book.Resource{}, errors.New("missing toc.ncx")
}

func parseResources(items []opf.ManifestItem, r *zip.ReadCloser, opfPath string) []book.Resource {
	var resources []book.Resource
	for _, item := range items {
		itemPath := absPath(opfPath, item.Href)
		if item.MediaType == ncx.MediaType {
			// This is the toc.ncx file, which we treat separately as metadata
		} else if item.MediaType == opf.MediaType {
			// This is the content.opf file, which we treat separately as metadata
		} else {
			resources = append(resources, book.Resource{
				ID:        item.ID,
				Path:      itemPath,
				MediaType: item.MediaType,
				Contents:  readFile(r, itemPath),
			})
		}
	}
	return resources
}

func parseSpineItems(itemRefs []opf.SpineItemRef) ([]book.SpineItem, error) {
	var spineItems []book.SpineItem
	for _, item := range itemRefs {
		var linear bool
		if item.Linear == "yes" || item.Linear == "" {
			linear = true
		} else if item.Linear == "no" {
			linear = false
		} else {
			return nil, fmt.Errorf("unexpected value for linear: %s", item.Linear)
		}
		spineItems = append(spineItems, book.SpineItem{
			ID:     item.IDRef,
			Linear: linear,
		})
	}
	return spineItems, nil
}

func parseTOCItems(navPoints []ncx.NavPoint, tocPath string) ([]book.TOCItem, error) {
	var tocItems []book.TOCItem
	for _, np := range navPoints {
		playOrder, err := strconv.Atoi(np.PlayOrder)
		if err != nil {
			return nil, err
		}
		children, err := parseTOCItems(np.NavPoints, tocPath)
		if err != nil {
			return nil, err
		}
		tocItems = append(tocItems, book.TOCItem{
			ID:        np.ID,
			PlayOrder: playOrder,
			Label:     np.Label,
			Href:      absPath(tocPath, np.Content.Src),
			Children:  children,
		})
	}
	return tocItems, nil
}

func parseUniqueID(pack opf.Package) (string, error) {
	for _, identifier := range pack.Metadata.Identifiers {
		if identifier.ID == pack.UniqueIdentifier {
			return identifier.Value, nil
		}
	}
	return "", errors.New("can't find unique identifier")
}

func parseCoverImageID(metas []opf.Meta, manifestItems []opf.ManifestItem) (string, error) {
	var coverMeta *opf.Meta
	for _, meta := range metas {
		if meta.Name == "cover" {
			coverMeta = &meta
			break
		}
	}
	if coverMeta == nil {
		return "", nil
	}
	for _, item := range manifestItems {
		if item.Href == coverMeta.Content || item.ID == coverMeta.Content {
			return item.ID, nil
		}
	}
	return "", errors.New("can't find cover manifest item")
}
