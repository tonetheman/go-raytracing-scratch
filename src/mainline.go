package main

import (
	"fmt"
	"math"
)

type Vec3 struct {
	x float64
	y float64
	z float64
}

func (v *Vec3) add(other *Vec3) Vec3 {
	return Vec3{v.x + other.x, v.y + other.y, v.z + other.z}
}
func (v *Vec3) multConst(c float64) Vec3 {
	return Vec3{v.x * c, v.y * c, v.z * c}
}
func (v *Vec3) mult(other *Vec3) Vec3 {
	return Vec3{other.x * v.x, other.y * v.y, other.z * v.z}
}
func (v *Vec3) length2() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}
func (v *Vec3) minus(other *Vec3) Vec3 {
	return Vec3{v.x - other.x, v.y - other.y, v.z - other.z}
}

// in place
func (v *Vec3) normalize() {
	var nor2 float64 = v.length2()
	if nor2 > 0 {
		var invNor float64 = 1 / math.Sqrt(nor2)
		v.x = v.x * invNor
		v.y = v.y * invNor
		v.z = v.z * invNor
	}
}

func (v *Vec3) dot(other *Vec3) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z
}

type Sphere struct {
	center                      Vec3
	radius, radius2             float64
	surfaceColor, emissionColor Vec3
	transparency, reflection    float64
}

func (s *Sphere) init(c *Vec3, r float64, sc *Vec3, refl float64,
	transp float64, ec *Vec3) {
	s.center = *c
	s.radius = r
	s.radius2 = r * r
	s.surfaceColor = *sc
	s.emissionColor = *ec
	s.transparency = transp
	s.reflection = refl
}

func (s *Sphere) intersect(rayorig *Vec3,
	raydir *Vec3, t0 *float64, t1 *float64) bool {
	var l Vec3 = s.center.minus(rayorig)
	var tca float64 = l.dot(raydir)
	if tca < 0 {
		return false
	}
	var d2 float64 = l.dot(&l) - tca*tca
	if d2 > s.radius2 {
		return false
	}
	var thc float64 = math.Sqrt(s.radius2 - d2)
	*t0 = tca - thc
	*t1 = tca + thc
	return true
}

var INF float64 = 1e8

func trace(rayorig *Vec3, raydir *Vec3, spheres []Sphere,
	depth int) Vec3 {
	var tnear float64 = INF
	var sphere *Sphere = nil
	for i := 0; i < 6; i++ {
		var t0 float64 = INF
		var t1 float64 = INF
		if spheres[i].intersect(rayorig, raydir, &t0, &t1) {
			if t0 < 0 {
				t0 = t1
			}
			if t0 < tnear {
				tnear = t0
				sphere = &spheres[i]
			}
		}
	}

	if sphere == nil {
		return Vec3{2, 2, 2}
	}

	// TODO: not done here:

}

func render(spheres [6]Sphere) {
	var width int = 640
	var height int = 480

	var image []Vec3 = make([]Vec3, width*height)
	fmt.Println(image)
	var counter int = 0
	// var pixel Vec3 = image[counter]
	var invWidth float64 = 1 / float64(width)
	var invHeight = 1 / float64(height)
	var fov float64 = 30.0
	var aspectratio float64 = float64(width) / float64(height)
	var angle = math.Tan(math.Pi * 0.5 * fov / 180.0)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var xx float64 = (2*((float64(x)+0.5)*invWidth) - 1) *
				angle * aspectratio
			var yy float64 = (1 - 2*
				((float64(y)+0.5)*invHeight)) * angle
			var raydir Vec3 = Vec3{xx, yy, -1}
			raydir.normalize()

			// TODO: not done here

			counter++
		}
	}

	fmt.Println("P6")
	fmt.Println(width, height)
	fmt.Println(255)
}

func main() {
	var spheres [6]Sphere
	spheres[0] = Sphere{Vec3{0.0, -10004, -20}, 10000, 2 * 10000,
		Vec3{0.20, 0.20, 0.20}, Vec3{0, 0, 0}, 0, 0.0}
	spheres[1] = Sphere{Vec3{0.0, 0, -20}, 4, 2 * 4,
		Vec3{1.0, 0.32, 0.36}, Vec3{0, 0, 0}, 1, 0.5}
	spheres[2] = Sphere{Vec3{5, -1, -15}, 2, 2 * 2,
		Vec3{0.9, 0.76, 0.46}, Vec3{0, 0, 0}, 1, 0.0}
	spheres[3] = Sphere{Vec3{5, 0, -25}, 3, 2 * 3,
		Vec3{0.65, 0.77, 0.97}, Vec3{0, 0, 0}, 1, 0.0}
	spheres[4] = Sphere{Vec3{-5.5, 0, -15}, 3, 2 * 3,
		Vec3{0.90, 0.90, 0.90}, Vec3{0, 0, 0}, 1, 0.0}

	// light
	spheres[5] = Sphere{Vec3{0.0, 20, -35}, 3, 2 * 3,
		Vec3{0.0, 0.0, 0.0}, Vec3{3, 3, 3}, 0, 0.0}

	render(spheres)

}
