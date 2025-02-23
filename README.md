# Emergy-Pack

## Dev Setup
1. `brew install portaudio`
2. `brew install ffmpeg`
3. `go get -u github.com/gordonklaus/portaudio` 
4. `go get -u github.com/go-audio/wav`


## Run
1. Record using `go run cmd/record/main.go`
2. Play audio only using `go run cmd/play/raw/main.go`
3. Play video using `go run cmd/play/video/main.go`