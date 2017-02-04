package scanner

type EntryPoint struct {
	Host string
	Port int
	Opened bool
}

type ScanReport map[string][]int
