# Capture

[![Build Status](https://github.com/ifabos/capture/actions/workflows/go.yml/badge.svg)](https://github.com/ifabos/capture/actions/workflows/go.yml)
[![LICENSE](https://img.shields.io/github/license/horsing/coder.svg)](https://github.com/ifabos/capture/blob/master/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ifabos/capture)](https://goreportcard.com/report/github.com/ifabos/capture)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2761/badge)](https://bestpractices.coreinfrastructure.org/projects/6232)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/ifabos/capture/badge)](https://securityscorecards.dev/viewer/?uri=github.com/ifabos/capture)
[![Codecov](https://img.shields.io/codecov/c/github/horsing/coder?style=flat-square&logo=codecov)](https://codecov.io/gh/horsing/coder)
[![CLOMonitor](https://img.shields.io/endpoint?url=https://clomonitor.io/api/projects/cncf/chubao-fs/badge)](https://clomonitor.io/projects/cncf/chubao-fs)
[![Release](https://img.shields.io/github/v/release/horsing/coder.svg?color=161823&style=flat-square&logo=smartthings)](https://github.com/ifabos/capture/releases)
[![Tag](https://img.shields.io/github/v/tag/horsing/coder.svg?color=ee8936&logo=fitbit&style=flat-square)](https://github.com/ifabos/capture/tags)

## Overview

Capture is a tool designed to monitor and record TCP packets transmitted over a network. It is useful for network administrators, security experts, and developers who need to analyze network traffic, diagnose network issues, and perform security audits. This program is written in Go, offering efficiency, scalability, and ease of configuration.

## Usage

### Install

```bash
go install github.com/ifabos/capture@latest # using @<version> to try development features
```

### Introduction

```text
Usage: capture
```

### Start default

```bash
capture
```

## License

Capture is licensed under the [MIT](https://opensource.org/license/mit).
For detail see [LICENSE](LICENSE).

## Note

The master branch may be in an unstable or even broken state during development. Please use releases instead of the
master branch in order to get a stable set of binaries.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ifabos/capture&type=Date)](https://star-history.com/#ifabos/capture&Date)
