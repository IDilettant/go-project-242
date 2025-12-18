package size

import "testing"

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name  string
		size  int64
		human bool
		want  string
	}{
		{"raw bytes", 123, false, "123B"},
		{"raw zero", 0, false, "0B"},

		{"human bytes", 123, true, "123B"},
		{"human zero", 0, true, "0B"},
		{"human below kb", 1023, true, "1023B"},

		{"exact kb", 1024, true, "1.0KB"},
		{"exact mb", 1024 * 1024, true, "1.0MB"},
		{"exact gb", 1024 * 1024 * 1024, true, "1.0GB"},

		{"fractional kb", 1536, true, "1.5KB"},
		{"fractional mb", 1234567, true, "1.2MB"},

		{"clamp to tb", 1 << 50, true, "1.0PB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSize(tt.size, tt.human); got != tt.want {
				t.Fatalf("FormatSize(%d, %v) = %q, want %q",
					tt.size, tt.human, got, tt.want)
			}
		})
	}
}

func TestFormatOutput(t *testing.T) {
	tests := []struct {
		name    string
		sizeStr string
		path    string
		want    string
	}{
		{name: "simple", sizeStr: "123B", path: "file.txt", want: "123B\tfile.txt"},
		{name: "human size", sizeStr: "1.2MB", path: "output.dat", want: "1.2MB\toutput.dat"},
		{name: "empty path", sizeStr: "0B", path: "", want: "0B\t"},
		{
			name:    "path with spaces",
			sizeStr: "24.0MB",
			path:    "my file.dat",
			want:    "24.0MB\tmy file.dat",
		},
		{
			name:    "already contains unit",
			sizeStr: "999KB",
			path:    "dir/name",
			want:    "999KB\tdir/name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatOutput(tt.sizeStr, tt.path)
			if got != tt.want {
				t.Fatalf("FormatOutput(%q, %q) = %q, want %q", tt.sizeStr, tt.path, got, tt.want)
			}
		})
	}
}
