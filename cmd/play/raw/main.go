package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gordonklaus/portaudio"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Read the raw file
	data, err := os.ReadFile("test/audio.raw")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert byte slice to int16 slice
	samples := make([]int16, len(data)/2)
	for i := 0; i < len(data)/2; i++ {
		samples[i] = int16(data[i*2]) | (int16(data[i*2+1]) << 8)
	}

	// Open an output stream (0 input, 1 output channel)
	stream, err := portaudio.OpenDefaultStream(0, 1, 16000, 1024, samples)
	if err != nil {
		fmt.Println("Error opening stream:", err)
		return
	}
	defer stream.Close()

	// Start and play
	if err := stream.Start(); err != nil {
		fmt.Println("Error starting stream:", err)
		return
	}
	fmt.Println("Playing...")

	// Write all samples (blocks until done)
	if err := stream.Write(); err != nil {
		fmt.Println("Error writing audio:", err)
		return
	}

	// Wait for playback to finish (rough estimate based on sample rate)
	<-time.After(time.Duration(len(samples)) * time.Second / 16000)
	fmt.Println("Done!")
}
