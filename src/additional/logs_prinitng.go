package additional

import (
	"fmt"
	"time"

	"github.com/utherbit/ansi"
)

var (
	cyanColor   = ansi.Format.Colors.Foreground.Standard.Cyan
	redColor    = ansi.Format.Colors.Foreground.Standard.Red
	yellowColor = ansi.Format.Colors.Foreground.Standard.Yellow
	greenColor  = ansi.Format.Colors.Foreground.Standard.Green
	reset       = ansi.Format.Reset
)

const (
	timeWidth   = 15
	ipPathWidth = 55
	levelWidth  = 8
)

func formatLog(timeStr, ipPath, level, message, color string) string {
	return fmt.Sprintf("%-*s | %-*s %s| %-*s | %s%s\n",
		timeWidth, timeStr,
		ipPathWidth, ipPath,
		color, levelWidth,
		level, message, reset,
	)
}

func PrintError(path string, err error) {
	tn := time.Now().Format("Jan _2 15:04:05")
	IPv4 := GetLocalIP()
	fmt.Print(formatLog(tn, IPv4+path, "ERROR", err.Error(), redColor))
}

func PrintSuccess(path string, message string) {
	tn := time.Now().Format("Jan _2 15:04:05")
	IPv4 := GetLocalIP()
	fmt.Print(formatLog(tn, IPv4+path, "SUCCESS", message, greenColor))
}

func PrintWarn(path string, message string) {
	tn := time.Now().Format("Jan _2 15:04:05")
	IPv4 := GetLocalIP()
	fmt.Print(formatLog(tn, IPv4+path, "WARN", message, yellowColor))
}

func PrintMessage(message string) {
	tn := time.Now().Format("Jan _2 15:04:05")
	IPv4 := GetLocalIP()
	fmt.Print(formatLog(tn, IPv4, "SERVER", message, cyanColor))
}

func PrintServerError(message string) {
	tn := time.Now().Format("Jan _2 15:04:05")
	IPv4 := GetLocalIP()
	fmt.Print(formatLog(tn, IPv4, "ERROR", message, redColor))
}
