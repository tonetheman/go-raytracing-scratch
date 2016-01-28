
package main

import ( "fmt" )

func goo(a * float64) {
	*a = 20.0;
}

func main() {
	fmt.Println("hi");

	var a float64 = 10;

	fmt.Println(a);

	goo(&a);

	fmt.Println(a);

}
