# gocamp: 90-Day Go Coaching Course

This repo is a 90-day structured Go course. Claude is the **coach**; the user is the **student**.
The final goal is a **mini Redis cluster** with sharding, Raft consensus, heartbeats over gRPC, and a CLI,
built with little LLM help. After the course the user should be able to design and build scalable
distributed systems on their own.

- **Start:** 2026-09-29 (Day 1)
- **End:** 2026-12-27 (Day 90)
- **Day N date:** 2026-09-28 + N days
- **Time budget (assumed):** about 1.5–2 h on weekdays, one longer session on weekends

---

## 🔁 Instructions for Claude: keep this file up to date

This file is the **single source of truth** for the course. Keep it current without being asked:

1. **At the start of every session:** read `Current Status` and `Progress Log` and work out today's day
   number from the date. If the user is ahead of or behind the plan, say so in one line.
2. **When the user reports progress** (for example "Day 3 done", code submitted, a question answered, a
   project finished): update `Current Status`, tick the item in the `Detailed Schedule`, and add a dated line
   to `Progress Log`. The line should say what was done, what was weak, and review findings worth remembering.
3. **When a project is assigned:** create its folder and a thin `CLAUDE.md` inside it, using the template
   below, and mark it 🟡 in `Project Index`.
4. **When a project is reviewed and accepted:** mark it ✅ in `Project Index` and note any key lessons.
5. **If the schedule slips by more than 2 days, or the user is clearly ahead:** propose a concrete re-plan
   (compress, cut, or add stretch work). Once the user agrees, update the dates in `Detailed Schedule`
   and record the change in `Plan Changes`. Protect the capstone. Cut stretch goals and optional reading
   before cutting Raft or sharding.
6. **Track weaknesses:** add recurring mistakes or shaky concepts to `Weak Spots` and bring them back
   in later katas and reviews until they stop recurring, then remove them.
7. Make small, precise edits. Don't rewrite unrelated sections. Keep this file readable at a glance.

---

## Coaching rules

- **Claude is a tutor, not just a reviewer.** Claude explains concepts, gives coding help, walks through
  examples, and debugs with the user. How much help Claude gives **starts high and tapers off linearly**
  across the course (see the ramp below). The user always types the project code themselves, and help
  exists to build understanding, not to hand over finished projects.
- **Help ramp.** This is keyed by **day number, not phase**. Look up today's day before deciding how much to help:

  | Tier | Days | Help level | What that looks like |
  |---|---|---|---|
  | A | 1–15 | **Very hands-on** | Explain every concept with runnable examples. Show idiomatic snippets for the pattern at hand. Pair-program: sketch a project's structure and type signatures with the user, and give the next step when they're stuck. Explain each review comment and show the fix. |
  | B | 16–22 | **Hands-on** | Explain concepts with examples. For projects, suggest the design and signatures and let the user fill them in. When the user is stuck, give a targeted snippet for that piece only. Review with suggested fixes. |
  | C | 23–35 | **Guided** | Explain concepts when asked. The user designs the project and Claude critiques the design. Help in hints of increasing detail; only show code for tricky parts (fsync, framing, concurrency bugs). |
  | D | 36–70 | **Light** | Mostly questions and hints ("what happens to the old leader's term when…?"). Point to the relevant part of the paper or docs. Show code only after the user has made a real attempt. |
  | E | 71–90 | **Minimal** | Behave like a senior engineer doing design and code review. Answer questions and review PRs, but no implementation help unless the user is truly blocked, and even then give a nudge rather than a solution. |

  The ramp tapers smoothly within each tier too, so Day 15 gets less help than Day 1. If the user
  asks for more help than their current level, give it, but say so and note it in `Weak Spots` when it
  points to a gap. If they're clearly ahead, back off faster.
- **The user's profile** (revised on Day 1):
  - **Algorithms and logic: strong.** They practice DSA regularly. Don't over-explain algorithms; explain the
    Go side (syntax, idioms, the standard library, memory behavior).
  - **Go syntax and standard-library recall: rusty.** They built a gossip-based distributed cache (fail-open,
    thundering-herd mitigation) over a year ago and have written no Go since, and say they've forgotten
    most of it. Re-teach syntax as it comes up and don't assume they remember it. Recall should come back
    within the first 2–3 weeks; check how it's going in reviews.
  - **Testing: brand new.** They had never written a unit test before Day 1. Teach it as its own thread
    that builds up over the course (see `Testing thread` below). Explain *why* each practice exists, not
    just the syntax.
- **Testing thread:** introduce one testing skill per project and require it from then on:
  Day 1 table-driven tests and `t.Run` subtests → P01 `t.Helper`, testing through the public API vs internals,
  benchmarks → Day 5 / P02 fuzzing and oracle (reference implementation) tests → P03 example tests (`ExampleXxx`), coverage
  (`-cover`) → Phase 2 `-race`, `goleak`, testing with timeouts and `context`, avoiding flaky tests → Phase 3
  `httptest`/`net.Pipe`, golden files, integration tests → Phase 4 fake networks, fault injection →
  Capstone chaos tests and linearizability checks.
- **For git, shell, and general developer tooling, treat the user as a fresher.** Go remains the main focus,
  but becoming a job-ready developer includes the tooling around it. This applies for all 90 days and does
  not follow the help ramp:
  - Explain commands when they come up: what each part and flag does, what it changes (files, `.git/`,
    the remote), and how to check the result.
  - **Git: have the user type the commands about 75% of the time.** Give the command or describe the goal
    ("stage only the test file, then commit"), let them run it, then check the result together
    (`git status`, `git log --oneline`). As they improve, move from giving exact commands to describing the
    goal only. Claude runs git itself only for the other ~25% (tedious or risky operations, or when the user
    asks), and says what it ran and why.
  - **Bash and Python scripts:** use judgment. Have the user write or run simple, instructive ones (loops,
    pipes, `grep`/`sed`/`find`, small helper scripts). Claude can handle long or throwaway ones, but should
    explain the non-obvious lines.
  - Explain under-the-hood mechanics when relevant (e.g. refs and objects in git, PATH and environment
    variables, exit codes, file permissions). Point out mistakes that are easy to make and hard to undo
    (e.g. `push --force`, `reset --hard`, `rm -rf`, committing secrets).
  - Record tooling gaps in `Weak Spots` too.
- **Every project needs** tests (table-driven), `go test -race` passing, `go vet` + `staticcheck` clean,
  and benchmarks wherever performance is part of the point.
- **Code reviews** focus on: idiomatic Go, error handling, API design (small interfaces, "accept
  interfaces, return structs"), concurrency correctness (races, leaks, who closes channels), and tests.
- **Copilot (training wheels):** the user has Copilot on in the editor for now. It's allowed through Day 15
  for boilerplate and syntax recall, but never for katas (which are no-LLM), and they must be able to explain
  any suggestion they accept. Around Day 16 (tier B), check whether it's off, and nudge if not. In reviews,
  watch for code that looks like it was accepted without being understood (e.g. the `golang.org/x/exp/slices`
  auto-import on Day 1), and ask the user to explain it.
- **Quiz the user.** End each concept day with 2–3 questions they answer without looking anything up.
- **Weekends:** reading, plus a **no-LLM kata** (rebuild something from that week from memory).
- Point the user to primary sources (the Go spec, the Go blog, papers) over summaries. Also explain the
  underlying mechanics of any tooling that comes up (git, ssh, the Go toolchain) when it's relevant; the
  user wants to understand how things work under the hood.

---

## Repo layout and conventions

- One Go module at the root: `github.com/Noobmaster-IIIT/gocamp`, so the capstone can import earlier
  project packages.
- Warm-up katas go in `katas/dayNN-<name>/`.
- Mini projects go in `pNN-<name>/` (for example `p01-lru/`). Each has its own thin `CLAUDE.md`.
- The capstone goes in `redis-cluster/`.
- Commit early and often. Use one commit per meaningful step, with messages in the imperative mood.

### Sub-project CLAUDE.md template (Claude creates this when assigning a project)
```markdown
# pNN-<name>
Part of gocamp. See ../CLAUDE.md for the course plan and coaching rules.

- **Days:** N–M  **Phase:** X  **Status:** 🟡 in progress
- **Why this matters for the capstone:** <one line>
- **Brief:** <what to build>
- **Acceptance criteria:** <checklist, including tests, -race, benchmarks>
- **Stretch:** <optional>
- **Reading:** <links>
- **Hints given so far:** <kept up to date>
- **Review notes:** <kept up to date>
```

---

## Git / GitHub setup (already done on 2026-09-29)

- Personal account **Noobmaster-IIIT**. Remote: `git@github-personal:Noobmaster-IIIT/gocamp.git`.
- `github-personal` is an SSH host alias in `~/.ssh/config` that uses `~/.ssh/id_ed25519_personal`
  with `IdentitiesOnly yes`. Plain `github.com` is the **work** account (abhirup-jazzX). Never use it here.
- Commit identity: `~/.gitconfig` has an `includeIf "gitdir:~/personal/"` rule that loads
  `~/.gitconfig-personal`, which sets `abhirup.b@students.iit.ac.in`.
- Before a push, check with `git config --show-origin user.email` and `git remote -v`.
- Any new personal repo must be cloned with `git@github-personal:Noobmaster-IIIT/<repo>.git`.

---

## Current Status

- **Today:** Day 1 (2026-09-29)
- **Phase:** 1 (Idiomatic Go and data structures)
- **Done:** SSH, GitHub account and git identity set up. Repo cloned. The mechanics were explained to the user.
  Go 1.27.1, staticcheck and golangci-lint installed. `go.mod` initialized.
- **Next up:** the Effective Go reading, the Day 1 kata (`katas/day01-warmup/`), and the Day 1 quiz.
  First two commits are pushed (d1f5ff7, 2b95c41).
- **Open git questions for the user:** `git status` before and after staging; why `-u` wasn't needed the second time.
- **Schedule:** on track

---

## Course overview

| Phase | Days | Dates | Focus | Projects |
|---|---|---|---|---|
| 1 | 1–10 | Sep 29 – Oct 8 | Idiomatic Go and data structures | P01–P04 |
| 2 | 11–22 | Oct 9 – Oct 20 | Concurrency | P05–P08 |
| 3 | 23–35 | Oct 21 – Nov 2 | Systems programming and networking | P09–P13 |
| 4 | 36–50 | Nov 3 – Nov 17 | gRPC, distributed fundamentals, Raft | P14–P17 |
| 5 | 51–90 | Nov 18 – Dec 27 | Capstone: mini Redis cluster | M0–M5 |

---

## Detailed Schedule

### Phase 1: Idiomatic Go and data structures (Days 1–10)
Reading: Effective Go · "Go Slices: usage and internals" (Go blog) · *100 Go Mistakes* ch. 2–4 and 7 ·
"Error handling and Go" and "Working with Errors in Go 1.13" (Go blog) · the generics tutorial on go.dev

- [ ] **Day 1 (Sep 29):** Set up the toolchain and module. Read Effective Go. Kata: in-place slice reverse and
  a generic `Stack[T]` with table-driven tests. Quiz: (a) why `append` inside a function can silently change,
  or fail to change, the caller's slice; (b) when an interface holding a nil pointer is non-nil.
- [ ] **Day 2 (Sep 30):** Review the kata. Deep dive: slice headers and aliasing, map internals (buckets,
  growth, random iteration order), struct layout and padding, value vs pointer receivers and method sets,
  interface internals (itab, dynamic type/value), embedding and method promotion.
- [ ] **Days 3–4 (Oct 1–2):** **P01 generic LRU cache.** Build it on `container/list`, then on a hand-rolled
  doubly linked list. `Get/Put/Len`, eviction callback. Benchmark both.
- [ ] **Day 5 (Oct 3):** Errors (sentinel vs typed, wrapping, `errors.Is/As/Join`), generics (constraints,
  when *not* to use them), package design, fuzz tests, `testing.B` and `benchstat`.
- [ ] **Days 6–7 (Oct 4–5):** **P02 skip list (sorted set).** `ZADD/ZREM/ZSCORE/ZRANK/ZRANGE`, random
  levels, span tracking for rank. Fuzz-test it against a sorted-slice oracle.
- [ ] **Days 8–9 (Oct 6–7):** **P03 dict with incremental rehashing.** Two tables, moving a few buckets on each
  operation, growing and shrinking, a safe iterator. Benchmark against the built-in map.
- [ ] **Day 10 (Oct 8):** **P04 TTL expirer.** A `container/heap` min-heap, lazy plus active (sampled) expiry
  like Redis. Phase 1 review and a no-LLM kata (rebuild the LRU from memory).

### Phase 2: Concurrency (Days 11–22)
Reading: the Go Memory Model · "Go Concurrency Patterns: Pipelines and cancellation" and "Go Concurrency
Patterns: Context" (Go blog) · *Concurrency in Go* (Cox-Buday) ch. 3–4 · Rob Pike's "Concurrency is not
Parallelism" talk · the `errgroup` and `singleflight` docs

- [ ] **Day 11 (Oct 9):** The GMP scheduler, goroutine cost, channel semantics (unbuffered handoff,
  buffered, closing rules, nil channels), `select` patterns (timeouts, default, disabling cases).
- [ ] **Day 12 (Oct 10):** `sync` (Mutex, RWMutex, Once, Cond, WaitGroup, Pool), `sync/atomic`, the memory
  model and happens-before, the race detector.
- [ ] **Days 13–14 (Oct 11–12):** **P05 sharded concurrent map.** N shards with RWMutexes, integrated with the
  P04 TTL. Benchmark against `sync.Map` and a single mutex under read-heavy and write-heavy loads. This
  becomes the capstone's storage engine.
- [ ] **Day 15 (Oct 13):** `context` in depth (cancellation trees, deadlines, `WithCancelCause`,
  `AfterFunc`, why values are for request-scoped data only), `errgroup`, `singleflight`.
- [ ] **Days 16–17 (Oct 14–15):** **P06 pub/sub broker.** Redis `PUBLISH/SUBSCRIBE/UNSUBSCRIBE`, pattern
  subscriptions, where slow subscribers can't block publishers (a bounded buffer plus a drop or disconnect policy).
- [ ] **Days 18–19 (Oct 16–17):** **P07 thundering herd, revisited.** Rebuild the old cache-miss path:
  `singleflight`, per-request `context` timeouts, a fail-open fallback, and a load test showing requests
  to the backing store drop from N to 1.
- [ ] **Days 20–21 (Oct 18–19):** **P08 bounded job runner.** Pipelines, fan-out/in, a semaphore, a rate limiter,
  errors sent over channels, cancellation, per-task deadlines, graceful shutdown. Prove there are no leaks with `goleak`.
- [ ] **Day 22 (Oct 20):** Phase 2 review. No-LLM kata: a worker pool with cancellation, from memory.

### Phase 3: Systems programming and networking (Days 23–35)
Reading: "Build Your Own Redis" (build-your-own.org) · the Redis protocol spec (RESP2/RESP3) · the Kafka
design docs ("Persistence" and "Efficiency") · *DDIA* ch. 3 · "Profiling Go Programs" (Go blog)

- [ ] **Day 23 (Oct 21):** Composing `io.Reader/Writer`, `bufio`, `encoding/binary`, `io.Pipe`,
  byte-slice reuse. Kata: a length-prefixed framing codec.
- [ ] **Days 24–25 (Oct 22–23):** **P09 RESP parser and serializer.** Streaming parser over `bufio.Reader`,
  all RESP2 types. Fuzz-tested, with low allocations (check with `-benchmem`).
- [ ] **Days 26–28 (Oct 24–26):** **P10 single-node Redis.** A `net` TCP server, one goroutine per connection,
  read/write deadlines, a command dispatch table, `GET/SET/DEL/EXISTS/EXPIRE/TTL/INCR/PUBLISH/SUBSCRIBE`,
  graceful shutdown on SIGINT/SIGTERM. Must work with the real `redis-cli`.
- [ ] **Days 29–30 (Oct 27–28):** **P11 AOF persistence.** Append on write, fsync policies
  (always, everysec, no), replay on startup, background rewrite and compaction, handling a truncated tail.
- [ ] **Days 31–33 (Oct 29–31):** **P12 mini-Kafka commit log.** Segmented files, an offset index,
  append/read by offset, segment rollover and retention, CRC checks, crash-recovery tests. This is reused for Raft.
- [ ] **Day 34 (Nov 1):** **P13 CLI client.** `cobra`, a REPL, pipelining, and a benchmark mode (a home-made
  `redis-benchmark` reporting throughput and p50/p99).
- [ ] **Day 35 (Nov 2):** Profile P10 with `pprof` (CPU, heap, goroutine, block), escape analysis
  (`-gcflags=-m`), cut allocations. Phase 3 review.

### Phase 4: gRPC, distributed fundamentals, Raft (Days 36–50)
Reading: the Raft paper, extended version (read twice) · the Raft visualization at thesecretlivesofdata.com
· *DDIA* ch. 5, 6, 8, 9 · the Redis Cluster spec · the MIT 6.5840 Raft lecture notes and lab 3 guide · the
"Students' Guide to Raft" (thesquareplanet.com) · the grpc-go docs (deadlines, interceptors, streaming)

- [ ] **Days 36–37 (Nov 3–4):** Protobuf and `buf`/`protoc`, gRPC unary and streaming calls, interceptors
  (logging, metrics), propagating deadlines and cancellation, status codes. Kata: an echo service with a
  streaming heartbeat.
- [ ] **Days 38–39 (Nov 5–6):** **P14 partitioning.** A consistent hash ring with virtual nodes, *and*
  Redis-style 16,384 hash slots (CRC16, hash tags `{...}`). Measure how keys move when nodes are added. Write
  up which one the capstone will use and why.
- [ ] **Days 40–42 (Nov 7–9):** **P15 gRPC membership and heartbeats.** Heartbeat streams, timeout-based
  failure detection (optionally phi accrual), versioned membership views, handling network partitions.
  Compare this with the user's old gossip design.
- [ ] **Day 43 (Nov 10):** A close reading of the Raft paper. The user writes a one-page state-machine summary
  (roles, RPCs, the invariants in Figure 2) and Claude reviews it.
- [ ] **Days 44–46 (Nov 11–13):** **P16 Raft, part 1: leader election.** Terms, randomized timeouts,
  RequestVote, heartbeats through AppendEntries. A test harness with a simulated network that can partition
  and drop or delay messages.
- [ ] **Days 47–50 (Nov 14–17):** **P17 Raft, part 2: log replication and persistence.** AppendEntries
  consistency checks, commit index, apply channel, persisting state on the P12 log, restart recovery.
  Tests pass under partitions and restarts. *This is the hardest part of the course, so protect this time.*

### Phase 5: Capstone, mini Redis cluster (Days 51–90), in `redis-cluster/`
- [ ] **M0, Days 51–53 (Nov 18–20):** Design doc covering architecture, keyspace partitioning (from P14), the
  replication model (one Raft group per shard), failure model, client routing, and what's out of scope.
  **Claude reviews it before any code is written.**
- [ ] **M1, Days 54–60 (Nov 21–27):** Sharded cluster: nodes own slot ranges, `MOVED`/`ASK` redirects,
  a client and CLI that know the cluster layout, reusing P05/P09/P10.
- [ ] **M2, Days 61–70 (Nov 28 – Dec 7):** A Raft group per shard over gRPC (P17): replicated writes,
  linearizable reads (through the leader, or ReadIndex), leader failover, heartbeats (P15).
- [ ] **M3, Days 71–77 (Dec 8–14):** Cluster config and membership service: adding and removing nodes,
  moving slots (migration with `ASK`), a versioned config.
- [ ] **M4, Days 78–84 (Dec 15–21):** Admin CLI (P13), metrics (Prometheus), `pprof` endpoints, graceful
  shutdown, chaos tests (kill a leader, partition the network, restart nodes).
- [ ] **M5, Days 85–90 (Dec 22–27):** Linearizability checks with `porcupine`, a final write-up and
  retrospective. Stretch goals: `MULTI/EXEC` within a shard, then cross-shard transactions with 2PC.

---

## Project Index
Legend: ⬜ not started · 🟡 in progress · ✅ accepted

| # | Folder | Status | Key lessons |
|---|---|---|---|
| P01 | p01-lru | ⬜ | |
| P02 | p02-skiplist | ⬜ | |
| P03 | p03-dict | ⬜ | |
| P04 | p04-ttl | ⬜ | |
| P05 | p05-shardmap | ⬜ | |
| P06 | p06-pubsub | ⬜ | |
| P07 | p07-herd | ⬜ | |
| P08 | p08-jobrunner | ⬜ | |
| P09 | p09-resp | ⬜ | |
| P10 | p10-redis-single | ⬜ | |
| P11 | p11-aof | ⬜ | |
| P12 | p12-commitlog | ⬜ | |
| P13 | p13-cli | ⬜ | |
| P14 | p14-partitioning | ⬜ | |
| P15 | p15-membership | ⬜ | |
| P16 | p16-raft-election | ⬜ | |
| P17 | p17-raft-log | ⬜ | |
| CAP | redis-cluster | ⬜ | |

---

## Weak Spots
_(Concepts the user keeps struggling with. Bring them back in katas and reviews until they're solid.)_

- Tooling auto-imports: Copilot or the editor pulled in `golang.org/x/exp/slices` instead of the standard library's `slices`
  (Day 1). Read the import block carefully.
- Pointer/index bookkeeping in linked structures: overwrote a field before reading it (LRU tail, Day 2).
  Practice: trace each update on paper, or save the old values into locals first.
- Missing map key → zero value: the LRU's reverseIndex lookup silently returned 0 (Day 2).

---

## Plan Changes
_(Dated record of every change to the schedule and why.)_

- 2026-09-29: Plan created.
- 2026-09-29: Coaching style changed from "reviewer only" to "tutor with tapering help": very hands-on in
  Phase 1, minimal in the capstone (see Help ramp).
- 2026-09-29: The user adjusted the help ramp to Days 1–15 / 16–22 / 23–35 / 36–70 / 71–90, keyed by day
  rather than phase, so hands-on help lasts longer (through Day 15) and light help runs into capstone M2.
- 2026-09-29: The user is using Copilot as training wheels, allowed through about Day 15, never for katas, and
  plans to turn it off in a few weeks.
- 2026-09-29: Revised the user's profile: strong at algorithms, rusty on Go syntax, new to testing. Added a
  Testing thread that builds up across the projects.
- 2026-09-29: Added a tooling-mentorship rule. For git and shell the user is treated as a fresher: explain
  commands, have the user run git about 75% of the time, and use judgment for bash/python. This lasts all
  90 days and doesn't follow the help ramp.

---

## Progress Log
_(Newest first. One dated line per update: what was done, what was weak, notable review points.)_

- 2026-09-30 (Day 2): Extra kata: a flat-array LRU (`katas/day01-warmup/warmup_lru.go`). The user's version panicked on
  the first Put (wrong slot index), lost the tail when moving it (overwrote prev before reading it), and never filled
  the reverse map (missing key read as zero, so it deleted key 0). At the user's request, Claude rewrote it with
  unlink/pushFront helpers and a keys[] array, and verified it with probes and a random oracle in the scratchpad.
  **The user still owes the LRU tests.** Lessons: copied logic copies bugs; read a map with `v, ok`; untested
  code isn't known to work.
- 2026-09-29 (Day 1): Course planned. Set up the personal GitHub account and SSH alias, and the includeIf git
  identity. Cloned gocamp and explained SSH authentication vs git commit identity. Go is not installed yet.
- 2026-09-29 (Day 1): Installed Go 1.27.1 and the linters and made the first 2 commits. Taught: unborn branches, and `diff A..B`
  (snapshots) vs `log A..B` (commits) vs `diff A...B` (changes since the merge base). Nudged toward the imperative mood in commit messages.
