package fps

import (
	"fmt"
	"net"
	"time"
	"os"
	"os/signal"
	"sync"
)

var (
	quit chan bool = make(chan bool, 2)
	wg sync.WaitGroup
)

func InitWorkers(n int) (workers chan net.Conn) {
	workers = make(chan net.Conn)

	for i := 0; i < n; i++ {
		go Worker(i, workers)
	}
	
	return 
}

func ServeForever(workerss... chan net.Conn) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<- sig

	signal.Stop(sig)

	fmt.Println("Got signal")

	quit <- true

	wg.Wait()

	<- quit

	for _, workers := range workerss {
		workers <- nil
		<- workers
	}
}

func Server(hostPort string, workers chan net.Conn) {
	wg.Add(1)
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
			server.SetDeadline(time.Now().Add(1e8))
			conn, err := server.Accept()

			if err != nil {


				if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
					select {
					case <- quit:
						server.Close()
						fmt.Println("Server Down,", hostPort)
						wg.Done()
						quit <- true
						return
					default:
						continue
					}
				}

				fmt.Println(err)
				continue
			}

			fmt.Println("New connection on", hostPort)

			workers <- conn;
		}
	}()
}

func Worker(i int, workers chan net.Conn) {
	for {
		conn := <- workers

		if conn == nil {
			fmt.Println("Worker done", i)
			workers <- nil
			return
		}

		fmt.Println("Worker", i)

		conn.Write([]byte("Hello there!\n"))
		conn.Close()
	}
}