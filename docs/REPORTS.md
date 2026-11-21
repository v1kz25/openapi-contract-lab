# Informes locales — issue #5

Alcance: [issue #5](https://github.com/v1kz25/openapi-contract-lab/issues/5), después del adaptador y las explicaciones. Renderizar resumen y lista de hallazgos como Markdown o HTML, offline y con salida reproducible. Go **1.27.4**, oasdiff **1.39.0**; biblioteca estándar, sin nuevas dependencias, servicios ni coste adicional obligatorio.

## Uso

Desde la raíz del repositorio, con el ejecutable compilado y oasdiff preparado según [DEVELOPMENT.md](DEVELOPMENT.md) y el spike:

```powershell
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe --format html --output .\dist\revision.html .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe --format markdown --output .\dist\revision.md .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
```

- `--format`: `json` (predeterminado), `markdown` o `html`. Las opciones van antes de BASE/HEAD.
- `--output`: archivo local **nuevo**, en un directorio que ya exista. La escritura evita sobrescribir cualquier archivo existente, incluidos los contratos. Para repetir la exportación, elegir otro nombre.
- Sin `--output`, se escribe en stdout. Para guardar HTML/Markdown desde PowerShell 5.1, usar `--output`: así no intervienen su decodificación de pipelines ni la codificación de `>`/`Out-File`.
- Archivos UTF-8 sin BOM, saltos LF. JSON mantiene los escapes Unicode compatibles con consolas antiguas. HTML declara `charset=utf-8` y Markdown conserva los caracteres Unicode.
- Desde el issue #6, una exportación completada devuelve 0 si supera la política o 3 si no la supera; en ambos casos se conserva el informe. Con `--output`, stdout permanece vacío. Formato/argumentos inválidos: 2. Fallo de motor, renderizado o escritura: 1. Ver [política](SEVERITY-POLICY.md); usar `--fail-on none` para exportar sin gate.

## Contenido y estabilidad

Ambos formatos incluyen motor/versión, base/head, recuentos por severidad, y cada hallazgo en el orden entregado por oasdiff: rule ID, nivel original, operación/ruta, mensaje, explicación, impacto, acción sugerida y ubicaciones. Fuentes o rutas ausentes se muestran como no reportadas. Los datos completos del evento siguen disponibles en JSON.

Una lista vacía muestra cero hallazgos y recuerda que eso no implica contratos idénticos: el changelog excluye cambios solo de descripción. No se deducen resultados del diff ni se modifica la clasificación del motor.

No hay fechas de generación, IDs aleatorios ni dependencias de la red en el contenido. Para un mismo informe normalizado se generan los mismos bytes; las rutas locales forman parte de los datos y pueden diferir entre máquinas. Se normalizan los finales de línea para evitar diferencias por `autocrlf` en Windows.

## Escape y contenido activo

- HTML usa `html/template` con escape contextual. Ningún texto de contratos, mensajes o explicaciones se interpreta como HTML confiable.
- Markdown escapa HTML y sintaxis inline, y aplana los saltos dentro de campos para evitar inyectar nuevos bloques o enlaces. Los textos se presentan como contenido, no como instrucciones ejecutables.
- Solo las URLs `file:///...` locales, sin host UNC, consulta ni esquema ejecutable, se convierten en enlaces. Se rechazan destinos `javascript:`, `data:`, HTTP(S), UNC y fragmentos ajenos a líneas. Las ubicaciones originales siguen visibles cuando se descarta un enlace.
- El HTML es un documento completo: CSS estático incluido, sin JavaScript, imágenes, fuentes remotas ni librerías externas. Incluye una política CSP que bloquea cargas de contenido por defecto. No inicia un servidor.
- Los enlaces a fuentes se abren por acción del usuario. El fragmento `#L12` identifica la línea; el posicionamiento automático depende del editor/visor.

## Prueba manual

Abrir PowerShell en la carpeta local `openapi-contract-lab` (no System32). Este bloque genera nombres nuevos y sirve en PowerShell 5.1:

```powershell
$prefijo = ".\dist\revision-$(Get-Date -Format yyyyMMdd-HHmmss-fff)"
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe --format html --output "$prefijo.html" .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
$LASTEXITCODE # 3: informe generado; breaking change con política predeterminada
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe --format markdown --output "$prefijo.md" .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
$LASTEXITCODE # 3: informe generado; breaking change con política predeterminada
Invoke-Item "$prefijo.html"
Get-Content "$prefijo.md" -Encoding UTF8
```

Comprobar título, acentos, un hallazgo/error, regla `request-parameter-became-required`, GET `/items`, explicación en español y ambas fuentes en línea 12/columna 11. En el HTML, comprobar que las rutas largas no desborden la página y que se pueda leer también en una ventana estrecha.

Para el informe vacío, sustituir `required.yaml` por `description.yaml` y usar otro nombre: debe mostrar cero hallazgos y el aviso sobre el alcance del changelog. Volver a exportar sobre un archivo existente debe fallar con exit 1 y conservar el archivo anterior.

## Verificación real — 2026-11-21

Entorno Windows x64 build 26300.9974, SDK Go 1.27.4 y oasdiff 1.39.0 verificados, descargas Go deshabilitadas.

- `go test ./...` con `OASDIFF_BINARY`: PASS. Cuatro snapshots: mixto/vacío × Markdown/HTML; generación repetida idéntica; escape de scripts, etiquetas, atributos y enlaces maliciosos; formatos inválidos; escritura exclusiva y protección de archivos existentes; exportación CLI real en los tres formatos.
- `go vet ./...`, `gofmt` y compilación: correctos.
- CLI: HTML y Markdown del caso breaking y HTML vacío generados. Reintento sobre el HTML existente: exit 1 y SHA-256 sin cambios. JSON predeterminado conserva la explicación y los acentos.
- PowerShell 5.1 con CP850: exportación directa de HTML/Markdown en UTF-8, explicación con acentos comprobada exactamente al leer los archivos.
- Snapshots y código revisados. **Revisión visual pendiente del usuario:** el navegador integrado rechazó abrir `file://` por su política de seguridad. No se abrió una vista previa ni se afirma validación visual. Los archivos locales de ejemplo están preparados en `dist/issue-5-preview.html`, `dist/issue-5-preview.md` y `dist/issue-5-empty.html`.

No se probó ARM64 ni otros sistemas operativos. No se implementa configuración de severidad, CI, Overlays ni otras fases. La PR se dirige a `develop`, requiere aprobación explícita para fusionar y no se avanza al siguiente issue.
