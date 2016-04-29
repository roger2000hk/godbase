package decimal

import (
	"fmt"
	"math/big"
)

type Value struct {
	data big.Int
	mult big.Int
}

func New(d, m int64) *Value {
	return new(Value).Init(d, m)
}

func (v *Value) Init(d, m int64) *Value {
	v.data.SetInt64(d)
	v.mult.SetInt64(m)
	return v
}

func (v *Value) AddInt64(l *Value, d, m int64) *Value {
	var dv big.Int
	dv.SetInt64(d)

	vm := l.mult.Int64()
	if l.data.Int64() == 0 {
		vm = m
	}

	if m != vm {
		var mv big.Int
		if vm > m {
			mv.SetInt64(vm / m)
			dv.Mul(&dv, &mv)
		} else {
			mv.SetInt64(m / vm)
			dv.Div(&dv, &mv)
		}
	}

	v.mult.SetInt64(vm)
	v.data.Add(&l.data, &dv)

	return v
}

func (v *Value) SubInt64(l *Value, d, m int64) *Value {
	var dv big.Int
	dv.SetInt64(d)

	vm := l.mult.Int64()
	if l.data.Int64() == 0 {
		vm = m
	}

	if m != vm {
		var mv big.Int
		if vm > m {
			mv.SetInt64(vm / m)
			dv.Mul(&dv, &mv)
		} else {
			mv.SetInt64(m / vm)
			dv.Div(&dv, &mv)
		}
	}

	v.mult.SetInt64(vm)
	v.data.Sub(&l.data, &dv)

	return v
}

func (v *Value) Frac() int64 {
	var res big.Int
	return (&res).Mod(&v.data, &v.mult).Int64()
}

func (v *Value) String() string {
	var res big.Int
	d, m := (&res).DivMod(&v.data, &v.mult, &v.mult)
	return fmt.Sprintf("%v.%v", d.Int64(), m.Int64())
}

func (v *Value) Trunc() int64 {
	var res big.Int
	return (&res).Div(&v.data, &v.mult).Int64()	
}
