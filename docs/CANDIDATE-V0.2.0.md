# Revisar el candidato experimental v0.2.0 — issue #33

> **Registro histórico de preparación.** v0.2.0 ya está [publicada como prerelease](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.2.0), tras #34/#35. El resto conserva el estado de #33 antes de autorizar su publicación; no describe una publicación pendiente actual. Para instalar/actualizar, usar [INSTALL.md](INSTALL.md); para próximas versiones, [RELEASE.md](RELEASE.md).

Estado: **candidato preparado para revisión**, no publicado. Rama `release/0.2.0`, PR hacia `main`; ninguna fusión, tag, release ni draft release está autorizada. v0.1.0 publicada permanece intacta. Ver [notas del candidato](releases/v0.2.0.md) para funciones, ejemplos, códigos y límites.

## Fuente y herramientas

Base revalidada: `develop` en `a048d2545bf1d2819c4450581603b95d297d70f0`. Los issues #21–#25 están completados y sus PRs fusionadas; #29 aporta el logo. El candidato solo cambia identificación de versión, expectativas específicas, cabecera de avisos, documentación y la resolución relativa del runner del spike. No altera motor, explicaciones, agrupación, JSON ni gate.

Para un checkout nuevo, seleccionar la rama del candidato, registrar el SHA y usar las herramientas fijadas según [INSTALL.md](INSTALL.md): Go **1.27.4**, oasdiff **1.39.0**. Los hashes de sus archivos oficiales son, respectivamente, `178ab98950ac69ab25069307dc6ade79f3816e414cce95b4da166a115cde1e7e` y `956a73df88e03efdf3fcc5f327446cc176fefe4adcea1978b0fef464a1ee8901`. No descargar ni extraer sobre herramientas existentes sin verificarlas.

```powershell
git clone --branch release/0.2.0 https://github.com/v1kz25/openapi-contract-lab.git
Set-Location .\openapi-contract-lab
git remote get-url origin
git rev-parse HEAD
```

La rama puede avanzar durante revisión: para reproducir una entrega, usar su SHA exacto en otro checkout, no asumir que el nombre de rama identifica bytes inmutables.

## Verificación offline desde la raíz

```powershell
$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
$env:GOENV = 'off'
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
$env:GOMODCACHE = Join-Path (Get-Location) 'dist/go-mod-cache'
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
./scripts/verify.ps1
./scripts/smoke.ps1
powershell.exe -NoProfile -ExecutionPolicy Bypass -File ./scripts/smoke.ps1
./spikes/change-kind/verify.ps1
```

Verify exige motor real y ejecuta tests sin caché, formato, vet y build. Incluye trece explicaciones, fallback, cinco pares 3.1, agrupación/orden/conservación de hallazgos, JSON/códigos 0/1/2/3, escape/CSP, enlaces y protección de archivos. El runner del spike requiere PowerShell 7: se comprueba además desde otra carpeta, con directorio nativo distinto, rutas relativas/absolutas y un destino existente que debe rechazarse.

La corrección del runner usa explícitamente `repoRoot` como base de `GetFullPath`; no cambia la señal investigada ni implementa un clasificador. El problema original se reproduce antes de corregirlo. No se regeneran snapshots para aceptar el candidato.

## Build reproducible y paquete

Después de verify, compilar con flags explícitos; las cachés anteriores son locales y las descargas siguen deshabilitadas:

```powershell
$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:GOAMD64 = 'v1'
$env:GOFLAGS = ''
$env:GOEXPERIMENT = ''
$env:GOWORK = 'off'
go build -trimpath -buildvcs=false -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
./dist/openapi-contract-lab.exe --version
Get-FileHash ./dist/openapi-contract-lab.exe -Algorithm SHA256
```

La verificación de reproducibilidad compara ese candidato con otra compilación del mismo SHA extraído mediante `git archive` a una carpeta distinta, manteniendo toolchain, flags y variables. Se comprueban bytes/hashes del ejecutable y metadatos de `go version -m`. No se reutilizan hashes del binario v0.1.0 de la auditoría previa, ni se promete un ZIP idéntico entre herramientas de compresión o fechas distintas.

El paquete Windows x64 se construye con una lista explícita: CLI, ejemplos 3.0.3, tres pares nuevos de explicaciones, contratos de unión/null y propiedad de respuesta 3.1, informes sintéticos agrupados, notas, LICENSE, NOTICE, THIRD_PARTY_NOTICES y BUILD.txt. BUILD registra SHA fuente, herramientas/flags y hash de la CLI, sin rutas personales. El ZIP no contiene oasdiff, SDK, cachés, logs ni datos privados.

`SHA256SUMS.txt` identifica los bytes reales del ZIP. El hash del ejecutable también se registra en BUILD.txt y se contrasta tras extraer. Se ejecuta smoke en PowerShell 5.1/7 usando **ese ejecutable extraído**, y se prueban los ejemplos incluidos. Artefactos y evidencias locales se entregan fuera de Git en `dist/`; los valores concretos, SHA y resultado de CI se registran en la PR del candidato, evitando referencias circulares al propio commit.

## Revisión y límites pendientes

La [aceptación de #23](https://github.com/v1kz25/openapi-contract-lab/issues/23#issuecomment-7070896921) documenta revisión visual previa del HTML en `a048d25`. El candidato no modifica el renderer; no se presenta esa auditoría como nueva inspección visual ni como CI del candidato. Tampoco se han probado otros sistemas/arquitecturas ni una instalación en otra máquina.

Antes de decidir integración/publicación: revisar diff frente a `develop` (estabilización) y frente a `main` (contenido acumulado), notas, BUILD, checksums, resultados locales y CI del **head exacto de esta PR**. Los éxitos de PRs anteriores no sustituyen esa CI.

Esta preparación termina con PR abierta y archivos locales para revisión. Solicitar aprobación explícita antes de fusionar a `main`, crear tag/release o publicar archivos. Solo después correspondería otra PR de sincronización a `develop`, también sujeta a aprobación. No habilitar auto-merge ni borrar ramas. Coste adicional obligatorio **0 €**.
