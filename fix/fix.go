package fix

import (
	"fmt"
	"math/big"
)

type Val struct {
	num big.Int
	denom big.Int
}

func New(dv, mv int64) *Val {
	var d, m big.Int
	d.SetInt64(dv)
	m.SetInt64(mv)
	var res Val
	res.Init(&d, &m)
	return &res
}

func (v *Val) AddFloat64(l Val, r float64) *Val {
	im := l.denom.Int64()
	return v.AddInt64(l, int64(r*float64(im)), im)
}

func (v *Val) AddInt64(l Val, d, m int64) *Val {
	var dv big.Int
	dv.SetInt64(d)
	lm := l.denom.Int64()
	lv := &l.num

	if m != lm {
		var mv big.Int

		if lm > m {
			mv.SetInt64(lm / m)
			dv.Mul(&dv, &mv)
		} else {
			mv.SetInt64(m / lm)
			lv.Mul(lv, &mv)
		}
	}

	v.denom.SetInt64(lm)
	v.num.Add(lv, &dv)

	if lm < m {
		var mv big.Int
		mv.SetInt64(m / lm)
		v.num.Div(&v.num, &mv)
	}

	return v
}

func (l *Val) Cmp(r Val) int {
	lm, lv := l.denom.Int64(), l.num
	rm, rv := r.denom.Int64(), r.num

	if lm != rm {
		var m big.Int

		if lm > rm {
			m.SetInt64(lm / rm)
			rv.Mul(&rv, &m)
		} else {
			m.SetInt64(rm / lm)
			lv.Mul(&lv, &m)
 		}
	}

	return lv.Cmp(&rv)
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

func (v *Val) Init(d, m *big.Int) *Val {
	v.num = *d
	v.denom = *m
	return v
}

func (v *Val) Num() big.Int {
	return v.num
}

func (v *Val) Scale(m int64) *Val {
	vm := v.denom.Int64()

	if m != vm {
		var mi big.Int

		if vm < m {
			mi.SetInt64(m / vm)
			v.num.Mul(&v.num, &mi)
		} else {
			mi.SetInt64(vm / m)
			v.num.Div(&v.num, &mi)
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
	im := l.denom.Int64()
	return v.SubInt64(l, int64(r*float64(im)), im)
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
