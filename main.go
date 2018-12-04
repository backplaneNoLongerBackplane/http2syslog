package main

import (
	"flag"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
)

var (
	port    = flag.String("port", os.Getenv("PORT"), "The http port on which to listen for log messages.")
	network = flag.String("network", "", "Network (eg: udp) to send syslog messages. Default network is the local syslog server.")
	addr    = flag.String("addr", "", "Address (host:port) of the destination syslog server.")
	tag     = flag.String("tag", "", "If left empty, the process name will be used.")
)

func main() {
	flag.Parse()

	log.Printf("Dialing %s://%s\n", *network, *addr)
	logger, err := syslog.Dial(*network, *addr, syslog.LOG_EMERG|syslog.LOG_KERN, *tag)
	if err != nil {
		log.Fatal("failed to dial syslog")
	}
	defer logger.Close()

	log.Printf("Listening at http://localhost:%s\n", *port)
	http.ListenAndServe(":"+*port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			return
		}
		defer r.Body.Close()

		msg, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Err(err.Error())
		}
		log.Println(string(msg))
		if err := logger.Info(string(msg)); err != nil {
			log.Println("error: ", err)
		}
	}))
}
