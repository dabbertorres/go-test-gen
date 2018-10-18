package main

import "testing"

func TestAdd(t *testing.T) {
	type (
		Input struct {
			x int
			y int
		}

		Output struct {
			int int
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			actual.int = Add(c.In.x, c.In.y)

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}

func TestAdd2(t *testing.T) {
	type (
		Input struct {
			x int
			y int
		}

		Output struct {
			int int
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			actual.int = Add2(c.In.x, c.In.y)

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}

func TestAdd3(t *testing.T) {
	type (
		Input struct {
			x int
			y int
			z int
		}

		Output struct {
			r int
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			actual.r = Add3(c.In.x, c.In.y, c.In.z)

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}

func TestAdd4(t *testing.T) {
	type (
		Input struct {
			x int
			y int
			z int
			w int
		}

		Output struct {
			a int
			b int
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			actual.a, actual.b = Add4(c.In.x, c.In.y, c.In.z, c.In.w)

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}

func TestAdd42(t *testing.T) {
	type (
		Input struct {
			x int
			y int
			z int
			w int
		}

		Output struct {
			a int
			b int
		}

		Case struct {
			Name   string
			In     Input
			Expect Output
		}
	)

	// TODO create test cases
	cases := []Case{}

	tester := func(c Case) func(*testing.T) {
		return func(t *testing.T) {
			var actual Output
			actual.a, actual.b = Add42(c.In.x, c.In.y, c.In.z, c.In.w)

			if actual != c.Expect {
				t.Errorf("expected %+v, actual: %+v\n", c.Expect, actual)
			}
		}
	}

	for _, c := range cases {
		t.Run(c.Name, tester(c))
	}
}
