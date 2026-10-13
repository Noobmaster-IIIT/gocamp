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
   number from the date. Don't comment on the schedule by default (see the pace rule below).
2. **When the user reports progress** (for example "Day 3 done", code submitted, a question answered, a
   project finished): update `Current Status`, tick the item in the `Detailed Schedule`, and add a dated line
   to `Progress Log`. The line should say what was done, what was weak, and review findings worth remembering.
3. **When a project is assigned:** create its folder and a thin `CLAUDE.md` inside it, using the template
   below, and mark it 🟡 in `Project Index`.
4. **When a project is reviewed and accepted:** mark it ✅ in `Project Index` and note any key lessons.
5. **Pace rule: fun first, timeline second** (the user asked for this on Day 15). This is a journey, not a deadline.
   Don't push the schedule, nag, or frame things around dates. Only when there's a **sustained** slowdown (about a
   week or more with little progress, or the projected finish slipping by more than about 7 days) mention it
   once, neutrally, with the number: *"At the current pace the capstone would finish around <date> (~N days
   later than planned)."* Then let the user decide. Offer a re-plan only if they want one. Record agreed changes
   in `Plan Changes` and update the dates in `Detailed Schedule`. When cutting, protect the capstone: drop stretch
   goals and optional reading before Raft or sharding. Watch for burnout signals (frustration, "chore" feelings)
   and lighten up: shorter tasks, more fun side quests, more factoids.
6. **Track weaknesses:** add recurring mistakes or shaky concepts to `Weak Spots` and bring them back
   in later katas and reviews until they stop recurring, then remove them.
7. **Keep `go_notes.md` up to date:** whenever Claude drops a factoid, explains a Go concept, or shows a useful
   signature or idiom, add it to the right section of `go_notes.md` (or add a new section and link it from
   Contents). Also update it whenever the user asks. Keep it in documentation style: headings, short bullets,
   runnable snippets, and no chat tone.
8. Make small, precise edits. Don't rewrite unrelated sections. Keep this file readable at a glance.

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
- **Quizzes: only on what has been taught.** Quiz the user (2–3 questions, from memory) only on concepts
  already covered in a lesson. **Present deep internals as factoids instead:** short explained
  "did you know" notes with a runnable example and no question attached. Examples are slice aliasing and
  cap, interfaces holding nil pointers, map growth, and escape analysis. They can turn into quiz questions
  later, once the concept has been taught and used in a project. Keep quizzes encouraging and aimed at the
  user's current level. (The user asked for this on Day 2: internals questions asked before teaching
  were discouraging.)
- **Weekends:** reading, plus a **no-LLM kata** (rebuild something from that week from memory).
- **The user sometimes picks a longer approach on purpose, to practice** (e.g. a recursive heap-invariant checker on
  Day 15 to write their first recursive Go function). Present simpler alternatives as optional ("FYI, there's also…"),
  not as fixes, unless the code is actually wrong.
- **Keep dropping factoids.** Short "did you know" notes about Go's design, history, internals, and
  idioms, even when the user has read the topic before. Repetition helps them and they enjoy it. Don't
  apologize for repeating things.
- Point the user to primary sources (the Go spec, the Go blog, papers) over summaries. Also explain the
  underlying mechanics of any tooling that comes up (git, ssh, the Go toolchain) when it's relevant; the
  user wants to understand how things work under the hood.

---

## Repo layout and conventions

- One Go module at the root: `github.com/Noobmaster-IIIT/gocamp`, so the capstone can import earlier
  project packages.
- `go_notes.md` holds the user's Go reference notes (stdlib cheat sheet and factoids). Claude maintains it.
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

- **Today:** Day 14 (2026-10-12). Last session was Day 4 (Oct 2), so there was a 10-day gap.
- **Phase:** 1 (Idiomatic Go and data structures)
- **Done:** SSH, GitHub account and git identity set up. Repo cloned. The mechanics were explained to the user.
  Go 1.27.1, staticcheck and golangci-lint installed. `go.mod` initialized.
- **Next up:** finish LC 703 with `Heap[T]` (the user is doing this on their own), plus the open heap review items
  (move the code to `heap.go`, zero the slot in Pop, pointer receiver on Len, add Peek). Then **Day 16: goroutines,
  channels, select**.
- **Already covered:** comma-ok (maps, type assertions), composite literals and `struct{}{}`,
  table-driven tests and `t.Run`, pointer vs value receivers, int sizes and explicit integer conversion.
  (Repeating these is welcome. The user likes seeing them again, so keep bringing them up.)
- **Schedule:** relaxed; see the pace rule. Re-planned on Day 15. Phase 2 starts Day 16 (Oct 14). The capstone shrinks from 40 to 35 days; the end date (Dec 27) is unchanged.

---

## Course overview

| Phase | Days | Dates | Focus | Projects |
|---|---|---|---|---|
| 1 | 1–15 | Sep 29 – Oct 13 | Idiomatic Go and data structures (done; DSA continues as self-study) | LRU, deque, heap katas |
| 2 | 16–27 | Oct 14 – Oct 25 | Concurrency | P05–P08 (+ TCP echo bridge) |
| 3 | 28–40 | Oct 26 – Nov 7 | Systems programming and networking (`net`) | P09–P13 |
| 4 | 41–55 | Nov 8 – Nov 22 | gRPC, distributed fundamentals, Raft | P14–P17 |
| 5 | 56–90 | Nov 23 – Dec 27 | Capstone: mini Redis cluster | M0–M5 |
---

## Detailed Schedule

### Phase 1: Idiomatic Go and data structures (Days 1–15) ✅
Condensed on Day 15 at the user's request. The rest of the DSA work (skip list, incremental-rehash dict) is
**self-study**, not coached. Reading: Effective Go (in progress).

- [x] **Day 1 (Sep 29):** Toolchain, module, git/SSH setup. Kata: `Reverse`, `Stack[T]`, first table-driven tests.
- [x] **Day 2 (Sep 30):** Flat-array LRU (LeetCode 146, accepted). Factoids: slice aliasing, interfaces holding nil.
- [x] **Day 4 (Oct 2):** Two generic deques (ring buffer, unrolled block list) behind a `Deque[T]` interface.
- [x] **Days 14–15 (Oct 12–13):** `container/heap` lesson, then a generic `Heap[T]` with tests written first (random test
  and recursive invariant checker). LC 703 in progress.

### Phase 2: Concurrency (Days 16–27)
Reading: the Go Memory Model · "Go Concurrency Patterns: Pipelines and cancellation" and "Go Concurrency
Patterns: Context" (Go blog) · *Concurrency in Go* (Cox-Buday) ch. 3–4 · Rob Pike's "Concurrency is not
Parallelism" talk · the `errgroup` and `singleflight` docs

- [ ] **Day 16 (Oct 14):** The GMP scheduler, goroutine cost, channel semantics (unbuffered handoff, buffered, closing
  rules, nil channels), `select` patterns (timeouts, default, disabling cases).
- [ ] **Day 17 (Oct 15):** `sync` (Mutex, RWMutex, Once, **Cond**, WaitGroup, Pool), `sync/atomic`, the memory model and
  happens-before, the race detector.
- [ ] **Days 18–19 (Oct 16–17):** **P05 sharded concurrent map with TTL.** N shards with RWMutexes. TTL expiry reuses the
  Day 14 `Heap[T]` (the old P04 is folded in here): lazy expiry plus an active expiry goroutine. Benchmark against
  `sync.Map` and a single mutex. This becomes the capstone's storage engine.
- [ ] **Day 20 (Oct 18):** `context` in depth (cancellation trees, deadlines, `WithCancelCause`, `AfterFunc`, why values
  are for request-scoped data only), `errgroup`, `singleflight`.
- [ ] **Days 21–22 (Oct 19–20):** **P06 pub/sub broker.** Redis `PUBLISH/SUBSCRIBE/UNSUBSCRIBE`, pattern subscriptions,
  where slow subscribers can't block publishers (a bounded buffer plus a drop or disconnect policy).
- [ ] **Days 23–24 (Oct 21–22):** **P07 thundering herd, revisited.** `singleflight`, per-request `context` timeouts,
  a fail-open fallback, and a load test showing requests to the backing store drop from N to 1.
- [ ] **Days 25–26 (Oct 23–24):** **P08 bounded job runner.** Pipelines, fan-out/in, a semaphore, a rate limiter, errors sent
  over channels, cancellation, graceful shutdown. Prove there are no leaks with `goleak`.
- [ ] **Day 27 (Oct 25):** Phase 2 review. **Bridge to `net`:** a TCP echo server with one goroutine per connection, plus
  context-driven graceful shutdown.

### Phase 3: Systems programming and networking (Days 28–40)
Reading: "Build Your Own Redis" (build-your-own.org) · the Redis protocol spec (RESP2/RESP3) · the Kafka
design docs ("Persistence" and "Efficiency") · *DDIA* ch. 3 · "Profiling Go Programs" (Go blog)

- [ ] **Day 28 (Oct 26):** Composing `io.Reader/Writer`, `bufio`, `encoding/binary`, `io.Pipe`, `net.Pipe` for tests.
  Kata: a length-prefixed framing codec.
- [ ] **Days 29–30 (Oct 27–28):** **P09 RESP parser and serializer.** Streaming over `bufio.Reader`, fuzz-tested,
  with low allocations.
- [ ] **Days 31–33 (Oct 29–31):** **P10 single-node Redis.** A `net` TCP server, read/write deadlines, a dispatch table,
  `GET/SET/DEL/EXISTS/EXPIRE/TTL/INCR/PUBLISH/SUBSCRIBE`, graceful shutdown. Must work with the real `redis-cli`.
- [ ] **Days 34–35 (Nov 1–2):** **P11 AOF persistence.** fsync policies, replay, rewrite/compaction, a truncated tail.
- [ ] **Days 36–38 (Nov 3–5):** **P12 mini-Kafka commit log.** Segments, an offset index, rollover, CRCs, crash recovery.
- [ ] **Day 39 (Nov 6):** **P13 CLI client.** `cobra`, a REPL, pipelining, a benchmark mode (p50/p99).
- [ ] **Day 40 (Nov 7):** Profile P10 with `pprof`, escape analysis, cut allocations. Phase 3 review.

### Phase 4: gRPC, distributed fundamentals, Raft (Days 41–55)
Reading: the Raft paper, extended version (read twice) · the Raft visualization at thesecretlivesofdata.com
· *DDIA* ch. 5, 6, 8, 9 · the Redis Cluster spec · the MIT 6.5840 Raft lecture notes and lab 3 guide · the
"Students' Guide to Raft" · the grpc-go docs (deadlines, interceptors, streaming)

- [ ] **Days 41–42 (Nov 8–9):** Protobuf/`buf`, gRPC unary and streaming calls, interceptors, deadlines, status codes.
  Kata: an echo service with a streaming heartbeat.
- [ ] **Days 43–44 (Nov 10–11):** **P14 partitioning.** A consistent hash ring with virtual nodes *and* 16,384 hash slots
  (CRC16, `{...}` hash tags). Measure key movement and choose one for the capstone.
- [ ] **Days 45–47 (Nov 12–14):** **P15 gRPC membership and heartbeats.** Failure detection, versioned views, partitions.
  Compare with the user's old gossip design.
- [ ] **Day 48 (Nov 15):** Close reading of the Raft paper. The user writes a one-page Figure 2 summary for review.
- [ ] **Days 49–51 (Nov 16–18):** **P16 Raft, part 1: leader election**, with a simulated, partitionable network.
- [ ] **Days 52–55 (Nov 19–22):** **P17 Raft, part 2: log replication and persistence** on the P12 log. *This is the
  hardest part of the course, so protect this time.*

### Phase 5: Capstone, mini Redis cluster (Days 56–90), in `redis-cluster/`
- [ ] **M0, Days 56–57 (Nov 23–24):** Design doc (architecture, partitioning, one Raft group per shard, failure model,
  routing, non-goals). **Claude reviews it before any code is written.**
- [ ] **M1, Days 58–63 (Nov 25–30):** Sharded cluster: slot ownership, `MOVED`/`ASK`, a client and CLI that know the
  cluster layout.
- [ ] **M2, Days 64–72 (Dec 1–9):** A Raft group per shard over gRPC: replicated writes, linearizable reads, failover.
- [ ] **M3, Days 73–78 (Dec 10–15):** Cluster config/membership: adding and removing nodes, moving slots, a versioned config.
- [ ] **M4, Days 79–84 (Dec 16–21):** Admin CLI, metrics, `pprof`, graceful shutdown, chaos tests.
- [ ] **M5, Days 85–90 (Dec 22–27):** `porcupine` linearizability checks, final write-up. Stretch: `MULTI/EXEC`, 2PC.

---

## Project Index
Legend: ⬜ not started · 🟡 in progress · ✅ accepted · ➖ dropped/merged

| # | Folder | Status | Key lessons |
|---|---|---|---|
| P01 | p01-lru | ✅ | Done as a kata (`katas/day01-warmup`, flat-array LRU, LC 146 accepted) |
| P02 | p02-skiplist | ➖ | Dropped on Day 15: self-study |
| P03 | p03-dict | ➖ | Dropped on Day 15: self-study |
| P04 | p04-ttl | ➖ | Folded into P05 (TTL via the Day 14 `Heap[T]`) |
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
- Linked/ring structures: think through what happens at the boundaries on the *pop* side too (a block runs empty,
  head == tail when full); the push side was right (Day 4).
- File placement: implementation went into a `_test.go` file twice (Day 1, Day 15). Code goes in `x.go`, tests in `x_test.go`.
- Zeroing vacated slots on pop: missed in the ring deque (Day 4) and the heap (Day 15), despite the Stack lesson.
- Pointer/index bookkeeping in linked structures: overwrote a field before reading it (LRU tail, Day 2).
  Practice: trace each update on paper, or save the old values into locals first.
- Missing map key → zero value: the LRU's reverseIndex lookup silently returned 0 (Day 2).
- Git's three areas (working tree / index / HEAD) and what `git status` shows in each (Day 2 quiz).

---

## Plan Changes
_(Dated record of every change to the schedule and why.)_

- 2026-09-29: Plan created.
- 2026-10-13 (Day 15): Pace rule: fun first and no timeline pressure. Only flag a sustained slowdown, once, as "at the current pace the project is delayed by ~N days".
- 2026-10-13 (Day 15): **Re-plan, at the user's request.** The rest of the DSA projects were dropped (P02 skip list and P03 dict
  become self-study; P04 TTL is folded into P05 using `Heap[T]`). Phase 1 closed. Phases 2–4 shifted +5 days (Phase 2 = Days
  16–27). The capstone went from 40 to 35 days (M0–M4 trimmed by a day each); the end date is still Dec 27. Added a TCP echo
  bridge to `net` at the end of Phase 2 because the user wants to get to `net` sooner. This also absorbs the 10-day gap.
- 2026-09-30: Changed quizzes to cover only taught material. Deep internals are now presented as explained factoids, not questions.
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

- 2026-10-13 (Day 15): **Solved LC 295 (Median from Data Stream)** with their own generic two-heap solution; all tests pass. The logic was correct; the issues were compile-level (duplicate heap in one package, `any` vs `cmp.Ordered`, unpacking `(T, bool)`, int vs float division). Learned: one package per directory.
- 2026-10-13 (Day 15): The heap tests pass after fixing the `for`-init loop (after three unsaved-file confusions; auto-save recommended). The user wants to skip further DSA and move on to concurrency and `net`, so the course was re-planned.
- 2026-10-13 (Day 15): Generic `Heap[T]` with tests written first (they skipped LC 703 for now). up/down are correct, and they
  added an O(n) `makeHeap` on their own. Tests: push/pop, min vs max, a recursive invariant checker, and a seeded random test.
  Issues: the implementation is in `kth_test.go` (second time code went in a _test file); the random test's `for got, ok := h.Pop(); ok; {}`
  only pops once (the user had flagged it); Pop doesn't zero the vacated slot; mixed receivers (Len uses a value receiver);
  Peek is missing; `min`/`max` shadow builtins.
- 2026-10-12 (Day 14): Came back after a 10-day gap. Started the priority-queue lesson (`container/heap`, then a generic `Heap[T]`, with tests first). Proposed a re-plan.
- 2026-10-02 (Day 4): The user built two generic deques themselves (a ring buffer, and an unrolled linked list of
  1024-slot blocks, like std::deque), behind a `Deque[T]` interface. It was work in progress. At their request, Claude
  fixed it and explained each fix. Bugs found: `resize`/`Resize` name mismatch (compile error); resize looped `i != tail`,
  which copies nothing when full (head == tail); popped slots not zeroed in the ring; linked PopBack read `tail`
  (off by one); pops never moved to or unlinked the next block (panic at index 1024); the first block started at 0, so
  PushFront allocated immediately; an unused `blocks` slice. Verified with a randomized oracle in the scratchpad.
  NOTE: the user was annoyed by a "doesn't compile" framing. WIP code is expected; don't phrase it as a failure.
- 2026-09-30 (Day 2): Day 1 quiz. `-u`/upstream: correct. `git status` before and after staging: wrong (said "clean"
  before staging). The `append` aliasing answer was half right (knew the header is copied, missed the shared backing array
  and cap). Interface holding a nil pointer: didn't know. Both internals questions came too early; they were
  turned into factoids and the quiz rule was changed. CLAUDE.md committed (ef16903).
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
