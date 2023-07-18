// Copyright © by Jeff Foley 2017-2023. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestGetDataSourceConfig(t *testing.T) {
	name := "test"
	c := NewConfig()

	// Add a DataSource to the Config object to use it in the test
	c.DataSrcConfigs = &DataSourceConfig{
		Datasources: []DataSource{
			{
				Name: name,
			},
		},
	}

	if dsc := c.GetDataSourceConfig(""); dsc != nil {
		t.Errorf("GetDataSourceConfig returned a non-nil value when provided an invalid argument")
	}

	if dsc := c.GetDataSourceConfig(name); dsc == nil || dsc.Name != name {
		t.Errorf("GetDataSourceConfig returned an error when provided a valid argument")
	}
}

func TestAddCredentials(t *testing.T) {
	name := "test"
	c := NewConfig()

	// Add a DataSource to the Config object to use it in the test
	c.DataSrcConfigs = &DataSourceConfig{
		Datasources: []DataSource{
			{
				Name: name,
			},
		},
	}

	dsc := c.GetDataSourceConfig(name)

	if err := dsc.AddCredentials("account1", Credentials{Username: "username", Password: "password"}); err != nil {
		t.Errorf("AddCredentials returned an error when provided an valid arguments: %v", err)
	}

	if dsc.Creds["account1"].Username != "username" {
		t.Errorf("AddCredentials failed to enter the new credentials into the data source configuration")
	}
}

func TestGetCredentials(t *testing.T) {
	c := NewConfig()

	// Add a DataSource with credentials to the Config object to use it in the test
	c.DataSrcConfigs = &DataSourceConfig{
		Datasources: []DataSource{
			{
				Name: "test",
				Creds: map[string]Credentials{
					"account1": {
						Username: "username",
						Password: "password",
					},
				},
			},
		},
	}
	dsc := c.DataSrcConfigs

	// Pass the name of the data source when calling GetCredentials
	if creds := dsc.GetCredentials("test"); creds == nil || creds.Username != "username" {
		t.Errorf("GetCredentials returned an error when provided a valid argument")
	}
}

func TestLoadDataSourceSettings(t *testing.T) {
	c := NewConfig()
	ymlData := `
datasources:
  - name: AlienVault
    ttl: 4320
    creds:
      account1:
        username: avuser
        password: avpass
  - name: BinaryEdge
    creds:
      account2:
        username: beuser
        password: bepass
global_options:
  minimum_ttl: 1440
`
	var dsConfig DataSourceConfig
	err := yaml.Unmarshal([]byte(ymlData), &dsConfig)
	if err != nil {
		t.Errorf("Failed to parse the data source settings: %v", err)
	}

	// Assign the unmarshalled DataSourceConfig to the Config struct
	c.DataSrcConfigs = &dsConfig

	dsc := c.GetDataSourceConfig("AlienVault")
	if dsc == nil {
		t.Errorf("Failed to load data source settings")
	}

	// Pass the name of the data source when calling GetCredentials
	if creds := c.DataSrcConfigs.GetCredentials("AlienVault"); creds == nil || creds.Username != "avuser" {
		t.Errorf("Failed to load data source credentials")
	}

	if c.DataSrcConfigs.GlobalOptions["minimum_ttl"] != 1440 {
		t.Errorf("Failed to load global options")
	}
}
