package tests

import (
	"fmt"

	"github.com/cucumber/godog"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
)

func (sc *scenarioCtx) serverHasLibraries(table *godog.Table) error {
	var items []map[string]string
	for i, row := range table.Rows {
		if i == 0 {
			continue // header
		}
		items = append(items, map[string]string{
			"Id":             row.Cells[0].Value,
			"Name":           row.Cells[1].Value,
			"CollectionType": row.Cells[2].Value,
		})
	}
	sc.mock.On("GET", "/Users/test-user-id/Views", 200, map[string]any{
		"Items": items,
	})
	return nil
}

func (sc *scenarioCtx) serverHasNoLibraries() error {
	sc.mock.On("GET", "/Users/test-user-id/Views", 200, map[string]any{
		"Items": []any{},
	})
	return nil
}

func (sc *scenarioCtx) iListLibraries() error {
	libs, err := ops.ListLibraries(sc.client)
	sc.lastErr = err
	sc.libraries = libs
	return nil
}

func (sc *scenarioCtx) iShouldGetNLibraries(n int) error {
	libs, ok := sc.libraries.([]ops.LibraryEntry)
	if !ok {
		if n == 0 && sc.libraries == nil {
			return nil
		}
		return fmt.Errorf("no library results available")
	}
	if len(libs) != n {
		return fmt.Errorf("expected %d libraries, got %d", n, len(libs))
	}
	return nil
}

func (sc *scenarioCtx) libraryShouldHaveType(name, expectedType string) error {
	libs, ok := sc.libraries.([]ops.LibraryEntry)
	if !ok {
		return fmt.Errorf("no library results available")
	}
	for _, l := range libs {
		if l.Name == name {
			if l.CollectionType == expectedType {
				return nil
			}
			return fmt.Errorf("library %q has type %q, expected %q", name, l.CollectionType, expectedType)
		}
	}
	return fmt.Errorf("library %q not found", name)
}
