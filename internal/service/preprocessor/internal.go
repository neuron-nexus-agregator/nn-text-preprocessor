package preprocessor

import (
	"regexp"
	"strings"
)

func (p *Preprocessor) clearHTML(text string, safeTextTags bool) string {
	switch safeTextTags {
	case true:
		return p.clearNonTextTags(text)
	default:
		return p.clearAllHTML(text)
	}
}

func (p *Preprocessor) clearAllHTML(text string) string {
	// Регулярное выражение для удаления всех HTML-тегов
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(text, "")
}

func (p *Preprocessor) clearNonTextTags(text string) string {
	// Регулярное выражение для всех HTML-тегов.
	// Группы:
	// 1: '/' если это закрывающий тег (опционально)
	// 2: Имя тега
	// 3: Атрибуты тега (опционально)
	allTagsRegex := regexp.MustCompile(`<(\/?)([a-zA-Z0-9]+)([^>]*)?>`)

	// Множество разрешенных имен тегов для быстрого поиска
	allowedTagNames := map[string]bool{
		"b": true, "i": true, "strong": true, "em": true, "u": true,
		"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
		"p":  true,
		"ul": true, "ol": true, "li": true,
		"table": true, "tr": true, "td": true, "th": true, "thead": true, "tbody": true, "tfoot": true,
	}

	// Используем ReplaceAllStringFunc для обработки каждого найденного тега
	return allTagsRegex.ReplaceAllStringFunc(text, func(tag string) string {
		matches := allTagsRegex.FindStringSubmatch(tag)
		if len(matches) < 3 {
			// Это не должно произойти, если regex корректен, но на всякий случай
			return "" // Удаляем некорректный тег
		}

		// matches[2] содержит имя тега (например, "div", "b", "script")
		tagName := strings.ToLower(matches[2]) // Приводим к нижнему регистру для сравнения

		// Проверяем, является ли имя тега разрешенным
		if allowedTagNames[tagName] {
			return tag // Если разрешен, сохраняем исходный тег
		}
		return "" // Если не разрешен, удаляем (заменяем на пустую строку)
	})
}

// func (p *Preprocessor) normalizeText(text string) string {
// 	if text == "" {
// 		return ""
// 	}
// 	// Convert to lowercase and remove non-graphic characters
// 	text = p.clearAllHTML(strings.ToLower(text))
// 	text = strings.Map(func(r rune) rune {
// 		if unicode.IsGraphic(r) {
// 			return r
// 		}
// 		return -1
// 	}, text)
// 	// Replace multiple spaces with single space and trim
// 	text = strings.Join(strings.Fields(text), " ")
// 	return strings.TrimSpace(text)
// }

// func (p *Preprocessor) cosineSimilarity(text1, text2 string) float64 {
// 	if text1 == "" || text2 == "" {
// 		return 0.0
// 	}

// 	text1 = p.normalizeText(text1)
// 	text2 = p.normalizeText(text2)

// 	// Split texts into words
// 	words1 := strings.Split(text1, " ")
// 	words2 := strings.Split(text2, " ")

// 	// Build word frequency maps
// 	freq1 := make(map[string]int)
// 	freq2 := make(map[string]int)
// 	for _, word := range words1 {
// 		freq1[word]++
// 	}
// 	for _, word := range words2 {
// 		freq2[word]++
// 	}

// 	// Calculate dot product and magnitudes
// 	dotProduct := 0.0
// 	magnitude1 := 0.0
// 	magnitude2 := 0.0
// 	for word, count1 := range freq1 {
// 		count2, exists := freq2[word]
// 		if exists {
// 			dotProduct += float64(count1 * count2)
// 		}
// 		magnitude1 += float64(count1 * count1)
// 	}
// 	for _, count2 := range freq2 {
// 		magnitude2 += float64(count2 * count2)
// 	}

// 	// Avoid division by zero
// 	if magnitude1 == 0 || magnitude2 == 0 {
// 		return 0.0
// 	}

// 	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
// }
