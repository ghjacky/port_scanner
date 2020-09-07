package main

import (
	"flag"
)

func main() {
	ips := flag.String("ips", "", "指定ip列表，以','分隔")
	cidr := flag.String("cidr", "", "指定ip段")
	ports := flag.String("ports", "", "指定port列表，以','分隔")
	portRange := flag.String("port-range", "", "指定port范围")
	writeTo := flag.String("write-to", "./res.txt", "指定结果写入文件")
	flag.Parse()
	var allPorts = new(PORTS)
	var allIPs = new(IPS)
	*allIPs = append(*allIPs, *NewIPS(*ips)...)
	*allIPs = append(*allIPs, *NewIPS(*cidr)...)
	*allPorts = append(*allPorts, *NewPorts(*ports)...)
	*allPorts = append(*allPorts, *NewPorts(*portRange)...)
	Scan(*allIPs, *allPorts, *writeTo)
}
