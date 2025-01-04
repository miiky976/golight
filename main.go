package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	BACKLIGHT_PATH string = "/sys/class/backlight/"
)

var (
	setValue int
	incValue int
)

func init() {
	flag.IntVar(&setValue, "set", -1, "Set the brightness value raw")
	flag.IntVar(&incValue, "inc", 0, "Increment brightness")
}

func main() {
	var err error
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Println("Run without args to get the actual brightness")
		flag.PrintDefaults()
	}
	var changed int
	dir := getDrivers(BACKLIGHT_PATH)
	flag.Parse()
	if setValue != -1 {
		err = set(setValue, dir[0])
		changed = setValue
	}
	if incValue == 0 {
		changed, err = inc(incValue, dir[0])
	} else {
		changed, err = inc(incValue, dir[0])
	}
	fmt.Println(changed)
	if err != nil {
		// TODO
		log.Fatal(err)
	}

}

func getDrivers(path string) []string {
	var drivers []string
	dirs, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range dirs {
		drivers = append(drivers, path+i.Name())
	}
	return drivers
}

func getMax(path string) (int, error) {
	maxFile, err := os.ReadFile(path + "/max_brightness")
	if err != nil {
		return 0, err
	}
	maxValue, err := strconv.Atoi(strings.Replace(string(maxFile), "\n", "", -1))
	if err != nil {
		return 0, err
	}
	return maxValue, nil
}

func getActual(path string) (int, error) {
	actualFile, err := os.ReadFile(path + "/actual_brightness")
	if err != nil {
		return 0, err
	}
	maxValue, err := strconv.Atoi(strings.Replace(string(actualFile), "\n", "", -1))
	if err != nil {
		return 0, err
	}
	return maxValue, nil
}

func set(value int, path string) error {
	maxim, err := getMax(path)
	if err != nil {
		return err
	}
	brightness, err := os.OpenFile(path+"/brightness", os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		return err
	}
	if value >= 0 && value <= maxim {
		_, err = brightness.Write([]byte(strconv.Itoa(value) + "\n"))
		if err != nil {
			return err
		}
	}
	err = brightness.Close()
	if err != nil {
		return err
	}
	return nil
}

func inc(value int, path string) (int, error) {
	maxim, err := getMax(path)
	if err != nil {
		return 0, err
	}
	actual, err := getActual(path)
	if err != nil {
		return 0, err
	}
	newValue := actual + value
	if newValue > maxim {
		set(maxim, path)
		return maxim, nil
	}
	if value < 0 && value > maxim {
		return actual, nil
	}
	set(newValue, path)
	return newValue, nil
}
