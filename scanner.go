package scanner

import "errors"

type Scanner struct {
	Address string
	ports   []int
}

func (scan *Scanner) AddPort(port int) {
	scan.ports = appendUnique(scan.ports, port)
}

func (scan *Scanner) AddPortCollection(ports []int) {
	for _, port := range ports {
		scan.AddPort(port)
	}
}

func (scan *Scanner) AddRange(starts int, ends int) error {
	if starts > ends {
		return errors.New("Invalid range")
	}

	for i := starts; i <= ends; i++ {
		scan.AddPort(i)
	}

	return nil
}

func appendUnique(slice []int, current int) []int {
	for _, element := range slice {
		if element == current {
			return slice
		}
	}

	return append(slice, current)
}

func (scan *Scanner) Run() (scanReport ScanReport, error error) {
	if error != nil {
		return
	}

	if len(scan.ports) > 0 {
		scanReport, error = CheckSpecificPorts(scan.Address, scan.ports)
	} else {
		scanReport, error = AllOpenedPorts(scan.Address)
	}

	if error != nil {
		return
	}

	return
}
