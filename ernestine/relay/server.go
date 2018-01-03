package relay

import (
  "encoding/gob"
  "fmt"
  "github.com/golang/glog"
  "net"
)

func HandleRelayedService(conn net.Conn, servicePort int, relayPort int) {
  // First open the two ports for listening, THEN tell the service
  // or else a potential race condition exists
  glog.V(2).Info("Opening port for relayed service:", servicePort)
  serviceLn, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
  if err != nil {
    glog.Exit("Listening:", err)
    return
  }
  glog.V(2).Info("Opening relayed port ", relayPort)
  relayLn, err := net.Listen("tcp", fmt.Sprintf(":%d", relayPort))
  if err != nil {
    glog.Exit("Listening:", err)
    return
  }
  glog.V(1).Infof("Ready to relay from port %d to port %d\n", servicePort, relayPort)

  // Ready to relay, let service know
  relayedService := gob.NewEncoder(conn)
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
    glog.V(2).Info("Accepted service connection on port ", relayPort)

    // TODO: When do I properly close the connection?
    // Or do I just let everything close itself when done? Hmm...
    go CopyStreams(clientConn, relayConn)
    go CopyStreams(relayConn, clientConn)
  }
}
