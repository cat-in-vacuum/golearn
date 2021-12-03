package expserv

import (
	"encoding/json"
	"fmt"
	"github.com/cat-in-vacuum/golearn/expclient"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"path"
	"time"
)

var stopSrvChan = make(chan struct{})

func Run() {

	startHTTP()
	log.Debug().Msg("test http  server has benn started")

}

func startHTTP() {
	client := expclient.New(&http.Client{})

	go func() {
		<-stopSrvChan
		log.Debug().Msg("app has been stopped")
		os.Exit(1)
	}()
	srv := mux.NewRouter()


	srv.HandleFunc("/debug/pprof/", pprof.Index)
	srv.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	srv.HandleFunc("/debug/pprof/profile", pprof.Profile)
	srv.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	srv.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	srv.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	srv.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	srv.Handle("/debug/pprof/block", pprof.Handler("block"))

	srv.Handle("/", rootHandler())
	srv.Handle("/test", handleTest()).Methods(http.MethodGet)
	srv.Handle("/test/off", handleTestServerOff()).Methods(http.MethodGet)
	srv.Handle("/test/post/{[0-9]+}", getPosts(client)).Methods(http.MethodGet)


	log.Fatal().Err(http.ListenAndServe(":8080", srv))
}

func handleTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	}
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("route : /")
		}

		log.Debug().Interface("payload", payload).Msg("payload")


		log.Debug().Interface("http_req", r).Msg("payload")
	}
}

func getPosts(client *expclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := path.Base(r.URL.Path)

		resp, err := client.GetPosts(id)
		if err != nil {
			log.Error().Err(err).Msg("")
			return
		}

		outPayload, err := json.Marshal(resp)
		if err != nil {
			log.Error().Err(err).Msg("")
		}

		w.Write(outPayload)
	}
}

func handleTestServerOff() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "shutting down..."}`))
		stopSrvChan <- struct{}{}
	}
}

func EmulateOne() {
	time.Sleep(time.Second*1)
	fmt.Println("EmulateOne")
}

func EmulateTwo() {
	time.Sleep(time.Second*1)
	fmt.Println("EmulateTwo")
}

func EmulateThree() {

}