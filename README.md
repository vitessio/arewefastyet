# [Arewefastyet](https://benchmark.vitess.io)

## Background

Pull Request after Pull Request, the Vitess codebase changes a lot.
We must ensure that the performance of the codebase is not diminishing over time.
Arewefastyet automatically tests the performance of Vitess by benchmarking it using several workloads.
The performance is compared against the main branch, release branches and recent git tags, along with custom git SHA.

### Pull Request needing benchmarks

When someone wants to know if a Pull Request will affect the performance of Vitess, one might wish to benchmark it before merging it. This can be done by setting the `Benchmark me` label to your Pull Request.
Arewefastyet will then start benchmarking the head commit of your Pull Request and to compare against the Pull Request's base.
