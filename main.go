package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asavoy/reprint/cmd"
)

func main() {
	fmt.Println(os.Args)
	switch os.Args[1] {
	case "run":
		err := cmd.Run(os.Args[2], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
	case "check":
		err := cmd.Check(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unexpected argument")
	}
}
