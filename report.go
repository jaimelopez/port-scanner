package scanner

type Host struct{ string }
type Port struct{ int }

type EntryPoint struct {
	Host
	Port
	Opened bool
}

type ScanReport map[Host][]Port
