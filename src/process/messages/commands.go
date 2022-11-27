package messages

import (
	"fmt"
	"strconv"
	"strings"
)

type RequestCommand struct {
	NeededReplies int
	ChildIds      []int
}

func IdentifyCommand(command string) CommandType {
	for _, cmdType := range []CommandType{Ask, Detect, Liberate} {
		if len(command) < len(string(cmdType)) {
			continue
		}
		if command[0:len(string(cmdType))] == string(cmdType) {
			return cmdType
		}
	}
	return Unknown
}

func ParseRequestCommand(command string, numPorts int) *RequestCommand {
	req := &RequestCommand{ChildIds: make([]int, 0)}
	nums := strings.Split(command[len(string(Ask))+1:], " ")
	fmt.Println(nums)
	for idx, numStr := range nums {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil
		}
		if idx == 0 {
			req.NeededReplies = num
			continue
		}
		if num <= 0 || num > numPorts {
			return nil
		}
		req.ChildIds = append(req.ChildIds, num)
	}
	if req.NeededReplies < 0 || req.NeededReplies > len(req.ChildIds) {
		return nil
	}
	return req
}

type LiberateCommand struct {
	ParentIds   []int
	LiberateAll bool
}

func ParseLiberateCommand(command string, numPorts int) *LiberateCommand {
	cmd := &LiberateCommand{LiberateAll: true, ParentIds: make([]int, 0)}
	nums := strings.Split(command[len(string(Liberate))+1:], " ")
	for _, numStr := range nums {
		cmd.LiberateAll = false
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil
		}
		if num <= 0 || num > numPorts {
			return nil
		}
		cmd.ParentIds = append(cmd.ParentIds, num)
	}
	return cmd
}
