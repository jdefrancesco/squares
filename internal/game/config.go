package game

const (
	// ScreenWidth is the fixed logical width of the playfield in pixels.
	ScreenWidth = 640
	// ScreenHeight is the fixed logical height of the playfield in pixels.
	ScreenHeight = 480

	initialPlayerSize  = 22.0
	maxPlayerSize      = 64.0
	playerHitScale     = 0.90
	playerRotationRate = 3.0

	growthScale = 0.05
	growthFlat  = 0.8

	maxSpawnsWithoutEdible = 2
	maxActiveDangers       = 6
	maxEntitySpeed         = 220.0
	spawnMargin            = 100.0
	spawnReactionTime      = 0.85
	escapeHorizon          = 1.5
	escapeMoveTime         = 0.35

	eatPulseDuration   = 0.18
	scorePopupDuration = 0.45
)
