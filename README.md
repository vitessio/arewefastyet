# [Arewefastyet](https://benchmark.vitess.io)
## Background

With the codebase of Vitess becoming larger and complex changes getting merged, we need to ensure our changes are not degrading the performance of Vitess.

## Benchmarking Tool

To solve the aforementioned issue, we use a tool named arewefastyet that automatically tests the performance of Vitess. The performance are measured through a set of benchmarks divided into two categories: `micro` and `macro`, the former focuses on unit-level functions, and the latter targets system-wide performance changes.

The GitHub repository where lies all of arewefastyet's code can be found [here: vitessio/arewefastyet](https://github.com/vitessio/arewefastyet).

## CRON Schedule

Our benchmarks run frequently based on three different CRON schedules that are defined in [this file](https://github.com/vitessio/arewefastyet/blob/main/config/prod/config.yaml) under the `web-cron-*` keys.

### Pull Request needing benchmarks

When a pull request affect the performance of Vitess, one might wish to benchmark it before merging it. This can be done by setting the `Benchmark me` label to your pull request.
The corresponding CRON schedule will be used to start benchmarking the head commit of your pull request and to compare against the pull request's base.

## Website

The performances of Vitess can be observed throughout different releases, git SHAs, and nightly builds on arewefastyet's website at [https://benchmark.vitess.io](https://benchmark.vitess.io).

The website lets us:

* See previous benchmarks.
* Search results for a specific git SHA.
* Compare two results for two git SHAs.
* See micro and macro benchmarks results throughout different releases.
* Compare performance between VTGate's v3 planner and Gen4 planner.