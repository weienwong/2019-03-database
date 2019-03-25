package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/weienwong/2019-03-database"
)

type Geo struct {
	Country string `json:"country"`
	North   string `json:"north"`
	South   string `json:"south"`
	West    string `json:"west"`
	East    string `json:"east"`
}

type Pop struct {
	Country    string `json:"country"`
	Population string `json:"population"`
}

func GetGeos() ([]Geo, error) {
	var geos []Geo
	f, err := os.Open("../geo.json")
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&geos)
	return geos, err
}

func GetPops() ([]Pop, error) {
	var pops []Pop
	f, err := os.Open("../pop.json")
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&pops)
	return pops, err
}

func parseFloat(s string) float64 {
	if strings.TrimSpace(s) == "" {
		return 0
	}

	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func main() {
	geos, err := GetGeos()
	if err != nil {
		panic(err)
	}

	pops, err := GetPops()
	if err != nil {
		panic(err)
	}

	geoPop := make(map[string]database.Country, len(geos))
	var nameLength int

	for _, g := range geos {
		var c database.Country

		if len(g.Country) > nameLength {
			nameLength = len(g.Country)
		}

		c.Length = uint64(len(g.Country))
		copy(c.Name[:], g.Country)
		c.North = parseFloat(g.North)
		c.South = parseFloat(g.South)
		c.West = parseFloat(g.West)
		c.East = parseFloat(g.East)

		geoPop[g.Country] = c
	}

	for _, p := range pops {

		c, ok := geoPop[p.Country]
		if !ok {
			c.Length = uint64(len(p.Country))
			copy(c.Name[:], p.Country)
		}

		c.Population, _ = strconv.ParseUint(p.Population, 10, 64)
		geoPop[p.Country] = c
	}

	fmt.Println(geoPop)
	fmt.Println(nameLength)

	fmt.Println(unsafe.Sizeof(database.Country{}))

	f, err := os.Create("mapping.sdg")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sorted := make([]database.Country, 0, len(geoPop))

	for k, b := range geoPop {
		sorted = append(sorted, b)
		_ = k
	}

	sort.Sort(database.ByName(sorted))
	binary.Write(f, binary.LittleEndian, sorted)
}
