package process

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/kenji-yamane/pq-deadlock/src"
	"github.com/kenji-yamane/pq-deadlock/src/clock"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
	"github.com/kenji-yamane/pq-deadlock/src/network"
	"github.com/kenji-yamane/pq-deadlock/src/process/messages"
)

type snapshotRecord struct {
	out     []int
	in      []int
	time    int
	blocked bool
	p       int
}

type localSnapshot struct {
	records map[int]snapshotRecord
}

type Process struct {
	id          int
	ports       []string
	wait        bool
	lastBlocked int
	in          []int
	out         []int
	p           int
	weight      float64
	snapshot    localSnapshot
	clock       clock.LogicalClock
	connections map[int]*net.UDPConn
}

func NewProcess(
	id int,
	ports []string,
	clock clock.LogicalClock,
) *Process {
	p := &Process{
		id:          id,
		ports:       ports,
		wait:        false,
		lastBlocked: -1,
		in:          make([]int, 0),
		out:         make([]int, 0),
		p:           0,
		weight:      1.0,
		snapshot: localSnapshot{
			records: make(map[int]snapshotRecord, 0),
		},
		clock: clock,
	}
	for i := 0; i <= len(ports); i++ {
		p.snapshot.records[i] = snapshotRecord{
			in:      make([]int, 0),
			out:     make([]int, 0),
			time:    0,
			blocked: false,
			p:       0,
		}
	}
	p.initConnections()
	return p
}

func (p *Process) initConnections() {
	connections := make(map[int]*net.UDPConn)
	for idx, port := range p.ports {
		if idx+1 == p.id {
			continue
		}
		conn := network.UdpConnect(port)
		connections[idx+1] = conn
	}
	p.connections = connections
	p.csConn = network.UdpConnect(src.SharedResourcePort)
	return
}

func (p *Process) Serve(ch chan string) {
	network.Serve(ch, p.ports[p.id-1])
}

func (p *Process) InterpretCommand(cmd string) {
	switch cmd {
	case strconv.Itoa(p.id):
		p.clock.InternalEvent()
	case ConsumeCmd:
		p.requestSharedResource()
	default:
		fmt.Println("invalid command, ignoring...")
	}
}

func (p *Process) requestSharedResource() {
	switch p.state {
	case Released:
		p.state = Wanted
		p.clock.InternalEvent()
		p.lastRequest = p.clock.GetClockStr()
		for id := 0; id < len(p.ports); id++ {
			if id+1 == p.id {
				continue
			}
			network.UdpSend(p.connections[id+1], messages.BuildRequestMessage(p.id, p.clock))
		}
	case Wanted:
		fmt.Println("x ignored")
	case Held:
		fmt.Println("x ignored")
	}
}

func (p *Process) InterpretMessage(msg string) {
	parsedMsg, err := messages.ParseMessage(msg)
	if err != nil {
		fmt.Println("invalid message, ignoring...")
		return
	}
	p.clock.ExternalEvent(parsedMsg.ClockStr)

	switch messages.MessageType(parsedMsg.Text) {
	case messages.Request:
		p.processRequest(parsedMsg)
	case messages.Reply:
		p.processReply()
	case messages.Consume:
		fmt.Printf("received %s, but I'm not a shared resource, ignoring...\n", messages.Consume)
	default:
	}
}

func (p *Process) processRequest(msg messages.Message) {
	switch p.state {
	case Released:
		network.UdpSend(p.connections[msg.SenderId], messages.BuildReplyMessage(p.id, p.clock))
	case Wanted:
		selectedId, err := p.clock.CompareClocks(p.lastRequest, msg.ClockStr, msg.SenderId)
		if err != nil {
			fmt.Println("invalid message, ignoring...")
			break
		}
		if selectedId == p.id {
			p.replyManager.EnqueueProcess(msg.SenderId)
		} else {
			network.UdpSend(p.connections[msg.SenderId], messages.BuildReplyMessage(p.id, p.clock))
		}
	case Held:
		p.replyManager.EnqueueProcess(msg.SenderId)
	}
}

func (p *Process) processReply() {
	if !p.replyManager.ReceiveReply() {
		return
	}
	p.state = Held
	network.UdpSend(p.csConn, messages.BuildConsumeMessage(p.id, p.clock))
	fmt.Println("entered cs")
	time.Sleep(5 * time.Second)
	p.state = Released
	fmt.Println("left cs")
	processesToReply := p.replyManager.Dequeue()
	for _, id := range processesToReply {
		network.UdpSend(p.connections[id], messages.BuildReplyMessage(p.id, p.clock))
	}
}

func (p *Process) closeConnections() {
	for _, conn := range p.connections {
		err := conn.Close()
		customerror.CheckError(err)
	}
	err := p.csConn.Close()
	customerror.CheckError(err)
}
