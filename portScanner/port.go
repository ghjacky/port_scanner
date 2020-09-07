package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type PORT uint16
type PORTS []PORT

func NewPorts(ps string) *PORTS {
	var ports = new(PORTS)
	if len(ps) == 0 {
		return &PORTS{}
	}
	portl := strings.Split(ps, ",")
	_portl := strings.Split(ps, "-")
	if len(portl) > 1 || (len(portl) == 1 && len(_portl) < 2) {
		for _, ps := range portl {
			p, e := strconv.Atoi(ps)
			if e != nil || p <= 0 || p > 65535 {
				log.Fatalf("端口输入有误: %d", p)
				return &PORTS{}
			}
			*ports = append(*ports, PORT(p))
		}
		return ports
	}
	if len(_portl) != 2 {
		fmt.Println("端口范围输入有误")
		return &PORTS{}
	}
	start, e1 := strconv.Atoi(_portl[0])
	stop, e2 := strconv.Atoi(_portl[1])
	if e1 != nil || e2 != nil || start <= 0 || stop <= 0 || start > 65535 || stop > 65535 || start > stop {
		fmt.Println("端口范围输入有误")
		return &PORTS{}
	}
	for i := start; i <= stop; i++ {
		*ports = append(*ports, PORT(i))
	}
	return ports
}
