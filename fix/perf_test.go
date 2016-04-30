package fix

import (
	//"fmt"
	"github.com/oguzbilgic/fpd"
	"github.com/shopspring/decimal"
	"math/rand"
	"testing"
)

func rndSlice(n int) []int {
	res := make([]int, n)

	for i := 0; i < n; i++ {
		res[i] = rand.Intn(n) + 1
	}

	return res
}

const testReps = 40000
var addNums = rndSlice(testReps)
var divNums = rndSlice(testReps)
var mulNums = rndSlice(testReps)
var subNums = rndSlice(testReps)

func BenchmarkFix(t *testing.B) {
	fv := *New(0, 100)

	for i := 0; i < testReps; i++ {
		fv.AddInt64(fv, int64(addNums[i]), 100)
		fv.DivInt64(fv, int64(divNums[i]), 100)
		fv.MulInt64(fv, int64(mulNums[i]), 100)
		fv.SubInt64(fv, int64(subNums[i]), 100)
	}

 	//fmt.Printf("fix: %v\n", &fv)
}

func BenchmarkFixScale(t *testing.B) {
	fv := *New(0, 10)

	for i := 0; i < testReps; i++ {
		fv.AddInt64(fv, int64(addNums[i]), 1000)
		fv.DivInt64(fv, int64(divNums[i]), 1000)
		fv.MulInt64(fv, int64(mulNums[i]), 1000)
		fv.SubInt64(fv, int64(subNums[i]), 1000)
	}
}

func BenchmarkFpd(t *testing.B) {
	fv := fpd.New(0, -2)

	for i := 0; i < testReps; i++ {
		fv = fv.Add(fpd.New(int64(addNums[i]), -2))
		fv = fv.Div(fpd.New(int64(divNums[i]), -2))
		fv = fv.Mul(fpd.New(int64(mulNums[i]), -2))
		fv = fv.Sub(fpd.New(int64(subNums[i]), -2))
	}

 	//fmt.Printf("fpd: %v\n", fv)
}

func BenchmarkFpdScale(t *testing.B) {
	fv := fpd.New(0, -1)

	for i := 0; i < testReps; i++ {
		fv = fv.Add(fpd.New(int64(addNums[i]), -3))
		fv = fv.Div(fpd.New(int64(divNums[i]), -3))		
		fv = fv.Mul(fpd.New(int64(mulNums[i]), -3))		
		fv = fv.Sub(fpd.New(int64(subNums[i]), -3))
	}
}

func BenchmarkDecimal(t *testing.B) {
	fv := decimal.New(0, -2)

	for i := 0; i < testReps; i++ {
		fv = fv.Add(decimal.New(int64(addNums[i]), -2))
		fv = fv.Div(decimal.New(int64(divNums[i]), -2))
		fv = fv.Mul(decimal.New(int64(mulNums[i]), -2))
		fv = fv.Sub(decimal.New(int64(subNums[i]), -2))
	}

 	//fmt.Printf("decimal: %v\n", fv)
}

func BenchmarkDecimalScale(t *testing.B) {
	fv := decimal.New(0, -1)

	for i := 0; i < testReps; i++ {
		fv = fv.Add(decimal.New(int64(addNums[i]), -3))
		fv = fv.Div(decimal.New(int64(divNums[i]), -3))
		fv = fv.Mul(decimal.New(int64(mulNums[i]), -3))
		fv = fv.Sub(decimal.New(int64(subNums[i]), -3))
	}
}
