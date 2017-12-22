package fraction

import (
	"testing"

	asrt "github.com/stretchr/testify/assert"
)

//Mul(Fraction) Fraction
//Div(Fraction) Fraction
//Add(Fraction) Fraction
//Sub(Fraction) Fraction
//Evaluate() int64

func TestNew(t *testing.T) {
	assert := asrt.New(t)

	assert.Equal(New(1), New(1, 1))
	assert.Equal(New(100), New(100, 1))
	assert.Equal(New(-1), New(-1, 1))
	assert.Equal(New(-100), New(-100, 1))
	assert.Equal(New(0), New(0, 1))

	// do not allow for more than 2 variables
	assert.Panics(func() { New(1, 1, 1) })
}

func TestNegativePositive(t *testing.T) {
	assert := asrt.New(t)

	f1 := New(100, 1)
	f2 := New(-100, -1)
	f3 := New(100, -1)
	f4 := New(-100, 1)

	assert.True(f1.Positive())
	assert.False(f1.Negative())

	assert.True(f2.Positive())
	assert.False(f2.Negative())

	assert.False(f3.Positive())
	assert.True(f3.Negative())

	assert.False(f4.Positive())
	assert.True(f4.Negative())
}

func TestSimplify(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		start, simplified Fraction
	}{
		{New(1), New(1)},
		{New(-100), New(-100)},
		{New(100, 100), New(1)},
		{New(10000000, 10000000), New(1)},
		{New(-10000000, 10000000), New(-1)},
		{New(-10000000, -10000000), New(1)},
		{New(7, 13), New(7, 13)},
		{New(100, 13), New(100, 13)},
		{New(4, 2), New(2, 1)},
		{New(69, 3), New(23)},
		{New(333, 106), New(333, 106)}, //pi :)
		{New(10000000, 2), New(5000000)},
	}

	for _, test := range tests {
		assert.True(test.simplified.Equal(test.start.Simplify()))
		assert.True(test.simplified.Inv().Equal(test.start.Inv().Simplify()))
	}
}

func TestEqualities(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		f1, f2     Fraction
		gt, lt, eq bool
	}{
		{New(0), New(0), false, false, true},
		{New(0, 100), New(0, 10000), false, false, true},
		{New(100), New(100), false, false, true},
		{New(-100), New(-100), false, false, true},
		{New(-100, -1), New(100), false, false, true},
		{New(-1, 1), New(1, -1), false, false, true},
		{New(1, -1), New(-1, 1), false, false, true},
		{New(3, 7), New(3, 7), false, false, true},

		{New(0), New(3, 7), false, true, false},
		{New(0), New(100), false, true, false},
		{New(-1), New(3, 7), false, true, false},
		{New(-1), New(100), false, true, false},
		{New(1, 7), New(100), false, true, false},
		{New(1, 7), New(3, 7), false, true, false},
		{New(-3, 7), New(-1, 7), false, true, false},

		{New(3, 7), New(0), true, false, false},
		{New(100), New(0), true, false, false},
		{New(3, 7), New(-1), true, false, false},
		{New(100), New(-1), true, false, false},
		{New(100), New(1, 7), true, false, false},
		{New(3, 7), New(1, 7), true, false, false},
		{New(-1, 7), New(-3, 7), true, false, false},
	}

	for _, test := range tests {
		assert.Equal(test.gt, test.f1.GT(test.f2))
		assert.Equal(test.lt, test.f1.LT(test.f2))
		assert.Equal(test.eq, test.f1.Equal(test.f2))
	}

}

func TestArithmatic(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		f1, f2                         Fraction
		resMul, resDiv, resAdd, resSub Fraction
	}{
		// f1    f2      MUL     DIV     ADD     SUB
		{New(0), New(0), New(0), New(0), New(0), New(0)},
		{New(1), New(0), New(0), New(0), New(1), New(1)},
		{New(0), New(1), New(0), New(0), New(1), New(-1)},
		{New(0), New(-1), New(0), New(0), New(-1), New(1)},
		{New(-1), New(0), New(0), New(0), New(-1), New(-1)},

		{New(1), New(1), New(1), New(1), New(2), New(0)},
		{New(-1), New(-1), New(1), New(1), New(-2), New(0)},
		{New(1), New(-1), New(-1), New(-1), New(0), New(2)},
		{New(-1), New(1), New(-1), New(-1), New(0), New(-2)},

		{New(3), New(7), New(21), New(3, 7), New(10), New(-4)},
		{New(2), New(4), New(8), New(1, 2), New(6), New(-2)},
		{New(100), New(100), New(10000), New(1), New(200), New(0)},

		{New(3, 2), New(3, 2), New(9, 4), New(1), New(3), New(0)},
		{New(3, 7), New(7, 3), New(1), New(9, 49), New(58, 21), New(-40, 21)},
		{New(1, 21), New(11, 5), New(11, 105), New(5, 231), New(236, 105), New(-226, 105)},
		{New(-21), New(3, 7), New(-9), New(-49), New(-144, 7), New(-150, 7)},
		{New(100), New(1, 7), New(100, 7), New(700), New(701, 7), New(699, 7)},
	}

	for _, test := range tests {
		assert.Equal(test.resMul, test.f1.Mul(test.f2), "f1 %v, f2 %v", test.f1, test.f2)
		assert.Equal(test.resAdd, test.f1.Add(test.f2), "f1 %v, f2 %v", test.f1, test.f2)
		assert.Equal(test.resSub, test.f1.Sub(test.f2), "f1 %v, f2 %v", test.f1, test.f2)

		if test.f2.GetNumerator() == 0 { // panic for divide by zero
			assert.Panics(func() { test.f1.Div(test.f2) })
		} else {
			assert.Equal(test.resDiv, test.f1.Div(test.f2), "f1 %v, f2 %v", test.f1, test.f2)
		}
	}
}

//func TestEvaluate(t *testing.T) {
//assert := asrt.New(t)
//}
