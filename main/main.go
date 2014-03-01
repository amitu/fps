package main

import "github.com/amitu/fps"

func main() {
	workers := fps.CreateWorkers("all", 10)
	policy := fps.PolicyFromFile("policy.xml")

	fps.Server("127.0.0.1:8000", workers, policy)
	fps.Server("127.0.0.1:8001", workers, policy)

	fps.ServeForever()
}