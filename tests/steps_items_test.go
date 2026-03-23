package tests

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
)

func (sc *scenarioCtx) serverHasMovies(table *godog.Table) error {
	var items []map[string]any
	for i, row := range table.Rows {
		if i == 0 {
			continue
		}
		year, _ := strconv.Atoi(row.Cells[2].Value)
		items = append(items, map[string]any{
			"Id":             row.Cells[0].Value,
			"Name":           row.Cells[1].Value,
			"Type":           "Movie",
			"ProductionYear": year,
		})
	}
	sc.mock.On("GET", "/Users/test-user-id/Items", 200, map[string]any{
		"Items":            items,
		"TotalRecordCount": len(items),
	})
	return nil
}

func (sc *scenarioCtx) serverHasItemsOfType(itemType string, table *godog.Table) error {
	var items []map[string]any
	for i, row := range table.Rows {
		if i == 0 {
			continue
		}
		items = append(items, map[string]any{
			"Id":   row.Cells[0].Value,
			"Name": row.Cells[1].Value,
			"Type": itemType,
		})
	}
	sc.mock.On("GET", "/Users/test-user-id/Items", 200, map[string]any{
		"Items":            items,
		"TotalRecordCount": len(items),
	})
	return nil
}

func (sc *scenarioCtx) serverHasItem(itemID, name string, year int) error {
	item := map[string]any{
		"Id":             itemID,
		"Name":           name,
		"Type":           "Movie",
		"ProductionYear": year,
	}
	// Register for both GetItem and UpdateItem GET
	sc.mock.On("GET", "/Users/test-user-id/Items/"+itemID, 200, item)
	sc.mock.On("GET", "/Items/"+itemID, 200, item)
	sc.mock.On("POST", "/Items/"+itemID, 204, nil)
	return nil
}

func (sc *scenarioCtx) iListMovies() error {
	result, err := ops.ListMovies(sc.client)
	sc.lastErr = err
	sc.itemsResult = result
	return nil
}

func (sc *scenarioCtx) iListItemsWithType(itemType string) error {
	result, err := ops.ListItems(sc.client, ops.ListItemsOptions{ItemType: itemType})
	sc.lastErr = err
	sc.itemsResult = result
	return nil
}

func (sc *scenarioCtx) iGetInfoForItem(itemID string) error {
	info, err := ops.GetItemInfo(sc.client, itemID)
	sc.lastErr = err
	sc.itemInfo = info
	return nil
}

func (sc *scenarioCtx) iShouldGetNItems(n int) error {
	result, ok := sc.itemsResult.(*ops.ItemsResult)
	if !ok {
		if n == 0 {
			return nil
		}
		return fmt.Errorf("no items result available")
	}
	if len(result.Items) != n {
		return fmt.Errorf("expected %d items, got %d", n, len(result.Items))
	}
	return nil
}

func (sc *scenarioCtx) itemShouldHaveYear(name string, year int) error {
	result, ok := sc.itemsResult.(*ops.ItemsResult)
	if !ok {
		return fmt.Errorf("no items result available")
	}
	for _, item := range result.Items {
		if item.Name == name {
			if item.ProductionYear == year {
				return nil
			}
			return fmt.Errorf("item %q has year %d, expected %d", name, item.ProductionYear, year)
		}
	}
	return fmt.Errorf("item %q not found", name)
}

func (sc *scenarioCtx) itemNameShouldBe(expected string) error {
	info, ok := sc.itemInfo.(*ops.ItemEntry)
	if !ok {
		return fmt.Errorf("no item info available")
	}
	if info.Name != expected {
		return fmt.Errorf("expected name %q, got %q", expected, info.Name)
	}
	return nil
}

func (sc *scenarioCtx) itemYearShouldBe(expected int) error {
	info, ok := sc.itemInfo.(*ops.ItemEntry)
	if !ok {
		return fmt.Errorf("no item info available")
	}
	if info.ProductionYear != expected {
		return fmt.Errorf("expected year %d, got %d", expected, info.ProductionYear)
	}
	return nil
}
