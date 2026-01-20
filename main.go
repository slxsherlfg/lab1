package main

import (
	"fmt"
	"math"
)

const PI = math.Pi

// ANGLE

type Angle struct {
	rad float64
}

func normalizeRad(x float64) float64 {
	t := math.Mod(x, 2*PI)
	if t < 0 {
		t += 2 * PI
	}
	return t
}

func AngleRad(r float64) Angle {
	return Angle{rad: normalizeRad(r)}
}

func AngleDeg(d float64) Angle {
	return AngleRad(d * PI / 180)
}

func (a *Angle) SetRad(r float64) {
	a.rad = normalizeRad(r)
}

func (a *Angle) SetDeg(d float64) {
	a.rad = normalizeRad(d * PI / 180)
}

func (a Angle) Rad() float64 {
	return a.rad
}

func (a Angle) Deg() float64 {
	return a.rad * 180 / PI
}

func (a Angle) Float64() float64 {
	return a.rad
}

func (a Angle) Int() int {
	return int(a.rad)
}

func (a Angle) String() string {
	return fmt.Sprintf("%.6f rad (%.6f°)", a.rad, a.Deg())
}

func (a Angle) Repr() string {
	return fmt.Sprintf("Angle(rad=%.6f)", a.rad)
}

func (a Angle) Equal(b Angle) bool {
	return math.Abs(normalizeRad(a.rad)-normalizeRad(b.rad)) < 1e-9
}

func (a Angle) Add(x interface{}) Angle {
	switch v := x.(type) {
	case Angle:
		return AngleRad(a.rad + v.rad)
	case float64:
		return AngleRad(a.rad + v)
	case int:
		return AngleRad(a.rad + float64(v))
	default:
		panic("unsupported type in Add")
	}
}

func (a Angle) Sub(x interface{}) Angle {
	switch v := x.(type) {
	case Angle:
		return AngleRad(a.rad - v.rad)
	case float64:
		return AngleRad(a.rad - v)
	case int:
		return AngleRad(a.rad - float64(v))
	default:
		panic("unsupported type in Sub")
	}
}

func (a Angle) Mul(k float64) Angle {
	return AngleRad(a.rad * k)
}

func (a Angle) Div(k float64) Angle {
	return AngleRad(a.rad / k)
}

//ANGLE RANGE

type AngleRange struct {
	Start, End   Angle
	IncludeStart bool
	IncludeEnd   bool
}

func NewAngleRange(start, end interface{}, includeStart, includeEnd bool) AngleRange {
	return AngleRange{
		Start:        convertToAngle(start),
		End:          convertToAngle(end),
		IncludeStart: includeStart,
		IncludeEnd:   includeEnd,
	}
}

func convertToAngle(v interface{}) Angle {
	switch x := v.(type) {
	case Angle:
		return x
	case float64:
		return AngleRad(x)
	case int:
		return AngleRad(float64(x))
	default:
		panic("неверный тип")
	}
}

func (r AngleRange) String() string {
	lb := "("
	rb := ")"
	if r.IncludeStart {
		lb = "["
	}
	if r.IncludeEnd {
		rb = "]"
	}
	return fmt.Sprintf("%s%s, %s%s", lb, r.Start.String(), r.End.String(), rb)
}

func (r AngleRange) Repr() string {
	return fmt.Sprintf("AngleRange(%v, %v, incStart=%v, incEnd=%v)",
		r.Start.Repr(), r.End.Repr(), r.IncludeStart, r.IncludeEnd)
}

func (r AngleRange) Len() float64 {
	d := normalizeRad(r.End.Rad() - r.Start.Rad())
	return d
}

func (r AngleRange) Equal(b AngleRange) bool {
	return r.Start.Equal(b.Start) &&
		r.End.Equal(b.End) &&
		r.IncludeStart == b.IncludeStart &&
		r.IncludeEnd == b.IncludeEnd
}

// проверка in
func (r AngleRange) ContainsAngle(a Angle) bool {
	s := r.Start.Rad()
	e := r.End.Rad()
	x := normalizeRad(a.Rad())

	if s < e {
		if x < s || x > e {
			return false
		}
	}
	if x == s && !r.IncludeStart {
		return false
	}
	if x == e && !r.IncludeEnd {
		return false
	}
	return true
}

func (r AngleRange) ContainsRange(b AngleRange) bool {
	return r.ContainsAngle(b.Start) && r.ContainsAngle(b.End)
}

func (r AngleRange) Add(x Angle) []AngleRange {
	return []AngleRange{
		NewAngleRange(r.Start.Add(x), r.End.Add(x), r.IncludeStart, r.IncludeEnd),
	}
}

func (r AngleRange) Sub(x Angle) []AngleRange {
	return []AngleRange{
		NewAngleRange(r.Start.Sub(x), r.End.Sub(x), r.IncludeStart, r.IncludeEnd),
	}
}

//  демонстрация работы

func main() {
	fmt.Println("ANGLE")

	a := AngleRad(3 * PI)
	b := AngleDeg(180)

	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("a == PI?", a.Equal(AngleRad(PI)))

	fmt.Println("a + 1:", a.Add(1.0))
	fmt.Println("b - 90°:", b.Sub(AngleDeg(90)))
	fmt.Println("b * 2:", b.Mul(2))
	fmt.Println("b / 2:", b.Div(2))

	fmt.Println("ANGLE RANGE")

	r := NewAngleRange(0, PI, true, false)
	fmt.Println("r:", r)
	fmt.Println("len(r):", r.Len())
	fmt.Println("contains 1 rad?", r.ContainsAngle(AngleRad(1)))
	fmt.Println("contains 3 rad?", r.ContainsAngle(AngleRad(3)))

	r2 := NewAngleRange(AngleDeg(30), AngleDeg(150), true, true)
	fmt.Println("r2 inside r?", r.ContainsRange(r2))

	res := r.Add(AngleDeg(45))
	fmt.Println("r + 45°:", res)

	x := Angle{rad: 3.141592}
	fmt.Println(r.Sub(x))
}
