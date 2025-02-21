package dto

import "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"

// URLMappings represents a mapping of slugs to their corresponding URL.
type URLMappings map[domain.Slug]domain.URLMapping

// URLMappingsCopy creates and returns a deep copy of the given URLMappings.
func URLMappingsCopy(m URLMappings) URLMappings {
	cp := make(URLMappings)
	for k, v := range m {
		cp[k] = v
	}

	return cp
}
