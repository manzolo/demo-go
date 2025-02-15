package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// DemoInfo rappresenta le informazioni relative al container.
type DemoInfo struct {
	ID         string `json:"id"`
	HostName   string `json:"hostname"`
	IP         string `json:"ip"`
	DateTime   string `json:"datetime"`
	AppVersion string `json:"app_version"`
}

// SystemInfo raccoglie alcune informazioni di sistema del server.
type SystemInfo struct {
	OS              string `json:"os"`
	Architecture    string `json:"architecture"`
	CPUs            int    `json:"cpus"`
	GoVersion       string `json:"go_version"`
	ProcessID       int    `json:"process_id"`
	WorkingDir      string `json:"working_directory"`
	CurrentTime     string `json:"current_time"`
	MemoryAllocated uint64 `json:"memory_allocated_bytes"`       // bytes attualmente allocati
	MemoryTotal     uint64 `json:"memory_total_allocated_bytes"` // bytes allocati totali nel tempo
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Configura il middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permetti tutte le origini
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoint per ottenere le info base del container.
	router.GET("/", getInfo)

	// Endpoint che riceve un parametro reale ed esegue un semplice calcolo.
	router.GET("/calc/:value", calcWithParam)

	// Endpoint che restituisce tutte le info di sistema.
	router.GET("/system", getSystemInfoHandler)

	// Endpoint che, in base al parametro ricevuto, restituisce una specifica informazione di sistema.
	router.GET("/system/:param", getSystemParamHandler)

	fmt.Println("Server start...")
	router.Run(":8080")
}

// getInfo restituisce le informazioni di base del container.
func getInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getInfoDetail())
}

// calcWithParam gestisce la richiesta che include un parametro reale ed esegue un calcolo (moltiplicazione per un fattore fisso).
func calcWithParam(c *gin.Context) {
	valueStr := c.Param("value")
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parametro non valido. Deve essere un numero reale."})
		return
	}

	// Esempio di calcolo: moltiplichiamo il valore per 1.5.
	const factor = 1.5
	result := value * factor

	response := gin.H{
		"input_value":        value,
		"calculation_factor": factor,
		"result":             result,
		"message":            fmt.Sprintf("Calcolo effettuato: %v x %v = %v", value, factor, result),
		"container_info":     getInfoDetail(),
	}
	c.IndentedJSON(http.StatusOK, response)
}

// getSystemInfoHandler restituisce tutte le informazioni di sistema.
func getSystemInfoHandler(c *gin.Context) {
	sysInfo, err := getSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore nel recupero delle info di sistema"})
		return
	}
	c.IndentedJSON(http.StatusOK, sysInfo)
}

// getSystemParamHandler restituisce una specifica informazione di sistema in base al parametro fornito.
func getSystemParamHandler(c *gin.Context) {
	param := c.Param("param")
	sysInfo, err := getSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Errore nel recupero delle info di sistema"})
		return
	}

	var response interface{}
	switch param {
	case "cpu":
		response = sysInfo.CPUs
	case "goversion":
		response = sysInfo.GoVersion
	case "os":
		response = sysInfo.OS
	case "arch":
		response = sysInfo.Architecture
	case "pid":
		response = sysInfo.ProcessID
	case "workingdir":
		response = sysInfo.WorkingDir
	case "mem":
		response = sysInfo.MemoryAllocated
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parametro non riconosciuto. Usa: cpu, goversion, os, arch, pid, workingdir, mem",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{param: response})
}

// getInfoDetail raccoglie le informazioni di base del container.
func getInfoDetail() *DemoInfo {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Unknown"
	}

	currentTime := time.Now()
	id := strconv.Itoa(rand.Intn(1000)) // ID casuale per la demo

	// Recupera le variabili d'ambiente relative al container.
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "1.0.0"
	}

	return &DemoInfo{
		ID:         id,
		HostName:   hostName,
		IP:         GetOutboundIP().String(),
		DateTime:   currentTime.Format("2006-01-02 15:04:05"),
		AppVersion: appVersion,
	}
}

// getSystemInfo raccoglie informazioni di sistema dal server.
func getSystemInfo() (*SystemInfo, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	cpus := runtime.NumCPU()
	goVersion := runtime.Version()
	pid := os.Getpid()
	wd, err := os.Getwd()
	if err != nil {
		wd = "Unknown"
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Lettura delle statistiche di memoria.
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	sysInfo := &SystemInfo{
		OS:              osName,
		Architecture:    arch,
		CPUs:            cpus,
		GoVersion:       goVersion,
		ProcessID:       pid,
		WorkingDir:      wd,
		CurrentTime:     currentTime,
		MemoryAllocated: memStats.Alloc,
		MemoryTotal:     memStats.TotalAlloc,
	}
	return sysInfo, nil
}

// GetOutboundIP restituisce l'indirizzo IP outbound della macchina.
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
