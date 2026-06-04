package ostuning

import "fmt"

func darwinPlan(includeNetworkTest bool) Plan {
	snapshotPath := "/var/db/files.com/os-tuning/high-throughput-upload.snapshot"
	var networkTests []Step
	if includeNetworkTest {
		networkTests = []Step{darwinNetworkQualityStep()}
	}

	return Plan{
		OS:           "darwin",
		Profile:      ProfileHighThroughputUpload,
		SnapshotPath: snapshotPath,
		Summary:      "Inspect and optionally raise macOS TCP socket buffer ceilings for high-throughput uploads. macOS has fewer safe global knobs than Linux.",
		UserSteps: []Step{
			{
				ID:          "darwin.inspect-tcp",
				Title:       "Inspect current TCP buffer settings",
				Description: "Read macOS socket buffer and TCP auto-buffer ceilings.",
				Privilege:   PrivilegeUser,
				Commands: []Command{
					posix("sysctl kern.ipc.maxsockbuf net.inet.tcp.sendspace net.inet.tcp.recvspace net.inet.tcp.autorcvbufmax net.inet.tcp.autosndbufmax"),
				},
				CanFailSoftly: true,
			},
			{
				ID:          "darwin.inspect-nofile",
				Title:       "Inspect open file descriptor limits",
				Description: "Read the current shell soft and hard nofile limits. Adaptive uploads can open many sockets and files, so a low soft limit can cap the process before TCP is the bottleneck.",
				Privilege:   PrivilegeUser,
				Commands: []Command{
					posix(openFileLimitInspectCommand()),
				},
				ExpectedOutcome: fmt.Sprintf("Soft nofile should be at least %d, with %d or higher preferred on dedicated high-throughput hosts.", MinimumOpenFileLimit, PreferredOpenFileLimit),
			},
		},
		NetworkTests: networkTests,
		SnapshotSteps: []Step{
			darwinSnapshotStep(snapshotPath),
		},
		AdminSteps: []Step{
			{
				ID:          "darwin.runtime-buffers",
				Title:       "Raise runtime TCP buffer ceilings",
				Description: "Raise socket and TCP auto-buffer ceilings for the current boot. Treat this as a benchmark-only change before deploying through MDM or a launch daemon.",
				Privilege:   PrivilegeAdministrator,
				Commands: []Command{
					posix("sysctl -w kern.ipc.maxsockbuf=16777216"),
					posix("sysctl -w net.inet.tcp.autorcvbufmax=16777216"),
					posix("sysctl -w net.inet.tcp.autosndbufmax=16777216"),
					posix("sysctl -w net.inet.tcp.sendspace=1048576"),
					posix("sysctl -w net.inet.tcp.recvspace=1048576"),
				},
				Verification: []Command{
					posix("sysctl kern.ipc.maxsockbuf net.inet.tcp.sendspace net.inet.tcp.recvspace net.inet.tcp.autorcvbufmax net.inet.tcp.autosndbufmax"),
				},
				RuntimeOnly:     true,
				CanFailSoftly:   true,
				ExpectedOutcome: "Supported macOS releases report larger TCP buffer ceilings for new sockets during the current boot.",
			},
		},
		RestoreSteps: []Step{
			darwinRestoreStep(snapshotPath),
		},
		Warnings: []string{
			"macOS TCP sysctls vary by release. Unsupported keys should be treated as a no-op rather than a deployment blocker.",
			"These changes are runtime-only in this plan. Prefer MDM or an audited launch daemon for persistence after validating a specific macOS release.",
			"Do not change unrelated keepalive, delayed ACK, MTU, or private sysctl values as a generic upload optimization.",
		},
		Notes: []string{
			"Expected gains are usually smaller than Linux because modern macOS already auto-tunes TCP windows and does not expose BBR as a generic sysctl.",
			"Run a before/after upload benchmark on the same network before recommending this to desktop users.",
			"Adaptive CLI uploads raise their own soft nofile limit to the preferred high-throughput limit at startup when macOS allows it. Persistent maxfiles changes should be deployed through MDM or an audited launch daemon after validating the target macOS release.",
		},
		References: []Reference{
			{Title: "Apple guidance on querying system features with sysctl", URL: "https://developer.apple.com/documentation/Apple-Silicon/addressing-architectural-differences-in-your-macos-code"},
			{Title: "Apple networking guidance warning about host-wide sysctl impact", URL: "https://developer.apple.com/library/archive/documentation/NetworkingInternetWeb/Conceptual/NetworkingOverview/CommonPitfalls/CommonPitfalls.html"},
		},
	}
}

func darwinNetworkQualityStep() Step {
	return Step{
		ID:          "darwin.measure-network",
		Title:       "Measure Apple network quality",
		Description: "Use Apple's networkQuality tool to record general Internet bandwidth and responsiveness against Apple's selected test endpoint. This is not a Files.com, S3, or agent upload benchmark.",
		Privilege:   PrivilegeUser,
		Commands: []Command{
			posix(`if command -v networkQuality >/dev/null 2>&1; then
  networkQuality -v
else
  printf 'networkQuality is not available on this macOS host; skipping active network measurement.\n'
fi`),
		},
		CanFailSoftly:   true,
		ExpectedOutcome: "When networkQuality is available, record Apple endpoint capacity and responsiveness for comparison.",
	}
}

func darwinSnapshotStep(snapshotPath string) Step {
	return Step{
		ID:          "darwin.snapshot",
		Title:       "Snapshot current high-throughput tuning values",
		Description: "Store the current Files.com-managed macOS TCP buffer values before applying runtime repair changes.",
		Privilege:   PrivilegeAdministrator,
		Commands: []Command{
			posix(fmt.Sprintf(`SNAPSHOT=%q
mkdir -p "$(dirname "$SNAPSHOT")"
{
  printf 'version=1\n'
  printf 'os=darwin\n'
  for key in kern.ipc.maxsockbuf net.inet.tcp.sendspace net.inet.tcp.recvspace net.inet.tcp.autorcvbufmax net.inet.tcp.autosndbufmax; do
    value=$(sysctl -n "$key" 2>/dev/null || true)
    if [ -n "$value" ]; then
      printf 'sysctl.%%s=%%s\n' "$key" "$value"
    fi
  done
} > "$SNAPSHOT"
chmod 600 "$SNAPSHOT"
printf 'Saved Files.com high-throughput OS tuning snapshot to %%s\n' "$SNAPSHOT"`, snapshotPath)),
		},
		ExpectedOutcome: "A pre-change snapshot exists for exact runtime restore when supported by the macOS release.",
	}
}

func darwinRestoreStep(snapshotPath string) Step {
	return Step{
		ID:          "darwin.restore",
		Title:       "Restore high-throughput tuning from snapshot or defaults",
		Description: "Restore saved macOS TCP buffer values when a snapshot exists; otherwise report the safest best-effort default restore path.",
		Privilege:   PrivilegeAdministrator,
		Commands: []Command{
			posix(fmt.Sprintf(`set -e
SNAPSHOT=%q
if [ -f "$SNAPSHOT" ]; then
  while IFS='=' read -r key value; do
    case "$key" in
      sysctl.kern.ipc.maxsockbuf|sysctl.net.inet.tcp.sendspace|sysctl.net.inet.tcp.recvspace|sysctl.net.inet.tcp.autorcvbufmax|sysctl.net.inet.tcp.autosndbufmax)
        sysctl -w "${key#sysctl.}=$value"
        ;;
    esac
  done < "$SNAPSHOT"
  printf 'Restored Files.com high-throughput OS tuning snapshot from %%s\n' "$SNAPSHOT"
else
  printf 'No Files.com high-throughput OS tuning snapshot found at %%s.\n' "$SNAPSHOT"
  printf 'This macOS plan only applies runtime sysctl changes, so rebooting is the safest best-effort default restore when no snapshot exists.\n'
fi`, snapshotPath)),
		},
		Verification: []Command{
			posix("sysctl kern.ipc.maxsockbuf net.inet.tcp.sendspace net.inet.tcp.recvspace net.inet.tcp.autorcvbufmax net.inet.tcp.autosndbufmax"),
		},
		ExpectedOutcome: "Host values return to the captured snapshot when available; otherwise the user receives the safest runtime-default restore guidance.",
	}
}
