package main

import "encoding/json"

// Message represents a network message exchanged between peers. It includes the
// operation type (Ops), an optional JSON-encoded payload, and the sender info.
type Message struct {
    Ops string `json:"ops"`

    Payload string `json:"payload"`
    Sender  Peer   `json:"sender"`
}

// Json serializes the message into JSON. Errors are ignored in this demo code.
func (m Message) Json() []byte {
    json, _ := json.Marshal(m)
    return json
}

// WithSender returns a shallow copy of the message with the provided sender set.
func (m Message) WithSender(sender Peer) Message {
    m.Sender = sender
    return m
}

// NewRegisterMessage creates a registration message that carries the sender's
// listening address and port so that peers can learn about a new node.
func NewRegisterMessage(listenAddr string, listenPort int) Message {

    payload, _ := json.Marshal(Peer{listenAddr, listenPort})
    return Message{
        Ops:     "register",
        Payload: string(payload),
    }
}
// NewReanounceMessage announces this node again to peers. Used after discovery
// to propagate our presence.
func NewReanounceMessage(listenAddr string, listenPort int) Message {

    payload, _ := json.Marshal(Peer{listenAddr, listenPort})
    return Message{
        Ops:     "reanounce",
        Payload: string(payload),
    }
}
// NewPropogateMessage is used to propagate information about a new peer to
// others.
func NewPropogateMessage(listenAddr string, listenPort int) Message {

    payload, _ := json.Marshal(Peer{listenAddr, listenPort})
    return Message{
        Ops:     "propogate",
        Payload: string(payload),
    }
}

// NewPromoteMessage starts a leader election (RequestVote equivalent in Raft).
func NewPromoteMessage() Message {
    return Message{
        Ops: "promote",
    }
}

// NewHeartbeatMessage represents a leader heartbeat to followers.
func NewHeartbeatMessage() Message {
    return Message{
        Ops: "beat",
    }
}

// NewVoteMessage represents a positive vote response to a candidate.
func NewVoteMessage() Message {
    return Message{
        Ops: "vote",
    }
}
