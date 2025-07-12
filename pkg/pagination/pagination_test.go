package pagination

import (
	"testing"
)

func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		limit      int
		totalItems int
		expect     struct {
			currentPage int
			pageSize    int
			totalPages  int
			hasNext     bool
			hasPrev     bool
		}
	}{
		{
			name:       "Normal middle page",
			page:       2,
			limit:      10,
			totalItems: 35,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{2, 10, 4, true, true},
		},
		{
			name:       "First page",
			page:       1,
			limit:      10,
			totalItems: 35,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{1, 10, 4, true, false},
		},
		{
			name:       "Last page",
			page:       4,
			limit:      10,
			totalItems: 35,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{4, 10, 4, false, true},
		},
		{
			name:       "Page out of range (too high)",
			page:       10,
			limit:      10,
			totalItems: 35,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{10, 10, 4, false, true},
		},
		{
			name:       "Zero page and limit",
			page:       0,
			limit:      0,
			totalItems: 35,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{1, 10, 4, true, false},
		},
		{
			name:       "Negative page and limit",
			page:       -1,
			limit:      -5,
			totalItems: 35,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{1, 10, 4, true, false},
		},
		{
			name:       "No items",
			page:       1,
			limit:      10,
			totalItems: 0,
			expect:     struct{currentPage, pageSize, totalPages int; hasNext, hasPrev bool}{1, 10, 1, false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentPage, pageSize, totalPages, hasNext, hasPrev := CalculatePagination(tt.page, tt.limit, tt.totalItems)
			if currentPage != tt.expect.currentPage {
				t.Errorf("currentPage: got %d, want %d", currentPage, tt.expect.currentPage)
			}
			if pageSize != tt.expect.pageSize {
				t.Errorf("pageSize: got %d, want %d", pageSize, tt.expect.pageSize)
			}
			if totalPages != tt.expect.totalPages {
				t.Errorf("totalPages: got %d, want %d", totalPages, tt.expect.totalPages)
			}
			if hasNext != tt.expect.hasNext {
				t.Errorf("hasNext: got %v, want %v", hasNext, tt.expect.hasNext)
			}
			if hasPrev != tt.expect.hasPrev {
				t.Errorf("hasPrev: got %v, want %v", hasPrev, tt.expect.hasPrev)
			}
		})
	}
} 