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

You can download the pre-built binary for your platform from this repository's [Releases](https://github.com/drsigned/substko/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

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

```bash
$ substko -h
```

| Flag                     | Description                              | Example                              |
| :----------------------- | :--------------------------------------- | :----------------------------------- |
| -c, --concurrency         | concurrency level (default: 20)         | `substko -l subdomains.txt -c 100`               |
| -f, --fingerprints        | path to fingerprints file               | `substko -l subdomains.txt -f fingerprints.json` |
| --https                   | force HTTPS connection (default: false) | `substko -l subdomains.txt --https`              |
| -l, --list                | targets list                            | `substko -l subdomains.txt`                      |
| -nc, --no-color           | no color mode (default: false)          | `substko -l subdomains.txt -nc`                  |
| -s, --silent              | silent mode                             | `substko -l subdomains.txt -s`                   |
| -t, --timeout             | HTTP timeout in seconds (default: 10)   | `substko -l subdomains.txt -t 10`                |
| -u, --update-fingerprints | download/update fingerprints            | `substko -u` |
| -v, --verbose             | verbose mode                            | `substko -l subdomains.txt -v`                   |

**Note:** domains can be also be provided by piping them into substko from stdin.
## Contribution

[Issues](https://github.com/drsigned/substko/issues) and [Pull Requests](https://github.com/drsigned/substko/pulls) are welcome!