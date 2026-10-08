package main

import "testing"

// TestConvertNumbers utilise plusieurs exemples pour vérifier les conversions hexadécimales
// et binaires, la ponctuation adjacente, la conservation des espaces et des retours à la ligne,
// ainsi que le comportement en l'absence de balise.
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
			got := convertNumbers(test.input)
			if got != test.want {
				t.Errorf("convertNumbers() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestConvertNumbersLeavesInvalidMarkersUnchanged vérifie que les balises numériques invalides
// restent dans le texte au lieu d'empêcher la conversion d'autres nombres.
func TestConvertNumbersLeavesInvalidMarkersUnchanged(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "marker without a preceding number",
			input: "(hex) files",
			want:  "(hex) files",
		},
		{
			name:  "invalid hexadecimal number",
			input: "not-a-number (hex)",
			want:  "not-a-number (hex)",
		},
		{
			name:  "convert valid number after invalid marker",
			input: "(hex) then 10 (bin)",
			want:  "(hex) then 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := convertNumbers(test.input)
			if got != test.want {
				t.Errorf("convertNumbers() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestConvertCase vérifie les trois modes de casse, leurs variantes avec un nombre de mots,
// la ponctuation, les retours à la ligne et le texte qui ne contient aucune balise.
func TestConvertCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "uppercase previous word",
			input: "Ready, set, go (up)!",
			want:  "Ready, set, GO!",
		},
		{
			name:  "lowercase previous word",
			input: "I should stop SHOUTING (low)",
			want:  "I should stop shouting",
		},
		{
			name:  "capitalize previous word",
			input: "welcome to the brooklyn bridge (cap)",
			want:  "welcome to the brooklyn Bridge",
		},
		{
			name:  "uppercase requested number of words",
			input: "This is so exciting (up, 2)",
			want:  "This is SO EXCITING",
		},
		{
			name:  "capitalize requested number of words",
			input: "welcome to the brooklyn bridge (cap, 2)",
			want:  "welcome to the Brooklyn Bridge",
		},
		{
			name:  "preserve whitespace and line breaks",
			input: "hello\nWORLD (low)",
			want:  "hello\nworld",
		},
		{
			name:  "leave text without markers unchanged",
			input: "No case changes here.",
			want:  "No case changes here.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := convertCase(test.input)
			if got != test.want {
				t.Errorf("convertCase() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestConvertCaseLeavesInvalidMarkersUnchanged vérifie que les balises mal formées ou impossibles
// à appliquer restent dans le résultat, sans empêcher le traitement d'une balise valide ultérieure.
func TestConvertCaseLeavesInvalidMarkersUnchanged(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "marker without a preceding word",
			input: "(up)",
			want:  "(up)",
		},
		{
			name:  "count exceeds preceding words",
			input: "hello (up, 2)",
			want:  "hello (up, 2)",
		},
		{
			name:  "zero count",
			input: "hello (up, 0)",
			want:  "hello (up, 0)",
		},
		{
			name:  "non-numeric count",
			input: "hello (up, many)",
			want:  "hello (up, many)",
		},
		{
			name:  "keep invalid marker and process a later valid one",
			input: "hello (up, 2) world (low)",
			want:  "hello (up, 2) world",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := convertCase(test.input)
			if got != test.want {
				t.Errorf("convertCase() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestFormatApostrophes vérifie que les espaces intérieurs aux apostrophes appariées
// sont retirés, sans modifier les apostrophes qui n'ont pas de paire.
func TestFormatApostrophes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove spaces around one quoted word",
			input: "They said ' awesome ' yesterday",
			want:  "They said 'awesome' yesterday",
		},
		{
			name:  "remove spaces around a quoted phrase",
			input: "He said ' I am ready ' today",
			want:  "He said 'I am ready' today",
		},
		{
			name:  "leave an unmatched apostrophe unchanged",
			input: "This ' quote is unmatched",
			want:  "This ' quote is unmatched",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatApostrophes(test.input)
			if got != test.want {
				t.Errorf("formatApostrophes() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestFixPunctuation vérifie le collage des signes au mot précédent,
// le maintien des groupes de signes et l'ajout d'un espace après eux.
func TestFixPunctuation(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "move punctuation away from following word",
			input: "over there ,and then BAMM !!",
			want:  "over there, and then BAMM!!",
		},
		{
			name:  "keep punctuation groups together",
			input: "I was thinking ... You were right !?",
			want:  "I was thinking... You were right!?",
		},
		{
			name:  "preserve line breaks after punctuation",
			input: "Hello ,\nworld !",
			want:  "Hello,\nworld!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := fixPunctuation(test.input)
			if got != test.want {
				t.Errorf("fixPunctuation() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestConvertArticle vérifie la transformation de "a" en "an" devant une voyelle
// ou h, ainsi que la casse, la ponctuation et les séparateurs d'origine.
func TestConvertArticle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "add n before a vowel",
			input: "a amazing story",
			want:  "an amazing story",
		},
		{
			name:  "add n before h",
			input: "a honest person",
			want:  "an honest person",
		},
		{
			name:  "preserve initial capitalization",
			input: "A unusual day",
			want:  "An unusual day",
		},
		{
			name:  "leave article before consonant unchanged",
			input: "a great day",
			want:  "a great day",
		},
		{
			name:  "preserve punctuation and whitespace",
			input: "A,  amazing!",
			want:  "An,  amazing!",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := convertArticle(test.input)
			if got != test.want {
				t.Errorf("convertArticle() = %q, want %q", got, test.want)
			}
		})
	}
}
