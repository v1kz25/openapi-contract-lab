# Adaptador local — issue #3

Alcance: [issue #3](https://github.com/v1kz25/openapi-contract-lab/issues/3), basado en el setup integrado por la PR #11. Entrada base/head de archivos locales y salida JSON normalizada. El motor sigue siendo oasdiff **1.39.0**; Go **1.27.4**, biblioteca estándar, coste adicional obligatorio **0 €**. No se añade dependencia Go, servicio ni motor de diff propio. El ejecutable mantiene la licencia [Apache-2.0](https://github.com/oasdiff/oasdiff/blob/v1.39.0/LICENSE); descarga/checksum y compatibilidad Windows comprobados en el [spike](OASDIFF-SPIKE.md).

## Comportamiento

`compare [--oasdiff PATH] BASE HEAD` comprueba que ambas entradas sean archivos regulares legibles, localiza el ejecutable y exige `oasdiff version 1.39.0`. Después ejecuta `changelog --format json --lang en --level INFO`. Una configuración temporal vacía explícita y la exclusión de variables `OASDIFF_*` aíslan la configuración del motor. Se pasan `--allow-external-refs=false`, `--open=false` y `--fetch=false`. No hay shell intermedia ni aceptación de URLs, stdin o rutas UNC como entradas.

El contexto limita versión y comparación a 30 segundos en total. La captura de stdout está limitada a 16 MiB y stderr a 64 KiB; los fallos devuelven un error, no un informe parcial. `*oasdiff.Error` conserva etapa y causa (entrada, ejecutable, versión, configuración, comparación o decodificación). Los errores de proceso incorporan su diagnóstico stderr. No se descarga ni instala oasdiff automáticamente.

La CLI devuelve 0 si el informe supera la política, 3 si se genera pero no supera `--fail-on` (por defecto `error`), 1 si falla la entrada/motor/salida y 2 para argumentos inválidos. **Un finding breaking no activa exit 1**: con la política predeterminada activa exit 3. `--fail-on none` permite exportarlo con exit 0. Ver [política vigente](SEVERITY-POLICY.md).

## Modelo y fidelidad

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

`changelog` incluye cambios que afectan a consumidores, no cada diferencia del documento: una modificación solo de descripción produce `findings: []`. Esto **no significa contratos idénticos**. No se llama a `diff` ni se crean hallazgos documentales propios. Los textos todavía son los del motor; las explicaciones propias corresponden al issue #4. Tampoco se afirma cobertura completa de OpenAPI 3.1/3.2 con estos fixtures 3.0.3.

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
