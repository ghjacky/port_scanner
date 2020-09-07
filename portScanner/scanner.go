package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func Scan(ips IPS, ports PORTS, writeTo string) {
	if len(ips) == 0 || len(ports) == 0 {
		return
	}
	var fileToWrite *os.File = nil
	if len(writeTo) != 0 {
		f, e := os.OpenFile(writeTo, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0644)
		if e != nil {
			fmt.Printf("无法写入文件：%s\n", writeTo)
			return
		}
		fileToWrite = f
		defer f.Close()
	}
	wg := sync.WaitGroup{}
	wg1 := sync.WaitGroup{}
	lock := sync.Mutex{}
	for _, _ip := range ips {
		ip := _ip.ToString()
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			//fmt.Printf("thread: %d on %s starting...\n", (i+1)%Thread, ip)
			if ipReachable(ip) {
				for _, port := range ports {
					wg1.Add(1)
					go func(port PORT) {
						defer wg1.Done()
						if portOpen(ip, uint16(port)) {
							if fileToWrite != nil {
								lock.Lock()
								_, _ = fileToWrite.WriteString(fmt.Sprintf("%s:%d\n", ip, port))
								lock.Unlock()
							}
							fmt.Printf("%s:%d is open\n", ip, port)
						}
					}(port)
				}
				wg1.Wait()
			}
			//fmt.Printf("thread: %d on %s finish !\n", (i+1)%Thread, ip)
		}(ip)
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
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 5*time.Second)
	if err == nil && conn != nil {
		defer conn.Close()
		b = true
	}
	return
}
