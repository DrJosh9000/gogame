package game

// This file has types for describing the whole game world at a high level,
// including the collection of landscapes and the connections between them.
// The world doesn't describe where the player is. The player object can
// be thought of as having coordinates (x,y,z,landscapeKey).

// world describes a collection of landscapes.
type world struct {
	landscapes map[string]*landscape
}

// landscape includes terrain, decorations, ...
type landscape struct {
	terrain, decorations complexBase
}
