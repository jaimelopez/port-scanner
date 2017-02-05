package scanner

type EntryPoint struct {
	Host string
	Port int
	Opened bool
}

type ScanReport map[string][]int

func (scanReport ScanReport) Compose(entries []EntryPoint) {
	for _, entryPoint := range entries {
		if !entryPoint.Opened {
			continue
		}

		scanReport[entryPoint.Host] = append(scanReport[entryPoint.Host], entryPoint.Port)
	}
}