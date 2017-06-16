package jpl

import (

)
// egalite a 'prec' pres
func equal(a, b, prec float64) bool {
	if 0 > prec {panic("Invalid precision")}
	ok := -prec < b - a && b - a < prec
	return ok
}
// type polynome
type Poly []float64
// valeur du polynome en 'x'
func (p Poly)Value(x float64) float64 {
	var xx, v float64 = 1, 0
	n := len(p) - 1
	for k := 0; k <= n; k++ {
		v += xx * p[k]
		xx *= x		
	}
	return v
}
// degre du polynome
func (p Poly)Degree() int {
	return len(p) - 1
}
// polynome derive
func (p Poly)Diff() Poly {
	n := len(p) - 1
	if 0 > n {
		panic("Invalid operation")
	}
	v := make(Poly,n)
	for k := 1; k <= n; k++ {
		v[k-1] = float64(k)*p[k]
	}
	return v
}
// polynome multiplie par 'f' * 'p'
func (p Poly)Mult(f float64) Poly {	// p * f
	v := make(Poly,len(p))
	for k := 0; k < len(p); k++ {
		v[k] = p[k] * f
	}
	return v
}
// polynome somme 'p' + 'a'
func (p Poly)AddPoly(a Poly) Poly {	// p + a
	lp := len(p)
	la := len(a)
	var v Poly
	if la < lp {
		v = make(Poly,lp)
		for k := 0; k < lp; k++ {
			v[k] = p[k]			
		}
		for k := 0; k < la; k++ {
			v[k] += a[k]
		}
	} else {
		v = make(Poly,la)
		for k := 0; k < la; k++ {
			v[k] = a[k]			
		}
		for k := 0; k < lp; k++ {
			v[k] += p[k]
		}
	}
	return v
}
// polynome soustraction 'p' - 'a'
func (p Poly)SubPoly(a Poly) Poly {	// p - a
	lp := len(p)
	la := len(a)
	var v Poly
	if la < lp {
		v = make(Poly,lp)
		for k := 0; k < lp; k++ {
			v[k] = p[k]			
		}
		for k := 0; k < la; k++ {
			v[k] -= a[k]
		}
	} else {
		v = make(Poly,la)
		for k := 0; k < la; k++ {
			v[k] = -a[k]			
		}
		for k := 0; k < lp; k++ {
			v[k] += p[k]
		}
	}
	return v
}
// polynome mutiplication 'p' * 'a'
func (p Poly)MultPoly(a Poly) Poly {
	lp := len(p)
	la := len(a)
	var v Poly = make(Poly,lp+la-1)
	for i := 0; i < lp; i++ {
		for j := 0; j < la; j++ {
			v[i+j] += p[i] * a[j]
		}
	}
	return v
}
// verification de l'egalite a la precision 'prec' pres 'p' == 'a'
func (p Poly)EqualPoly(a Poly, prec float64) bool {	// p = a
	if 0 > prec {panic("Invalid negative prec")}
	lp := len(p)
	la := len(a)
	v := true
	if la < lp {
		for k := 0; k < la && v; k++ {
			v = v && equal(a[k],p[k],prec)
		}
		for k := la+1; k < lp && v; k++ {
			v = v && equal(0,p[k],prec)	
		}
	} else {
		for k := 0; k < lp && v; k++ {
			v = v && equal(p[k],a[k],prec)	
		}
		for k := lp+1; k < la && v; k++ {
			v = v && equal(0,a[k],prec)	
		}
	}
	return v
}

