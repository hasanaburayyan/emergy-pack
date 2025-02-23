package main

import (
	"fmt"
	"os/exec"
)

func openMP4() {
	cmd := exec.Command("open", "test/recording.mp4")
	if err := cmd.Run(); err != nil {
		fmt.Println("Error opening test/recording.mp4:", err)
	} else {
		fmt.Println("Opened test/recording.mp4 in default player")
	}
}

func main() {
	openMP4()
}
