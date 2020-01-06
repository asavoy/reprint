package cmd

import (
	"github.com/asavoy/reprint/clean"
	"github.com/asavoy/reprint/epub"
)

func Run(inPath, outPath string) error {
	book, err := epub.Read(inPath)
	if err != nil {
		return err
	}
	err = clean.Clean(&book)
	if err != nil {
		return err
	}
	err = epub.Write(outPath, book)
	if err != nil {
		return err
	}
	return nil
}
