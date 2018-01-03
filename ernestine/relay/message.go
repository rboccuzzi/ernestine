package relay

type Action int

const (
  NewService Action = iota
  NewConnectionOnService
  EndConnectionOnService // Not sure if I need this
)

type Message struct {
  Action Action
  Data   string
}
