package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gordonklaus/portaudio"
)

func captureVoIPAudio(done chan bool) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	// Create file for raw audio
	file, err := os.Create("test/audio.raw")
	if err != nil {
		fmt.Println("Error creating test/audio.raw:", err)
		return
	}
	defer file.Close()

	// Open default stream
	buffer := make([]int16, 512)
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, len(buffer), buffer)
	if err != nil {
		fmt.Println("Error opening default stream:", err)
		return
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		fmt.Println("Error starting stream:", err)
		return
	}
	defer stream.Stop()

	fmt.Println("Capturing VoIP audio from default device and saving to test/audio.raw...")
	for {
		select {
		case <-done:
			return
		default:
			if err := stream.Read(); err != nil {
				fmt.Println("Error reading audio:", err)
				return
			}
			fmt.Printf("VoIP: %d samples from default device\n", len(buffer))
			for _, sample := range buffer {
				if _, err := file.Write([]byte{byte(sample & 0xff), byte(sample >> 8)}); err != nil {
					fmt.Println("Error writing to test/audio.raw:", err)
					return
				}
			}
		}
	}
}

func recordAV(done chan bool) {
	cmd := exec.Command("ffmpeg",
		"-f", "avfoundation",
		"-framerate", "30",
		"-video_size", "1280x720",
		"-i", "0:0", // Default video and audio devices
		"-c:v", "libx264",
		"-c:a", "aac",
		"-ar", "16000",
		"-af", "volume=10",
		"-movflags", "+faststart",
		"-y",
		"test/recording.mp4",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting FFmpeg:", err)
		return
	}

	fmt.Println("Recording video + audio to test/recording.mp4 using default devices...")
	<-done
	fmt.Println("Stopping FFmpeg...")
	stdin.Write([]byte("q"))
	stdin.Close()
	if err := cmd.Wait(); err != nil {
		fmt.Println("FFmpeg exited with error:", err)
	}
	fmt.Println("Recording stopped")
}

func main() {
	done := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go captureVoIPAudio(done)
	go recordAV(done)

	fmt.Println("Running... Press Ctrl+C to stop")
	<-sigs
	close(done)
	time.Sleep(5 * time.Second) // Extended time for FFmpeg to finalize
}
