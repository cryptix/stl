package stl

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// This file defines the Triangle data type, the building block for Solid

// Triangle represents single triangles used in Solid.Triangles. The vertices
// have to be ordered counter-clockwise when looking at their outside surface.
// The vector Normal is orthogonal to the triangle, pointing outside, and
// has length 1. This is redundant but included in the STL format in order to
// avoid recalculation.
type Triangle struct {
	// Normal vector of triangle, should be normalized...
	Normal mgl64.Vec3

	// Vertices of triangle in right hand order.
	// I.e. from the front the triangle's vertices are ordered counterclockwise
	// and the normal vector is orthogonal to the front pointing outside.
	Vertices [3]mgl64.Vec3

	// 16 bits of attributes. Not available in ASCII format. Could be used
	// for color selection, texture selection, refraction etc. Some tools ignore
	// this field completely, always writing 0 on export.
	Attributes uint16
}

// Calculate the normal vector using the right hand rule
func (t *Triangle) calculateNormal() mgl64.Vec3 {
	// The normal is calculated by normalizing the result of
	// (V0-V2) x (V1-V2)
	return t.Vertices[0].Sub(t.Vertices[2]).
		Cross(t.Vertices[1].Sub(t.Vertices[2])).
		Normalize()
}

func (t *Triangle) Area() float64 {
	a := t.Vertices[0].Sub(t.Vertices[1]).Len()
	b := t.Vertices[1].Sub(t.Vertices[2]).Len()
	c := t.Vertices[2].Sub(t.Vertices[0]).Len()
	s := (a + b + c) / 2
	a = s - a
	b = s - b
	c = s - c
	return math.Sqrt(s * a * b * c)
}

func (t *Triangle) Quadric() mgl64.Mat4 {
	n := t.calculateNormal()

	x, y, z := t.Vertices[0].X(), t.Vertices[0].Y(), t.Vertices[0].Z()
	a, b, c := n.X(), n.Y(), n.Z()
	d := -a*x - b*y - c*z
	return mgl64.Mat4{
		a * a, a * b, a * c, a * d,
		a * b, b * b, b * c, b * d,
		a * c, b * c, c * c, c * d,
		a * d, b * d, c * d, d * d,
	}
}

// Recalculate the redundant normal vector using the right hand rule
func (t *Triangle) RecalculateNormal() {
	t.Normal = t.calculateNormal()
}

// Applies a 4x4 transformation matrix to every vertex
// and recalculates the normal
func (t *Triangle) transform(transformationMatrix mgl64.Mat4) {
	t.transformNR(transformationMatrix)
	t.RecalculateNormal()
}

// Applies a 4x4 transformation matrix to every vertex
// without recalculating the normal afterwards
func (t *Triangle) transformNR(transformationMatrix mgl64.Mat4) {
	m3 := transformationMatrix.Mat3()
	for i := 0; i < 3; i++ {
		t.Vertices[i] = m3.Mul3x1(t.Vertices[i])
	}
	/*
		t.Vertices[x] = transformationMatrix.MultVec3(t.Vertices[x])
	*/
}

// Returns true if at least two vertices are exactly equal, meaning
// this is a line, or even a dot.
func (t *Triangle) hasEqualVertices() bool {
	return t.Vertices[0].ApproxEqual(t.Vertices[1]) ||
		t.Vertices[0].ApproxEqual(t.Vertices[2]) ||
		t.Vertices[1].ApproxEqual(t.Vertices[2])
}

// Checks if normal matches vertices using right hand rule, with
// numerical tolerance for angle between them given by tol in radians.
func (t *Triangle) checkNormal(tol float64) bool {
	calculatedNormal := t.calculateNormal()
	return Angle(t.Normal, calculatedNormal) < tol
}

func Angle(vec, o mgl64.Vec3) float64 {
	lenProd := vec.Len() * o.Len()
	if lenProd == 0 {
		return 0
	}
	cosAngle := vec.Dot(o) / lenProd
	// Numerical correction
	if cosAngle < -1 {
		cosAngle = -1
	} else if cosAngle > 1 {
		cosAngle = 1
	}

	return math.Acos(cosAngle)
}
