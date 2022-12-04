package journal

type listCategory int

const (
	states = iota
	usParks
	provinces
	caParks
)

func (lc listCategory) String() string {
	switch lc {
	case states:
		return "* **U.S. States**"
	case usParks:
		return "* **National Parks**"
	case provinces:
		return "* **Canadian Provinces**"
	case caParks:
		return "* **Canadian National Parks**"
	default:
		return ""
	}
}
