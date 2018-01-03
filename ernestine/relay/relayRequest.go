package relay

import (
  "encoding/gob"
  "fmt"
  "github.com/golang/glog"
  "net"
)

// ServiceConnectionHandler is the callback that will be called from ClientConnect
// once a client to the service being relayed has initiated contact to the relay server.
// Treat this socket as you would a direct connection to the client
type ServiceConnectionHandler func(conn net.Conn)

// RelayServiceRequest is given the connection string to the relay server, and will manage
// all relay server communication. Each time a client connects to the relayed server
// for this service, a callback to the handler will be made.
func RelayServiceRequest(connection string, handler ServiceConnectionHandler) {
  glog.V(1).Info("Connecting to relay on ", connection)
  relayServerConn, err := net.Dial("tcp", connection)
  if err != nil {
    glog.Exit("Connecting to relay server ", err)
  }
  defer relayServerConn.Close()

  // Keep processing Messages on control channel
  dec := gob.NewDecoder(relayServerConn)
  var mess Message

  for {
    readErr := dec.Decode(&mess)
    if readErr != nil {
      glog.Exit("Reading ", readErr, ", mess is ", mess)
    }
    glog.V(2).Info("Read Message: ", mess)

    switch mess.Action {
    case NewService:
      fmt.Println("Established relay address: ", mess.Data)
    case NewConnectionOnService:
      conn, err := net.Dial("tcp", mess.Data)
      if err != nil {
        glog.Exit("Connecting to relay server ", err)
      }
      go handler(conn)

    case EndConnectionOnService:
      // Not sure what I need to do here. Probably use channel from above to blast all routines

    default:
      glog.Warning("Unknown Action: ", mess.Action)
    }
  }
}
