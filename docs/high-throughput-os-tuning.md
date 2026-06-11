# High-throughput OS tuning

The Go SDK exposes high-throughput host tuning as a structured plan in
`github.com/Files-com/files-sdk-go/v3/lib/ostuning`. The CLI command is only a
thin renderer for that SDK package:

```sh
files-cli os-tuning high-throughput plan
files-cli os-tuning high-throughput verify
files-cli os-tuning high-throughput verify --os darwin --include-network-test
files-cli os-tuning high-throughput repair --os linux --interface ens4
files-cli os-tuning high-throughput repair --apply
files-cli os-tuning high-throughput restore
files-cli os-tuning high-throughput restore --apply
files-cli os-tuning high-throughput repair --os windows --interface "Ethernet 2" --commands-only
```

The plan separates user-level inspection commands from privileged changes. This
is intentional so CLI and SDK callers can present the same recommendations with
their own approval and elevation workflow.

## Workflow

Use the same flow that host-management tools commonly use:

1. `verify`: runs non-mutating checks on the current host.
2. `plan`: dry-run view of the full workflow, including user-level checks,
   privileged changes, warnings, notes, and references.
3. `repair`: dry-run view of the snapshot and privileged changes. With
   `--apply`, the CLI first stores a pre-change snapshot when it has enough
   privilege, then executes every repair step allowed by the current process
   privileges. If some steps need root or Administrator rights, they are skipped,
   listed, and paired with the exact elevated follow-up command. For Linux and
   macOS, re-run that command with `sudo`. For Windows, re-run it from an
   Administrator PowerShell session.
4. `restore`: dry-run view of the rollback workflow. With `--apply`, the CLI
   restores the saved snapshot when one exists. If no snapshot is available, it
   performs the safest best-effort default restore for the OS and clearly reports
   what could not be inferred.
5. `--include-network-test`: opt-in active bandwidth measurement where the OS
   plan supports it. `verify --include-network-test` records the current state.
   `repair --apply --include-network-test` records a before measurement, applies
   the repair steps allowed by the current process privileges, then records an
   after measurement only when privileged repair commands actually ran. This is
   off by default so `verify` remains a lightweight inspection workflow.
6. `--commands-only`: script-oriented output for administrators, deployment
   tooling, MDM, docs, or rendering checks for another OS.

The SDK mirrors this shape through `VerificationSteps`, `RepairSteps`,
`RestorePlanSteps`, `RunSteps`, `RunStepsWithElevation`, and `RunCommands`. The
command strings do not embed `sudo`; callers decide whether to run elevated now,
skip privileged steps, or hand an elevated follow-up command to a separate
installer flow.

## Snapshots and restore

`repair --apply` writes a pre-change snapshot before making privileged changes
when the command is already running as root or Administrator. The snapshot is a
host-local safety net for values that do not have a universal default, such as
NIC ring sizes or adapter offload state.

Snapshot locations:

- Linux: `/var/lib/files.com/os-tuning/high-throughput-upload.snapshot`
- macOS: `/var/db/files.com/os-tuning/high-throughput-upload.snapshot`
- Windows:
  `%ProgramData%\Files.com\os-tuning\high-throughput-upload-snapshot.json`

`restore --apply` prefers the snapshot. Without a snapshot, Linux removes the
Files.com-managed sysctl/module-load files and reloads the remaining system
configuration. macOS reports that rebooting is the safest default restore
because this plan only applies runtime sysctls. Windows reapplies supported TCP
defaults and reports that adapter RSS/RSC driver defaults require a snapshot.

## Linux

The Linux plan focuses on the settings that produced the largest benchmark
improvement for adaptive uploads:

- BBR congestion control when the kernel supports it.
- Larger TCP send and receive buffer ceilings.
- Wider ephemeral port range for high concurrent connection counts, using
  `10000 65535` to avoid overlapping common registered service ports.
- `tcp_slow_start_after_idle=0` so short pauses do not reset the congestion
  window as aggressively.
- MTU probing.
- Runtime-only NIC RX/TX ring increases when the driver supports them.
- Open file descriptor limit inspection and persistent `nofile` remediation so
  adaptive upload can use high HTTP part concurrency without being capped by a
  low process descriptor limit.

The privileged Linux plan writes `/etc/sysctl.d/99-files-high-throughput.conf`
and `/etc/security/limits.d/99-files-high-throughput-nofile.conf`.

Adaptive CLI uploads also attempt to raise the current process soft `nofile`
limit to the preferred high-throughput limit at transfer startup. That bounded
runtime raise helps current CLI runs, while the persistent limits file helps new
shells and SDK callers that start outside the CLI. PAM limits do not
automatically apply to already-running shells or all systemd services. For
service-managed SDK callers, set an equivalent `LimitNOFILE` in the systemd
service unit or drop-in.

Other Files.com products may apply their own Linux UDP buffer tuning for QUIC.
The high-throughput upload plan uses a later sysctl file and higher shared
socket buffer ceilings. The two can coexist; the later high-throughput file wins
for the shared `rmem_max`/`wmem_max` maxima.

## macOS

The macOS plan is intentionally conservative. It provides inspection commands
for TCP buffer ceilings and current `nofile` limits, plus runtime-only candidate
`sysctl` changes for socket and TCP buffer ceilings. macOS has fewer safe global
knobs than Linux, and supported TCP sysctls vary by release, so unsupported keys
should be treated as non-fatal.

Adaptive CLI uploads attempt to raise the current process soft `nofile` limit to
the preferred high-throughput limit at transfer startup on macOS as well.
Persistent `maxfiles` changes should be deployed through MDM or an audited
launch daemon only after validating the target macOS release.

The optional `--include-network-test` flag adds Apple's `networkQuality` command
when available. This is a general Internet quality test against Apple's selected
test endpoint, not a Files.com or S3 upload benchmark. `verify` runs it
once. `repair --apply` runs it before the host-wide repair changes and runs the
post-repair measurement only when privileged repair commands actually ran. That
active bandwidth test is intentionally not part of the default `verify` path
because it consumes network capacity.

## Windows

The Windows plan focuses on supported Windows TCP and adapter controls:

- Inspect `netsh interface tcp show global`.
- Inspect `Get-NetTCPSetting`, `Get-NetAdapterRss`, and `Get-NetAdapterRsc`.
- Restore receive window auto-tuning to `Normal`.
- Enable global RSS.
- Enable per-adapter RSS and IPv4 RSC where supported.

The plan does not force a congestion provider globally. That should be validated
against the Windows version, NIC driver, and network path before deployment.
IPv6 RSC is intentionally outside this tuning plan.
