package main

func Add(x, y int) int {
	return x + y
}

func Add2(x int, y int) int {
	return x + y
}

func Add3(x int, y, z int) (r int) {
	r = x + y + z
	return
}

func Add4(x, y int, z, w int) (a, b int) {
	a = x + y
	b = z + w
	return
}

func Add42(x, y int, z, w int) (a int, b int) {
	a = x + y
	b = z + w
	return
}
