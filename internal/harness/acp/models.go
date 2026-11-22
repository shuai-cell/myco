package acp

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/chenhg5/myco/internal/harness"
)

type modelOption struct {
	Value   string        `json:"value"`
	Name    string        `json:"name"`
	Options []modelOption `json:"options"`
}

type sessionConfig struct {
	SessionID     string            `json:"sessionId"`
	ConfigOptions []json.RawMessage `json:"configOptions"`
}

type selectConfig struct {
	ID      string        `json:"id"`
	Current string        `json:"currentValue"`
	Options []modelOption `json:"options"`
}

func (config sessionConfig) models() (*harness.ModelCatalog, map[string]string) {
	for _, raw := range config.ConfigOptions {
		var option selectConfig
		if json.Unmarshal(raw, &option) != nil || option.ID != "model" {
			continue
		}
		catalog := &harness.ModelCatalog{}
		values := make(map[string]string)
		add := func(choice modelOption, group string) {
			if choice.Value == "" {
				return
			}
			digest := sha256.Sum256([]byte(choice.Value))
			id := fmt.Sprintf("dsh:%x", digest[:16])
			if _, ok := values[id]; ok {
				return
			}
			values[id] = choice.Value
			name := choice.Name
			if name == "" {
				name = choice.Value
			}
			if group != "" {
				name = name + " · " + group
			}
			if len(name) > 80 {
				name = name[:80]
				for !utf8.ValidString(name) {
					name = name[:len(name)-1]
				}
			}
			catalog.Choices = append(catalog.Choices, harness.ModelChoice{ID: id, Name: name})
			if choice.Value == option.Current {
				catalog.Current = id
				last := len(catalog.Choices) - 1
				catalog.Choices[0], catalog.Choices[last] = catalog.Choices[last], catalog.Choices[0]
			}
		}
		for _, item := range option.Options {
			if len(item.Options) == 0 {
				add(item, "")
			}
			for _, choice := range item.Options {
				add(choice, item.Name)
			}
		}
		return catalog, values
	}
	return nil, nil
}
