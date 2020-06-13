package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

func main() {
	http.HandleFunc("/memo", memoHandler)
	http.HandleFunc("/cpu", cpuHandler)
	http.HandleFunc("/proc", procHandler)
	http.HandleFunc("/procs", procsHandler)
	http.ListenAndServe(":8080", nil)
}

func procsHandler(w http.ResponseWriter, r *http.Request) {

	procs, err := GetProcsInfo()

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
	//información del estado de la memoria
	var memoInfo MemoInfo
	err := GetMemInfo(&memoInfo)

	if err != nil {
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
