package mesh

import "fmt"
import "github.com/ungerik/go3d/float64/vec3"
import "testing"

func TestTriangleContainsPoint(t *testing.T) {
	tri := Triangle{
		vec3.T{-1, -1, 0},
		vec3.T{1, -1, 0},
		vec3.T{0, 1, 0},
	}

	if !tri.ContainsPoint(vec3.T{0, 0, 0}) {
		t.Fatal("Point determined to be not in triangle when it actually is")
	}

	if tri.ContainsPoint(vec3.T{0, 3, 0}) {
		t.Fatal("Point determined to be in triangle when it actually isn't")
	}
}

func TestPlaneLineIntersection(t *testing.T) {
	plane := Plane{Normal: vec3.T{1, 1, 1}}

	validCases := map[Line]interface{}{
		Line{vec3.T{-1, -1, -1}, vec3.T{2, 2, 2}}:    vec3.T{0, 0, 0},
		Line{vec3.T{2, 0, -2}, vec3.T{-2, 0, 2}}:     nil,
		Line{vec3.T{2, 2, 2}, vec3.T{5, 5, 5}}:       nil,
		Line{vec3.T{-2, -2, -2}, vec3.T{-5, -5, -5}}: nil,
	}

	for line, intersection := range validCases {
		pt := plane.IntersectLine(line)
		if pt == nil {
			if intersection != nil {
				t.Fatal("Line-Plane intersection failed")
			}
			continue
		} else if *pt != intersection {
			t.Fatal("Line-Plane intersection failed")
		}
	}
}

func TestTriangleIntersection(t *testing.T) {
	a := Triangle{
		vec3.T{-1, -1, 0},
		vec3.T{1, -1, 0},
		vec3.T{0, 1, 0},
	}

	//b := Triangle{
	//vec3.T{2, 0, 5},
	//vec3.T{2, 0, -5},
	//vec3.T{-10, 0, 2},
	//}

	//fmt.Println(a.IntersectTriangle(b))
	//fmt.Println(b.IntersectTriangle(a))

	c := Triangle{
		vec3.T{0.6, 0, -5},
		vec3.T{0.6, 0, 5},
		vec3.T{-0.4, 1, 0},
	}

	//fmt.Println(a.IntersectTriangle(c))
	fmt.Println(c.IntersectTriangle(a))
}

func TestCoplanarTriangles(t *testing.T) {
	a := Triangle{
		vec3.T{-1, -1, 0},
		vec3.T{1, -1, 0},
		vec3.T{0, 1, 0},
	}

	b := Triangle{
		vec3.T{-1, 0.5, 0},
		vec3.T{1, 0.5, 0},
		vec3.T{0, -1, 0},
	}

	fmt.Println(a.IntersectTriangle(b))
}

func TestIntersectionRing(t *testing.T) {
	ibuf := IndexBuffer{
		Vertices: []vec3.T{
			{-5, 0, -5}, // 0
			{5, 0, -5},  // 1
			{0, 0, 5},   // 2

			{-2, 3, -2}, // 3
			{2, 3, -2},  // 4
			{0, 3, 2},   // 5
			{0, -2, 0},  // 6
		},
		Faces: []Face{
			// Base triangle
			{0, 1, 2},

			// Top triangles
			{3, 6, 4},
			{5, 4, 6},
			{5, 6, 3},
		},
	}

	abuf := ArrayBuffer{}
	abuf.ConvertFrom(&ibuf)

	vertStore := make(map[vec3.T]bool)
	for _, tri := range abuf[1:] {
		line, _ := abuf[0].IntersectTriangle(tri)
		fmt.Println(line)
		vertStore[line[0]] = true
		vertStore[line[1]] = true
	}

	fmt.Println()
	fmt.Println(vertStore)
}
