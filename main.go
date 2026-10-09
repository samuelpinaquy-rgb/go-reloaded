package main

import (
	"os"
	"strconv"
	"strings"
)

// main vérifie les deux arguments de la commande, qui indiquent le fichier d'entrée et le fichier de sortie.
// Elle lit le texte, applique les conversions numériques, les balises de casse,
// les apostrophes, la ponctuation et la règle de "a/an", puis écrit le résultat.
// Une erreur de lecture ou d'écriture est affichée dans la sortie d'erreur et arrête le traitement.
func main() {
	if len(os.Args) != 3 {
		os.Stdout.WriteString("Usage: go run . input.txt output.txt\n")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		os.Stderr.WriteString("Erreur lors de la lecture du fichier : " + err.Error() + "\n")
		os.Exit(1)
	}

	text := string(content)
	text = convertNumbers(text)
	text = convertCase(text)
	text = formatApostrophes(text)
	text = fixPunctuation(text)
	text = convertArticle(text)

	if err := os.WriteFile(outputFile, []byte(text), 0644); err != nil {
		os.Stderr.WriteString("Erreur lors de la création du fichier : " + err.Error() + "\n")
		os.Exit(1)
	}

	os.Stdout.WriteString(text + "\n")

	os.Stdout.WriteString("Fichier de sortie : " + outputFile + "\n")
}

// convertNumbers cherche les balises (hex) et (bin), puis convertit le mot qui précède
// respectivement depuis la base 16 ou la base 2 vers la base 10.
// Elle conserve les séparateurs d'origine, y compris les espaces et retours à la ligne,
// et garde la ponctuation qui suit la balise. Une balise sans nombre valide devant elle
// est conservée telle quelle afin que le reste du texte puisse tout de même être traité.
func convertNumbers(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
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
		result.WriteString(separators[i])
		if i+1 < len(words) {
			base, punctuation, isMarker := numberMarker(words[i+1])
			if !isMarker {
				result.WriteString(word)
				continue
			}
			number, err := strconv.ParseInt(word, base, 64)
			if err != nil {
				result.WriteString(word)
				continue
			}
			result.WriteString(strconv.FormatInt(number, 10))
			result.WriteString(punctuation)
			i++
			continue
		}

		result.WriteString(word)
	}
	result.WriteString(separators[len(words)])
	return result.String()
}

// numberMarker vérifie si un mot commence par (hex) ou (bin), éventuellement suivi de ponctuation.
// Si c'est une balise reconnue, elle renvoie la base à utiliser, la ponctuation à conserver
// et true. Sinon, elle renvoie des valeurs vides et false.
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

// convertCase applique les balises (up), (low) et (cap) aux mots qui les précèdent.
// Un nombre indiqué dans la balise étend la transformation aux mots précédents demandés.
// La fonction conserve les séparateurs du texte et retire uniquement les balises valides ;
// une balise qui ne peut pas être appliquée reste inchangée avec son texte.
func convertCase(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	separators := make([]string, len(words)+1)
	position := 0
	for i, word := range words {
		start := position + strings.Index(text[position:], word)
		separators[i] = text[position:start]
		position = start + len(word)
	}
	separators[len(words)] = text[position:]

	omitted := make([]bool, len(words))
	for i, word := range words {
		if omitted[i] {
			continue
		}

		markerWord := word
		countWord := -1
		if strings.HasSuffix(word, ",") && i+1 < len(words) && strings.HasSuffix(words[i+1], ")") {
			markerWord += words[i+1]
			countWord = i + 1
		}

		mode, count, punctuation, isMarker := caseMarker(markerWord)
		if !isMarker {
			continue
		}

		previous := make([]int, 0, count)
		for j := i - 1; j >= 0 && len(previous) < count; j-- {
			if !omitted[j] {
				previous = append(previous, j)
			}
		}
		if len(previous) != count {
			continue
		}

		omitted[i] = true
		if countWord >= 0 {
			omitted[countWord] = true
		}
		for _, index := range previous {
			switch mode {
			case "up":
				words[index] = strings.ToUpper(words[index])
			case "low":
				words[index] = strings.ToLower(words[index])
			case "cap":
				words[index] = capitalize(words[index])
			}
		}
		if punctuation != "" {
			words[previous[0]] += punctuation
		}
	}

	var result strings.Builder
	for i, word := range words {
		if !omitted[i] {
			result.WriteString(separators[i])
			result.WriteString(word)
		}
	}
	result.WriteString(separators[len(words)])
	return result.String()
}

// caseMarker analyse une balise de casse, par exemple (up), (cap, 2) ou (low)!.
// Elle renvoie le mode demandé, le nombre de mots à transformer, la ponctuation éventuelle
// et true si la balise est valide. Pour une balise inconnue ou mal formée, elle renvoie false.
func caseMarker(word string) (string, int, string, bool) {
	if !strings.HasPrefix(word, "(") {
		return "", 0, "", false
	}

	end := strings.Index(word, ")")
	if end == -1 {
		return "", 0, "", false
	}

	content := word[1:end]
	punctuation := word[end+1:]
	if strings.Trim(punctuation, ".,!?;:'") != "" {
		return "", 0, "", false
	}

	parts := strings.Split(content, ",")
	mode := strings.TrimSpace(parts[0])
	if mode != "up" && mode != "low" && mode != "cap" {
		return "", 0, "", false
	}
	if len(parts) > 2 {
		return "", 0, "", false
	}
	if len(parts) == 1 {
		return mode, 1, punctuation, true
	}

	count, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || count < 1 {
		return "", 0, "", false
	}
	return mode, count, punctuation, true
}

// capitalize met la première lettre du mot en majuscule et toutes les lettres suivantes en minuscules.
// Un mot vide est renvoyé tel quel afin de ne pas tenter d'accéder à un caractère inexistant.
func capitalize(word string) string {
	if word == "" {
		return word
	}
	letters := []rune(word)
	return strings.ToUpper(string(letters[0])) + strings.ToLower(string(letters[1:]))
}

// formatApostrophes enlève les espaces entre les apostrophes simples appariées
// et les mots qu'elles entourent. Les espaces à l'extérieur des apostrophes
// restent inchangés ; une apostrophe sans paire est laissée telle quelle.
func formatApostrophes(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	separators := make([]string, len(words)+1)
	position := 0
	for i, word := range words {
		start := position + strings.Index(text[position:], word)
		separators[i] = text[position:start]
		position = start + len(word)
	}
	separators[len(words)] = text[position:]

	opening := -1
	for i, word := range words {
		if word != "'" {
			continue
		}
		if opening == -1 {
			opening = i
			continue
		}

		separators[opening+1] = ""
		separators[i] = ""
		opening = -1
	}

	var result strings.Builder
	for i, word := range words {
		result.WriteString(separators[i])
		result.WriteString(word)
	}
	result.WriteString(separators[len(words)])
	return result.String()
}

// fixPunctuation colle les signes .,!?;: au mot précédent.
// Les signes consécutifs restent groupés, et un espace les sépare du texte suivant,
// même si le signe et le mot suivant étaient collés dans le texte d'entrée.
func fixPunctuation(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
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
	for i, word := range words {
		result.WriteString(separators[i])
		previousWasPunctuation := false
		for _, character := range word {
			if strings.ContainsRune(".,!?;:", character) {
				trimTrailingWhitespace(&result)
				result.WriteRune(character)
				previousWasPunctuation = true
				continue
			}

			if previousWasPunctuation {
				result.WriteByte(' ')
			}
			result.WriteRune(character)
			previousWasPunctuation = false
		}
	}
	result.WriteString(separators[len(words)])
	return result.String()
}

// trimTrailingWhitespace retire les espaces de fin déjà écrits afin de coller
// le signe de ponctuation au texte précédent, tout en préservant les espaces initiaux.
func trimTrailingWhitespace(text *strings.Builder) {
	current := text.String()
	trimmed := strings.TrimRight(current, " \t\r\n")
	if len(trimmed) != len(current) && trimmed != "" {
		text.Reset()
		text.WriteString(trimmed)
	}
}

// convertArticle transforme l'article anglais "a" en "an" devant une voyelle ou h.
// Elle conserve la majuscule de "A", les signes de ponctuation après l'article
// et les séparateurs présents entre les mots.
func convertArticle(text string) string {
	words := strings.Fields(text)
	if len(words) < 2 {
		return text
	}

	separators := make([]string, len(words)+1)
	position := 0
	for i, word := range words {
		start := position + strings.Index(text[position:], word)
		separators[i] = text[position:start]
		position = start + len(word)
	}
	separators[len(words)] = text[position:]

	for i := 0; i < len(words)-1; i++ {
		article := strings.TrimRight(words[i], ".,!?;:")
		if article != "a" && article != "A" {
			continue
		}

		nextWord := strings.TrimLeft(words[i+1], ".,!?;:'\"([{")
		if nextWord == "" {
			continue
		}
		firstLetter := strings.ToLower(nextWord[:1])
		if !strings.Contains("aeiouh", firstLetter) {
			continue
		}

		if article == "A" {
			words[i] = "An" + strings.TrimPrefix(words[i], "A")
		} else {
			words[i] = "an" + strings.TrimPrefix(words[i], "a")
		}
	}

	var result strings.Builder
	for i, word := range words {
		result.WriteString(separators[i])
		result.WriteString(word)
	}
	result.WriteString(separators[len(words)])
	return result.String()
}
func splitWordsAndSeparators(text string) (words, separators []string) {
	words = strings.Fields(text)
	if len(words) == 0 {
		return
	}

	separators = make([]string, len(words)+1)
	position := 0
	for i, word := range words {
		start := position + strings.Index(text[position:], word)
		separators[i] = text[position:start]
		position = start + len(word)
	}
	separators[len(words)] = text[position:]
	return
}
