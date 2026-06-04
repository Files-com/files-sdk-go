package ostuning

import (
	"fmt"
	"strings"
)

func windowsPlan(adapterName string) Plan {
	adapter := adapterName
	if adapter == "" {
		adapter = "<adapter-name>"
	}
	adapterLiteral := powerShellString(adapter)
	snapshotPath := `%ProgramData%\Files.com\os-tuning\high-throughput-upload-snapshot.json`

	return Plan{
		OS:            "windows",
		Profile:       ProfileHighThroughputUpload,
		InterfaceName: adapterName,
		SnapshotPath:  snapshotPath,
		Summary:       "Inspect and normalize Windows TCP auto-tuning, RSS, and receive segment coalescing settings for high-throughput upload clients.",
		UserSteps: []Step{
			{
				ID:          "windows.inspect-tcp",
				Title:       "Inspect TCP global and template settings",
				Description: "Check receive window auto-tuning, RSS, RSC, and active TCP templates.",
				Privilege:   PrivilegeUser,
				Commands: []Command{
					commandPrompt("netsh interface tcp show global"),
					powershell("Get-NetTCPSetting | Select-Object SettingName,AutoTuningLevelLocal,CongestionProvider | Format-Table -AutoSize"),
				},
			},
			{
				ID:          "windows.inspect-adapter",
				Title:       "Inspect adapter RSS and RSC",
				Description: "Check whether the selected network adapter has receive-side scaling and receive segment coalescing enabled.",
				Privilege:   PrivilegeUser,
				Commands: []Command{
					powershell("Get-NetAdapter | Select-Object Name,Status,LinkSpeed | Format-Table -AutoSize"),
					powershell("Get-NetAdapterRss | Format-Table -AutoSize"),
					powershell("Get-NetAdapterRsc | Format-Table -AutoSize"),
				},
				CanFailSoftly: true,
			},
		},
		SnapshotSteps: []Step{
			windowsSnapshotStep(adapterLiteral),
		},
		AdminSteps: []Step{
			{
				ID:          "windows.tcp-autotuning",
				Title:       "Enable normal TCP receive window auto-tuning and RSS",
				Description: "Restore Windows TCP receive window auto-tuning to Normal and enable global RSS. This avoids common under-tuned Windows hosts.",
				Privilege:   PrivilegeAdministrator,
				Commands: []Command{
					commandPrompt("netsh interface tcp set global autotuninglevel=normal rss=enabled"),
					powershell("Set-NetTCPSetting -AutoTuningLevelLocal Normal"),
				},
				Verification: []Command{
					commandPrompt("netsh interface tcp show global"),
					powershell("Get-NetTCPSetting | Select-Object SettingName,AutoTuningLevelLocal,CongestionProvider | Format-Table -AutoSize"),
				},
				ExpectedOutcome: "Receive Window Auto-Tuning Level is Normal and Receive-Side Scaling is enabled.",
			},
			{
				ID:          "windows.adapter-rss-rsc",
				Title:       "Enable adapter RSS and RSC when supported",
				Description: "Enable per-adapter receive-side scaling and receive segment coalescing. Replace the adapter name with the active upload interface.",
				Privilege:   PrivilegeAdministrator,
				Commands: []Command{
					powershell(fmt.Sprintf("Enable-NetAdapterRss -Name %s", adapterLiteral)),
					powershell(fmt.Sprintf("Enable-NetAdapterRsc -Name %s -IPv4", adapterLiteral)),
				},
				Verification: []Command{
					powershell(fmt.Sprintf("Get-NetAdapterRss -Name %s | Format-List", adapterLiteral)),
					powershell(fmt.Sprintf("Get-NetAdapterRsc -Name %s | Format-List", adapterLiteral)),
				},
				CanFailSoftly:   true,
				ExpectedOutcome: "Adapters that support RSS/RSC report RSS and IPv4 RSC enabled.",
			},
		},
		RestoreSteps: []Step{
			windowsRestoreStep(),
		},
		Warnings: []string{
			"Run PowerShell as Administrator for privileged steps.",
			"Do not force a congestion provider globally without validating the Windows version and network path. This plan keeps the OS default congestion provider.",
			"Adapter offload support depends on the NIC driver and virtualization platform; unsupported adapter commands should be treated as non-fatal.",
		},
		Notes: []string{
			"Windows Server documentation treats Normal receive-window auto-tuning as the default setting for most scenarios.",
			"RSS is most relevant on multi-core hosts where network processing must scale across CPUs.",
		},
		References: []Reference{
			{Title: "Microsoft network adapter performance tuning", URL: "https://learn.microsoft.com/en-us/windows-server/networking/technologies/network-subsystem/net-sub-performance-tuning-nics"},
			{Title: "Microsoft netsh interface tcp commands", URL: "https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/netsh-interface"},
		},
	}
}

func windowsSnapshotStep(adapterLiteral string) Step {
	return Step{
		ID:          "windows.snapshot",
		Title:       "Snapshot current high-throughput tuning values",
		Description: "Store the current Windows TCP and selected adapter offload values before applying repair changes.",
		Privilege:   PrivilegeAdministrator,
		Commands: []Command{
			powershell(fmt.Sprintf(`$Snapshot = Join-Path $env:ProgramData 'Files.com\os-tuning\high-throughput-upload-snapshot.json'
$AdapterName = %s
New-Item -ItemType Directory -Path (Split-Path $Snapshot) -Force | Out-Null
$State = [ordered]@{
  Version = 1
  OS = 'windows'
  AdapterName = $AdapterName
  TcpSettings = @(Get-NetTCPSetting | Select-Object SettingName,AutoTuningLevelLocal)
  GlobalOffload = Get-NetOffloadGlobalSetting -ErrorAction SilentlyContinue | Select-Object ReceiveSideScaling
  AdapterRss = if ($AdapterName -and $AdapterName -ne '<adapter-name>') { Get-NetAdapterRss -Name $AdapterName -ErrorAction SilentlyContinue | Select-Object Name,Enabled } else { $null }
  AdapterRsc = if ($AdapterName -and $AdapterName -ne '<adapter-name>') { Get-NetAdapterRsc -Name $AdapterName -ErrorAction SilentlyContinue | Select-Object Name,IPv4Enabled } else { $null }
}
$State | ConvertTo-Json -Depth 6 | Set-Content -Path $Snapshot -Encoding UTF8
Write-Host "Saved Files.com high-throughput OS tuning snapshot to $Snapshot"`, adapterLiteral)),
		},
		ExpectedOutcome: "A pre-change snapshot exists for exact restore of supported Windows TCP and adapter offload values.",
	}
}

func windowsRestoreStep() Step {
	return Step{
		ID:          "windows.restore",
		Title:       "Restore high-throughput tuning from snapshot or defaults",
		Description: "Restore saved Windows TCP and selected adapter offload values when a snapshot exists; otherwise apply the safest supported TCP defaults.",
		Privilege:   PrivilegeAdministrator,
		Commands: []Command{
			powershell(`$ErrorActionPreference = 'Stop'
$Snapshot = Join-Path $env:ProgramData 'Files.com\os-tuning\high-throughput-upload-snapshot.json'
if (Test-Path $Snapshot) {
  $State = Get-Content -Path $Snapshot -Raw | ConvertFrom-Json
  foreach ($Setting in @($State.TcpSettings)) {
    if ($Setting.SettingName -and $Setting.AutoTuningLevelLocal) {
      Set-NetTCPSetting -SettingName $Setting.SettingName -AutoTuningLevelLocal $Setting.AutoTuningLevelLocal -ErrorAction Stop
    }
  }
  if ($State.GlobalOffload -and $State.GlobalOffload.ReceiveSideScaling) {
    Set-NetOffloadGlobalSetting -ReceiveSideScaling $State.GlobalOffload.ReceiveSideScaling -ErrorAction Stop
  }
  if ($State.AdapterRss -and $State.AdapterRss.Name) {
    if ([System.Convert]::ToBoolean($State.AdapterRss.Enabled)) {
      Enable-NetAdapterRss -Name $State.AdapterRss.Name -ErrorAction Stop
    } else {
      Disable-NetAdapterRss -Name $State.AdapterRss.Name -ErrorAction Stop
    }
  }
  if ($State.AdapterRsc -and $State.AdapterRsc.Name) {
    $RscEnabled = [System.Convert]::ToBoolean($State.AdapterRsc.IPv4Enabled)
    if ($RscEnabled) {
      Enable-NetAdapterRsc -Name $State.AdapterRsc.Name -IPv4 -ErrorAction Stop
    } else {
      Disable-NetAdapterRsc -Name $State.AdapterRsc.Name -IPv4 -ErrorAction Stop
    }
  }
  Write-Host "Restored Files.com high-throughput OS tuning snapshot from $Snapshot"
} else {
  Write-Host "No Files.com high-throughput OS tuning snapshot found at $Snapshot."
  netsh interface tcp set global autotuninglevel=normal rss=enabled
  if ($LASTEXITCODE -ne 0) {
    throw "netsh interface tcp set global failed with exit code $LASTEXITCODE"
  }
  Set-NetTCPSetting -AutoTuningLevelLocal Normal -ErrorAction Stop
  Write-Host 'Applied supported Windows TCP defaults. Adapter RSS/RSC driver defaults cannot be inferred without a snapshot.'
}`),
		},
		Verification: []Command{
			commandPrompt("netsh interface tcp show global"),
			powershell("Get-NetTCPSetting | Select-Object SettingName,AutoTuningLevelLocal,CongestionProvider | Format-Table -AutoSize"),
			powershell("Get-NetAdapterRss | Format-Table -AutoSize"),
			powershell("Get-NetAdapterRsc | Format-Table -AutoSize"),
		},
		ExpectedOutcome: "Host values return to the captured snapshot when available, or supported Windows TCP defaults are applied while adapter defaults are reported as snapshot-dependent.",
	}
}

func powerShellString(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
