# Adaptive Upload V2 File-Admission Cap Decision

## Decision

Adaptive upload V2 defaults the multi-file job admission cap to `128` files:

```go
manager.AdaptiveUploadV2ConcurrentFiles = 128
```

The adaptive part-concurrency cap remains `1024`. The part controller still owns HTTP part concurrency and settles at the measured active target. The file-admission cap only controls how many files in a multi-file upload job may be active enough to feed the part scheduler.

## Why This Exists

Adaptive upload has two different concurrency concerns:

- Part concurrency: how many HTTP upload parts may be in flight across the job.
- File admission: how many files may be active at the same time to keep enough parts available.

If file admission is too low, many-file jobs can starve the adaptive part scheduler, especially when each file has only a small number of parts. If file admission is too high, the CLI/SDK creates extra file-level work, API pressure, memory churn, and scheduling overhead without increasing part throughput.

The cap therefore needs to be high enough to feed direct-S3 uploads on a fast network, but not so high that it recreates the static V1 behavior of admitting far more work than the adaptive scheduler can use.

## Benchmark Environment

The benchmark compared only the file-admission cap while keeping adaptive part concurrency unchanged.

VM:

- GCP project: `expanded-curve-341004`
- Zone: `us-west2-c`
- VM name: `files-cli-tier1-upload-bench`
- Actual machine type: `n2-standard-32`
- Network: Tier_1 with GVNIC
- OS: Ubuntu 24.04 LTS
- CPU: 32 vCPU

The original preferred `c3d-standard-30-lssd` was unavailable in `us-west2-c`, and `c3-standard-44` exceeded the project global CPU quota. `c3-standard-22` was rejected because it did not have enough vCPUs for Tier_1. `n2-standard-32` was the largest practical Tier_1/GVNIC VM within the project quota.

Network note:

- Tier_1 documents high internal egress for eligible VM sizes, but direct-to-S3 uses external egress. Results were therefore evaluated from observed upload throughput, not from the internal Tier_1 headline number.

## OS Tuning

The first matrix was run on default Ubuntu settings. Those defaults were not suitable for a high-throughput upload benchmark:

- Soft file descriptor limit: `1024`
- TCP congestion control: `cubic`
- Queueing discipline: `fq_codel`
- `tcp_slow_start_after_idle`: `1`
- Socket buffer maximums: `212992`
- Ephemeral port range: `32768 60999`
- NIC RX/TX rings: `1024`

Before the final decision matrix, the VM was tuned with:

```sh
files-cli-dev os-tuning high-throughput repair --apply --interface ens4
```

The benchmark shell also raised its soft file descriptor limit to `1048576`.
That nofile dependency is now part of the committed tuning model: adaptive CLI
uploads attempt to raise their own soft nofile limit to the preferred
high-throughput limit at startup, and the high-throughput OS tuning plan
inspects the current limit and can persist a Linux PAM limits file for new
sessions.

Verified tuned state:

- TCP congestion control: `bbr`
- Queueing discipline: `fq`
- `tcp_slow_start_after_idle`: `0`
- `net.core.rmem_max`: `134217728`
- `net.core.wmem_max`: `134217728`
- `net.ipv4.tcp_rmem`: `4096 1048576 134217728`
- `net.ipv4.tcp_wmem`: `4096 1048576 134217728`
- Ephemeral port range: `10000 65535`
- `tcp_mtu_probing`: `1`
- `somaxconn`: `65535`
- NIC RX/TX rings: `2048`
- Soft `nofile` limit: `1048576` for the benchmark shell and upload process

## Test Matrix

Direct-to-S3 uploads only.

Datasets:

| Dataset | File count | File size | Approx total |
|---|---:|---:|---:|
| `1000x64MiB` | 1000 | 64 MiB | 64 GiB |
| `200x200MiB` | 200 | 200 MiB | 39 GiB |

File-admission caps:

- `300`
- `128`
- `50`

Each cap was run three times per dataset in interleaved order:

```text
run 1: 300, 128, 50
run 2: 128, 50, 300
run 3: 50, 300, 128
```

The hidden benchmark flag was used so the same binary could vary only file admission:

```sh
files-cli-dev upload \
  --adaptive-concurrency \
  --adaptive-upload-v2-file-concurrency=<cap> \
  --connection-metrics \
  --ignore-version-check \
  --debug=<debug-log> \
  <dataset> \
  <remote-path>
```

## Tuned Results

Median throughput after high-throughput OS tuning:

| Dataset | File cap | Median throughput | Best throughput | Median elapsed | Avg CPU | Max RSS | Result vs `300` |
|---|---:|---:|---:|---:|---:|---:|---:|
| `1000x64MiB` | `300` | 20.458 Gbit/s | 21.387 Gbit/s | 26.2 s | 179.0% | 105.2 MiB | baseline |
| `1000x64MiB` | `128` | 21.683 Gbit/s | 22.042 Gbit/s | 24.8 s | 207.6% | 88.7 MiB | +6.0% |
| `1000x64MiB` | `50` | 17.713 Gbit/s | 17.798 Gbit/s | 30.3 s | 163.2% | 82.5 MiB | -13.4% |
| `200x200MiB` | `300` | 20.254 Gbit/s | 20.338 Gbit/s | 16.6 s | 199.5% | 94.9 MiB | baseline |
| `200x200MiB` | `128` | 20.622 Gbit/s | 20.635 Gbit/s | 16.3 s | 185.3% | 89.4 MiB | +1.8% |
| `200x200MiB` | `50` | 19.530 Gbit/s | 20.383 Gbit/s | 17.2 s | 184.2% | 82.7 MiB | -3.6% |

The adaptive part controller reported `adaptive_final_target=150`, `adaptive_peak_target=150`, and `adaptive_peak_running=150` for all runs. That confirms the benchmark isolated file admission while holding the adaptive part target behavior constant.

## Interpretation

`128` is the best default from this run:

- It beat `300` on both tuned datasets by median throughput.
- It used less memory than `300`.
- It preserved enough admitted files to feed the adaptive part scheduler for the `1000x64MiB` workload.
- It avoided the underfeeding seen at `50` for many 64 MiB files.

`50` can be acceptable for larger files, because each active file contributes more parts, but it is too low for many-file workloads with small-to-medium files. `300` can feed the scheduler, but after OS tuning it did not produce better median throughput and used more memory.

## Follow-Up

Keep the hidden `--adaptive-upload-v2-file-concurrency` flag for benchmark and diagnostic runs. It should not be exposed as normal CLI UX unless future data shows users need a public cap separate from `--concurrent-connection-limit`.

If future workloads show a recurring split between small-file and large-file behavior, the next step should be dynamic file admission based on observed planned parts per active file and current adaptive part target. Until then, `128` is the measured default.
