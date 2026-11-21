# Desarrollo — openapi-contract-lab

## Estado

Módulo `github.com/v1kz25/openapi-contract-lab`, entrada en `cmd/openapi-contract-lab` y manejo de comandos en `internal/cli`. El issue #3 añade `internal/oasdiff` y `compare`: [modelo y verificación](OASDIFF-ADAPTER.md). El issue #4 incorpora `internal/explain`: [reglas y explicaciones](EXPLANATIONS.md). El issue #5 añade `internal/report`: [informes locales Markdown/HTML](REPORTS.md). El issue #6 incorpora la [política de severidad](SEVERITY-POLICY.md). El issue #7 añade [CI y regresión](CI.md).

## Entorno propuesto

- Windows 11; Git, GitHub CLI (`gh`), Codex en desktop/VS Code/CLI.
- Go 1.27.4 fijado en `go.mod`, sin módulos externos ni `go.sum` necesario. El adaptador sigue la decisión del [spike oasdiff](OASDIFF-SPIKE.md): binario local 1.39.0, Apache-2.0.
- Sin hosting, cuentas de servicios comerciales, API keys ni infraestructura de pago.

## Flujo de trabajo

1. Leer AGENTS.md, roadmap y el issue activo.
2. Crear `feature/<issue>-<nombre>` desde `develop` para una única tarea (Gitflow).
3. Implementar sólo su alcance y tests; revisar dependencias/licencias.
4. Ejecutar y registrar pruebas; comprobar Git diff y la identidad Git.
5. Proponer Pull Request; no fusionarla sin permiso.
6. Actualizar el README únicamente para funciones realmente disponibles.
7. Incluir en cada PR los pasos de comprobación manual y sus resultados esperados cuando apliquen. Avisar al usuario cuando deba probar un flujo; arrancar la aplicación local para facilitarlo cuando exista dicho flujo. La ayuda/versión inicial puede comprobarse opcionalmente desde terminal.

## Tests y comandos

Desde la raíz, con Go 1.27.4 en PATH. Estas variables solo afectan a la sesión PowerShell; las pruebas no necesitan descargas:

```powershell
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
gofmt -l cmd internal
go test ./...
go vet ./...
go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
./dist/openapi-contract-lab.exe --help
./dist/openapi-contract-lab.exe --version
./dist/openapi-contract-lab.exe diff
$LASTEXITCODE # 2: diff aún no es un comando
```

`gofmt -l` no debe listar archivos. Ayuda/version y comparaciones que superan la política terminan con código 0; argumentos inválidos con 2; fallos de escritura, entrada o motor con 1. Una comparación completada que no supera `--fail-on` devuelve 3, conservando el informe. El adaptador proporciona `*oasdiff.Error` con etapa y causa original.

Para incluir las pruebas con el motor real (no se descarga durante los tests):

```powershell
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
go test ./...
```

Sin esta variable se ejecutan los tests unitarios y golden; `TestIntegration` se omite con un mensaje explícito. Las pruebas reales usan las entradas sintéticas del spike, prueban errores de YAML/referencias externas y comprueban aislamiento de configuración. No es necesario configurar CI para ejecutarlas.

Con Go/gofmt en PATH y `OASDIFF_BINARY` definido, `./scripts/verify.ps1` ejecuta las comprobaciones usadas en CI y rechaza omitir el motor. `.gitattributes` fija LF para los archivos Go; si un árbol anterior ya tiene CRLF, normalizarlo con `gofmt -w cmd internal` antes de comprobar el formato.

El issue #8 añade [instalación reproducible](INSTALL.md), [preparación de publicación](RELEASE.md) y `./scripts/smoke.ps1`. El smoke usa el ejecutable compilado y oasdiff local, sin SDK ni descargas; se ejecuta también en CI después de compilar. No crea tags, releases ni paquetes.

Con la misma variable también se ejecutan los diez casos sintéticos de `internal/explain`, uno por regla documentada, y se comprueba que el motor emite los IDs/niveles esperados y las fuentes enlazadas. Los snapshots de mensajes se revisan en `internal/explain/testdata/`; no hay actualización automática de expectativas.

`internal/report/testdata/` contiene snapshots revisables de HTML y Markdown para un informe mixto y otro vacío. Los tests de escape y enlaces se ejecutan sin motor. Con `OASDIFF_BINARY`, los tests de CLI también exportan JSON/Markdown/HTML reales a archivos UTF-8 en un directorio temporal. La exportación usa creación exclusiva y no sobrescribe archivos existentes.

### SDK portátil en Windows x64

Si no hay Go instalado, obtener el [ZIP oficial Go 1.27.4](https://go.dev/dl/go1.27.4.windows-amd64.zip), cuyo SHA-256 publicado en [go.dev](https://go.dev/dl/?mode=json) es `178ab98950ac69ab25069307dc6ade79f3816e414cce95b4da166a115cde1e7e`. Comprobarlo antes de extraer:

```powershell
New-Item -ItemType Directory -Force dist | Out-Null
Invoke-WebRequest 'https://go.dev/dl/go1.27.4.windows-amd64.zip' -OutFile dist/go1.27.4.windows-amd64.zip
$sdkHash = (Get-FileHash dist/go1.27.4.windows-amd64.zip -Algorithm SHA256).Hash
if ($sdkHash -ne '178ab98950ac69ab25069307dc6ade79f3816e414cce95b4da166a115cde1e7e') { throw 'Go checksum mismatch' }
Expand-Archive dist/go1.27.4.windows-amd64.zip -DestinationPath dist/_tools/go1.27.4
$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path
go version
```

No requiere instalación global, privilegios de administrador ni oasdiff para ayuda/versión. La descarga inicial necesita Internet; compilación y pruebas usan la biblioteca estándar y pueden hacerse offline. El SDK conserva su [licencia BSD-3-Clause](https://github.com/golang/go/blob/go1.27.4/LICENSE) en `dist/`, que permanece ignorado por Git. Coste adicional obligatorio: 0 €. No se añade framework CLI, servicio, CI ni dependencia de terceros.

El SDK se extrae bajo `_tools` porque `go test ./...` y `go vet ./...` ignoran los directorios que empiezan por `_`. `.gitignore` solo controla Git: extraer el SDK directamente en `dist/go1.27.4` haría que Go intentase recorrer también los tests del compilador.

### Validación del issue #2 — 2026-11-21

Ejecutada en Windows x64, build 26300.9974, PowerShell 7.6.9, Go 1.27.4 portátil con checksum comprobado. No se modificó el PATH global. Con `GOTOOLCHAIN=local`, `GOPROXY=off` y `GOSUMDB=off`:

| Comando | Resultado real |
| --- | --- |
| `go list ./...` | Solo los paquetes `cmd/openapi-contract-lab` e `internal/cli` |
| `go test ./...` | PASS, exit 0; `internal/cli` en 0.513 s, entrada main sin tests propios |
| `go vet ./...` | Sin diagnósticos, exit 0 |
| `gofmt -l cmd internal` | Sin archivos pendientes, exit 0 |
| `go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab` | Binario Windows generado, exit 0 |
| Ejecutable con `--help` / `--version` | Ayuda / `openapi-contract-lab dev`, exit 0 |
| Ejecutable con `diff` | `unknown command "diff"; use --help for usage`, stderr, exit 2 esperado |

Los tests cubren las cuatro formas de ayuda (incluida invocación sin argumentos), dos de versión, cinco entradas inválidas y un fallo de escritura simulado. No invocan oasdiff ni servicios. El primer `go test ./...` falló al recorrer el SDK extraído directamente en `dist/`; se movió a `dist/_tools/` y se repitió la verificación completa correctamente. No se oculta ese fallo de preparación. No se probó ARM64 ni Linux/macOS.

## Gitflow

- `main`: versiones estables; rama predeterminada. No trabajar directamente en ella después de la preparación inicial.
- `develop`: integración de la próxima versión; punto de partida de tareas y spikes.
- `feature/<issue>-<nombre>`: una tarea pequeña o spike del issue activo, creada desde `develop`; PR hacia `develop`. No avanzar al siguiente issue.
- `release/<version>`: estabilización de una versión desde `develop`; solo correcciones, pruebas y documentación. PR hacia `main` y sincronización posterior hacia `develop` mediante otra PR.
- `hotfix/<issue>-<nombre>`: corrección urgente de una versión publicada, desde `main`; PR hacia `main` y otra hacia `develop` para propagar la corrección.
- Preparación de ramas y PRs autorizada dentro del issue activo. Toda fusión requiere aprobación explícita del usuario. No habilitar auto-merge.
- No hacer force push, rebase de ramas compartidas, ni eliminar ramas. Si la interfaz ofrece eliminar la rama tras fusionar, desactivar esa opción.
- Usar Git y GitHub CLI estándar; no se necesita instalar una extensión Gitflow ni añadir servicios o pipelines.
- Antes de una PR: revisar diff, pruebas relevantes, identidad Git y ausencia de secretos. Registrar resultados y bloqueos reales.
