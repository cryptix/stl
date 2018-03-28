package stl

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func BenchmarkMultMat4(b *testing.B) {
	r := mgl32.Mat4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
	}
	s := mgl32.Mat4{
		9, 2, 1, -6,
		8, 3, 0, -5,
		7, 4, -1, -4,
		6, 5, -2, -3,
	}
	for i := 0; i < b.N; i++ {
		_ = r.Mul4(s)
	}
}

func TestMultMat4(t *testing.T) {
	m := mgl32.Mat4{
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
		1, 2, 3, 4,
	}
	var r mgl32.Mat4
	r = m.Mul4(mgl32.Ident4())
	if m != r {
		t.Errorf("Result: %v, Expected: %v", r, m)
	}
}

func BenchmarkMultVec3(b *testing.B) {
	m := mgl32.Mat4{
		9, 2, 1, -6,
		8, 3, 0, -5,
		7, 4, -1, -4,
		6, 5, -2, -3,
	}
	v := mgl32.Vec3{-1000, 234, 1000}
	for i := 0; i < b.N; i++ {
		_ = m.Mul4x1(v.Vec4(0))
	}
}

func TestMultVec3(t *testing.T) {
	m := mgl32.Mat4{
		1, 0, 0, 1000,
		0, 2, 0, 500,
		0, 0, 1, 250,
		0, 0, 0, 1,
	}
	v := mgl32.Vec4{1, 1, 1, 1}
	r := m.Transpose().Mul4x1(v)
	expected := mgl32.Vec4{1001, 502, 251, 1}
	if r != expected {
		t.Errorf("MultVec3 result: %v, expected: %v", r, expected)
	}
}
