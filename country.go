package database

import "bytes"

const (
	NameField       = "country"
	PopulationField = "population"
	NorthField      = "north"
	SouthField      = "south"
	EastField       = "east"
	WestField       = "west"
)

const KeySize = 56

// Country is a single database row.
type Country struct {
	Length                   uint64
	Name                     [KeySize]byte
	Population               uint64
	North, South, East, West float64
}

type ByName []Country

func (by ByName) Len() int {
	return len(by)
}

func (by ByName) Less(i, j int) bool {
	return bytes.Compare(by[i].Name[:], by[j].Name[:]) < 0
}

func (by ByName) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}
