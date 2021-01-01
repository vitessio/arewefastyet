# Introduction 
The purpose of this project is to do a benchmark run when ever there is a push. The background activity is fairly simple, we create our own bare metal server. Once this server is created we run a bunch of ansibles(for sysbench) and once the run is complete we read the results and store them in a mysql instance. Once the following operations are complete we take down the server. 

### We use the Packet API to create and kill the bare metal server which we used to run the benchmarks on.

## Index 
1. Installation
2. [Api](Api.md) 
