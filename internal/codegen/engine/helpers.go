package engine

import (
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GetHelperFunctions returns template helper functions
func GetHelperFunctions() template.FuncMap {
	return template.FuncMap{
		"pascalCase": pascalCase,
		"camelCase":  camelCase,
		"snakeCase":  snakeCase,
		"kebabCase":  kebabCase,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"title":      titleCase,
		"plural":     plural,
		"singular":   singular,
		"join":       strings.Join,
		"split":      strings.Split,
		"contains":   strings.Contains,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
	}
}

// titleCase converts a string to Title Case using golang.org/x/text/cases
func titleCase(s string) string {
	caser := cases.Title(language.English)
	return caser.String(s)
}

// pascalCase converts a string to PascalCase
func pascalCase(s string) string {
	if s == "" {
		return ""
	}

	// Split by common delimiters
	parts := splitByDelimiters(s)

	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Capitalize first letter, lowercase rest
		runes := []rune(part)
		result.WriteRune(unicode.ToUpper(runes[0]))
		result.WriteString(strings.ToLower(string(runes[1:])))
	}

	return result.String()
}

// camelCase converts a string to camelCase
func camelCase(s string) string {
	pascal := pascalCase(s)
	if pascal == "" {
		return ""
	}

	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// snakeCase converts a string to snake_case
func snakeCase(s string) string {
	parts := splitByDelimiters(s)
	return strings.ToLower(strings.Join(parts, "_"))
}

// kebabCase converts a string to kebab-case
func kebabCase(s string) string {
	parts := splitByDelimiters(s)
	return strings.ToLower(strings.Join(parts, "-"))
}

// splitByDelimiters splits a string by common delimiters
func splitByDelimiters(s string) []string {
	// Replace common delimiters with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, ".", " ")

	// Split by spaces and capital letters
	var parts []string
	var current strings.Builder

	for i, r := range s {
		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else if unicode.IsUpper(r) && i > 0 {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			current.WriteRune(r)
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// plural returns the plural form of a word (simple implementation)
func plural(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}
	return s + "s"
}

// singular returns the singular form of a word (simple implementation)
func singular(s string) string {
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}
	if strings.HasSuffix(s, "es") {
		if strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "ches") || strings.HasSuffix(s, "shes") {
			return s[:len(s)-2]
		}
		return s[:len(s)-1]
	}
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}
	return s
}
