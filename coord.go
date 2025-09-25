package main

import (
    "encoding/json"
    "fmt"
    "time"
)

// Coordinator encapsulates the core node state and behavior for a Raft-like
// coordination process. It manages peer membership, election timeouts,
// leader election, and heartbeats.
type Coordinator struct {
    // Self represents this node's network identity.
    Self  Peer
    // Peers contains the known peers in the cluster (excluding Self).
    Peers Peers
    // Timer is the election timeout timer. When it fires, a new election starts.
    Timer *RaftTimer

    // heartBeat is an internal signal channel used to indicate receipt of a
    // heartbeat (from leader) or to trigger leader heartbeats to followers.
    heartBeat   chan struct{}
    // isElected indicates whether this node believes it is the current leader.
    isElected   bool
    // isCandidate indicates whether this node is currently campaigning.
    isCandidate bool

    // votes tracks the number of votes received during an ongoing election.
    votes int
}

// NewCoordinator constructs and initializes a Coordinator for the given host
// and port. It starts background goroutines to handle heartbeats and election
// timeouts and returns the ready-to-use coordinator instance.
func NewCoordinator(host string, port int) *Coordinator {
    c := &Coordinator{
        Peers:     make(Peers, 0),
        Self:      Peer{host, port},
        Timer:     NewTimer(),
        heartBeat: make(chan struct{}, 10),
    }

    go func() {
        for {
            select {
            case <-c.heartBeat:
                c.Timer.Reset()
                if c.isElected {
                    c.sendHeartbeat()
                }
            case <-c.Timer.TimeoutSignal():
                c.isCandidate = true
                c.PromoteSelf()
            }
        }
    }()

    go func() {
        ticker := time.Tick(50 * time.Millisecond)
        for {
            <-ticker
            if c.isElected {
                c.heartBeat <- struct{}{}
            }
        }
    }()

    c.Timer.StartTimer()

    return c
}

// ProcessMessage is the entry point for handling all inbound messages from the
// network. It routes messages based on their operation type and updates
// Coordinator state accordingly.
func (c *Coordinator) ProcessMessage(msg Message) error {

    switch msg.Ops {
    case "register", "reanounce", "propogate":
        c.peerInfoReceived(msg)
    case "beat":
        c.heartBeat <- struct{}{}
    case "promote":
        if !c.isCandidate && !c.isElected {
            c.sendVote(msg.Sender)
        }
    case "vote":
        c.votes++
        if c.votes >= c.Peers.Quorum() {
            fmt.Println("I am the elected leader.")
            c.isElected = true
        }
    }

    return nil
}

// broadcastMessage sends the given message to all known peers.
func (c *Coordinator) broadcastMessage(msg Message) error {
    return c.Peers.BroadcastMessage(msg)
}

// peerInfoReceived handles membership-related messages (register, reanounce,
// and propogate). It updates the local peers list and re-announces self when
// appropriate to help new nodes discover the full cluster.
func (c *Coordinator) peerInfoReceived(msg Message) error {
    var p Peer
    if err := json.Unmarshal([]byte(msg.Payload), &p); err != nil {
        return err
    }

    if c.Self.Equals(p) {
        return nil
    }

    switch msg.Ops {
    case "register":
        c.broadcastMessage(NewPropogateMessage(p.Host, p.Port).WithSender(c.Self))
        c.registerPeer(p)
        c.reanounceSelf()
    case "reanounce":
        c.registerPeer(p)
    case "propogate":
        c.registerPeer(p)
        c.reanounceSelf()
    }
    c.Peers.PrintInfo()
    return nil
}
// registerPeer adds a peer to the local peer set if it is not already present.
func (c *Coordinator) registerPeer(peer Peer) error {
    if !c.Peers.Contains(peer) {
        c.Peers = append(c.Peers, peer)
    }

    return nil
}

// joinCluster attempts to join an existing cluster by registering with the
// specified peer, advertising our listen address and port.
func (c *Coordinator) joinCluster(port int, peer Peer) error {
    fmt.Printf("Joining peer on %s:%d\n", peer.Host, peer.Port)
    return peer.SendMessage(NewRegisterMessage("127.0.0.1", port).WithSender(c.Self))
}

// reanounceSelf broadcasts this node's identity so that peers can update their
// membership views.
func (c *Coordinator) reanounceSelf() error {
    return c.broadcastMessage(NewReanounceMessage(c.Self.Host, c.Self.Port).WithSender(c.Self))
}

// PromoteSelf starts an election by self-voting and requesting votes from peers.
func (c *Coordinator) PromoteSelf() error {
    fmt.Println("Starting election...")
    c.votes = 1
    return c.broadcastMessage(NewPromoteMessage().WithSender(c.Self))
}

// sendHeartbeat broadcasts a heartbeat message to all peers. Only the leader
// should send heartbeats.
func (c *Coordinator) sendHeartbeat() error {
    return c.broadcastMessage(NewHeartbeatMessage().WithSender(c.Self))
}

// sendVote sends a vote message to the specified candidate.
func (c *Coordinator) sendVote(target Peer) error {
    return target.SendMessage(NewVoteMessage().WithSender(c.Self))
}
