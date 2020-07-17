package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"sync"

	"github.com/sparrc/go-ping"
)

type IPAddress struct {
	name    string
	address string
}

func main() {

	x := readAddresses()
	// for _, y := range x{
	// 	fmt.Println(y.address, y.name)
	// }
	trackPing(x)
}

func pingTest() {
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		fmt.Println(err)
	}
	pinger.SetPrivileged(true)

	pinger.Count = 3
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	fmt.Println(stats)
}

func readAddresses(inputNames ...string) (outputAddresses []IPAddress) {

	raw, err := ioutil.ReadFile("knownaddresses.txt")
	if err != nil {
		panic(err)
	}

	var ipSlice []IPAddress
	addressFileEntries := strings.Split(string(raw), "\n")

	for _, entry := range addressFileEntries {
		x := strings.Split(entry, " ")
		tempAddress := IPAddress{x[0], strings.Trim(x[1], "\n\r")}
		ipSlice = append(ipSlice, tempAddress)
	}

	if len(inputNames) != 0 {
		var outputIPSlice []IPAddress
		for _, entry := range ipSlice {
			for _, input := range inputNames {
				if input == entry.name {
					outputIPSlice = append(outputIPSlice, entry)
				}
			}
		}
		ipSlice = outputIPSlice
	}
	return ipSlice

}

func logResult() {

}

func logFile() (fileName string) {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModeDir)
	}
	currentTime := time.Now()
	timeFormatted := currentTime.Format("2006-01-02_15-04-05")
	fileName = "logs/" + timeFormatted + ".txt"
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	l, err := f.WriteString("Log started at: " + currentTime.Format("2006-01-02 15:04:05\n"))
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	f.Close()
	log.Println(l, "bytes written successfully")
	return fileName
}

func trackPing(addresses []IPAddress) {
	fileName := logFile()
	fmt.Println(fileName)
	var wg sync.WaitGroup

	for _, address := range addresses {
		wg.Add(1)
		go trackAddress(address, fileName, &wg)
	}
	wg.Wait()

}

func trackAddress(target IPAddress, logFile string, wg *sync.WaitGroup) {

	defer wg.Done()
	timeStart := time.Now()
	var logger []*ping.Statistics
	for {
		pinger, err := ping.NewPinger(target.address)
		if err != nil {
			fmt.Println(err)
		}
		pinger.SetPrivileged(true)
		pinger.Count = 3
		pinger.Run()                 // blocks until finished
		stats := pinger.Statistics() // get send/receive/rtt stats
		fmt.Println(target.name, stats.AvgRtt)
		logger = append(logger, stats)
		if len(logger) == 6 {

			//log shit
			//empty log
			appendLog(pingSliceToString(target, logger, timeStart), logFile) 
			logger = nil
		}
		time.Sleep(10)
	}

}

func pingSliceToString(target IPAddress, log []*ping.Statistics, timeStart time.Time) string {
	formattedTime := strings.Split(timeStart.Format("2006-01-02 3:4:5 PM"), " ")
	clockTime := formattedTime[1] + formattedTime[2]
	var tempString string
	tempString += fmt.Sprintf("%s %s %v %v %v %v %v %v", 
					target.name, clockTime, log[0].AvgRtt, log[1].AvgRtt, log[2].AvgRtt, 
					log[3].AvgRtt, log[4].AvgRtt, log[5].AvgRtt)

	return tempString
}

func appendLog(message string, logFileName string) {
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Opening Error:", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(message + "\n"); err != nil {
		log.Println(err)
	}

}
