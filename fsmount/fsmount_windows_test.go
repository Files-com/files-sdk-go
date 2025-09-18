package fsmount

import "testing"

func TestValidMountPoint(t *testing.T) {
	tests := []struct {
		mountPoint string
		valid      bool
	}{
		{"", true},
		{"C:", true},
		{"Z:", true},
		{"A:", true},
		{"1:", false},
		{"C", false},
		{"CC:", false},
		{"C::", false},
		{"C:/", false},
		{"C:\\", false},
		{"C:extra", false},
	}

	for _, test := range tests {
		err := validateMountPoint(test.mountPoint)
		if (err == nil) != test.valid {
			t.Errorf("validateMountPoint(%q) = %v; want valid=%v", test.mountPoint, err, test.valid)
		}
	}
}
