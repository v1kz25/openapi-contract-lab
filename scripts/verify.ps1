# Run from the repository root with Go 1.27.4 and gofmt on PATH.
$ErrorActionPreference = 'Stop'
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
if (!$env:OASDIFF_BINARY -or !(Test-Path -LiteralPath $env:OASDIFF_BINARY -PathType Leaf)) {
    throw 'Set OASDIFF_BINARY to the local oasdiff 1.39.0 executable; integration tests must not be skipped.'
}
$goVersion = & go version
if ($LASTEXITCODE -ne 0 -or $goVersion -notmatch '^go version go1\.27\.4 ') { throw 'Go 1.27.4 required' }
$engineVersion = & $env:OASDIFF_BINARY --version
if ($LASTEXITCODE -ne 0 -or $engineVersion.Trim() -ne 'oasdiff version 1.39.0') { throw 'oasdiff 1.39.0 required' }
Write-Output $goVersion
Write-Output $engineVersion
$unformatted = & gofmt -l cmd internal
if ($LASTEXITCODE -ne 0) { throw 'gofmt failed' }
if ($unformatted) { throw ('Run gofmt on: ' + ($unformatted -join ', ')) }
& go test -count=1 ./...
if ($LASTEXITCODE -ne 0) { throw 'go test failed' }
& go vet ./...
if ($LASTEXITCODE -ne 0) { throw 'go vet failed' }
New-Item -ItemType Directory -Force dist | Out-Null
& go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
if ($LASTEXITCODE -ne 0) { throw 'go build failed' }
Write-Output 'PASS: format, tests with local engine, vet and build'
