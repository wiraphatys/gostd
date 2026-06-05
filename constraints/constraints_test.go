package constraints

import "testing"

func TestLessGreater(t *testing.T) {
	if !Less(1, 2) {
		t.Error("Less(1, 2) should be true")
	}
	if Less(2, 1) {
		t.Error("Less(2, 1) should be false")
	}
	if Less(1, 1) {
		t.Error("Less(1, 1) should be false (strict)")
	}
	if !Greater(2, 1) {
		t.Error("Greater(2, 1) should be true")
	}
	if Greater(1, 2) {
		t.Error("Greater(1, 2) should be false")
	}

	// Strings are ordered too.
	if !Less("apple", "banana") {
		t.Error("Less(\"apple\", \"banana\") should be true")
	}
}

func TestMinMax(t *testing.T) {
	if got := Min(3, 7); got != 3 {
		t.Errorf("Min(3, 7) = %d, want 3", got)
	}
	if got := Max(3, 7); got != 7 {
		t.Errorf("Max(3, 7) = %d, want 7", got)
	}
	if got := Min(7, 3); got != 3 {
		t.Errorf("Min(7, 3) = %d, want 3", got)
	}
	if got := Max(7, 3); got != 7 {
		t.Errorf("Max(7, 3) = %d, want 7", got)
	}
	if got := Min(5.5, 5.5); got != 5.5 {
		t.Errorf("Min(5.5, 5.5) = %v, want 5.5", got)
	}
	if got := Max("a", "z"); got != "z" {
		t.Errorf("Max(\"a\", \"z\") = %q, want \"z\"", got)
	}
}

// Compile-time checks that the named numeric types satisfy their constraints.
func acceptSigned[T Signed](v T) T     { return v }
func acceptUnsigned[T Unsigned](v T) T { return v }
func acceptFloat[T Float](v T) T       { return v }
func acceptNumber[T Number](v T) T     { return v }
func acceptOrdered[T Ordered](v T) T   { return v }
func acceptComplex[T Complex](v T) T   { return v }
func acceptInteger[T Integer](v T) T   { return v }

func TestConstraintMembership(t *testing.T) {
	if acceptSigned(int8(-1)) != -1 {
		t.Error("int8 should satisfy Signed")
	}
	if acceptUnsigned(uint16(1)) != 1 {
		t.Error("uint16 should satisfy Unsigned")
	}
	if acceptFloat(float32(1.5)) != 1.5 {
		t.Error("float32 should satisfy Float")
	}
	if acceptInteger(int64(9)) != 9 {
		t.Error("int64 should satisfy Integer")
	}
	if acceptNumber(3.0) != 3.0 {
		t.Error("float64 should satisfy Number")
	}
	if acceptOrdered("x") != "x" {
		t.Error("string should satisfy Ordered")
	}
	if acceptComplex(complex(1, 2)) != complex(1, 2) {
		t.Error("complex128 should satisfy Complex")
	}
}
