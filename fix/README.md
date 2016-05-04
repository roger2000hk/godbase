# fix
#### a speedy, comparable fixed-point for Go

### why?
I needed a fixed-point type, but couldn't find one that felt just right. Further investigation showed that big.Ints are not comparable which means using them this deep down the stack is an extra responsibility and limitation that I'm not really comfortable with.

### what?
I ended up implementing a value struct using int64 numerator and denominators. This means that the upper representable limit is max(int64)/denominator. The api resembles the big apis, with pointer receivers for results. 

### performance
I implemented a basic benchmark loop that adds, subs, muls & divs random numbers; for fix, github.com/oguzbilgic/fpd and github.com/shopspring/decimal. The Scale flavor rescales each value. I would like to increase the number of iterations further to get a valid comparison, but both fpd and decimal start running into weird edge cases and timeouts.

```
BenchmarkFix-4                  2000000000               0.00 ns/op
BenchmarkFixScale-4             2000000000               0.00 ns/op
BenchmarkFpd-4                  2000000000               0.09 ns/op
BenchmarkFpdScale-4             2000000000               0.12 ns/op
BenchmarkDecimal-4              2000000000               0.34 ns/op
BenchmarkDecimalScale-4         2000000000               0.34 ns/op
```

### status
fix is one of the column types supported by godbase. At the moment it only implements enough functionality to prove my point that this is a good idea. Add/Div/Mul/Sub and basic scaling, conversion and printing is supported.