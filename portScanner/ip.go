package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type IP [4]byte
type IPS []IP

func NewIPS(is string) *IPS {
	if len(is) == 0 {
		return &IPS{}
	}
	var ip = new(IP)
	var ips = new(IPS)
	il := strings.Split(is, ",")
	_il := strings.Split(is, "/")
	if len(il) > 1 || (len(il) == 1 && len(_il) < 2) {
		for _, is := range il {
			ip = newIP(is)
			if ip != nil {
				*ips = append(*ips, *ip)
			}
		}
		return ips
	}
	return genIPS(_il)
}

func newIP(is string) *IP {
	var ip = new(IP)
	ipNodeList := strings.Split(is, ".")
	if len(ipNodeList) != 4 {
		log.Fatalln("输入的IP有误")
		return nil
	}
	for i, node := range ipNodeList {
		n, e := strconv.ParseUint(node, 10, 8)
		if e != nil || ((n == 0 || n == 255) && (i == 3 || i == 0)) {
			log.Fatalln("输入的ip有误")
			return nil
		}
		ip[i] = uint8(n)
	}
	return ip
}

func genIPS(il []string) *IPS {
	var ips = new(IPS)
	if len(il) != 2 {
		log.Fatalln("输入的ip段有误")
		return &IPS{}
	}
	cidrip := il[0]
	mask, err := strconv.ParseUint(il[1], 10, 8)
	if err != nil || mask >= 32 || mask < 16 {
		fmt.Println("输入的ip段有误，或不支持的网段")
		return &IPS{}
	}
	ip := newIP(cidrip)
	if ip == nil {
		return &IPS{}
	}
	quotient := mask / 8
	remainder := mask % 8
	lastNetNode := ip[quotient] >> (8 - remainder) << (8 - remainder)
	ip[quotient] = lastNetNode
	switch quotient {
	case 2:
		for i := lastNetNode; ; i++ {
			if i > 255 {
				break
			}
			ip[2] = i
			for n := 1; i < 255; i++ {
				ip[3] = uint8(n)
				*ips = append(*ips, *ip)
			}
		}
	case 3:
		for i := lastNetNode; ; i++ {
			if i == 0 {
				continue
			}
			if i == 255 {
				break
			}
			ip[3] = i
			*ips = append(*ips, *ip)
		}
	}
	return ips
}

func (ip IP) ToString() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}
