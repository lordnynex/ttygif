package main

import (
	"fmt"
	"image/gif"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sugyan/ttyread"
)

// GifGenerator type
type GifGenerator struct {
	Speed  float64
	NoLoop bool
}

// NewGifGenerator returns GifGenerator instance
func NewGifGenerator() *GifGenerator {
	return &GifGenerator{Speed: 1.0}
}

// Generate writes to outFile an animated GIF
func (g *GifGenerator) Generate(inFile string, outFile string) (err error) {
	tempDir, err := ioutil.TempDir("", "ttygif")
	if err != nil {
		log.Println("FNR: 28")
		return
	}
	defer os.RemoveAll(tempDir)

	var (
		delays    []int
		timestamp ttyread.TimeVal
	)
	worker := NewWorker()
	// play and capture
	err = Play(inFile, func(diff ttyread.TimeVal) (err error) {
		delay := int(float64(diff.Sec*1000000+diff.Usec)/g.Speed) / 10000
		if delay == 0 {
			return nil
		}
		delays = append(delays, delay)
		timestamp = timestamp.Add(diff)

		// capture and append to images
		imgPath := filepath.Join(tempDir, fmt.Sprintf("%03d_%06d", timestamp.Sec, timestamp.Usec))
		fileType, err := CaptureImage(imgPath)
		if err != nil {
			log.Println("FNR: 53")
			return
		}
		worker.AddTargetFile(imgPath, fileType)
		return nil
	})
	if err != nil {
		log.Println("FNR: 60")
		return
	}
	// get paletted images from capture files
	progress := make(chan struct{})
	go func() {
	Loop:
		for {
			select {
			case <-time.Tick(time.Millisecond * 500):
				print(".")
			case <-progress:
				break Loop
			}
		}
		print("\r")
	}()
	images, err := worker.GetAllImages()
	if err != nil {
		log.Println("FNR: 79")
		return
	}
	close(progress)

	// generate GIF file
	file, err := os.Create(outFile)
	if err != nil {
		log.Println("FNR: 87")
		return
	}
	defer file.Close()
	opts := gif.GIF{
		Image: images,
		Delay: delays,
	}
	if g.NoLoop {
		opts.LoopCount = 1
	}
	err = gif.EncodeAll(file, &opts)
	if err != nil {
		log.Println("FNR: 100")
		return
	}
	return nil
}
