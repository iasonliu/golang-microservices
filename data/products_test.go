package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "nics",
		Price: 1.00,
		SKU:   "adb-xxds-dfasdf",
	}
	v := NewValidation()
	err := v.Validate(p)

	if err != nil {
		t.Fatal(err)
	}
}
