# HTTP-to-SYSLOG Service

This is a small program that will accept log messages over http and proxy them to a syslog server.

### Usage

```sh
$ PORT=5000 go run main.go -proto "udp" -syslog "logsN.papertrailapp.com:XXXXX"
$ curl -X POST -d '{"Foo":"Bar"}' http://localhost:5000
```