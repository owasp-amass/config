package config

import (
	"fmt"
	"strings"
)

// Transformation represents an individual transofmration with optional priority & confidence.
type Transformation struct {
	From       string   `yaml:"-" json:"-"`
	To         string   `yaml:"-" json:"-"`
	Priority   int      `yaml:"priority,omitempty" json:"priority,omitempty"`
	Confidence int      `yaml:"confidence,omitempty" json:"confidence,omitempty"`
	Exclude    []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

// loadTransformSettings processes the Transformations map from the configuration,
// assigning structured data to each Transformation based on its key.
// Each key is parsed into 'From' and 'To' segments, representing the origin and target
// of the transformation, respectively, which are then stored in the corresponding Transformation struct.
//
// Additionally, the function ensures that transformation rules are consistent:
// if a 'none' rule (indicating no transformation) is set for a 'From' type, no other
// transformation should be defined for the same 'From'. This is checked by maintaining
// two maps: one for 'From' types with 'none' transformations and another for those with valid ones.
//
// Additionally, it applies a global confidence level to all transformations that don't have a confidence value set.
//
// The function operates directly on the pointers within the Transformations map, modifying the
// original Transformation structs it points to. Proper error handling is included for nil pointers
// and invalid key formats.
//
// Returns an error if it encounters an invalid key format or conflicting transformation rules.
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

		// Prevent nil pointer dereference by checking if the transformation is nil.
		if transformation == nil {
			return fmt.Errorf("nil transformation for key: %s", key)
		}

		// Spit the key into 'From' and 'To' components.
		if err := transformation.Split(key); err != nil {
			return fmt.Errorf("error when splitting the key: %w", err)
		}

		// Apply the global confidence if no specific confidence is set for this transformation.
		if transformation.Confidence == 0 {
			transformation.Confidence = globalConfidence
		}

		err := c.ValidateTransform(transformation, fromWithValid, fromWithNone)
		if err != nil {
			return err
		}
	}

	// If the loop completes with no conflicts, the function returns nil, indicating success.
	return nil
}

// Split splits the key into 'From' and 'To' components, expecting a "->" delimiter.
// Requires a non-nil Transformation pointer and a valid key format. Example: FQDN->IP
func (t *Transformation) Split(key string) error {

	if t.From != "" && t.To != "" {
		return nil // Already split
	}

	// Split the key into 'From' and 'To' components, expecting a "->" delimiter.
	parts := strings.Split(key, "->")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format: %s", key)
	}

	// Assign the 'From' and 'To' values to the Transformation struct.
	t.From = parts[0]
	t.To = parts[1]
	return nil
}

// ValidateTransform checks the validity of a given transformation with respect to
// previously registered transformations. The function ensures that there are no conflicts
// between transformations with 'none' (indicating no action) and other valid transformations
// for the same 'From' type.
//
// Parameters:
// - transform: The Transformation object to be validated.
// - fromWithValid: A map tracking 'From' types that have at least one valid transformation defined.
// - fromWithNone: A map tracking 'From' types that have a 'none' transformation indicating no processing.
//
// The function works as follows:
//  1. If the 'To' field of the transformation is 'none', it checks if there are any valid transformations
//     for the same 'From' type. If such transformations exist, an error is returned.
//  2. For other valid transformations, it checks if a 'none' transformation is already registered
//     for the same 'From' type. If such a transformation exists, an error is returned.
//
// The function updates the maps fromWithValid and fromWithNone based on the validation results.
//
// Returns:
// - nil if the transformation is valid.
// - An error if a conflict is detected or the transformation is otherwise invalid.
func (c *Config) ValidateTransform(transform *Transformation, fromWithValid, fromWithNone map[string]bool) error {
	// Check for a 'none' transformation, which indicates that no further processing is required for this 'From' type.
	if transform.To == "none" {
		// Conflict arises if there's already a valid transformation for this 'From'.
		if fromWithValid[transform.From] {
			return fmt.Errorf("invalid config: 'none' specified after a valid transformation for 'From' type: %s. 'None' should be the only transformation", transform.From)
		}
		fromWithNone[transform.From] = true
	} else { // For other valid transformations.
		// Conflict arises if a 'none' transformation is already registered for this 'From'.
		if fromWithNone[transform.From] {
			return fmt.Errorf("invalid config: valid transformation specified after 'none' for 'From' type: %s. 'None' should be the only transformation", transform.From)
		}
		// Mark this 'From' as having a valid transformation.
		fromWithValid[transform.From] = true
	}

	return nil
}
