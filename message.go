package raft

import "encoding/json"

type message struct {
	Ops string `json:"ops"`

	Payload string `json:"payload"`
	Sender  peer   `json:"sender"`
}

func (m message) Json() []byte {
	json, _ := json.Marshal(m)
	return json
}

func (m message) WithSender(sender peer) message {
	m.Sender = sender
	return m
}

func NewRegisterMessage(listenAddr string, listenPort int) message {

	payload, _ := json.Marshal(peer{listenAddr, listenPort})
	return message{
		Ops:     "register",
		Payload: string(payload),
	}
}
func NewReanounceMessage(listenAddr string, listenPort int) message {

	payload, _ := json.Marshal(peer{listenAddr, listenPort})
	return message{
		Ops:     "reanounce",
		Payload: string(payload),
	}
}
func NewPropogateMessage(listenAddr string, listenPort int) message {

	payload, _ := json.Marshal(peer{listenAddr, listenPort})
	return message{
		Ops:     "propogate",
		Payload: string(payload),
	}
}

func NewPromoteMessage() message {
	return message{
		Ops: "promote",
	}
}

func NewHeartbeatMessage() message {
	return message{
		Ops: "beat",
	}
}

func NewVoteMessage() message {
	return message{
		Ops: "vote",
	}
}
