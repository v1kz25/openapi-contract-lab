# Run from the repository root. No SDK or downloads are needed by this script.
param(
    [string]$Binary = '.\dist\openapi-contract-lab.exe',
    [string]$Engine = '.\dist\oasdiff-1.39.0\oasdiff.exe'
)
$ErrorActionPreference = 'Stop'
$Binary = (Resolve-Path -LiteralPath $Binary).Path
$Engine = (Resolve-Path -LiteralPath $Engine).Path
$base = (Resolve-Path '.\spikes\oasdiff\fixtures\base.yaml').Path
$required = (Resolve-Path '.\spikes\oasdiff\fixtures\required.yaml').Path
$description = (Resolve-Path '.\spikes\oasdiff\fixtures\description.yaml').Path
$directory = Join-Path '.\dist' ('smoke-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $directory | Out-Null

function Invoke-Smoke {
    param([string]$Name, [string[]]$NativeArgs, [int]$ExpectedCode)
    # Windows PowerShell 5.1 turns redirected native stderr into error records.
    # Check the native exit code explicitly, including expected policy failures.
    $ErrorActionPreference = 'Continue'
    $errorPath = Join-Path $directory ($Name + '.stderr')
    $output = @(& $Binary @NativeArgs 2> $errorPath)
    $code = $LASTEXITCODE
    if ($code -ne $ExpectedCode) { throw "$Name exit=$code, expected=$ExpectedCode" }
    $diagnostic = Get-Content -LiteralPath $errorPath -Raw
    if ($code -eq 0 -and $diagnostic) { throw "$Name unexpected stderr" }
    if ($code -eq 3 -and $diagnostic -notmatch 'policy:') { throw "$Name missing policy diagnostic" }
    if (($code -eq 1 -or $code -eq 2) -and $diagnostic -notmatch 'error:') { throw "$Name missing error diagnostic" }
    Write-Host "$Name : exit $code"
    return ($output -join "`n")
}

$helpText = Invoke-Smoke 'help' @('--help') 0
if ($helpText -notmatch 'Usage:' -or $helpText -notmatch 'compare') { throw 'Help regression' }
$version = Invoke-Smoke 'version' @('--version') 0
if ($version.Trim() -ne 'openapi-contract-lab v1.0.0') { throw 'Unexpected release version' }
$compareHelp = Invoke-Smoke 'compare-help' @('compare','--help') 0
if ($compareHelp -notmatch '--fail-on') { throw 'Policy help missing' }
$common = @('compare','--oasdiff',$Engine)
$json = Invoke-Smoke 'breaking' ($common + @($base,$required)) 3
$report = $json | ConvertFrom-Json
if ($report.findings.Count -ne 1 -or $report.findings[0].ruleId -ne 'request-parameter-became-required' -or $report.findings[0].compatibility -ne 'breaking') { throw 'Breaking report regression' }
$summary = 'El par' + [char]0x00e1 + 'metro pasa de opcional a obligatorio.'
if ($report.findings[0].explanation.summary -ne $summary) { throw 'Unicode explanation regression' }
$json = Invoke-Smoke 'compatible' ($common + @($required,$base)) 0
$report = $json | ConvertFrom-Json
if ($report.findings.Count -ne 1 -or $report.findings[0].ruleId -ne 'request-parameter-became-optional' -or $report.findings[0].compatibility -ne 'non-breaking') { throw 'Compatible report regression' }
$json = Invoke-Smoke 'description' ($common + @($base,$description)) 0
if (($json | ConvertFrom-Json).findings.Count -ne 0) { throw 'Description report regression' }
$json = Invoke-Smoke 'no-gate' ($common + @('--fail-on','none',$base,$required)) 0
if (($json | ConvertFrom-Json).findings.Count -ne 1) { throw 'Gate filtered findings' }

foreach ($format in @('json','markdown','html')) {
    $path = Join-Path $directory ('report with spaces.' + $format)
    $arguments = $common + @('--format',$format,'--output',$path,$base,$required)
    $output = Invoke-Smoke ('export-' + $format) $arguments 3
    if ($output) { throw 'Export polluted stdout' }
    $data = [IO.File]::ReadAllBytes((Resolve-Path -LiteralPath $path).Path)
    $decoder = New-Object Text.UTF8Encoding($false, $true)
    $text = $decoder.GetString($data)
    if ($text[0] -eq [char]0xfeff) { throw 'Unexpected BOM' }
    if ($text -notmatch 'request-parameter-became-required') { throw 'Export missing finding' }
    if ($format -eq 'json') {
        if (($text | ConvertFrom-Json).findings[0].explanation.summary -ne $summary) { throw 'JSON export regression' }
    } elseif (!$text.Contains($summary)) { throw 'UTF-8 export regression' }
    $hash = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash
    $null = Invoke-Smoke ('overwrite-' + $format) $arguments 1
    if ((Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash -ne $hash) { throw 'Existing report changed' }
}
$null = Invoke-Smoke 'invalid-policy' ($common + @('--fail-on','invalid',$base,$required)) 2
$null = Invoke-Smoke 'missing-input' ($common + @($base,(Join-Path $directory 'missing.yaml'))) 1
Copy-Item -LiteralPath $base -Destination (Join-Path $directory 'base.yaml')
Copy-Item -LiteralPath $required -Destination (Join-Path $directory 'required.yaml')
$copyBase = Join-Path $directory 'base.yaml'
$copyHead = Join-Path $directory 'required.yaml'
foreach ($original in @($copyBase,$copyHead)) {
    $hash = (Get-FileHash -LiteralPath $original -Algorithm SHA256).Hash
    $null = Invoke-Smoke ('protect-' + [IO.Path]::GetFileName($original)) ($common + @('--output',$original,$copyBase,$copyHead)) 1
    if ((Get-FileHash -LiteralPath $original -Algorithm SHA256).Hash -ne $hash) { throw 'Input file changed' }
}
Write-Output ('PASS smoke; PowerShell ' + $PSVersionTable.PSVersion.ToString())
exit 0
