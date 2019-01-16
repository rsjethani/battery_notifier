package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rsjethani/sysinfo"
)

var argWatch bool
var argThreshold uint
var argLowInterval time.Duration
var argNormalInterval time.Duration

func init() {
	flag.BoolVar(&argWatch, "w", false, "continously watch battery status")
	flag.DurationVar(&argLowInterval, "l", time.Minute*2, "battery check interval during low battery")
	flag.DurationVar(&argNormalInterval, "n", time.Minute*5, "battery check interval during good battery")
	flag.UintVar(&argThreshold, "t", 20, "threshold below which battery capacty would be considered critical")
}

func getBatteryStatus() (uint, string, error) {
	info, err := sysinfo.GetInfo("hardware", "battery")
	if err != nil {
		return 0, "", err
	}
	c, _ := info.Attribute(0, "CAPACITY")
	capacity, _ := c.(uint)
	s, _ := info.Attribute(0, "STATUS")
	state, _ := s.(string)

	return capacity, state, nil
}

func main() {
	flag.Parse()

	for {
		capacity, state, err := getBatteryStatus()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(capacity, state)

		// If -w not given, exit
		if !argWatch {
			break
		}

		sleepInterval := argNormalInterval
		if state == "Discharging" && capacity < argThreshold {
			sleepInterval = argLowInterval
			err = sendNotification(capacity, state)
			if err != nil {
				fmt.Println("an error on displaying notification occured: %s", err)
				os.Exit(2)
			}
		}
		time.Sleep(sleepInterval)
	}
}
