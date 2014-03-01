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

const timeout time.Duration = 1e8

func CreateWorkers(name string, n int) (workers chan net.Conn) {
	workers = make(chan net.Conn)

	for i := 0; i < n; i++ {
		go worker(name, i, workers)
	}
	
	return 
}

func ServeForever() {
	waitForCtrlC()

	fmt.Println("Got Ctrl-C, Time to quit.")

	beginQuitServers()
	wgServers.Wait()

	beginQuitWorkers()
	wgWorkers.Wait()

	fmt.Println("All done, goodbye.")
}

func Server(hostPort string, workers chan net.Conn) {
	wgServers.Add(1)
	go func() {
		addr, err := net.ResolveTCPAddr("tcp", hostPort)
		if err != nil {
			fmt.Println(err)
			return
		}

		server, err := net.ListenTCP("tcp", addr)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Listening on", hostPort)

		for {
			server.SetDeadline(time.Now().Add(timeout))
			conn, err := server.Accept()

			if err != nil {


				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					if timeToQuitServers() {						
						server.Close()
						fmt.Println("Server Down,", hostPort)
						wgServers.Done()
						return
					} 
					continue
				}

				fmt.Println(err)
				continue
			}

			fmt.Println("New connection on", hostPort)

			workers <- conn;
		}
	}()
}

func worker(name string, i int, workers chan net.Conn) {
	wgWorkers.Add(1)
	for {
		select {
		case conn := <- workers:
			fmt.Println("Worker", name, i)
			fmt.Fprintf(
				conn, "\nHello there!\n\n%s:%d welcomes you.\n\n", name, i,
			)
			conn.Close()
		case <- time.After(timeout):
			if timeToQuitWorkers() {
				fmt.Println("Worker done", name, i)
				wgWorkers.Done()	
				return				
			}
		}
	}
}