# HTTP-to-SYSLOG Service

This is a small example program that will accept log messages over http and proxy them to a remote syslog server (such as Papertrail).

### Usage

```sh
$ go run main.go -port "5000" -network "udp" -addr "logsN.papertrailapp.com:XXXXX" -tag "myapp"
$ curl -X POST -d '{"Foo":"Bar"}' http://localhost:5000
```