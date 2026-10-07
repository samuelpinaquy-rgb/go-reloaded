package main

import (
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		os.Stdout.WriteString("Usage: go run . input.txt output.txt\n")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		os.Stderr.WriteString("Erreur lors de la lecture du fichier : " + err.Error() + "\n")
		return
	}

	text := string(content)
	text, err = convertNumbers(text)
	if err != nil {
		os.Stderr.WriteString("Erreur lors de la conversion des nombres : " + err.Error() + "\n")
		return
	}

	if err := os.WriteFile(outputFile, []byte(text), 0644); err != nil {
		os.Stderr.WriteString("Erreur lors de la création du fichier : " + err.Error() + "\n")
		return
	}

	os.Stdout.WriteString(text + "\n")

	os.Stdout.WriteString("Fichier de sortie : " + outputFile + "\n")
}

func convertNumbers(text string) (string, error) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text, nil
	}

	separators := make([]string, len(words)+1)
	position := 0
	for i, word := range words {
		start := position + strings.Index(text[position:], word)
		separators[i] = text[position:start]
		position = start + len(word)
	}
	separators[len(words)] = text[position:]

	var result strings.Builder
	for i := 0; i < len(words); i++ {
		word := words[i]
		if _, _, isMarker := numberMarker(word); isMarker {
			return "", strconv.ErrSyntax
		}

		result.WriteString(separators[i])
		if i+1 < len(words) {
			base, punctuation, isMarker := numberMarker(words[i+1])
			if !isMarker {
				result.WriteString(word)
				continue
			}
			number, err := strconv.ParseInt(word, base, 64)
			if err != nil {
				return "", err
			}
			result.WriteString(strconv.FormatInt(number, 10))
			result.WriteString(punctuation)
			i++
			continue
		}

		result.WriteString(word)
	}
	result.WriteString(separators[len(words)])
	return result.String(), nil
}

func numberMarker(word string) (int, string, bool) {
	marker := "(hex)"
	base := 16
	if strings.HasPrefix(word, "(bin)") {
		marker = "(bin)"
		base = 2
	} else if !strings.HasPrefix(word, marker) {
		return 0, "", false
	}

	punctuation := strings.TrimPrefix(word, marker)
	if strings.Trim(punctuation, ".,!?;:'") != "" {
		return 0, "", false
	}
	return base, punctuation, true
}
