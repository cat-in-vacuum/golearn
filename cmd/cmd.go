package main

import (
	"github.com/cat-in-vacuum/golearn/klepman"
)

var isTraceEnabled = true

type T struct {
	S string
}

func main() {
	klepman.Run()
	//os.RunLock()
	//expserv.Run()
	//examples.Run()
	//algoritms.Run()
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

