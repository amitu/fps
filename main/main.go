package main

import (
	"github.com/amitu/fps"
	"flag"
	"fmt"
	"strings"
)	

func main() {
	nworkers := flag.Int("workers", 10, "Number of workers to use.")

	flag.Parse()	

	workers := fps.CreateWorkers("all", *nworkers)
	fmt.Printf("Started %d workers.\n", *nworkers)

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