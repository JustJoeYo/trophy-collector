package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestLiveItemsPayload_ContainsSlotTypeAndCost(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live API contract test in short mode")
	}

	assetsURL := "https://assets.deadlock-api.com"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, assetsURL+"/v2/items", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("live request failed: %v", err)
	}

	var payload []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(payload) == 0 {
		t.Fatalf("expected non-empty items payload")
	}

	foundSlotType := false
	foundCost := false
	for _, item := range payload {
		if v, ok := item["item_slot_type"]; ok {
			if s, ok := v.(string); ok && s != "" {
				foundSlotType = true
			}
		}

		if _, ok := item["cost"]; ok {
			foundCost = true
		}
		if foundSlotType && foundCost {
			break
		}

	}

	if !foundSlotType {
		t.Fatalf("live payload missing non-empty item_slot_type")
	}
	if !foundCost {
		t.Fatalf("live payload missing cost key")
	}
}
