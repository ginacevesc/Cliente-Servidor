package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func servery(Can []chan uint64, outcha []chan bool) {

	catch := make(chan bool)
	asign := make(chan uint64, 2)
	post := make(chan net.Conn)

	s, err := net.Listen("tcp", ":9999")

	if err != nil {

		fmt.Println(err)

		return

	}
	go master(Can, outcha, asign, post, catch)

	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		post <- c
	}
}

func handleClient(c net.Conn, tunel chan uint64, out chan bool, asign chan uint64, catch chan bool) {

	out <- true
	items := [2]uint64{<-tunel, <-tunel}
	err := gob.NewEncoder(c).Encode(items)

	if err != nil {
		fmt.Println(err)
	}

	for {

		err := gob.NewDecoder(c).Decode(&items[1])

		if err != nil {
			asign <- items[0]
			asign <- items[1]
			catch <- true
			return
		}
	}
}

func Proceso(id uint64, i uint64, tunel chan uint64, out chan bool) {
	for {
		select {
		case <-out:
			tunel <- id
			tunel <- i
			return
		default:
			fmt.Println(id, " .- ", i)
			i += 1
			if id == 4 {
				fmt.Println("~~~~~~~~~~~~")
			}
			time.Sleep(time.Millisecond * 500)
		}

	}

}

func master(Can []chan uint64, outcha []chan bool, asign chan uint64, post chan net.Conn, catch chan bool) {
	for {
		select {
		case <-catch:
			items := [2]uint64{<-asign, <-asign}
			Can = append(Can, make(chan uint64, 2))
			outcha = append(outcha, make(chan bool))
			go Proceso(items[0], items[1], Can[len(Can)-1], outcha[len(outcha)-1])
		case c := <-post:
			go handleClient(c, Can[0], outcha[0], asign, catch)
			Can = Can[1:]
			outcha = outcha[1:]

		}
	}
}

func main() {
	var idCount uint64 = 0
	var input string
	Can := make([]chan uint64, 5)
	outcha := make([]chan bool, 5)

	for i := 0; i < 5; i++ {
		Can[i] = make(chan uint64, 2)
		outcha[i] = make(chan bool)
		go Proceso(idCount, 0, Can[i], outcha[i])
		idCount++
	}

	go servery(Can, outcha)

	fmt.Scanln(&input)

}
