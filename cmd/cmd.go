package main

import (
	"github.com/cat-in-vacuum/golearn/algoritms"
	"github.com/cat-in-vacuum/golearn/examples"
	"github.com/cat-in-vacuum/golearn/expserv"
)

var isTraceEnabled = true

type T struct {
	S string
}

func main() {
	expserv.Run()
	examples.Run()
	algoritms.Run()
	/*go func() {
		time.Sleep(time.Second * 1)
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Err(err).Err(err).Msg("conn failure")


	//expserv.Run()
	//expserv.Run()

	// go examples.Run()
	*/
}

