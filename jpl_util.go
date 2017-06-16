// ftp://ssd.jpl.nasa.gov/pub/eph/planets/
//https://hpiers.obspm.fr/iers/bul/bulc/Leap_Second.dat
package jpl

import (
	"os"
	"io"
	"fmt"
	"strings"
	"errors"
	"sync"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
// type pour la lecture successive des plusieurs fichiers comme s'il s'agissait d'un seul
type multiFile struct {
	files []*os.File
}
// lecture
func (r *multiFile) Read(p []byte) (n int, err error) {
	for len(r.files) > 0 {
		n, err = r.files[0].Read(p)
		if 0 == n && io.EOF == err {
			r.files[0].Close()
			err = nil
			r.files = r.files[1:]
			continue
		}
		return
	}
	return 0, io.EOF
}
// fermeture
func (r *multiFile) Close() error {
	for len(r.files) > 0 {
		r.files[0].Close()
		r.files = r.files[1:]
	}
	return nil
}
// lire un flottant
func read_double(r io.Reader, eol bool) (f float64, err error) {
	var s string
	if eol {
		_, err = fmt.Fscanf(r,"%s\n",&s)
	} else {
		_, err = fmt.Fscanf(r,"%s",&s)
	}
	if nil != err {
		return
	}
	fmt.Sscanf(strings.Replace(s,"D","e",1),"%f",&f)
	return
}
// ouvrir les fichiers
func jpl_open(dir string, files []string) (io.ReadCloser, error) {
	if !strings.HasSuffix(dir,"/") {
		dir += "/"
	}
	m := &multiFile{make([]*os.File,len(files))}
	for i,file := range files {
		reader, err := os.Open(dir+file)
		if nil != err {
			return nil, err
		}
		m.files[i] = reader
	}
	return m, nil
}
// lire les parametres 'ksize' et 'ncoeff'
func jpl_read_param(r io.Reader) (ksize, ncoeff int, err error) {
	var s string
	_, err = fmt.Fscanf(r,"KSIZE=%d NCOEFF=%d",&ksize,&ncoeff)
	if nil != err {
		return
	}
	_, err = fmt.Fscanf(r,"%s",&s)
	return
}
// lire le numero de groupe
func jpl_read_group_number(r io.Reader) (group int, err error) {
	var s string
	_, err = fmt.Fscanf(r,"GROUP %d",&group)
	if nil != err {
		return
	}
	_, err = fmt.Fscanf(r,"%s",&s)
	return
}
// lire le groupe 1010
func jpl_read_group_1010(r io.Reader) (ephemeris string, jed_start, jed_final float64, err error) {
	var s string
	var dt_start,dt_final,hr_start,hr_final string
	_, err = fmt.Fscanf(r,"JPL Planetary Ephemeris %s",&ephemeris)
	if nil != err {
		return
	}
	_, err = fmt.Fscanf(r,"Start Epoch: JED=%f %s %s\n",&jed_start,&dt_start,&hr_start)
	if nil != err {
		return
	}
	_, err = fmt.Fscanf(r,"Final Epoch: JED=%f %s %s\n",&jed_final,&dt_final,&hr_final)
	if nil != err {
		return
	}
	_, err = fmt.Fscanf(r,"%s",&s)
	return
}
// lire le groupe 1030
func jpl_read_group_1030(r io.Reader) (jd_min, jd_max, jd_int float64, err error) {
	var s string
	_, err = fmt.Fscanf(r,"%f %f %f",&jd_min,&jd_max,&jd_int)
	if nil != err {
		return
	}
	_, err = fmt.Fscanf(r,"%s",&s)
	if nil != err {
		return
	}
	return
}
// lire le groupe 1040
func jpl_read_group_1040(r io.Reader) (const_nam []string, err error) {
	var s string
	var const_num int
	_, err = fmt.Fscanf(r,"   %d",&const_num)
	if nil != err {
		return
	}
	const_nam = make([]string,const_num)
	for i := 0; i < const_num; i++ {
		if 9 == i % 10 {
			_, err = fmt.Fscanf(r,"%s\n",&const_nam[i])
		} else {
			_, err = fmt.Fscanf(r,"%s",&const_nam[i])
		}
		if nil != err {
			return
		}
	}
	_,err = fmt.Fscanf(r,"%s",&s)
	return
}
// lire le groupe 1041
func jpl_read_group_1041(r io.Reader) (const_val []float64, err error) {
	var s string
	var const_num int
	_, err = fmt.Fscanf(r,"%d",&const_num)
	if nil != err {
		return
	}
	const_val = make([]float64,const_num)
	
	for i := 0; i < const_num; i++ {
		const_val[i], err = read_double(r,2 == i % 3)
		if nil != err {
			return
		}
	}
	_, err = fmt.Fscanf(r,"%s",&s)
	return
}
// lire le groupe 1041
func jpl_read_group_1050(r io.Reader) (ipt [12][3]int,lpt ,rpt ,tpt [3]int, err error) {
	var s string
	for i := 0; i < 3; i++ {
		_, err = fmt.Fscanf(r,"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n",
			&ipt[0][i],&ipt[1][i],&ipt[2][i],&ipt[3][i],
			&ipt[4][i],&ipt[5][i],&ipt[6][i],&ipt[7][i],
			&ipt[8][i],&ipt[9][i],&ipt[10][i],&ipt[11][i],
			&lpt[i],&rpt[i],&tpt[i])
		if nil != err {
			return
		}
	}
	_, err = fmt.Fscanf(r,"%s",&s)
	return
}
// lire le groupe 1070
func jpl_read_group_1070(r io.Reader) (rnum int, rval []float64, err error) {
	var ctr int
	_, err = fmt.Fscanf(r,"%d %d",&rnum, &ctr)
	if nil != err {
		return
	}
	rval = make([]float64,ctr)
	for i := 0; i < ctr; i++ {
		rval[i], err = read_double(r,2 == i % 3)
		if nil != err {
			return
		}
	}
	return
}
// creer un pool
var db_pool = sync.Pool{New:
	func() interface{} {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
				Configuration.SQL.User,
				Configuration.SQL.Pwd,
				Configuration.SQL.Host,
				Configuration.SQL.Port,
				Configuration.SQL.Db))
		if err != nil {
			panic(err)
		}
		if err = db.Ping(); err != nil {
			panic(err)
		}
		return db
	},
}
// creer la base de donnees
func Jpl_createdb() {
	var (
	drop_de432_table = "DROP TABLE IF EXISTS %s;"
	create_de432_tables = map[string]string{
		"de432_header":`
		CREATE TABLE IF NOT EXISTS de432_header (
		  body_id INT NOT NULL,
		  body_name VARCHAR(10) NOT NULL,
		  jd0 DOUBLE NOT NULL,
		  jd1 DOUBLE NOT NULL,
		  dt DOUBLE NOT NULL,
		  ncf INT NOT NULL,
		  PRIMARY KEY (body_id)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1
		;
	`,
		"de432_consts" : `
		CREATE TABLE IF NOT EXISTS de432_consts (
		  const_name VARCHAR(6) NOT NULL,
		  val DOUBLE NOT NULL,
		  PRIMARY KEY (const_name)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1
		;
	`,
		"de432_coeffs" : `
		CREATE TABLE IF NOT EXISTS de432_coeffs (
		  body_id int(11) NOT NULL,
		  jd0 double NOT NULL,
		  jd1 double NOT NULL,
		  pv0 blob,
		  pv1 blob,
		  pv2 blob,
		  pv3 blob,
		  pv4 blob,
		  pv5 blob,
		  PRIMARY KEY (body_id,jd0)
		) ENGINE=InnoDB DEFAULT CHARSET=latin1
		PARTITION BY KEY (body_id)
		PARTITIONS 11
		;
	`,
	}
	)
	var err error
	db := db_pool.Get().(*sql.DB)
	defer db_pool.Put(db)

	generate := false
	for table, create_table := range create_de432_tables {
		_, err = db.Query(fmt.Sprintf("SELECT * FROM %s LIMIT 1;",table))
		if nil != err {
			generate = true
			_, err = db.Exec(fmt.Sprintf(drop_de432_table,table))
			if nil != err {
				//log.Fatal(err)
				panic(err)
			}
			_, err = db.Exec(create_table)
			if nil != err {
				//log.Fatal(err)
				panic(err)
			}
		}
	}
	if generate {
		jpl_generatedb()
	}
}
// initialiser la base de donnees
func jpl_generatedb() {
	files := []string{"header_571.432",
			"ascp01550.432",
			"ascp01650.432",
			"ascp01750.432",
			"ascp01850.432",
			"ascp01950.432",
			"ascp02050.432",
			"ascp02150.432",
			"ascp02250.432",
			"ascp02350.432",
			"ascp02450.432",
			"ascp02550.432",
	}
	r, err := jpl_open(Configuration.InputDir+"/de432",files)
	if nil != err {panic(err)}
	defer r.Close()

	db := db_pool.Get().(*sql.DB)
	defer db_pool.Put(db)

	d_header, err := db.Prepare("DELETE FROM de432_header WHERE body_id = ?;")
	if err != nil {
		panic(err)
//		log.Fatal(err)
	}
	defer d_header.Close()
	i_header, err := db.Prepare("INSERT INTO de432_header (body_id, body_name, jd0, jd1, dt, ncf) VALUES (?, ?, ?, ?, ?, ?);")
	if nil != err {
//		log.Fatal(err)
		panic(err)
	}
	defer i_header.Close()

	d_consts, err := db.Prepare("DELETE FROM de432_consts WHERE const_name = ?;")
	if err != nil {
		panic(err)
//		log.Fatal(err)
	}
	defer d_consts.Close()
	i_consts, err := db.Prepare("INSERT INTO de432_consts(const_name, val) VALUES (?,?);")
	if nil != err {
//		log.Fatal(err)
		panic(err)
	}
	defer i_consts.Close()

	d_coeffs, err := db.Prepare("DELETE FROM de432_coeffs WHERE body_id = ? AND jd0 = ?;")
	if err != nil {
		panic(err)
//		log.Fatal(err)
	}
	defer d_coeffs.Close()
	i_coeffs, err := db.Prepare("INSERT INTO de432_coeffs (body_id, jd0, jd1, pv0, pv1, pv2, pv3, pv4, pv5) VALUES (?,?,?,?,?,?,?,?,?);")
	if nil != err {
//		log.Fatal(err)
		panic(err)
	}
	defer i_coeffs.Close()

	_/*ksize*/, _/*ncoeff*/, err = jpl_read_param(r)

	_/*g1010*/, err = jpl_read_group_number(r)
	_/*ephemeris*/,_/*jed_start*/,_/*jed_final*/, err = jpl_read_group_1010(r)

	_/*g1030*/, err = jpl_read_group_number(r)
	jd_min, jd_max, _/*jd_int*/, err := jpl_read_group_1030(r)

	_/*g1040*/, err = jpl_read_group_number(r)
	const_nam, err := jpl_read_group_1040(r)

	_/*g1041*/, err = jpl_read_group_number(r)
	const_val, err := jpl_read_group_1041(r)

	_/*g1050*/, err = jpl_read_group_number(r)
	ipt, lpt, rpt, tpt, err := jpl_read_group_1050(r)
	fmt.Println("ipt",ipt)
	fmt.Println("lpt",lpt)
	fmt.Println("rpt",rpt)
	fmt.Println("tpt",tpt)
	var au, emrat float64
	for i, cname := range const_nam {
		if ("AU" == cname) {
			au = const_val[i]
		}
		if ("EMRAT" == cname) {
			emrat = const_val[i]
		}
		if _, err := d_consts.Exec(cname); nil != err {
			panic(err)
//			log.Fatal(err)
		}
		if _, err := i_consts.Exec(cname,const_val[i]); nil != err {
			panic(err)
//			log.Fatal(err)
		}
	}
	for i := 0; i < 11; i++ {	// bodies
		by := Jpl_body(i+1)
		dt := float64(32 / ipt[i][2])
		ncf := ipt[i][1]
		if _, err := d_header.Exec(by); nil != err {
			panic(err)
//			log.Fatal(err)
		}
		if _, err := i_header.Exec(by,by.String(),jd_min,jd_max,dt,ncf); nil != err {
			panic(err)
//			log.Fatal(err)
		}
		if MOON == by {
			if _, err := d_header.Exec(EARTH); nil != err {
				panic(err)
//				log.Fatal(err)
			}
			if _, err := i_header.Exec(EARTH,EARTH.String(),jd_min,jd_max,dt,ncf); nil != err {
				panic(err)
//				log.Fatal(err)
			}
		}
	}

	fmt.Println("au",au,"emrat",emrat)
	// coefficients
	_/*g1070*/, err = jpl_read_group_number(r)
	for {
		rnum, rval, err := jpl_read_group_1070(r)
		if io.EOF == err {
			break
		} else if nil != err {
			panic(err)
		} else {
			rec := &jpl_record_t{num:rnum,rec:&rval,ipt:&ipt}
			polys := rec.jpl_polys(au,emrat)
			for i := 0; i < 12; i++ {	// bodies
				by := Jpl_body(i+1)
				pvm, ok := polys[by]
				for j := 0; j < len(pvm) && ok; j++ {	// granules
					p := pvm[j]
					b, err := p.Blob()
					if nil != err {
						panic(err)
					}
					if _, err := d_coeffs.Exec(by,p.jd0); nil != err {
						panic(err)
//						log.Fatal(err)
					}
					if _, err := i_coeffs.Exec(by,p.jd0,p.jd1,b[0],b[1],b[2],b[3],b[4],b[5]); nil != err {
						panic(err)
//						log.Fatal(err)
					}
				}
			}
		}
	}
}
// lire les polynomes
func jpl_read_poly(body Jpl_body, jd float64) (*jpl_poly_t, error) {
	var (	
		body_id int
		jd0, jd1 float64
		b [6][]byte 
	)
	db := db_pool.Get().(*sql.DB)
	defer db_pool.Put(db)

	s_coeffs, err := db.Prepare("SELECT body_id, jd0, jd1, pv0, pv1, pv2, pv3, pv4, pv5 FROM de432_coeffs WHERE body_id = ? AND jd0 <= ? AND jd1 > ?;")
	if err != nil {
		panic(err)
//		log.Fatal(err)
	}
	defer s_coeffs.Close()
	if err := s_coeffs.QueryRow(body,jd,jd).Scan(&body_id, &jd0, &jd1, &b[0], &b[1], &b[2], &b[3], &b[4], &b[5]); nil != err {
		return nil, err
	}
	poly, err := new_poly_t(body,jd0,jd1,&b)
	return poly, err
}
// lire les entetes
func Jpl_read_header() (map[Jpl_body]jpl_header_t, error) {
	db := db_pool.Get().(*sql.DB)
	defer db_pool.Put(db)
	s_header, err := db.Prepare("SELECT body_id, body_name, jd0, jd1, dt, ncf FROM de432_header;")
	if err != nil {
		panic(err)
//		log.Fatal(err)
	}
	defer s_header.Close()
	rows, err := s_header.Query()
	if nil != err {
		return nil, err
	}
	defer rows.Close()
	headers := make(map[Jpl_body]jpl_header_t,11)
	for rows.Next() {
		var (	
			h jpl_header_t
		)
		if err := rows.Scan(&h.body_id,&h.body_name,&h.jd0,&h.jd1,&h.dt,&h.ncf); nil != err {
			return nil, err
		}
		headers[h.body_id] = h
	}
	return headers, nil
}
// lire les constantes
func Jpl_read_consts() (Jpl_consts_t, error) {
	db := db_pool.Get().(*sql.DB)
	defer db_pool.Put(db)
	s_consts, err := db.Prepare("SELECT const_name, val FROM de432_consts;")
	if err != nil {
		panic(err)
//		log.Fatal(err)
	}
	defer s_consts.Close()
	rows, err := s_consts.Query()
	if nil != err {
		return nil, err
	}
	defer rows.Close()
	consts := make(Jpl_consts_t)
	for rows.Next() {
		var (
			name string
			val float64
		)
		if err := rows.Scan(&name,&val); nil != err {
			return nil, err
		}
		consts[name] = val
	}
	return consts, nil
}
// obtenir une constante particuliere
func Jpl_get_const(name string) (float64, error) {
	var val float64
	db := db_pool.Get().(*sql.DB)
	defer db_pool.Put(db)
	s_consts, err := db.Prepare("SELECT val FROM de432_consts WHERE const_name = ?;")
	if err != nil {
		return val, err
	}
	defer s_consts.Close()
	if err := s_consts.QueryRow(name).Scan(&val); nil != err {
		return val, err
	}
	return val, nil
}
// creer une éphéméride
func Jpl_eph(body Jpl_body, jd float64) (*[6]float64, error) {
	if 0 > body || 12 < body {
		return nil, errors.New(fmt.Sprintf("Invalid body %d", body))
	}
	if SSB == body {
		return &[6]float64{0,0,0,0,0,0}, nil
	}
	poly, err := jpl_read_poly(body, jd)
	if nil != err {
		return nil, err
	}
	pv, err := poly.Values(jd)
	if nil != err {
		return nil, err
	}
	if MOON == body || EARTH == body {
		p3, err := jpl_read_poly(EMB, jd)
		if nil != err {
			return nil, err
		}
		pv3, err := p3.Values(jd)
		if nil != err {
			return nil, err
		}
		for i := 0; i < 6; i++ {
			pv[i] += pv3[i]
		}
	} 
	return pv, err
}
