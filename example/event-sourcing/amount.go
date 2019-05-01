package main

type Amount int64

func (a Amount) Subtract(other Amount) Amount {
	return a - other
}

func (a Amount) IsPositive() bool {
	return a > 0
}

func (a Amount) LargerThanEq(other Amount) bool {
	return a >= other
}
