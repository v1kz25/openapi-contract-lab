# Instalación y actualización — Windows x64

Versión vigente: [**v0.2.0 experimental publicada**](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.2.0), no estable. Solo Windows x64, PowerShell 5.1 y 7 verificados. Go 1.27.4 solo es necesario para **compilar desde fuente**; para ejecutar basta la CLI del ZIP y oasdiff 1.39.0 externo. No hay instalador, autoupdater ni instalación global.

El [candidato v1.0.0 de #40](releases/v1.0.0.md) aún no está publicado: usar únicamente sus archivos locales y SHA/checksums registrados en la PR para revisarlo. No cambiar las URLs/hashes de descarga de v0.2.0 por una v1.0.0 inexistente. [Preparación reproducible del candidato](RELEASE.md#candidato-v100--issue-40); las instrucciones siguientes siguen reproduciendo v0.2.0 publicada.

Se conservan las notas y bytes históricos de v0.1.0/v0.2.0; las frases de candidato pendiente dentro del ZIP v0.2.0 describen su preparación anterior a publicar. El estado vigente se comprueba en la release y aquí, no reescribiendo archivos publicados. Licencia Apache-2.0; conservar LICENSE, NOTICE y THIRD_PARTY_NOTICES.txt (runtime/stdlib Go BSD-3-Clause). oasdiff Apache-2.0 se descarga por separado, conservando sus avisos; no se redistribuye dentro de nuestro ZIP.

## Usar el ZIP sin SDK

Requisitos de ejecución: PowerShell 5.1 o 7 x64 y `tar` de Windows para extraer el motor. No se necesita Git ni Go para este bloque. Ejecutarlo en una carpeta donde se pueda escribir; crea un directorio nuevo, sin sobrescribir una instalación anterior. Solo las descargas iniciales necesitan Internet; ayuda, comparación e informes funcionan localmente sin credenciales ni coste obligatorio. Los checksums verifican bytes, no son una firma digital.

```powershell
$ErrorActionPreference = 'Stop'
$installRoot = Join-Path (Get-Location) ('ocl-install-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $installRoot | Out-Null
$zipName = 'openapi-contract-lab_0.2.0_windows_amd64.zip'
$zipPath = Join-Path $installRoot $zipName
$sumsPath = Join-Path $installRoot 'SHA256SUMS.txt'
$releaseBase = 'https://github.com/v1kz25/openapi-contract-lab/releases/download/v0.2.0/'
Invoke-WebRequest -UseBasicParsing ($releaseBase + $zipName) -OutFile $zipPath
Invoke-WebRequest -UseBasicParsing ($releaseBase + 'SHA256SUMS.txt') -OutFile $sumsPath
$sumLines = [regex]::Matches([IO.File]::ReadAllText($sumsPath), '(?m)^([0-9a-fA-F]{64})  openapi-contract-lab_0\.2\.0_windows_amd64\.zip\r?$')
$zipHash = (Get-FileHash -LiteralPath $zipPath -Algorithm SHA256).Hash
if ($sumLines.Count -ne 1 -or $zipHash -ne $sumLines[0].Groups[1].Value -or $zipHash -ne '6c2bfd76f4c123399029e908bb58f73a40b073fe5f19b7eb4069e1af9aff4c81') { throw 'CLI checksum mismatch' }
$cliFolder = Join-Path $installRoot 'v0.2.0'
Expand-Archive -LiteralPath $zipPath -DestinationPath $cliFolder
$cli = Join-Path $cliFolder 'openapi-contract-lab.exe'
$exeHash = (Get-FileHash -LiteralPath $cli -Algorithm SHA256).Hash
if ($exeHash -ne 'a5eb2523bce835e8a7ab33796c0baa78bed0e9406796b894fff0c3658f1a15f7') { throw 'CLI executable checksum mismatch' }
$motorFolder = Join-Path $installRoot 'oasdiff-1.39.0'
New-Item -ItemType Directory -Path $motorFolder | Out-Null
$motorArchive = Join-Path $installRoot 'oasdiff.tar.gz'
Invoke-WebRequest -UseBasicParsing 'https://github.com/oasdiff/oasdiff/releases/download/v1.39.0/oasdiff_1.39.0_windows_amd64.tar.gz' -OutFile $motorArchive
if ((Get-FileHash -LiteralPath $motorArchive -Algorithm SHA256).Hash -ne '956a73df88e03efdf3fcc5f327446cc176fefe4adcea1978b0fef464a1ee8901') { throw 'oasdiff checksum mismatch' }
tar -xzf $motorArchive -C $motorFolder
if ($LASTEXITCODE -ne 0) { throw 'oasdiff extraction failed' }
$motor = Join-Path $motorFolder 'oasdiff.exe'
& $cli --help
if ($LASTEXITCODE -ne 0) { throw 'CLI help failed' }
& $cli --version
if ($LASTEXITCODE -ne 0) { throw 'CLI version failed' }
& $motor --version
if ($LASTEXITCODE -ne 0) { throw 'oasdiff version failed' }
```

Esperado: `openapi-contract-lab v0.2.0` y `oasdiff version 1.39.0`. `BUILD.txt` incluido identifica fuente `f6e7f68ac31d81b8e23cff0e18b41b0e19197c35`; su árbol coincide con el tag publicado, aunque el commit del tag sea el merge `72ebbd87558bd73ae138696fd9dc62a58b6e2c04`.

```powershell
& $cli compare --oasdiff $motor (Join-Path $cliFolder 'examples/base.yaml') (Join-Path $cliFolder 'examples/required.yaml')
if ($LASTEXITCODE -ne 3) { throw 'Expected completed breaking report, exit 3' }
& $cli compare --oasdiff $motor (Join-Path $cliFolder 'examples/required.yaml') (Join-Path $cliFolder 'examples/base.yaml')
if ($LASTEXITCODE -ne 0) { throw 'Expected compatible report, exit 0' }
$reportPath = Join-Path $installRoot 'revision.html'
& $cli compare --oasdiff $motor --format html --output $reportPath (Join-Path $cliFolder 'examples/base.yaml') (Join-Path $cliFolder 'examples/required.yaml')
if ($LASTEXITCODE -ne 3) { throw 'Expected HTML generated before the gate' }
```

Abrir `revision.html` localmente si se desea; no hay servidor. El ZIP incluye también ejemplos `breaking/`, `openapi31/` y dos informes agrupados sintéticos. Los enlaces `file:///synthetic/` de estos últimos son ilustrativos, no fuentes incluidas.

## Actualizar y volver a la versión anterior

Antes de cambiar una instalación que ya funciona, guardar sus rutas explícitas:

```powershell
$previousCli = $cli
$previousMotor = $motor
```

Volver a ejecutar el bloque de instalación con la versión **publicada y verificada** que se quiera seleccionar; los valores actuales son exclusivamente para v0.2.0, no sustituir el nombre de versión manteniendo sus hashes. Cada ejecución crea una carpeta nueva. Revisar notas y [compatibilidad CLI/JSON](DEVELOPMENT.md#compatibilidad-del-contrato-v1), verificar hashes/versión y probar ejemplos antes de usar contratos propios. No copiar el ejecutable sobre la carpeta antigua, no mover informes ni contratos, no modificar PATH global ni borrar la versión previa. `$cli`/`$motor` seleccionan los nuevos binarios solo en esta sesión.

Para revertir, usar las rutas guardadas y comprobar la versión:

```powershell
$cli = $previousCli
$motor = $previousMotor
& $cli --version
if ($LASTEXITCODE -ne 0) { throw 'Previous CLI unavailable' }
```

La reversión no convierte informes nuevos a un esquema antiguo ni deshace cambios de contratos. Conservar versión/procedencia de cada informe; v0.1/v0.2 son experimentales. Si hay error de escritura, elegir otro nombre según [recuperación de informes](REPORTS.md#fallos-de-escritura-y-recuperación--issue-38).

## Compilar desde fuente

Solo esta alternativa requiere Git, Go 1.27.4 y espacio para el SDK portátil. Los pasos siguientes reproducen **v0.2.0 publicada**, no los cambios posteriores de `develop` ni una futura v1 estable. Para empaquetar desde un SHA concreto, seguir [RELEASE.md](RELEASE.md#procedimiento-reproducible-de-empaquetado).

## 1. Código y carpeta de trabajo

En una carpeta de trabajo nueva, no en System32:

```powershell
git clone --branch v0.2.0 https://github.com/v1kz25/openapi-contract-lab.git
Set-Location .\openapi-contract-lab
git remote get-url origin
git rev-parse HEAD
```

Guardar el SHA mostrado junto con la evidencia de compilación. El tag publicado `v0.2.0` identifica esa versión en detached HEAD; no se trabaja sobre él. Para reproducir otra revisión concreta, usar su SHA aprobado en un checkout limpio (`git switch --detach SHA_APROBADO`) y comprobarlo con `git rev-parse HEAD`. Si ya se tiene el repositorio, no volver a clonarlo ni sobrescribir cambios. Para desarrollar, seguir Gitflow desde `develop` según [DEVELOPMENT.md](DEVELOPMENT.md).

## 2. Go y oasdiff fijados

Ejecutar una vez en un checkout nuevo. Si las carpetas de herramientas ya existen, verificar sus versiones y hashes antes de reutilizarlas; no repetir la extracción sobre una instalación ajena.

```powershell
$ErrorActionPreference = 'Stop'
New-Item -ItemType Directory -Force dist | Out-Null
Invoke-WebRequest 'https://go.dev/dl/go1.27.4.windows-amd64.zip' -OutFile dist/go1.27.4.windows-amd64.zip
$sdkHash = (Get-FileHash dist/go1.27.4.windows-amd64.zip -Algorithm SHA256).Hash
if ($sdkHash -ne '178ab98950ac69ab25069307dc6ade79f3816e414cce95b4da166a115cde1e7e') { throw 'Go checksum mismatch' }
Expand-Archive dist/go1.27.4.windows-amd64.zip -DestinationPath dist/_tools/go1.27.4
$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path

New-Item -ItemType Directory dist/oasdiff-1.39.0 | Out-Null
Invoke-WebRequest 'https://github.com/oasdiff/oasdiff/releases/download/v1.39.0/oasdiff_1.39.0_windows_amd64.tar.gz' -OutFile dist/oasdiff-1.39.0/oasdiff.tar.gz
$engineHash = (Get-FileHash dist/oasdiff-1.39.0/oasdiff.tar.gz -Algorithm SHA256).Hash
if ($engineHash -ne '956a73df88e03efdf3fcc5f327446cc176fefe4adcea1978b0fef464a1ee8901') { throw 'oasdiff checksum mismatch' }
tar -xzf dist/oasdiff-1.39.0/oasdiff.tar.gz -C dist/oasdiff-1.39.0
if ($LASTEXITCODE -ne 0) { throw 'Extraction failed' }
go version
.\dist\oasdiff-1.39.0\oasdiff.exe --version
```

Esperado: `go version go1.27.4 windows/amd64` y `oasdiff version 1.39.0`. Fuentes: [Go, versiones y SHA-256](https://go.dev/dl/?mode=json), [release y checksums oasdiff](https://github.com/oasdiff/oasdiff/releases/tag/v1.39.0). Conservar las licencias extraídas (Go BSD-3-Clause y oasdiff Apache-2.0). El SDK queda bajo `_tools` para que Go no recorra sus tests. No se modifica el PATH global.

## 3. Comprobar y compilar offline

```powershell
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
$env:GOMODCACHE = Join-Path (Get-Location) 'dist/go-mod-cache'
$env:GOENV = 'off'
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
.\scripts\verify.ps1

$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:GOAMD64 = 'v1'
$env:GOFLAGS = ''
$env:GOEXPERIMENT = ''
$env:GOWORK = 'off'
go build -trimpath -buildvcs=false -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
if ($LASTEXITCODE -ne 0) { throw 'Build failed' }
Get-FileHash dist/openapi-contract-lab.exe -Algorithm SHA256
.\scripts\smoke.ps1
```

`verify.ps1` ejecuta tests con el motor real, análisis y compilación; el último `go build` prepara el binario con flags explícitos. Según [Go build](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies), `-trimpath` elimina rutas del sistema de archivos y `-buildvcs=false` omite los metadatos Git; por eso hay que conservar el SHA fuente por separado. No se promete el mismo hash entre versiones, plataformas o flags distintos. Los scripts devuelven 0 si terminan correctamente; un fallo detiene la comprobación.

## 4. Comparar y abrir un informe

```powershell
.\dist\openapi-contract-lab.exe --help
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
$LASTEXITCODE # 3: JSON completo y breaking change
$destino = '.\dist\revision-' + [guid]::NewGuid().ToString('N') + '.html'
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe --format html --output $destino .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
$LASTEXITCODE # 3: HTML generado antes del gate
Invoke-Item $destino
```

Para un cambio compatible, invertir `required.yaml` y `base.yaml`: exit 0 y regla `request-parameter-became-optional`. Para el caso documental, sustituir `required.yaml` por `description.yaml`: exit 0, cero hallazgos, sin afirmar que los contratos sean idénticos. Ver [política](SEVERITY-POLICY.md) para `--fail-on none`, y [formatos](REPORTS.md) para Markdown/JSON. Flags antes de BASE/HEAD; `--output` exige un archivo nuevo y evita problemas de codificación de `>` en PowerShell 5.1.

No se debe compartir un informe real sin revisar su contenido: incluye rutas absolutas y datos de los contratos. Los ejemplos anteriores solo usan fixtures sintéticos. El SDK no es necesario para ejecutar el binario ya construido; sí hace falta oasdiff 1.39.0 local para comparar.

El ámbito verificado es OpenAPI 3.0.3 y [cinco pares concretos 3.1.0](OPENAPI31-REGRESSION.md), no cobertura completa 3.1 ni soporte verificado 3.2. Son archivos locales revisados; sin URL/UNC/stdin, referencias externas, cuota de RAM ni sandbox. Ver [límites del motor](OASDIFF-ADAPTER.md#límites-y-alcance-de-seguridad--issue-38). La aplicación no sube informes automáticamente; mensajes/eventos también pueden contener datos del contrato. Para reportar un problema, usar [mantenimiento](DEVELOPMENT.md#mantenimiento-y-reporte-de-problemas) con reproducción sintética, nunca subir especificaciones privadas por defecto.
