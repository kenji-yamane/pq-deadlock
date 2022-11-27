package process

import (
	"fmt"
	"github.com/kenji-yamane/pq-deadlock/src/clock"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
	"github.com/kenji-yamane/pq-deadlock/src/math"
	"github.com/kenji-yamane/pq-deadlock/src/network"
	"github.com/kenji-yamane/pq-deadlock/src/process/messages"
	"net"
)

type snapshotRecord struct {
	out     []int
	in      []int
	time    int
	blocked bool
	replies int
}

type localSnapshot struct {
	records map[int]*snapshotRecord
}

type Process struct {
	id          int
	ports       []string
	wait        bool
	lastBlocked int
	in          []int
	out         []int
	replies     int
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
	process := &Process{
		id:          id,
		ports:       ports,
		wait:        false,
		lastBlocked: -1,
		in:          make([]int, 0),
		out:         make([]int, 0),
		replies:     0,
		weight:      1.0,
		snapshot: localSnapshot{
			records: make(map[int]*snapshotRecord, 0),
		},
		clock: clock,
	}
	for i := 0; i <= len(ports); i++ {
		process.snapshot.records[i] = &snapshotRecord{
			in:      make([]int, 0),
			out:     make([]int, 0),
			time:    0,
			blocked: false,
			replies: 0,
		}
	}
	process.initConnections()
	return process
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
	return
}

func (p *Process) Serve(ch chan string) {
	network.Serve(ch, p.ports[p.id-1])
}

func (p *Process) InterpretCommand(cmdString string) {
	cmd := messages.IdentifyCommand(cmdString)
	switch cmd {
	case messages.Ask:
		reqCmd := messages.ParseRequestCommand(cmdString, len(p.ports))
		p.requestPOutOfQ(reqCmd)
	case messages.Liberate:
		libCmd := messages.ParseLiberateCommand(cmdString, len(p.ports))
		p.replyToParents(libCmd)
	case messages.Detect:
		p.snapshotInitiate()
	default:
		fmt.Println("invalid command, ignoring...")
	}
}

func (p *Process) requestPOutOfQ(req *messages.RequestCommand) {
	p.replies = req.NeededReplies
	p.wait = true
	p.lastBlocked = p.clock.GetTicks()
	for _, outId := range req.ChildIds {
		p.out = append(p.out, outId)
		p.sendMessage(outId, messages.Request, 0, p.id, p.clock.GetTicks())
	}
}

func (p *Process) replyToParents(cmd *messages.LiberateCommand) {
	if cmd.LiberateAll {
		for _, parentId := range p.in {
			p.sendMessage(parentId, messages.Reply, 0, p.id, p.clock.GetTicks())
		}
		p.in = make([]int, 0)
		return
	}
	for _, parentId := range cmd.ParentIds {
		if math.Contains(p.in, parentId) {
			math.RemoveFrom(p.in, parentId)
			p.sendMessage(parentId, messages.Reply, 0, p.id, p.clock.GetTicks())
		}
	}
}

func (p *Process) InterpretMessage(msg string) {
	parsedMsg, err := messages.ParseMessage(msg)
	if err != nil {
		fmt.Println("invalid message, ignoring...")
		return
	}

	switch messages.MessageType(parsedMsg.Type) {
	case messages.Request:
		p.processRequest(parsedMsg)
		p.clock.ExternalEvent(parsedMsg.ClockStr)
	case messages.Reply:
		p.processReply(parsedMsg)
		p.clock.ExternalEvent(parsedMsg.ClockStr)
	case messages.Cancel:
		p.clock.ExternalEvent(parsedMsg.ClockStr)
	case messages.Flood:
		p.processFlood(parsedMsg.SenderId, parsedMsg.InitiatorId, parsedMsg.InitiatedAt, parsedMsg.Weight)
	case messages.Echo:
		p.processEcho(parsedMsg.SenderId, parsedMsg.InitiatorId, parsedMsg.InitiatedAt, parsedMsg.Weight)
	case messages.Short:
		p.processShort(parsedMsg.SenderId, parsedMsg.InitiatedAt, parsedMsg.Weight)
	default:
		fmt.Println("invalid message, ignoring...")
	}
}

func (p *Process) processRequest(msg messages.Message) {
	if !math.Contains(p.in, msg.SenderId) {
		p.in = append(p.in, msg.SenderId)
	}
}

func (p *Process) processReply(msg messages.Message) {
	if !math.Contains(p.out, msg.SenderId) {
		return
	}
	math.RemoveFrom(p.out, msg.SenderId)
	p.replies--
	if p.replies > 0 {
		return
	}
	p.wait = false
	for _, outId := range p.out {
		p.sendMessage(outId, messages.Cancel, 0, p.id, p.clock.GetTicks())
	}
	p.out = make([]int, 0)
}

func (p *Process) closeConnections() {
	for _, conn := range p.connections {
		err := conn.Close()
		customerror.CheckError(err)
	}
}

func (p *Process) snapshotInitiate() {
	p.weight = 0
	p.snapshot.records[p.id].time = p.clock.GetTicks()
	p.snapshot.records[p.id].out = p.out
	p.snapshot.records[p.id].blocked = true
	p.snapshot.records[p.id].in = []int{}
	p.snapshot.records[p.id].replies = p.replies

	for _, outId := range p.out {
		p.sendMessage(outId, messages.Flood, 1.0/float64(len(p.out)), p.id, p.clock.GetTicks())
	}
}

func (p *Process) sendMessage(outId int, messageType messages.MessageType, weight float64, initId int, initiatedAt int) {
	network.UdpSend(p.connections[outId], messages.BuildMessage(p.id, p.clock, messageType, weight, initId, initiatedAt))
}

func (p *Process) processFlood(j int, init int, initiatedAt int, weight float64) {
	snapshot := p.snapshot.records[init]

	if snapshot.time < initiatedAt && math.Contains(p.in, j) {
		snapshot.out = p.out
		snapshot.in = make([]int, 0)
		snapshot.in = append(snapshot.in, j)
		snapshot.blocked = p.wait
		if p.wait {
			snapshot.replies = p.replies
			for _, outId := range p.out {
				p.sendMessage(outId, messages.Flood, weight/float64(len(p.out)), init, initiatedAt)
			}

		}
		if p.wait == false {
			snapshot.replies = 0
			p.sendMessage(j, messages.Echo, weight, init, initiatedAt)
			snapshot.in = math.RemoveFrom(snapshot.in, j)
		}
	}

	if snapshot.time < initiatedAt && !math.Contains(p.in, j) {
		p.sendMessage(j, messages.Echo, weight, init, initiatedAt)
	}

	if snapshot.time == initiatedAt && !math.Contains(p.in, j) {
		p.sendMessage(j, messages.Echo, weight, init, initiatedAt)
	}

	if snapshot.time == initiatedAt && math.Contains(p.in, j) {
		if snapshot.blocked == false {
			p.sendMessage(j, messages.Echo, weight, init, initiatedAt)
		}

		if snapshot.blocked == true {

			if !math.Contains(snapshot.in, j) {
				snapshot.in = append(snapshot.in, j)
			}

			p.sendMessage(init, messages.Short, weight, init, initiatedAt)

		}

	}

}

func (p *Process) processEcho(j int, init int, initiatedAt int, weight float64) {
	snapshot := p.snapshot.records[init]

	if snapshot.time > initiatedAt {
		return
	}

	if snapshot.time < initiatedAt {
		fmt.Println("Cannot happened, just happened")
		return
	}

	if snapshot.time == initiatedAt {
		snapshot.out = math.RemoveFrom(snapshot.out, j)
		if snapshot.blocked == false {
			p.sendMessage(init, messages.Short, weight, init, initiatedAt)
		}
		if snapshot.blocked {
			snapshot.replies--
			if snapshot.replies == 0 {
				snapshot.blocked = false

				if init == p.id {
					fmt.Println("No deadlock, we rock! :)") // TO DO fazer o que aqui?
					return
				}

				for _, inId := range snapshot.in {
					p.sendMessage(inId, messages.Echo, weight/float64(len(snapshot.in)), init, initiatedAt)
				}

			}

			if snapshot.replies != 0 {
				p.sendMessage(init, messages.Short, weight, init, initiatedAt)
			}
		}
	}
}

func (p *Process) processShort(init int, initiatedAt int, weight float64) {
	snapshot := p.snapshot.records[init]

	if initiatedAt < p.lastBlocked {
		return
	}

	if initiatedAt > p.lastBlocked {
		fmt.Println("Not possible, possible! Pls, give up!")
		return
	}

	if initiatedAt == p.lastBlocked && snapshot.blocked == false {
		return
	}

	if initiatedAt == p.lastBlocked && snapshot.blocked {
		p.weight += weight
	}

	if p.weight > 0.999 {
		fmt.Println("DeadLock!")
		return
	}
}
