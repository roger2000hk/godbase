package fix

import (
	"fmt"
	"math/big"
)

type Val struct {
	denom, num big.Int
}

func New(nv, dv int64) *Val {
	var n, d big.Int
	n.SetInt64(nv)
	d.SetInt64(dv)
	var res Val
	res.Init(&n, &d)
	return &res
}

func (v *Val) Add(l Val, r Val) *Val {
	return v.AddBig(l, r.num, r.denom)
}

func (v *Val) AddFloat64(l Val, rv float64) *Val {
	var fv, rfv big.Float
	fv.SetInt(&l.denom)
	rfv.SetFloat64(rv)
	fv.Mul(&fv, &rfv)

	var iv big.Int
	fv.Int(&iv)

	return v.AddBig(l, iv, l.denom)
}

func (v *Val) AddInt64(l Val, rnv, rdv int64) *Val {
	var rn, rd big.Int
	rn.SetInt64(rnv)
	rd.SetInt64(rdv)

	return v.AddBig(l, rn, rd)
}

func (v *Val) AddBig(l Val, rn, rd big.Int) *Val {
	ldv, rdv := l.denom.Int64(), rd.Int64()
	ln := l.num

	if ldv != rdv {
		scaleUp(&ln, ldv, &rn, rdv)
	}

	v.denom.SetInt64(ldv)
	v.num.Add(&ln, &rn)

	if ldv < rdv {
		scaleDown(&v.num, rdv/ldv)
	}

	return v
}

func (l *Val) Cmp(r Val) int {
	ln, ldv := l.num, l.denom.Int64()
	rn, rdv := r.num, r.denom.Int64()

	if ldv != rdv {
		scaleUp(&ln, ldv, &rn, rdv)
	}

	return ln.Cmp(&rn)
}

func (v *Val) Denom() big.Int {
	return v.denom
}

func (v *Val) Div(l Val, r Val) *Val {
	return v.DivBig(l, r.num, r.denom)
}

func (v *Val) DivBig(l Val, rn, rd big.Int) *Val {
	ln := l.num

	v.denom = l.denom
	v.num.Mul(&ln, &rd)
	v.num.Div(&ln, &rn)

	return v
}

func (v *Val) DivFloat64(l Val, rv float64) *Val {
	var fv, rfv big.Float
	fv.SetInt(&l.denom)
	rfv.SetFloat64(rv)
	fv.Mul(&fv, &rfv)

	var iv big.Int
	fv.Int(&iv)

	return v.DivBig(l, iv, l.denom)
}

func (v *Val) DivInt64(l Val, rnv, rdv int64) *Val {
	var rn, rd big.Int
	rn.SetInt64(rnv)
	rd.SetInt64(rdv)

	return v.DivBig(l, rn, rd)
}

func (v *Val) Float64() float64 {
	var res big.Int
	res.Div(&v.num, &v.denom)
	iv := float64(res.Int64())
	res.Mul(&res, &v.denom)
	return iv + float64(v.num.Int64()-res.Int64())/float64(v.denom.Int64())
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

func (v *Val) Mul(l Val, r Val) *Val {
	return v.MulBig(l, r.num, r.denom)
}

func (v *Val) MulBig(l Val, rn, rd big.Int) *Val {
	ln := l.num

	v.denom = l.denom
	v.num.Mul(&ln, &rn)
	v.num.Div(&ln, &rd)

	return v
}

func (v *Val) MulFloat64(l Val, rv float64) *Val {
	var fv, rfv big.Float
	fv.SetInt(&l.denom)
	rfv.SetFloat64(rv)
	fv.Mul(&fv, &rfv)
	
	var iv big.Int
	fv.Int(&iv)
	
	return v.MulBig(l, iv, l.denom)
}

func (v *Val) MulInt64(l Val, rnv, rdv int64) *Val {
	var rn, rd big.Int
	rn.SetInt64(rnv)
	rd.SetInt64(rdv)
	
	return v.MulBig(l, rn, rd)
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

func (v *Val) Sub(l Val, r Val) *Val {
	return v.SubBig(l, r.num, r.denom)
}

func (v *Val) SubBig(l Val, rn, rd big.Int) *Val {
	ldv, rdv := l.denom.Int64(), rd.Int64()
	ln := l.num

	if ldv != rdv {
		scaleUp(&ln, ldv, &rn, rdv)
	}

	v.denom.SetInt64(ldv)
	v.num.Sub(&ln, &rn)

	if ldv < rdv {
		scaleDown(&v.num, rdv/ldv)
	}

	return v
}

func (v *Val) SubFloat64(l Val, rv float64) *Val {
	var fv, rfv big.Float
	fv.SetInt(&l.denom)
	rfv.SetFloat64(rv)
	fv.Mul(&fv, &rfv)

	var iv big.Int
	fv.Int(&iv)

	return v.SubBig(l, iv, l.denom)
}

func (v *Val) SubInt64(l Val, rnv, rdv int64) *Val {
	var rn, rd big.Int
	rn.SetInt64(rnv)
	rd.SetInt64(rdv)

	return v.SubBig(l, rn, rd)
}

func (v *Val) Trunc() int64 {
	var res big.Int
	return (&res).Div(&v.num, &v.denom).Int64()
}

func scaleDown(n *big.Int, fv int64) {
	var f big.Int
	f.SetInt64(fv)
	n.Div(n, &f)
}

func scaleUp(ln *big.Int, ldv int64, rn *big.Int, rdv int64) {
	var f big.Int

	if ldv < rdv {
		f.SetInt64(rdv / ldv)
		ln.Mul(ln, &f)
	} else {
		f.SetInt64(ldv / rdv)
		rn.Mul(rn, &f)
	}
}
