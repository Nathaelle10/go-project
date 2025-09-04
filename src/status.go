package main

import (
	"encoding/json"
	"net/http"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	psnet "github.com/shirou/gopsutil/v3/net"
)

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

func StatusHandler(w http.ResponseWriter, r *http.Request) {
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
	cores, _ := cpu.Counts(true)
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
