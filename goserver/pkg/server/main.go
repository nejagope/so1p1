package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//MemoInfo Representa el estado de la memoria RAM
type MemoInfo struct {
	TotalKb     uint64
	AvailableKb uint64
	UsedKb      uint64

	TotalMb     float64
	AvailableMb float64
	UsedMb      float64
}

//CPUInfo Representa el estado de la CPU
type CPUInfo struct {
	Used float64
}

func main() {
	http.HandleFunc("/memo", memoHandler)
	http.HandleFunc("/cpu", cpuHandler)
	http.ListenAndServe(":8080", nil)
}

func memoHandler(w http.ResponseWriter, r *http.Request) {
	var memoInfo MemoInfo
	err := GetMemInfo(&memoInfo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	memoInfo.UsedKb = memoInfo.TotalKb - memoInfo.AvailableKb
	memoInfo.AvailableMb = float64(memoInfo.AvailableKb) / 1000
	memoInfo.TotalMb = float64(memoInfo.TotalKb) / 1000
	memoInfo.UsedMb = memoInfo.TotalMb - memoInfo.AvailableMb

	js, err := json.Marshal(memoInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func cpuHandler(w http.ResponseWriter, r *http.Request) {
	idle0, total0 := getCPUSample()
	time.Sleep(1 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpuUsage := (totalTicks - idleTicks) / totalTicks

	cpuInfo := CPUInfo{cpuUsage}
	js, err := json.Marshal(cpuInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

//GetMemInfo Updates s with current values, usign the pid stored in the Stat
func GetMemInfo(m *MemoInfo) error {
	var err error

	path := filepath.Join("/proc/meminfo")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		n := strings.Index(text, ":")
		if n == -1 {
			continue
		}

		key := text[:n] // metric
		data := strings.Split(strings.Trim(text[(n+1):], " "), " ")
		if len(data) == 2 {
			if data[1] == "kB" {
				value, err := strconv.ParseUint(data[0], 10, 64)
				if err != nil {
					continue
				}

				if key == "MemTotal" {
					m.TotalKb = value
				} else if key == "MemAvailable" {
					m.AvailableKb = value
				}
			}
		}

	}
	return nil

}
