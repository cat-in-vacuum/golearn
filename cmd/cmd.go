package main

import (
	"github.com/cat-in-vacuum/golearn/examples"
)

var isTraceEnabled = true

type T struct {
	S string
}

func main() {
	// expserv.Run()

	/*go func() {
		time.Sleep(time.Second * 1)
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Err(err).Err(err).Msg("conn failure")

		} else {
			log.Debug().Err(err).Msg("conn failure")
			conn.Write([]byte("ok"))
		}
	}()*/

	//expserv.Run()

	// go examples.Run()
	examples.Run()
}
