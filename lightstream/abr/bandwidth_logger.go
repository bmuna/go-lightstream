// client/collector.go
package abr

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func MeasureBandwidth() float64 {
	url := "https://speed.hetzner.de/1MB.bin"

	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("❌ Error downloading:", err)
		return 0
	}
	defer resp.Body.Close()

	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Println("❌ Error reading body:", err)
		return 0
	}

	elapsed := time.Since(start).Seconds()
	bits := float64(n * 8)
	kbps := bits / elapsed / 1000
	return kbps
}

func LogToFile(kbps float64, file *os.File) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("%s,%.2f\n", timestamp, kbps)
	file.WriteString(logLine)
}
