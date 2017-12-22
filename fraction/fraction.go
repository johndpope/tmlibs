package fraction

import (
	wire "github.com/tendermint/go-wire"
)

// XXX test fractions!

// Fraction -  basic fraction functionality
// TODO better name that Fraction?
type Fraction interface {
	GetNumerator() int64
	GetDenominator() int64
	RectifySign() Fraction
	Inv() Fraction
	Simplify() Fraction
	Negative() bool
	Positive() bool
	GT(Fraction) bool
	LT(Fraction) bool
	Equal(Fraction) bool
	Mul(Fraction) Fraction
	Div(Fraction) Fraction
	Add(Fraction) Fraction
	Sub(Fraction) Fraction
	Evaluate() int64
}

// fraction - basic fraction
type fraction struct {
	Numerator, Denominator int64
}

var _ Fraction = fraction{} // enforce at compile time
var _ = wire.RegisterInterface(struct{ Fraction }{}, wire.ConcreteType{fraction{}, 0x01})

// New - create a new fraction object
func New(Numerator int64, Denominator ...int64) Fraction {
	switch len(Denominator) {
	case 0:
		return fraction{Numerator, 1}
	case 1:
		return fraction{Numerator, Denominator[0]}
	default:
		panic("improper use of NewFraction, can only have one denominator")
	}
}

// GetNumerator - return the Numerator
func (f fraction) GetNumerator() int64 {
	return f.Numerator
}

// GetDenominator - return the Denominator
func (f fraction) GetDenominator() int64 {
	return f.Denominator
}

// RectifySign - make any negative sign exist only in Numerator
func (f fraction) RectifySign() Fraction {
	if f.Denominator < 0 {
		return fraction{-1 * f.Numerator, -1 * f.Denominator}
	}
	return f
}

// Inv - Inverse
func (f fraction) Inv() Fraction {
	return fraction{f.Denominator, f.Numerator}
}

// Simplify - find the greatest common Denominator, divide
func (f fraction) Simplify() Fraction {
	gcd := f.Numerator
	for d := f.Denominator; d != 0; {
		gcd, d = d, gcd%d
	}
	f.Numerator /= gcd
	f.Denominator /= gcd
	return f.RectifySign()
}

// Negative - is the fractior negative
func (f fraction) Negative() bool {
	switch {
	case f.Numerator > 0:
		if f.Denominator > 0 {
			return false
		}
		return true
	case f.Numerator < 0:
		if f.Denominator < 0 {
			return false
		}
		return true
	}
	return false
}

// Positive - is the fraction positive
func (f fraction) Positive() bool {
	switch {
	case f.Numerator > 0:
		if f.Denominator > 0 {
			return true
		}
		return false
	case f.Numerator < 0:
		if f.Denominator < 0 {
			return true
		}
		return false
	}
	return false
}

// Equal - test if two Fractions are equal, does not simplify
func (f fraction) Equal(f2 Fraction) bool {
	if f.Numerator == 0 {
		return f2.GetNumerator() == 0
	}
	f1 := f.RectifySign()
	f2 = f2.RectifySign()
	return ((f1.GetNumerator() == f2.GetNumerator()) &&
		(f1.GetDenominator() == f2.GetDenominator()))
}

// GT - greater than
func (f fraction) GT(f2 Fraction) bool {
	return f.Sub(f2).Positive()
}

// LT - less than
func (f fraction) LT(f2 Fraction) bool {
	return f.Sub(f2).Negative()
}

// Mul - multiply
func (f fraction) Mul(f2 Fraction) Fraction {
	return fraction{
		f.Numerator * f2.GetNumerator(),
		f.Denominator * f2.GetDenominator(),
	}.Simplify()
}

// Div - divide
func (f fraction) Div(f2 Fraction) Fraction {
	if f2.GetNumerator() == 0 {
		panic("fraction divide by zero error!")
	}
	return fraction{
		f.Numerator * f2.GetDenominator(),
		f.Denominator * f2.GetNumerator(),
	}.Simplify()
}

// Add - add without simplication
func (f fraction) Add(f2 Fraction) Fraction {
	if f.Denominator == f2.GetDenominator() {
		return fraction{
			f.Numerator + f2.GetNumerator(),
			f.Denominator,
		}.Simplify()
	}
	return fraction{
		f.Numerator*f2.GetDenominator() + f2.GetNumerator()*f.Denominator,
		f.Denominator * f2.GetDenominator(),
	}.Simplify()
}

// Sub - subtract without simplication
func (f fraction) Sub(f2 Fraction) Fraction {
	if f.Denominator == f2.GetDenominator() {
		return fraction{
			f.Numerator - f2.GetNumerator(),
			f.Denominator,
		}.Simplify()
	}
	return fraction{
		f.Numerator*f2.GetDenominator() - f2.GetNumerator()*f.Denominator,
		f.Denominator * f2.GetDenominator(),
	}.Simplify()
}

// Evaluate - evaluate the fraction using bankers rounding
func (f fraction) Evaluate() int64 {

	d := f.Numerator / f.Denominator // always drops the decimal
	if f.Numerator%f.Denominator == 0 {
		return d
	}

	// evaluate the remainder using bankers rounding
	remainderDigit := (f.Numerator * 10 / f.Denominator) - (d * 10) // get the first remainder digit

	isFinalDigit := (f.Numerator*10%f.Denominator == 0) // is this the final digit in the remainder?
	if isFinalDigit && (remainderDigit == 5 || remainderDigit == -5) {
		return d + (d % 2) // always rounds to the even number
	}
	if remainderDigit >= 5 {
		d++
	}
	return d
}
