package main

import (
	"fmt"
	"os"
	"path"

	"github.com/asavoy/reprint/cmd"
)

func main() {
	switch len(os.Args) {
	case 1:
		fmt.Printf("usage: %s source.epub [fixed.epub]\n", os.Args[0])
	case 2:
		inPath := os.Args[1]
		outDir := path.Dir(inPath)
		outFilename := fmt.Sprintf("%s.reprint.epub", path.Base(inPath))
		outPath := path.Join(outDir, outFilename)
		if inPath == outPath {
			fmt.Println("error: inPath and outPath are the same!")
		}
		err := cmd.Run(inPath, outPath)
		if err != nil {
			fmt.Println("error:", err)
		}
	case 3:
		inPath := os.Args[1]
		outPath := os.Args[2]
		fmt.Println("error: inPath and outPath are the same!")
		err := cmd.Run(inPath, outPath)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}
