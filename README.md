# Billables Buddy

[![CI](https://github.com/jordanleven/billables-buddy/actions/workflows/ci.yml/badge.svg)](https://github.com/jordanleven/billables-buddy/actions/workflows/ci.yml)

Billables Buddy is the easiest way to get a consolidated look at your tracked hours. Designed to be used with [xbar], it will automatically update your tracked hours against your forecasted schedule to calculate your progress towards your scheduled hours.

![BillablesBuddyHero](/assets/BillablesBuddyHero.png)

## Features

Using your schedule in Forecast, and your tracked hours in Harvest, Billables Buddy has the following features:

1. Displays your current status in the menu bar to tell you whether you're on track, ahead, falling behind, or over billable hours.
1. Gives you a breakdown of your total, billable, and non-billable hours.
1. Shows your expected hours that are updated throughout the day based on your starting time.
1. Automatically updates every five minutes to keep your data up-to-date.

## Getting Started

1. To start, make sure you have [Go] downloaded to your machine. You'll need it to compile your personalized version of Billables Buddy.
1. Next, copy the contents of `.env-sample` to `.env`. This is where we'll store your authentication tokens.
1. Next, generate your API credentials at [Harvest][harvest_api].
1. Lastly, after creating API tokens, copy each of the three credentials in the newly-created `.env` file attributed to the appropriate variable.

## Running locally

To run locally, simply run `go run .` in your command line. This will output your statistics in the CLI.

## Building and deploying xbar Plugin

Run `./build.sh` to generate your binary plugin. The compiled version will be output as `billablesbuddy.1m.goc` and automatically copied to the xbar plugins directory.

[xbar]: https://github.com/matryer/xbar
[Go]: https://golang.org/doc/install
[harvest_api]: https://id.getharvest.com/developers
