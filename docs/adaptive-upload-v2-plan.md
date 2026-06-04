# Adaptive Upload V2 Plan

## Purpose

The current upload path (`targets/go/file/uploadio.go` and `targets/go/file/manager/main.go`) uses a static `ConstrainedWorkGroup` semaphore sized by the operator. A single integer concurrency cap cannot serve the workload range the SDK needs to handle:

- Low scale (consumer internet, single-digit Mbps, lossy links): the default of 50 parallel parts overwhelms the link, triggers TLS-handshake storms, and produces worse throughput than 4-8 parallel parts would.
- Medium scale (office Gbps, daily batches): the default is roughly correct.
- High scale (10 Gbps WAN, multi-TB transfers, Apple-shaped workloads): direct-to-S3 paths may need 100-300 active parts to fill the pipe, while FIW and agent-mediated paths usually need lower clean concurrency. The controller must classify the target, choose sane bounds automatically, and ramp into the stable point instead of treating one global integer as correct everywhere.

A static default cannot satisfy all three. An adaptive concurrency controller can: it observes queue wait time, grows when there is headroom, shrinks when contention rises, and settles at the right level for each environment automatically.

This plan adds adaptive concurrency to the upload path as a new package (`file/uploadv2`) that runs alongside the existing implementation. V2 is explicit opt-in for SDK callers; CLI uploads default to V2 with `--adaptive-concurrency=false` as the opt-out. V1 remains available as the rollback path.

## Stages

This work ships in two stages. The V2 package structure, the `AdaptiveConcurrencyManager`, the sample-collection model, and the Vegas algorithm are the same in both. The difference is scope, defaults, and which environments are validated.

### Stage 1: adaptive upload core for all upload targets

Adds one adaptive upload engine that works across the three major CLI/SDK upload paths:

- Direct-to-S3 signed URLs, the dominant path and the path most likely to need 10 Gbps scaling.
- Files Integration Worker upload URLs, a more direct HTTP path than the agent but still service-mediated.
- Agent-mediated upload URLs, the most constrained path because requests flow through proxy load control, data lanes, and libp2p.

The agent work remains important because it is the strictest overload case and gives useful back-pressure behavior requirements, but it is not the global tuning target.

In scope for Stage 1:

- Loadcontrol module extraction (prerequisite for both stages).
- `lib.AdaptiveConcurrencyManager` implementing `ConcurrencyManager`.
- `file/uploadv2` package with adaptive upload.
- Target-aware dynamic part-size planner.
- Sample emission at part completion, with outcome mapping for S3, FIW, and agent paths.
- Automatic upload target classification from the resolved upload URL.
- CLI `--adaptive-concurrency` flag (default on; `false` opts out).
- Stderr telemetry.
- Defaults tuned by target class: S3 can grow higher, FIW sits in the middle, agent is capped lower and must respect proxy back-pressure. See Vegas Configuration.
- Opt-in checksum trailer support behind `upload-v2-checksum-trailer`; disabled by default and only applied for supported V2 destinations, currently direct AWS S3, when the upload URL signs the required trailer headers.
- Unit tests and parity tests against V1.

Stage 1 deliberately excludes default-on advanced destination-specific optimizations such as trailing checksums, throughput-reactive part sizing, and procurement-grade JSON reports. Basic direct-to-S3 adaptive concurrency and deterministic part sizing are in Stage 1. AWS checksum trailers may be tested through the separate `upload-v2-checksum-trailer` feature flag for supported destinations without changing the V1 upload route or the default V2 wire format.

Target: ~2 weeks of focused work after loadcontrol module extraction lands.

### Stage 2: high-scale direct-to-S3 optimization (Apple-shaped workloads)

Extends the Stage 1 S3 behavior with the features needed to compete with Aspera on the three Apple scale bands (low / medium / high) when uploading directly to S3 without FIW or the agent in the path.

Additions in Stage 2:

- Default-on trailing-checksum support after signed-URL support and compatibility validation (RFC 7230 HTTP trailers for native uploads, `aws-chunked` with `x-amz-trailer` for direct S3).
- Throughput-informed part-size tuning, if Stage 1 telemetry shows the static planner is leaving meaningful performance on the table.
- Pre-flight bandwidth probe to seed Vegas's initial target.
- Scale-band synthetic tests for low / medium / high simulated links.
- Optional CLI flag like `--transfer-profile {low,medium,high,auto}` to pre-select bounds, defaulting to auto.
- Structured JSON telemetry for procurement-grade evaluation reports.
- S3-specific optimizations: transfer acceleration support, BBR documentation runbook, kernel buffer tuning notes.

Stage 2 is additive. Stage 1 callers do not need to change to consume Stage 2 features; they are opt-in via additional flags and config fields.

Target: ~3 weeks of additional work after Stage 1 ships.

### What both stages share

- V2 package layout in `file/uploadv2/`.
- `AdaptiveConcurrencyManager` and its `loadcontrol.Bulkhead` + `algos.Vegas` internals.
- Sample-collection model and outcome mapping.
- V1 untouched throughout both stages.
- The same backwards-compat and promotion-criteria philosophy.

## Design Constraints

- V2 is a sibling package to the V1 upload code. V1 callers, tests, and behavior are not changed by this work.
- V2 implements the existing `lib.ConcurrencyManager` interface internally so existing upload helpers can be reused unchanged.
- V2's outer API mirrors V1's surface so callers switch packages without restructuring code.
- Adaptive behavior is driven by the `loadcontrol` Go module. The SDK takes a new dependency on it.
- The Vegas controller, bulkhead primitive, and sample model come from `loadcontrol`. No reimplementation.
- Defaults must work without operator tuning across low, medium, and high scale.
- V2's failure mode in pathological cases (no samples ever arrive, time goes backwards, bounded retries) must be no worse than V1's static behavior.
- Telemetry emits to stderr only in v1 of v2. Structured output (JSON, metrics endpoint) is a follow-up.

## Scope

| Item | Stage 1 | Stage 2 |
|---|---|---|
| Loadcontrol module extraction | ✅ | — (prerequisite shipped in Stage 1) |
| `lib.AdaptiveConcurrencyManager` | ✅ | — |
| `file/uploadv2` package | ✅ | extended |
| Target-aware part-size planner | ✅ | extended if telemetry proves need |
| Sample emission at part completion | ✅ | extended with new outcome paths |
| Vegas-tuned bounded concurrency | ✅ | extended with higher ceilings |
| Upload target classification | ✅ | extended |
| CLI `--adaptive-concurrency` flag | ✅ | — |
| Stderr telemetry (key=value lines) | ✅ | — |
| Unit and parity tests | ✅ | — |
| Scale-band synthetic tests | — | ✅ |
| Trailing-checksum support | feature-flagged supported-destination path, currently direct AWS S3 | default-on after validation |
| Throughput-reactive part sizing | — | optional |
| Pre-flight bandwidth probe | — | ✅ |
| `--transfer-profile` CLI flag | — | ✅ |
| Structured JSON telemetry | — | ✅ |
| High-throughput OS tuning command | Linux inspect/apply/verify command | macOS/Windows expansion |

Out of scope for both stages:

- Download path. V2 is upload-only initially.
- Modifying V1 or removing it before promotion criteria are met.
- Pre-emptive promotion of V2 to default behavior.

## Prerequisite

`loadcontrol` must be importable as a Go module. Today it lives in `files-integration-worker/lib/loadcontrol/`. Three viable options:

1. Promote to its own module (`github.com/Files-com/loadcontrol`). Both files-integration-worker and the SDK import from there.
2. Vendor the package into the SDK repo.
3. Cross-link via `go.work` during development; ship via option 1 later.

Option 1 is the right answer. The module extraction is bookkeeping (no code changes) but it is the gating step for any SDK integration.

## Package Layout

```
targets/go/file/uploadv2/
├── doc.go              # package overview, usage example
├── uploader.go         # V2 Uploader, mirrors file/uploader.go entry points
├── uploadio.go         # V2 uploadIO, mirrors file/uploadio.go Run/waitOnParts logic
├── partsizer.go        # target-aware dynamic part-size planner
├── sample.go           # sample emission helpers
├── telemetry.go        # stderr periodic output
└── *_test.go           # parity tests, scale-band tests

targets/go/lib/
└── adaptiveconcurrency.go   # AdaptiveConcurrencyManager (implements ConcurrencyManager)
```

The adaptive manager lives in `lib/` because it implements an existing interface that lives in `lib/`. The V2 upload logic lives in its own subpackage for isolation.

## Public API

V2 mirrors V1's public surface. The intent is that a caller can switch from `file.UploadWithResume(...)` to `uploadv2.UploadWithResume(...)` with no other code changes.

```go
package uploadv2

// Same shape as file.UploaderParams plus adaptive-specific fields.
type UploaderParams struct {
    file.UploaderParams       // embedded; reuses V1 params for everything else
    AdaptiveConfig AdaptiveConfig
}

type AdaptiveConfig struct {
    Enabled         bool            // default true; set false to fall back to V1 manager behavior
    TargetClass     TargetClass     // default auto; s3, fiw, agent, generic
    TransferProfile TransferProfile // default auto; low, medium, high, auto
    PartSize        PartSizeConfig  // default auto; target-aware planner
    InitialTarget  int             // default from TargetClass/Profile
    MinTarget      int             // default 2
    MaxTarget      int             // default from TargetClass/Profile
    MaxRampStep    int             // default from TargetClass/Profile
    RetryAfterRespect bool         // default true; 429/503 Retry-After pauses new acquires
    WaitFloor      time.Duration   // default from TargetClass/Profile
    WaitCeiling    time.Duration   // default from TargetClass/Profile
    MinSamples     int             // default 8
    Cooldown       time.Duration   // default 2s
    WaitHalfLife   time.Duration   // default 5s
    Telemetry      io.Writer       // optional; nil disables stderr output
}

type TargetClass string
const (
    TargetClassAuto    TargetClass = "auto"
    TargetClassS3      TargetClass = "s3"
    TargetClassFIW     TargetClass = "fiw"
    TargetClassAgent   TargetClass = "agent"
    TargetClassGeneric TargetClass = "generic"
)

type TransferProfile string
const (
    TransferProfileAuto   TransferProfile = "auto"
    TransferProfileLow    TransferProfile = "low"
    TransferProfileMedium TransferProfile = "medium"
    TransferProfileHigh   TransferProfile = "high"
)

type PartSizeConfig struct {
    Mode              PartSizeMode // default auto; fixed, auto
    FixedBytes        int64        // used when Mode=fixed
    UnknownSizeBytes  int64        // default from target class; used when size is nil
    MaxParts          int64        // default from target class; S3=10000, FIW/agent=0 unlimited
    MinPartBytes      int64        // default from target class; S3=5 MiB, FIW/agent=0
    MaxPartBytes      int64        // default from target class; S3=5 GiB, FIW/agent bounded by client memory policy
    PreferredPartBytes int64       // target/profile preference when constraints allow
}

type PartSizeMode string
const (
    PartSizeModeAuto  PartSizeMode = "auto"
    PartSizeModeFixed PartSizeMode = "fixed"
)

func UploadWithResume(ctx context.Context, params UploaderParams) (file.UploadResumable, error)
type Uploader = file.Uploader
```

Existing V1 types (`file.UploadResumable`, `file.Uploader`, `file.Job`) are reused via type aliases or embedding where possible. New code is the minimum needed to swap in adaptive concurrency.

## Concurrency Manager Implementation

`lib.AdaptiveConcurrencyManager` implements `lib.ConcurrencyManager`:

```go
type AdaptiveConcurrencyManager struct {
    bulkhead     *loadcontrol.Bulkhead
    orchestrator *loadcontrol.Orchestrator
    sink         *sampleCollector
}

func NewAdaptiveConcurrencyManager(cfg AdaptiveConfig) *AdaptiveConcurrencyManager
```

The manager:

- Constructs a `loadcontrol.Bulkhead` with `InitialTarget` as initial capacity and `MaxTarget` as the upper bound. `MaxTarget` is never the startup capacity.
- Constructs an `algos.Vegas` controller with the band/cooldown/samples settings.
- Constructs a `loadcontrol.Orchestrator` that ties them together and runs the periodic tick.
- Exposes the `ConcurrencyManager` interface methods (`Wait`, `Done`, `WaitWithContext`, `RunningCount`, etc.).
- Translates V1's acquire/release pattern into loadcontrol's Ticket lifecycle.
- Holds a sample collector that the V2 uploader code calls at part completion.
- Gates worker creation/acquisition through the adaptive bulkhead. The CLI should not spawn hundreds of active HTTP requests and then hope the proxy absorbs the burst; excess work waits locally until the controller admits it.
- Applies a client-side back-pressure pause when a response carries `Retry-After`. During the pause, existing in-flight parts may finish, but new acquires wait until the retry deadline or context cancellation.

## Target Classification

The SDK should infer `TargetClass` from the resolved upload URL when the caller leaves it as `auto`:

| Target class | Detection | Share of expected CLI usage | Default behavior |
|---|---|---:|---|
| `s3` | Signed URL host/path matches S3, S3 accelerate, or known S3-compatible host pattern | ~65% | Highest ceiling; optimize for direct object-store throughput |
| `fiw` | Files Integration Worker upload URL or Files.com service-mediated upload host that is not agent proxy | ~25% | Middle ceiling; service-mediated HTTP path with back-pressure |
| `agent` | Agent proxy upload URL, remote server agent-v2 mount, or explicit metadata from begin-upload response | ~10% | Lower ceiling; must cooperate with proxy/agent load control |
| `generic` | Unknown endpoint | fallback | Conservative direct-HTTP behavior |

The target class chooses defaults only. The same adaptive controller, outcome mapping, retry handling, telemetry, and API surface are shared. A caller may override `TargetClass` or `TransferProfile`, but normal CLI users should not need to.

## Part Size Strategy

Part sizing is independent from adaptive concurrency. Concurrency is a live control loop because queue wait, back-pressure, and transport failures are only visible during the transfer. Part size should be deterministic for a given upload target and known file size. Changing part size on the fly based on short-window throughput adds another controller, makes resume state harder to reason about, and can produce unstable behavior when the concurrency controller is already adapting.

Stage 1 should add a small target-aware planner:

```go
type PartSizer interface {
    Next(offset int64, partNumber int) (size int64, final bool, ok bool)
}

func NewPartSizer(target TargetClass, totalSize *int64, cfg PartSizeConfig) PartSizer
```

The planner replaces the V2 copy of `ByteOffset` sizing logic. V1 remains untouched.

### Known-size uploads

When the input size is known, compute the part size before scheduling work:

1. Start with the target's size-tier preferred size.
2. Apply the target's protocol constraints.
3. Ensure the chosen size keeps the upload within the target's max part count when one exists.
4. Emit a final smaller part when needed.

For direct S3, the planner must obey S3 multipart constraints:

- Minimum part size is 5 MiB except for the final part.
- Maximum part count is 10,000.
- Maximum object size is 5 TiB.
- Maximum individual part size is 5 GiB.
- For known-size uploads, choose `max(size_tier_preference, 5 MiB, ceil(total_size / 10000))`, round up to a MiB boundary, and fail or fall back if that would exceed the S3 part-size/object-size limits.

S3 known-size tier preferences:

| Total size | Preferred part size | Why |
|---:|---:|---|
| `< 1 GiB` | 16 MiB | Keeps enough parts for adaptive concurrency without inflating small uploads |
| `1 GiB - < 16 GiB` | 32 MiB | Reduces request overhead while leaving hundreds of parts for multi-GiB uploads |
| `16 GiB - < 512 GiB` | 64 MiB | Baseline high-throughput size for large S3 uploads without excessive request count |
| `512 GiB - < 1 TiB` | 128 MiB | Lowers request count for very large uploads |
| `1 TiB - < 2 TiB` | 256 MiB | Keeps huge uploads efficient without jumping directly to maximum-sized parts |
| `>= 2 TiB` | 512 MiB, then larger if required by `ceil(total_size / 10000)` | Balances request overhead with S3's 10,000 part ceiling |

For known-size direct-S3 uploads, V2 may reduce the tier preference when the aggregate job workload would otherwise underfeed the adaptive controller. The planner estimates the job's valid upload bytes, targets `initial S3 target * 8` planned parts, rounds down to a power-of-two MiB part size, and never goes below 8 MiB or S3's 5 MiB minimum. The selected size remains fixed for that upload session. This gives short multi-file jobs and single large files enough schedulable parts to ramp, while large aggregate jobs keep the larger tier size and avoid unnecessary request churn.

For FIW and agent-mediated uploads, the contract is less strict:

- No protocol-level 5 MiB minimum.
- No S3-style 10,000 part limit.
- Part size is a client efficiency choice, not a correctness boundary.
- Defaults should avoid tiny parts that amplify request overhead, but they do not need to inflate part size just to satisfy S3.

Known-size defaults for non-S3 targets:

| Target class | Total size | Preferred part size | Min part size | Max parts | Notes |
|---|---:|---:|---:|---:|---|
| `fiw` | `< 16 GiB` | 32 MiB | none | unlimited | Balanced request overhead and retry cost for normal FIW uploads |
| `fiw` | `16 GiB - < 256 GiB` | 64 MiB | none | unlimited | Larger request unit for sustained high-throughput uploads |
| `fiw` | `>= 256 GiB` | 128 MiB | none | unlimited | Avoids excessive request churn for huge FIW uploads |
| `agent` / `generic` | `< 8 GiB` | 16 MiB | none | unlimited | Smaller retry unit; cooperates with stricter proxy envelope |
| `agent` / `generic` | `8 GiB - < 128 GiB` | 32 MiB | none | unlimited | Better efficiency for larger agent-mediated uploads |
| `agent` / `generic` | `>= 128 GiB` | 64 MiB | none | unlimited | Conservative huge-upload request reduction |

These are startup defaults, not a live throughput controller. The chosen part size stays fixed once the upload session starts so resume state, part offsets, ETag ordering, and retry cost remain predictable. The concurrency controller adapts during the transfer; the part-size planner keeps the unit of retry and the protocol envelope stable.

### Unknown-size uploads

Unknown-size inputs should not attempt high-scale optimization. Without total size, the client cannot prove an S3 upload will stay below 10,000 parts. The planner should use a simple target default and preserve today's streaming behavior:

- S3/direct signed URL: use the existing S3-safe static sequence or a conservative fixed default that matches the server contract. If the stream exceeds what can be represented within S3 constraints, fail clearly rather than silently creating an invalid multipart upload.
- FIW/agent/generic: use the target's unknown-size default, initially 5 MiB or 8 MiB, because there is no protocol max-part-count pressure.

Unknown-size uploads are about correctness, bounded memory, and predictable retry units. They are not the path to optimize for 10 Gbps.

### On-the-fly resizing

Do not resize parts mid-upload in Stage 1. Throughput-reactive part sizing should remain a Stage 2 experiment only if telemetry shows it is needed. If added later, it should use coarse, deterministic epochs rather than per-part twitching:

- Only increase part size after enough successful samples at a stable concurrency target.
- Never decrease below already-advertised protocol constraints.
- Preserve resume metadata so retries use the same offset and length as the original part.
- Do not couple the part-size decision directly to the same short-window samples that drive Vegas, or the two controllers can fight each other.

The trick: V1's `ConstrainedWorkGroup` has `Wait()/Done()` without per-call sample data. V2 needs to attach samples on `Done()`. Two options:

1. Add a `DoneWithSample(s loadcontrol.Sample)` method to a V2-specific subinterface. V2 upload code uses this; V1 uploadio code still uses `Done()`.
2. Use a per-goroutine context to carry sample data; `Done()` reads from the context.

Option 1 is cleaner and explicit. V2 code calls `DoneWithSample`; V1 code is unchanged.

## Sample Collection

After each part upload completes (in V2's `uploadio.go`), emit a sample to the adaptive manager:

```go
sample := loadcontrol.Sample{
    Outcome:        outcomeFromHTTPResult(resp, err),
    SlotWaitNs:     acquired.Sub(requested).Nanoseconds(),
    SlotOccupancyNs: time.Since(acquired).Nanoseconds(),
    Bytes:          partSize,
    InFlight:       manager.RunningCount(),
    At:             time.Now(),
}
manager.DoneWithSample(sample)
```

Outcome mapping (consistent with the proxy-side mapping in `data_plane_route_limit_conn.go`):

| Condition | Outcome | Vegas treatment |
|---|---|---|
| 2xx response | `OutcomeSuccess` | Normal sample; feeds slot-wait and goodput EWMAs |
| 429 response | `OutcomeServerBackPressure` | Honor `Retry-After`; shrink target; feed back-pressure rate |
| 503 response with `Retry-After`, or known agent busy response | `OutcomeServerBackPressure` | Same as 429; this is overload, not an ordinary server bug |
| Other 5xx response | `OutcomeServerError` | Excluded from controller (bottleneck is past our layer) |
| `net.OpError` or RST_STREAM | `OutcomeStreamReset` | Strong shrink trigger; bypasses EWMA smoothing |
| Context deadline or timeout | `OutcomeTimeout` | Sample with capped RTT contribution |
| Other transport error | `OutcomeConnectionLost` | Strong shrink trigger; same high-signal treatment as `OutcomeStreamReset` |

`OutcomeConnectionLost` and `OutcomeStreamReset` now receive identical treatment in the controller. The distinction is preserved in telemetry for operator debugging, but the algorithmic response is the same: immediate multiplicative shrink and a short jittered cool-off before replacement work is admitted. This matches the proxy-side controller behavior so client and server adapt symmetrically when either side observes connection-level failure.

`Retry-After` is authoritative for overload responses. Parse both delta-seconds and HTTP-date forms. Clamp negative values to zero, apply full jitter, and cap the local pause at a small bounded value (for example 30 seconds) so one bad header cannot freeze an upload indefinitely. The retry deadline should be visible in telemetry.

The `Sample.SlotOccupancyNs` field is auto-populated by `loadcontrol.Ticket.enrich` at release time (acquire-time to release-time delta). The V2 emitter does not need to compute it explicitly. The field is not consumed by Vegas in the client; it exists for the same reason it does on the proxy side, in case future client-side estimators (e.g., predicting agent-side queue state from observed hold times) need it.

Helper: `outcomeFromHTTPResult(resp *http.Response, err error) loadcontrol.Outcome` in `sample.go`.

## Vegas Configuration

Defaults are tuned for CLI workloads (shorter sessions than the proxy) and selected by upload target class. The defaults must be good enough without operator tuning; config knobs are escape hatches.

### Stage 1 auto defaults

| Parameter | Value | Rationale |
|---|---|---|
| `InitialTarget` | 8 | Useful immediately on medium scale; low enough for low scale to not flood |
| `MinTarget` | 2 | Floor; even worst link benefits from 2 over 1 |
| `RetryAfterRespect` | true | 429/503 overload responses are explicit load-control signals |
| `MinSamples` | 8 | Lower than proxy (20); CLI sessions are shorter, needs faster convergence |
| `Cooldown` | 2 s | Faster than proxy (5 s); CLI users see results sooner |
| `WaitHalfLife` | 5 s | Faster than proxy (10 s); shorter session timescale |

Target-specific bounds:

| Target class | `MaxTarget` | `MaxRampStep` | `WaitFloor` | `WaitCeiling` | Rationale |
|---|---:|---:|---:|---:|---|
| `s3` | 1024 hard cap with 150 soft growth ceiling | measured enterprise soft ceiling | throughput/loss driven below cap | throughput/loss driven | Dominant path; direct object-store uploads start at the measured 150-active-part C3/S3 plateau. Stage 1 VM testing validates the default plateau for the tested workloads, not field growth above 150. Higher headroom stays available, but the soft ceiling only unlocks after sustained bytes and throughput prove the extra connection footprint is worth probing. |
| `fiw` | 192 | 6 | 25 ms | 225 ms | Service-mediated HTTP path; higher than agent, lower than direct S3 until FIW-specific load behavior proves more headroom |
| `agent` | 128 | 4 | 30 ms | 250 ms | Proxy/agent path has the strictest server-side envelope and explicit load-control feedback |
| `generic` | 64 | 4 | 30 ms | 250 ms | Unknown endpoint fallback favors stability |

`TransferProfile` modifies only the ceiling:

| Profile | Effect |
|---|---|
| `low` | cap `MaxTarget` at 16 |
| `medium` | cap `MaxTarget` at 64 |
| `high` | use the target-class ceiling |
| `auto` | use the target-class ceiling and let Vegas settle below it |

The direct-S3 Stage 1 default uses `InitialTarget=150` and a `150` soft growth ceiling because cloud C3-to-S3 validation of `--adaptive-concurrency` with no numeric tuning matched or stayed within 5% of the best static baseline for `20x200MiB`, `200x200MiB`, and `1x20GiB`. The closest result was the `200x200MiB` workload at roughly 4.90% slower, so this should be treated as a thin no-regression margin rather than a comfortable win. Those field workloads are all below the default 64 GiB ceiling-unlock threshold, so they validate the default 150 plateau and shrink path, not the grow-above-150 field path. The above-150 path is unit-tested and remains available through `MaxTarget=1024`, but it only unlocks after enough bytes and measured throughput indicate a larger transfer on a fast path. This is still adaptive for slower links: the controller can shrink on latency, failures, 429/503, or transport back-pressure. The HTTP transport uses the soft ceiling for benchmark-sized workloads and opens the larger cap for workloads that are large enough to justify high-ceiling probing.

S3 keeps idle HTTP retention capped below the active ceiling. Directory/job uploads must share one adjusted V2 HTTP client per target class so many-file jobs do not create one transport pool per file. This avoids accumulating hundreds of idle established connections on many-file jobs.

Adaptive upload V2 uses a separate file-admission cap for multi-file jobs. The default is `128` files while the part-concurrency manager remains capped at `1024` and adapts toward the active target. The benchmark decision record is [Adaptive Upload V2 File-Admission Cap Decision](adaptive-upload-v2-file-concurrency-cap.md). In short, a tuned Tier_1/GVNIC VM benchmark compared file-admission caps `300`, `128`, and `50` while holding adaptive part concurrency constant. Cap `128` produced the best median throughput on both tested datasets while using less memory than `300`; cap `50` underfed the `1000x64MiB` workload.

### Diagnostic tuning knobs

V2 keeps normal users on automatic defaults: `--adaptive-concurrency` without numeric tuning is the product contract and the benchmark success path. Numeric benchmark exploration must not require rebuilding the CLI, so the CLI exposes hidden S3-only diagnostic flags that map to `file.UploadV2Tuning` and only take effect when adaptive upload V2 is enabled:

| Area | Hidden flags |
|---|---|
| Startup and floor | `--adaptive-upload-v2-s3-initial-target`, `--adaptive-upload-v2-s3-adaptive-floor` |
| Growth cadence | `--adaptive-upload-v2-s3-grow-every`, `--adaptive-upload-v2-s3-grow-step` |
| Throughput probe behavior | `--adaptive-upload-v2-s3-throughput-window`, `--adaptive-upload-v2-s3-throughput-min-gain-percent`, `--adaptive-upload-v2-s3-probe-min-windows`, `--adaptive-upload-v2-s3-probe-floor-target`, `--adaptive-upload-v2-s3-probe-floor-rate-bps`, `--adaptive-upload-v2-s3-probe-plateau-target`, `--adaptive-upload-v2-s3-throughput-shrink-percent`, `--adaptive-upload-v2-s3-throughput-hold-windows`, `--adaptive-upload-v2-s3-probe-min-gain-per-target-percent`, `--adaptive-upload-v2-s3-probe-loss-tolerance-percent` |
| S3 soft ceiling | `--adaptive-upload-v2-s3-growth-ceiling`, `--adaptive-upload-v2-s3-growth-ceiling-probe-bytes`, `--adaptive-upload-v2-s3-growth-ceiling-probe-successes`, `--adaptive-upload-v2-s3-growth-ceiling-probe-rate-bps` |
| Latency protection | `--adaptive-upload-v2-s3-latency-queue-high`, `--adaptive-upload-v2-s3-latency-growth-queue-high` |
| Part-size planning | `--adaptive-upload-v2-s3-part-size-mib`, `--adaptive-upload-v2-s3-workload-bytes`, `--adaptive-upload-v2-s3-workload-target-part-multiplier`, `--adaptive-upload-v2-s3-workload-min-part-size-mib`, `--adaptive-upload-v2-s3-workload-scan-wait-ms` |
| Ready runway | `--adaptive-upload-ready-runway-parts`, `--adaptive-upload-ready-runway-bytes` |
| File admission | `--adaptive-upload-v2-file-concurrency` |

These are intentionally hidden escape hatches for benchmark runs, customer diagnostics, and rollout tuning. They should not be promoted as normal CLI UX, and changing the values should never change V1 upload behavior.

### High-Throughput OS Tuning

The SDK exposes a reusable high-throughput OS tuning workflow in `lib/ostuning` and the CLI surfaces it as `files-cli os-tuning high-throughput`. On Linux, the command supports plan, verify, repair, and restore modes for BBR, `fq`, TCP buffer limits, ephemeral port range, `tcp_slow_start_after_idle`, and supported NIC ring sizes. The CLI command is intentionally outside upload execution; users opt into host-level repair separately from the upload command, and SDK callers such as Desktop Helper can reuse the same inspect/plan/repair/restore API.

Recent local proxy/agent testing with the split-plane controller showed:

- `128` instant diagnostic workers became clean after the proxy queue default moved to `256`: no stream failures, no Retry-After responses, no connection resets, and roughly 8 Gbit/s on a Mac development loopback path.
- `200` instant workers completed and reached roughly 9 Gbit/s, but produced stream-reset/server-error/timeout samples. This is above the clean default target for an instant agent-path burst.
- `500` instant workers did not crash the agent, but it produced deeper queueing, more reset/error samples, and lower throughput. The right client behavior at this level is local throttling and adaptive reduction, not insisting that all requested work run concurrently.

For the agent path, the proxy's per-connection slot cap is bounded by transport: TCP and QUIC start at `16` streams per data connection and may grow to `32` only after the proxy is pinned at the lane ceiling and measured pressure remains. The proxy's default max lane count is `16`, but it starts low and earns its way up through AIMD plus bounded pressure growth. The client should therefore expect clean agent-path equilibrium commonly around `16-128` active parts, depending on path quality and file size, and should treat values above `128` as non-default escape hatches that require measured proof.

For the agent target, `MaxTarget=128` provides headroom for scenarios where the proxy has grown the lane count well past default (e.g., 4 lanes × 32 QUIC slots = 128 streams), but the client should never push past the proxy's actual envelope. The proxy is the source of truth for what the agent will accept; the client adapts to what the proxy admits.

If a user supplies `--concurrent-connection-limit=75`, adaptive mode should not open 75 HTTP requests immediately. It should set `MaxTarget=75`, start at `InitialTarget=8`, and ramp up by at most the target-class `MaxRampStep` per cooldown while success/goodput/queue wait justify the increase. If the server returns 429/503 with `Retry-After`, or the client observes connection loss/reset, the controller shrinks and pauses new work before retrying.

When the CLI's `--concurrent-connection-limit` is set in adaptive mode, it becomes `MaxTarget` for both stages, overriding the stage default. It is a ceiling, not a target and not startup concurrency. Existing customers who need old static behavior can disable adaptive mode.

SDK jobs that do not provide an explicit manager share one default process-wide manager. This protects Desktop, FIW, and other multi-job callers from multiplying file, part, and directory-listing concurrency every time they start a separate upload job. The shared default manager is initialized atomically once per process. For adaptive uploads, a nil manager remains a job-scheduling and V1-fallback default, not a V2 upload cap, so V2 still uses target-specific concurrency caps. V2 adaptive upload learning is also shared across jobs by target class, effective max cap, and manager-relevant tuning so a directory job, individual file jobs, and retry jobs all benefit from the same learned concurrency envelope. Passing an explicit manager remains the escape hatch for a caller that needs a separate cap or tenant-specific limiter.

## Telemetry

Periodic stderr output, default every 5 seconds while a transfer is in progress:

```
adaptive-upload: target=24 active=22 queued=0 wait_p95=12ms goodput=180MB/s outcome_success=487 outcome_429=12 outcome_5xx=0
```

Disabled when `AdaptiveConfig.Telemetry == nil`. Caller can supply `os.Stderr`, an `io.Discard`, or a custom writer.

Format: one line per emission, space-separated key=value, machine-parseable. Include at least: `target_class`, `transfer_profile`, `target`, `active`, `queued`, `wait_p95`, `goodput`, `outcome_success`, `outcome_429`, `outcome_503_busy`, `outcome_stream_reset`, `outcome_connection_lost`, `outcome_timeout`, `retry_after_active`, and `retry_after_ms`.

Part-size telemetry should be included at transfer start and when periodic telemetry emits: `part_size`, `part_size_mode`, `known_size`, `part_number`, and `estimated_total_parts` when the total size is known.

## Coordination with Server-side Load Control

Different upload targets expose different server-side bottlenecks. The client uses the same local controller for all of them, but it must interpret target-specific feedback correctly.

### Direct-to-S3

S3 is the dominant path and should be the high-throughput default target. The client should expect fewer explicit back-pressure signals than the agent path, so it relies more heavily on observed goodput, queue wait, timeouts, and transport errors. `MaxTarget=1024` remains available as hard headroom, but Stage 1 starts with a `150` soft ceiling because it is the measured stable C3/S3 point; higher growth requires sustained bytes and throughput, and Stage 2 should add a pre-flight bandwidth probe before raising that soft ceiling automatically.

### Files Integration Worker

FIW is a service-mediated HTTP path without libp2p, so it should tolerate more client concurrency than the agent path but less than direct S3 by default. It should honor ordinary HTTP overload feedback (`429`, `503`, `Retry-After`) and transport failures the same way as the agent path. The initial FIW ceiling of `192` is deliberately below direct S3 until FIW-specific validation proves a higher default is clean.

### Agent Proxy

The proxy ships adaptive controllers default-on with independent disable flags for fast rollback (see `agent-proxy-dynamic-concurrency-plan.md` and `agent-proxy-adaptive-back-pressure-plan.md`). The proxy operates at two levels:

- Per-connection slot cap: Vegas-style controller, bounded by transport. TCP/QUIC start at `16` concurrent streams and can grow to `32` when lane count is pinned and pressure remains.
- Per-peer data connection count: AIMD-flavored counter plus bounded pressure shortcut. The proxy starts with low lane count, can grow toward `16` data connections, and uses a bounded queue (`256` by default) to absorb short bursts while new lanes register.

The V2 client cooperates with the proxy through the same neutral `Outcome` enum and the same Vegas-based algorithm in the loadcontrol package. When the proxy issues a 429/503 + Retry-After, the V2 client's sample (`OutcomeServerBackPressure`) drives a shrink decision and pauses new work; when the client backs off, the proxy's queue wait drops and its controller stabilizes. Both ends converge on the same operating point.

The V2 client deliberately does not try to predict the proxy's controller state. It reacts to the signals the proxy emits (429s, transport-level failures, Retry-After hints) the same way the proxy reacts to signals from the agent. Symmetric back-pressure all the way down.

Despite the proxy being default-on, SDK uploads stay opt-in so existing integrations and shared tenant services keep their current static manager behavior. The CLI can default to adaptive because it owns its transfer manager and exposes `--adaptive-concurrency=false` as the per-command rollback path.

## Backwards Compatibility

- V1 (`file.UploadWithResume`, `file.Uploader`, the existing `manager.Manager` and `ConstrainedWorkGroup`) is untouched.
- Existing SDK consumers see no change unless they explicitly opt in with `UploadWithV2()` or `UploaderParams.AdaptiveConcurrency`.
- `FeatureFlagAdaptiveUploadV2` by itself does not upgrade existing SDK uploads. It remains a product/config gate, but upload code still requires an explicit per-upload or per-job opt-in.
- After opt-in, adaptive upload uses the global shared adaptive manager by default so multi-file and multi-job SDK usage can learn from the broader workload.
- When an SDK caller explicitly supplies a `manager.Manager` with adaptive upload enabled, that manager is treated as an isolation boundary and V2 uses a per-upload adaptive manager capped by the supplied manager.
- V2's `AdaptiveConfig.Enabled = false` falls back to V1's `ConstrainedWorkGroup` behavior using `MaxTarget` as the static cap. Provides a per-call disable without removing the V2 import.
- Default `AdaptiveConfig` values are tuned to behave at least as well as V1's defaults at medium scale, where V1 is roughly correct.
- The CLI's `--adaptive-concurrency` flag selects V2 by default; pass `--adaptive-concurrency=false` to use the V1/static upload path.

## Test Strategy

Three test layers. Stage 1 owns target classification, parity, and basic integration for S3/FIW/agent. Stage 2 adds deeper high-scale direct-S3 validation and S3-only features.

### Stage 1 tests

**Unit tests** in `file/uploadv2/`:

- Adaptive manager grows under low queue wait + sustained demand.
- Adaptive manager shrinks under high queue wait.
- Adaptive manager respects `MinTarget` and `MaxTarget` bounds.
- Part-size planner respects S3 minimum part size, maximum part count, maximum object size, and final-part exception.
- Part-size planner does not apply S3 minimum or max-part inflation to FIW/agent/generic targets.
- Unknown-size S3 uploads use the conservative static sequence and fail clearly when the stream exceeds representable multipart bounds.
- Unknown-size FIW/agent uploads use the target default without S3 max-part-count behavior.
- Sample emission produces correct outcomes for each HTTP status / error type.
- `AdaptiveConfig.Enabled = false` uses static behavior; identical to V1.
- Telemetry output is parseable and respects the writer config.

**Parity tests** comparing V1 and V2 against the same mock server:

- Small upload (<10 MB, single PUT path): V2 wall-clock within 10% of V1.
- Medium upload (100 MB, ~20 parts): V2 wall-clock within 10% of V1.
- Large upload (10 GB, ~500 parts) with simulated server back-pressure: V2 has fewer 429s and similar or better wall-clock than V1.

**Target-class integration tests** (added in Stage 1):

- Direct-to-S3 signed URL mock: starts at the S3 initial target or known part count, can sustain 100-150 active parts when the server is healthy, and shrinks cleanly when overload or transport loss appears.
- FIW-like mock: emits ordinary HTTP `429`/`503` + `Retry-After` under load. V2 should adapt down, settle, and produce at least as good wall-clock as V1 with significantly fewer overload responses observed.
- Target classifier tests for S3, S3 accelerate, FIW/service-mediated URLs, agent proxy URLs, and unknown generic URLs.
- Round-trip against a mock proxy that emits 429 + Retry-After under load. V2 should adapt down, settle, and produce at least as good wall-clock as V1 with significantly fewer 429s observed.
- Round-trip against a mock proxy that behaves like the tuned split-plane proxy: starts at 1 data lane × 16 streams, has queue depth 256, can grow lanes toward 16, and can grow streams toward 32 only after lanes are pinned. With `MaxTarget=128`, the client should ramp without initial 429s, settle without sitting at the ceiling forever, and record zero connection-reset samples in the healthy case.
- Synthetic overload with `MaxTarget=200` and `MaxTarget=500`: the client must keep failures bounded by shrinking/pacing locally. Success criteria are no unbounded goroutine growth, no retry storm, no repeated immediate 429 loop, and telemetry showing back-pressure or reset/loss samples caused target reduction.

### Stage 2 tests

**Scale-band synthetic tests** (skipped by default; run with `-tags scale`):

- Low-scale link (4 Mbps simulated, 100ms RTT, 1% packet loss): V2 should settle at 4-8 parallel parts; throughput within 80% of theoretical max for the link.
- Medium-scale link (1 Gbps simulated, 10ms RTT, 0% loss): V2 should settle 30-70 parallel parts; throughput within 90% of link capacity.
- High-scale direct-to-S3 link (10 Gbps simulated, 50ms RTT, 0% loss): Stage 1 should hold the measured 150-active-part range without manual tuning. Stage 2 may raise the ceiling above 150 only after a bandwidth probe or comparable telemetry shows the extra connection footprint produces meaningful goodput. This expectation does not apply to FIW or agent-mediated paths, where the server-side envelope and back-pressure decide the clean concurrency.

The synthetic tests use a mock HTTP server with bandwidth/latency injection to approximate each band.

**Trailing-checksum integration tests:** wire-format correctness against a fake server that consumes HTTP trailers (and `aws-chunked` for the S3 path). Mismatch detection. Per-part verification. V2 must only use the trailer path for supported destinations, currently direct AWS S3, when the upload URL signs `content-encoding`, `x-amz-content-sha256`, `x-amz-decoded-content-length`, `x-amz-sdk-checksum-algorithm`, and `x-amz-trailer`.

**Throughput-reactive part-size tests, if implemented:** correct epoch transitions, resume metadata stability, multipart-completion composite hash verification, and no oscillation when concurrency is also adapting.

## Promotion Criteria

Two transitions to track.

### Stage 1 → Stage 2

Stage 2 work can begin in parallel with Stage 1 validation, but Stage 2 features should not be enabled by default until Stage 1 has met these criteria:

1. Stage 1 parity test suite passes consistently (10 consecutive runs, no flakes).
2. Stage 1 target-class integration tests pass for S3, FIW, and agent paths.
3. CLI default-on validation passes for representative direct-S3, FIW, and agent-path workloads, with `--adaptive-concurrency=false` confirmed as the rollback path.
4. Telemetry shows the controller converges, doesn't oscillate, and doesn't sit at target-class `MaxTarget` permanently except when the target is demonstrably ceiling-limited.
5. Agent-path overload validation shows `MaxTarget=200` and `MaxTarget=500` are handled by client-side throttling and target reduction, not by repeated immediate 429s, connection-reset loops, or agent instability.
6. Direct-S3 validation shows the default S3 ceiling reaches the measured enterprise throughput band without manual `--concurrent-connection-limit` tuning.

When those hold, Stage 1 is considered shipped. Stage 2 features can be enabled in parallel without invalidating Stage 1's correctness.

### V2 → V1 promotion

V2 replaces V1 (V1 deprecated, eventually removed) when all of the following hold:

1. Both stages have shipped.
2. Stage 1 promotion criteria above are met for S3, FIW, and agent-path workloads.
3. Stage 2 scale-band tests pass for low / medium / high simulated links.
4. Stage 2 direct-to-S3 validation shows checksum and optimization behavior does not regress transfer success rate, resume behavior, or wall-clock time.
5. Documented runbook for operator-visible behaviors.

When those hold, a follow-up change:

- Renames `file.UploadWithResume` to `file.UploadWithResumeV1` and deprecates it.
- Promotes `uploadv2.UploadWithResume` to `file.UploadWithResume`.
- Removes `uploadv2/` after one release of overlap.
- Removes `ConstrainedWorkGroup` after V1 callers are gone.

## Rollout

### Stage 1 rollout (adaptive upload core)

1. Land `loadcontrol` as an importable Go module. No SDK code changes yet.
2. Land `lib.AdaptiveConcurrencyManager` with unit tests. No callers yet.
3. Land target classification and target-specific defaults for S3, FIW, agent, and generic URLs.
4. Land the V2 target-aware part-size planner with unit tests. No caller default change yet.
5. Land `file/uploadv2/` with Stage 1 defaults. Parity tests against V1. SDK callers remain explicit opt-in.
6. Land CLI `--adaptive-concurrency` flag. Default on with `--adaptive-concurrency=false` rollback.
7. Internal validation against S3, FIW, and agent-path mocks, including back-pressure and retry-after behavior.
8. Validate CLI default-on behavior across S3-heavy, FIW-heavy, and agent-heavy workloads.
9. Keep SDK adaptive upload explicit opt-in until customer and service owners choose where shared adaptive learning is appropriate.

Stage 1 ends here. The CLI has one adaptive upload controller that works across direct S3, FIW, and agent-mediated uploads.

### Stage 2 rollout (direct-to-S3 optimization)

11. Introduce `--transfer-profile` flag if Stage 1 field telemetry shows profile overrides are still useful.
12. Validate the V2-only `upload-v2-checksum-trailer` feature flag against backend-generated direct-S3 URLs that sign the required AWS trailer headers. The direct-S3 trailer algorithm is fixed to CRC32C until the upload response can carry a signed algorithm value.
13. Evaluate throughput-reactive part-size epochs only if the Stage 1 planner is measurably insufficient.
14. Add pre-flight bandwidth probe (optional, behind a flag initially).
15. Add scale-band synthetic test suite (`-tags scale`).
16. Internal validation across simulated low / medium / high links.
17. Field validation with Apple POC (or comparable customer with mixed-scale workload).
18. Add structured JSON telemetry for procurement-grade reports.
19. Promote V2 to default upload path per the promotion criteria.

## Pitfalls

- **Cold start.** First N samples have low signal. `MinSamples=8` is the gate. Sessions shorter than ~10 seconds may never adapt; behavior falls back to `InitialTarget=8`. Acceptable since most production workloads run longer.
- **Mock server idiosyncrasies in tests.** Mock HTTP servers in Go tests don't have realistic queueing characteristics. Parity tests should not over-fit to mock behavior; scale-band tests are the more meaningful validation.
- **Outcome misclassification.** If `outcomeFromHTTPResult` returns the wrong outcome (e.g., labels a 5xx as `ServerBackPressure`), Vegas misreacts. Test exhaustively against each error class.
- **The V1/V2 split temptation.** Keep V2-only high-scale features behind explicit feature flags while V2 is under validation. The `upload-v2-checksum-trailer` flag is additive and must not alter V1 upload behavior or the default V2 wire format.
- **MaxTarget semantics.** If `--concurrent-connection-limit` is left at 50, V2's `MaxTarget` is 50 and high-scale customers get the same ceiling as V1. The flag becomes the ceiling, not the target. In adaptive mode it must not become startup concurrency.
- **Per-job default managers.** Default SDK jobs must not create independent managers. Desktop/FIW can start multiple directory or file jobs concurrently; without a shared default manager, aggregate concurrency multiplies across jobs and defeats the adaptive controller.
- **Two controllers fighting.** If part size and concurrency both adapt from the same short-window throughput samples, the system can oscillate and make regressions hard to explain. Stage 1 part sizing should be deterministic.
- **Applying S3 constraints everywhere.** FIW and agent uploads do not need S3's 5 MiB minimum or 10,000-part correction. Applying those constraints globally makes retry units larger than necessary and hides target-specific behavior.
- **Unknown-size over-optimization.** Unknown-size uploads cannot guarantee S3 part-count correctness from total size. Keep the algorithm simple, bounded, and explicit about failure when a stream exceeds representable bounds.
- **Adaptive output that hovers at MaxTarget.** Indicates the ceiling is the bottleneck, not the controller's optimum. Telemetry should make this obvious (target stuck at config max). Document the diagnostic.
- **Goroutine count under high MaxTarget.** With `MaxTarget=256` and 256 parallel parts, the client process runs 256+ goroutines plus their HTTP transport state. Test client memory footprint at the high end before flipping defaults.
- **Retry-After ignored.** This is the failure mode observed in CLI testing against the agent path and it applies to FIW/service-mediated paths too: the client keeps feeding work while the server is saying it is busy. Treat this as a release blocker for Stage 1.
- **Instant-burst overfitting.** The proxy can absorb very large bursts, but the SDK should not require it to. The client should ramp and hold work locally; the proxy queue is a safety buffer, not the primary scheduling mechanism.

## References

- Loadcontrol package (currently in `files-integration-worker/lib/loadcontrol/`).
- Loadcontrol algorithms plan: `files-integration-worker/docs/loadcontrol-algorithms-plan.md`.
- V1 upload path: `targets/go/file/uploader.go`, `targets/go/file/uploadio.go`.
- V1 concurrency primitive: `targets/go/lib/constrainedwaitgroup.go`.
- V1 manager: `targets/go/file/manager/main.go`.
