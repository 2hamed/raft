package raft

import "encoding/json"

type Message struct {
	Ops string `json:"ops"`

	Payload string `json:"payload"`
	Sender  Peer   `json:"sender"`
}

func (m Message) Json() []byte {
	json, _ := json.Marshal(m)
	return json
}

func (m Message) WithSender(sender Peer) Message {
	m.Sender = sender
	return m
}

func NewRegisterMessage(listenAddr string, listenPort int) Message {

	payload, _ := json.Marshal(Peer{listenAddr, listenPort})
	return Message{
		Ops:     "register",
		Payload: string(payload),
	}
}
func NewReanounceMessage(listenAddr string, listenPort int) Message {

	payload, _ := json.Marshal(Peer{listenAddr, listenPort})
	return Message{
		Ops:     "reanounce",
		Payload: string(payload),
	}
}
func NewPropogateMessage(listenAddr string, listenPort int) Message {

	payload, _ := json.Marshal(Peer{listenAddr, listenPort})
	return Message{
		Ops:     "propogate",
		Payload: string(payload),
	}
}

func NewPromoteMessage() Message {
	return Message{
		Ops: "promote",
	}
}

func NewHeartbeatMessage() Message {
	return Message{
		Ops: "beat",
	}
}

func NewVoteMessage() Message {
	return Message{
		Ops: "vote",
	}
}
