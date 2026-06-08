# Adaptive Download V2 Output Baseline

Adaptive download V2 uses a keep-size temp-file output path for native known-size ranged downloads:

- Create or open the `.download` temp file on the destination volume.
- Preallocate it to the final file size before scheduling range workers.
- Fill the copy buffer before writing so small HTTP/TLS reads are coalesced into larger `WriteAt` calls.
- Stream each coalesced range chunk directly into its final byte offset with `WriteAt`.
- On incomplete exit, truncate the temp file back to the contiguous downloaded prefix before returning.
- Finalize with the existing temp-file move path only after all ranges complete.

Accepted tuning for this MR:

| Path | VM | Disk | Object Size | Elapsed | Throughput | Decision |
| --- | --- | --- | ---: | ---: | ---: | --- |
| Adaptive V2 download with keep-size temp-file preallocation | n2-highcpu-32, us-west2-c | 4 local NVMe RAID0 | 100 GiB | 94.9s | 9.05 Gbps | Previous accepted baseline |
| Adaptive V2 download with coalesced 1 MiB `WriteAt` chunks | n2-highcpu-32, us-east4-a | 4 local NVMe RAID0 | 100 GiB | 66.1s | 12.99 Gbps | Accepted tuning |
| Adaptive V2 download with coalesced 1 MiB `WriteAt` chunks | z3-highmem-22-highlssd, us-east4-a | 3 Titanium local NVMe RAID0 | 100 GiB | 35.2s median | 24.38 Gbps median | Confirms the output path can reach about 25 Gbps when local disk is not the limiter |

The accepted result used 64 MiB S3 ranges, adaptive start target 150, peak running parts 162, and zero part failures. Coalescing short response-body reads before `WriteAt` reduced write syscall pressure while preserving the same temp-file, resume, cancellation, and exact byte-count correctness model. Later tuning should preserve this output design unless benchmarks show a clear improvement without increasing correctness risk for resume, cancellation, or non-file writer compatibility paths.

Agent-proxy download URLs use the same V2 preallocated temp-file output path when adaptive download concurrency is enabled. Agent downloads use the V2 adaptive controller and V2 agent part-size tiers instead of the legacy 5 MiB part size and fixed 15-worker pool. The agent-proxy path currently keeps the legacy ranged downloader's 15-part starting point because VM benchmarks showed that starting higher produced more proxy streams but lower throughput on the tested agent path. The Download V2 completion log includes adaptive target, peak running, throughput-probe, backpressure, latency, and queue-estimate fields so future agent/proxy tuning can be based on observed controller behavior instead of only elapsed time.

The exact benchmark setup, VM commands, metrics collection, and validation checks are documented in [Adaptive Download V2 Benchmark Reproduction](adaptive-download-v2-benchmark-reproduction.md).

Global HTTP/2 validation:

| Path | VM | Disk | Object Size | Elapsed | Throughput | Decision |
| --- | --- | --- | ---: | ---: | ---: | --- |
| Coalesced V2 download with default pooled transport attempting HTTP/2 | n2-highcpu-32, us-east4-a | 4 local NVMe RAID0 | 100 GiB | 70.8s | 12.13 Gbps | No single-object S3 download improvement over accepted tuning |
| Coalesced V2 download with HTTP/2 default, S3 upload requests forced to HTTP/1.1 | n2-highcpu-32, us-east4-a | 4 local NVMe RAID0 | 100 GiB | 66.7s | 12.87 Gbps | Preserves accepted tuning performance while allowing HTTP/2 for Files.com API calls |

The HTTP/2-enabled runs kept the S3 download path above 10 Gbps and had zero part failures. They did not improve this single native/S3 object benchmark because the data transfer URL negotiated HTTP/1.1 with S3; the primary expected benefit is for Files.com API calls on hosts that negotiate HTTP/2.

Connection metrics must remain available when the pooled transport splits S3 upload requests onto an HTTP/1.1-only transport. The S3 upload transport is cloned from the main pooled transport and uses the same stats-wrapped `DialContext`, so CLI `--connection-metrics` still sees both API and data connections through `GetConnectionStatsFromClient`. HTTP/2 API requests are counted as their underlying TCP connections, not individual HTTP/2 streams, matching the existing "Open Connections" metric contract.
