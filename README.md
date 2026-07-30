# Gen data version v1

## Help
```
Gen data to databases with configurable concurrency and batch size.

Usage:
  gendata [command] [flags]

Available Commands:
  mysql        Generate data to MySQL database (alias: sql)
  postgres     Generate data to PostgreSQL database (alias: pg, postgresql)
  clickhouse   Generate data to ClickHouse database (alias: ck)
  completion   Generate the autocompletion script for the specified shell
  help         Help about any command

Global Flags:
  --concurrency       int       Number of concurrent write workers (default: 1)
  --batchsize         int       Number of records to write in each batch (default: 1000)
  --repeatcount       int       Number of times to repeat the batch write (default: 10)
                                Only used in write mode when --duration = 0
  --mode              string    Run mode: write | read | mixed (default: write)
  --readconcurrency   int       Number of concurrent read workers (default: same as --concurrency)
  --duration          duration  Run for a fixed duration, e.g. 60s, 2m (default: 0)
                                Required for read/mixed mode; overrides --repeatcount in write mode
  --debug             bool      Enable debug log output

Database Connection Flags (required for mysql, postgres, clickhouse):
  --host             string          Database host
  --port             int             Database port
  --user             string          Database username
  --password         string          Database password
  --dbname           string          Database name
  --table            string          Database table (default: gendata_table)
  									 The actual table and its structure are fixed;
									 This option changes the default table name used for insertions
  --additionalargs   stringToString  Additional connection arguments (key=value format)
									 Test: --additionalargs "sslmode=disable,charset=utf8mb4"
Note:
The actual table and its structure are fixed; the --table option merely changes the default table name used for insertions.
```

## Truncate Table
```SQL
TRUNCATE TABLE gendata_table;
TRUNCATE TABLE gendata_table RESTART IDENTITY;
```

## Test
```bash
gendata postgres --host=127.0.0.1 --port=5432 \
  --user=postgres --password=xxxxx \
  --concurrency=1 --batchsize=3000 --repeatcount=3

Completed: Con = 0, Count = 0, InsertTime = 0.386s, PerS = 7765 row/s
Completed: Con = 0, Count = 1, InsertTime = 0.167s, PerS = 18017 row/s
Completed: Con = 0, Count = 2, InsertTime = 0.141s, PerS = 21287 row/s
Total: Time = 0.694s, Mean = 12972 row/s, Max = 21287 row/s, Min = 7765 row/s
```

## Modes

`--mode` controls the workload:

| mode   | behavior                                   | needs data first?   | requires `--duration` |
|--------|--------------------------------------------|---------------------|-----------------------|
| write  | pure insert throughput (default)           | no                  | no                    |
| read   | point select by `user_id` (unique index)   | yes (run write first)| yes                  |
| mixed  | concurrent read + write, seeds initial rows| no                  | yes                   |

Read/write ratio in `mixed` is set by `--concurrency` (writers) vs `--readconcurrency` (readers).

Latency percentiles (P50/P95/P99) are reported for both read and write.

### Read (point select)
```bash
# 1) first write some data
gendata postgres --host=127.0.0.1 --port=5432 --user=postgres --password=xxxxx \
  --concurrency=4 --batchsize=3000 --repeatcount=10

# 2) then run point-select load for 60s
gendata postgres --host=127.0.0.1 --port=5432 --user=postgres --password=xxxxx \
  --mode=read --readconcurrency=32 --duration=60s
```

### Mixed (read + write)
```bash
# 8 writers + 32 readers for 60s (≈20% write / 80% read)
gendata postgres --host=127.0.0.1 --port=5432 --user=postgres --password=xxxxx \
  --mode=mixed --concurrency=8 --readconcurrency=32 --batchsize=1000 --duration=60s
```

### Sample mixed output
```
Completed: Con = 0, Count = 0, InsertTime = 0.386s, PerS = 7765 row/s
...

=== Write ===
Total: Time = 8.120s, Rows = 80000, Mean = 9852 row/s, Max = 12000 row/s, Min = 6500 row/s
BatchLatency: P50 = 91.200ms, P95 = 142.500ms, P99 = 168.000ms (samples = 80)

=== Read ===
Total: Time = 57.300s, Queries = 1843200, QPS = 32167
QueryLatency: P50 = 0.420ms, P95 = 2.100ms, P99 = 5.800ms (samples = 500000)
```
