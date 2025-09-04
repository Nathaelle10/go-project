package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"
)

func ScanPort(protocol, hostname string, port int) bool {
	target := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, target, 300*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func ScanHandler(w http.ResponseWriter, r *http.Request) {
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
