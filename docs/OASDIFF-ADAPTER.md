# Adaptador local — issue #3

Alcance: [issue #3](https://github.com/v1kz25/openapi-contract-lab/issues/3), basado en el setup integrado por la PR #11. Entrada base/head de archivos locales y salida JSON normalizada. El motor sigue siendo oasdiff **1.39.0**; Go **1.27.4**, biblioteca estándar, coste adicional obligatorio **0 €**. No se añade dependencia Go, servicio ni motor de diff propio. El ejecutable mantiene la licencia [Apache-2.0](https://github.com/oasdiff/oasdiff/blob/v1.39.0/LICENSE); descarga/checksum y compatibilidad Windows comprobados en el [spike](OASDIFF-SPIKE.md).

## Comportamiento

`compare [--oasdiff PATH] BASE HEAD` comprueba que ambas entradas sean archivos regulares legibles, localiza el ejecutable y exige `oasdiff version 1.39.0`. Después ejecuta `changelog --format json --lang en --level INFO`. Una configuración temporal vacía explícita y la exclusión de variables `OASDIFF_*` aíslan la configuración del motor. Se pasan `--allow-external-refs=false`, `--open=false` y `--fetch=false`. No hay shell intermedia ni aceptación de URLs, stdin o rutas UNC como entradas.

El contexto limita versión y comparación a 30 segundos en total. La captura de stdout está limitada a 16 MiB y stderr a 64 KiB; los fallos devuelven un error, no un informe parcial. `*oasdiff.Error` conserva etapa y causa (entrada, ejecutable, versión, configuración, comparación o decodificación). Los errores de proceso incorporan su diagnóstico stderr. No se descarga ni instala oasdiff automáticamente.

La CLI devuelve 0 si el informe supera la política, 3 si se genera pero no supera `--fail-on` (por defecto `error`), 1 si falla la entrada/motor/salida y 2 para argumentos inválidos. **Un finding breaking no activa exit 1**: con la política predeterminada activa exit 3. `--fail-on none` permite exportarlo con exit 0. Ver [política vigente](SEVERITY-POLICY.md).

### Límites y alcance de seguridad — issue #38

El uso soportado son **archivos locales revisados**, con el motor fijado de procedencia comprobada. No se ofrece un servicio para contratos arbitrarios hostiles. El rechazo de URL/UNC y referencias externas más el aislamiento de configuración restringen la invocación; no constituyen un sandbox del ejecutable, auditoría de tráfico ni garantía general de que un binario malicioso no acceda a red/archivos.

- El deadline de **30 s** cubre conjuntamente versión y comparación del motor; no incluye renderizado/escritura posterior de la CLI. Un contexto padre más corto también se respeta. `WaitDelay=2 s` limita la espera adicional de proceso/pipes conforme a Go; no es una cuota total de recursos de la aplicación.
- La captura por invocación acepta hasta **16 MiB stdout / 64 KiB stderr**; exceder cualquiera produce error, nunca un informe truncado aceptado. El diagnóstico identifica el límite y conserva la causa de proceso aunque haya salida no cero. No se garantiza cancelación inmediata al alcanzar el límite: el deadline sigue acotando la ejecución.
- No hay límite propio de tamaño de entrada, cuota de RAM/CPU del motor ni cuota de memoria total de la CLI: parsear/enriquecer/renderizar puede necesitar más memoria que los bytes capturados. No se prometen manejo seguro de cualquier tamaño/profundidad, aislamiento de procesos descendientes ni defensa universal frente a recursos agotados.
- Ante fallo de motor/versión/decodificación, la CLI devuelve **1**, sin informe en stdout ni creación del archivo de salida; prevalece incluso con `--fail-on none`. Una escritura posterior puede dejar una salida nueva parcial, como declara el contrato CLI. Ver [recuperación de informes](REPORTS.md#fallos-de-escritura-y-recuperación--issue-38).

Los tests de #38 usan un ejecutable Go sintético en `internal/oasdiff/testdata/engine/main.go` compilado localmente, sin descargas. Prueban pipes reales exactamente en el límite y por encima (incluida terminación no cero), JSON incompleto, fallo de proceso, cancelación activa de comparación y deadline durante versión. La comprobación de proceso terminado usa un handle Windows abierto mientras el helper está vivo: prueba que ese helper directo terminó, no que cualquier árbol de procesos arbitrario quede aislado. El padre de los tests usa deadline de 3 s para ese caso y tolerancia de 8 s, sin esperar 30 s por prueba; el deadline de producción no cambia. Las pruebas CLI recorren ambos umbrales `error`/`none`, tres formatos y salidas a stdout/destino nuevo/existente/entradas, comprobando conservación por hash. Son pruebas de fronteras, no evidencia semántica nueva: el corpus con oasdiff 1.39.0 real se mantiene aparte.

## Modelo y fidelidad

La referencia completa del JSON enriquecido emitido por la CLI es el [contrato JSON v1](#contrato-json-v1--issue-37). Las estructuras Go internas no son una API pública de biblioteca. Este contrato describe el comportamiento actual a estabilizar; no convierte v0.2.0 en una versión estable.

El informe lleva `engine`, `engineVersion`, `base`, `head` y `findings`. Base/head son las rutas absolutas realmente enviadas al motor. Cada hallazgo contiene:

| Campo | Procedencia |
| --- | --- |
| `ruleId`, `upstreamLevel`, `message` | `id`, `level`, `text`, sin reinterpretar el ID o el texto |
| `operation`, `operationId`, `path`, `section` | Campos homónimos cuando existen; no se exige ruta a eventos de componentes |
| `baseSource`, `headSource` | `baseSource`, `revisionSource`, con archivo/línea/columna originales |
| `fingerprint` | Fingerprint del motor, sin recalcular |
| `upstream` | Evento JSON completo, incluidos comment, atributos, disclaimers o campos desconocidos |

El mapeo de severidad se basa exclusivamente en [los niveles upstream fijados](https://github.com/oasdiff/oasdiff/blob/v1.39.0/checker/rules/level.go):

| Nivel | `severity` | `compatibility` |
| --- | --- | --- |
| 3 / ERR | `error` | `breaking` |
| 2 / WARN | `warning` | `potentially-breaking` |
| 1 / INFO | `info` | `non-breaking` |
| Cualquier otro entero | `unknown` | `unknown` |

Un ID nuevo se conserva con el nivel que indique el motor; no requiere una lista local de reglas. Un nivel desconocido se conserva sin clasificarlo como compatible. Se rechazan respuestas que no sean arrays, eventos sin ID/nivel entero, tipos incompatibles y ubicaciones negativas. Campos adicionales permanecen en `upstream`. La representación JSON puede cambiar de espacios/escapes al serializar; sus valores se conservan.

En Windows, el motor observado devuelve archivos como `/Users/.../base.yaml`, sin la letra de unidad. El adaptador conserva esa localización original y añade por separado `base`/`head` con las rutas Windows completas. No inventa ni sustituye rutas de fuentes desconocidas. Los golden con motor real normalizan únicamente los nombres locales para compararlos con la evidencia portable del spike.

`changelog` incluye cambios que afectan a consumidores, no cada diferencia del documento: una modificación solo de descripción produce `findings: []`. Esto **no significa contratos idénticos**. No se llama a `diff` ni se crean hallazgos documentales propios. Los mensajes originales son los del motor; las explicaciones propias corresponden a los issues #4/#21. Tampoco se afirma cobertura completa de OpenAPI 3.1/3.2 con estos fixtures 3.0.3.

## Contrato JSON v1 — issue #37

Un informe es un objeto JSON; se consume cuando la CLI termina con **0 o 3**, nunca asumir informe completo con 1/2. Campos propios no usan `null` para representar ausencias. `findings` y `explanation.locations` son arrays, incluso vacíos (`[]`). Las rutas pueden contener datos locales privados y deben revisarse antes de compartir el informe.

| Campo del informe | Tipo | Presencia y significado |
| --- | --- | --- |
| `engine` | string | Requerido; `oasdiff` |
| `engineVersion` | string | Requerido; `1.39.0`, versión comprobada del motor, no del esquema ni de la CLI |
| `base`, `head` | string | Requeridos en una comparación CLI completada: rutas absolutas enviadas al motor. Los modelos internos/ejemplos normalizados pueden omitirlos si no tienen contexto; eso no cambia la garantía de la CLI. |
| `findings` | array de Finding | Requerido; orden upstream, conserva repeticiones; vacío no demuestra igualdad de contratos |

| Campo de Finding | Tipo | Presencia y significado |
| --- | --- | --- |
| `ruleId` | string no vacío | Requerido; ID original, no catálogo cerrado |
| `severity` | string | Requerido: `error`, `warning`, `info` o `unknown` |
| `compatibility` | string | Requerido: `breaking`, `potentially-breaking`, `non-breaking` o `unknown`, según tabla de niveles anterior |
| `upstreamLevel` | integer | Requerido; entero original aceptado por el adaptador, incluidos niveles desconocidos; no inferir clasificación a partir del texto |
| `message` | string | Omitido si vacío; texto original `text` del motor |
| `operation`, `operationId`, `path`, `section` | string | Cada uno omitido si vacío; no inferir que todo hallazgo tiene operación o ruta |
| `baseSource`, `headSource` | Source (objeto) | Cada uno omitido si no existe fuente; head procede de `revisionSource`, sin sustituir la ruta original |
| `fingerprint` | string | Omitido si vacío; valor opaco del motor, no una identidad durable garantizada entre motores/versiones |
| `upstream` | objeto JSON | Requerido; evento original completo, incluidos campos desconocidos, arrays y null de upstream |
| `explanation` | Explanation (objeto) | Requerido; enriquecimiento determinista propio |

| Objeto / campo | Tipo | Presencia y significado |
| --- | --- | --- |
| Source.`file` | string no vacío | Requerido; ruta que reportó el motor; no necesariamente ruta absoluta Windows con unidad |
| Source.`line`, Source.`column` | integer positivo | Omitidos si cero/no informado; posiciones originales basadas en 1, no se inventa posición si falta |
| Explanation.`kind` | string | Requerido: `template` o `fallback` |
| Explanation.`summary`, `.impact`, `.action` | string | Requeridos; prosa española para revisión, no identificadores ni clasificación global de compatibilidad |
| Explanation.`locations` | array de Location | Requerido; cero a dos fuentes, base antes de head; no crear fuente ausente |
| Location.`side` | string | Requerido: `base` o `head` |
| Location.`file` | string no vacío | Requerido; fuente original, también si no se puede enlazar |
| Location.`line`, `.column` | integer positivo | Omitidos si cero/no informado; ausencia no equivale a línea/columna 0 |
| Location.`url` | string | Omitido si no se resuelve con seguridad; URL `file:///...` de la entrada conocida, fragmento `#L<n>` solo si hay línea. No es una promesa de que el visor salte a ella. |

Una fuente ausente, una fuente con solo archivo y una fuente con URL son estados distintos. No reemplazar ausencias por `null`, `false`, ruta vacía o posiciones inventadas. Valores desconocidos de nivel se conservan y se normalizan a `unknown` en ambos campos de clasificación; con política activa bloquean, con `none` no. Un ID desconocido con nivel conocido conserva la clasificación de ese nivel. `fallback` tampoco equivale a severidad desconocida: indica falta de explicación verificada para ese motor/regla/nivel.

Los consumidores deben leer por nombre de campo, tolerar campos opcionales adicionales e ignorar campos que no necesiten. Deben conservar el orden de arrays que utilicen; el orden de claves, espacios, escapes equivalentes y prosa no son significativos. `upstream` es contenido opaco y extensible del motor: se preservan sus valores (incluidos arrays/repeticiones/campos desconocidos), sin prometer representación textual idéntica. No se deben convertir números upstream indiscriminadamente a coma flotante si se necesita su precisión. El orden emitido por una versión futura del motor no se deduce del actual.

### Ejemplos de consumidor

Vacío, con rutas sintéticas ilustrativas (en Windows la CLI utiliza las rutas absolutas reales):

```json
{"engine":"oasdiff","engineVersion":"1.39.0","base":"C:\\synthetic\\base.yaml","head":"C:\\synthetic\\head.yaml","findings":[]}
```

Finding de plantilla sin fuentes; el siguiente es un evento **sintético de contrato**, no una captura adicional de oasdiff. La prosa se abrevia para mostrar la forma; textos reales y fixtures verificados están en [EXPLANATIONS.md](EXPLANATIONS.md).

```json
{"ruleId":"request-parameter-became-required","severity":"error","compatibility":"breaking","upstreamLevel":3,"upstream":{"id":"request-parameter-became-required","level":3},"explanation":{"kind":"template","summary":"El parámetro pasa de opcional a obligatorio.","impact":"Revisar consumidores.","action":"Revisar el cambio.","locations":[]}}
```

Finding fallback de nivel desconocido, también sintético: `baseSource.file` existe pero no se puede enlazar a una entrada conocida; no se inventan línea, columna, URL ni fuente head. Campos extra de upstream conservados:

```json
{"ruleId":"future-rule","severity":"unknown","compatibility":"unknown","upstreamLevel":99,"baseSource":{"file":"unknown.yaml"},"upstream":{"id":"future-rule","level":99,"baseSource":{"file":"unknown.yaml"},"extra":{"values":[2,1,2],"note":null}},"explanation":{"kind":"fallback","summary":"No hay explicación verificada.","impact":"Consultar datos originales.","action":"Revisar manualmente.","locations":[{"side":"base","file":"unknown.yaml"}]}}
```

La versión normativa se denomina **contrato JSON v1**, ligada a la futura CLI 1.x, sin campo nuevo `schemaVersion`. La decisión, migración y actualización upstream están en [DEVELOPMENT.md](DEVELOPMENT.md#compatibilidad-del-contrato-v1). `engineVersion` nunca debe usarse como versión del esquema.

## Comprobación manual

Abrir PowerShell **en la raíz del repositorio** (si se abre en System32, cambiar primero de carpeta con `Set-Location 'RUTA_DEL_REPOSITORIO'`). Compilar según [DEVELOPMENT.md](DEVELOPMENT.md), preparar oasdiff según el spike y ejecutar:

```powershell
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\required.yaml
$LASTEXITCODE # 3; un finding request-parameter-became-required, error/breaking
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe .\spikes\oasdiff\fixtures\base.yaml .\spikes\oasdiff\fixtures\description.yaml
$LASTEXITCODE # 0; findings: []
.\dist\openapi-contract-lab.exe compare --oasdiff .\dist\oasdiff-1.39.0\oasdiff.exe .\spikes\oasdiff\fixtures\required.yaml .\spikes\oasdiff\fixtures\base.yaml
$LASTEXITCODE # 0; request-parameter-became-optional, info/non-breaking
```

En el primer caso, comprobar GET `/items`, nivel original 3 y ambas fuentes en línea 12, columna 11. El ejecutable se compiló y ejecutó en el workspace para facilitar la revisión; termina al imprimir el JSON, sin servidor que mantener arrancado.

## Verificación histórica del issue #3 — 2026-11-21

Esta evidencia precede a la política del issue #6. Sus códigos 0 para breaking no describen el comportamiento actual; la comprobación manual anterior usa la política vigente.

En Windows x64 build 26300.9974, SDK portátil Go 1.27.4 y oasdiff 1.39.0 ya comprobados, con descargas Go deshabilitadas:

- `go test ./...` con `OASDIFF_BINARY` fijado: correcto. Incluye golden breaking/descripción/evento desconocido; casos malformados, niveles INFO/WARN/ERR/desconocido, errores de entrada/motor/versión/JSON, cancelación y límite de salida.
- Integración con motor real: parámetro obligatorio → breaking; descripción → sin hallazgos; parámetro que vuelve a opcional → info/non-breaking. Se conserva la localización y el fingerprint coincide con la evidencia del spike. Configuración de usuario inválida y variable `OASDIFF_CONFIG` no alteran el resultado. YAML malformado y `$ref` externo local fallan con diagnóstico.
- `go vet ./...` y `go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab`: correctos.
- Ejecutable compilado: los tres ejemplos anteriores producen los resultados esperados con exit 0; motor inexistente produce exit 1.
- Invocación desde `powershell.exe -NoProfile -NonInteractive` (Windows PowerShell 5.1): comparación breaking correcta, JSON legible y exit 0.
- La primera prueba de integración detectó la forma de rutas sin unidad que emite oasdiff en Windows. Se corrigió la expectativa del test, se añadieron base/head al informe y se repitió la verificación completa correctamente.

No hay servicios de ejecución externos ni llamadas de IA. La configuración offline evita las referencias externas, pero no es una auditoría de red del ejecutable. No se probó ARM64 ni otros sistemas operativos. La fusión de la PR requiere aprobación explícita; no se continúa al issue #4.
