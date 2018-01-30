# crypto-celeb-monitor
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)]()
[![Build Status](https://travis-ci.org/typenil/crypto-celeb-monitor.svg?branch=master)](https://travis-ci.org/typenil/crypto-celeb-monitor)
[![Test Coverage](https://api.codeclimate.com/v1/badges/368430c7858f2a9afaac/test_coverage)](https://codeclimate.com/github/typenil/crypto-celeb-monitor/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/368430c7858f2a9afaac/maintainability)](https://codeclimate.com/github/typenil/crypto-celeb-monitor/maintainability)

Golang app to detect new celebrities on CryptoCelebrities.co.

Personally, I set up a cron job that runs every minute. There's a post request to https://pushover.net/ that I use to send myself mobile push notifications. You have to set `alertToken` and `alertUser` environment variables to get that to work.

It all works quite well, but it's pretty pointless when every transaction you send to CryptoCelebrities fails unless the price is already too high for anyone to be interested.

Here's this really specific software that isn't nearly as profitable as I hoped it would be.
