package device

import (
	"fmt"

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
	src al.Source

	buffer *audio.PCMBuffer
}

func NewSink(buffer *audio.PCMBuffer) (*Sink, error) {
	// TODO(jbd): only mono and stereo (8 and 16 bits are allowed)
	// return error if format is different.
	al.OpenDevice()
	s := (al.GenSources(1)[0])
	if code := al.Error(); code != 0 {
		return nil, fmt.Errorf("device: cannot generate an audio source [err=%x]", code)
	}
	return &Sink{src: s, buffer: buffer}, nil
}

func (s *Sink) Play() {
	buf := al.GenBuffers(1)[0]
	f := alFormat(s.buffer.Format)
	buf.BufferData(f, s.buffer.Bytes, int32(s.buffer.Format.SampleRate))
	s.src.QueueBuffers(buf)
	al.PlaySources(s.src)
}

func (s *Sink) Pause() {
	// TODO(jbd): return error, dont' return until EOF.
	al.PauseSources(s.src)
}

func (s *Sink) Stop() {
	// TODO(jbd): return error, dont' return until EOF.
	al.StopSources(s.src)
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
