# Spike: señales sobre diferencias — issue #25

## Decisión

**Reducir la propuesta y posponer su integración pública.** Con oasdiff **1.39.0**, una invocación adicional a `summary --format json` permite obtener la señal explícita `diff`. Es viable proponer «diferencias detectadas / no detectadas / desconocido», con alcance ligado al motor y su configuración. No hay evidencia suficiente para prometer las categorías generales «idéntico», «solo documental» y «contractual».

Este spike no cambia la aplicación, JSON, CLI, códigos de salida, adaptador, clasificación ni explicaciones. La CLI actual sigue advirtiendo que cero hallazgos no demuestra igualdad. No se implementa la señal propuesta en v0.2.0 durante este issue. Go 1.27.4 y oasdiff 1.39.0 se conservan, sin dependencias nuevas, servicios, API de pago ni coste adicional obligatorio.

## Matriz ejecutada

Cinco pares sintéticos OpenAPI 3.0.3. Cuatro reutilizan los contratos del spike inicial; se añade un JSON con el mismo contenido representado con otro orden de claves y formato. Es el motor el que compara; no se usan texto ni hashes de los contratos para decidir su significado.

| Par conocido por construcción | `diff -f json` | `summary -f json` | `changelog -f json` | `diff --fail-on-diff` |
| --- | --- | --- | --- | --- |
| Base frente a sí misma | `{}` | `diff: false` | `[]` | 0 |
| YAML frente a JSON equivalente, compacto y reordenado | `{}` | `diff: false` | `[]` | 0 |
| Solo descripción de parámetro | Objeto con `description.from/to` | `diff: true` | `[]` | 1 |
| Obligatorio → opcional | Objeto con `required: true → false` | `diff: true` | `request-parameter-became-optional`, INFO/1 | 1 |
| Opcional → obligatorio | Objeto con `required: false → true` | `diff: true` | `request-parameter-became-required`, ERR/3 | 1 |

Sin `--fail-on-diff` ni `--fail-on`, los quince comandos de la matriz terminan con **0**, incluso cuando hay diferencias o ERR. Los códigos anteriores son del motor, no los códigos 0/1/2/3 de nuestra CLI. La política actual de la aplicación se mantiene.

Las salidas completas, con fuentes relativas y sin rutas personales, están en [spikes/change-kind/evidence](../spikes/change-kind/evidence/). [cases.json](../spikes/change-kind/cases.json) registra los cinco pares y sus expectativas. [commands.json](../spikes/change-kind/evidence/commands.json) registra argumentos, códigos, tamaños y tiempos reales; los tiempos son una observación local, no snapshots reproducibles ni umbrales de rendimiento.

Control adicional: `diff --exclude-elements description` oculta el único cambio del par documental y devuelve `{}`. Demuestra que filtros y configuración afectan a la señal; no se usa ese resultado para deducir una clasificación general. Descartar también todas las extensiones no sería una prueba segura de que el resto es solo documentación.

## Qué aporta cada superficie del motor

Se inspeccionó el tag **v1.39.0**, resuelto al commit `21b8a57a0a192e7478b764ec0164b5694cce3393`. No se copió código ni se añadió el módulo Go:

- [CLI changelog](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/internal/changelog.go): calcula el diff y conoce `Empty()` antes de ejecutar los checks. Pasa esa señal al formatter, pero la salida JSON CLI sigue siendo el array de eventos.
- [Formatter JSON](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/formatters/format_json.go) y [RenderOpts](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/formatters/types.go): la API del formatter tiene una opción `WrapInObject` con `changes` y `diff_empty`. La ruta CLI inspeccionada no activa esa opción y la ayuda no ofrece un flag para ella. No se propone recurrir a un servicio externo.
- [CLI summary](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/internal/summary.go) y [modelo Summary](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/diff/summary.go): ofrecen el booleano explícito `diff` y recuentos agregados. No ofrecen una categoría documental/contractual ni severidad por detalle. Las descripciones y los cambios contractuales del corpus tienen los mismos recuentos.
- [CLI diff](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/internal/diff.go): expone diferencias estructuradas sin clasificación de compatibilidad; al serializar JSON excluye `endpoints`. No conviene interpretar sus claves como una lista cerrada de categorías semánticas.
- [API diff](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/diff/diff.go): `GetWithOperationsSourcesMap` permite obtener el diff completo y `Empty()`. Cambiar a integración como módulo implicaría dependencias y una nueva decisión. El ejemplo público con `GetPathsDiff` tiene alcance de paths, por lo que no se debe tratar sin más como un sustituto del documento completo.

La comparación sucede después de cargar/normalizar las especificaciones y depende de [configuración](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/diff/config.go), opciones y cobertura del motor. Una ausencia de diferencias detectadas no certifica igualdad textual ni equivalencia de todos los comportamientos posibles.

## Información en la frontera actual

`internal/oasdiff/model.go` acepta el array JSON de changelog. Conserva cada evento completo en `upstream`, incluidos campos desconocidos; no se ha detectado pérdida de campos de esos eventos. La señal global `diff_empty` no llega en ese array, y tampoco llegan cambios que no disparan un check, como la descripción del corpus. El adaptador no puede reconstruir esa información a partir de `[]`.

No se puede identificar «solo formato» frente a «idéntico» con las señales observadas: ambos dan `{}`, `diff:false` y `[]`. Tampoco se puede identificar de forma general «solo documental» por tener `diff:true` y `[]`: eso solo indica diferencias sin eventos del catálogo, no una prueba exhaustiva de que no sean contractuales.

## Propuesta para una tarea futura, no implementada

Si se decide continuar, limitar la ampliación a un objeto JSON opcional y aditivo, por ejemplo `comparison.differenceStatus`, con valores:

- `no-differences-detected`: `summary.diff == false` con el motor/configuración verificados; no usar «idénticos».
- `differences-detected`: `summary.diff == true`, independientemente del número o severidad de eventos.
- `unknown`: análisis de esta señal no solicitado o evidencia insuficiente; nunca convertirlo en ausencia de diferencias.

La migración exacta necesita otro issue: los consumidores actuales seguirían leyendo los campos existentes; la ausencia del objeto en informes antiguos no sería `false`. Conservar los hallazgos y la política como fuente de los códigos de gate. Fallos operativos de una señal solicitada seguirían siendo errores (exit 1), no un éxito silencioso con `unknown`. No se diseña aquí una nueva opción CLI ni se añaden campos.

Una futura implementación necesitaría tests de JSON aditivo/consumidores anteriores, ambas señales con y sin eventos, desconocido/no solicitado, errores y timeouts, límites de salida, entradas cambiadas entre invocaciones y config/flags idénticos. No etiquetar como «solo documental» usando una lista propia de claves del diff: eso sería ampliar la semántica sin evidencia suficiente.

## Tiempo, salida y seguridad

Medición local Windows x64/PowerShell 7.6.9, cinco pares pequeños, arranque del proceso incluido, captura registrada en commands.json:

| Comando | Tiempo observado por par | Máximo stdout UTF-8 |
| --- | --- | --- |
| `changelog` | 29–30 ms | 368 bytes |
| `diff` | 28–33 ms | 189 bytes |
| `summary` | 27–32 ms | 76 bytes |

`summary` es la candidata para una única invocación extra: menor superficie de interpretación y salida pequeña en este corpus. Sigue recalculando/cargando el diff; no evita ese coste. No se midieron contratos grandes, pico de memoria ni otras plataformas; no se extrapolan estos tiempos.

La aplicación actual limita 30 segundos totales de comparación, stdout a 16 MiB y stderr a 64 KiB. Una futura llamada extra debe compartir el presupuesto total de tiempo y fijar un presupuesto conjunto de captura, no duplicar automáticamente los límites. Necesita validar versión una sola vez, config vacía explícita, las mismas entradas/opciones y todos los bloqueos de URL/UNC/ref externas; entre procesos existe riesgo de leer versiones diferentes de un archivo que cambie. No se alteran esas protecciones aquí.

El runner del spike usa exclusivamente fixtures sintéticos, subprocesos sin shell y timeout de 30 segundos por prueba. No implementa la captura incremental limitada de la aplicación; no es un nuevo adaptador de producción. No habilita `--open`, usa `--fetch=false` y `--allow-external-refs=false`. El config inválido ambiental se inyecta solo al hijo y no afecta al config vacío explícito.

## Ejecución y resultado — 2026-11-22

```powershell
./spikes/change-kind/verify.ps1
```

Requiere PowerShell 7 y oasdiff 1.39.0 local, sin SDK, descarga ni nueva instalación. Acepta `-Binary` y `-OutputDirectory`; rechaza un directorio de salida existente. Genera nueva evidencia en `dist/`, compara las quince salidas con las revisadas y registra mediciones nuevas sin actualizar expectativas versionadas.

Estabilización del candidato #33: las rutas relativas de `OutputDirectory` se resuelven explícitamente contra la raíz del repositorio, incluso si el directorio nativo del proceso difiere de `Get-Location`. Las rutas absolutas se conservan. Se reprodujo el comportamiento anterior y se comprobaron ruta relativa, absoluta y rechazo de salida existente invocando desde otra carpeta. No cambia ninguna señal ni clasificación de la aplicación.

**PASS: 25 pruebas por ejecución**, repetidas con las mismas salidas: versión, quince comparaciones, cinco controles de `--fail-on-diff`, una exclusión de descripción y tres bloqueos de referencias externas. También comprueba que las entradas sigan idénticas byte a byte. Comparar outputs con evidencia congelada comprueba reproducibilidad, no clasifica contratos mediante texto/hashes.

El control externo usa una referencia local existente y fue rechazado por los tres comandos antes de cargarla. El motor devolvió **123** y stdout vacío, contrastado con [internal/errors.go](https://github.com/oasdiff/oasdiff/blob/21b8a57a0a192e7478b764ec0164b5694cce3393/internal/errors.go). La primera ejecución del runner había supuesto 100 y falló; se contrastó el resultado real y el código fijado, se corrigió a 123 y se repitieron las comprobaciones completas. No se cambió una salida esperada para ocultar una diferencia de contenido.

No hace falta una prueba manual para aceptar este spike: son investigación, fixtures y evidencia reproducible; no cambia un flujo de usuario. La CI existente verifica la aplicación y la regresión, pero no ejecuta automáticamente este runner de investigación; sus 25 controles se ejecutaron localmente. No se publican versión, clasificación nueva, SARIF ni API; las siguientes decisiones requieren acotar otra tarea.

Regresión de la aplicación tras preparar el spike: `scripts/verify.ps1` PASS con Go 1.27.4 y motor real (tests sin caché, gofmt, vet y build); `scripts/smoke.ps1` PASS en PowerShell 7.6.9, con códigos 0/1/2/3 y protección de entradas/salidas. La CI del commit exacto queda registrada en su PR.
