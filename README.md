# Billables Buddy

[![CI](https://github.com/jordanleven/billables-buddy/actions/workflows/ci.yml/badge.svg)](https://github.com/jordanleven/billables-buddy/actions/workflows/ci.yml)

Billables Buddy is the easiest way to get a consolidated look at your tracked hours. Designed to be used with [xbar], it will automatically update your tracked hours against your forecasted schedule to calculate your progress towards your scheduled hours.

![BillablesBuddyHero](/assets/BillablesBuddyHero.png)

## Features

Using your schedule in Forecast, and your tracked hours in Harvest, Billables Buddy has the following features:

- Displays your current status in the menu bar to tell you whether you're on-track, ahead, falling behind, or over billable hours.
- Tells you the time of your estimated end-of-day based on tracked hours.
- A tracking percentage shows your actual to expected billable hours, along with the delta. Need more detail? Click on it to go directly to the week view in Harvest.
- Gives you a breakdown of your total, billable, and non-billable hours.
- Shows your expected hours that are updated throughout the day based on your starting time.
- Automatically updates every five minutes to keep your data up-to-date.

### The Status

The status is the heart of Billables Buddy and is designed to give you a clear indication of your tracking against expected billable hours. Based on your progress, it will show one of three emoji:

Emoji | Status
------------ | -------------
:white_check_mark: | On-track, ahead of, or within fifteen minutes of current billable expectations
:x: | At least 15 minutes behind current billable expectations
:stop_sign: | Over total billable expectations

## Getting Started

1. To start, make sure you have [Go] downloaded to your machine. You'll need it to compile your personalized version of Billables Buddy.
1. Next, copy the contents of `.env-sample` to `.env`. This is where we'll store your authentication tokens.
1. Next, generate your API credentials at [Harvest][harvest_api].
1. Lastly, after creating API tokens, copy each of the three credentials in the newly-created `.env` file attributed to the appropriate variable.

### Running locally

To run locally, simply run `go run .` in your command line. This will output your statistics in the CLI.

### Building and deploying xbar Plugin

Run `./build.sh` to generate your binary plugin. The compiled version will be output as `billablesbuddy.1m.goc` and automatically copied to the xbar plugins directory.

## Definitions

**Actual hours**\
Your actual hours are the current summation of billable, non-billable, or total hours based on entries in Harvest.

**Expected hours**\
Your expected hours are the tracking summation of billable, non-billable, or total hours based on the current date and time. These are calculated throughout the day using the schedule in Forecast and your progression through an assumed eight-hour workday (spanning a total of nine hours) starting at the earliest entry in Harvest.

Since expected hours are calculating using this "8 working hours over a nine hour workday" method to normally distribute expectations, Billables Buddy will only expect 89% (8/9) of any working minute to be tracked. This means that over the course of 60 minutes, Billables Buddy will only expect around 53 minutes to be logged to either a billable or non-billable project.

In the graph below, you can see a graphic representation of how this works. The user begins work at 9:00, and as they continue to log 100% of their time (which is 11% more than Billables Buddy would expect) the delta between actual and expect hours will drift further apart as time goes on and the user gets ahead of their expected hours. Eventually, when the user takes a break, the delta will decrease, and then converge with the expected hours once the workday is completed.

![ExpectedHoursGraph](/assets/ExpectedHoursGraph.png)

[xbar]: https://github.com/matryer/xbar
[Go]: https://golang.org/doc/install
[harvest_api]: https://id.getharvest.com/developers
