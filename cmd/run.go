package cmd

import (
	"github.com/asavoy/reprint/clean"
	"github.com/asavoy/reprint/epub"
)

func Run(inFile, outFile string) error {
	book, err := epub.Read(inFile)
	if err != nil {
		return err
	}
	err = clean.Clean(&book)
	if err != nil {
		return err
	}
	err = epub.Write(outFile, book)
	if err != nil {
		return err
	}
	return nil
}