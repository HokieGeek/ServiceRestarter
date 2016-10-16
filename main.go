package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func restartService(log *log.Logger, serviceName string) {
	log.Printf("Restarting %s\n", serviceName)
	cmd := exec.Command("/usr/bin/systemctl", "restart", serviceName)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	} else {
		log.Printf("systemctl: %s\n", string(out))
	}
}

func Handler(log *log.Logger, id string, w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}

	if err := r.Body.Close(); err != nil {
		log.Panic(err)
	}

	// TODO: need to check that id and service are available
	if r.URL.Query()["id"][0] == id {
		restartService(log, r.URL.Query()["service"][0])
	} else {
		log.Printf("Ignoring request with id: %s\n", r.URL.Query()["id"])
	}
}

func main() {
	// logger := log.New(os.Stdout, "", log.Lshortfile)
	logger := log.New(os.Stdout, "", 0)
	portPtr := flag.Int("port", 8081, "Specify the port to use")
	flag.Parse()

	if flag.NArg() < 1 {
		logger.Panic("Not enough arguments")
	}
	id := flag.Arg(0)

	logger.Printf("Using identifier: %s\n", id)

	mux := http.NewServeMux()

	mux.HandleFunc("/ServiceRestart",
		func(w http.ResponseWriter, r *http.Request) {
			Handler(logger, id, w, r)
		})

	http.ListenAndServe(fmt.Sprintf(":%d", *portPtr), mux)
}

// vim: set foldmethod=marker:
