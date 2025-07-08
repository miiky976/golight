package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"strconv"
)

const (
	BACKLIGHT_PATH   = "/sys/class/backlight/"
	BACKLIGHT_ACTUAL = "brightness"
	BACKLIGHT_MAX    = "max_brightness"
	BACKLIGHT        = "brightness"
	DESC             = "Modify backlight brightness."
	USAGE            = "Usage of %s:\n"
	USAGE_DATA       = "Run without args to get the actual brightness level"
)

var (
	setValue int
	id       int
	idp      int
	driver   string
	list     bool
	raw      bool
)

func init() {
	flag.IntVar(&setValue, "set", -1, "Set the brightness value raw")
	flag.IntVar(&id, "id", 0, "Increment or Decrease brightness")
	flag.IntVar(&idp, "idp", 0, "Increment or Decrease brightness by Percentage")
	flag.StringVar(&driver, "driver", "", "Number of driver to modify, by default it modifies the first (literally the index)")
	flag.BoolVar(&list, "list", false, "List the available drivers")
	flag.BoolVar(&raw, "raw", false, "Return the value as raw (by default its percentage)")
}

func main() {
	flag.Usage = func() {
		fmt.Println(DESC)
		fmt.Printf(USAGE, os.Args[0])
		fmt.Println(USAGE_DATA)
		flag.PrintDefaults()
	}
	drivers, err := getDrivers(BACKLIGHT_PATH)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	printDrivers(drivers)
	setDriver(drivers)
	setLevel()
	idLevel()
	idpLevel()
	if raw {
		fmt.Println(getLevel())
	} else {
		p := (getLevel() * 100) / getMax()
		fmt.Println(p)
	}
}

func printDrivers(drivers []string) {
	if list {
		for i, d := range drivers {
			fmt.Printf("%d - %s\n", i, d)
		}
		os.Exit(0)
	}
}

func setDriver(drivers []string) {
	if driver == "" {
		driver = drivers[0]
	} else {
		d, err := strconv.Atoi(driver)
		if err != nil {
			panic(err)
		}
		if d < 0 && d > len(drivers) {
			panic("Cmon that's not a valid index")
		}
		driver = drivers[d]
	}
}

func getLevel() int {
	act := path.Join(driver, BACKLIGHT_ACTUAL)
	data, err := os.ReadFile(act)
	if err != nil {
		panic(err)
	}
	val := string(data[:len(data)-1])
	actual, _ := strconv.Atoi(val)
	return actual
}

func getMax() int {
	mx := path.Join(driver, BACKLIGHT_MAX)
	data, err := os.ReadFile(mx)
	if err != nil {
		panic(err)
	}
	val := string(data[:len(data)-1])
	maxim, _ := strconv.Atoi(val)
	return maxim
}

func getDrivers(path string) ([]string, error) {
	var drivers []string
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, i := range dirs {
		drivers = append(drivers, path+i.Name())
	}
	return drivers, nil
}

func setLevel() {
	if setValue < 0 {
		return
	}
	set := path.Join(driver, BACKLIGHT)
	data, err := os.OpenFile(set, os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	setValue = min(getMax(), setValue)
	_, err = data.Write([]byte(strconv.Itoa(setValue)))
	if err != nil {
		panic(err)
	}
}

func idLevel() {
	if id == 0 {
		return
	}
	actual := getLevel()
	actual += id
	setValue = max(0, actual)
	setLevel()
}

func idpLevel() {
	if idp == 0 {
		return
	}
	maxim := getMax()
	actual := getLevel()
	pt := float64(maxim) / 100
	val := pt * float64(idp)
	// This is to make it a little more accurate
	// because when 255 is the maximum it leads to inconsistencies when rounding
	if idp < 0 {
		actual += int(math.Floor(val))
	} else {
		actual += int(math.Ceil(val))
	}
	setValue = max(0, actual)
	setLevel()
}
