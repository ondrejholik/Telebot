package telebot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

// Gps struct with Latitude, Longtitude
type Gps struct {
	Lat float32
	Lon float32
}

// Village name, gps point lat, lon
type Village struct {
	Name  string
	Point Gps
}

// VillageArr is there only, because json.Marshal
type VillageArr struct {
	Arr []Village
}

// AreVillagesInDB -- check if villages has been initialized(if exists)
func AreVillagesInDB(db *badger.DB) bool {
	var found = true
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("villages"))
		if err != nil {
			found = false
		}
		return nil
	})
	if err != nil {
		return false
	}
	return found
}

// ClosestVillage find closest village to point
func ClosestVillage(db *badger.DB, lat float32, lon float32) string {

	log.Println("Getting villages")
	var villagesJsonCopy []byte
	var villages VillageArr
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("villages"))
		if err != nil {
			log.Println("Villages not found")
			log.Panic(err)
		}
		err = item.Value(func(val []byte) error {
			villagesJsonCopy = append([]byte{}, val...)
			return err
		})

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(villagesJsonCopy, &villages)
	if err != nil {
		log.Println(err)
	}

	var closestVillage Village
	var current float32
	minDistance := float32(math.Inf(1))
	for _, x := range villages.Arr {
		current = HarvestineDistance(x.Point, Gps{Lat: lat, Lon: lon})
		if current < minDistance {
			closestVillage = x
			minDistance = current
		}
	}
	return fmt.Sprintf("Your new location has been set to:\n%s which is %g km from you", closestVillage.Name, minDistance)
}

func degreesToRadians(d float32) float64 {
	return float64(d) * math.Pi / 180
}

func HarvestineDistance(x, y Gps) float32 {
	const earthRaidusKm = 6371
	lat1 := degreesToRadians(x.Lat)
	lon1 := degreesToRadians(x.Lon)
	lat2 := degreesToRadians(y.Lat)
	lon2 := degreesToRadians(y.Lon)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km := c * earthRaidusKm

	return float32(km)
}

// ParseToVillage take slice of string splitted by delimiter and convert to Village struct
func ParseToVillage(row []string) Village {
	var lat, lon float32
	val, err := strconv.ParseFloat(row[1], 32)
	if err != nil {
		log.Println(err)
	}
	lat = float32(val)
	val, err = strconv.ParseFloat(row[2], 32)
	if err != nil {
		log.Println(err)
	}
	lon = float32(val)
	return Village{Name: row[0], Point: Gps{Lat: lat, Lon: lon}}
}

// LoadVillagesToDB load to badgerDB
func LoadVillagesToDB(db *badger.DB) {
	var villages []Village

	dat, err := ioutil.ReadFile("assets/villages.csv")
	if err != nil {
		log.Panic(err)
	}
	r := csv.NewReader(strings.NewReader(string(dat)))
	r.Comma = ';'

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		village := ParseToVillage(record)
		villages = append(villages, village)
	}

	err = db.Update(func(txn *badger.Txn) error {
		vlgsStruct := &VillageArr{
			Arr: villages,
		}
		encoded, err := json.Marshal(vlgsStruct)
		if err != nil {
			log.Println(err)
		}
		err = txn.Set([]byte("villages"), []byte(encoded))
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	log.Println("Villages in DB")
}
