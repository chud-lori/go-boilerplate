package pagination

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit, totalItems int) (currentPage, pageSize, totalPages int, hasNext, hasPrev bool) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Calculate total pages (ceiling division)
	totalPages = (totalItems + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Calculate pagination flags
	hasNext = page < totalPages
	hasPrev = page > 1

	return page, limit, totalPages, hasNext, hasPrev
} 