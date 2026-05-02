// Command cooksense-server is the HTTP API server for the CookSense application.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate", "seed":
			fmt.Printf("%s: not implemented\n", os.Args[1])
			return
		}
	}
	fmt.Println("cooksense-server starting")
}
