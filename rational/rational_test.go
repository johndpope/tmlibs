package rational

import (
	"encoding/json"
	"testing"

	asrt "github.com/stretchr/testify/assert"
	rqr "github.com/stretchr/testify/require"
)

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

func TestNewFromDecimal(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		decimalStr string
		expErr     bool
		exp        Rational
	}{
		{"0", false, New(0)},
		{"1", false, New(1)},
		{"1.1", false, New(11, 10)},
		{"0.75", false, New(3, 4)},
		{"0.8", false, New(4, 5)},
		{"0.11111", false, New(11111, 100000)},
		{".", true, Rat{}},
		{".0", true, Rat{}},
		{"1.", true, Rat{}},
		{"foobar", true, Rat{}},
		{"0.foobar", true, Rat{}},
		{"0.foobar.", true, Rat{}},
	}

	for _, test := range tests {

		res, err := NewFromDecimal(test.decimalStr)
		if test.expErr {
			assert.NotNil(err, test.decimalStr)
		} else {
			assert.Nil(err)
			assert.True(res.Equal(test.exp))
		}

		// negative test
		res, err = NewFromDecimal("-" + test.decimalStr)
		if test.expErr {
			assert.NotNil(err, test.decimalStr)
		} else {
			assert.Nil(err)
			assert.True(res.Equal(test.exp.Mul(New(-1))))
		}
	}
}

func TestEqualities(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		r1, r2     Rational
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
		assert.Equal(test.gt, test.r1.GT(test.r2))
		assert.Equal(test.lt, test.r1.LT(test.r2))
		assert.Equal(test.eq, test.r1.Equal(test.r2))
	}

}

func TestArithmatic(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		r1, r2                         Rational
		resMul, resDiv, resAdd, resSub Rational
	}{
		// r1    r2      MUL     DIV     ADD     SUB
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
		assert.True(test.resMul.Equal(test.r1.Mul(test.r2)), "r1 %v, r2 %v", test.r1.GetRat(), test.r2.GetRat())
		assert.True(test.resAdd.Equal(test.r1.Add(test.r2)), "r1 %v, r2 %v", test.r1.GetRat(), test.r2.GetRat())
		assert.True(test.resSub.Equal(test.r1.Sub(test.r2)), "r1 %v, r2 %v", test.r1.GetRat(), test.r2.GetRat())

		if test.r2.Num() == 0 { // panic for divide by zero
			assert.Panics(func() { test.r1.Quo(test.r2) })
		} else {
			assert.True(test.resDiv.Equal(test.r1.Quo(test.r2)), "r1 %v, r2 %v", test.r1.GetRat(), test.r2.GetRat())
		}
	}
}

func TestEvaluate(t *testing.T) {
	assert := asrt.New(t)

	tests := []struct {
		r1  Rational
		res int64
	}{
		{New(0), 0},
		{New(1), 1},
		{New(1, 4), 0},
		{New(1, 2), 0},
		{New(3, 4), 1},
		{New(5, 6), 1},
		{New(3, 2), 2},
		{New(5, 2), 2},
		{New(6, 11), 1},  // 0.545-> 1 even though 5 is first decimal and 1 not even
		{New(17, 11), 2}, // 1.545
		{New(5, 11), 0},
		{New(16, 11), 1},
		{New(113, 12), 9},
	}

	for _, test := range tests {
		assert.Equal(test.res, test.r1.Evaluate(), "%v", test.r1)
		assert.Equal(test.res*-1, test.r1.Mul(New(-1)).Evaluate(), "%v", test.r1.Mul(New(-1)))
	}
}

func TestSerialization(t *testing.T) {
	assert, require := asrt.New(t), rqr.New(t)

	r := New(1, 3)

	rMarshal, err := json.Marshal(r)
	require.Nil(err)

	var rUnmarshal Rat
	err = json.Unmarshal(rMarshal, &rUnmarshal)
	require.Nil(err)

	//panic(fmt.Sprintf("debug rUnmarshal: %v\n", rUnmarshal))
	assert.True(r.Equal(rUnmarshal), "original: %v, unmarshalled: %v", r, rUnmarshal)
}
