# lsm-tree

A log-structured merge-tree storage engine written from scratch in Go, as a hands-on exercise for Chapter 3 ("Storage and Retrieval") of *Designing Data-Intensive Applications*.

The goal is to build up the standard LSM pipeline piece by piece: write-ahead log → in-memory memtable → sorted string tables on disk → compaction.

## Current state

Working today:

- **Write-ahead log** — every `Put`/`Delete` is appended to `data/wal.log` before it touches memory, so writes survive a crash.
- **Crash recovery** — `db.Open` replays the entire WAL on startup and rebuilds the memtable from it.
- **In-memory memtable** — a plain `map[string][]byte` backing `Put`, `Get`, `Delete`, and `Count`.
- **Flush trigger** — once the memtable holds more than 5 entries, `db.Put` calls `Flush`.

Not done yet (see [Roadmap](#roadmap)): `Flush` currently sorts the keys and prints them rather than writing an SSTable, so nothing is ever persisted beyond the WAL, and the memtable is never cleared or truncated.

## Layout

```
cmd/lsm/main.go      entry point / scratch driver
db/db.go             DB type: Put, Get, Delete
db/open.go           Open: WAL recovery into a fresh memtable
db/sstable.go        Flush (stub — sorts and prints keys)
memtable/memtable.go map-backed memtable
wal/wal.go           append-only log: Put, Delete, Replay
record/record.go     Record type + WAL line parser
data/wal.log         the log itself (created on first run)
```

## Running it

```sh
go run ./cmd/lsm
```

`main.go` is a scratch driver: it opens the DB at `./data`, writes six keys, and exits. The sixth write crosses the flush threshold, so you should see `Memtable Full` followed by the sorted key list. Run it twice and the replay path kicks in — the WAL is re-read on open, so the memtable starts full and the flush fires on the first write instead.

To start clean, delete `data/wal.log`.

## WAL format

One record per line, pipe-delimited:

```
PUT | key | value
DELETE | key
```

`record.ParseWal` splits on `|` and trims each field.

## Roadmap

Roughly in the order the book introduces them:

1. **Make `Flush` write a real SSTable.** Serialize the sorted entries to a segment file on disk, then clear the memtable and truncate the WAL.
2. **Read path through SSTables.** `db.Get` only checks the memtable right now; it needs to fall back to segment files, newest first.
3. **Tombstones.** `memtable.Delete` deletes the map key outright. Once deleted keys can still exist in older segments, a delete has to write a tombstone marker instead, and reads have to stop at it.
4. **Sparse index per segment,** so a lookup doesn't scan the whole file.
5. **Compaction and merging** of segment files, dropping overwritten values and tombstones.
6. **Bloom filters** to avoid touching segments that can't hold the key.

## Known rough edges

- `db.Put` and `db.Delete` `panic` on WAL write errors instead of returning them.
- Neither the WAL file nor the DB is ever closed — there's no `DB.Close`, and writes rely on `os.File` buffering plus process exit rather than an explicit `Sync`.
- A value containing `|` or a leading/trailing space won't round-trip through the WAL parser.
- `Flush`'s error return is discarded by the caller.
- Nothing is safe for concurrent use; there are no locks anywhere.
- The flush threshold (5 entries) is hardcoded in `db.Put`.
- No tests yet.
