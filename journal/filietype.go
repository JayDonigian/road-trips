package journal

type fileType int

const (
	entry = iota
	dayMap
	bikeMap
	totalMap
)

func (ft fileType) format() string {
	switch ft {
	case entry:
		return "%s/entries/%s.md"
	case dayMap:
		return "%s/maps/day/%s.png"
	case bikeMap:
		return "%s/maps/bike/%s.png"
	case totalMap:
		return "%s/maps/total/%s-total.png"
	default:
		return ""
	}
}

func (ft fileType) formatPathRelativeToEntry() string {
	switch ft {
	case entry:
		return "%s.md"
	case dayMap:
		return "../maps/day/%s.png"
	case bikeMap:
		return "../maps/bike/%s.png"
	case totalMap:
		return "../maps/total/%s-total.png"
	default:
		return ""
	}
}
