package main

import (
  "flag"
  "fmt"
  "github.com/golang/glog"
  "github.com/rboccuzzi/ernestine/relay"
  "net"
)

var numberHandled = 0

const echoBufferSize = 1024

// Callback function that implements my server, doing whatever it needs to do
// Which in this case is "echo back", which incidentally I needed to make
// that CopyStreams functionality as part of the relay server. /Shrug
func echoServerFunctionality(conn net.Conn) {
  numberHandled++
  glog.V(1).Infof("echoing %d ringy dingy\n", numberHandled)
  defer conn.Close()
  relay.CopyStreams(conn, conn)
}

func main() {
  var relayPort int
  var relayHost string
  flag.IntVar(&relayPort, "port", 9100, "The port to listen to connections for new services to relay")
  flag.StringVar(&relayHost, "host", "localhost", "he host running the relay server")
  flag.Parse()

  connectionString := fmt.Sprintf("%s:%d", relayHost, relayPort)
  relay.RelayServiceRequest(connectionString, echoServerFunctionality)
}
