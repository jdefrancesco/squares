package game

// Entity is an axis-aligned square with a center, side length, and fixed velocity.
// Its size relative to the player determines whether it is edible or dangerous.
type Entity struct {
	x    float64
	y    float64
	size float64
	vx   float64
	vy   float64
}

// inputState captures controls independently of the window for deterministic steps.
type inputState struct {
	x, y                                                    float64
	mouseDown, pauseDown, escapeDown, restartDown, quitDown bool
}
