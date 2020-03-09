package engine

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type SoundEngine struct {
	on        bool
	soundFile *os.File
	buffer    *beep.Buffer
	isPlaying bool
}

func (se *SoundEngine) Init(on bool, filePath string) error {
	se.on = on

	// Open MP3 File when engine is on
	var err error
	if on {
		se.soundFile, err = os.Open(filePath)
		if err != nil {
			return err
		}

		streamer, format, err := mp3.Decode(se.soundFile)
		if err != nil {
			return err
		}

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		se.buffer = beep.NewBuffer(format)
		se.buffer.Append(streamer)
		streamer.Close()
	}

	return nil
}

func (se *SoundEngine) Play() error {
	if se.on {

		if se.isPlaying {
			return nil
		}

		se.isPlaying = true
		speaker.Play(se.buffer.Streamer(0, se.buffer.Len()))
		se.isPlaying = false
	}
	return nil
}

func (se *SoundEngine) IsPlaying() bool {
	return se.isPlaying
}

func (se *SoundEngine) Close() {
	defer se.soundFile.Close()
}
