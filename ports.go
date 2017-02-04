package scanner

import (
	"fmt"
	"net"
	"time"
	"github.com/jaimelopez/port-scanner/resolver"
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

	scanReport, error = SpecificPorts(address, ports)

	return
}

func SpecificPorts(address string, ports []int) (scanReport ScanReport, error error) {
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
			entryPoint.Opened = IsPortOpened(entryPoint.Host.string, entryPoint.Port.int)

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

	scanReport = makeReport(processedEntries)

	return
}

func makeReport(entries []EntryPoint) ScanReport {
	scanReport := make(ScanReport)

	for _, entryPoint := range entries {
		if entryPoint.Opened {
			scanReport[entryPoint.Host] = append(scanReport[entryPoint.Host], entryPoint.Port)
		}
	}

	return scanReport
}

func ListAllHostsAndPorts(address string, ports []int) (list []EntryPoint, error error) {
	hosts, error := resolver.Address(address)

	if error != nil {
		return
	}

	for _, host := range hosts {
		for _, port := range ports {
			entryPoint := EntryPoint{}
			entryPoint.Host = Host{host}
			entryPoint.Port = Port{port}

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
