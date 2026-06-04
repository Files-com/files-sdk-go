package ostuning

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

type Privilege string

const (
	PrivilegeUser          Privilege = "user"
	PrivilegeAdministrator Privilege = "administrator"
)

type Shell string

const (
	ShellPOSIX      Shell = "posix"
	ShellPowerShell Shell = "powershell"
	ShellCommand    Shell = "cmd"
)

type Command struct {
	Shell       Shell
	CommandLine string
}

type Step struct {
	ID              string
	Title           string
	Description     string
	Privilege       Privilege
	Commands        []Command
	Verification    []Command
	RequiresReboot  bool
	RuntimeOnly     bool
	CanFailSoftly   bool
	ExpectedOutcome string
}

type Reference struct {
	Title string
	URL   string
}

type Plan struct {
	OS            string
	Profile       string
	InterfaceName string
	SnapshotPath  string
	Summary       string
	UserSteps     []Step
	NetworkTests  []Step
	SnapshotSteps []Step
	AdminSteps    []Step
	RestoreSteps  []Step
	Warnings      []string
	Notes         []string
	References    []Reference
}

type RunOptions struct {
	Stdout        io.Writer
	Stderr        io.Writer
	DryRun        bool
	StopOnError   bool
	BeforeStep    func(Step)
	BeforeCommand func(Command)
}

type CommandResult struct {
	Command Command
	DryRun  bool
	Err     error
}

type StepResult struct {
	Step                Step
	CommandResults      []CommandResult
	SkippedForPrivilege bool
	Err                 error
	SoftFailed          bool
}

type Options struct {
	OS                 string
	InterfaceName      string
	IncludeNetworkTest bool
}

const ProfileHighThroughputUpload = "high-throughput-upload"

const (
	// MinimumOpenFileLimit is the lowest soft nofile value OS tuning treats as
	// acceptable for high-throughput adaptive uploads.
	MinimumOpenFileLimit = 8192
	// PreferredOpenFileLimit is the recommended soft nofile value for hosts
	// dedicated to high-throughput transfer workloads.
	PreferredOpenFileLimit = 65536
)

type OpenFileLimitResult struct {
	Supported  bool
	BeforeSoft uint64
	BeforeHard uint64
	AfterSoft  uint64
	Changed    bool
}

var supportedOS = []string{"linux", "darwin", "windows"}

func SupportedOS() []string {
	return slices.Clone(supportedOS)
}

func HighThroughputUploadPlan(options Options) (Plan, error) {
	targetOS := normalizeOS(options.OS)
	if targetOS == "" {
		targetOS = runtime.GOOS
	}

	switch targetOS {
	case "linux":
		return linuxPlan(options.InterfaceName), nil
	case "darwin":
		return darwinPlan(options.IncludeNetworkTest), nil
	case "windows":
		return windowsPlan(options.InterfaceName), nil
	default:
		return Plan{}, fmt.Errorf("unsupported OS %q for high-throughput tuning; supported values are %s", targetOS, strings.Join(supportedOS, ", "))
	}
}

func (p Plan) Steps() []Step {
	steps := make([]Step, 0, len(p.UserSteps)+len(p.NetworkTests)+len(p.SnapshotSteps)+len(p.AdminSteps))
	steps = append(steps, p.UserSteps...)
	steps = append(steps, p.NetworkTests...)
	steps = append(steps, p.SnapshotSteps...)
	steps = append(steps, p.AdminSteps...)
	return steps
}

func (p Plan) VerificationSteps() []Step {
	steps := slices.Clone(p.UserSteps)
	steps = append(steps, p.NetworkTests...)
	for _, step := range p.AdminSteps {
		if len(step.Verification) == 0 {
			continue
		}
		steps = append(steps, Step{
			ID:              step.ID + ".verify",
			Title:           "Verify " + lowerFirst(step.Title),
			Description:     "Verify the expected state for: " + step.Title,
			Privilege:       PrivilegeUser,
			Commands:        slices.Clone(step.Verification),
			ExpectedOutcome: step.ExpectedOutcome,
		})
	}
	return steps
}

func (p Plan) RepairSteps() []Step {
	steps := make([]Step, 0, len(p.NetworkTests)*2+len(p.SnapshotSteps)+len(p.AdminSteps))
	steps = append(steps, p.RepairPreflightSteps()...)
	steps = append(steps, p.RepairChangeSteps()...)
	steps = append(steps, p.RepairPostChangeSteps()...)
	return steps
}

func (p Plan) RepairPreflightSteps() []Step {
	return labeledNetworkTests(p.NetworkTests, "before repair")
}

func (p Plan) RepairChangeSteps() []Step {
	steps := slices.Clone(p.SnapshotSteps)
	steps = append(steps, p.AdminSteps...)
	return steps
}

func (p Plan) RepairPostChangeSteps() []Step {
	return labeledNetworkTests(p.NetworkTests, "after repair")
}

func (p Plan) RestorePlanSteps() []Step {
	return slices.Clone(p.RestoreSteps)
}

func CurrentProcessElevated() bool {
	return currentProcessElevated()
}

func StepRequiresElevation(step Step) bool {
	return step.Privilege == PrivilegeAdministrator
}

func SplitStepsByElevation(steps []Step, elevated bool) (runnable []Step, skipped []Step) {
	for _, step := range steps {
		if StepRequiresElevation(step) && !elevated {
			skipped = append(skipped, step)
			continue
		}
		runnable = append(runnable, step)
	}
	return runnable, skipped
}

func RunSteps(ctx context.Context, steps []Step, options RunOptions) []StepResult {
	return RunStepsWithElevation(ctx, steps, CurrentProcessElevated(), options)
}

func RunStepsWithElevation(ctx context.Context, steps []Step, elevated bool, options RunOptions) []StepResult {
	results := make([]StepResult, 0, len(steps))
	for _, step := range steps {
		stepResult := StepResult{Step: step}
		if StepRequiresElevation(step) && !elevated {
			stepResult.SkippedForPrivilege = true
			results = append(results, stepResult)
			continue
		}

		if options.BeforeStep != nil {
			options.BeforeStep(step)
		}

		stopAfterStep := false
		for _, command := range step.Commands {
			commandResults := RunCommands(ctx, []Command{command}, RunOptions{
				Stdout:        options.Stdout,
				Stderr:        options.Stderr,
				DryRun:        options.DryRun,
				StopOnError:   true,
				BeforeCommand: options.BeforeCommand,
			})
			stepResult.CommandResults = append(stepResult.CommandResults, commandResults...)
			if len(commandResults) == 0 || commandResults[0].Err == nil {
				continue
			}

			stepResult.Err = commandResults[0].Err
			if step.CanFailSoftly {
				stepResult.SoftFailed = true
				continue
			}
			stopAfterStep = options.StopOnError
			if options.StopOnError {
				break
			}
		}
		results = append(results, stepResult)
		if stopAfterStep {
			return results
		}
	}
	return results
}

func CommandsForSteps(steps []Step) []Command {
	var commands []Command
	for _, step := range steps {
		commands = append(commands, step.Commands...)
	}
	return commands
}

func RunCommands(ctx context.Context, commands []Command, options RunOptions) []CommandResult {
	results := make([]CommandResult, 0, len(commands))
	for _, command := range commands {
		if options.BeforeCommand != nil {
			options.BeforeCommand(command)
		}
		result := CommandResult{Command: command, DryRun: options.DryRun}
		if !options.DryRun {
			result.Err = runCommand(ctx, command, options.Stdout, options.Stderr)
		}
		results = append(results, result)
		if result.Err != nil && options.StopOnError {
			return results
		}
	}
	return results
}

func normalizeOS(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return ""
	case "mac", "macos", "osx":
		return "darwin"
	case "win":
		return "windows"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func labeledNetworkTests(steps []Step, label string) []Step {
	if len(steps) == 0 {
		return nil
	}

	labeled := make([]Step, 0, len(steps))
	idSuffix := strings.ReplaceAll(label, " ", "-")
	for _, step := range steps {
		clone := step
		clone.ID = step.ID + "." + idSuffix
		clone.Title = step.Title + " " + label
		if clone.Description != "" {
			clone.Description = strings.TrimRight(clone.Description, ".") + " " + networkTestLabelDescription(label)
		}
		labeled = append(labeled, clone)
	}
	return labeled
}

func networkTestLabelDescription(label string) string {
	switch label {
	case "before repair":
		return "before applying host-wide repair changes."
	case "after repair":
		return "after applying host-wide repair changes so the result can be compared with the baseline."
	default:
		return label + "."
	}
}

func posix(command string) Command {
	return Command{Shell: ShellPOSIX, CommandLine: strings.TrimSpace(command)}
}

func powershell(command string) Command {
	return Command{Shell: ShellPowerShell, CommandLine: strings.TrimSpace(command)}
}

func commandPrompt(command string) Command {
	return Command{Shell: ShellCommand, CommandLine: strings.TrimSpace(command)}
}

func runCommand(ctx context.Context, command Command, stdout io.Writer, stderr io.Writer) error {
	program := "sh"
	args := []string{"-c", command.CommandLine}
	switch command.Shell {
	case ShellPowerShell:
		program = "powershell"
		args = []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", command.CommandLine}
	case ShellCommand:
		program = "cmd"
		args = []string{"/C", command.CommandLine}
	}

	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
