# Explicaciones deterministas — issues #4 y #21

Alcance: [issue #4](https://github.com/v1kz25/openapi-contract-lab/issues/4), después del adaptador de la PR #12. `internal/explain` enriquece la salida de `compare` con textos propios en español y enlaces locales. Go 1.27.4, oasdiff 1.39.0, biblioteca estándar; ninguna dependencia nueva ni servicio/API de IA. Coste adicional obligatorio: **0 €**. Se mantienen las licencias y la decisión de integración documentadas en el spike.

El [issue #21](https://github.com/v1kz25/openapi-contract-lab/issues/21) añadió tres plantillas a las diez publicadas en v0.1.0. La ampliación se incluye en [v0.2.0 experimental publicada](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.2.0); no altera los archivos de v0.1.0, sus diez explicaciones, el esquema JSON ni los códigos de salida. Las [notas v0.2.0](releases/v0.2.0.md) conservan la redacción histórica del candidato #33.

## Contrato

Cada hallazgo conserva todos los campos del adaptador, incluido `upstream`, y añade:

| Campo de `explanation` | Contenido |
| --- | --- |
| `kind` | `template` para una combinación verificada de motor/versión/rule ID/nivel; `fallback` en los demás casos |
| `summary` | Qué cambia según esa regla |
| `impact` | Consecuencia contractual del hallazgo individual, sin asegurar el comportamiento real del servidor ni la compatibilidad global de la API |
| `action` | Comprobación o coordinación sugerida; no aplica cambios automáticos |
| `locations` | Fuentes base/head recibidas del motor, con archivo/línea/columna y URL local cuando puede resolverse sin adivinar |

No se analiza ni se extraen parámetros del texto libre del motor. El mensaje original proporciona el nombre o valor concreto; las plantillas describen la regla.

La CLI serializa el JSON usando únicamente bytes ASCII y escapes Unicode estándar (`\u00e1`, por ejemplo). Así, PowerShell 5.1 puede leer el pipeline con CP850, Windows-1252 o UTF-8 antes de `ConvertFrom-Json` sin corromper los acentos. Los valores Unicode, rutas y eventos upstream se conservan al parsear; los caracteres fuera del plano básico se representan con pares sustitutos UTF-16. No se modifica la configuración de la consola del usuario.

Los trece pares verificados y publicados en v0.2.0 son:

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
| `request-parameter-type-changed` | 3 / ERR | Tipo del parámetro query string → integer |
| `request-body-became-required` | 3 / ERR | Cuerpo JSON de POST opcional → obligatorio |
| `response-required-property-removed` | 3 / ERR | Propiedad `id` obligatoria retirada del esquema de respuesta 200 |

Las pruebas usan trece pares OpenAPI 3.0.3 mínimos y seguros en directorios temporales (`internal/explain/cases_test.go`): diez generados y tres pares persistidos en `internal/explain/testdata/breaking/` para reproducir manualmente los casos nuevos. No se extrapola cobertura de dialectos 3.1/3.2 a partir de ellos. Cada explicación trata ese hallazgo: una adición y una retirada simultáneas pueden producir varios hallazgos distintos.

Las nuevas plantillas describen el impacto contractual y una revisión sugerida; no garantizan lo que el servidor hace en ejecución. `request-parameter-type-changed` puede señalar tipo o formato, por lo que remite a los valores del mensaje original sin analizarlos ni afirmar que todo cambio de tipo es incompatible. En el par concreto inverso integer → string, oasdiff emite `request-parameter-type-generalized`, nivel 1, y se conserva el fallback; no se reinterpreta su semántica. Los otros inversos comprobados son `request-body-became-optional` y `response-required-property-added`, ambos nivel 1 y fallback. No se añaden plantillas para esas tres reglas inversas.

### Fallback y severidad

Una regla sin plantilla, un nivel distinto del par verificado o una versión de motor diferente recibe texto neutral que remite al mensaje y los datos originales. Esto incluye reglas conocidas cuya severidad se reduzca por contexto. No se sustituye INFO por ERR ni se inventa una causa para un nivel desconocido. Los campos `severity`, `compatibility` y `upstreamLevel` permanecen intactos.

Una lista de hallazgos vacía sigue siendo `[]`: no se inventa una explicación para el cambio de descripción. Las explicaciones no cambian la [política vigente](SEVERITY-POLICY.md): exit 3 si el informe no supera el umbral (por defecto `error`), 0 si lo supera o con `--fail-on none`, 1 para error operativo y 2 para argumentos inválidos.

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
$LASTEXITCODE # 3: informe generado; breaking change con política predeterminada
$informe.findings[0].explanation | Format-List kind, summary, impact, action
$informe.findings[0].explanation.locations | Format-List *
```

Esperado: `kind: template`, resumen «El parámetro pasa de opcional a obligatorio.», impacto y acción en español; dos ubicaciones con línea 12, columna 11 y enlaces `file:///...#L12`. Los campos originales siguen indicando error/breaking y nivel 3. Para probar la explicación INFO, invertir required/base. El cambio de descripción sigue devolviendo cero hallazgos.

Un caso fuera de las trece plantillas, usando el par sintético inverso de tipo:

```powershell
$fixtures = '.\internal\explain\testdata\breaking'
$otro = .\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe "$fixtures\parameter-head.json" "$fixtures\parameter-base.json" | ConvertFrom-Json
$LASTEXITCODE # 0: el motor clasifica este inverso como INFO
$otro.findings[0].explanation | Format-List kind, summary, impact, action
```

Esperado: regla upstream `request-parameter-type-generalized`, `kind: fallback`, texto neutral y clasificación info/non-breaking conservada. Para una regla desconocida o para cualquiera de las tres nuevas con nivel no verificado, versión u otro motor, también se conserva el fallback; no se cambian datos upstream.

### Comprobar las tres nuevas explicaciones

Tras compilar con las herramientas fijadas y desde la raíz del repositorio:

```powershell
$fixtures = '.\internal\explain\testdata\breaking'
foreach ($caso in @('parameter', 'body', 'response')) {
    $informe = .\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe "$fixtures\$caso-base.json" "$fixtures\$caso-head.json" | ConvertFrom-Json
    $LASTEXITCODE # 3: informe completo y breaking con política predeterminada
    $informe.findings[0] | Select-Object ruleId, upstreamLevel, severity, compatibility
    $informe.findings[0].explanation | Format-List kind, summary, impact, action
}
$salida = '.\dist\body-' + [guid]::NewGuid().ToString('N') + '.html'
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe --format html --output $salida "$fixtures\body-base.json" "$fixtures\body-head.json"
$LASTEXITCODE # 3, archivo HTML completo; Invoke-Item $salida permite abrirlo
```

Esperado por cada par: un hallazgo de su ID de la tabla, nivel 3, error/breaking, `kind: template`, impacto y acción en español y fuentes reales del motor. Al añadir `--fail-on none` antes de los archivos, exit 0 y el mismo contenido. Los informes siguen escapando datos y rechazando sobrescritura. Los tests comprueban los tres formatos; no se modificó la presentación del renderer.

## Verificación del issue #21 — 2026-11-21

- Catálogo `checks changelog -f json` de oasdiff 1.39.0 y tres pares sintéticos reales: los tres IDs/niveles 3 de la tabla confirmados; inversos comprobados por separado, sin inferir simetría.
- `scripts/verify.ps1`: PASS con Go 1.27.4 y motor real, tests sin caché, formato, vet y build. Cachés Go fijadas en `dist/` y descargas deshabilitadas. La verificación base anterior pasó, pero emitió un aviso de permiso al intentar escribir la caché de módulos global; fijar también `GOMODCACHE` local evitó el aviso.
- Trece mensajes en snapshot; los diez originales conservados exactamente. Fallback de las tres nuevas reglas comprobado para niveles 1/2/desconocido, versión distinta y otro motor; datos upstream intactos.
- Integración CLI: tres reglas × JSON/Markdown/HTML, exit 3 predeterminado y 0 con `none`, contenido idéntico entre políticas, archivos existentes y entradas conservados. Se mantienen los snapshots de renderer y tests de seguridad originales.
- `scripts/smoke.ps1`: PASS en PowerShell 7.6.9 y Windows PowerShell 5.1.26100.9818. El ejecutable compilado produce los tres nuevos hallazgos y exportaciones UTF-8 en ambos entornos.

Solo Windows x64 comprobado. Las nuevas fixtures son 3.0.3, no amplían por sí solas la cobertura de 3.1. La CI del commit y su resultado se registran en la PR; no se confunden con pruebas locales. Sin publicar versión ni fusionar la PR, y sin avanzar al siguiente issue funcional.

## Verificación histórica del issue #4 — 2026-11-21

Evidencia anterior a los informes y la política de #5/#6; se conserva como registro histórico. Las instrucciones activas anteriores reflejan la política actual.

- `go test ./...`, con `OASDIFF_BINARY` fijado y descargas Go deshabilitadas: PASS. Diez rule IDs/niveles contrastados con oasdiff real; snapshots de mensajes y fallback, preservación del hallazgo original, enlaces con espacios/Unicode/#, omisión de unidad Windows y fuentes desconocidas/ausentes.
- `go vet ./...`, `gofmt` y compilación: correctos.
- CLI compilada: explicación y enlaces del parámetro obligatorio correctos; descripción con cero hallazgos; cambio de tipo con fallback neutral.
- La comprobación inicial con Windows PowerShell 5.1 y codificación heredada UTF-8 no cubría las consolas CP850. El usuario detectó acentos corruptos; se reprodujo exactamente el fallo con CP850 antes de corregir la serialización JSON.
- Tras la corrección: Windows PowerShell 5.1 con CP850, Windows-1252 y UTF-8, mensajes template/fallback y ubicaciones comprobados exactamente. JSON de la CLI verificado byte a byte como ASCII; tests de roundtrip con acentos, rutas Unicode, emoji, CJK, escapes literales y caracteres de control, además de errores de escritura.
- La primera integración falló por las rutas URI-escapadas del directorio temporal y por un fixture JSON de respuesta inválido. Se corrigieron ambos, se añadió validación del JSON sintético y se repitió la verificación completa correctamente.

Entorno: Windows x64 build 26300.9974, Go 1.27.4 y oasdiff 1.39.0 ya comprobados. No se verificaron ARM64 ni otros sistemas operativos. Pendiente de revisión y prueba manual del usuario antes de solicitar fusión. No se avanza al siguiente issue.
