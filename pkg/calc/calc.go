/*
Package calc provides functions for performing mathematical calculations.
*/

package calc

import (
	"context"
	"fmt"
	"math/big"
)

// Factorial calculates the factorial of a number n using the formula n! = n * (n-1) * (n-2) * ... * 1
// The step parameter is used to calculate the factorial in steps of step numbers at a time, i.e.,
// n! = n * (n-step) * (n-2*step) * ...
func Factorial(ctx context.Context, n *big.Float, step int) (*big.Float, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		result := big.NewInt(1)
		for i := new(big.Int).Set(intN); i.Cmp(big.NewInt(0)) > 0; i.Sub(i, stepInt) {
			result.Mul(result, i)
		}

		return big.NewFloat(0).SetInt(result), nil
	}

	return nil, fmt.Errorf("factorial is only defined for integers")
}

// GreatestCommonDivisor calculates the greatest common divisor of two numbers x and y using the Euclidean algorithm
func GreatestCommonDivisor(ctx context.Context, x, y *big.Float) (*big.Float, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	xInt, accuracy := x.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("GCD function requires integer arguments")
	}

	yInt, accuracy := y.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("GCD function requires integer arguments")
	}

	return big.NewFloat(0).SetInt(new(big.Int).GCD(nil, nil, xInt, yInt)), nil
}

// LeastCommonMultiple calculates the least common multiple of two numbers x and y using the formula LCM(x, y) = x * y / GCD(x, y)
func LeastCommonMultiple(ctx context.Context, x, y *big.Float) (*big.Float, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	xInt, accuracy := x.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("LCM function requires integer arguments")
	}

	yInt, accuracy := y.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("LCM function requires integer arguments")
	}

	gcd := new(big.Int).GCD(nil, nil, xInt, yInt)
	if gcd.Cmp(big.NewInt(0)) == 0 {
		return big.NewFloat(0), nil
	}

	// Calculate LCM(x, y) = abs(x * y) / GCD(x, y)
	lcm := big.NewInt(0).Mul(xInt, yInt)
	lcm = lcm.Quo(lcm, gcd)
	lcm = lcm.Abs(lcm) // Ensure LCM is positive

	return big.NewFloat(0).SetInt(lcm), nil
}

// Pow calculates base^exp for big.Float values using exp(ln(base) * exp)
func Pow(ctx context.Context, base, exponent *big.Float) (*big.Float, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

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
		if err := ctx.Err(); err != nil {
			return nil, err
		}

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
