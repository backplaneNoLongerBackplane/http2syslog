package main

import (
	"bufio"
	"flag"
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

	/*
		App logs go to stdout and are prefixed with the line numbers.
	*/
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	/*
		Dial the syslog server and get a writer-logger...
	*/
	log.Printf("Dialing syslog server: %s://%s\n", *network, *addr)
	sysLog, err := syslog.Dial(*network, *addr, syslog.LOG_INFO|syslog.LOG_USER, *tag)
	if err != nil {
		log.Fatal("failed to dial syslog")
	}
	defer sysLog.Close()

	/*
		...then listen on HTTP and handle log requests.
	*/
	log.Printf("Listening for logs over http at http://localhost:%s\n", *port)
	http.ListenAndServe(":"+*port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			Ignore all requests that do not POST a body...
		*/
		if r.Body == nil || r.Method != "POST" {
			return
		}
		defer r.Body.Close()

		/*
			...otherwise scan each line as a log message...
		*/
		scanner := bufio.NewScanner(r.Body)
		for scanner.Scan() {
			msg := scanner.Text()
			log.Println(msg)
			/*
				...and send each message to the syslog server.
			*/
			if err := sysLog.Info(msg); err != nil {
				log.Println("error: failed to send log message:", err)
			}
		}
		/*
			If there are any errors encountered scanning the request body, log them here.
		*/
		if err := scanner.Err(); err != nil {
			log.Println("error: failed to scan log messages:", err)
			if err := sysLog.Err("http2syslog: " + err.Error()); err != nil {
				log.Println("error: failed to send error message:", err)
			}
		}
	}))
}
