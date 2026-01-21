package types

import (
	"fmt"
	"math"
)

// Vec2 represents a 2-dimensional vector with float64 components.
type Vec2 struct {
	X, Y float64
}

func (v Vec2) String() string {
	return fmt.Sprintf("Vec2{X: %f, Y: %f}", v.X, v.Y)
}

// NewVector2D constructs a new Vector2D with the given x and y components.
func NewVec2(x, y float64) Vec2 {
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
func (v Vec2) Mul(scalar float64) Vec2 {
	return Vec2{v.X * scalar, v.Y * scalar}
}

// LengthSq returns the squared magnitude (length) of the vector (v.X*v.X + v.Y*v.Y).
// This is computationally cheaper than Len() as it avoids the square root calculation.
// Useful for comparing vector lengths.
func (v Vec2) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Len returns the magnitude (length) of the vector (sqrt(v.X*v.X + v.Y*v.Y)).
// If you only need to compare lengths, use LengthSq() for better performance.
func (v Vec2) Len() float64 {
	// Directly use LengthSq() for clarity and potential minor optimization reuse.
	return math.Sqrt(v.LengthSq())
}

// Dot returns the dot product (scalar product) of vectors v and other.
// The dot product is calculated as v.X*other.X + v.Y*other.Y.
// Useful for calculating angles between vectors or projecting one vector onto another.
func (v Vec2) Dot(other Vec2) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) Rotate(angleRad float64) Vec2 {
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)
	newX := v.X*cosA - v.Y*sinA
	newY := v.X*sinA + v.Y*cosA
	return Vec2{X: newX, Y: newY}
}

// Normalize returns a unit vector (a vector with length 1) pointing in the same direction as v.
// If the original vector v has a length of 0, it returns a zero vector {0, 0}.
// It does not modify the original vector v.
func (v Vec2) Normalize() Vec2 {
	lenSq := v.LengthSq()
	if lenSq == 0 {
		// Cannot normalize a zero-length vector, return zero vector.
		return Vec2{0, 0}
	}
	len := math.Sqrt(lenSq) // Calculate length only if non-zero
	return Vec2{v.X / len, v.Y / len}
}
