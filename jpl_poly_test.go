package jpl

import (
	"testing"
)

const prec = 1e-13 // precision maximum
func equal_f(a, b float64) (ok bool) {
	ok = -prec < b-a && b-a < prec
	return ok
}
func polyTchebT(n int, x float64) float64 {
	if 0 > n {
		panic("invalid parameter n")
	}
	if -1 > x || 1 < x {
		panic("invalid parameter x")
	}
	var p float64
	switch n {
	case 0:
		p = 1
	case 1:
		p = x
	default:
		p = 2*x*polyTchebT(n-1, x) - polyTchebT(n-2, x)
	}
	return p
}
func polyTchebU(n int, x float64) float64 {
	if 0 > n {
		panic("invalid parameter n")
	}
	if -1 > x || 1 < x {
		panic("invalid parameter x")
	}
	var p float64
	switch n {
	case 0:
		p = 1
	case 1:
		p = 2 * x
	default:
		p = 2*x*polyTchebU(n-1, x) - polyTchebU(n-2, x)
	}
	return p
}

var dataTest = []float64{-1, -0.5, -0.25, 0, 0.25, 0.5, 1}

func TestTchebT(t *testing.T) {
	nmax := Configuration.TchebNmax
	for _, x := range dataTest {
		for n := 0; n <= nmax; n++ {
			expected := polyTchebT(n, x)
			actual := TchebT[n].Value(x)
			if !equal_f(expected, actual) {
				t.Errorf("Tcheb1[%d](%f): expected %.13f, actual %.13f", n, x, expected, actual)
			}
		}
	}
}

func TestTchebU(t *testing.T) {
	nmax := Configuration.TchebNmax
	for _, x := range dataTest {
		for n := 0; n <= nmax; n++ {
			expected := polyTchebU(n, x)
			actual := TchebU[n].Value(x)
			if !equal_f(expected, actual) {
				t.Errorf("Tcheb1[%d](%f): expected %.13f, actual %.13f", n, x, expected, actual)
			}
		}
	}
}

func TestPoly(t *testing.T) {
	const prec = 1e-15
	var actual, expected Poly
	p1 := Poly{1, 2, 3}
	p2 := Poly{4, 5}
	p3 := Poly{4, 13, 22, 15} // = p1 * p2
	p4 := Poly{5, 7, 3}       // = p1 + p2
	p5 := Poly{-3, -3, 3}     // = p1 - p2
	p6 := Poly{3, 3, -3}      // = p2 - p1
	p7 := Poly{2, 4, 6}       // =  2 * p1
	actual = p1.MultPoly(p2)
	expected = p3
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("MultPoly p1*p2: expected %v, actual %v", expected, actual)
	}
	actual = p2.MultPoly(p1)
	expected = p3
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("MultPoly p2*p1: expected %v, actual %v", expected, actual)
	}
	actual = p1.AddPoly(p2)
	expected = p4
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("AddPoly p1+p2: expected %v, actual %v", expected, actual)
	}
	actual = p2.AddPoly(p1)
	expected = p4
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("AddPoly p2+p1: expected %v, actual %v", expected, actual)
	}
	actual = p1.SubPoly(p2)
	expected = p5
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("SubPoly p1-p2: expected %v, actual %v", expected, actual)
	}
	actual = p2.SubPoly(p1)
	expected = p6
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("SubPoly p1-p2: expected %v, actual %v", expected, actual)
	}
	actual = p1.Mult(2)
	expected = p7
	if !expected.EqualPoly(actual, prec) {
		t.Errorf("Mult 2*p1: expected %v, actual %v", expected, actual)
	}
}
