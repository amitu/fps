package fps

import (
	"os"
	"os/signal"
	"sync"
	"io/ioutil"
)

var (
	quitServers bool = false
	quitWorkers bool = false
	quitLock sync.Mutex
)

func waitForCtrlC() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<- sig

	signal.Stop(sig)
}

func timeToQuitServers() bool {
	quitLock.Lock()
	defer quitLock.Unlock()

	return quitServers
}

func timeToQuitWorkers() bool {
	quitLock.Lock()
	defer quitLock.Unlock()

	return quitWorkers
}

func beginQuitServers() {
	quitLock.Lock()
	defer quitLock.Unlock()
	
	quitServers = true	
}

func beginQuitWorkers() {
	quitLock.Lock()
	defer quitLock.Unlock()
	
	quitWorkers = true	
}

func PolicyFromFile(f string) []byte {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	return b
}
