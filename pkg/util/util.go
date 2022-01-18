package util

import (
	"log"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

func IPRangeSize(ipRange string) uint64 {
	_, ipnet, err := net.ParseCIDR(ipRange)
	if err != nil {
		log.Printf("unable to parse ip cidr %v,  %v", ipRange, err)
		return 0
	}
	return cidr.AddressCount(ipnet)
}

func AllocatableIps(ipRange string, count int) int {

	cidrCapacity := IPRangeSize(ipRange)

	return int(cidrCapacity) - count
}
