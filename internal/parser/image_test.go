package parser

import "testing"

func TestImageName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"docker.io/library/alpine",
			"library/alpine",
		},
		{
			"ghcr.io/sagernet/sing-box",
			"sagernet/sing-box",
		},
		{
			"alpine",
			"library/alpine",
		},
	}

	for _, tt := range tests {
		got := ImageName(tt.input)

		if got != tt.want {
			t.Errorf("ImageName(%q) = %q, want %q",
				tt.input,
				got,
				tt.want,
			)
		}
	}
}
