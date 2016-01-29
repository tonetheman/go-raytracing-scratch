
package main

import ( "fmt" )

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

func main() {
	fmt.Println("hi");

	var v Vec3 = Vec3{0,0,0}
	var v2 Vec3 = Vec3{1,1,1}
	var v3 Vec3 = v.add(v2).multConst(3)

}
