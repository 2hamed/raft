package main

import "encoding/json"

type Message struct {
	Ops string `json:"ops"`

	Payload string `json:"payload"`
}

func (m Message) Json() []byte {
	json, _ := json.Marshal(m)
	return json
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

