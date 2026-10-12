package game

import (
	"encoding/binary"
	"math"
	"testing"
)

// TestEatSoundHasStereoFade checks the generated sound format and click-free ends.
func TestEatSoundHasStereoFade(t *testing.T) {
	data := eatSoundPCM(44100)
	frames := len(data) / 8
	if frames < 3000 || frames > 6000 || len(data)%8 != 0 {
		t.Fatalf("unexpected stereo sound length: %d", len(data))
	}
	nonzero := false
	for i := 0; i < frames; i++ {
		left := math.Float32frombits(binary.LittleEndian.Uint32(data[i*8:]))
		right := math.Float32frombits(binary.LittleEndian.Uint32(data[i*8+4:]))
		if left != right || math.Abs(float64(left)) > 0.3 {
			t.Fatal("sound should be quiet, valid stereo PCM")
		}
		nonzero = nonzero || left != 0
		if (i == 0 || i == frames-1) && math.Abs(float64(left)) > 0.001 {
			t.Fatal("sound should fade to silence at both ends")
		}
	}
	if !nonzero {
		t.Fatal("eating sound should not be silent")
	}
}
