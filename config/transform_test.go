package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

// Mock YAML inputs for various test cases
var validYAML = []byte(`
options:
  confidence: 50 # default confidence level for all transformations unless otherwise specified

transformations:
  FQDN->IP:
    priority: 1
    confidence: 80
  FQDN->WHOIS:
    priority: 2
  FQDN->ALL: 
    exclude: [TLS,FQDN]
  IP->IP:
    priority: 1
    confidence: 80
  IP->WHOIS:
    priority: 2
  IP->TLS:
    # leaving both priority and confidence out

`)

var conflictingYAML = []byte(`
options:
  confidence: 50

transformations:
  FQDN->IP:
    priority: 1
    confidence: 80
  FQDN->none:
    priority: 2
  FQDN->ALL: 
    exclude: [TLS,FQDN]
  IP->IP:
    priority: 1
    confidence: 80
  IP->WHOIS:
    priority: 2
  IP->TLS:
    # leaving both priority and confidence out
`)

var invalidKeyYAML = []byte(`
options:
  confidence: 50

transformations:
  FQDN-IP:
    priority: 1
`)

// Utility function to unmarshal YAML and load transformation settings
func prepareConfig(yamlInput []byte) (*Config, error) {
	conf := NewConfig()
	err := yaml.Unmarshal(yamlInput, conf)
	if err != nil {
		return nil, err
	}
	err = conf.loadTransformSettings(conf)
	return conf, err
}

func TestLoadTransformSettings(t *testing.T) {
	// Test with valid YAML input
	t.Run("valid YAML and settings", func(t *testing.T) {
		conf, err := prepareConfig(validYAML)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if conf.Transformations["FQDN->WHOIS"].Confidence != 50 {
			t.Errorf("Expected confidence to be set to global value")
		}
		// Add debugging logs
		t.Logf("Configuration: %v", conf)
	})

	// Test with conflicting 'none' transformation
	t.Run("conflicting transformations", func(t *testing.T) {
		_, err := prepareConfig(conflictingYAML)
		if err == nil {
			t.Fatalf("Expected error due to conflicting 'none' transformation, got nil")
		}
		// Add debugging logs
		t.Logf("Error: %v", err)
	})

	// Test with invalid key format in YAML
	t.Run("invalid key format", func(t *testing.T) {
		_, err := prepareConfig(invalidKeyYAML)
		if err == nil {
			t.Fatalf("Expected error due to invalid key format, got nil")
		}
		// Add debugging logs
		t.Logf("Error: %v", err)
	})
}
