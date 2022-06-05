package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

func Start(canal chan uint64, id uint64, value uint64) {
	for {
		fmt.Println(id, " -- ", value)
		value += 1
		canal <- value
		time.Sleep(time.Millisecond * 500)
	}
}

func client() {
	var res [2]uint64
	c, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewDecoder(c).Decode(&res)

	if err != nil {
		fmt.Println(err)
	} else {
		canal := make(chan uint64)
		go Start(canal, res[0], res[1])

		for {
			res[1] = <-canal
			err := gob.NewEncoder(c).Encode(res[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func main() {
	go client()
	var input string
	fmt.Scanln(&input)
}
