# Explicaciones deterministas — issue #4

Alcance: [issue #4](https://github.com/v1kz25/openapi-contract-lab/issues/4), después del adaptador de la PR #12. `internal/explain` enriquece la salida de `compare` con textos propios en español y enlaces locales. Go 1.27.4, oasdiff 1.39.0, biblioteca estándar; ninguna dependencia nueva ni servicio/API de IA. Coste adicional obligatorio: **0 €**. Se mantienen las licencias y la decisión de integración documentadas en el spike.

## Contrato

Cada hallazgo conserva todos los campos del adaptador, incluido `upstream`, y añade:

| Campo de `explanation` | Contenido |
| --- | --- |
| `kind` | `template` para una combinación verificada de motor/versión/rule ID/nivel; `fallback` en los demás casos |
| `summary` | Qué cambia según esa regla |
| `impact` | Consecuencia contractual del hallazgo individual, sin asegurar el comportamiento real del servidor ni la compatibilidad global de la API |
| `action` | Comprobación o coordinación sugerida; no aplica cambios automáticos |
| `locations` | Fuentes base/head recibidas del motor, con archivo/línea/columna y URL local cuando puede resolverse sin adivinar |

No se analiza ni se extraen parámetros del texto libre del motor. El mensaje original proporciona el nombre o valor concreto; las plantillas describen la regla. Los diez pares verificados son:

| Rule ID | Nivel upstream verificado | Cambio sintético real |
| --- | --- | --- |
| `request-parameter-became-required` | 3 / ERR | Parámetro opcional → obligatorio |
| `request-parameter-became-optional` | 1 / INFO | Obligatorio → opcional |
| `request-parameter-became-enum` | 3 / ERR | Cadena sin enum → enum `[red]` |
| `request-parameter-became-not-enum` | 1 / INFO | Enum `[red]` → cadena sin enum |
| `request-parameter-enum-value-added` | 1 / INFO | Enum de petición `[red]` → `[red, blue]` |
| `request-parameter-enum-value-removed` | 3 / ERR | Enum de petición `[red, blue]` → `[red]` |
| `request-parameter-became-nullable` | 1 / INFO | Esquema de parámetro string → string con `nullable: true` |
| `request-parameter-became-not-nullable` | 3 / ERR | Cambio anterior invertido |
| `response-property-enum-value-added` | 3 / ERR | Enum de propiedad de respuesta `[red]` → `[red, blue]` |
| `response-property-enum-value-removed` | 1 / INFO | Cambio anterior invertido |

Las pruebas generan diez pares OpenAPI 3.0.3 mínimos y seguros en directorios temporales (`internal/explain/cases_test.go`). No se extrapola cobertura de dialectos 3.1/3.2 a partir de ellos. Cada explicación trata ese hallazgo: una adición y una retirada simultáneas pueden producir varios hallazgos distintos.

### Fallback y severidad

Una regla sin plantilla, un nivel distinto del par verificado o una versión de motor diferente recibe texto neutral que remite al mensaje y los datos originales. Esto incluye reglas conocidas cuya severidad se reduzca por contexto. No se sustituye INFO por ERR ni se inventa una causa para un nivel desconocido. Los campos `severity`, `compatibility` y `upstreamLevel` permanecen intactos.

Una lista de hallazgos vacía sigue siendo `[]`: no se inventa una explicación para el cambio de descripción. Los códigos de salida siguen siendo operativos, con 0 para una comparación correcta incluso si hay breaking changes. No se implementa la política CI del issue #6 ni los informes del issue #5.

## Enlaces

Se utiliza `file:///...#L12`, con caracteres especiales escapados y línea/columna también como datos separados. El fragmento identifica la línea; el posicionamiento automático depende del editor o visor que abra el enlace. No se presupone un editor instalado ni se lanzan programas automáticamente.

Solo se enlaza una fuente cuya ruta completa coincida con base/head. Se contempla la forma sin unidad que emite oasdiff en Windows y, si hace falta, una sola decodificación de su ruta URI-escapada. La URL se construye desde la entrada local conocida; `file`, `line` y `column` originales se conservan. No se resuelven fuentes por su mero nombre de archivo. Una fuente desconocida conserva sus datos sin `url`; una fuente ausente no genera una ubicación inventada.

## Fuentes comprobadas

Catálogo real de la versión fijada, consultado el 2026-11-21:

```powershell
.\dist\oasdiff-1.39.0\oasdiff.exe checks changelog -f json
```

Las reglas, niveles y efectos se contrastaron con ese catálogo y con los diez casos ejecutados. Referencias primarias adicionales: [obligatoriedad de parámetros](https://github.com/oasdiff/oasdiff/blob/v1.39.0/checker/check_request_parameter_required_value_updated.go), [valores enum de parámetros](https://github.com/oasdiff/oasdiff/blob/v1.39.0/checker/check_request_parameter_enum_value_updated.go), [adición de valores en respuestas](https://github.com/oasdiff/oasdiff/blob/v1.39.0/checker/check_response_property_enum_value_added.go). No se copió implementación upstream; las plantillas son texto propio. La dependencia oasdiff conserva [Apache-2.0](https://github.com/oasdiff/oasdiff/blob/v1.39.0/LICENSE).

## Prueba manual

Desde PowerShell en la carpeta local `openapi-contract-lab` (si aparece System32, entrar primero con `Set-Location` a la ruta real del repositorio):

```powershell
$informe = .\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml | ConvertFrom-Json
$LASTEXITCODE # 0
$informe.findings[0].explanation | Format-List kind, summary, impact, action
$informe.findings[0].explanation.locations | Format-List *
```

Esperado: `kind: template`, resumen «El parámetro pasa de opcional a obligatorio.», impacto y acción en español; dos ubicaciones con línea 12, columna 11 y enlaces `file:///...#L12`. Los campos originales siguen indicando error/breaking y nivel 3. Para probar la explicación INFO, invertir required/base. El cambio de descripción sigue devolviendo cero hallazgos.

Un caso fuera de las diez plantillas, generado solo con datos sintéticos:

```powershell
$base = Get-Content .\spikes\oasdiff\fixtures\base.yaml -Raw
$base.Replace('type: string', 'type: integer') | Set-Content .\dist\issue-4-type-change.yaml -Encoding UTF8
$otro = .\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe .\spikes\oasdiff\fixtures\base.yaml .\dist\issue-4-type-change.yaml | ConvertFrom-Json
$otro.findings[0].explanation | Format-List kind, summary, impact, action
```

Esperado: regla upstream `request-parameter-type-changed`, `kind: fallback`, texto neutral y clasificación error/breaking conservada. El ejecutable está compilado y se ha ejecutado en el workspace; termina al imprimir el JSON.

## Verificación real — 2026-11-21

- `go test ./...`, con `OASDIFF_BINARY` fijado y descargas Go deshabilitadas: PASS. Diez rule IDs/niveles contrastados con oasdiff real; snapshots de mensajes y fallback, preservación del hallazgo original, enlaces con espacios/Unicode/#, omisión de unidad Windows y fuentes desconocidas/ausentes.
- `go vet ./...`, `gofmt` y compilación: correctos.
- CLI compilada: explicación y enlaces del parámetro obligatorio correctos; descripción con cero hallazgos; cambio de tipo con fallback neutral.
- Windows PowerShell 5.1: JSON parseado y resumen en español con acentos conservados, comprobados guardando el texto y comparándolo exactamente.
- La primera integración falló por las rutas URI-escapadas del directorio temporal y por un fixture JSON de respuesta inválido. Se corrigieron ambos, se añadió validación del JSON sintético y se repitió la verificación completa correctamente.

Entorno: Windows x64 build 26300.9974, Go 1.27.4 y oasdiff 1.39.0 ya comprobados. No se verificaron ARM64 ni otros sistemas operativos. Pendiente de revisión y prueba manual del usuario antes de solicitar fusión. No se avanza al siguiente issue.
