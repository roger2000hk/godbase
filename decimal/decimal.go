package decimal

import (
	"fmt"
	"math/big"
)

type Value struct {
	data big.Int
	mult big.Int
}

func Cons(d, m big.Int) *Value {
	return new(Value).Init(d, m)
}

func New(dv, mv int64) (res Value) {
	var d, m big.Int
	d.SetInt64(dv)
	m.SetInt64(mv)
	res.Init(d, m)
	return res
}

func (v *Value) AddFloat64(l *Value, r float64) *Value {
	im := l.mult.Int64()
	return v.AddInt64(l, int64(r*float64(im)), im)
}

func (v *Value) AddInt64(l *Value, d, m int64) *Value {
	var dv big.Int
	dv.SetInt64(d)
	lm := l.mult.Int64()

	if m != lm {
		lv := &l.data
		var mv big.Int

		if lm > m {
			mv.SetInt64(lm / m)
			dv.Mul(&dv, &mv)
		} else {
			mv.SetInt64(m / lm)
			lv.Mul(lv, &mv)
		}
	}

	v.mult.SetInt64(lm)
	v.data.Add(&l.data, &dv)

	if lm < m {
		var mv big.Int
		mv.SetInt64(m / lm)
		v.data.Mul(&v.data, &mv)
	}

	return v
}

func (v *Value) Data() big.Int {
	return v.data
}

func (l *Value) Cmp(r *Value) int {
	lm, lv := l.mult.Int64(), l.data
	rm, rv := r.mult.Int64(), r.data

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

func (v *Value) Float64() float64 {
	var res big.Int
	res.Div(&v.data, &v.mult)
	iv := float64(res.Int64())
	res.Mul(&res, &v.mult)
	return iv + float64(v.data.Int64() - res.Int64()) / float64(v.mult.Int64())
}

func (v *Value) Frac() int64 {
	var res big.Int
	return (&res).Mod(&v.data, &v.mult).Int64()
}

func (v *Value) Init(d, m big.Int) *Value {
	v.data = d
	v.mult = m
	return v
}

func (v *Value) Mult() big.Int {
	return v.mult
}

func (v *Value) String() string {
	var res big.Int
	d, m := res.DivMod(&v.data, &v.mult, &v.mult)
	return fmt.Sprintf("%v.%v", d.Int64(), m.Int64())
}

func (v *Value) SubFloat64(l *Value, r float64) *Value {
	im := l.mult.Int64()
	return v.SubInt64(l, int64(r*float64(im)), im)
}

func (v *Value) SubInt64(l *Value, d, m int64) *Value {
	var dv big.Int
	dv.SetInt64(d)
	lm := l.mult.Int64()
	lv := l.data
	
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

	v.mult.SetInt64(lm)
	v.data.Sub(&lv, &dv)

	if lm < m {
		var mv big.Int
		mv.SetInt64(m / lm)
		v.data.Div(&v.data, &mv)
	}

	return v
}

func (v *Value) Trunc() int64 {
	var res big.Int
	return (&res).Div(&v.data, &v.mult).Int64()	
}
