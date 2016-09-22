package main

import (
	"log"
	"os"
	"time"

	"github.com/mattetti/audio/device"
	"github.com/mattetti/audio/wav"
)

func main() {
	f, err := os.Open("/Users/jbd/Desktop/chord.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		log.Fatal(err)
	}

	dev, err := device.NewSink(buf)
	if err != nil {
		log.Fatal(err)
	}
	dev.Play()

	time.Sleep(5 * time.Second)
	dev.Pause()
	time.Sleep(time.Second)
	dev.Play()

	time.Sleep(time.Minute)
}
