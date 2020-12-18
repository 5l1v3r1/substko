# substko

![Made with Go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![Maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/substko.svg?style=flat&color=0040ff)](https://github.com/drsigned/substko/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/substko.svg?style=flat&color=0040ff)](https://github.com/drsigned/substko/issues?q=is:issue+is:closed) [![License](https://img.shields.io/badge/License-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/substko/blob/master/LICENSE) [![Author](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

substko is a subdomain takeover discovery tool written in Go. It takes in a list of FQDNs (Fully Qualified Domain Names) from stdin or via `-l` argument, then for each FQDN it:
* Checks for CNAME Subdomain Takeover
    * a dangling CNAME pointing to a non-existent domain name
    * a dangling CNAME pointing to a third party service that can be taken over.
* Checks for NS Subdomain Takeover

## Resources


* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From Source](#from-source)
    * [From Github](#from-github)
* [Usage](#usage)

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/drsigned/substko/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

substko requires **go1.14+** to install successfully. Run the following command to get the repo

```bash
$ GO111MODULE=on go get github.com/drsigned/substko/cmd/substko
```

#### From Github

```bash
$ git clone https://github.com/drsigned/substko.git; cd substko/cmd/substko/; go build; mv substko /usr/local/bin/; substko -h
```

## Usage

To display help message for substko use the `-h` flag:

```
$ substko -h

           _         _   _
 ___ _   _| |__  ___| |_| | _____
/ __| | | | '_ \/ __| __| |/ / _ \
\__ \ |_| | |_) \__ \ |_|   < (_) |
|___/\__,_|_.__/|___/\__|_|\_\___/ v1.1.0

USAGE:
  substko [OPTIONS]

OPTIONS:
  -c               concurrency level (default: 20)
  -f               path to fingerprints file
  -https           force HTTPS connection (default: false)
  -l               targets list
  -nc              no color mode (default: false)
  -silent          silent mode
  -timeout         HTTP timeout in seconds (default: 10)
  -u               download/update fingerprints
  -v               verbose mode

```

**Note:** domains can be also be provided by piping them into substko from stdin.
## Contribution

[Issues](https://github.com/drsigned/substko/issues) and [Pull Requests](https://github.com/drsigned/substko/pulls) are welcome!