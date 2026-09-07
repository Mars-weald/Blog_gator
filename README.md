Blog_aggregator does what it says on the tin: aggregates blogs in the CLI. 

_Go_ and _Postgres_ are required to run the 'gator.

Install using ```go install```, then run using ```gator``` on the CLI.
```register``` a user, then get aggregating!
Feeds can be added to the database with ```addfeed {name} {url}``` and can be followed with ```follow {url}```

The ```agg``` func is a continuous loop that scrapes the feed database at the time interval you pass to it (```1m```, ```10s```, etc)
