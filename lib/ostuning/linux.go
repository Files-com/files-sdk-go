package ostuning

import (
	"fmt"
	"strings"
)

func linuxPlan(interfaceName string) Plan {
	iface := posixSingleQuote(interfaceName)
	if interfaceName == "" {
		iface = "${FILES_HIGH_THROUGHPUT_INTERFACE:-$(ip route show default 0.0.0.0/0 | awk '{print $5; exit}')}"
	}
	snapshotPath := "/var/lib/files.com/os-tuning/high-throughput-upload.snapshot"

	return Plan{
		OS:            "linux",
		Profile:       ProfileHighThroughputUpload,
		InterfaceName: interfaceName,
		SnapshotPath:  snapshotPath,
		Summary:       "Tune Linux TCP and NIC queue defaults for sustained high-throughput uploads over many parallel HTTPS connections.",
		UserSteps: []Step{
			{
				ID:          "linux.inspect-tcp",
				Title:       "Inspect current TCP settings",
				Description: "Read the active congestion control, queueing, socket buffer, port range, and idle slow-start settings.",
				Privilege:   PrivilegeUser,
				Commands: []Command{
					posix("sysctl net.ipv4.tcp_congestion_control net.core.default_qdisc net.ipv4.tcp_slow_start_after_idle net.core.rmem_max net.core.wmem_max net.ipv4.tcp_rmem net.ipv4.tcp_wmem net.ipv4.ip_local_port_range net.ipv4.tcp_mtu_probing"),
				},
			},
			{
				ID:          "linux.inspect-nic",
				Title:       "Inspect NIC driver, channels, and ring buffers",
				Description: "Check whether the active interface exposes multiple queues and larger RX/TX rings.",
				Privilege:   PrivilegeUser,
				Commands: []Command{
					posix(fmt.Sprintf("IFACE=%s; if [ -z \"$IFACE\" ]; then IFACE=$(ip route show default 0.0.0.0/0 | awk '{print $5; exit}'); fi; ip -s link show \"$IFACE\"; ethtool -i \"$IFACE\"; ethtool -l \"$IFACE\"; ethtool -g \"$IFACE\"", iface)),
				},
				CanFailSoftly: true,
			},
		},
		SnapshotSteps: []Step{
			linuxSnapshotStep(iface, snapshotPath),
		},
		AdminSteps: []Step{
			{
				ID:          "linux.persist-tcp",
				Title:       "Persist high-throughput TCP defaults",
				Description: "Enable BBR where available, increase socket buffer ceilings, widen ephemeral ports, keep congestion windows warm after idle periods, and enable MTU probing.",
				Privilege:   PrivilegeAdministrator,
				Commands: []Command{
					posix("modprobe tcp_bbr || true"),
					posix("modprobe sch_fq || true"),
					posix("printf '%s\\n' tcp_bbr sch_fq | tee /etc/modules-load.d/files-high-throughput.conf >/dev/null"),
					posix(`tee /etc/sysctl.d/99-files-high-throughput.conf >/dev/null <<'EOF'
net.ipv4.tcp_congestion_control = bbr
net.core.default_qdisc = fq
net.ipv4.tcp_slow_start_after_idle = 0
net.core.rmem_max = 134217728
net.core.wmem_max = 134217728
net.ipv4.tcp_rmem = 4096 1048576 134217728
net.ipv4.tcp_wmem = 4096 1048576 134217728
net.ipv4.ip_local_port_range = 10000 65535
net.ipv4.tcp_mtu_probing = 1
net.core.somaxconn = 65535
EOF`),
					posix("sysctl --system"),
				},
				Verification: []Command{
					posix("sysctl net.ipv4.tcp_congestion_control net.core.default_qdisc net.ipv4.tcp_slow_start_after_idle net.core.rmem_max net.core.wmem_max net.ipv4.ip_local_port_range net.ipv4.tcp_mtu_probing"),
				},
				ExpectedOutcome: "New TCP connections should use BBR when the kernel supports it, with larger send/receive buffers and a wider local port range.",
			},
			{
				ID:          "linux.nic-rings",
				Title:       "Raise NIC RX/TX ring buffers when supported",
				Description: "Increase hardware queue rings on the active interface. This is runtime-only and may need to be reapplied after reboot or interface recreation.",
				Privilege:   PrivilegeAdministrator,
				Commands: []Command{
					posix(fmt.Sprintf("IFACE=%s; if [ -n \"$IFACE\" ] && command -v ethtool >/dev/null 2>&1; then ethtool -G \"$IFACE\" rx 2048 tx 2048 || true; fi", iface)),
				},
				Verification: []Command{
					posix(fmt.Sprintf("IFACE=%s; if [ -n \"$IFACE\" ] && command -v ethtool >/dev/null 2>&1; then ethtool -g \"$IFACE\"; fi", iface)),
				},
				RuntimeOnly:     true,
				CanFailSoftly:   true,
				ExpectedOutcome: "Interfaces that support 2048-entry rings report RX and TX current values of 2048.",
			},
		},
		RestoreSteps: []Step{
			linuxRestoreStep(iface, snapshotPath),
		},
		Warnings: []string{
			"These settings are host-wide. Apply them only on hosts dedicated to high-throughput transfer workloads or after validating other workloads on the same host.",
			"BBR requires kernel support. If tcp_bbr is unavailable, sysctl --system will report the failure and the host will retain its existing congestion control.",
			"Some cloud NICs keep a multiqueue root qdisc with per-queue fq_codel even when net.core.default_qdisc is set to fq. Do not force a single root qdisc unless you have benchmarked that host.",
		},
		Notes: []string{
			"Run the user inspection steps before and after the privileged steps.",
			"Existing TCP connections keep their existing congestion control; rerun upload tests with new connections after applying changes.",
			"The NIC ring step is intentionally soft-fail because not every driver exposes configurable rings.",
			"Files Agent already has Linux UDP buffer tuning for QUIC in lib/agent/linuxudpbuffer. That package writes /etc/sysctl.d/80-files-agent-udp-buffers.conf with at least 7,500,000 byte net.core.rmem_max and net.core.wmem_max values. This TCP upload plan writes /etc/sysctl.d/99-files-high-throughput.conf with higher buffer ceilings, so both can coexist on agent hosts and the later high-throughput file wins for shared rmem/wmem maxima.",
		},
		References: []Reference{
			{Title: "Linux tcp(7) congestion-control sysctls", URL: "https://man7.org/linux/man-pages/man7/tcp.7.html"},
			{Title: "Linux IP sysctl documentation", URL: "https://docs.kernel.org/networking/ip-sysctl.html"},
		},
	}
}

func linuxSnapshotStep(iface string, snapshotPath string) Step {
	return Step{
		ID:          "linux.snapshot",
		Title:       "Snapshot current high-throughput tuning values",
		Description: "Store the current Files.com-managed sysctl and NIC ring values before applying repair changes.",
		Privilege:   PrivilegeAdministrator,
		Commands: []Command{
			posix(fmt.Sprintf(`SNAPSHOT=%q
IFACE=%s
mkdir -p "$(dirname "$SNAPSHOT")"
{
  printf 'version=1\n'
  printf 'os=linux\n'
  printf 'interface=%%s\n' "$IFACE"
  for key in net.ipv4.tcp_congestion_control net.core.default_qdisc net.ipv4.tcp_slow_start_after_idle net.core.rmem_max net.core.wmem_max net.ipv4.tcp_rmem net.ipv4.tcp_wmem net.ipv4.ip_local_port_range net.ipv4.tcp_mtu_probing net.core.somaxconn; do
    value=$(sysctl -n "$key" 2>/dev/null || true)
    if [ -n "$value" ]; then
      printf 'sysctl.%%s=%%s\n' "$key" "$value"
    fi
  done
  if [ -n "$IFACE" ] && command -v ethtool >/dev/null 2>&1; then
    ethtool -g "$IFACE" 2>/dev/null | awk 'BEGIN { current=0 } /^Current hardware settings:/ { current=1; next } current == 1 && $1 == "RX:" { print "nic.rx="$2 } current == 1 && $1 == "TX:" { print "nic.tx="$2 }'
  fi
} > "$SNAPSHOT"
chmod 600 "$SNAPSHOT"
printf 'Saved Files.com high-throughput OS tuning snapshot to %%s\n' "$SNAPSHOT"`, snapshotPath, iface)),
		},
		ExpectedOutcome: "A pre-change snapshot exists for exact restore when supported by the host.",
	}
}

func linuxRestoreStep(iface string, snapshotPath string) Step {
	return Step{
		ID:          "linux.restore",
		Title:       "Restore high-throughput tuning from snapshot or defaults",
		Description: "Restore saved sysctl and NIC ring values when a snapshot exists; otherwise remove Files.com tuning files and reload the remaining system configuration.",
		Privilege:   PrivilegeAdministrator,
		Commands: []Command{
			posix(fmt.Sprintf(`set -e
SNAPSHOT=%q
IFACE=%s
rm -f /etc/sysctl.d/99-files-high-throughput.conf /etc/modules-load.d/files-high-throughput.conf
if [ -f "$SNAPSHOT" ]; then
  RX=""
  TX=""
  while IFS='=' read -r key value; do
    case "$key" in
      interface)
        if [ -z "$IFACE" ]; then IFACE="$value"; fi
        ;;
      sysctl.net.ipv4.tcp_congestion_control|sysctl.net.core.default_qdisc|sysctl.net.ipv4.tcp_slow_start_after_idle|sysctl.net.core.rmem_max|sysctl.net.core.wmem_max|sysctl.net.ipv4.tcp_rmem|sysctl.net.ipv4.tcp_wmem|sysctl.net.ipv4.ip_local_port_range|sysctl.net.ipv4.tcp_mtu_probing|sysctl.net.core.somaxconn)
        sysctl -w "${key#sysctl.}=$value"
        ;;
      nic.rx)
        RX="$value"
        ;;
      nic.tx)
        TX="$value"
        ;;
    esac
  done < "$SNAPSHOT"
  if [ -n "$IFACE" ] && [ -n "$RX" ] && [ -n "$TX" ] && command -v ethtool >/dev/null 2>&1; then
    ethtool -G "$IFACE" rx "$RX" tx "$TX"
  fi
  printf 'Restored Files.com high-throughput OS tuning snapshot from %%s\n' "$SNAPSHOT"
else
  printf 'No Files.com high-throughput OS tuning snapshot found at %%s. Removing Files.com tuning files and reloading system configuration.\n' "$SNAPSHOT"
  sysctl --system
  printf 'Runtime-only NIC ring settings cannot be restored without a snapshot; reboot or apply known driver defaults if needed.\n'
fi`, snapshotPath, iface)),
		},
		Verification: []Command{
			posix("sysctl net.ipv4.tcp_congestion_control net.core.default_qdisc net.ipv4.tcp_slow_start_after_idle net.core.rmem_max net.core.wmem_max net.ipv4.ip_local_port_range net.ipv4.tcp_mtu_probing"),
			posix(fmt.Sprintf("IFACE=%s; if [ -n \"$IFACE\" ] && command -v ethtool >/dev/null 2>&1; then ethtool -g \"$IFACE\"; fi", iface)),
		},
		ExpectedOutcome: "Host values return to the captured snapshot when available, or the Files.com sysctl/module-load files are removed and remaining system configuration is reloaded.",
	}
}

func posixSingleQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
