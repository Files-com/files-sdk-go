package ostuning

import (
	"context"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHighThroughputUploadPlanLinux(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "linux", InterfaceName: "ens4"})
	require.NoError(t, err)
	require.Equal(t, "linux", plan.OS)
	require.Len(t, plan.UserSteps, 2)
	require.Empty(t, plan.NetworkTests)
	require.Len(t, plan.SnapshotSteps, 1)
	require.Len(t, plan.AdminSteps, 2)
	require.Len(t, plan.RestoreSteps, 1)
	require.Equal(t, "/var/lib/files.com/os-tuning/high-throughput-upload.snapshot", plan.SnapshotPath)
	require.Contains(t, plan.AdminSteps[0].Commands[3].CommandLine, "net.ipv4.tcp_congestion_control = bbr")
	require.Contains(t, plan.AdminSteps[0].Commands[3].CommandLine, "net.ipv4.ip_local_port_range = 10000 65535")
	require.NotContains(t, plan.AdminSteps[0].Commands[4].CommandLine, "sudo")
	require.Contains(t, plan.AdminSteps[1].Commands[0].CommandLine, "ens4")
	require.Contains(t, plan.AdminSteps[1].Commands[0].CommandLine, "IFACE='ens4'")
	require.NotContains(t, plan.AdminSteps[1].Commands[0].CommandLine, "sudo")
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "sysctl --system")
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "set -e")
	require.NotContains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "|| true")
	require.False(t, plan.RestoreSteps[0].CanFailSoftly)
}

func TestHighThroughputUploadPlanLinuxEscapesInterfaceName(t *testing.T) {
	interfaceName := "eth$Primary's`name`"
	plan, err := HighThroughputUploadPlan(Options{OS: "linux", InterfaceName: interfaceName})
	require.NoError(t, err)

	expectedAssignment := "IFACE=" + posixSingleQuote(interfaceName)
	require.Contains(t, plan.UserSteps[1].Commands[0].CommandLine, expectedAssignment)
	require.Contains(t, plan.SnapshotSteps[0].Commands[0].CommandLine, expectedAssignment)
	require.Contains(t, plan.AdminSteps[1].Commands[0].CommandLine, expectedAssignment)
	require.Contains(t, plan.AdminSteps[1].Verification[0].CommandLine, expectedAssignment)
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, expectedAssignment)
	require.Contains(t, plan.RestoreSteps[0].Verification[1].CommandLine, expectedAssignment)
	require.NotContains(t, plan.AdminSteps[1].Commands[0].CommandLine, `IFACE="eth$Primary`)
}

func TestHighThroughputUploadPlanDarwin(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "macos"})
	require.NoError(t, err)
	require.Equal(t, "darwin", plan.OS)
	require.Len(t, plan.UserSteps, 1)
	require.Empty(t, plan.NetworkTests)
	require.Len(t, plan.SnapshotSteps, 1)
	require.Len(t, plan.AdminSteps, 1)
	require.Len(t, plan.RestoreSteps, 1)
	require.NotContains(t, CommandsForSteps(plan.UserSteps)[0].CommandLine, "networkQuality")
	require.Contains(t, plan.AdminSteps[0].Commands[0].CommandLine, "kern.ipc.maxsockbuf")
	require.True(t, plan.AdminSteps[0].RuntimeOnly)
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "rebooting is the safest")
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "set -e")
	require.NotContains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "|| true")
	require.False(t, plan.RestoreSteps[0].CanFailSoftly)
}

func TestHighThroughputUploadPlanDarwinNetworkQualityIsOptIn(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "darwin", IncludeNetworkTest: true})
	require.NoError(t, err)
	require.Len(t, plan.UserSteps, 1)
	require.Len(t, plan.NetworkTests, 1)
	require.Contains(t, plan.NetworkTests[0].Commands[0].CommandLine, "networkQuality -v")

	verificationSteps := plan.VerificationSteps()
	require.Contains(t, verificationSteps[1].Title, "Measure Apple network quality")
	require.NotContains(t, verificationSteps[1].Title, "before repair")

	repairSteps := plan.RepairSteps()
	require.Contains(t, repairSteps[0].Title, "before repair")
	require.Contains(t, repairSteps[len(repairSteps)-1].Title, "after repair")
	require.Contains(t, repairSteps[0].Commands[0].CommandLine, "networkQuality -v")
	require.Contains(t, repairSteps[len(repairSteps)-1].Commands[0].CommandLine, "networkQuality -v")
}

func TestHighThroughputUploadPlanWindows(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "windows", InterfaceName: "Ethernet 2"})
	require.NoError(t, err)
	require.Equal(t, "windows", plan.OS)
	require.Len(t, plan.UserSteps, 2)
	require.Empty(t, plan.NetworkTests)
	require.Len(t, plan.SnapshotSteps, 1)
	require.Len(t, plan.AdminSteps, 2)
	require.Len(t, plan.RestoreSteps, 1)
	require.Equal(t, ShellCommand, plan.AdminSteps[0].Commands[0].Shell)
	require.Contains(t, plan.AdminSteps[1].Commands[0].CommandLine, "'Ethernet 2'")
	require.Contains(t, plan.AdminSteps[1].Commands[1].CommandLine, "-IPv4")
	require.Contains(t, plan.SnapshotSteps[0].Commands[0].CommandLine, "high-throughput-upload-snapshot.json")
	require.Contains(t, plan.SnapshotSteps[0].Commands[0].CommandLine, "IPv4Enabled")
	require.NotContains(t, plan.SnapshotSteps[0].Commands[0].CommandLine, "IPv6Enabled")
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "ConvertFrom-Json")
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "ErrorAction Stop")
	require.NotContains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "SilentlyContinue")
	require.Contains(t, plan.RestoreSteps[0].Commands[0].CommandLine, "-IPv4")
	require.False(t, plan.RestoreSteps[0].CanFailSoftly)
}

func TestHighThroughputUploadPlanWindowsEscapesPowerShellAdapterName(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "windows", InterfaceName: "Ethernet $Primary's"})
	require.NoError(t, err)

	require.Contains(t, plan.AdminSteps[1].Commands[0].CommandLine, "-Name 'Ethernet $Primary''s'")
	require.Contains(t, plan.AdminSteps[1].Commands[1].CommandLine, "-Name 'Ethernet $Primary''s' -IPv4")
	require.Contains(t, plan.SnapshotSteps[0].Commands[0].CommandLine, "$AdapterName = 'Ethernet $Primary''s'")
	require.NotContains(t, plan.AdminSteps[1].Commands[0].CommandLine, `"Ethernet $Primary's"`)
}

func TestHighThroughputUploadPlanDefaultOS(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{})
	if slicesContains(SupportedOS(), runtime.GOOS) {
		require.NoError(t, err)
		require.Equal(t, runtime.GOOS, plan.OS)
		return
	}
	require.Error(t, err)
}

func TestHighThroughputUploadPlanUnsupportedOS(t *testing.T) {
	_, err := HighThroughputUploadPlan(Options{OS: "plan9"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported OS")
}

func TestPlanStepsPreservesUserThenAdminOrder(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "linux"})
	require.NoError(t, err)
	steps := plan.Steps()
	require.Len(t, steps, len(plan.UserSteps)+len(plan.SnapshotSteps)+len(plan.AdminSteps))
	require.True(t, strings.HasPrefix(steps[0].ID, "linux.inspect"))
	require.Equal(t, plan.SnapshotSteps[0].ID, steps[len(plan.UserSteps)].ID)
	require.Equal(t, plan.AdminSteps[0].ID, steps[len(plan.UserSteps)+len(plan.SnapshotSteps)].ID)
}

func TestPlanVerificationAndRepairSteps(t *testing.T) {
	plan, err := HighThroughputUploadPlan(Options{OS: "linux"})
	require.NoError(t, err)

	verificationSteps := plan.VerificationSteps()
	require.Greater(t, len(verificationSteps), len(plan.UserSteps))
	require.Contains(t, verificationSteps[len(plan.UserSteps)].ID, ".verify")
	require.Equal(t, plan.AdminSteps[0].ExpectedOutcome, verificationSteps[len(plan.UserSteps)].ExpectedOutcome)

	repairSteps := plan.RepairSteps()
	require.Len(t, repairSteps, len(plan.SnapshotSteps)+len(plan.AdminSteps))
	require.Equal(t, plan.SnapshotSteps[0].ID, repairSteps[0].ID)
	require.Equal(t, plan.AdminSteps[0].ID, repairSteps[1].ID)

	restoreSteps := plan.RestorePlanSteps()
	require.Len(t, restoreSteps, len(plan.RestoreSteps))
	require.Equal(t, plan.RestoreSteps[0].ID, restoreSteps[0].ID)
}

func TestSplitStepsByElevation(t *testing.T) {
	steps := []Step{
		{ID: "user", Privilege: PrivilegeUser},
		{ID: "admin", Privilege: PrivilegeAdministrator},
	}

	runnable, skipped := SplitStepsByElevation(steps, false)
	require.Len(t, runnable, 1)
	require.Equal(t, "user", runnable[0].ID)
	require.Len(t, skipped, 1)
	require.Equal(t, "admin", skipped[0].ID)

	runnable, skipped = SplitStepsByElevation(steps, true)
	require.Len(t, runnable, 2)
	require.Empty(t, skipped)
}

func TestRunStepsWithElevationSkipsPrivilegedSteps(t *testing.T) {
	steps := []Step{
		{ID: "user", Privilege: PrivilegeUser, Commands: []Command{posix("echo user")}},
		{ID: "admin", Privilege: PrivilegeAdministrator, Commands: []Command{posix("echo admin")}},
	}

	results := RunStepsWithElevation(context.Background(), steps, false, RunOptions{DryRun: true, StopOnError: true})
	require.Len(t, results, 2)
	require.False(t, results[0].SkippedForPrivilege)
	require.Len(t, results[0].CommandResults, 1)
	require.True(t, results[1].SkippedForPrivilege)
	require.Empty(t, results[1].CommandResults)
}

func TestRunStepsWithElevationHardFailureStops(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("posix command execution is not available on Windows test hosts")
	}

	steps := []Step{
		{ID: "restore", Privilege: PrivilegeUser, Commands: []Command{posix("exit 17")}},
		{ID: "after", Privilege: PrivilegeUser, Commands: []Command{posix("echo should-not-run")}},
	}

	results := RunStepsWithElevation(context.Background(), steps, true, RunOptions{StopOnError: true})
	require.Len(t, results, 1)
	require.Equal(t, "restore", results[0].Step.ID)
	require.Error(t, results[0].Err)
	require.False(t, results[0].SoftFailed)
}

func TestRunStepsWithElevationSoftFailureContinues(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("posix command execution is not available on Windows test hosts")
	}

	steps := []Step{
		{ID: "optional", Privilege: PrivilegeUser, CanFailSoftly: true, Commands: []Command{posix("exit 17")}},
		{ID: "after", Privilege: PrivilegeUser, Commands: []Command{posix("echo ok")}},
	}

	results := RunStepsWithElevation(context.Background(), steps, true, RunOptions{StopOnError: true})
	require.Len(t, results, 2)
	require.Equal(t, "optional", results[0].Step.ID)
	require.Error(t, results[0].Err)
	require.True(t, results[0].SoftFailed)
	require.Equal(t, "after", results[1].Step.ID)
	require.NoError(t, results[1].Err)
}

func TestRunCommandsDryRun(t *testing.T) {
	results := RunCommands(context.Background(), []Command{posix("exit 1")}, RunOptions{DryRun: true, StopOnError: true})
	require.Len(t, results, 1)
	require.True(t, results[0].DryRun)
	require.NoError(t, results[0].Err)
}

func slicesContains(values []string, value string) bool {
	for _, current := range values {
		if current == value {
			return true
		}
	}
	return false
}
