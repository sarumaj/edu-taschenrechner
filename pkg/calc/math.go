package calc

import (
	"fmt"
	"math/big"
)

// Factorial calculates the factorial of a number n using the formula n! = n * (n-1) * (n-2) * ... * 1
// The step parameter is used to calculate the factorial in steps of step numbers at a time, i.e.,
// n! = n * (n-step) * (n-2*step) * ...
func Factorial(n *big.Float, step int) (*big.Float, error) {
	zero := big.NewFloat(0)
	one := big.NewFloat(1)

	if n.Cmp(zero) < 0 {
		return nil, fmt.Errorf("factorial of negative number is not defined")
	}

	if n.Cmp(zero) == 0 {
		return one, nil
	}

	stepInt := big.NewInt(int64(step))
	if intN, accuracy := n.Int(nil); accuracy == big.Exact {
		result := big.NewInt(1)
		for i := new(big.Int).Set(intN); i.Cmp(big.NewInt(0)) > 0; i.Sub(i, stepInt) {
			result.Mul(result, i)
		}

		return big.NewFloat(0).SetInt(result), nil
	}

	return nil, fmt.Errorf("factorial is only defined for integers")
}

// Pow calculates base^exp for big.Float values using exp(ln(base) * exp)
func Pow(base, exponent *big.Float) (*big.Float, error) {
	if base.Cmp(big.NewFloat(0)) == 0 && exponent.Cmp(big.NewFloat(0)) == 0 {
		return nil, fmt.Errorf("0^0 is undefined")
	}

	// Handle simple cases
	zero := big.NewFloat(0)
	one := big.NewFloat(1)

	if exponent.Cmp(zero) == 0 {
		return one, nil
	}

	if exponent.Cmp(one) == 0 {
		return base, nil
	}

	// Handle integer exponents directly
	oneInt := big.NewInt(1)
	if intExp, accuracy := exponent.Int(nil); accuracy == big.Exact {
		result := big.NewFloat(1)
		for i := big.NewInt(0); i.Cmp(big.NewInt(0).Abs(intExp)) < 0; i.Add(i, oneInt) {
			result = result.Mul(result, base)
		}

		if intExp.Sign() < 0 { // Negative exponent: take reciprocal
			result = result.Quo(one, result)
		}

		return result, nil
	}

	return nil, fmt.Errorf("non-integer exponents are not supported")
}
