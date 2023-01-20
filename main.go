package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DemoInfo represents data.
type DemoInfo struct {
	ID       string `json:"id"`
	HostName string `json:"hostname"`
	IP       string `json:"ip"`
	DateTime string `json:"datetime"`
}

// albums slice to seed record album data.

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", getInfo)
	fmt.Println("Server start...")
	router.Run(":8080")
}

// getAlbums responds with the list of all albums as JSON.
func getInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getInfoDetail())
}

func getInfoDetail() *DemoInfo {
	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	currentTime := time.Now()
	var id = strconv.Itoa(rand.Int())
	return &DemoInfo{ID: id, HostName: hostName, IP: GetOutboundIP().String(), DateTime: currentTime.Format("2006.01.02 15:04:05")}

}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
