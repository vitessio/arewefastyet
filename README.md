# [Arewefastyet](https://benchmark.vitess.io)

Pull Request after Pull Request, the Vitess codebase changes a lot.
We must ensure that the performance of the codebase is not diminishing over time.
Arewefastyet automatically tests the performance of Vitess by benchmarking it using several workloads.
The performance is compared against the main branch, release branches and recent git tags, along with custom git SHA.

## Pull Request needing benchmarks

When someone wants to know if a Pull Request will affect the performance of Vitess, one might wish to benchmark it before merging it. This can be done by setting the `Benchmark me` label to your Pull Request.
Arewefastyet will then start benchmarking the head commit of your Pull Request and to compare against the Pull Request's base.

## How to run

Arewefastyet uses Docker and Docker Compose to easily run on any environment. You will need to install both tools before running arewefastyet.

Moreover, some secrets are required to run arewefastyet correctly which can be provided by a maintainer of Vitess.
Those secrets will allow you to connect to the arewefastyet database, to connect to the remote benchmarking server etc.

### Locally

```
docker compose build
docker compose up
```

### Production

```
docker compose -f docker-compose.prod.yml build
docker compose -f docker-compose.prod.yml up
```

### Formatting

```
npm run format
```