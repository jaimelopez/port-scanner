package scanner

import "errors"

type Scanner struct {
	Address string
	Ports   []int
}

func (scan *Scanner) AddPort(port int) {
	scan.Ports = append(scan.Ports, port)
}

func (scan *Scanner) AddPortCollection(ports []int) {
}

func (scan *Scanner) AddRange(starts int, ends int) error {
	if starts > ends {
		return errors.New("Invalid range")
	}

	for i := starts; i <= ends; i++ {
		scan.AddPort(i)
	}
}

func (scan *Scanner) Run() (scanReport ScanReport, error error) {
	if error != nil {
		return
	}

	if len(scan.Ports) > 0 {
		scanReport, error = SpecificPorts(scan.Address, scan.Ports)
	} else {
		scanReport, error = AllOpenedPorts(scan.Address)
	}

	if error != nil {
		return
	}

	return
}
