ssllabs
==========

[![GitHub release](https://img.shields.io/github/release/keltia/ssllabs.svg)](https://github.com/keltia/ssllabs/releases)
[![GitHub issues](https://img.shields.io/github/issues/keltia/ssllabs.svg)](https://github.com/keltia/ssllabs/issues)
[![Go Version](https://img.shields.io/badge/go-1.10-blue.svg)](https://golang.org/dl/)
[![Build Status](https://travis-ci.org/github.com/keltia/ssllabs.svg?branch=master)](https://travis-ci.org/github.com/keltia/ssllabs)
[![GoDoc](http://godoc.org/github.com/keltia/ssllabs?status.svg)](http://godoc.org/github.com/keltia/ssllabs)
[![SemVer](http://img.shields.io/SemVer/2.0.0.png)](https://semver.org/spec/v2.0.0.html)
[![License](https://img.shields.io/pypi/l/Django.svg)](https://opensource.org/licenses/BSD-2-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/keltia/ssllabs)](https://goreportcard.com/report/github.com/keltia/ssllabs)

Go wrapper for the [SSLLabs](https://ssllabs.com/) API for testing TLS parameters of a given website.

## Requirements

* Go >= 1.10

`keltia.net/ssllabs` is a Go module (you can use either Go 1.10 with `vgo` or 1.11+).  The API exposed follows the Semantic Versioning scheme to guarantee a consistent API compatibility.

## USAGE

There is a small example program included in `cmd/ssllabs` to either show the grade of a given site or JSON dump of the detailed report.

You can use [`jq`](https://stedolan.github.io/jq/) to display the output of `ssllabs -d <site>` in a colorised way:

    ssllabs -d www.ssllabs.com | jq .

## API Usage

As with many API wrappers, you will need to first create a client with some optional configuration, then there are two main functions:

``` go
    // Simplest way
    c, _ := ssllabs.NewClient()
    grade, err := c.GetScore("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }


    // With some options, timeout at 15s, caching for 10s and debug-like verbosity
    cnf := ssllabs.Config{
        Timeout:15,
        Retries:3,
        Log:2,
    }
    c, err := ssllabs.NewClient(cnf)
    report, err := c.GetScore("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
```

For the `GetDetailedReport()` call, the raw JSON object will be returned (and presumably handled by `jq`).

``` go
    // Simplest way
    c, _ := ssllabs.NewClient()
    report, err := c.GetDetailedReport("example.com")
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    fmt.Printf("Full report:\n%v\n", report)
```


OPTIONS

| Option  | Type | Description |
| ------- | ---- | ----------- |
| Timeout | int  | time for connections (default: 10s) |
| Log     | int  | 1: verbose, 2: debug (default: 0) |
| Retries | int  | Number of retries when not FINISHED (default: 5) |
| Refresh | bool | Force refresh of the sites (default: false) |


## Using behind a web Proxy

Dependency: proxy support is provided by my `github.com/keltia/proxy` module.

UNIX/Linux:

```
    export HTTP_PROXY=[http://]host[:port] (sh/bash/zsh)
    setenv HTTP_PROXY [http://]host[:port] (csh/tcsh)
```

Windows:

```
    set HTTP_PROXY=[http://]host[:port]
```

The rules of Go's `ProxyFromEnvironment` apply (`HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY`, lowercase variants allowed).

If your proxy requires you to authenticate, please create a file named `.netrc` in your HOME directory with permissions either `0400` or `0600` with the following data:

    machine proxy user <username> password <password>

and it should be picked up. On Windows, the file will be located at

    %LOCALAPPDATA%\ssllabs\netrc

## License

The [BSD 2-Clause license](https://github.com/keltia/ssllabs/blog/master/LICENSE.md).

# Contributing

This project is an open Open Source project, please read `CONTRIBUTING.md`.

# References

[SSLLabs API documentation](https://github.com/ssllabs/ssllabs-scan/blob/master/ssllabs-api-docs-v3.md)

# Feedback

We welcome pull requests, bug fixes and issue reports.

Before proposing a large change, first please discuss your change by raising an issue.