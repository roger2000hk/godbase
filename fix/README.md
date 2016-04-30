# fix
### speedy fixed-point math for Go

## why?
I needed a fixed-point type, but couldn't find one with a sane api and without weird edge cases.

## status
fix is one of the column types supported by godbase. At the moment it only implements enough functionality to prove my point that this is the right way to go. Add/Div/Mul/Sub and basic scaling, conversion and printing is supported.

## performance
I implemented a basic benchmark loop that adds, subs, muls & divs random numbers; for fix, github.com/oguzbilgic/fpd and github.com/shopspring/decimal. The Scale flavor rescales each value. I would like to go further, but both fpd and decimal start running into weird edge cases when the number of iters is increased.

```
BenchmarkFix-4                  2000000000               0.02 ns/op
BenchmarkFixScale-4             2000000000               0.03 ns/op
BenchmarkFpd-4                  2000000000               0.10 ns/op
BenchmarkFpdScale-4             2000000000               0.14 ns/op
BenchmarkDecimal-4              2000000000               0.33 ns/op
BenchmarkDecimalScale-4         2000000000               0.34 ns/op
```