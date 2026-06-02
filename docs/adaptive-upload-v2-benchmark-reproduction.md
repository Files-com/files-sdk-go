# Adaptive Upload V2 Benchmark Reproduction

This document describes the benchmark shape used to compare static upload concurrency against adaptive upload V2. It intentionally avoids environment-specific connection details such as VM names, IP addresses, SSH keys, cloud project names, bucket names, and local profile names.

## Environment

- Use a Linux host with enough CPU and network capacity to exercise S3 upload throughput. A small multi-core cloud VM is sufficient for default validation; larger instances may be useful when testing higher network ceilings.
- Install the current CLI build under test and a baseline CLI build from the target branch or release being compared.
- Configure normal Files.com CLI authentication with access to a destination that begins normal direct-S3 multipart uploads.
- Use an empty remote test directory or prefix for each benchmark run.
- Enable CLI debug logging so upload V1 and V2 scheduler summaries, adaptive telemetry, HTTP client limits, and retry observations are captured.

## Test Data

Create fixed-size binary inputs on local disk. The same source files should be reused for baseline and adaptive runs.

Recommended datasets:

| Dataset | Shape | Purpose |
| --- | --- | --- |
| `20x200` | 20 files, 200 MiB each | Short multi-file ramp and scheduling overhead |
| `200x200` | 200 files, 200 MiB each | Sustained multi-file throughput and adaptive stability |
| `1x20GiB` | 1 file, 20 GiB | Single large-file ramp and per-file part scheduling |

Example data generation pattern:

```sh
mkdir -p ./upload-sources/20x200 ./upload-sources/200x200 ./upload-sources/1x20GiB

for i in $(seq -w 1 20); do
  dd if=/dev/urandom of="./upload-sources/20x200/file-${i}.bin" bs=1M count=200 status=none
done

for i in $(seq -w 1 200); do
  dd if=/dev/urandom of="./upload-sources/200x200/file-${i}.bin" bs=1M count=200 status=none
done

dd if=/dev/urandom of="./upload-sources/1x20GiB/file-001.bin" bs=1M count=20480 status=none
```

## Metrics Collection

Capture at least these metrics per run:

- Wall-clock elapsed time.
- Effective throughput in MiB/s and Gbps, calculated from total source bytes and elapsed time.
- CLI process CPU percent sampled during the upload.
- CLI process resident memory sampled during the upload.
- Network transmit and receive byte deltas from the host network interface.
- Maximum observed HTTPS connection count for the CLI process.
- CLI debug logs containing upload V1/V2 scheduler summaries and adaptive telemetry.

A generic sampler can poll once per second while the CLI process is running:

```sh
while kill -0 "${CLI_PID}" 2>/dev/null; do
  date +%s
  ps -p "${CLI_PID}" -o %cpu=,rss=
  ss -tanp 2>/dev/null | grep -c "${CLI_PROCESS_NAME}" || true
  cat "/sys/class/net/${NET_IFACE}/statistics/tx_bytes"
  cat "/sys/class/net/${NET_IFACE}/statistics/rx_bytes"
  sleep 1
done
```

Use an equivalent command on hosts that do not provide `ss` or `/sys/class/net`.

## Baseline Runs

Run static uploads without adaptive upload V2. Test representative concurrency values so the adaptive default is compared against the best static result, not only against the existing default.

Recommended static limits:

- `50`
- `100`
- `150`
- `200`

Generic command shape:

```sh
files-cli upload \
  --debug \
  --max-concurrent-connections "${STATIC_LIMIT}" \
  "./upload-sources/${DATASET}" \
  "remote/path/${RUN_ID}/static-${STATIC_LIMIT}/${DATASET}"
```

Record one row per dataset and static limit.

## Adaptive Runs

Run the same datasets with adaptive upload V2 enabled. Do not pass exact tuning values for the primary default comparison. User-provided concurrency, when present, should be treated only as a maximum cap.

Generic command shape:

```sh
files-cli upload \
  --debug \
  --adaptive-concurrency \
  "./upload-sources/${DATASET}" \
  "remote/path/${RUN_ID}/adaptive-default/${DATASET}"
```

Diagnostic tuning flags may be used for development-only exploration, but those runs should be labeled separately and should not replace the default adaptive comparison.

## Comparison

For each dataset, compare adaptive default against the best static elapsed time observed for that same dataset.

Calculate:

```text
time_delta_percent = ((adaptive_elapsed_seconds - best_static_elapsed_seconds) / best_static_elapsed_seconds) * 100
throughput_delta_percent = ((adaptive_gbps - best_static_gbps) / best_static_gbps) * 100
```

The target result for this tuning pass is adaptive default at least matching, or landing within `5%` slower than, the best static result for:

- `20x200`
- `200x200`
- `1x20GiB`

## Validation Checks

Review the debug log for each adaptive run and confirm:

- `upload v2 enabled` appears for the uploaded files.
- The target class is `s3` for direct-S3 validation.
- `adaptive_max_target` reflects the V2 adaptive headroom.
- `adaptive_growth_ceiling` reflects the default S3 soft ceiling.
- `adaptive_growth_ceiling_unlocked` is only true when the workload and throughput criteria justify probing above the soft ceiling.
- `upload_max_conns_per_host` is raised above the SDK default and matches the expected transport cap for the workload.
- Retry, back-pressure, and throughput-probe counters are present in completion telemetry.

Also confirm all test VMs or benchmark hosts are stopped or torn down when the run is finished if they are not needed for ongoing work.
