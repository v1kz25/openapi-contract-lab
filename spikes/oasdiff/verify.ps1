# Requires PowerShell 7. Run from any directory; no downloads or installation.
param(
    [string]$Binary = './dist/oasdiff-1.39.0/oasdiff.exe',
    [string]$OutputDirectory = './dist/oasdiff-evidence'
)

$ErrorActionPreference = 'Stop'
$repoRoot = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
Push-Location $repoRoot
try {
    $binaryPath = (Resolve-Path -LiteralPath $Binary).Path
    $version = & $binaryPath --version
    if ($LASTEXITCODE -ne 0 -or $version -ne 'oasdiff version 1.39.0') {
        throw "Expected oasdiff 1.39.0, received: $version"
    }
    New-Item -ItemType Directory -Force -Path $OutputDirectory | Out-Null
    $rows = [System.Collections.Generic.List[object]]::new()
    $config = @('--config', 'spikes/oasdiff/offline.yaml')
    $base = 'spikes/oasdiff/fixtures/base.yaml'

    function Invoke-Probe([string]$Name, [string[]]$Arguments, [int]$ExpectedExit = 0) {
        $start = [System.Diagnostics.ProcessStartInfo]::new($binaryPath)
        $start.UseShellExecute = $false
        $start.RedirectStandardOutput = $true
        $start.RedirectStandardError = $true
        foreach ($argument in $Arguments) { $start.ArgumentList.Add($argument) }
        $process = [System.Diagnostics.Process]::Start($start)
        try {
            $stdoutTask = $process.StandardOutput.ReadToEndAsync()
            $stderrTask = $process.StandardError.ReadToEndAsync()
            $process.WaitForExit()
            $stdout = $stdoutTask.GetAwaiter().GetResult()
            $stderr = $stderrTask.GetAwaiter().GetResult()
            $rows.Add([pscustomobject]@{
                name = $Name; command = 'oasdiff ' + ($Arguments -join ' ')
                exit = $process.ExitCode; stderr = $stderr
            })
            if ($process.ExitCode -ne $ExpectedExit -or $stderr.Length -gt 0) {
                throw "$Name failed: exit=$($process.ExitCode), stderr=$stderr"
            }
            # Normalize line endings and terminal blank lines for versioned evidence.
            $stdout = $stdout.Replace("`r`n", "`n").TrimEnd("`n") + "`n"
            [System.IO.File]::WriteAllText((Join-Path $OutputDirectory $Name), $stdout)
            return $stdout
        } finally { $process.Dispose() }
    }

    foreach ($fixture in @('base', 'required', 'description')) {
        $result = Invoke-Probe "$fixture.validate.json" (@('validate', "spikes/oasdiff/fixtures/$fixture.yaml", '-f', 'json') + $config)
        if ($result.Trim() -ne '[]') { throw "Unexpected validation findings: $fixture" }
    }
    foreach ($case in @('required', 'description')) {
        $revision = "spikes/oasdiff/fixtures/$case.yaml"
        foreach ($mode in @('breaking', 'changelog', 'diff')) {
            foreach ($format in @('json', 'yaml', 'markdown')) {
                $extension = if ($format -eq 'markdown') { 'md' } else { $format }
                $null = Invoke-Probe "$case.$mode.$extension" (@($mode, $base, $revision, '-f', $format) + $config)
            }
        }
        $expected = if ($case -eq 'required') { 1 } else { 0 }
        $null = Invoke-Probe "$case.fail-on.json" (@('breaking', $base, $revision, '-f', 'json', '--fail-on', 'ERR') + $config) $expected
    }

    $breaking = @(Get-Content (Join-Path $OutputDirectory 'required.breaking.json') -Raw | ConvertFrom-Json)
    if ($breaking.Count -ne 1 -or $breaking[0].id -ne 'request-parameter-became-required' -or
        $breaking[0].level -ne 3 -or $breaking[0].operation -ne 'GET' -or $breaking[0].path -ne '/items') {
        throw 'Unexpected breaking classification'
    }
    foreach ($mode in @('breaking', 'changelog')) {
        if ((Get-Content (Join-Path $OutputDirectory "description.$mode.json") -Raw).Trim() -ne '[]') {
            throw "Description edit unexpectedly reported by $mode"
        }
    }
    $diff = Get-Content (Join-Path $OutputDirectory 'description.diff.json') -Raw | ConvertFrom-Json
    $description = $diff.paths.modified.'/items'.operations.modified.GET.parameters.modified.query.filter.description
    if ($description.from -ne 'Filter items.' -or $description.to -ne 'Filter items by name.') {
        throw 'Description edit missing from full diff'
    }
    $rows | ConvertTo-Json -Depth 5 | Set-Content (Join-Path $OutputDirectory 'commands.json') -Encoding utf8NoBOM
    Write-Output "$version; $($rows.Count) commands and classification assertions passed. Evidence: $OutputDirectory"
} finally { Pop-Location }
