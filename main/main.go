package main

import "github.com/amitu/fps"

func main() {
	workers := fps.InitWorkers(10)

	fps.Server("127.0.0.1:8000", workers)
	fps.Server("127.0.0.1:8001", workers)

	fps.ServeForever(workers)
}