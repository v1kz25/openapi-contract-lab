# Instalación desde fuente — Windows x64

Esta guía reproduce la versión publicada **v0.1.0**. Para revisar el candidato aún no publicado **v0.2.0**, usar [CANDIDATE-V0.2.0.md](CANDIDATE-V0.2.0.md); sus herramientas siguen fijadas a Go 1.27.4 y oasdiff 1.39.0, pero su ejecutable muestra `v0.2.0`.

Instrucciones del issue #8, actualizadas para [v0.1.0 experimental publicada](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.1.0), issue #18 completado. `--version` muestra `openapi-contract-lab v0.1.0`. Licencia propia Apache-2.0, autorizada por el propietario. No se empaqueta oasdiff ni existe instalador. La versión publicada tiene ZIP Windows x64 y checksums.

Para usar el ZIP publicado `openapi-contract-lab_0.1.0_windows_amd64.zip`, comprobar su SHA-256 contra `SHA256SUMS.txt` antes de extraerlo en una carpeta nueva. Incluye la CLI, ejemplos y avisos; no requiere SDK para ejecutarse. Descargar oasdiff por separado como indica el paso 2 y usar `--oasdiff` con su ruta. Las notas incluidas describen los comandos. Los pasos siguientes permiten compilar desde fuente.

Requisitos: Windows x64, Git, PowerShell 5.1 o 7, `tar` y espacio para el SDK portátil. Solo obtener código y herramientas necesita Internet. Después, compilación, comparación e informes son locales, sin credenciales ni coste adicional obligatorio. Las descargas son de Go y oasdiff oficiales.

## 1. Código y carpeta de trabajo

En una carpeta de trabajo nueva, no en System32:

```powershell
git clone --branch v0.1.0 https://github.com/v1kz25/openapi-contract-lab.git
Set-Location .\openapi-contract-lab
git remote get-url origin
git rev-parse HEAD
```

Guardar el SHA mostrado junto con la evidencia de compilación. El tag publicado `v0.1.0` identifica esa versión en detached HEAD; no se trabaja sobre él. Para reproducir otra revisión concreta, usar su SHA aprobado en un checkout limpio (`git switch --detach SHA_APROBADO`) y comprobarlo con `git rev-parse HEAD`. Si ya se tiene el repositorio, no volver a clonarlo ni sobrescribir cambios. Para desarrollar, seguir Gitflow desde `develop` según [DEVELOPMENT.md](DEVELOPMENT.md).

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
