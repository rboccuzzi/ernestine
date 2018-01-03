package relay

import (
  "encoding/gob"
  "fmt"
  "github.com/golang/glog"
  "net"
)

func HandleRelayedService(conn net.Conn, servicePort int, relayPort int) {
  relayedService := gob.NewEncoder(conn)

  // First open the port, THEN tell service (or else a potential race condition)
  glog.V(2).Info("Accepted a relay request. Opening port for relayed service:", servicePort)
  serviceLn, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
  if err != nil {
    glog.Exit("Listening:", err)
    return
  }

  // Again, first open the port THEN tell service (to avoid potential race condition)
  // TODO: I think I need a different port for each client connection
  // or else I could cross-connect under rapid connection load
  relayLn, err := net.Listen("tcp", fmt.Sprintf(":%d", relayPort))
  if err != nil {
    glog.Exit("Listening:", err)
    return
  }
  glog.V(2).Info("Ready to relay on port ", relayPort)

  newServiceMessage := Message{NewService, fmt.Sprintf("localhost:%d", servicePort)}
  err = relayedService.Encode(newServiceMessage)
  if err != nil {
    glog.Exit("Encode error: ", err)
  }

  for {
    clientConn, err := serviceLn.Accept()
    if err != nil {
      glog.Exit("In Accept: ", err)
    }
    glog.V(2).Info("Accepted client connection on port ", servicePort)
    defer clientConn.Close()

    newConnectionMessage := Message{NewConnectionOnService, fmt.Sprintf("localhost:%d", relayPort)}
    err = relayedService.Encode(newConnectionMessage)
    if err != nil {
      glog.Exit("Encode error: ", err)
    }

    relayConn, err := relayLn.Accept()
    if err != nil {
      glog.Exit("In Accept: ", err)
    }
    // TODO: When do I properly close the connection?
    go CopyStreams(clientConn, relayConn)
    go CopyStreams(relayConn, clientConn)
  }
}
