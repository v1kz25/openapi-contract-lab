# Spike oasdiff — issue #1

Verificado el 2026-11-21. Alcance: [issue #1](https://github.com/v1kz25/openapi-contract-lab/issues/1), abierto y con título «Validar integración local de oasdiff y formatos de salida». Sin comentarios adicionales. La [PR Gitflow #9](https://github.com/v1kz25/openapi-contract-lab/pull/9) estaba fusionada hacia `develop` antes de empezar. Rama `feature/1-oasdiff`, creada con árbol limpio desde `develop` actualizado, commit `3b76fbee1af9734743630b53ad263d3459ddb93e`.

## Decisión

Elegir **binario local oasdiff v1.39.0 + JSON** como frontera de integración inicial para la futura aplicación Go. Usar `breaking` para incompatibilidades y `diff` cuando se necesiten diferencias documentales. `changelog` no recoge todos los cambios no breaking: una descripción modificada no genera un hallazgo suyo. No deducir una regla propia a partir de un diff ni interpretar una lista vacía como contratos idénticos.

| Opción | Ventajas | Costes y riesgos observados |
| --- | --- | --- |
| Módulo Go fijado a v1.39.0 | Objetos tipados, sin subproceso ni binario separado; acceso a checker y localizaciones | Acoplamiento a API Go y parser; requisitos de toolchain y dependencias transitivas. La release elimina una función pública (`diff.SchemaAcceptsNull`): la API cambia incluso entre releases 1.x. Compilación no verificada aquí |
| Binario fijado a v1.39.0 — elegido | Ejecutado realmente en Windows x64, sin SDK Go para el spike; mismo motor semántico; contrato JSON y esquema inspeccionables | Distribuir/resolver ejecutable por arquitectura; controlar versión, códigos de salida y errores; formato y semántica requieren regresión al actualizar |

La API Go existe: `load.NewSpecInfo`, `diff.GetPathsDiff` y `checker.CheckBackwardCompatibility` aparecen en el [ejemplo oficial fijado](https://github.com/oasdiff/oasdiff/blob/v1.39.0/checker/example_test.go). `diff.Get` por sí solo no clasifica breaking changes. Este spike inspecciona esa ruta, pero no la compila ni implementa un adaptador. La elección queda sustentada por la ejecución del binario; no se afirma equivalencia comprobada entre ambas integraciones.

## Versiones y entorno

- [Última release estable consultada](https://github.com/oasdiff/oasdiff/releases/tag/v1.39.0): **v1.39.0**, publicada el 2026-11-17. Binario ejecutado: `oasdiff version 1.39.0`.
- Artefacto: `oasdiff_1.39.0_windows_amd64.tar.gz`. SHA-256 calculado y coincidente con el digest de GitHub: `956a73df88e03efdf3fcc5f327446cc176fefe4adcea1978b0fef464a1ee8901`.
- PowerShell **7.6.9**, Git **2.56.0.windows.4**, gh **2.104.0**. Host x64, Windows NT **10.0.26300.0**; registro `DisplayVersion=26H2`, build **26300.9974**. La API de compatibilidad/registro devuelve «Windows 10» y CIM fue denegado; se conserva la identificación real, sin inventar una edición de Windows 11. El binario funciona en este host; ARM64 no se probó.
- `Get-Command go,gh,oasdiff` solo encontró gh; `go version` falló por comando no encontrado. **No se pudo compilar ni ejecutar tests Go**; no se instaló un SDK. El [go.mod fijado](https://github.com/oasdiff/oasdiff/blob/v1.39.0/go.mod) declara `go 1.26.0` y `toolchain go1.27.2`. [Go estable actual consultado](https://go.dev/dl/?mode=json): **1.27.4**, candidato fijado para una futura comprobación de compilación, no instalado ni validado aquí.
- El [build oficial](https://github.com/oasdiff/oasdiff/blob/v1.39.0/.goreleaser.yml) incluye Windows amd64/arm64 con `CGO_ENABLED=0`. No se necesita Docker ni un compilador C para ejecutar el binario.

## Dos ejemplos reproducibles

Fixtures propios y sintéticos OpenAPI **3.0.3**, sin servidores, credenciales ni referencias externas. Base compartida: `spikes/oasdiff/fixtures/base.yaml`.

| Revisión frente a base | Única edición | Resultado real |
| --- | --- | --- |
| `required.yaml` | `filter`, query de `GET /items`: `required: false` → `true` | Un hallazgo `request-parameter-became-required`, `level: 3` (ERR), localizado en línea 12 de ambas entradas |
| `description.yaml` | `Filter items.` → `Filter items by name.` | Sin hallazgos breaking ni changelog; `diff` conserva el cambio documental |

Preparación, desde la raíz del repositorio en PowerShell 7 x64 (solo la descarga requiere Internet):

```powershell
New-Item -ItemType Directory -Force dist/oasdiff-1.39.0 | Out-Null
Invoke-WebRequest 'https://github.com/oasdiff/oasdiff/releases/download/v1.39.0/oasdiff_1.39.0_windows_amd64.tar.gz' -OutFile dist/oasdiff-1.39.0/oasdiff.tar.gz
$archiveHash = (Get-FileHash dist/oasdiff-1.39.0/oasdiff.tar.gz -Algorithm SHA256).Hash
if ($archiveHash -ne '956a73df88e03efdf3fcc5f327446cc176fefe4adcea1978b0fef464a1ee8901') { throw 'Checksum mismatch' }
tar -xzf dist/oasdiff-1.39.0/oasdiff.tar.gz -C dist/oasdiff-1.39.0
./dist/oasdiff-1.39.0/oasdiff.exe --version
./spikes/oasdiff/verify.ps1
```

El script no descarga nada: usa el ejecutable ya disponible y genera evidencia en `dist/oasdiff-evidence` (ignorado por Git). Acepta `-Binary` y `-OutputDirectory`. Ejecución directa de cada ejemplo:

```powershell
./dist/oasdiff-1.39.0/oasdiff.exe breaking spikes/oasdiff/fixtures/base.yaml spikes/oasdiff/fixtures/required.yaml --config spikes/oasdiff/offline.yaml -f json --fail-on ERR
# exit 1; id=request-parameter-became-required, level=3
./dist/oasdiff-1.39.0/oasdiff.exe breaking spikes/oasdiff/fixtures/base.yaml spikes/oasdiff/fixtures/description.yaml --config spikes/oasdiff/offline.yaml -f json --fail-on ERR
# exit 0; []
```

## Comparación de salidas reales

| Salida | `breaking` / `changelog` | `diff` | Adecuación |
| --- | --- | --- | --- |
| JSON | Array de hallazgos: `id`, `text`, `level` numérico, operación, ruta, sección, fuentes y fingerprint en el caso breaking; `[]` para descripción | Objeto anidado `paths.modified...`, valores `from`/`to`; sin severidad ni ID de regla | Elegido para lectura por Go con biblioteca estándar. `oasdiff schema` ofrece JSON Schema 2020-12 para breaking/changelog, no para el diff completo |
| YAML | Misma estructura de hallazgos observada; `[]` para descripción | Incluye `paths` **y también `endpoints`**, con clave compuesta método/ruta. JSON omite esta segunda sección | Legible, pero necesita parser adicional y el diff no es intercambiable con JSON |
| Markdown | Agrupa por operación; caso breaking con icono `:warning:` pero sin ID, nivel numérico ni fuentes. Para descripción avisa que los contratos difieren sin hallazgos | Lista el parámetro modificado y sus valores anteriores/nuevos | Útil para revisión humana; no usar iconos o frases para clasificar ni parsear como contrato de máquina |

Las salidas están en [spikes/oasdiff/evidence](../spikes/oasdiff/evidence); [commands.json](../spikes/oasdiff/evidence/commands.json) registra cada comando, exit code y stderr. Los `.md` de esa carpeta son salida del motor, no documentación redactada. Se normalizan CRLF a LF y se deja un único salto de línea al final; no se modifica el contenido de los hallazgos.

**Estabilidad:** las dos ejecuciones con esta versión reprodujeron los 24 archivos de evidencia exactamente tras la normalización indicada. No prueba estabilidad entre versiones. Fijar versión y aceptar campos opcionales/nuevos en la futura integración; no usar el texto localizado como identificador. El esquema solo tipa `level` como entero; los valores 3/2/1 se contrastaron con [checker/rules/level.go](https://github.com/oasdiff/oasdiff/blob/v1.39.0/checker/rules/level.go). El [formatter](https://github.com/oasdiff/oasdiff/blob/v1.39.0/formatters/changes.go) omite campos vacíos y no serializa `IsBreaking`. Un fingerprint no sustituye al ID de regla ni a la severidad.

## Dependencias, licencias y coste

No se añade módulo Go ni dependencia de aplicación. El ejecutable descargado y su LICENSE quedan en `dist/`, sin versionar ni redistribuir. No se copió código de implementación del upstream.

| Componente inspeccionado | Versión fijada por upstream | Licencia comprobada | Necesidad |
| --- | --- | --- | --- |
| [oasdiff LICENSE](https://github.com/oasdiff/oasdiff/blob/v1.39.0/LICENSE) | 1.39.0 | Apache-2.0 | Único motor de diff/clasificación |
| [kin-openapi LICENSE](https://github.com/getkin/kin-openapi/blob/v0.153.0/LICENSE) | 0.153.0 | MIT | Parser usado por oasdiff; relevante al evaluar módulo Go |
| [go-yaml LICENSE](https://github.com/yaml/go-yaml/blob/v3.0.6/LICENSE) | go.yaml.in/yaml/v3 3.0.6 | MIT para archivos derivados de libyaml y Apache-2.0 para los restantes | Entrada/salida YAML en upstream |
| [Cobra LICENSE](https://github.com/spf13/cobra/blob/v1.10.2/LICENSE.txt) | 1.10.2 | Apache-2.0 | CLI upstream |

Es una revisión de los componentes relevantes, **no una auditoría completa del árbol transitivo**. El go.mod upstream lista más módulos, incluido `cloud.google.com/go`; su presencia no obliga a contratar servicios para ejecutar estos comandos locales. Antes de empaquetar/redistribuir un producto habrá que inventariar las dependencias realmente incluidas, sus licencias y avisos. No se ha decidido la licencia propia del repositorio.

Coste monetario adicional obligatorio: **0 €**. Ejemplos y verificador locales, sin cuentas comerciales, tarjetas, APIs de IA ni runners. Tras obtener el binario, los ejemplos funcionan sin servicios externos: rutas locales y configuración explícita `allow-external-refs: false`; no se usan `--open`, URLs de entrada ni `--fetch`. No se hizo una captura de tráfico ni una auditoría general de red del ejecutable.

## Verificación y límites

- Ejecutados `git fetch origin`, `git switch develop`, `git pull --ff-only origin develop` y `git switch -c feature/1-oasdiff`: base actualizada sin conflictos. El primer fetch fue denegado por el sandbox; el reintento autorizado funcionó. La descarga también necesitó ejecución fuera de la restricción DNS del sandbox.
- SHA-256 verificado y versión ejecutada; ayudas reales de `breaking`, `diff` y `validate` inspeccionadas; `oasdiff schema` ejecutado e inspeccionado.
- Cada ejecución del verificador: **23 comandos** (3 validaciones, 18 comparaciones de 2 casos × 3 subcomandos × 3 formatos, 2 pruebas `--fail-on ERR`). Sin stderr; validaciones `[]`, exit 0. Las comparaciones sin `--fail-on` salen con 0 incluso con breaking changes; con `--fail-on ERR`, required sale con 1 y description con 0.
- Assertions de ID, nivel, método/ruta, ausencia de hallazgos documentales y valores del diff: correctas. Segunda ejecución completa: **24 archivos idénticos**, incluido registro de comandos. Se compararon también los fixtures: una sola línea modificada en cada revisión.
- Diff completo y archivos staged revisados; `git diff --cached --check` correcto tras normalizar los saltos finales de la evidencia. Búsqueda de patrones de secretos y revisión manual sin hallazgos. Identidad Git existente de author/committer comprobada, sin modificarla ni añadir trailers de IA.
- No hay tests Go, compilación Go ni pruebas de integración por módulo ejecutadas: falta `go` en PATH. No se presenta este spike como producto terminado ni se exige al usuario una prueba manual de aplicación. La revisión del diff y repetición opcional de ejemplos pueden hacerse antes de la PR a `develop`; la prueba manual del flujo completo corresponderá cuando exista dicho flujo, también antes de publicar una release.
- Estos ejemplos no validan OpenAPI 3.1/3.2. La [documentación fijada](https://github.com/oasdiff/oasdiff/blob/v1.39.0/docs/OPENAPI-31.md) declara soporte 3.1 y lectura 3.2, con limitaciones en `$dynamicRef`/`$dynamicAnchor` y `components.pathItems`. No extrapolar cobertura completa ni compatibilidad entre dialectos a partir de esta prueba 3.0.3.

El issue queda técnicamente documentado con sus límites. Publicación de rama/PR hacia `develop` pendiente de autorización del usuario; no fusionar, no habilitar auto-merge ni borrar ramas. No avanzar al issue de inicialización.
