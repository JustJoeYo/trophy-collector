package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetItems_decodesItemSlotType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/items" {
			t.Fatalf("expected path /v2/items, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"class_name": "weapon_core",
				"name": "Basic Magazine",
				"item_slot_type": "weapon"
			}
		]`))
	}))
	defer ts.Close()

	client := NewDeadlockClient("http:..unused-base-url", ts.URL)

	var payload []map[string]any
	if err := client.(*deadlockClient).fetch(context.Background(), ts.URL+"/v2/items", &payload); err != nil {
	}

	rawSlot, ok := payload[0]["item_slot_type"]
	if !ok {
		t.Fatalf("expevted item_slot_type key in payload")
	}

	slotType, ok := rawSlot.(string)
	if !ok {
		t.Fatalf("expevted item_slot_type to be string, got %T", rawSlot)
	}

	if slotType != "weapon" {
		t.Fatalf("expected item_slot_type weapon, got %q", slotType)
	}
}

func TestGetItemsPayload_ContainsCost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/items" {
			t.Fatalf("expevted path /v2/items, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"id": 1,
				"class_name": "weapon_core",
				"name": "Basic Magazine",
				"item_slot_type": "weapon",
				"cost": 800
			}
		]`))
	}))
	defer ts.Close()

	client := NewDeadlockClient("http://unused-base-url", ts.URL)

	var payload []map[string]any
	if err := client.(*deadlockClient).fetch(context.Background(), ts.URL+"/v2/items", &payload); err != nil {
		t.Fatalf("expected 1 item, got %d", len(payload))
	}

	rawCost, ok := payload[0]["cost"]
	if !ok {
		t.Fatalf("expevted cost key in payload, got Keys: %#v", payload[0])
	}

	cost, ok := rawCost.(float64)
	if !ok {
		t.Fatalf("expected cost to be numeric, got %T", rawCost)
	}

	if cost != 800 {
		t.Fatalf("expected cost 800, got %v", cost)
	}
}
