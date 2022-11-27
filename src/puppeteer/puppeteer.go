package puppeteer

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
	"net"
	"time"

	"github.com/kenji-yamane/pq-deadlock/src/network"
)

type Puppeteer struct {
	ports       []string
	connections map[int]*net.UDPConn
	commands    []Command
}

func NewPuppeteer(
	ports []string,
) *Puppeteer {
	puppeteer := &Puppeteer{
		ports:    ports,
		commands: make([]Command, 0),
	}
	puppeteer.initConnections()
	return puppeteer
}

func (p *Puppeteer) AddCommand(id int, line string, await time.Duration) *Puppeteer {
	p.commands = append(p.commands, Command{
		Id:    id,
		Line:  line,
		Await: await,
	})
	return p
}

func (p *Puppeteer) Execute() {
	for _, cmd := range p.commands {
		msg := cmd.ToMessage()
		msgStr, err := json.MarshalIndent(msg, "", "\t")
		if err != nil {
			customerror.CheckError(fmt.Errorf("error serializing"))
		}
		network.UdpSend(p.connections[cmd.Id], string(msgStr))
		time.Sleep(cmd.Await)
	}
}

func (p *Puppeteer) initConnections() {
	connections := make(map[int]*net.UDPConn)
	for idx, port := range p.ports {
		conn := network.UdpConnect(port)
		connections[idx+1] = conn
	}
	p.connections = connections
	return
}
