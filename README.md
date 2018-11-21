# charitycase

This repo contains code to help calculate the amount of government grants a list
of charities has received.

It takes two types of file:

1. `-c`: A CSV containing the list of charities to find grant totals for
1. `-d`: Path to SQLite database containing entries from open990.com containing data from Form 990 Part VIII
   line 1e


## Making the `data.csv`

1. Download "GovernmentGrants.zip", unzip
1. Download "GovernmentGrantsAmt.zip", unzip
1. `tail -n+2 GovernmentGrants/GovernmentGrants.csv > data.csv`, to remove the
   copyright notice
1. `tail -n+3 GovernmentGrantsAmt/GovernmentGrantsAmt.csv >> data.csv`, to
   remove the copyright notice and schema header, and append it to data.csv

## Making the `grants.db`

```sh
ch $ sqlite3 grants.db
SQLite version 3.19.3 2017-06-27 16:48:08
Enter ".help" for usage hints.
sqlite> .mode csv
sqlite> .import "data.csv" grants
sqlite> .quit
```
