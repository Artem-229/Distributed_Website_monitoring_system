package worker

import "math/rand"

type Region struct {
	Name      string
	MinOffset int
	MaxOffset int
}

func (r Region) RandomOffset() float64 {
	return float64(r.MinOffset + rand.Intn(r.MaxOffset-r.MinOffset+1))
}

var Regions = []Region{
	{Name: "Russia",         MinOffset: 5,   MaxOffset: 60},
	{Name: "USA",            MinOffset: 120, MaxOffset: 250},
	{Name: "China",          MinOffset: 200, MaxOffset: 380},
	{Name: "Central Europe", MinOffset: 20,  MaxOffset: 100},
	{Name: "Asia Pacific",   MinOffset: 150, MaxOffset: 300},
}
