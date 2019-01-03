package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	client := connect()

	switch argument(1) {
	case "add":
		add(client, "cities", cities)
		add(client, "peaks", peaks)
	case "lookup":
		lookup(client, argument(2))
	case "find":
		find(client, argument(2))
	case "export":
		export(client, argument(2))
	case "flush":
		flushDB(client)
	default:
		usage()
	}
}

func connect() *redis.Client {
	// create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})

	// use the PING command to check the connection
	err := client.Ping().Err()
	if err != nil {
		panic(err)
	}

	return client
}

func add(client *redis.Client, key string, locations []*redis.GeoLocation) {
	// add locations to the database
	err := client.GeoAdd(key, locations...).Err()
	if err != nil {
		panic(err)
	}
}

func lookup(client *redis.Client, name string) {
	// look up the coordinates for a city
	positions, err := client.GeoPos("cities", name).Result()
	if err != nil {
		panic(err)
	}
	pos := positions[0]
	fmt.Printf("coordinates for %s: %f, %f\n", name, pos.Longitude, pos.Latitude)
}

func find(client *redis.Client, name string) {
	// look up coordinates for the city
	positions, err := client.GeoPos("cities", name).Result()
	if err != nil {
		panic(err)
	}
	pos := positions[0]

	// find closest peaks
	query := &redis.GeoRadiusQuery{
		Radius:   200,
		Unit:     "km",
		WithDist: true,
		Sort:     "ASC",
	}
	values, err := client.GeoRadius("peaks", pos.Longitude, pos.Latitude, query).Result()
	if err != nil {
		panic(err)
	}

	// print peaks
	fmt.Printf("Peaks closest to %s in the database:\n", name)
	for i, v := range values {
		fmt.Printf("(%d) %s, %.0f km\n", i+1, v.Name, v.Dist)
	}
}

func export(client *redis.Client, filename string) {
	// open CSV file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(file)

	// write header and data
	w.Write([]string{"name", "lat", "lon", "marker-color"})
	exportLocations(client, w, "cities", "#CD0000")
	exportLocations(client, w, "peaks", "#0000CD")

	// finish up writing CSV
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func exportLocations(client *redis.Client, w *csv.Writer, key, color string) {
	// get all members
	locations, err := client.ZRange(key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	for _, name := range locations {
		// load coordinates
		positions, err := client.GeoPos(key, name).Result()
		if err != nil {
			panic(err)
		}
		pos := positions[0]

		// write CSV line
		err = w.Write([]string{
			name,
			fmt.Sprintf("%f", pos.Latitude),
			fmt.Sprintf("%f", pos.Longitude),
			color,
		})
		if err != nil {
			panic(err)
		}
	}
}

func flushDB(client *redis.Client) {
	// delete all keys
	err := client.FlushDB().Err()
	if err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Println("Usage: go run main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("	add              add sample data for cities and peaks")
	fmt.Println("	lookup <city>    look up a city")
	fmt.Println("	find  <city>     find the peaks closest to a city")
	fmt.Println("	export <file>    export to CSV file suitable for geojson.io")
	fmt.Println("	flush            clear database")
	os.Exit(1)
}

func argument(i int) string {
	if len(os.Args) < i+1 {
		usage()
	}
	return os.Args[i]
}

var cities = []*redis.GeoLocation{
	{
		Name:      "Zurich",
		Latitude:  47.3775499,
		Longitude: 8.4666755,
	},
	{
		Name:      "Milan",
		Latitude:  45.462889,
		Longitude: 9.0376498,
	},
	{
		Name:      "Geneva",
		Latitude:  46.2050836,
		Longitude: 6.1090692,
	},
	{
		Name:      "Salzburg",
		Latitude:  47.802904,
		Longitude: 12.9863905,
	},
	{
		Name:      "Nice",
		Latitude:  43.7032932,
		Longitude: 7.1827775,
	},
}

var peaks = []*redis.GeoLocation{
	{
		Name:      "Mont Blanc",
		Latitude:  45.8326504,
		Longitude: 6.8476653,
	},
	{
		Name:      "Monte Rosa",
		Latitude:  45.9370551,
		Longitude: 7.8501157,
	},
	{
		Name:      "Matterhorn",
		Latitude:  45.9766029,
		Longitude: 7.6409423,
	},
	{
		Name:      "Grossglockner",
		Latitude:  12.6946761,
		Longitude: 47.0741846,
	},
	{
		Name:      "Wildspitze",
		Latitude:  46.8854563,
		Longitude: 10.8497499,
	},
}
