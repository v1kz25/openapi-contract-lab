# Publicación y preparación de candidatos

## Estado vigente

[**v1.0.0 estable publicada**](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v1.0.0), no draft ni prerelease, desde el 2026-11-22. PR #46 fusionada con autorización explícita; tag anotado sobre `main` `342aaae133b3807e92b487c51a4316a8dbf3cbce`, árbol idéntico al candidato aprobado `3939de88994dfbd5240c9e1217e9f415a54b539e`. ZIP SHA-256 `5f7273b2ef1046e8c4da7c19e42b35be385865a2174f5e8450603a392ad16061`, ejecutable `f072fbb75fd6bbbf3a65696995ae3c7c596de902b83f6f208574ad7c495981be`: descargados tras publicar, coincidentes con los aprobados y smoke PS5.1/7 correcto. [CI del candidato exacto](https://github.com/v1kz25/openapi-contract-lab/actions/runs/45040822531) SUCCESS.

El propietario autorizó también finalizar la actualización documental y sincronizar hacia `develop` mediante PR con comprobaciones. No se cambian tag, notas empaquetadas ni bytes publicados: los textos de revisión dentro del ZIP registran su preparación histórica. Para instalar, [INSTALL.md](INSTALL.md). Las siguientes referencias v0.2.0 y a #40 son historia/procedimiento, no una publicación pendiente actual.

[v0.2.0 experimental está publicada](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.2.0) desde el 2026-11-22, como prerelease bajo **Apache-2.0**, con ZIP Windows x64 y SHA256SUMS.txt, sin redistribuir oasdiff. #33 completado, PR #34 integrada en `main` y #35 sincronizada a `develop`. El tag desreferenciado es `72ebbd87558bd73ae138696fd9dc62a58b6e2c04`. [v0.1.0](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.1.0) sigue disponible como versión histórica; sus notas y artefactos se conservan.

Las [notas v0.2.0](releases/v0.2.0.md), BUILD.txt y [guía #33](CANDIDATE-V0.2.0.md) registran la preparación anterior a publicar; sus frases «no publicado» son históricas. No modificar tags, archivos ni notas publicadas para corregir ese estado. Usar [INSTALL.md](INSTALL.md) para instalación/actualización vigente. `develop` contiene estabilización posterior de #37/#38, y `release/1.0.0` prepara el candidato siguiente; no confundir ninguno con bytes de una versión ya publicada.

## Candidato v1.0.0 — issue #40

**Registro histórico de preparación, ya completada y publicada por #41.** Los siguientes párrafos conservan el contexto de revisión anterior a autorizar/publicar; el estado vigente y la evidencia de publicación están arriba.

`release/1.0.0` parte de `develop` aprobado `b3f0d3c72758ffe318ae96ca7f0aec81219001da`. #36–#39 ya revisados e integrados; sus cierres anteriores a `main` se verifican expresamente porque GitHub no cierra automáticamente los issues al fusionar a `develop`. El candidato solo ajusta identificación/expectativas a v1.0.0, cabecera de avisos y documentación. [Notas, contrato, migración y límites](releases/v1.0.0.md). La versión publicada sigue siendo v0.2.0; **no existe publicación v1.0.0**.

La PR hacia `main` registrará la checklist de prepublicación: alcance, contrato, regresión/seguridad, instalación/actualización, bytes/avisos/checksums y mantenimiento sobre el SHA exacto. La casilla de aprobación de ROADMAP queda pendiente para #41. No confundir preparación correcta con publicación autorizada ni completar aquella casilla anticipadamente.

Para reproducir un paquete de revisión: guardar el bloque PowerShell de la sección siguiente en `dist/package-local.ps1`, seleccionar el SHA exacto del candidato registrado en su PR y ejecutar desde la raíz con Go/motor fijados. Ejemplo (sustituir únicamente el SHA por los 40 caracteres comprobados):

```powershell
$candidateSHA = 'SHA_DE_40_CARACTERES_DE_LA_PR'
if ($candidateSHA -notmatch '^[0-9a-fA-F]{40}$') { throw 'Set the reviewed candidate SHA first' }
./dist/package-local.ps1 -packageVersion '1.0.0' -sourceRef $candidateSHA
```

El helper es una copia literal del documento, no un requisito privado sin versionar. Genera todo en directorios nuevos, no publica nada. La evidencia e inventario/checksums finales se registran en la PR y BUILD para evitar referencias circulares al propio commit. Las verificaciones de UI de #23 son históricas y reutilizables porque las plantillas/snapshots no cambian. Verificar ambos diffs (`develop...release/1.0.0` para preparación y `main...release/1.0.0` para todo lo acumulado) antes de aprobar. La integración/publicación y PR posterior a `develop` requieren autorizaciones propias.

## Proceso para una próxima versión

El [roadmap](ROADMAP.md) describe propuestas, no una publicación aprobada. Antes de una nueva release: revisar alcance/compatibilidad, fijar versión y notas, crear `release/<version>` desde `develop` aprobado, verificar su SHA exacto con tests/CI/smoke y revisar los archivos que se distribuirían. La licencia propia ya es Apache-2.0; volver a revisar avisos de terceros solo según los componentes distribuidos.

Preparar PR hacia `main` y solicitar autorización explícita de fusión/publicación. Solo después crear tag y release sobre el commit aprobado y preparar otra PR de sincronización a `develop`. Sin auto-merge, force push ni borrado de ramas. Esta iteración no crea una release ni modifica tags existentes.

## Procedimiento reproducible de empaquetado

El perfil predeterminado del bloque sigue siendo el **ensayo local sobre v0.2.0 publicada** comprobado en #39; no sustituye su ZIP oficial. #40 reutiliza la misma receta con los parámetros v1.0.0/SHA del candidato y notas correspondientes. #41 exige autorización explícita de publicación. No ejecutar comandos de tag/release como parte del procedimiento.

Desde la raíz, con Go 1.27.4 en PATH y oasdiff 1.39.0 local según INSTALL. Se extrae el SHA en dos carpetas nuevas, se construye el mismo ejecutable y se empaqueta una lista explícita de **20 archivos**. LICENSE/NOTICE/THIRD_PARTY_NOTICES permanecen íntegros; el motor y sus dependencias no están enlazados ni incluidos. Los archivos propios usados son fixtures sintéticos, no informes reales.

```powershell
param(
    [string]$packageVersion = '0.2.0',
    [string]$sourceRef = 'v0.2.0'
)
$ErrorActionPreference = 'Stop'
if ($packageVersion -notmatch '^\d+\.\d+\.\d+$') { throw 'Expected a numeric release version' }
$repoRoot = (Get-Location).Path
$sourceCommit = (& git rev-parse --verify --end-of-options ($sourceRef + '^{commit}')).Trim()
if ($LASTEXITCODE -ne 0 -or $sourceCommit -notmatch '^[0-9a-fA-F]{40}$') { throw 'Expected an exact source commit' }
$sourceTree = (& git rev-parse ($sourceCommit + '^{tree}')).Trim()
if ($LASTEXITCODE -ne 0) { throw 'Source tree unavailable' }
$env:GOTOOLCHAIN = 'local'; $env:GOPROXY = 'off'; $env:GOSUMDB = 'off'
$env:GOENV = 'off'; $env:GOTELEMETRY = 'off'; $env:GOWORK = 'off'
$env:GOCACHE = Join-Path $repoRoot 'dist/go-cache'
$env:GOMODCACHE = Join-Path $repoRoot 'dist/go-mod-cache'
$env:CGO_ENABLED = '0'; $env:GOOS = 'windows'; $env:GOARCH = 'amd64'
$env:GOAMD64 = 'v1'; $env:GOFLAGS = ''; $env:GOEXPERIMENT = ''
$goVersion = (& go version).Trim()
if ($LASTEXITCODE -ne 0 -or $goVersion -ne 'go version go1.27.4 windows/amd64') { throw 'Go 1.27.4 Windows x64 required' }
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
$work = Join-Path $repoRoot ('dist/_tools/packaging-' + $packageVersion + '-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $work | Out-Null
$sourceZip = Join-Path $work 'source.zip'
& git archive --format=zip ('--output=' + $sourceZip) $sourceCommit
if ($LASTEXITCODE -ne 0) { throw 'Source archive failed' }
$sourceA = Join-Path $work 'source-a'; $sourceB = Join-Path $work 'source-b'
Expand-Archive -LiteralPath $sourceZip -DestinationPath $sourceA
Expand-Archive -LiteralPath $sourceZip -DestinationPath $sourceB
$binaryA = Join-Path $work 'binary-a.exe'; $binaryB = Join-Path $work 'binary-b.exe'
Push-Location $sourceA
try {
    ./scripts/verify.ps1
    if ($LASTEXITCODE -ne 0) { throw 'Selected source verification failed' }
    & go build -trimpath -buildvcs=false -o $binaryA ./cmd/openapi-contract-lab
    if ($LASTEXITCODE -ne 0) { throw 'First build failed' }
} finally { Pop-Location }
Push-Location $sourceB
try {
    & go build -trimpath -buildvcs=false -o $binaryB ./cmd/openapi-contract-lab
    if ($LASTEXITCODE -ne 0) { throw 'Second build failed' }
} finally { Pop-Location }
$binaryHash = (Get-FileHash $binaryA -Algorithm SHA256).Hash.ToLowerInvariant()
if ($binaryHash -ne (Get-FileHash $binaryB -Algorithm SHA256).Hash.ToLowerInvariant()) { throw 'Builds differ' }
if ((& $binaryA --version).Trim() -ne ('openapi-contract-lab v' + $packageVersion) -or $LASTEXITCODE -ne 0) { throw 'Version mismatch' }
& go version -m $binaryA | Out-File (Join-Path $work 'build-metadata.txt') -Encoding utf8
if ($LASTEXITCODE -ne 0) { throw 'Build metadata failed' }
$stage = Join-Path $work 'package'
New-Item -ItemType Directory -Path $stage | Out-Null
$files = [ordered]@{
    'LICENSE' = 'LICENSE'; 'NOTICE' = 'NOTICE'; 'THIRD_PARTY_NOTICES.txt' = 'THIRD_PARTY_NOTICES.txt'
    'RELEASE_NOTES.md' = ('docs/releases/v' + $packageVersion + '.md')
    'examples/base.yaml' = 'spikes/oasdiff/fixtures/base.yaml'
    'examples/required.yaml' = 'spikes/oasdiff/fixtures/required.yaml'
    'examples/description.yaml' = 'spikes/oasdiff/fixtures/description.yaml'
    'examples/reports/grouped.html' = 'docs/examples/issue-23-grouped.html'
    'examples/reports/grouped.md' = 'docs/examples/issue-23-grouped.md'
}
foreach ($name in @('parameter-base','parameter-head','body-base','body-head','response-base','response-head')) {
    $files['examples/breaking/' + $name + '.json'] = 'internal/explain/testdata/breaking/' + $name + '.json'
}
foreach ($name in @('base','union-null','response-property')) {
    $files['examples/openapi31/' + $name + '.json'] = 'internal/cli/testdata/openapi31/' + $name + '.json'
}
foreach ($entry in $files.GetEnumerator()) {
    $target = Join-Path $stage $entry.Key
    New-Item -ItemType Directory -Force -Path (Split-Path $target -Parent) | Out-Null
    Copy-Item -LiteralPath (Join-Path $sourceA $entry.Value) -Destination $target
}
Copy-Item -LiteralPath $binaryA -Destination (Join-Path $stage 'openapi-contract-lab.exe')
$inventory = @($files.Keys) + @('BUILD.txt','openapi-contract-lab.exe') | Sort-Object
$buildText = @(
    ('OpenAPI Contract Lab v' + $packageVersion + ' - local review package; not published by this procedure'),
    'Source: https://github.com/v1kz25/openapi-contract-lab',
    ('SourceCommit: ' + $sourceCommit), ('SourceTree: ' + $sourceTree),
    ('Toolchain: ' + $goVersion), 'ExternalEngine: oasdiff 1.39.0 (not included)',
    'CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOAMD64=v1 GOFLAGS= GOEXPERIMENT= GOWORK=off',
    'Build: go build -trimpath -buildvcs=false -o openapi-contract-lab.exe ./cmd/openapi-contract-lab',
    ('BinarySHA256: ' + $binaryHash), 'Two source directories: identical executable bytes.',
    'Windows x64 only; contract/status in release notes; limited OpenAPI corpus.',
    'LICENSE/NOTICE: Apache-2.0; Go BSD-3-Clause in THIRD_PARTY_NOTICES.txt.',
    'Merge, tag and publication require independent explicit approval.',
    ('Files: ' + ($inventory -join ', '))
) -join "`n"
[IO.File]::WriteAllText((Join-Path $stage 'BUILD.txt'), $buildText + "`n", [Text.UTF8Encoding]::new($false))
$zipName = 'openapi-contract-lab_' + $packageVersion + '_windows_amd64.zip'
$zipPath = Join-Path $work $zipName
Compress-Archive -Path (Join-Path $stage '*') -DestinationPath $zipPath -CompressionLevel Optimal
Add-Type -AssemblyName System.IO.Compression.FileSystem
$archive = [IO.Compression.ZipFile]::OpenRead($zipPath)
try {
    $actual = @($archive.Entries | Where-Object Name | ForEach-Object { $_.FullName.Replace('\','/') } | Sort-Object)
    if (@(Compare-Object $inventory $actual).Count -ne 0) { throw 'Unexpected ZIP inventory' }
} finally { $archive.Dispose() }
$zipHash = (Get-FileHash $zipPath -Algorithm SHA256).Hash.ToLowerInvariant()
[IO.File]::WriteAllText((Join-Path $work 'SHA256SUMS.txt'), $zipHash + '  ' + $zipName + "`n", [Text.Encoding]::ASCII)
$extracted = Join-Path $work 'extracted'
Expand-Archive -LiteralPath $zipPath -DestinationPath $extracted
foreach ($name in $inventory) {
    if ((Get-FileHash (Join-Path $stage $name) -Algorithm SHA256).Hash -ne (Get-FileHash (Join-Path $extracted $name) -Algorithm SHA256).Hash) { throw ('Extracted file changed: ' + $name) }
}
Push-Location $sourceA
try {
    ./scripts/smoke.ps1 -Binary (Join-Path $extracted 'openapi-contract-lab.exe') -Engine $env:OASDIFF_BINARY
    if ($LASTEXITCODE -ne 0) { throw 'Extracted CLI smoke failed' }
    powershell.exe -NoProfile -ExecutionPolicy Bypass -File ./scripts/smoke.ps1 -Binary (Join-Path $extracted 'openapi-contract-lab.exe') -Engine $env:OASDIFF_BINARY
    if ($LASTEXITCODE -ne 0) { throw 'Extracted CLI PS5.1 smoke failed' }
} finally { Pop-Location }
```

Registrar SHA fuente/tree, versión Go/motor, flags, SHA del ejecutable/ZIP e inventario y resultados. Este ensayo no garantiza ZIP idéntico entre fechas/herramientas de compresión: sus hashes se calculan sobre bytes reales. No utilizar su ZIP como si fuera el oficial; descargar este desde INSTALL. Antes de una release futura, repetir ejemplos del paquete, inspeccionar BUILD/avisos/inventario y privacidad de fuentes/logs. No empaquetar SDK, oasdiff, cachés, helpers de test, informes de contratos reales ni credenciales. Los hashes no son firma digital.

## Registro histórico de preparación del issue #8

El resto de este documento conserva la evidencia y la lista de decisiones pendientes **en aquel momento**, antes de #18. Las referencias a `dev` o licencia/publicación pendiente son históricas, no el estado actual.

Este issue prepara instrucciones y comprobaciones. **No publica una versión estable, no crea tags, no cambia la licencia ni distribuye binarios.** La CLI conserva `openapi-contract-lab dev`. [Instalación y ejemplos](INSTALL.md) describen funciones disponibles, sin anunciar un instalador ni una GUI.

## Evidencia real — 2026-11-21

Base de código: `fbc97e61bfc1e7ceec46772a2a60a2a9f5c0dbc4` de `develop`, con los scripts/documentación de preparación del issue #8; el SHA completo de la revisión se registra en su PR. No hay cambios en el código Go de producción.

- Host Windows x64, build **26300.9974** (API Windows NT 10.0.26300); objetivo Windows 11. No se infiere otra edición a partir de la API de compatibilidad.
- Go **1.27.4**, oasdiff **1.39.0**; hashes de los archivos oficiales comprobados y coincidentes con [INSTALL.md](INSTALL.md). Herramientas portátiles existentes; no se simuló una instalación limpia del sistema operativo ni se repitió la descarga del SDK en este issue.
- `scripts/verify.ps1`: PASS con `OASDIFF_BINARY`, tests reales sin caché, formato, `go vet` y compilación. Sin descargas de módulos ni toolchains.
- `scripts/smoke.ps1` sobre el ejecutable compilado, en PowerShell **5.1.26100.9818** y **7.6.9**: ayuda/versión/compare-help 0; breaking 3; compatible/documental/none 0; exportaciones JSON/Markdown/HTML 3 con archivo completo UTF-8 y stdout vacío; reintentos de escritura 1 y SHA-256 sin cambios; política inválida 2; entrada ausente 1. Acentos JSON y UTF-8 comprobados, incluidas rutas de salida con espacios. Logs locales en `dist/smoke-*`, sin publicarlos.
- Dos compilaciones con código copiado en directorios distintos, Go/GOOS/GOARCH/GOAMD64/CGO/flags fijados según INSTALL: **bytes idénticos**, SHA-256 `a18cde9bf7a8618bea985ea6a0a2436ecca60ddb4cff45833222331e831dc72c`. Es evidencia local de esta revisión, no un checksum de release ni una garantía para cualquier entorno. `go version -m` confirmó Go, trimpath, Windows/amd64, GOAMD64=v1 y CGO deshabilitado.
- CI añade el smoke tras compilar. Su ejecución real y revisión del diff se registran en la PR; no se confunde la prueba local con el runner alojado.

Los smoke automatizan comandos del ejecutable real: no hay que pedir al usuario que los repita. La revisión visual del HTML de breaking/vacío ya fue aportada por el usuario durante el issue #5 (PR #14); este issue no cambia su presentación ni afirma nuevas comprobaciones visuales. No se verificó ARM64, Linux/macOS ni una instalación en otra máquina.

## Límites que deben acompañar una publicación

- oasdiff 1.39.0 es el único motor semántico y debe obtenerse por separado. Referencias externas, URLs/UNC y servicios remotos deshabilitados. Timeout de comparación: 30 segundos.
- Fixtures mínimos 3.0.3 y 3.1.0; cobertura limitada, sin validación completa de 3.1 ni soporte verificado de 3.2. Diez explicaciones específicas y fallback neutral para otras combinaciones.
- Un informe sin hallazgos no implica contratos idénticos. Códigos de salida propios 0/1/2/3; JSON/CLI en desarrollo. Rutas locales forman parte de los informes y deben revisarse antes de compartirlos.
- Coste obligatorio 0 €; ejecución local y CI estándar gratuita en repositorio público. Sin APIs de IA, cuentas comerciales ni secretos nuevos. Go sin módulos externos; oasdiff Apache-2.0 y Go BSD-3-Clause ya documentados.

## Revisión antes de tag/release

1. Aprobar la PR de preparación hacia `develop`. No equivale a aprobar la publicación.
2. Resolver la **licencia propia pendiente** con decisión explícita del propietario. Si se redistribuyera oasdiff u otro componente, revisar además inventario, licencias y avisos de los binarios incluidos; esta preparación no los empaqueta.
3. Acordar versión, alcance estable y notas de release después de revisar las limitaciones. No se reserva ni crea un tag en este issue. Si se solicita que `--version` identifique una versión publicada, preparar y verificar ese cambio antes de publicar; hoy devuelve `dev`.
4. Crear `release/<version>` desde el `develop` aprobado para estabilizar según Gitflow; ejecutar CI y smoke sobre su SHA exacto. PR hacia `main`, con revisión y autorización explícita para fusionar; después otra PR de sincronización hacia `develop`. No borrar ramas.
5. Solo tras esa revisión y aprobación explícita de publicación, definir el tag sobre el commit aprobado de `main` y la release correspondiente. Si se acuerdan artefactos gratuitos, fijar fuente/toolchain/flags, generar checksums y avisos, y probar exactamente los archivos que se publicarán. No usar runners facturables ni habilitar auto-merge.

La licencia y las decisiones finales de versión/publicación son pendientes de la futura release, no motivos para inventar una licencia, tag o paquete en esta tarea de preparación. El issue termina con documentación, smoke verificados y PR para revisión.
