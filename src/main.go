// package src
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	netFromNet "net"
	"net/http"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
)

//go:embed www/*
var content embed.FS

// --- Structures et Fonctions API ---

type SystemStatus struct {
	CPUUsage    float64                `json:"cpu_usage"`
	CPUCores    int                    `json:"cpu_cores"`
	CPUModel    string                 `json:"cpu_model"`
	Memory      *mem.VirtualMemoryStat `json:"memory"`
	Swap        *mem.SwapMemoryStat    `json:"swap"`
	Disk        *disk.UsageStat        `json:"disk"`
	Network     []psnet.IOCountersStat `json:"network"`
	Load        *load.AvgStat          `json:"load"`
	Host        *host.InfoStat         `json:"host"`
	Uptime      uint64                 `json:"uptime"`
	Connections int                    `json:"connections"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		http.Error(w, "Erreur récupération CPU", http.StatusInternalServerError)
		return
	}
	cpuInfo, err := cpu.Info()
	if err != nil || len(cpuInfo) == 0 {
		http.Error(w, "Erreur récupération CPU info", http.StatusInternalServerError)
		return
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		http.Error(w, "Erreur récupération nombre de cores", http.StatusInternalServerError)
		return
	}
	vmem, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()
	diskUsage, _ := disk.Usage("/")
	netIO, _ := psnet.IOCounters(true)
	loadAvg, _ := load.Avg()
	hostInfo, _ := host.Info()
	uptime, _ := host.Uptime()
	conns, _ := psnet.Connections("tcp")

	status := SystemStatus{
		CPUUsage:    cpuPercent[0],
		CPUCores:    cores,
		CPUModel:    cpuInfo[0].ModelName,
		Memory:      vmem,
		Swap:        swap,
		Disk:        diskUsage,
		Network:     netIO,
		Load:        loadAvg,
		Host:        hostInfo,
		Uptime:      uptime,
		Connections: len(conns),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func ScanPort(protocol, hostname string, port int) bool {
	target := hostname + ":" + strconv.Itoa(port)
	conn, err := netFromNet.DialTimeout(protocol, target, 300*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		http.Error(w, "Paramètre ?ip= requis", http.StatusBadRequest)
		return
	}

	const startPort = 20
	const endPort = 1024
	const maxConcurrency = 100

	type portResult struct {
		Port int  `json:"port"`
		Open bool `json:"open"`
	}

	results := make(chan portResult, endPort-startPort+1)
	sem := make(chan struct{}, maxConcurrency)

	for port := startPort; port <= endPort; port++ {
		sem <- struct{}{}
		go func(p int) {
			defer func() { <-sem }()
			open := ScanPort("tcp", ip, p)
			results <- portResult{Port: p, Open: open}
		}(port)
	}

	openPorts := []int{}
	for i := startPort; i <= endPort; i++ {
		res := <-results
		if res.Open {
			openPorts = append(openPorts, res.Port)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ip":        ip,
		"openPorts": openPorts,
	})
}

// Middleware CORS
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- Serveurs HTTP ---

func startAPIServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", statusHandler)
	mux.HandleFunc("/api/scan", scanHandler)

	fmt.Println("API démarrée sur http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", withCORS(mux)))
}

func startWebServer() {
	//fs := http.FileServer(http.Dir("./www"))
	http.Handle("/", http.FileServer(http.FS(content)))
	//http.Handle("/", fs)

	port := 8080
	fmt.Printf("Dashboard dispo sur http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	go startAPIServer() // API sur :9090
	startWebServer()    // UI sur :8080 (bloquant)
}
