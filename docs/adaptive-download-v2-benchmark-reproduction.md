# Adaptive Download V2 Benchmark Reproduction

This document describes the benchmark shape used to validate adaptive download V2 for large native object-storage downloads. It is intended to make the performance claims reproducible without relying on local shell history.

The benchmark downloads one known-size, S3-native 100 GiB file through the generated CLI with adaptive download concurrency enabled. It records end-to-end CLI throughput, process CPU/RSS, system CPU/memory, NIC receive throughput, disk write throughput, and V2 debug logs.

## Build

Generate Go and CLI targets, then build the Linux benchmark binary:

```bash
/Users/dustin/.rbenv/bin/rbenv exec bundle exec ruby files-sdk-generator.rb \
  --file /Users/dustin/.codex/worktrees/7e86/files-sdk-generator/swagger_doc.json \
  --target go,cli

cd generated/cli
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -o /tmp/files-cli-dev-download-v2 .
```

## Benchmark Object

Use a 100 GiB file in native S3-backed storage. The accepted runs used:

```text
/tmp/adaptive-download-bench/20260605T013723Z-s3-native/single-100GiB
```

If the object must be seeded again, create a 100 GiB file on the benchmark VM local SSD and upload it to an empty S3-native remote prefix before running download benchmarks. Seeding may use adaptive upload V2, but seed upload throughput is not part of the download benchmark result.

## VM Shapes

Primary accepted benchmark:

```bash
gcloud compute instances create files-cli-download-east-bench \
  --project=expanded-curve-341004 \
  --zone=us-east4-a \
  --machine-type=n2-highcpu-32 \
  --network-performance-configs=total-egress-bandwidth-tier=TIER_1 \
  --network-interface=stack-type=IPV4_ONLY,nic-type=GVNIC \
  --image-project=ubuntu-os-cloud \
  --image-family=ubuntu-2404-lts-amd64 \
  --boot-disk-size=50GB \
  --boot-disk-type=pd-balanced \
  --local-ssd=interface=NVME \
  --local-ssd=interface=NVME \
  --local-ssd=interface=NVME \
  --local-ssd=interface=NVME \
  --maintenance-policy=TERMINATE \
  --restart-on-failure \
  --provisioning-model=STANDARD
```

Faster local SSD validation:

```bash
gcloud compute instances create files-cli-download-z3-bench \
  --project=expanded-curve-341004 \
  --zone=us-east4-a \
  --machine-type=z3-highmem-22-highlssd \
  --network-interface=stack-type=IPV4_ONLY,nic-type=GVNIC \
  --image-project=ubuntu-os-cloud \
  --image-family=ubuntu-2404-lts-amd64 \
  --boot-disk-size=50GB \
  --boot-disk-type=pd-balanced \
  --maintenance-policy=MIGRATE \
  --restart-on-failure \
  --provisioning-model=STANDARD
```

`z3-highmem-32-highlssd` was the preferred Titanium local SSD target, but the benchmark project quota only allowed 24 Z3-family vCPUs in `us-east4`.

## VM Setup

Copy the benchmark binary and CLI profile without printing API keys:

```bash
gcloud compute scp /tmp/files-cli-dev-download-v2 VM_NAME:/tmp/files-cli-dev \
  --project=expanded-curve-341004 \
  --zone=us-east4-a

gcloud compute ssh VM_NAME \
  --project=expanded-curve-341004 \
  --zone=us-east4-a \
  --command 'mkdir -p ~/.config'

gcloud compute scp --recurse "$HOME/.config/files-cli" VM_NAME:~/.config/files-cli \
  --project=expanded-curve-341004 \
  --zone=us-east4-a
```

Install metric tools, apply the TCP settings, and create RAID0 across all local NVMe devices. For `n2-highcpu-32`, require four local NVMe devices. For the z3 validation, use all non-boot local NVMe devices exposed by the machine shape.

```bash
sudo apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y mdadm sysstat jq moreutils

sudo modprobe tcp_bbr || true
sudo tee /etc/sysctl.d/99-files-bench.conf >/dev/null <<'SYSCTL'
net.core.default_qdisc=fq
net.ipv4.tcp_congestion_control=bbr
net.ipv4.tcp_slow_start_after_idle=0
net.ipv4.tcp_rmem=4096 87380 134217728
net.ipv4.tcp_wmem=4096 65536 134217728
net.core.rmem_max=134217728
net.core.wmem_max=134217728
net.core.netdev_max_backlog=250000
SYSCTL
sudo sysctl --system >/tmp/sysctl-apply.log

mapfile -t disks < <(ls /dev/disk/by-id/google-local-nvme-ssd-* 2>/dev/null | sort)
if [[ ${#disks[@]} -eq 0 ]]; then
  mapfile -t disks < <(lsblk -dnpo NAME,TYPE,TRAN,MODEL,MOUNTPOINT | awk '$2 == "disk" && $3 == "nvme" && $4 != "nvme_card-pd" && $5 == "" {print $1}' | sort)
fi
if [[ ${#disks[@]} -eq 0 ]]; then
  echo "expected local NVMe disks" >&2
  lsblk -o NAME,TYPE,SIZE,MODEL,TRAN,MOUNTPOINT >&2
  exit 1
fi

raid_devices="${#disks[@]}"
stripe_width=$((128 * raid_devices))
sudo mdadm --create /dev/md0 --level=0 --raid-devices="$raid_devices" --chunk=512K "${disks[@]}" --force
sudo mkfs.ext4 -F -E stride=128,stripe-width="$stripe_width" /dev/md0
sudo mkdir -p /mnt/bench
sudo mount -o noatime,nodiratime /dev/md0 /mnt/bench
sudo chmod 0777 /mnt/bench

mkdir -p /mnt/bench/bin /mnt/bench/logs /mnt/bench/metrics /mnt/bench/results /mnt/bench/downloads
cp /tmp/files-cli-dev /mnt/bench/bin/files-cli-dev
chmod +x /mnt/bench/bin/files-cli-dev

{
  echo "date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "kernel=$(uname -a)"
  echo "machine=$(curl -sf -H Metadata-Flavor:Google http://metadata.google.internal/computeMetadata/v1/instance/machine-type || true)"
  echo "zone=$(curl -sf -H Metadata-Flavor:Google http://metadata.google.internal/computeMetadata/v1/instance/zone || true)"
  echo "interface=$(ip -o link show | awk -F': ' '$2 !~ /^lo/ {print $2; exit}')"
  echo "local_nvme_count=$raid_devices"
  echo "local_nvme_devices=${disks[*]}"
  sysctl net.core.default_qdisc net.ipv4.tcp_congestion_control net.ipv4.tcp_slow_start_after_idle net.ipv4.tcp_rmem net.ipv4.tcp_wmem net.core.rmem_max net.core.wmem_max net.core.netdev_max_backlog
  lsblk -o NAME,TYPE,SIZE,MODEL,TRAN,MOUNTPOINT
  findmnt /mnt/bench
} > /mnt/bench/results/environment.txt
```

## Run

Run each benchmark with a unique label:

```bash
label="east-n2-32-coalesced"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-${label}"
remote="/tmp/adaptive-download-bench/20260605T013723Z-s3-native/single-100GiB"
dest="/mnt/bench/downloads/${run_id}"
debug_log="/mnt/bench/logs/${run_id}.debug.log"
output_log="/mnt/bench/results/${run_id}.output.txt"
summary="/mnt/bench/results/${run_id}.summary.tsv"
pidstat_log="/mnt/bench/metrics/${run_id}.pidstat.log"
vmstat_log="/mnt/bench/metrics/${run_id}.vmstat.log"
sar_log="/mnt/bench/metrics/${run_id}.sar.log"
iostat_log="/mnt/bench/metrics/${run_id}.iostat.log"

mkdir -p "$dest" /mnt/bench/logs /mnt/bench/results /mnt/bench/metrics
iface="$(ip -o link show | awk -F': ' '$2 !~ /^lo/ {print $2; exit}')"

vmstat 1 > "$vmstat_log" &
vmstat_pid=$!
sar -n DEV 1 > "$sar_log" &
sar_pid=$!
iostat -dxm 1 > "$iostat_log" &
iostat_pid=$!

start_ns="$(date +%s%N)"
set +e
/mnt/bench/bin/files-cli-dev download \
  --adaptive-concurrency \
  --connection-metrics \
  --ignore-version-check \
  --debug="$debug_log" \
  "$remote" \
  "$dest" > "$output_log" 2>&1 &
cli_pid=$!
pidstat -h -r -u -p "$cli_pid" 1 > "$pidstat_log" &
pidstat_pid=$!
wait "$cli_pid"
exit_code=$?
set -e
end_ns="$(date +%s%N)"

kill "$pidstat_pid" "$vmstat_pid" "$sar_pid" "$iostat_pid" 2>/dev/null || true
wait "$pidstat_pid" "$vmstat_pid" "$sar_pid" "$iostat_pid" 2>/dev/null || true

elapsed_s="$(awk -v s="$start_ns" -v e="$end_ns" 'BEGIN { printf "%.3f", (e-s)/1000000000 }')"
bytes="$(find "$dest" -type f -printf '%s\n' | awk '{sum+=$1} END {printf "%.0f", sum}')"
gbps="$(awk -v b="$bytes" -v s="$elapsed_s" 'BEGIN { if (s > 0) printf "%.3f", (b*8)/(s*1000000000); else print "0.000" }')"
nic_avg_peak="$(awk -v iface="$iface" '$0 ~ iface && $0 !~ /Average/ { rx=$5+0; if(rx>0){sum+=rx; n++; if(rx>max)max=rx}} END{ if(n>0) printf "%.3f\t%.3f", sum/n*8/1000000, max*8/1000000; else printf "0.000\t0.000" }' "$sar_log")"
cpu_rss="$(awk '/files-cli-dev/ && $0 !~ /Command/ { usr+=$4; sys+=$5; cpu+=$8; n++; if($8>pcpu)pcpu=$8; if($12>rss)rss=$12 } END{ if(n>0) printf "%.2f\t%.2f\t%.2f\t%.2f\t%.0f", cpu/n, pcpu, usr/n, sys/n, rss; else printf "0.00\t0.00\t0.00\t0.00\t0" }' "$pidstat_log")"
disk_stats="$(awk '$1=="md0" && $9+0>0 {sum+=$9; n++; if($9>max)max=$9; util=$NF+0; if(util>maxu)maxu=util} END{ if(n>0) printf "%.2f\t%.2f\t%.2f", sum/n, max, maxu; else printf "0.00\t0.00\t0.00" }' "$iostat_log")"

{
  printf "run_id\t%s\n" "$run_id"
  printf "exit_code\t%s\n" "$exit_code"
  printf "elapsed_s\t%s\n" "$elapsed_s"
  printf "bytes\t%s\n" "$bytes"
  printf "gbps\t%s\n" "$gbps"
  printf "iface\t%s\n" "$iface"
  printf "ens_rx_avg_gbps\t%s\n" "$(cut -f1 <<<"$nic_avg_peak")"
  printf "ens_rx_peak_gbps\t%s\n" "$(cut -f2 <<<"$nic_avg_peak")"
  printf "cpu_avg_percent\t%s\n" "$(cut -f1 <<<"$cpu_rss")"
  printf "cpu_peak_percent\t%s\n" "$(cut -f2 <<<"$cpu_rss")"
  printf "usr_avg_percent\t%s\n" "$(cut -f3 <<<"$cpu_rss")"
  printf "sys_avg_percent\t%s\n" "$(cut -f4 <<<"$cpu_rss")"
  printf "rss_peak_kb\t%s\n" "$(cut -f5 <<<"$cpu_rss")"
  printf "md0_write_avg_MBps\t%s\n" "$(cut -f1 <<<"$disk_stats")"
  printf "md0_write_peak_MBps\t%s\n" "$(cut -f2 <<<"$disk_stats")"
  printf "md0_util_peak\t%s\n" "$(cut -f3 <<<"$disk_stats")"
} > "$summary"

cat "$summary"
exit "$exit_code"
```

## Validate

The debug log must show:

- `download_v2_enabled: true`
- `download_v2_target: s3`
- `download_v2_output_mode: preallocated_temp_file_write_at`
- `download_v2_part_size: 67108864`
- `download_v2_adaptive_start: 150`
- `download_v2_adaptive_peak_running` around `162` for the accepted 100 GiB runs
- `download_v2_adaptive_failure_total: 0`
- `download_v2_contiguous_size: 107374182400`

Example check:

```bash
rg "download v2 start|download v2 finish" "$debug_log"
```

Archive results before deleting or reusing the VM:

```bash
tar -czf "/tmp/${label}-download-results.tgz" -C /mnt/bench logs metrics results
```

Copy the archive back locally with `gcloud compute scp`.

## Accepted Results

| Scenario | VM / storage | Elapsed | Throughput | Notes |
| --- | --- | ---: | ---: | --- |
| Coalesced V2 download | `n2-highcpu-32`, Tier 1 `gVNIC`, 4 local NVMe RAID0 | 66.134s | 12.989 Gbps | Accepted tuning for this MR |
| HTTP/2 default, S3 uploads forced HTTP/1.1 | `n2-highcpu-32`, Tier 1 `gVNIC`, 4 local NVMe RAID0 | 66.730s | 12.873 Gbps | Preserved accepted S3 download performance |
| Faster local SSD validation | `z3-highmem-22-highlssd`, `gVNIC`, 3 Titanium local NVMe RAID0 | 35.239s median | 24.376 Gbps median | Shows the output path can reach about 25 Gbps when local disk is not the limiter |

## Cleanup

Delete or stop benchmark VMs when they are no longer needed. Do not leave expensive benchmark VMs running after a completed benchmark unless they are being reused for follow-up work.
