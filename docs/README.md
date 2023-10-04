# The OAM Configuration Repository

<p align="center">
  <img src="https://github.com/owasp-amass/amass/blob/master/images/amass_video.gif">
</p>

[![GoDoc](https://pkg.go.dev/badge/github.com/owasp-amass/config?utm_source=godoc)](https://pkg.go.dev/github.com/owasp-amass/config/config)
[![Follow on Twitter](https://img.shields.io/twitter/follow/owaspamass.svg?logo=twitter)](https://twitter.com/owaspamass)
[![Chat on Discord](https://img.shields.io/discord/433729817918308352.svg?logo=discord)](https://discord.gg/HNePVyX3cp)

This configuration file repo serves the purpose of parsing the configuration used for tools under the OWASP AMASS framework.

`oam_i2y` also resides in this repository. It is a tool that converts the legacy INI configuration file into two YAML files that can be used by the framework.

## Users' Guide
For a more detailed guide on using the `configuration file` and `oam_i2y` as an OAM user, please check out:
- [Configuration Users' Guide](./user_guide.md)
- [oam_i2y Users' Guide](./oam_i2y_user_guide.md)

## Installing oam_i2y [![Go Version](https://img.shields.io/github/go-mod/go-version/owasp-amass/config)](https://golang.org/dl/) 

### From Source

1. Install [Go](https://golang.org/doc/install) and setup your Go workspace
2. Download `oam_i2y` by running `go install -v github.com/owasp-amass/config/cmd/...@master` or `go install -v github.com/owasp-amass/config/cmd/oam_i2y@master`
3. At this point, the binary should be in `$GOPATH/bin`

### Local Install

1. Install [Go](https://golang.org/doc/install) and setup your Go workspace
2. Use git to clone the repository: `git clone https://github.com/owasp-amass/config`
    - At this point, a directory called `config` should be made
3. Go into the `config` directory by running `cd config`, and then build the desired program by running `go build ./cmd/oam_i2y`
4. **Enjoy!** The binary will reside in your current working directory, which should be the `config` directory.

## Corporate Supporters

[![ZeroFox Logo](../images/zerofox_logo.png)](https://www.zerofox.com/) [![IPinfo Logo](../images/ipinfo_logo.png)](https://ipinfo.io/) [![WhoisXML API Logo](../images/whoisxmlapi_logo.png)](https://www.whoisxmlapi.com/)

## Contributing [![Contribute Yes](https://img.shields.io/badge/contribute-yes-brightgreen.svg)](./CONTRIBUTING.md) [![Chat on Discord](https://img.shields.io/discord/433729817918308352.svg?logo=discord)](https://discord.gg/HNePVyX3cp)

We are always happy to get new contributors on board! Join our [Discord Server](https://discord.gg/HNePVyX3cp) to discuss current project goals.

## Troubleshooting [![Chat on Discord](https://img.shields.io/discord/433729817918308352.svg?logo=discord)](https://discord.gg/HNePVyX3cp)

If you need help with the usage of the configuration, please join our [Discord server](https://discord.gg/HNePVyX3cp) where community members can best help you.

**Please avoid opening GitHub issues for support requests or questions!**

## Licensing [![License](https://img.shields.io/badge/license-apache%202-blue)](https://www.apache.org/licenses/LICENSE-2.0)

This program is free software: you can redistribute it and/or modify it under the terms of the [Apache license](LICENSE). OWASP Amass and any contributions are Copyright © by Jeff Foley 2017-2023. Some subcomponents have separate licenses.