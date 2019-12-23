package book

import "errors"

// Model an ebook
type Book struct {
	Title        string
	Identifier   string
	Creators     []string
	Publisher    string
	Language     string
	Subjects     []string
	Rights       string
	Source       string
	Dates        []Date
	Metas        []Meta
	Resources    []Resource
	SpineItems   []SpineItem
	CoverImageID string
	TOCItems     []TOCItem
}

type Resource struct {
	ID        string
	Path      string
	MediaType string
	Contents  []byte
}

type SpineItem struct {
	Linear bool
	ID     string // Matches some Resource.ID
}

type TOCItem struct {
	ID        string // Matches some Resource.ID
	PlayOrder int
	Label     string
	Href      string // May have URL fragment
	Children  []TOCItem
}

type Meta struct {
	Name    string
	Content string
}

type Date struct {
	Event string
	Value string
}

func (b Book) GetResource(path string) (Resource, error) {
	for _, r := range b.Resources {
		if r.Path == path {
			return r, nil
		}
	}
	return Resource{}, errors.New("can't find resource at path " + path)
}
