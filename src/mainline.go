package main

import ( "fmt" 
	"math" )

type Vec3 struct {
	x float64
	y float64
	z float64
}

func (v *Vec3) add(other *Vec3) Vec3 {
	return Vec3{v.x+other.x,v.y+other.y,v.z+other.z}
}
func (v *Vec3) multConst(c float64) Vec3 {
	return Vec3{v.x*c,v.y*c,v.z*c}
}
func (v *Vec3) mult(other *Vec3) Vec3 {
	return Vec3{other.x*v.x,other.y*v.y,other.z*v.z}
}
func (v *Vec3) length2() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z;
}
func (v *Vec3) minus(other * Vec3) Vec3 {
	return Vec3{v.x-other.x,v.y-other.y,v.z-other.z}
}

// in place
func (v *Vec3) normalize() {
	var nor2 float64 = v.length2();
	if (nor2>0) {
		var invNor float64 = 1 / math.Sqrt(nor2);
		v.x = v.x * invNor
		v.y = v.y * invNor
		v.z = v.z * invNor
	}
}

func (v *Vec3) dot(other * Vec3) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z;
}

type Sphere struct {
	center Vec3
	radius,radius2 float64
	surfaceColor, emissionColor Vec3
	transparency, reflection float64
}

func (s *Sphere) intersect(rayorig * Vec3, 
	raydir * Vec3, t0 *float64, t1 *float64) bool {
	var l Vec3 = s.center.minus(rayorig);
	var tca float64 = l.dot(raydir);
	if (tca<0) {
		return false;
	}
	var d2 float64 = l.dot(&l) - tca * tca;
	if (d2 > s.radius2) {
		return false;
	}
	var thc float64 = math.Sqrt(s.radius2 - d2);
	// TODO: how do you do this in golang?
	&t0 = tca - thc;
	&t1 = tca + thc;
	return true;
}

func main() {

	var v Vec3
	v.x = 1
	v.y = 2
	v.z = 3
	
	fmt.Println(v)

	v.normalize()

	fmt.Println(v)
}
