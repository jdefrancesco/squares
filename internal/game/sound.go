package game

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// eatSoundPCM generates a quiet, ascending 100-millisecond chime as stereo
// float32 little-endian PCM, with an envelope to avoid clicks at the endpoints.
func eatSoundPCM(sampleRate int) []byte {
	frames := sampleRate / 10
	data := make([]byte, frames*8)
	for i := 0; i < frames; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(frames-1)
		envelope := math.Sin(math.Pi*progress) * math.Exp(-3*progress)
		phase := 2 * math.Pi * (660*t + 1800*t*t)
		sample := float32(0.22 * math.Sin(phase) * envelope)
		bits := math.Float32bits(sample)
		binary.LittleEndian.PutUint32(data[i*8:], bits)
		binary.LittleEndian.PutUint32(data[i*8+4:], bits)
	}
	return data
}

// playEatSound initializes sound lazily inside the running game and reuses one
// player, keeping rapid eating feedback from accumulating overlapping audio.
func (g *Game) playEatSound() {
	if g.eatSound == nil {
		ctx := audio.CurrentContext()
		if ctx == nil {
			ctx = audio.NewContext(44100)
		}
		g.eatSound = ctx.NewPlayerF32FromBytes(eatSoundPCM(ctx.SampleRate()))
	}
	if err := g.eatSound.Rewind(); err != nil {
		log.Printf("Squares: could not replay eating sound: %v", err)
		return
	}
	g.eatSound.Play()
}
