package main

import (
  "flag"
  "fmt"
  "github.com/golang/glog"
  "github.com/rboccuzzi/ernestine/relay"
  "net"
)

func main() {
  var masterPort int
  var servicePort int
  var relayPort int
  flag.IntVar(&masterPort, "port", 9100, "The port to listen to connections for new services to relay")
  flag.IntVar(&servicePort, "service", 9200, "Where to start the service port range")
  flag.IntVar(&relayPort, "relay", 9300, "Where to start the relay port range")
  flag.Parse()

  glog.V(1).Info("Listening on Port: ", masterPort)
  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", masterPort))
  if err != nil {
    glog.Exit("Listen:", err)
  }

  var servicesHandled = 0
  for {
    glog.V(2).Info("Accepting")
    conn, err := ln.Accept()
    if err != nil {
      glog.Exit("Accept:", err)
    }
    defer conn.Close()
    servicesHandled++
    glog.V(1).Infof("relaying %d ringy dingy\n", servicesHandled)
    go relay.HandleRelayedService(conn, servicePort, relayPort)
    servicePort++
    relayPort++
  }
}
