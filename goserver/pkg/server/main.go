package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//ProcsInfo representa la sumatoria de procesos del so por estado
type ProcsInfo struct {
	Runing   int
	Sleeping int
	Stoped   int
	Zombie   int
	Other    int
	Total    int
}

//ProcInfo Representa el estado de un proceso del so
type ProcInfo struct {
	Pid                       int
	Nombre                    string
	Usuario                   string
	EstadoID                  string
	Estado                    string
	MemoriaUtiliada           uint64
	PorcentajeMemoriaUtiliada float64
}

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

//función main
func main() {
	http.HandleFunc("/memo", memoHandler)
	http.HandleFunc("/cpu", cpuHandler)
	http.HandleFunc("/proc", procHandler)
	http.HandleFunc("/procs", procsHandler)
	http.HandleFunc("/procsinfo", procsInfoHandler)
	http.HandleFunc("/killproc", killProcHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//request handlers

func killProcHandler(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["pid"]
	pid, _ := strconv.Atoi(keys[0])
	var procInfo ProcInfo
	err := GetProcInfo(&procInfo, pid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, errCmd := exec.Command("kill", strconv.Itoa(pid)).Output()

	if errCmd != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(procInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

var gettingProcsInfo bool = false
var cachedProcsInfo []ProcInfo

func procsInfoHandler(w http.ResponseWriter, r *http.Request) {
	var procs []ProcInfo
	var err error
	if gettingProcsInfo {
		procs = cachedProcsInfo
		err = nil
	} else {
		procs, err = GetProcsInfo()
		cachedProcsInfo = procs
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var procsInfo ProcsInfo

	for _, proc := range procs {
		procsInfo.Total++

		switch proc.EstadoID {
		case "R":
			procsInfo.Runing++
		case "S":
			procsInfo.Sleeping++
		case "Z":
			procsInfo.Zombie++
		case "T":
			procsInfo.Stoped++
		default:
			procsInfo.Other++
		}
	}

	js, err := json.Marshal(procsInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func procsHandler(w http.ResponseWriter, r *http.Request) {

	var procs []ProcInfo
	var err error
	if gettingProcsInfo {
		procs = cachedProcsInfo
		err = nil
	} else {
		procs, err = GetProcsInfo()
		cachedProcsInfo = procs
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(procs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func procHandler(w http.ResponseWriter, r *http.Request) {

	keys, _ := r.URL.Query()["pid"]
	pid, _ := strconv.Atoi(keys[0])
	var procInfo ProcInfo
	err := GetProcInfo(&procInfo, pid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(procInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

//Util functions

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

//GetMemInfo obtiene el estado de la memoria RAM
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

//GetProcsInfo obtiene el estado de los procesos del SO
func GetProcsInfo() ([]ProcInfo, error) {
	gettingProcsInfo = true
	//información del estado de la memoria
	var memoInfo MemoInfo
	err := GetMemInfo(&memoInfo)

	if err != nil {
		gettingProcsInfo = false
		return nil, err
	}

	//listado de procesos que se retornará
	var procs []ProcInfo

	directory := "/proc/"

	// abrir directorio
	outputDirRead, err := os.Open(directory)

	// obtener contenido del directorio
	outputDirFiles, err := outputDirRead.Readdir(0)

	if err != nil {
		gettingProcsInfo = false
		return nil, err
	}

	// iterar sobre el contenido del directorio
	for outputIndex := range outputDirFiles {
		info := outputDirFiles[outputIndex]

		if info.IsDir() {
			pid, err := strconv.Atoi(info.Name())

			if err == nil {
				var proc ProcInfo
				errGet := GetProcInfo(&proc, pid)

				if errGet == nil {
					proc.PorcentajeMemoriaUtiliada = float64(proc.MemoriaUtiliada) / float64(memoInfo.AvailableKb)
					procs = append(procs, proc)
				}
			}
		}
	}
	gettingProcsInfo = false
	return procs, nil
}

/*
type ProcInfo struct {
	Pid             int
	Nombre          string
	Usuario         string
	EstadoID        string
	Estado          string
	MemoriaUtiliada uint64
}
*/

//GetProcInfo obtiene el estado de un proceso del SO
func GetProcInfo(p *ProcInfo, pid int) error {
	p.Pid = pid

	var err error

	path := filepath.Join("/proc/" + strconv.Itoa(pid) + "/status")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()

		n := strings.Index(text, ":")

		key := text[:n] // metric

		if key == "Name" {
			p.Nombre = strings.TrimSpace(text[(n + 1):])

		} else if key == "State" {
			data := strings.Split(strings.TrimSpace(text[(n+1):]), " ")
			p.EstadoID = data[0]
			p.Estado = data[1]

		} else if key == "VmRSS" {
			data := strings.Split(strings.TrimSpace(text[(n+1):]), " ")
			value, err := strconv.ParseUint(data[0], 10, 64)

			if err != nil {
				return nil
			}

			p.MemoriaUtiliada = value
		}

	}

	out, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(pid)).Output()
	p.Usuario = strings.TrimSpace(string(out))

	return nil

}
