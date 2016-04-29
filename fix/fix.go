package fix

import (
	"fmt"
	"math/big"
)

type Val struct {
	num big.Int
	denom big.Int
}

func New(nv, dv int64) *Val {
	var n, d big.Int
	n.SetInt64(nv)
	d.SetInt64(dv)
	var res Val
	res.Init(&n, &d)
	return &res
}

func (v *Val) AddFloat64(l Val, r float64) *Val {
	id := l.denom.Int64()
	return v.AddInt64(l, int64(r*float64(id)), id)
}

func (v *Val) AddInt64(l Val, rnv, rdv int64) *Val {
	var rn big.Int
	rn.SetInt64(rnv)

	ldv := l.denom.Int64()
	ln := &l.num

	if rdv != ldv {
		var f big.Int

		if ldv > rdv {
			f.SetInt64(ldv / rdv)
			rn.Mul(&rn, &f)
		} else {
			f.SetInt64(rdv / ldv)
			ln.Mul(ln, &f)
		}
	}

	v.denom.SetInt64(ldv)
	v.num.Add(ln, &rn)

	if ldv < rdv {
		var f big.Int
		f.SetInt64(rdv / ldv)
		v.num.Div(&v.num, &f)
	}

	return v
}

func (l *Val) Cmp(r Val) int {
	ln, ld := l.num, l.denom.Int64()
	rn, rd := r.num, r.denom.Int64()

	if ld != rd {
		var f big.Int

		if ld > rd {
			f.SetInt64(ld / rd)
			rn.Mul(&rn, &f)
		} else {
			f.SetInt64(rd / ld)
			ln.Mul(&ln, &f)
 		}
	}

	return ln.Cmp(&rn)
}

func (v *Val) Denom() big.Int {
	return v.denom
}

func (v *Val) Float64() float64 {
	var res big.Int
	res.Div(&v.num, &v.denom)
	iv := float64(res.Int64())
	res.Mul(&res, &v.denom)
	return iv + float64(v.num.Int64() - res.Int64()) / float64(v.denom.Int64())
}

func (v *Val) Frac() int64 {
	var res big.Int
	return (&res).Mod(&v.num, &v.denom).Int64()
}

func (v *Val) Init(n, d *big.Int) *Val {
	v.num = *n
	v.denom = *d
	return v
}

func (v *Val) Num() big.Int {
	return v.num
}

func (v *Val) Scale(d int64) *Val {
	vd := v.denom.Int64()

	if d != vd {
		var f big.Int

		if vd < d {
			f.SetInt64(d / vd)
			v.num.Mul(&v.num, &f)
		} else {
			f.SetInt64(vd / d)
			v.num.Div(&v.num, &f)
 		}
	}
	
	return v
}

func (v *Val) String() string {
	var res big.Int
	d, m := res.DivMod(&v.num, &v.denom, &v.denom)
	return fmt.Sprintf("%v.%v", d.Int64(), m.Int64())
}

func (v *Val) SubFloat64(l Val, r float64) *Val {
	ld := l.denom.Int64()
	return v.SubInt64(l, int64(r*float64(ld)), ld)
}

func (v *Val) SubInt64(l Val, d, m int64) *Val {
	var dv big.Int
	dv.SetInt64(d)
	lm := l.denom.Int64()
	lv := l.num
	
	if m != lm {
		var mv big.Int

		if lm > m {
			mv.SetInt64(lm / m)
			dv.Mul(&dv, &mv)
		} else {
			mv.SetInt64(m / lm)
			lv.Mul(&lv, &mv)
 		}
	}

	v.denom.SetInt64(lm)
	v.num.Sub(&lv, &dv)

	if lm < m {
		var mv big.Int
		mv.SetInt64(m / lm)
		v.num.Div(&v.num, &mv)
	}

	return v
}

func (v *Val) Trunc() int64 {
	var res big.Int
	return (&res).Div(&v.num, &v.denom).Int64()	
}