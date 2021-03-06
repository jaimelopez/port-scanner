package scanner

import (
	"fmt"
	"net"
	"time"
)

const (
	REQUEST_TIMEOUT     = time.Duration(time.Millisecond * 400)
	REQUEST_CONCURRENCY = 200
	PORT_STARTS         = 1
	PORT_ENDS           = 65535
)

func AllOpenedPorts(address string) (scanReport ScanReport, error error) {
	ports := []int{}

	for i := PORT_STARTS; i <= PORT_ENDS; i++ {
		ports = append(ports, i)
	}

	scanReport, error = CheckSpecificPorts(address, ports)

	return
}

func CheckSpecificPorts(address string, ports []int) (scanReport ScanReport, error error) {
	entryPointCollection, error := ListAllHostsAndPorts(address, ports)

	if error != nil {
		return
	}

	concurrency := make(chan uint, REQUEST_CONCURRENCY)
	defer close(concurrency)

	channel := make(chan EntryPoint)
	defer close(channel)

	for _, entryPoint := range entryPointCollection {
		concurrency <- 1

		go func(entryPoint EntryPoint) {
			entryPoint.Opened = IsPortOpened(entryPoint.Host, entryPoint.Port)

			<-concurrency
			channel <- entryPoint
		}(entryPoint)
	}

	processedEntries := []EntryPoint{}

	for entryPoint := range channel {
		processedEntries = append(processedEntries, entryPoint)

		if len(entryPointCollection) == len(processedEntries) {
			break
		}
	}

	scanReport = ScanReport{}
	scanReport.Compose(processedEntries)

	return
}

func ListAllHostsAndPorts(address string, ports []int) (list []EntryPoint, error error) {
	hosts, error := ResolveAddress(address)

	if error != nil {
		return
	}

	for _, host := range hosts {
		for _, port := range ports {
			entryPoint := EntryPoint{Host: host, Port: port}

			list = append(list, entryPoint)
		}
	}

	return
}

func IsPortOpened(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	tcpAddr, error := net.ResolveTCPAddr("tcp4", address)

	if error != nil {
		return false
	}

	connection, error := net.DialTimeout("tcp", tcpAddr.String(), REQUEST_TIMEOUT)

	if error != nil {
		return false
	}

	connection.Close()

	return true
}
