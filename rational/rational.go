package rational

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// Rational - big Rat with additional functionality
type Rational interface {
	GetRat() *big.Rat
	Num() int64
	Denom() int64
	GT(Rational) bool
	LT(Rational) bool
	Equal(Rational) bool
	IsZero() bool
	Inv() Rational
	Mul(Rational) Rational
	Quo(Rational) Rational
	Add(Rational) Rational
	Sub(Rational) Rational
	Evaluate() int64
}

// Rat - extend big.Rat
type Rat struct {
	*big.Rat `json:"rat"`
}

var _ Rational = Rat{} // enforce at compile time

// XXX Is this even possible right now? or are we stuck using the struct?
//var _ = wire.RegisterInterface(struct{ Rational }{}, wire.ConcreteType{Rat{}, 0x01})

// New - create a new Rat from integers
func New(Numerator int64, Denominator ...int64) Rational {
	switch len(Denominator) {
	case 0:
		return Rat{big.NewRat(Numerator, 1)}
	case 1:
		return Rat{big.NewRat(Numerator, Denominator[0])}
	default:
		panic("improper use of New, can only have one denominator")
	}
}

//NewFromDecimal - create a rational from decimal string or integer string
func NewFromDecimal(decimalStr string) (f Rational, err error) {

	// first extract any negative symbol
	neg := false
	if string(decimalStr[0]) == "-" {
		neg = true
		decimalStr = decimalStr[1:]
	}

	str := strings.Split(decimalStr, ".")

	var numStr string
	var denom int64 = 1
	switch len(str) {
	case 1:
		if len(str[0]) == 0 {
			return f, fmt.Errorf("not a decimal string")
		}
		numStr = str[0]
	case 2:
		if len(str[0]) == 0 || len(str[1]) == 0 {
			return f, fmt.Errorf("not a decimal string")
		}
		numStr = str[0] + str[1]
		len := int64(len(str[1]))
		denom = new(big.Int).Exp(big.NewInt(10), big.NewInt(len), nil).Int64()
	default:
		return f, fmt.Errorf("not a decimal string")
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return f, err
	}

	if neg {
		num *= -1
	}

	return Rat{big.NewRat(int64(num), denom)}, nil
}

//nolint
func (r Rat) GetRat() *big.Rat         { return r.Rat }                                     // GetRat - get big.Rational
func (r Rat) Num() int64               { return r.Rat.Num().Int64() }                       // Num - return the numerator
func (r Rat) Denom() int64             { return r.Rat.Denom().Int64() }                     // Denom  - return the denominator
func (r Rat) IsZero() bool             { return r.Num() == 0 }                              // IsZero - Is the Rat equal to zero
func (r Rat) Equal(r2 Rational) bool   { return r.Rat.Cmp(r2.GetRat()) == 0 }               // Equal - rationals are equal
func (r Rat) GT(r2 Rational) bool      { return r.Rat.Cmp(r2.GetRat()) == 1 }               // GT - greater than
func (r Rat) LT(r2 Rational) bool      { return r.Rat.Cmp(r2.GetRat()) == -1 }              // LT - less than
func (r Rat) Inv() Rational            { return Rat{new(big.Rat).Inv(r.Rat)} }              // Inv - inverse
func (r Rat) Mul(r2 Rational) Rational { return Rat{new(big.Rat).Mul(r.Rat, r2.GetRat())} } // Mul - multiplication
func (r Rat) Quo(r2 Rational) Rational { return Rat{new(big.Rat).Quo(r.Rat, r2.GetRat())} } // Quo - quotient
func (r Rat) Add(r2 Rational) Rational { return Rat{new(big.Rat).Add(r.Rat, r2.GetRat())} } // Add - addition
func (r Rat) Sub(r2 Rational) Rational { return Rat{new(big.Rat).Sub(r.Rat, r2.GetRat())} } // Sub - subtraction

// Evaluate - evaluate the rational using bankers rounding
func (r Rat) Evaluate() int64 {

	num := r.Num()
	denom := r.Denom()

	d := num / denom // always drops the decimal
	if num%denom == 0 {
		return d
	}

	// evaluate the remainder using bankers rounding
	remainderDigit := (num * 10 / denom) - (d * 10) // get the first remainder digit
	isFinalDigit := (num*10%denom == 0)             // is this the final digit in the remainder?

	switch {
	case isFinalDigit && (remainderDigit == 5 || remainderDigit == -5):
		return d + (d % 2) // always rounds to the even number
	case remainderDigit >= 5:
		d++
	case remainderDigit <= -5:
		d--
	}
	return d
}

//___________________________________________________________________________________

// RatMarshal - Marshable Rational Struct
type RatMarshal struct {
	Numerator, Denominator int64
}

// MarshalJSON - custom implementation of JSON Marshal
func (r Rat) MarshalJSON() ([]byte, error) {
	return json.Marshal(RatMarshal{r.Num(), r.Denom()})
}

// UnmarshalJSON - custom implementation of JSON Unmarshal
func (r *Rat) UnmarshalJSON(data []byte) error {
	ratMar := new(RatMarshal)
	err := json.Unmarshal(data, ratMar)
	if err != nil {
		return err
	}
	r.Rat = big.NewRat(ratMar.Numerator, ratMar.Denominator)
	return nil
}
