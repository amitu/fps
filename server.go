package fps

import (
	"fmt"
	"net"
	"time"
	"sync"
)

var (
	wgWorkers, wgServers sync.WaitGroup
)

const timeout time.Duration = 2e8
const network_timeout time.Duration = 1e10

func CreateWorkers(name string, n int, strict bool) (workers chan work) {
	workers = make(chan work)

	for i := 0; i < n; i++ {
		go worker(name, i, workers, strict)
	}
	
	return 
}

func Shutdown() {
	beginQuitServers()
	wgServers.Wait()

	beginQuitWorkers()
	wgWorkers.Wait()

	fmt.Println("All done, goodbye.")
}

func ServeForever() {
	WaitForCtrlC()

	fmt.Println("Got Ctrl-C, time to quit.")
	Shutdown()
}

type work struct {
	conn net.Conn
	policy []byte
}

func Server(hostPort string, workers chan work, policy []byte) {
	wgServers.Add(1)
	go func() {
		addr, err := net.ResolveTCPAddr("tcp", hostPort)
		if err != nil {
			panic(err)
		}

		server, err := net.ListenTCP("tcp", addr)

		if err != nil {
			panic(err)
		}

		fmt.Println("Listening on", hostPort)

		for {
			server.SetDeadline(time.Now().Add(timeout))
			conn, err := server.Accept()

			if err != nil {


				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					if timeToQuitServers() {						
						server.Close()
						fmt.Println("Server down,", hostPort)
						wgServers.Done()
						return
					} 
					continue
				}

				fmt.Println(err)
				continue
			}

			fmt.Println("New connection on", hostPort)

			workers <- work{conn, policy};
		}
	}()
}

func worker(name string, i int, workers chan work, strict bool) {
	wgWorkers.Add(1)
	for {
		select {
		case work := <- workers:
			fmt.Println("Worker", name, i)
			if strict {
				// TODO:
				// should read data from socket 
				// and do some validation on it
				// strict client sends this first:
				// 		<policy-file-request/>\0
			}
			work.conn.SetDeadline(time.Now().Add(network_timeout))
			_, err := work.conn.Write(work.policy)
			if err != nil {
				fmt.Println(err)
			}
			work.conn.Close()
		case <- time.After(timeout):
			if timeToQuitWorkers() {
				fmt.Println("Worker done", name, i)
				wgWorkers.Done()	
				return				
			}
		}
	}
}