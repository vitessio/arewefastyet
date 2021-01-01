# Api 
Run API server (must have api key in header eg: curl -X GET 'http://127.0.0.1:5000/allresults' -H 'api-key:b0wewer')

```
api_key: <api key>
```

## Run benchmark [GET] 
run benchmark and notify result on slack channel

```
<ip-address>:<port no>/run?commit=<Commit hash>
```

## Run benchmark scheduler [GET] 
run benchmark on specified time everyday and notify result on slack channel
```
<ip-address>:<port no>/run_scheduler?time=<Server time>
```

## View server time [GET]
returns server time
```
<ip-address>:<port no>/servertime
```

## View all results [GET]
returns JSON of all benchmark results in the database
```
<ip-address>:<port no>/allresults 
```

## Filter results [GET]
Filters and returns result based on argument given
- n = all natural numbers 
```
<ip-address>:<port no>/filer_result?date=<reverse order for mysql>&commit=<commit hash>&commit=<commit hash>&..n&test_no=<int>
```

## Webhook [POST]
Triggers benchmark run on current HEAD (Called from github)  
```
<ip-address>:<port no>/webhook
```