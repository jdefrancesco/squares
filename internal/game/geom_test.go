package game

import "testing"

func TestClamp(t *testing.T) {
	cases := []struct {
		name      string
		v, lo, hi float64
		want      float64
	}{
		{"below", -1, 0, 10, 0},
		{"inside", 5, 0, 10, 5},
		{"above", 11, 0, 10, 10},
	}

	for _, tc := range cases {
		if got := clamp(tc.v, tc.lo, tc.hi); got != tc.want {
			t.Fatalf("%s: clamp(%v,%v,%v)=%v want %v", tc.name, tc.v, tc.lo, tc.hi, got, tc.want)
		}
	}
}

func TestClampInt(t *testing.T) {
	cases := []struct {
		name      string
		v, lo, hi int
		want      int
	}{
		{"below", -1, 0, 10, 0},
		{"inside", 5, 0, 10, 5},
		{"above", 11, 0, 10, 10},
	}

	for _, tc := range cases {
		if got := clampInt(tc.v, tc.lo, tc.hi); got != tc.want {
			t.Fatalf("%s: clampInt(%v,%v,%v)=%v want %v", tc.name, tc.v, tc.lo, tc.hi, got, tc.want)
		}
	}
}

func TestSquareIntersectsSquare(t *testing.T) {
	a := Entity{x: 0, y: 0, size: 10}

	b := Entity{x: 4, y: 0, size: 10}
	if !squareIntersectsSquare(a, b) {
		t.Fatalf("expected squares to intersect")
	}

	c := Entity{x: 10, y: 0, size: 10} // touching
	if !squareIntersectsSquare(a, c) {
		t.Fatalf("expected squares touching edges to intersect")
	}

	d := Entity{x: 10.001, y: 0, size: 10}
	if squareIntersectsSquare(a, d) {
		t.Fatalf("expected squares not to intersect")
	}
}

func TestCircleIntersectsSquare(t *testing.T) {
	sq := Entity{x: 0, y: 0, size: 2} // bounds [-1,1]

	c1 := Entity{x: 0, y: 0, size: 2} // r=1
	if !circleIntersectsSquare(c1, sq) {
		t.Fatalf("expected circle center inside square to intersect")
	}

	c2 := Entity{x: 10, y: 10, size: 2}
	if circleIntersectsSquare(c2, sq) {
		t.Fatalf("expected far circle not to intersect")
	}

	// Tangent to right edge at (1,0): center at (2,0), r=1.
	c3 := Entity{x: 2, y: 0, size: 2}
	if !circleIntersectsSquare(c3, sq) {
		t.Fatalf("expected tangent circle to count as intersecting")
	}
}
