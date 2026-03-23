package tests

import "github.com/jrogala/jellyfin-cli/pkg/ops"

func (sc *scenarioCtx) iUpdateItemWithName(itemID, name string) error {
	sc.lastErr = ops.UpdateItem(sc.client, itemID, map[string]any{"Name": name})
	return nil
}

func (sc *scenarioCtx) serverAcceptsIdentify(itemID string) error {
	sc.mock.On("POST", "/Items/"+itemID+"/Identify", 204, nil)
	return nil
}

func (sc *scenarioCtx) iIdentifyItem(itemID, name, imdb string) error {
	providerIDs := map[string]string{"Imdb": imdb}
	sc.lastErr = ops.IdentifyItem(sc.client, itemID, name, 0, providerIDs)
	return nil
}

func (sc *scenarioCtx) serverAcceptsRefresh(itemID string) error {
	sc.mock.On("POST", "/Items/"+itemID+"/Refresh", 204, nil)
	return nil
}

func (sc *scenarioCtx) iRefreshMetadata(itemID string) error {
	sc.lastErr = ops.RefreshMetadata(sc.client, itemID)
	return nil
}

func (sc *scenarioCtx) serverAcceptsLibraryScan() error {
	sc.mock.On("POST", "/Library/Refresh", 204, nil)
	return nil
}

func (sc *scenarioCtx) iScanLibraries() error {
	sc.lastErr = ops.ScanLibrary(sc.client)
	return nil
}
