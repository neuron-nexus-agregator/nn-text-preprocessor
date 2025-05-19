package preprocessor

import "regexp"

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
	// Регулярное выражение для удаления всех HTML-тегов, кроме разрешенных
	allowedTags := regexp.MustCompile(`<(\/?(b|i|strong|em|u|h[1-6]|p|ul|ol|li|table|tr|td|th|thead|tbody|tfoot))[^>]*>`)
	allTags := regexp.MustCompile(`<[^>]*>`)

	// Удаляем все теги, кроме разрешенных
	text = allowedTags.ReplaceAllStringFunc(text, func(tag string) string {
		return tag // Сохраняем разрешенные теги
	})

	// Удаляем оставшиеся теги
	return allTags.ReplaceAllString(text, "")
}
