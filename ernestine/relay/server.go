package relay

import (
  "encoding/gob"
  "fmt"
  "github.com/golang/glog"
  "net"
)

func handleErr(err error, msg string) {
  if err != nil {
    glog.Exit(msg, err)
  }
}

func HandleRelayedService(conn net.Conn, servicePort int, relayPort int) {
  // First open the two ports for listening, THEN tell the service
  // or else a potential race condition exists
  glog.V(2).Info("Opening port for relayed service:", servicePort)
  serviceLn, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
  handleErr(err, "Listening:")

  glog.V(2).Info("Opening relayed port ", relayPort)
  relayLn, err := net.Listen("tcp", fmt.Sprintf(":%d", relayPort))
  handleErr(err, "Listening:")

  glog.V(1).Infof("Ready to relay from port %d to port %d\n", servicePort, relayPort)

  // Ready to relay, let service know
  relayedService := gob.NewEncoder(conn)
  newServiceMessage := Message{NewService, fmt.Sprintf("%d", servicePort)}
  err = relayedService.Encode(newServiceMessage)
  handleErr(err, "Encode error: ")

  for {
    clientConn, err := serviceLn.Accept()
    handleErr(err, "In Accept: ")
    glog.V(2).Info("Accepted client connection on port ", servicePort)
    defer clientConn.Close()

    newConnectionMessage := Message{NewConnectionOnService, fmt.Sprintf("%d", relayPort)}
    err = relayedService.Encode(newConnectionMessage)
    handleErr(err, "Encode error: ")

    relayConn, err := relayLn.Accept()
    handleErr(err, "In Accept: ")
    glog.V(2).Info("Accepted service connection on port ", relayPort)

    // TODO: When do I properly close the connection?
    // Or do I just let everything close itself when done? Hmm...
    go CopyStreams(clientConn, relayConn)
    go CopyStreams(relayConn, clientConn)
  }
}
