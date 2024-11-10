package dto

import "github.com/patraden/ya-practicum-go-shortly/internal/app/domain"

type URLMappings map[domain.Slug]domain.URLMapping

func URLMappingsCopy(m URLMappings) URLMappings {
	cp := make(URLMappings)
	for k, v := range m {
		cp[k] = v
	}

	return cp
}
