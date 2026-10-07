package main

import "testing"

func TestConvertNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "convert hexadecimal number",
			input: "1E (hex) files were added",
			want:  "30 files were added",
		},
		{
			name:  "convert binary number",
			input: "It has been 10 (bin) years",
			want:  "It has been 2 years",
		},
		{
			name:  "convert markers followed by punctuation",
			input: "42 (hex), then 10 (bin).",
			want:  "66, then 2.",
		},
		{
			name:  "preserve whitespace and line breaks",
			input: " 42 (hex)\n\t10 (bin) ",
			want:  " 66\n\t2 ",
		},
		{
			name:  "leave text without markers unchanged",
			input: "No numbers to convert.",
			want:  "No numbers to convert.",
		},
		{
			name:  "leave empty text unchanged",
			input: "",
			want:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := convertNumbers(test.input)
			if err != nil {
				t.Fatalf("convertNumbers() returned an unexpected error: %v", err)
			}
			if got != test.want {
				t.Errorf("convertNumbers() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestConvertNumbersReturnsErrorForInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "marker without a preceding number",
			input: "(hex) files",
		},
		{
			name:  "invalid hexadecimal number",
			input: "not-a-number (hex)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := convertNumbers(test.input)
			if err == nil {
				t.Fatal("convertNumbers() did not return an error")
			}
		})
	}
}
