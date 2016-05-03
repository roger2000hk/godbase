package fix

// package godbase/fix implements speedy, comparable fixed-points with int64 nums/denoms 

import (
	"fmt"
)

type Val struct {
	denom, num int64
}

func New(n, d int64) *Val {
	return new(Val).Init(n, d)
}

func (self *Val) Add(l Val, r Val) *Val {
	return self.AddInt64(l, r.num, r.denom)
}

func (self *Val) AddFloat64(l Val, r float64) *Val {
	return self.AddInt64(l, int64(r * float64(l.denom)), l.denom)
}

func (self *Val) AddInt64(l Val, rn, rd int64) *Val {
	if l.denom != rd {
		if l.denom < rd {
			l.num *= rd / l.denom
		} else {
			rn *= l.denom / rd
		}
	}

	self.denom = l.denom
	self.num = l.num + rn

	if l.denom < rd {
		self.num /= (rd / l.denom)
	}

	return self
}

func (l *Val) Cmp(r Val) int {
	ln, ld := l.num, l.denom
	rn, rd := r.num, r.denom

	if ld != rd {
		if ld < rd {
			ln *= rd / ld
		} else {
			rn *= ld / rd
		}
	}

	if ln < rn {
		return -1
	}

	if ln == rn {
		return 0
	}

	return 1
}

func (self *Val) Denom() int64 {
	return self.denom
}

func (self *Val) Div(l Val, r Val) *Val {
	return self.DivInt64(l, r.num, r.denom)
}

func (self *Val) DivFloat64(l Val, r float64) *Val {
	return self.DivInt64(l, int64(r * float64(l.denom)), l.denom)
}

func (self *Val) DivInt64(l Val, rn, rd int64) *Val {
	self.denom = l.denom
	self.num = l.num * rn / rd
	return self
}

func (self *Val) Float64() float64 {
	tr, fr := self.Frac()
	return float64(tr) + float64(fr) / float64(self.denom)
}

func (self *Val) Frac() (int64, int64) {
	trunc := self.num / self.denom
	frac := (self.num - trunc * self.denom)
	return trunc, frac
}

func (self *Val) Init(n, d int64) *Val {
	self.num = n
	self.denom = d
	return self
}

func (self *Val) Mul(l Val, r Val) *Val {
	return self.MulInt64(l, r.num, r.denom)
}

func (self *Val) MulFloat64(l Val, r float64) *Val {
	return self.MulInt64(l, int64(r * float64(l.denom)), l.denom)
}

func (self *Val) MulInt64(l Val, rn, rd int64) *Val {
	self.denom = l.denom
	self.num = l.num /rn * rd
	return self
}

func (self *Val) Num() int64 {
	return self.num
}

func (self *Val) Scale(d int64) *Val {
	if d != self.denom {
		if self.denom < d {
			self.num *= d / self.denom
		} else {
			self.num /= self.denom / d
		}
	}

	return self
}

func (self *Val) String() string {
	trunc, frac := self.Frac()
	return fmt.Sprintf("%v.%v", trunc, frac)
}

func (self *Val) Sub(l Val, r Val) *Val {
	return self.AddInt64(l, -r.num, r.denom)
}

func (self *Val) SubFloat64(l Val, rv float64) *Val {
	return self.AddFloat64(l, -rv)
}

func (self *Val) SubInt64(l Val, rn, rd int64) *Val {
	return self.AddInt64(l, -rn, rd)
}
