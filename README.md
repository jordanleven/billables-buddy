# Billables Buddy

[![CI](https://github.com/jordanleven/billables-buddy/actions/workflows/ci.yml/badge.svg)](https://github.com/jordanleven/billables-buddy/actions/workflows/ci.yml)

A simple project that will evaluate weekly progress towards Sparkbox billable hours to be used with [BitBar].

## Getting Started

1. To start, make sure you have [Go] downloaded to your machine. You'll need it in order to compile your personalized version of Billables Buddy.
1. Next, copy the contents of `.env-sample` to `.env`. This is where we'll store your personal authentication tokens.
1. Next, generate your API credentials at [Harvest][harvest_api].
1. Lastly, after creating API tokens, copy each of the three credentials in the newly-created `.env` file attributed to the appropriate variable.

## Running locally

To run locally, simply run `go run .` in your command line. This will output your stats in the CLI.

## Building and deploying BitBar Plugin

Run `./build` to generate your binary package. The compiled version will be output as `billablesbuddy.1m.goc` and automatically copied to the BitBar plugins directory.

[BitBar]: https://github.com/matryer/xbar
[Go]: https://golang.org/doc/install
[harvest_api]: https://id.getharvest.com/developers
