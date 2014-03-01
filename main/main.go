package main

import (
	"github.com/amitu/fps"
	"flag"
	"fmt"
	"strings"
)	

func main() {
	nworkers := flag.Int("workers", 10, "Number of workers to use.")
	strict := flag.Bool("strict", false, "Should the server be strict?")

	flag.Parse()

	workers := fps.CreateWorkers("all", *nworkers, *strict)
	fmt.Printf("Started %d workers, strict=%t.\n", *nworkers, *strict)

	for _, server := range flag.Args() {
		parts := strings.SplitN(server, ":", 2)

		if len(parts) != 2 {
			fmt.Println("Bad argument:", server)
			return
		}

		fps.Server(parts[1], workers, fps.PolicyFromFile(parts[0]))
	}

	fps.ServeForever()
}