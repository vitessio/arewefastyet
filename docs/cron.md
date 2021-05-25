# Arewefastyet cron and regression testing

## Sources
In arewefastyet, we run benchmark for multiple reasons, these reasons are named as “sources”. We use the four following sources:

- New commits to the master branch: **cron**
- New tags: **cron_tags_{tag_name}**
- New commits to release branches: **cron_release_{release_name}**
- New commits to pull requests: **cron_pr**
- Base commit for a pull request: **cron_pr_base**

## Cron and Schedule
For each of the listed sources, we initiate a cron job that will periodically create new benchmarks. 
The cron schedule used for those jobs is defined through one of the CLI’s flags, which defaults to every day at midnight.

## Regressions and Notifications
After running and analyzing a benchmark, we can determine that the result is a regression. 
However, regression will be evaluated differently based on the benchmark’s source. 
The different ways of evaluating a regression are:


| Source | How to compare results |
| ---- | -------------- |
| **cron** |  Compare result against latest **cron** and latest release |
| **cron_tags** | None |
| **cron_release** | Compare result against the previous tag on this release branch and last **cron_release** |
| **cron_pr** | Compare result against the pull request base’s reference refered by **cron_pr_base** |

After the comparison and if we have detected a regression, we send a notification on the dedicated Slack channel. 
The notification includes the performance difference calculated in %, the benchmark UUID, and the 
different commit SHAs used for the comparison.

Notification is always sent upon regression, however, cron_pr benchmarks will always issue a 
new notification, the notification will be formatted based on whether we have a regression or not.

