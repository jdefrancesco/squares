package game

import "math"

// spawnEntityWithDifficulty appends a food or danger square outside a random edge.
// Difficulty is in [0, 1]. Crowd limits, reaction time, and escape checks replace
// rejected threats with food, while the opening seconds always offer food.
func (g *Game) spawnEntityWithDifficulty(d float64) {
	d = clamp(d, 0, 1)
	dangers := 0
	for _, e := range g.ents {
		if e.size >= g.player.size {
			dangers++
		}
	}
	limit := min(maxActiveDangers, 2+int(4*d))
	edible := g.elapsed < 2 || g.spawnsSinceEdible >= maxSpawnsWithoutEdible || dangers >= limit || g.random() < 0.70-0.32*d
	size := math.Max(10, g.player.size*(0.45+g.random()*0.30))
	if !edible {
		size = math.Min(112, g.player.size*(1.15+g.random()*0.55)+6)
	}
	e := g.makeSpawn(size, d)
	if !edible && !g.hasEscapeRoute(e) {
		edible = true
		e = g.makeSpawn(math.Max(10, g.player.size*0.60), d)
	}
	if edible {
		g.spawnsSinceEdible = 0
	} else {
		g.spawnsSinceEdible++
	}
	g.ents = append(g.ents, e)
}

// makeSpawn chooses an edge and a fixed trajectory toward the current player,
// slowing incoming squares when necessary to preserve minimum reaction time.
func (g *Game) makeSpawn(size, d float64) Entity {
	margin := spawnMargin + size/2
	e := Entity{size: size}
	switch int(g.random() * 4) {
	case 0:
		e.x, e.y = g.random()*ScreenWidth, -margin
	case 1:
		e.x, e.y = g.random()*ScreenWidth, ScreenHeight+margin
	case 2:
		e.x, e.y = -margin, g.random()*ScreenHeight
	default:
		e.x, e.y = ScreenWidth+margin, g.random()*ScreenHeight
	}
	dx, dy := g.player.x-e.x, g.player.y-e.y
	distance := math.Hypot(dx, dy)
	speed := math.Min(maxEntitySpeed, 95+110*d+g.random()*15)
	e.vx, e.vy = dx/distance*speed, dy/distance*speed
	playerHit := g.player
	playerHit.size *= playerHitScale
	if at, hit := squareCollisionTime(e, playerHit, spawnReactionTime); hit && at > 0 {
		scale := at / (spawnReactionTime + 0.01)
		e.vx, e.vy = e.vx*scale, e.vy*scale
	}
	return e
}

// hasEscapeRoute checks eight nearby destinations plus staying still. A route
// must remain clear while moving there and while waiting through the 1.5-second
// lookahead. Food is ignored; this is a local spawn safeguard, not a guarantee
// against later movement into danger.
func (g *Game) hasEscapeRoute(candidate Entity) bool {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			distance := math.Max(90, g.player.size*1.5)
			margin := g.player.size / math.Sqrt2
			x := clamp(g.player.x+float64(dx)*distance, margin, ScreenWidth-margin)
			y := clamp(g.player.y+float64(dy)*distance, margin, ScreenHeight-margin)
			if g.escapeRouteOpen(x, y, candidate) {
				return true
			}
		}
	}
	return false
}

// escapeRouteOpen checks the moving and stationary portions of one escape route
// against the proposed square and every currently dangerous square.
func (g *Game) escapeRouteOpen(x, y float64, candidate Entity) bool {
	player := g.player
	player.size *= playerHitScale
	player.vx, player.vy = (x-player.x)/escapeMoveTime, (y-player.y)/escapeMoveTime
	for i := -1; i < len(g.ents); i++ {
		e := candidate
		if i >= 0 {
			e = g.ents[i]
		}
		if e.size < g.player.size {
			continue
		}
		if sweptSquareCollision(player, e, escapeMoveTime) {
			return false
		}
		e.x, e.y = e.x+e.vx*escapeMoveTime, e.y+e.vy*escapeMoveTime
		atRest := Entity{x: x, y: y, size: player.size}
		if sweptSquareCollision(atRest, e, escapeHorizon-escapeMoveTime) {
			return false
		}
	}
	return true
}

// playfieldEntry predicts the point where an off-screen square's center enters
// the playfield, allowing a warning marker without revealing its whole trajectory.
func playfieldEntry(e Entity) (float64, float64, bool) {
	if e.x >= 0 && e.x <= ScreenWidth && e.y >= 0 && e.y <= ScreenHeight {
		return 0, 0, false
	}
	enter, leave := 0.0, math.Inf(1)
	for _, axis := range [][3]float64{{e.x, e.vx, ScreenWidth}, {e.y, e.vy, ScreenHeight}} {
		pos, velocity, limit := axis[0], axis[1], axis[2]
		if velocity == 0 {
			if pos < 0 || pos > limit {
				return 0, 0, false
			}
			continue
		}
		first, last := -pos/velocity, (limit-pos)/velocity
		if first > last {
			first, last = last, first
		}
		enter, leave = math.Max(enter, first), math.Min(leave, last)
		if enter > leave {
			return 0, 0, false
		}
	}
	return e.x + e.vx*enter, e.y + e.vy*enter, true
}
