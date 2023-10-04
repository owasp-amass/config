# oam_i2y Users' Guide
![Network graph](../images/network_06092018.png "Amass Network Mapping")

----

The `oam_i2y` command is primarily used to convert legacy INI configuration files into the new YAML format. This is essential for users who want to transition from an older version of amass to a newer one while retaining their configurations. The YAML configuration will be compliant with all future developments in regard to OAM. It is in the user's best interest to transition to YAML for better compatibility and support.

To view documentation for the configuration file, please refer to the [Configuration Users' Guide](./user_guide.md).

## oam_i2y Usage table

The following table shows the usage of `oam_i2y`:

| Flag | Description | Example |
|------|-------------|---------|
| -ini | Path to the INI configuration file | oam_i2y -ini config.ini |
| -cf  | YAML configuration file name (default = oam_config.yaml) | oam_i2y -ini config.ini -cf example_config.yaml |
| -df  | YAML data sources file name (default = oam_datasources.yaml ) | oam_i2y -ini config.ini -df example_datasources.yaml

Users can also specify the file path (doesn't have to be just the file name) in the arguments so it doesn't have to be in the cwd.

Example:
```bash
oam_i2y -ini config.ini -cf ../../config.yaml -df datasrc.yaml
```

## oam_i2y Examples

The following examples show how to use `oam_i2y` and pointers that users should be aware of:

**Example 1: Basic Conversion**
```bash
oam_i2y -ini config.ini
```
This command will convert the given `config.ini` file into a compliant YAML format. The resultant files will be named `oam_config.yaml` and `oam_datasources.yaml` by default.

**Example 2: Specifying Output YAML Configuration File Name**
```bash
oam_i2y -ini config.ini -cf new_config.yaml
```
Here, the INI configuration `config.ini` is converted into `new_config.yaml`, allowing users to name the resultant YAML file as desired. If there are any data source credentials in the INI, it will output into a file called `oam_datasources.yaml` in the current working directory. 

**Example 3: Specifying Paths for Both Configuration and Data Sources**
```bash
oam_i2y -ini config.ini -cf ../myconfigs/new_config.yaml -df ../mydatasources/data.yaml
```
In addition to specifying the output name for the configuration, this command also defines where the data sources YAML file should be stored. Both paths are relative to the current working directory.

<span style="font-size: 1.25em;">**Remember that when using `oam_i2y`, the path to the data source configuration file will already be populated in the `datasource` value of the configuration file. Unless you're relocating the data source file, there's no need to adjust the `datasource` value in your configuration.** </span>
