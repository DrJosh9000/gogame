package game

import (
	"time"
)

type Collider struct {
	Base
	x, y int
	fx, fy, dx, dy, ddx, ddy float64
	lastUpdate time.Duration
	
	Update chan time.Duration
}