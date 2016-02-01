package vec3

import "math"

type Vec3 struct {
	x float64
	y float64
	z float64
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.x + other.x, v.y + other.y, v.z + other.z}
}
func (v Vec3) MultConst(c float64) Vec3 {
	return Vec3{v.x * c, v.y * c, v.z * c}
}
func (v Vec3) Mult(other Vec3) Vec3 {
	return Vec3{other.x * v.x, other.y * v.y, other.z * v.z}
}
func (v Vec3) Length2() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}
func (v Vec3) Minus(other Vec3) Vec3 {
	return Vec3{v.x - other.x, v.y - other.y, v.z - other.z}
}

// in place
func (v *Vec3) Normalize() {
	var nor2 float64 = v.Length2()
	if nor2 > 0 {
		var invNor float64 = 1 / math.Sqrt(nor2)
		v.x = v.x * invNor
		v.y = v.y * invNor
		v.z = v.z * invNor
	}
}

func (v Vec3) Dot(other Vec3) float64 {
	return v.x*other.x + v.y*other.y + v.z*other.z
}
