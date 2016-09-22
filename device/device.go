package device

import (
	"fmt"
	"time"

	"github.com/mattetti/audio"

	"golang.org/x/mobile/exp/audio/al"
)

const (
	none = iota
	play
	pause
	stop
)

type Sink struct {
	src *al.Source

	prevControl int
	control     int
}

func NewSink(buffer *audio.PCMBuffer) (*Sink, error) {
	// TODO(jbd): only mono and stereo (8 and 16 bits are allowed)
	// return error if format is different.
	if err := al.OpenDevice(); err != nil {
		return nil, err
	}

	s := (al.GenSources(1)[0])
	if code := al.Error(); code != 0 {
		return nil, fmt.Errorf("device: cannot generate an audio source [err=%x]", code)
	}
	sink := &Sink{src: &s}
	// make sure there is always n buffers in the main source
	go sink.run(buffer)
	return sink, nil
}

func (s *Sink) run(buffer *audio.PCMBuffer) {
	buf := al.GenBuffers(1)[0]
	for {
		src := *s.src
		if n := src.BuffersProcessed(); n > 0 || src.BuffersQueued() < 1 {
			f := alFormat(buffer.Format)
			buf.BufferData(f, buffer.AsBytes(), int32(buffer.Format.SampleRate))
			src.QueueBuffers(buf)
			continue
		}
		// control playing, pausing and stopping in the main loop.
		if s.prevControl != s.control {
			switch s.control {
			case play:
				al.PlaySources(src)
			case pause:
				al.PauseSources(src)
			case stop:
				al.StopSources(src)
			}
			s.prevControl = s.control
		}
		// TODO(jbd): if EOF, return.
		time.Sleep(50 * time.Millisecond)
	}
}

func (s *Sink) Play() {
	// TODO(jbd): return error, dont' return until EOF.
	s.control = play
}

func (s *Sink) Pause() {
	// TODO(jbd): return error, dont' return until EOF.
	s.control = pause
}

func (s *Sink) Stop() {
	// TODO(jbd): return error, dont' return until EOF.
	s.control = stop
}

func alFormat(f *audio.Format) uint32 {
	switch f.NumChannels {
	case 1:
		switch f.BitDepth {
		case 8:
			return al.FormatMono8
		case 16:
			return al.FormatMono16
		}
	case 2:
		switch f.BitDepth {
		case 8:
			return al.FormatStereo8
		case 16:
			return al.FormatStereo16
		}
	}
	panic("unsupported format")
}
