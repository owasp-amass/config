package config

import (
	"errors"
	"fmt"
	"strings"

	oam "github.com/owasp-amass/open-asset-model"
)

// Transformation represents an individual transofmration with optional priority & confidence.
type Transformation struct {
	From       string   `yaml:"-" json:"-"`
	To         string   `yaml:"-" json:"-"`
	Priority   int      `yaml:"priority,omitempty" json:"priority,omitempty"`
	Confidence int      `yaml:"confidence,omitempty" json:"confidence,omitempty"`
	Exclude    []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

/*
loadTransformSettings processes the Transformations map from the configuration,
assigning structured data to each Transformation based on its key.
Each key is parsed into 'From' and 'To' segments, representing the origin and target
of the transformation, respectively, which are then stored in the corresponding Transformation struct.
*/
func (c *Config) loadTransformSettings(cfg *Config) error {
	// Map to track 'From' types that have a 'none' transformation, indicating no processing should occur.
	fromWithNone := make(map[string]bool)

	// Map to track 'From' types that have at least one valid transformation defined.
	fromWithValid := make(map[string]bool)

	// Retrieve the global confidence from the Options, if it's set.
	var globalConfidence int
	if gc, ok := c.Options["confidence"]; ok {
		globalConfidence, _ = gc.(int) // Assume it's an int; ignore the error otherwise.
	}

	// Iterate through each transformation rule defined in the configuration.
	for key, transformation := range c.Transformations {

		// Initialize transformation if nil
		if transformation == nil {
			transformation = &Transformation{}      // default struct
			c.Transformations[key] = transformation // assign it back to the map
		}

		// Spit the key into 'From' and 'To' components.
		if err := transformation.Split(key); err != nil {
			return fmt.Errorf("error when splitting the key: %w", err)
		}

		// Apply the global confidence if no specific confidence is set for this transformation.
		if transformation.Confidence == 0 {
			transformation.Confidence = globalConfidence
		}

		err := transformation.Validate(fromWithValid, fromWithNone)
		if err != nil {
			return err
		}
	}

	// If the loop completes with no conflicts, the function returns nil, indicating success.
	return nil
}

// Split splits the key into 'From' and 'To' components, expecting a "->" delimiter.
// Requires a non-nil Transformation pointer and a valid key format. Example: FQDN->IPaddress
func (t *Transformation) Split(key string) error {
	if t.From != "" && t.To != "" {
		t.From = strings.ToLower(t.From)
		t.To = strings.ToLower(t.To)
		return nil // Already split
	}

	// Split the key into 'From' and 'To' components, expecting a "->" delimiter.
	parts := strings.Split(key, "->")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key delimiter: %s", key)
	}

	// Assign the 'From' and 'To' values to the Transformation struct.
	t.From = strings.ToLower(parts[0])
	t.To = strings.ToLower(parts[1])
	return nil
}

/*
ValidateTransform checks the validity of a given transformation with respect to OAM &
previously registered transformations. The function ensures OAM compliance & that there are no conflicts
between transformations with 'none' (indicating no action) and other valid transformations
for the same 'From' type.
*/
func (t *Transformation) Validate(fromWithValid, fromWithNone map[string]bool) error {
	tfound := false
	ffound := false
	// Check if "From" and "To" is OAM compliant
	for _, a := range oam.AssetList {
		a := strings.ToLower(string(a))
		if t.From == a {
			ffound = true
		}
		if t.To == a || t.To == "none" || t.To == "all" {
			tfound = true
		}
		// Used to prevent unnecessary iterations
		if tfound && ffound {
			break
		}
	}

	if !ffound {
		return fmt.Errorf("invalid 'From' type: %s does not comply with OAM", t.From)
	}
	if !tfound {
		return fmt.Errorf("invalid 'To' type: %s does not comply with OAM", t.To)
	}

	// Check for a 'none' transformation, which indicates that no further processing is required for this 'From' type.
	if t.To == "none" {
		// Conflict arises if there's already a valid transformation for this 'From'.
		if fromWithValid[t.From] {
			return fmt.Errorf("invalid config: 'none' specified after a valid transformation for 'From' type: %s. 'None' should be the only transformation", t.From)
		}
		fromWithNone[t.From] = true
	} else { // For other valid transformations.
		// Conflict arises if a 'none' transformation is already registered for this 'From'.
		if fromWithNone[t.From] {
			return fmt.Errorf("invalid config: valid transformation specified after 'none' for 'From' type: %s. 'None' should be the only transformation", t.From)
		}
		// Mark this 'From' as having a valid transformation.
		fromWithValid[t.From] = true
	}

	return nil
}

// CheckTransformations checks if the given 'From' type has a valid transformation to any of the given 'To' types.
func (c *Config) CheckTransformations(from string, tos ...string) (map[string]struct{}, error) {
	lower := strings.ToLower(from)
	tomap := make(map[string]struct{})
	results := make(map[string]struct{})

	for _, v := range tos {
		t := strings.ToLower(v)
		tomap[t] = struct{}{}
	}

	for _, transform := range c.Transformations {
		if lower == transform.From {
			if transform.To == "all" {
				excludes := make(map[string]struct{})
				for _, e := range transform.Exclude {
					excludes[strings.ToLower(e)] = struct{}{}
				}

				for k := range tomap {
					if _, found := excludes[k]; !found {
						results[k] = struct{}{}
					}
				}
				continue
			} else if _, found := tomap[transform.To]; found {
				results[transform.To] = struct{}{}
			}
		}
	}

	if len(results) == 0 {
		return nil, errors.New("zero transformation matches in the session config")
	}
	return results, nil
}
