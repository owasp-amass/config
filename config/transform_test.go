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
  FQDN->IPAddress:
    priority: 1
    confidence: 80
  FQDN->WHOIS:
    priority: 2
  FQDN->ALL: 
    exclude: [RIRORG,FQDN]
  IPAddress->IPAddress:
    priority: 1
    confidence: 80
  IPAddress->WHOIS:
    priority: 2
  IPAddress->RIRORG:
    # leaving both priority and confidence out

`)

var conflictingNoneYAML = []byte(`
options:
  confidence: 50

transformations:
  FQDN->IPAddress:
    priority: 1
    confidence: 80
  FQDN->none:
    priority: 2
  FQDN->ALL: 
    exclude: [TLS,FQDN]
  IPAddress->IPAddress:
    priority: 1
    confidence: 80
  IPAddress->WHOIS:
    priority: 2
  IPAddress->TLS:
    # leaving both priority and confidence out
`)

var conflictingNoneYAML2 = []byte(`
options:
  confidence: 50

transformations:
  FQDN->none:
    priority: 2
  FQDN->IPAddress:
    priority: 1
    confidence: 80
  FQDN->ALL: 
    exclude: [TLS,FQDN]
  IPAddress->IPAddress:
    priority: 1
    confidence: 80
  IPAddress->WHOIS:
    priority: 2
  IPAddress->TLS:
    # leaving both priority and confidence out
`)

var invalidKeyYAML = []byte(`
options:
  confidence: 50

transformations:
  FQDN-IPAddress:
    priority: 1
`)

var nonOAMtoYAML = []byte(`
options:
  confidence: 50 # default confidence level for all transformations unless otherwise specified

transformations:
  FQDN->IPAddress:
    priority: 1
    confidence: 80
  FQDN->Amass:
    priority: 2
  FQDN->ALL: 
    exclude: [RIRORG,FQDN]
`)

var nonOAMfromYAML = []byte(`
options:
  confidence: 50 # default confidence level for all transformations unless otherwise specified

transformations:
  FQDN->IPAddress:
    priority: 1
    confidence: 80
  Amass->WHOIS:
    priority: 2
  FQDN->ALL: 
    exclude: [RIRORG,FQDN]
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
	t.Run("conflicting transformations - none after", func(t *testing.T) {
		_, err := prepareConfig(conflictingNoneYAML)
		if err == nil {
			t.Fatalf("Expected error due to conflicting 'none' transformation, got nil")
		}
		// Add debugging logs
		t.Logf("Error: %v", err)
	})

	// Test with conflicting 'none' transformation
	t.Run("conflicting transformations - none before", func(t *testing.T) {
		_, err := prepareConfig(conflictingNoneYAML2)
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

	// Test with non-OAM compliant 'to' transformation
	t.Run("non-OAM compliant 'to' transformation", func(t *testing.T) {
		_, err := prepareConfig(nonOAMtoYAML)
		if err == nil {
			t.Fatalf("Expected error due to non-OAM compliant 'to' transformation, got nil")
		}
		// Add debugging logs
		t.Logf("Error: %v", err)
	})

	// Test with non-OAM compliant 'from' transformation
	t.Run("non-OAM compliant 'from' transformation", func(t *testing.T) {
		_, err := prepareConfig(nonOAMfromYAML)
		if err == nil {
			t.Fatalf("Expected error due to non-OAM compliant 'from' transformation, got nil")
		}
		// Add debugging logs
		t.Logf("Error: %v", err)
	})

}
func TestCheckTransformations(t *testing.T) {
	conf := NewConfig()
	conf.Transformations = map[string]*Transformation{
		"FQDN->IPAddress": {
			From:       "fqdn",
			To:         "ipaddress",
			Priority:   1,
			Confidence: 80,
		},
		"FQDN->WHOIS": {
			From:     "fqdn",
			To:       "whois",
			Priority: 2,
		},
		"FQDN->ALL": {
			From:    "fqdn",
			To:      "all",
			Exclude: []string{"tls", "fqdn"},
		},
		"WHOIS->ALL": {
			From:    "whois",
			To:      "all",
			Exclude: []string{"fqdn"},
		},
	}
	conf2 := NewConfig()
	conf2.Transformations = map[string]*Transformation{
		"FQDN->IPAddress": {
			From: "fqdn",
			To:   "ipaddress",
		},
	}

	tests := []struct {
		name       string
		from       string
		tos        []string
		expectErr  bool
		errMessage string
	}{
		{
			name:      "Valid transformation",
			from:      "fqdn",
			tos:       []string{"ipaddress"},
			expectErr: false,
		},
		{
			name:      "No match",
			from:      "fqdn",
			tos:       []string{"rirorg"},
			expectErr: false},
		{
			name:      "Transformation to 'all'",
			from:      "fqdn",
			tos:       []string{"tls", "rirorg"},
			expectErr: false,
		},
		{
			name:      "Transformation with excluded targets",
			from:      "fqdn",
			tos:       []string{"ip", "tls"},
			expectErr: false,
		},
		{
			name:      "No \"from\" matches with config",
			from:      "ip",
			tos:       []string{"tls", "rirorg"},
			expectErr: true,
		},
		{
			name:      "No \"to\" matches with config",
			from:      "whois",
			tos:       []string{"fqdn"},
			expectErr: true,
		},
		{
			name:      "Nil \"to\" matches with config",
			from:      "fqdn",
			tos:       []string{"rirorg"},
			expectErr: true,
		},
	}

	var err error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name != "Nil \"to\" matches with config" {
				err = conf.CheckTransformations(tt.from, tt.tos...)
			} else {
				err = conf2.CheckTransformations(tt.from, tt.tos...)
			}
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		expected  *Transformation
		expectErr bool
	}{
		{
			name:      "Valid key1",
			key:       "FQDN->IPAddress",
			expected:  &Transformation{From: "fqdn", To: "ipaddress"},
			expectErr: false,
		},
		{
			name:      "Valid key2",
			key:       "FQDN->IPAddress",
			expected:  &Transformation{From: "fqdn", To: "ipaddress"},
			expectErr: false,
		},
		{
			name:      "Invalid key delimiter",
			key:       "FQDN-IPAddress",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Empty key",
			key:       "",
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf := &Transformation{}
			if tt.name == "Valid key1" {
				tf.From = "FQDN"
				tf.To = "IPAddress"
			}
			err := tf.Split(tt.key)
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tf.From != tt.expected.From || tf.To != tt.expected.To {
					t.Errorf("Expected From: %s, To: %s, got From: %s, To: %s", tt.expected.From, tt.expected.To, tf.From, tf.To)
				}
			}
		})
	}
}

func TestCheckTransformResult(t *testing.T) {
	var result bool
	c := &Config{
		results: map[string]*resultTransform{
			"fqdn": {
				to: map[string]struct{}{
					"ipaddress": {},
					"whois":     {},
					"rirorg":    {},
				},
			},
			"ipaddress": {
				to: map[string]struct{}{
					"fqdn":   {},
					"whois":  {},
					"rirorg": {},
				},
			},
		},
	}
	cNil := &Config{}
	cNil2 := &Config{
		results: map[string]*resultTransform{
			"fqdn": {}}}

	tests := []struct {
		name     string
		from     string
		to       string
		expected bool
	}{
		{
			name:     "Existing transformation",
			from:     "fqdn",
			to:       "ipaddress",
			expected: true,
		},
		{
			name:     "Non-existing transformation",
			from:     "fqdn",
			to:       "tls",
			expected: false,
		},
		{
			name:     "Nil results",
			from:     "fqdn",
			to:       "ip",
			expected: false,
		},
		{
			name:     "Nil results \"to\"",
			from:     "fqdn",
			to:       "ip",
			expected: false,
		},
		{
			name:     "Nil results \"from\"",
			from:     "ip",
			to:       "fqdn",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Nil results" {
				result = cNil.CheckTransformResult(tt.from, tt.to)
			} else if tt.name == "Nil results \"to\"" {
				result = cNil2.CheckTransformResult(tt.from, tt.to)
			} else {
				result = c.CheckTransformResult(tt.from, tt.to)
			}
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
