package types

import (
	"fmt"
)

// Vec2 represents a 2-dimensional vector with float64 components.
type Vec2 struct {
	X, Y int
}

func (v Vec2) String() string {
	return fmt.Sprintf("Vec2{X: %f, Y: %f}", v.X, v.Y)
}

// NewVector2D constructs a new Vector2D with the given x and y components.
func NewVec2(x, y int) Vec2 {
	return Vec2{x, y}
}

// Add returns the vector sum of v and other (v + other).
// It does not modify the original vector v.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

// Sub returns the vector difference of v and other (v - other).
// It does not modify the original vector v.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

// Mul returns the vector v scaled by the given scalar factor (v * scalar).
// It does not modify the original vector v.
func (v Vec2) Mul(scalar int) Vec2 {
	return Vec2{v.X * scalar, v.Y * scalar}
}

// LengthSq returns the squared magnitude (length) of the vector (v.X*v.X + v.Y*v.Y).
// This is computationally cheaper than Len() as it avoids the square root calculation.
// Useful for comparing vector lengths.
func (v Vec2) LengthSq() int {
	return v.X*v.X + v.Y*v.Y
}

// Dot returns the dot product (scalar product) of vectors v and other.
// The dot product is calculated as v.X*other.X + v.Y*other.Y.
// Useful for calculating angles between vectors or projecting one vector onto another.
func (v Vec2) Dot(other Vec2) int {
	return v.X*other.X + v.Y*other.Y
}
