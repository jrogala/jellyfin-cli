package tests

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
)

func (sc *scenarioCtx) serverReturnsSearchResults(query string, table *godog.Table) error {
	var hints []map[string]any
	for i, row := range table.Rows {
		if i == 0 {
			continue
		}
		year, _ := strconv.Atoi(row.Cells[3].Value)
		hints = append(hints, map[string]any{
			"Id":             row.Cells[0].Value,
			"Name":           row.Cells[1].Value,
			"Type":           row.Cells[2].Value,
			"ProductionYear": year,
		})
	}
	sc.mock.On("GET", "/Search/Hints", 200, map[string]any{
		"SearchHints":      hints,
		"TotalRecordCount": len(hints),
	})
	return nil
}

func (sc *scenarioCtx) serverReturnsNoSearchResults(query string) error {
	sc.mock.On("GET", "/Search/Hints", 200, map[string]any{
		"SearchHints":      []any{},
		"TotalRecordCount": 0,
	})
	return nil
}

func (sc *scenarioCtx) iSearchFor(query string, limit int) error {
	result, err := ops.Search(sc.client, query, limit)
	sc.lastErr = err
	sc.searchResult = result
	return nil
}

func (sc *scenarioCtx) iShouldGetNSearchResults(n int) error {
	result, ok := sc.searchResult.(*ops.SearchResult)
	if !ok {
		if n == 0 && sc.searchResult == nil {
			return nil
		}
		return fmt.Errorf("no search results available")
	}
	if len(result.Hints) != n {
		return fmt.Errorf("expected %d search results, got %d", n, len(result.Hints))
	}
	return nil
}

func (sc *scenarioCtx) searchResultShouldHaveType(name, expectedType string) error {
	result, ok := sc.searchResult.(*ops.SearchResult)
	if !ok {
		return fmt.Errorf("no search results available")
	}
	for _, h := range result.Hints {
		if h.Name == name {
			if h.Type == expectedType {
				return nil
			}
			return fmt.Errorf("search result %q has type %q, expected %q", name, h.Type, expectedType)
		}
	}
	return fmt.Errorf("search result %q not found", name)
}
