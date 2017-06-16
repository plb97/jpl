package jpl

import (
	"bytes"
	"encoding/binary"
	"errors"
	"encoding/json"
	"fmt"
	"sort"
)

var TchebT = []Poly{{1},{0,1},}
var TchebTd = []Poly{{},{1},}
var TchebU = []Poly{{1},{0,2},}

func init() {
	nmax := Configuration.TchebNmax
	p2 := Poly{0,2}
	for n := 2; n <= nmax; n++ {
		p := p2.MultPoly(TchebT[n-1]).SubPoly(TchebT[n-2])
		TchebT = append(TchebT,p)
		p = p.Diff()
		TchebTd = append(TchebTd,p)
		p = p2.MultPoly(TchebU[n-1]).SubPoly(TchebU[n-2])
		TchebU = append(TchebU,p)
	}
}

//type jpl_coord int
//const (
//	X = iota
//	Y
//	Z
//	Xd
//	Yd
//	Zd
//)
//func (o jpl_coord) String() string {
//	switch o {
//		case X:	return "X"
//		case Y:	return "Y"
//		case Z:	return "Z"
//		case Xd:	return "Xd"
//		case Yd:	return "Yd"
//		case Zd:	return "Zd"
//		default:	return "UNKNOWN"
//	}
//}
//
/*	JPL
      1 = mercury           8 = neptune
      2 = venus             9 = pluto
      3 = earth            10 = moon
      4 = mars             11 = sun
      5 = jupiter          12 = solar-system barycenter
      6 = saturn           13 = earth-moon barycenter
      7 = uranus           14 = nutations in longitude and obliquity
                                   15 = librations (if they exist on the file)
*/
type Jpl_body int
const (
	SSB Jpl_body = iota			// 0 solar-system barycenter
	MERCURY 				// 1
	VENUS					// 2
	EMB					// 3 earth-moon barycenter
	MARS					// 4
	JUPITER					// 5
	SATURN					// 6
	URANUS					// 7
	NEPTUNE					// 8
	PLUTO					// 9
	MOON					// 10
	SUN					// 11
	EARTH					// 12
//	NUTATIONS				// 14 in longitude and obliquity
//	LIBRATIONS				// 15 if they exist on the file
)
func (o Jpl_body)String() string {
	switch o {
		case MERCURY:		return "MERCURY"
		case VENUS:		return "VENUS"
		case EMB:		return "EMB"
		case MARS:		return "MARS"
		case JUPITER:		return "JUPITER"
		case SATURN:		return "SATURN"
		case URANUS:		return "URANUS"
		case NEPTUNE:		return "NEPTUNE"
		case PLUTO:		return "PLUTO"
		case MOON:		return "MOON"
		case SUN:		return "SUN"
		case EARTH:		return "EARTH"
//		case SSB:		return "SSB"
//		case NUTATIONS:		return "NUTATIONS"
//		case LIBRATIONS:	return "LIBRATIONS"
		default: 			return "UNKNOWN"
	}
}
// structure des polynomes
type jpl_poly_t struct {
	body Jpl_body
	jd0  float64
	jd1  float64
	pv   [6]Poly
}
// creer une structure polynome
func new_poly_t(body Jpl_body, jd0 float64, jd1 float64, b *[6][]byte) (*jpl_poly_t, error) {
	const float64_size = 8
	p := jpl_poly_t{body:body,jd0:jd0,jd1:jd1,}
	for i := 0; i < 6; i++ {
		buf := bytes.NewReader(b[i])
		n := len(b[i]) / float64_size
		p.pv[i] = make([]float64,n)
		for j := 0; j < n; j++ {
			if err := binary.Read(buf,binary.LittleEndian,&p.pv[i][j]); nil != err {
				return nil, err
			}
		}
	}
	return &p, nil
}
// format binaire 'little endian'
func (o *jpl_poly_t) Blob() (*[6][]byte, error) {
	r := (*o).pv // 6 polynomes de Tchebychev
	var b [6][]byte
	buf := new(bytes.Buffer)
	for i := 0; i < 6; i++ {
		buf.Reset()
		p := r[i]
		for j := 0; j < len(p); j++ {
			if err := binary.Write(buf,binary.LittleEndian,p[j]); nil != err {
				return nil, err
			}
		}
		b[i] = make([]byte,buf.Len())
		copy(b[i],buf.Bytes())
	}
	return &b,nil
}
// calcul de la valeur en fonction du jour Julien
func (o *jpl_poly_t) Values(jd float64) (*[6]float64, error) {
	if o.jd0 > jd || o.jd1 < jd {
		return nil, errors.New("Invalid jd") 
	}
	var pv [6]float64
	t := 2*(jd - o.jd0) / (o.jd1 - o.jd0) - 1
	for i := 0; i < 6; i++ {
		pv[i] = o.pv[i].Value(t)
	}
	return &pv, nil
}
// format 'JSON'
func (o *jpl_poly_t) JsonString() string {
	var buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	buf.WriteString("{\"body\":")
	enc.Encode(o.body)
	buf.Truncate(buf.Len()-1)
	buf.WriteString(",\"jd0\":")
	enc.Encode(o.jd0)
	buf.Truncate(buf.Len()-1)
	buf.WriteString(",\"jd1\":")
	enc.Encode(o.jd1)
	buf.Truncate(buf.Len()-1)
	buf.WriteString(",\"pv\":")
	enc.Encode(o.pv)
	buf.Truncate(buf.Len()-1)
	buf.WriteString("}")
	return buf.String()
}
// table des polynomes
type jpl_polys_t map[Jpl_body][]jpl_poly_t
// enregistrement
type jpl_record_t struct {
	num int
	rec *[]float64
	ipt *[12][3]int
}
// representation sous forme de chaine
func (o jpl_record_t) String() string {
	return fmt.Sprintf("[%d %v %v]",o.num,*o.rec,*o.ipt)
}
// jour Julien inferieur
func (o *jpl_record_t) Jd0() float64 {
	r := *o.rec
	return r[0]
}
// jour Julien superieur
func (o *jpl_record_t) Jd1() float64 {
	r := *o.rec
	return r[1]
}
// calcul de la table des polynomes
func (o *jpl_record_t) jpl_polys(au, emrat float64) jpl_polys_t {
	r := *o.rec
	f := 1 / (1 + emrat)
	jd0 := r[0]
	jd1 := r[1]
	polys := make(jpl_polys_t,12)
	for i := 0; i < 11; i ++ {						// body
		b := Jpl_body(i+1)
		ipt := (*o.ipt)[i]
		dt := (jd1 - jd0) / float64(ipt[2])
		a := 2 / dt
		for j := 0; j < ipt[2]; j++ { 					// granule
			pv := jpl_poly_t{body:b,jd0:jd0+float64(j)*dt,jd1:jd0+float64(j+1)*dt}
			for k := 0; k < 3; k++ { 					// coord
				n:= ipt[0] + ipt[1]*(3*j+k) - 1
				c := r[n:n+ipt[1]]
				pv.pv[k] = TchebT[0].Mult(c[0]/au)
				pv.pv[k+3] = TchebTd[0].Mult(c[0]/au)
				for l := 1; l < ipt[1]; l++ {			// coef
					pv.pv[k] = pv.pv[k].AddPoly(TchebT[l].Mult(c[l]/au))			// au
					pv.pv[k+3] = pv.pv[k+3].AddPoly(TchebTd[l].Mult(a*c[l]/au))		// au/d
				}
			}
			if MOON == b {
				pve := jpl_poly_t{body:EARTH,jd0:jd0+float64(j)*dt,jd1:jd0+float64(j+1)*dt}
				pvm := jpl_poly_t{body:MOON,jd0:jd0+float64(j)*dt,jd1:jd0+float64(j+1)*dt}
				for k := 0; k < 6; k++ {
					pve.pv[k] = pv.pv[k].Mult(-f)
					pvm.pv[k] = pv.pv[k].Mult(1-f)
					if !pv.pv[k].EqualPoly(pvm.pv[k].SubPoly(pve.pv[k]),1.e-15) {
						panic("Earth Moon error")
					}
				}
				polys[EARTH] = append(polys[EARTH],pve)
				polys[MOON] = append(polys[MOON],pvm)
			} else {
				polys[b] = append(polys[b],pv)
			}
		}
	}
	return polys
}
// structure representant les entetes
type jpl_header_t struct {
	body_id   Jpl_body
	body_name string
	jd0       float64
	jd1       float64
	dt        float64
	ncf       int
}
// numero de planete
func (h jpl_header_t) Body() Jpl_body {
	return h.body_id
}
// intitule de la planete
func (h jpl_header_t) Body_name() string {
	return h.body_name
}
// jour Julien inferieur
func (h jpl_header_t) Jd0() float64 {
	return h.jd0
}
// jour Julien superieur
func (h jpl_header_t) Jd1() float64 {
	return h.jd1
}
// nombre de dates 'jd0' et 'jd1'
func (h jpl_header_t) Dt() float64 {
	return h.dt
}
// nombre de coefficients
func (h jpl_header_t) Ncf() int {
	return h.ncf
}
// representation sous forme de chaine
func (h jpl_header_t) String() string {
	return fmt.Sprintf("[body:%v jd0:%v jd1:%v dt:%v ncf:%v]",h.body_id,h.jd0,h.jd1,h.dt,h.ncf)
}
// table des constantes
type Jpl_consts_t map[string]float64
func (c Jpl_consts_t) Keys() []string {
	lk := make([]string,0)
	for k := range c {
		lk = append(lk,k)
	}
	sort.Strings(lk)
	return lk
}
// obtenir une constante particuliere
func (c Jpl_consts_t) Get(k string) float64 {
	return c[k]
}
// representation sous forme d'une chaine
func (c Jpl_consts_t) String() string {
	return fmt.Sprintf("%v",c)
}
