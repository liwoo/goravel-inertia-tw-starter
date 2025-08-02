package tests

// Minimal type-state pattern example
type Builder[State any] struct {
	value string
}

// States
type needsA struct{}
type needsB struct{}
type canBuild struct{}

func NewBuilder() *Builder[needsA] {
	return &Builder[needsA]{}
}

func (b *Builder[needsA]) SetA(v string) *Builder[needsB] {
	return &Builder[needsB]{value: v}
}

func (b *Builder[needsB]) SetB(v string) *Builder[canBuild] {
	return &Builder[canBuild]{value: b.value + v}
}

func (b *Builder[canBuild]) Build() string {
	return b.value
}

// Test function
func TestMinimalTypeState() {
	// This should compile
	result := NewBuilder().SetA("a").SetB("b").Build()
	_ = result
	
	// These should NOT compile if uncommented:
	// NewBuilder().Build() // Error: Build undefined
	// NewBuilder().SetB("b") // Error: SetB undefined
	// NewBuilder().SetA("a").Build() // Error: Build undefined
}