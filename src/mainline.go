package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

const SphereCount int = 6

var background_count int

type Vec3 struct {
	x float64
	y float64
	z float64
}

func (v Vec3) add(other Vec3) Vec3 {
	return Vec3{v.x + other.x, v.y + other.y, v.z + other.z}
}
func (v Vec3) multConst(c float64) Vec3 {
	return Vec3{v.x * c, v.y * c, v.z * c}
}
func (v Vec3) mult(other Vec3) Vec3 {
	return Vec3{other.x * v.x, other.y * v.y, other.z * v.z}
}
func (v Vec3) length2() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}
func (v Vec3) minus(other Vec3) Vec3 {
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

func (v Vec3) dot(other Vec3) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z
}

type Sphere struct {
	center                      Vec3
	radius, radius2             float64
	surfaceColor, emissionColor Vec3
	transparency, reflection    float64
}

func (s *Sphere) intersect(rayorig Vec3,
	raydir Vec3, t0 *float64, t1 *float64) bool {
	// the next two lines are an experiment
	// removed types to see what they look like
	var l = s.center.minus(rayorig) // took out type here
	var tca = l.dot(raydir)         // took out type here

	if tca < 0 {
		return false
	}
	var d2 float64 = l.dot(l) - tca*tca
	if d2 > s.radius2 {
		return false
	}
	var thc float64 = math.Sqrt(s.radius2 - d2)
	*t0 = tca - thc
	*t1 = tca + thc
	return true
}

var Inf float64 = 1e8
var MaxRayDepth int = 5

func mix(a float64, b float64, mix float64) float64 {
	return b*mix + a*(1-mix)
}

func trace(tracer bool, rayorig Vec3, raydir Vec3, spheres [6]Sphere,
	depth int) Vec3 {
	var tnear float64 = Inf
	var sphere *Sphere = nil
	for i := 0; i < SphereCount; i++ {
		var t0 float64 = Inf
		var t1 float64 = Inf
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

	if tracer {
		fmt.Println("tnear sphere", tnear, sphere)
	}
	if sphere == nil {
		background_count++
		return Vec3{2, 2, 2}
	}

	var surfaceColor Vec3 = Vec3{0, 0, 0}
	var phit Vec3 = rayorig.add(raydir).multConst(tnear)
	var nhit Vec3 = phit.minus(sphere.center)
	nhit.normalize()
	var bias float64 = 1e-4
	var inside bool = false
	if raydir.dot(nhit) > 0 {
		nhit = nhit.multConst(-1)
		inside = true
	}

	if (sphere.transparency > 0 ||
		sphere.reflection > 0) && depth < MaxRayDepth {

		var facingratio float64 = -raydir.dot(nhit)
		var fresneleffect float64 = mix(math.Pow(1-facingratio, 3), 1, 0.1)

		var refldir Vec3 = raydir.minus(nhit).multConst(2).multConst(raydir.dot(nhit))
		refldir.normalize()

		var reflection Vec3 = trace(tracer, phit.add(nhit).multConst(bias), refldir, spheres, depth+1)
		var refraction Vec3 = Vec3{0, 0, 0}
		if sphere.transparency > 0 {
			var ior float64 = 1.1
			var eta float64
			if inside {
				eta = ior
			} else {
				eta = 1 / ior
			}
			var cosi float64 = -nhit.dot(raydir)
			var k float64 = 1 - eta*eta*(1-cosi*cosi)
			var refrdir Vec3 = raydir.multConst(eta).add(nhit).multConst(eta*cosi - math.Sqrt(k))
			refldir.normalize()
			refraction = trace(tracer, phit.minus(nhit.multConst(bias)), refrdir, spheres, depth+1)
		}

		_tmp0 := reflection.multConst(fresneleffect)
		_tmp1 := refraction.multConst((1 - fresneleffect) * sphere.transparency)
		_tmp2 := sphere.surfaceColor.mult(_tmp1)
		surfaceColor = _tmp0.add(_tmp2)

	} else {

		// diffuse object
		for i := 0; i < 6; i++ {
			if spheres[i].emissionColor.x > 0 {
				var transmission Vec3 = Vec3{1, 1, 1}
				var lightDirection Vec3 = spheres[i].center.minus(phit)
				lightDirection.normalize()

				for j := 0; j < 6; j++ {
					if i != j {
						var t0, t1 float64
						if spheres[i].intersect(phit.add(nhit).multConst(bias),
							lightDirection, &t0, &t1) {
							transmission = Vec3{0, 0, 0}
							break
						}
					}
				}

				_tmp0 := surfaceColor.add(sphere.surfaceColor)
				_tmp0a := _tmp0.mult(transmission)
				_tmp1 := spheres[i].emissionColor.multConst(math.Max(0, nhit.dot(lightDirection)))
				surfaceColor = _tmp0a.mult(_tmp1)

			}
		}
	}

	return surfaceColor.add(sphere.emissionColor)
}

func writePPMHeader(outf *os.File) {
	fheader := []byte("P6\n")
	binary.Write(outf, binary.LittleEndian, fheader)
	fheader2 := []byte("640 480\n")
	binary.Write(outf, binary.LittleEndian, fheader2)
	fheader3 := []byte("255\n")
	binary.Write(outf, binary.LittleEndian, fheader3)
}

func render(spheres [SphereCount]Sphere) {
	var width int = 640
	var height int = 480

	// image is of type []Vec3
	var image = make([]Vec3, width*height)
	var counter int
	// var pixel Vec3 = image[counter]
	var invWidth float64 = 1 / float64(width)
	var invHeight = 1 / float64(height)
	var fov float64 = 30.0
	var aspectratio float64 = float64(width) / float64(height)
	var angle = math.Tan(math.Pi * 0.5 * fov / 180.0)
	//var onetime bool = true

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var xx float64 = (2*((float64(x)+0.5)*invWidth) - 1) *
				angle * aspectratio
			var yy float64 = (1 - 2*
				((float64(y)+0.5)*invHeight)) * angle
			var raydir Vec3 = Vec3{xx, yy, -1}
			raydir.normalize()
			//image[counter] = trace(Vec3{0, 0, 0}, raydir, spheres, 0)
			//image[counter] = trace(false, Vec3{0, 0, 0}, raydir, spheres, 0)
			if (x == 309 && y == 57) || (y == 309 && x == 57) {
				//if onetime && image[counter].x != 2 && image[counter].y != 2 && image[counter].z != 2 {
				image[counter] = trace(true, Vec3{0, 0, 0}, raydir, spheres, 0)
				fmt.Println("pixel( ", x, y, counter, image[counter])

				//onetime = false
				//}
			} else {
				image[counter] = trace(false, Vec3{0, 0, 0}, raydir, spheres, 0)
			}
			counter++
		}
	}

	// stupid go crap here... just screw the type
	outf, _ := os.Create("out.ppm")

	// oh but wait you really the type to write this function...
	writePPMHeader(outf)

	counter = 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			binary.Write(outf, binary.LittleEndian,
				uint8(math.Min(1, image[counter].x)*255))
			binary.Write(outf, binary.LittleEndian,
				uint8(math.Min(1, image[counter].y)*255))
			binary.Write(outf, binary.LittleEndian,
				uint8(math.Min(1, image[counter].z)*255))

			counter = counter + 1
		}

	}
	outf.Close()
}

func make_sphere(center Vec3, radius float64, surface_color Vec3,
	reflectivity float64, transparency float64, emission_color Vec3) Sphere {
	return Sphere{center, radius, 2.0 * radius, surface_color, emission_color, reflectivity, transparency}
}

func main() {
	var spheres [SphereCount]Sphere
	spheres[0] = make_sphere(Vec3{0.0, -10004, -20},
		10000, Vec3{0.20, 0.20, 0.20},
		0.0, 0.0, Vec3{0, 0, 0})
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

	fmt.Println("background_count", background_count)
}
