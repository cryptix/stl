package stl

// Tests for the Solid data type.
// Could be more exhaustive.

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-test/deep"
)

func TestTransform(t *testing.T) {
	sOrig := makeTestSolid()
	s := makeTestSolid()
	s.Transform(mgl32.Ident4())
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after identity transformation")
		t.Log("Expected:\n", *sOrig)
		t.Log("Found:\n", *s)
	}
}

func TestScale(t *testing.T) {
	sOrig := makeTestSolid()
	s := makeTestSolid()
	s.Scale(2)
	s.Scale(0.5)
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after successive scaling * 2 * 0.5")
	}
}

func TestStretch(t *testing.T) {
	sOrig := makeTestSolid()
	s := makeTestSolid()
	s.Stretch(mgl32.Vec3{1, 2, 1})
	s.Stretch(mgl32.Vec3{1, 0.5, 1})
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after successive Y scaling * 2 * 0.5")
	}
}

func TestTranslate(t *testing.T) {
	sOrig := makeTestSolid()
	s := makeTestSolid()
	s.Translate([3]float32{1, 2, 4})
	s.Translate([3]float32{-1, -2, -4})
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after translation and inverse translation")
		t.Log("Expected:\n", sOrig)
		t.Log("Found:\n", s)
	}
}

func makeBrokenTestSolid() *Solid {
	return &Solid{
		Name:    "Simple",
		IsAscii: true,
		Triangles: []Triangle{
			// This triangle is the black sheep
			{
				Normal: mgl32.Vec3{0, 0, -1},
				Vertices: [3]mgl32.Vec3{
					// The edge V0 -> V1 is in the wrong direction, V0 and V1 are swapped
					{0, 1, 0},
					{0, 0, 0},
					// This vertex should be {1, 0, 0} - now there is a gap
					{0, 0, 0},
				},
			},
			// For this triangle there is no counter-edge for V0 -> V1
			{
				Normal: mgl32.Vec3{0, -1, 0},
				Vertices: [3]mgl32.Vec3{
					{0, 0, 0},
					{1, 0, 0},
					{0, 0, 1},
				},
			},
			// For this triangle there is no counter-edge for V1 -> V2
			{
				Normal: mgl32.Vec3{0.57735, 0.57735, 0.57735},
				Vertices: [3]mgl32.Vec3{
					{0, 0, 1},
					{1, 0, 0},
					{0, 1, 0},
				},
			},
			// The edge V2 -> V0 has a duplicate in triangle 0
			{
				Normal: mgl32.Vec3{-1, 0, 0},
				Vertices: [3]mgl32.Vec3{
					{0, 0, 0},
					{0, 0, 1},
					{0, 1, 0},
				},
			},
		},
	}
}

func TestValidate(t *testing.T) {
	solid := makeBrokenTestSolid()
	errors := solid.Validate()
	if errors[0] == nil || !errors[0].HasEqualVertices {
		t.Error("Failed to detect HasEqualVertices in triangle 0")
	}
	if errors[1] == nil || errors[1].EdgeErrors[0] == nil ||
		!errors[1].EdgeErrors[0].HasNoCounterEdge() {
		t.Error("Failed to detect missing counter-edge in triangle 1, edge 0")
	}
	if errors[2] == nil || errors[2].EdgeErrors[1] == nil ||
		!errors[2].EdgeErrors[1].HasNoCounterEdge() {
		t.Error("Failed to detect missing counter-edge in triangle 2, edge 1")
	}
	if errors[3] == nil || errors[3].EdgeErrors[2] == nil ||
		!errors[3].EdgeErrors[2].IsUsedInOtherTriangles() ||
		errors[3].EdgeErrors[2].SameEdgeTriangles[0] != 0 {
		t.Error("Failed to detect edge duplicate of triangle 3, edge 2 in triangle 0")
	}
}

func TestMeasure(t *testing.T) {
	testSolid := makeTestSolid()
	measure := testSolid.Measure()
	if !measure.Len.ApproxEqual(mgl32.Vec3{1, 1, 1}) {
		t.Errorf("Expected Len: [1 1 1], found: %v", measure.Len)
	}
}

func TestRotate(t *testing.T) {
	sOrig := makeTestSolid()
	s := makeTestSolid()
	s.Rotate(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 1}, 0)
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after rotation around z-axis with 0 angle")
		t.Error(deep.Equal(sOrig, s))
	}

	s.Rotate(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 1}, HalfPi)
	s.Rotate(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 1}, -HalfPi)
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after two rotations around z-axis cancelling each other out")
		t.Error(deep.Equal(sOrig, s))
	}

	s.Rotate(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1}, HalfPi)
	s.Rotate(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1}, -HalfPi)
	if !sOrig.sameOrderAlmostEqual(s) {
		t.Error("Not equal after two rotations cancelling each other out")
		t.Log("Expected:\n", sOrig)
		t.Log("Found:\n", s)
	}
}

func BenchmarkTransform(b *testing.B) {
	b.StopTimer()
	solid, err := ReadFile(testFilenameComplexASCII)
	if err != nil {
		b.Fatal(err)
	}
	//pos := mgl32.Vec3{30, 10, 10}
	axis := mgl32.Vec3{1, 1, 1}
	angle := float32(HalfPi / 4)
	rotationMatrix := mgl32.HomogRotate3D(angle, axis)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		solid.Transform(rotationMatrix)
	}
}
