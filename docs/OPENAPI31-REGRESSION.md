# Regresión OpenAPI 3.1 — issue #24

Corpus acotado de cinco pares **3.1.0**, observado con **oasdiff 1.39.0** y Go **1.27.4** en Windows x64. Aumenta la evidencia sobre casos concretos; no certifica compatibilidad completa con 3.1 ni soporte de 3.2 u otras plataformas. No cambia código de producción, plantillas, adaptador, esquema JSON, política ni herramientas. Coste adicional obligatorio: **0 €**.

## Pares y resultados observados

`internal/cli/testdata/openapi31/cases.json` fija IDs, mensajes originales, fingerprints, niveles, clasificación, explicación y posiciones de fuente observadas. No se actualiza automáticamente. Los tres primeros pares reutilizan `../regression/`; solo se añaden tres contratos JSON, compartiendo la base de los dos últimos pares.

| Caso | Base → head | Resultado de `changelog` | Explicación / exit predeterminado |
| --- | --- | --- | --- |
| Idéntico | `regression/base.yaml` → mismo archivo | `[]` | Ninguna / 0 |
| Documental | `regression/base.yaml` → `regression/description.yaml` | `[]` | Ninguna / 0 |
| Obligatoriedad | `regression/base.yaml` → `regression/required.yaml` | `request-parameter-became-required`, nivel 3 | Plantilla existente / 3 |
| Unión con null | `openapi31/base.json` → `openapi31/union-null.json` | `request-parameter-became-not-nullable`, nivel 3 | Plantilla existente / 3 |
| Propiedad de respuesta | `openapi31/base.json` → `openapi31/response-property.json` | `response-property-type-changed`, nivel 3 | Fallback neutral / 3 |

Todos los hallazgos observados son de GET `/items`, sección `paths`, severidad `error` y compatibilidad `breaking`; son datos del motor fijado, no una reclasificación local. `--fail-on none` devuelve 0 con el mismo contenido en los tres formatos.

El caso de unión cambia **`type: ["string", "null"]` a `type: "string"`**, manteniendo `required: false`. Se comprobó con el motor real, sin extrapolar el `nullable` de 3.0. Solo demuestra este par concreto: no prueba `oneOf`, `anyOf`, uniones generales ni cómo serializa un servidor un parámetro query nulo. La [especificación 3.1.0, Schema Object](https://spec.openapis.org/oas/v3.1.0.html#schema-object) es la referencia de dialecto, no un resultado de nuestro motor.

La propiedad `id` de respuesta 200 pasa de `string` a `integer`, sin retirarla ni cambiar su obligatoriedad. La regla no tiene plantilla local y conserva el fallback; no se añade una explicación ni se infiere comportamiento del servidor. Los dos casos sin hallazgos muestran el mismo aviso de alcance: `[]` no demuestra que los contratos sean idénticos.

Posiciones contrastadas, base → head (línea:columna): obligatoriedad **12:11 → 12:11**; unión **11:41 → 11:22**; propiedad **21:41 → 21:41**. Las rutas de máquina no se versionan en las expectativas.

## Comprobaciones automatizadas

`TestOpenAPI31CorpusIntegration`, con `OASDIFF_BINARY` definido:

- Ejecuta los cinco pares en el motor real; exporta cada uno a JSON, HTML y Markdown UTF-8 sin BOM y con LF.
- Comprueba los códigos 0/3, el diagnóstico de política y que `none` no cambie el contenido.
- Contrasta las expectativas registradas y conserva todos los campos JSON y el evento upstream frente al informe enriquecido del motor. La comparación conserva el orden de arrays; ignora únicamente formato JSON y escapes Unicode equivalentes.
- Comprueba explicación/fallback, operación, regla, ambas fuentes y sus enlaces locales en los formatos humanos. Las entradas copiadas a un directorio con espacios, Unicode y `#` permanecen idénticas.
- Mantiene la suite anterior: YAML inválido, referencias externas bloqueadas, errores operativos, política, protección de archivos y escape/CSP. No se regeneraron snapshots.

Desde la raíz, con las herramientas locales fijadas según [DEVELOPMENT.md](DEVELOPMENT.md):

```powershell
$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
$env:GOMODCACHE = Join-Path (Get-Location) 'dist/go-mod-cache'
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
go test -count=1 ./internal/cli -run TestOpenAPI31CorpusIntegration -v
./scripts/verify.ps1
./scripts/smoke.ps1
```

Sin `OASDIFF_BINARY` el test informa que se omite; `verify.ps1` rechaza esa omisión. No hace falta una prueba manual para comprobar estos resultados: la regresión está automatizada y no cambia la presentación. Opcionalmente puede exportarse `base.json` frente a `union-null.json` usando la CLI y un destino nuevo, para ver un hallazgo de nivel 3 y exit 3.

## Revisión de dependencias de plataforma

- `internal/oasdiff/client.go` usa rutas locales y rechaza URL/UNC; el binario Windows y los scripts usan `.exe`/PowerShell. La CI existente verifica solo Windows x64.
- `internal/explain/locations.go` contempla unidades Windows, comparación sin distinguir mayúsculas en esa plataforma y fuentes del motor que omiten la unidad o vienen escapadas como URI. Los tests nuevos verifican el enlace completo, no una coincidencia por nombre de archivo.
- Los enlaces `file:///...#L<n>` conservan su política local; el salto a línea depende del visor/editor. No se verifican navegadores o editores de otras plataformas.
- Se usan `filepath` y temporales del sistema en los tests. Eso no demuestra soporte multiplataforma; el test específico de unidades Windows ya se omite fuera de Windows. No se añaden matrices CI, binarios ni nuevas dependencias.

## Evidencia — 2026-11-22

Observación directa con `changelog --config <archivo local con {}> --format json --lang en --level INFO --allow-external-refs=false --open=false --fetch=false -- BASE HEAD`, y después mediante el adaptador y la CLI. Un primer intento de observación usó el nombre incorrecto `--external-refs`; el motor lo rechazó con exit 100, sin producir resultados. Se corrigió a `--allow-external-refs=false` y los cinco pares se ejecutaron correctamente antes de fijar expectativas.

El test dirigido de los cinco pares y quince exportaciones pasó. `scripts/verify.ps1`: PASS (formato, todos los tests con motor real sin caché, vet y build). `scripts/smoke.ps1`: PASS en PowerShell 7.6.9 y Windows PowerShell 5.1.26100.9818, incluidos errores, códigos 0/1/2/3 y protección de archivos. La CI del SHA exacto se registra en la PR. Las evidencias de v0.1.0 permanecen históricas y no cambian sus archivos publicados. No se ejecuta el spike #25 como parte de este issue.
