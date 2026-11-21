# CI gratuita y regresión — issue #7

Alcance: [issue #7](https://github.com/v1kz25/openapi-contract-lab/issues/7). `.github/workflows/ci.yml` se ejecuta en PR hacia `develop` y `main`, con un único job en `windows-2025` x64 y límite de 15 minutos. No publica releases, paquetes, comentarios ni artefactos; no realiza despliegues.

## Coste, permisos y dependencias

Repositorio público comprobado el 2026-11-21. GitHub ofrece los runners hospedados **estándar** gratis para repositorios públicos: [runners y etiquetas disponibles](https://docs.github.com/en/actions/reference/runners/github-hosted-runners), [facturación](https://docs.github.com/en/actions/concepts/billing-and-usage). Se usa la etiqueta estándar `windows-2025`, sin runner de mayor tamaño ni servicios de pago. El job solo se habilita si el repositorio es público; si se convierte en privado, se omite y habrá que revisar esta decisión. Coste monetario adicional obligatorio: **0 €**.

El evento es `pull_request`, no `pull_request_target`. Permiso mínimo `contents: read`; no necesita secretos configurados, PAT, tarjeta ni credenciales externas. GitHub proporciona su token efímero para las acciones oficiales; checkout no lo conserva para los tests. No se interpolan títulos/cuerpos de PR en comandos. No hay caché compartida, subida de informes ni lectura de archivos privados. Solo se ejecuta el contenido del checkout en una máquina efímera; las trazas de pruebas y compilación serán públicas, por lo que todos los fixtures deben seguir siendo sintéticos.

| Componente | Versión / anclaje | Licencia y necesidad |
| --- | --- | --- |
| Go | 1.27.4 | BSD-3-Clause, compilador y biblioteca estándar; ya usado por el proyecto |
| oasdiff Windows amd64 | 1.39.0, SHA-256 del archivo `956a73df88e03efdf3fcc5f327446cc176fefe4adcea1978b0fef464a1ee8901` | Apache-2.0, motor local existente |
| actions/checkout | v7.0.1, `3d3c42e5aac5ba805825da76410c181273ba90b1` | MIT, obtener el código de la PR |
| actions/setup-go | v7.1.0, `d8cbb8045d7f43706b9177cceb36c1f91a6644d5` | MIT, preparar Go exacto en el runner |

Fuentes comprobadas: [release y checksums de oasdiff](https://github.com/oasdiff/oasdiff/releases/tag/v1.39.0), [licencia checkout](https://github.com/actions/checkout/blob/3d3c42e5aac5ba805825da76410c181273ba90b1/LICENSE), [licencia setup-go](https://github.com/actions/setup-go/blob/d8cbb8045d7f43706b9177cceb36c1f91a6644d5/LICENSE). Ambas acciones usan Node 24 según sus `action.yml` en esos commits y el runner hospedado prepara su ejecución; la aplicación sigue siendo Go sin módulos externos.

Las descargas iniciales requieren Internet: checkout, preparación de Go y archivo oficial oasdiff. Antes de extraer oasdiff se comprueba el SHA-256 fijado, contrastado con `checksums.txt` de la release. Las pruebas y compilación no descargan módulos ni toolchains (`GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`). No se imprime el entorno completo ni tokens. La aplicación conserva su funcionamiento offline.

## Comprobaciones y corpus

`scripts/verify.ps1` exige Go 1.27.4 y oasdiff 1.39.0, rechaza la ausencia de `OASDIFF_BINARY` y ejecuta:

1. `gofmt -l cmd internal` sin archivos pendientes.
2. `go test -count=1 ./...`, con motor real y sin caché de resultados.
3. `go vet ./...`.
4. `go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab`.

El script falla en cuanto una comprobación falla. `.gitattributes` mantiene los Go con LF para que el formato sea consistente al hacer checkout en Windows.

`internal/cli/testdata/regression/` añade contratos mínimos OpenAPI **3.1.0**, YAML inválido intencionado y `cases.json` con expectativas revisables. El test CLI comprueba cinco casos: opcional a obligatorio (exit 3), obligatorio a opcional (exit 0), descripción (exit 0 y cero hallazgos), idéntico (igual) y YAML inválido (exit 1 sin informe). Para los hallazgos se verifican ID, nivel, severidad, compatibilidad, explicación, operación/ruta y ambas ubicaciones.

Se conservan las pruebas OpenAPI 3.0.3 del spike, diez reglas verificadas, referencias externas deshabilitadas, mapeo de niveles desconocidos, política y snapshots/escape de informes. No se regeneran expectativas automáticamente. No se afirma compatibilidad completa con 3.1 ni se valida 3.2 mediante este corpus limitado.

## Equivalente local

Desde la raíz, con el SDK y oasdiff preparados según [DEVELOPMENT.md](DEVELOPMENT.md):

```powershell
$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
./scripts/verify.ps1
```

No hace falta repetir comandos manualmente para aceptar el issue: el agente ejecuta este bloque antes de entregar la PR. No hay una interfaz nueva que arrancar o revisar.

## Evidencia y límites — 2026-11-21

- Ejecución local Windows x64: `scripts/verify.ps1` PASS con Go 1.27.4 y oasdiff 1.39.0. Tests sin caché en los cuatro paquetes, incluyendo los cinco casos nuevos con motor real; análisis, formato y compilación correctos.
- Primera ejecución local: el control de formato detectó CRLF de un checkout anterior. Tras `gofmt -w cmd internal`, se repitió el script completo con éxito. No se cambiaron contenidos de código de producción.
- Workflow revisado con [actionlint 1.7.12](https://github.com/rhysd/actionlint/releases/tag/v1.7.12), MIT: sin diagnósticos, exit 0. Herramienta local opcional en `dist/_tools/`, fuera de la aplicación y de la CI; ZIP Windows amd64 verificado contra su checksum oficial `6e7241b51e6817ea6a047693d8e6fed13b31819c9a0dd6c5a726e1592d22f6e9`.
- La ejecución hospedada se comprueba y registra en la PR; la validación local no se presenta como una ejecución de GitHub Actions. No se probó ARM64 ni Linux/macOS.

La imagen `windows-2025` recibe actualizaciones aunque su etiqueta esté fijada. Los SHAs y hashes fijados requieren mantenimiento explícito; una descarga inaccesible o un checksum distinto bloquea la CI. El workflow comprueba la PR, pero no configura protecciones de rama. Preparar una release o cambiar la licencia del proyecto queda fuera de este issue.
