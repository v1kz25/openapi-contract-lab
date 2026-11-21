#requires -Version 7.0
# Synthetic spike only; never called by the application. No downloads.
param(
    [string]$Binary = './dist/oasdiff-1.39.0/oasdiff.exe',
    [string]$OutputDirectory = ('./dist/change-kind-' + [guid]::NewGuid().ToString('N'))
)
$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
Push-Location $repoRoot
try {
    $binaryPath = (Resolve-Path -LiteralPath $Binary).Path
    $outputPath = [IO.Path]::GetFullPath($OutputDirectory)
    # Do not overwrite previous observations or contracts.
    New-Item -ItemType Directory -Path $outputPath -ErrorAction Stop | Out-Null
    $poison = Join-Path $outputPath 'ambient-poison.yaml'
    [IO.File]::WriteAllText($poison, "invalid: [`n")
    $config = 'spikes/change-kind/offline.yaml'
    $rows = [Collections.Generic.List[object]]::new()
    $cases = @(Get-Content -LiteralPath 'spikes/change-kind/cases.json' -Raw | ConvertFrom-Json)
    if ($cases.Count -ne 5) { throw 'Expected five reviewed pairs' }
    $original = @{}
    foreach ($file in @('spikes/change-kind/fixtures/external-ref.yaml','spikes/change-kind/fixtures/external.yaml') + @($cases.base) + @($cases.head)) {
        $original[$file] = [Convert]::ToBase64String([IO.File]::ReadAllBytes((Join-Path $repoRoot $file)))
    }

    function Invoke-Probe([string]$Name, [string[]]$NativeArgs, [int]$ExpectedExit = 0, [switch]$ExternalError) {
        $start = [Diagnostics.ProcessStartInfo]::new($binaryPath)
        $start.WorkingDirectory = $repoRoot
        $start.UseShellExecute = $false
        $start.RedirectStandardOutput = $true
        $start.RedirectStandardError = $true
        foreach ($key in @($start.Environment.Keys)) {
            if ($key.StartsWith('OASDIFF_', [StringComparison]::OrdinalIgnoreCase)) { $null = $start.Environment.Remove($key) }
        }
        # Only the child sees the invalid ambient config. Explicit --config must win.
        $start.Environment['OASDIFF_CONFIG'] = $poison
        foreach ($argument in $NativeArgs) { $start.ArgumentList.Add($argument) }
        $timer = [Diagnostics.Stopwatch]::StartNew()
        $process = [Diagnostics.Process]::Start($start)
        try {
            $stdoutTask = $process.StandardOutput.ReadToEndAsync()
            $stderrTask = $process.StandardError.ReadToEndAsync()
            if (!$process.WaitForExit(30000)) { $process.Kill($true); throw "$Name exceeded 30 seconds" }
            $stdout = $stdoutTask.GetAwaiter().GetResult()
            $stderr = $stderrTask.GetAwaiter().GetResult()
            $timer.Stop()
            $rows.Add([pscustomobject]@{
                name = $Name; command = 'oasdiff ' + ($NativeArgs -join ' ')
                exit = $process.ExitCode; elapsedMs = $timer.ElapsedMilliseconds
                stdoutBytes = [Text.Encoding]::UTF8.GetByteCount($stdout)
                stderrBytes = [Text.Encoding]::UTF8.GetByteCount($stderr)
            })
            if ($process.ExitCode -ne $ExpectedExit) { throw "$Name exit=$($process.ExitCode), expected=$ExpectedExit : $stderr" }
            if ($ExternalError) {
                if ($stderr -notmatch 'disallowed external reference' -or $stdout.Trim()) { throw "$Name did not reject external refs" }
            } elseif ($stderr) { throw "$Name unexpected stderr: $stderr" }
            $stdout = $stdout.Replace("`r`n", "`n").TrimEnd("`n") + "`n"
            $extension = if ($Name -eq 'version') { '.txt' } elseif ($ExternalError) { '.stdout' } else { '.json' }
            [IO.File]::WriteAllText((Join-Path $outputPath ($Name + $extension)), $stdout)
            if ($ExternalError) { [IO.File]::WriteAllText((Join-Path $outputPath ($Name + '.stderr.txt')), $stderr.Replace("`r`n","`n")) }
            return $stdout
        } finally { $process.Dispose() }
    }

    $version = Invoke-Probe 'version' @('--config',$config,'--version')
    if ($version.Trim() -ne 'oasdiff version 1.39.0') { throw 'oasdiff 1.39.0 required' }
    $common = @('--config',$config,'--format','json','--allow-external-refs=false','--fetch=false')
    foreach ($case in $cases) {
        foreach ($mode in @('diff','summary','changelog')) {
            $arguments = @($mode) + $common
            if ($mode -eq 'changelog') { $arguments += @('--open=false','--lang','en','--level','INFO') }
            $arguments += @('--',$case.base,$case.head)
            $name = $case.name + '.' + $mode
            $result = Invoke-Probe $name $arguments
            $expected = [IO.File]::ReadAllText((Join-Path $PSScriptRoot ('evidence/' + $name + '.json'))).Replace("`r`n", "`n")
            if ($result -cne $expected) { throw "$name differs from reviewed evidence; inspect, do not auto-update" }
            if ($mode -eq 'summary') {
                $summary = $result | ConvertFrom-Json
                if ($summary.diff -isnot [bool] -or $summary.diff -ne $case.diff) { throw "$name unexpected diff signal" }
            } elseif ($mode -eq 'changelog') {
                $changes = @(ConvertFrom-Json -InputObject $result)
                if ($changes.Count -ne $case.ids.Count) { throw "$name unexpected changelog count" }
                for ($i=0; $i -lt $changes.Count; $i++) {
                    if ($changes[$i].id -ne $case.ids[$i] -or $changes[$i].level -ne $case.levels[$i]) { throw "$name unexpected ID/level" }
                }
            }
        }
        $expectedExit = if ($case.diff) { 1 } else { 0 }
        $null = Invoke-Probe ($case.name + '.fail-on-diff') (@('diff') + $common + @('--fail-on-diff','--',$case.base,$case.head)) $expectedExit
    }
    $description = $cases | Where-Object name -eq 'description'
    $filtered = Invoke-Probe 'description.exclude-description' (@('diff') + $common + @('--exclude-elements','description','--',$description.base,$description.head))
    if ($filtered.Trim() -ne '{}') { throw 'Description exclusion probe changed' }
    foreach ($mode in @('diff','summary','changelog')) {
        $arguments = @($mode) + $common
        if ($mode -eq 'changelog') { $arguments += '--open=false' }
        $null = Invoke-Probe ('external-ref.'+$mode) ($arguments + @('--','spikes/oasdiff/fixtures/base.yaml','spikes/change-kind/fixtures/external-ref.yaml')) 123 -ExternalError
    }
    foreach ($file in $original.Keys) {
        if ([Convert]::ToBase64String([IO.File]::ReadAllBytes((Join-Path $repoRoot $file))) -cne $original[$file]) { throw "Input modified: $file" }
    }
    [IO.File]::WriteAllText((Join-Path $outputPath 'commands.json'), (($rows | ConvertTo-Json -Depth 5).Replace("`r`n","`n") + "`n"))
    Write-Output "PASS: $($rows.Count) probes, reviewed evidence, isolated config, blocked external refs and unchanged inputs. $outputPath"
} finally { Pop-Location }
