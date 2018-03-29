package stl

// This file contains tests for the Vec3 data type

import (
	"testing"

	"github.com/go-gl/mathgl/mgl64"
)

func TestVec3Angle(t *testing.T) {
	v := mgl64.Vec3{1, 0, 0}
	tol := 0.00005
	testV := []mgl64.Vec3{
		mgl64.Vec3{0, 1, 0},
		mgl64.Vec3{0, -1, 0},
		mgl64.Vec3{-1, 0, 0},
		mgl64.Vec3{-1, 1, 0},
		mgl64.Vec3{-1, -1, 0}}
	expected := []float64{
		HalfPi,
		HalfPi,
		Pi,
		HalfPi + QuarterPi,
		HalfPi + QuarterPi}
	for i, tv := range testV {
		r := Angle(v, tv)
		if !almostEqual64(expected[i], r, tol) {
			t.Errorf("angle(%v, %v) = %g Pi, expected %g Pi", v, tv, r/Pi, expected[i]/Pi)
		}
	}
}
