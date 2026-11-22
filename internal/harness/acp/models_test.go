package acp

import (
	"encoding/json"
	"testing"
)

func TestGroupedAndFlatModelChoices(t *testing.T) {
	var state sessionConfig
	err := json.Unmarshal([]byte(`{"configOptions":[{"id":"other","currentValue":true},{"id":"model","type":"select","currentValue":"route beta","options":[{"value":"alpha","name":"Alpha"},{"group":"p","name":"Provider","options":[{"value":"route beta","name":"Beta"}]}]}]}`), &state)
	if err != nil {
		t.Fatal(err)
	}
	catalog, values := state.models()
	if len(catalog.Choices) != 2 || catalog.Choices[0].Name != "Beta · Provider" || values[catalog.Current] != "route beta" || catalog.Choices[0].ID != catalog.Current {
		t.Fatalf("unexpected catalog: %+v", catalog)
	}
}
