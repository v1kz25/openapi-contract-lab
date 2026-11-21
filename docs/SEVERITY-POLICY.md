# Política de severidad — issue #6

Alcance: [issue #6](https://github.com/v1kz25/openapi-contract-lab/issues/6). Gate local por severidad normalizada por el adaptador, sin cambiar reglas, clasificación, explicaciones ni contenido del informe de oasdiff. Go 1.27.4 y oasdiff 1.39.0 permanecen fijados. Sin nuevas dependencias; biblioteca estándar Go (BSD-3-Clause) y binario oasdiff (Apache-2.0) ya documentados. Coste adicional obligatorio: 0 €.

## Configuración

`compare --fail-on error|warning|info|none`, antes de BASE/HEAD. Predeterminado: `error`. Solo se aceptan esos valores exactos; una cadena vacía o desconocida es un error de argumentos antes de ejecutar el motor. Sin archivo de configuración ni variables de entorno adicionales.

| `--fail-on` | Severidades que bloquean |
| --- | --- |
| `error` | error y desconocidas |
| `warning` | warning, error y desconocidas |
| `info` | info, warning, error y desconocidas |
| `none` | Ninguna |

Los niveles conocidos siguen siendo los del adaptador: 3 = error/breaking, 2 = warning/potentially-breaking, 1 = info/non-breaking. Un informe vacío supera cualquier política. El umbral `info` bloquea también cambios compatibles: es una elección explícita del usuario. No se inventan reglas ni se deduce compatibilidad del texto de una explicación. Un ID nuevo con nivel conocido conserva su clasificación; una severidad desconocida bloquea por precaución, sin reclasificar el hallazgo como breaking.

## Códigos de salida

| Código | Significado |
| --- | --- |
| 0 | Ayuda/versión o comparación que supera la política |
| 1 | Error de entrada, motor, renderizado o escritura |
| 2 | Argumentos o configuración CLI inválidos |
| 3 | Comparación e informe completados; política no superada |

Primero se completa el informe JSON/Markdown/HTML y su escritura; después se aplica el gate. Con código 3 el informe sigue disponible y se emite un diagnóstico `policy:` en stderr. stdout contiene solo el informe, o queda vacío al usar `--output`. Los fallos operativos tienen prioridad sobre el gate, incluso con `--fail-on none`; no hay garantía de informe completo con código 1.

**Cambio respecto a las versiones anteriores:** un breaking change ahora devuelve 3 por defecto, en lugar de 0. Para conservar el comportamiento anterior, usar `--fail-on none`. No utilizar el código de `go run` para distinguir códigos de la aplicación: ejecutar el binario compilado.

## Ejemplos reproducibles en PowerShell

Desde la raíz del repositorio, con los binarios locales preparados:

```powershell
$motor = '.\dist\oasdiff-1.39.0\oasdiff.exe'
$base = '.\spikes\oasdiff\fixtures\base.yaml'
$required = '.\spikes\oasdiff\fixtures\required.yaml'
& .\dist\openapi-contract-lab.exe compare --oasdiff $motor $base $required
$LASTEXITCODE # 3: breaking; JSON completo en stdout
& .\dist\openapi-contract-lab.exe compare --oasdiff $motor $required $base
$LASTEXITCODE # 0: parámetro pasa a opcional, info/non-breaking
& .\dist\openapi-contract-lab.exe compare --oasdiff $motor --fail-on info $required $base
$LASTEXITCODE # 3: umbral info solicitado
& .\dist\openapi-contract-lab.exe compare --oasdiff $motor --fail-on none $base $required
$LASTEXITCODE # 0: exportación sin gate, hallazgo conservado
& .\dist\openapi-contract-lab.exe compare --oasdiff $motor $base .\spikes\oasdiff\fixtures\description.yaml
$LASTEXITCODE # 0: cero hallazgos; no implica contratos idénticos
```

En automatización, capturar `$LASTEXITCODE` inmediatamente después de ejecutar el binario y propagarlo con `exit $codigo`. No basta con que se haya creado el archivo: también puede haberse creado con código 3. No se añade aún un workflow; corresponde al issue #7.

## Validación

Tests unitarios: matriz completa de los cuatro umbrales, severidades conocidas/desconocidas, lista vacía/mixta y valores inválidos. Tests de integración optativos con `OASDIFF_BINARY`: breaking predeterminado y cada umbral, cambio compatible, documental e idéntico; igualdad de bytes con/sin gate en los tres formatos; exportación real con código 3; prioridad del fallo de escritura y fallo de motor con `none`.

No hay un flujo visual nuevo. Los comandos y sus códigos se comprueban por el agente antes de entregar la PR; no se requiere repetirlos manualmente. No se instala CI ni se modifica el motor de diff.

### Resultados reales — 2026-11-21

Windows x64, PowerShell 5.1.26100.9818, Go 1.27.4 portátil y oasdiff 1.39.0. Sin descargas durante la verificación (`GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`, `GOCACHE` local en `dist/go-cache`).

- `OASDIFF_BINARY` apunta al motor local y `go test -count=1 ./...`: PASS en CLI, explicaciones, adaptador e informes; sin caché ni omisión de integración.
- `go vet ./...`: exit 0, sin diagnósticos. `go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab`: exit 0. `gofmt -l cmd internal`: sin archivos pendientes tras normalizar los finales de línea locales; sin cambios de contenido en los archivos ajenos al issue.
- Ejecutable real desde PowerShell 5.1: breaking predeterminado = 3; compatible = 0; compatible con umbral info = 3; breaking con none = 0; documental = 0. JSON válido y acentos verificados tras `ConvertFrom-Json`.
- Exportación JSON/Markdown/HTML con breaking: código 3 y archivo generado. Reintento sobre cada archivo: código 1 y SHA-256 sin cambios. Política inválida: 2; ejecutable de motor ausente con `none`: 1.
- `compare --help` documenta opciones y los cuatro códigos. No se probó ARM64, Linux/macOS ni una plataforma CI externa. Los niveles warning/desconocido de la política se cubren con datos sintéticos unitarios; no se afirma haber obtenido esos niveles de un contrato real en esta verificación.

Los artefactos y el script de comprobación PowerShell permanecen en `dist/` ignorado; la evidencia publicada no contiene rutas personales ni datos privados.
