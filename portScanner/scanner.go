package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"sync"
	"time"
)

func Scan(ips IPS, ports PORTS, ipThread, portThread int) {
	if len(ips) == 0 || len(ports) == 0 {
		return
	}
	if ipThread == 0 {
		ipThread = 100
	}
	if portThread == 0 {
		portThread = 100
	}
	wg := sync.WaitGroup{}
	for i, _ip := range ips {
		ip := _ip.ToString()
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			if ipReachable(ip) {
				wg1 := sync.WaitGroup{}
				for n, port := range ports {
					wg1.Add(1)
					go func(port uint16) {
						defer wg1.Done()
						if portOpen(ip, port) {
							fmt.Printf("%s:%d is open\n", ip, port)
						}
					}(uint16(port))
					if n+1%portThread == 0 {
						wg1.Wait()
					}
				}
				wg1.Wait()
			}
		}(ip)
		if i+1%ipThread == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
}

func ipReachable(ip string) (b bool) {
	var c = make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		b = <-c
		defer wg.Done()
	}()
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		log.Printf("ip：%s 有误", ip)
		return false
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		c <- true
	}
	p.OnIdle = func() {
		c <- false
	}
	err = p.Run()
	if err != nil {
		log.Printf("ip: %s 不可达: %s", ip, err)
		close(c)
		return false
	}
	close(c)
	wg.Wait()
	return
}

func portOpen(ip string, port uint16) (b bool) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3*time.Second)
	if err == nil && conn != nil {
		defer conn.Close()
		b = true
	}
	return
}
