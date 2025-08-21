package tests

func TestMinimalBroken() {
	// This should NOT compile - calling Build too early
	result := NewBuilder().Build()
	_ = result
}
