# Desarrollo — openapi-contract-lab

## Estado

Módulo `github.com/v1kz25/openapi-contract-lab`, entrada en `cmd/openapi-contract-lab` y manejo de comandos en `internal/cli`. El issue #3 añade `internal/oasdiff` y `compare`: [modelo y verificación](OASDIFF-ADAPTER.md). El issue #4 incorpora `internal/explain`: [reglas y explicaciones](EXPLANATIONS.md). El issue #5 añade `internal/report`: [informes locales Markdown/HTML](REPORTS.md). El issue #6 incorpora la [política de severidad](SEVERITY-POLICY.md). El issue #7 añade [CI y regresión](CI.md).

## Entorno propuesto

- Windows 11; Git, GitHub CLI (`gh`), Codex en desktop/VS Code/CLI.
- Go 1.27.4 fijado en `go.mod`, sin módulos externos ni `go.sum` necesario. El adaptador sigue la decisión del [spike oasdiff](OASDIFF-SPIKE.md): binario local 1.39.0, Apache-2.0.
- Sin hosting, cuentas de servicios comerciales, API keys ni infraestructura de pago.

## Flujo de trabajo

1. Leer AGENTS.md, roadmap y el issue activo.
2. Crear `feature/<issue>-<nombre>` desde `develop` para una única tarea (Gitflow).
3. Implementar sólo su alcance y tests; revisar dependencias/licencias.
4. Ejecutar y registrar pruebas; comprobar Git diff y la identidad Git.
5. Proponer Pull Request; no fusionarla sin permiso.
6. Actualizar el README únicamente para funciones realmente disponibles.
7. Incluir en cada PR los pasos de comprobación manual y sus resultados esperados cuando apliquen. Avisar al usuario cuando deba probar un flujo; arrancar la aplicación local para facilitarlo cuando exista dicho flujo. La ayuda/versión inicial puede comprobarse opcionalmente desde terminal.

## Compatibilidad del contrato v1

Decisión del [issue #37](https://github.com/v1kz25/openapi-contract-lab/issues/37): **contrato normativo JSON v1 ligado a CLI 1.x, documentado sin identificador embebido**. No se añade `schemaVersion`, comando, dependencia ni cambio de comportamiento. No hay necesidad demostrada de negociar varios esquemas dentro del mismo ejecutable; añadir un discriminador ahora no aportaría esa negociación. La referencia está en [OASDIFF-ADAPTER.md](OASDIFF-ADAPTER.md#contrato-json-v1--issue-37) y la [matriz pública CLI](SEVERITY-POLICY.md#referencia-pública-cli--contrato-v1-issue-37). Las garantías de compatibilidad 1.x entrarán en vigor al publicar v1.0.0, no antes: el ejecutable actual sigue identificándose como v0.2.0 experimental.

Migración desde v0.2.0: no cambia la estructura, defaults ni códigos. Los informes antiguos sin campo de esquema siguen leyéndose con los mismos campos. Un consumidor que archiva informes debe guardar fuera del JSON la versión CLI obtenida con `--version` y sus expectativas de contrato; un informe aislado sin esa procedencia no certifica qué versión de la CLI lo generó. No interpretar la ausencia de `schemaVersion` como una versión desconocida numérica, `false` o `engineVersion`; tampoco distinguir futuras versiones incompatibles mediante heurísticas de campos. Si aparece necesidad real de esquemas múltiples, diseñar antes una migración explícita mediante issue/PR independiente.

| Cambio en CLI 1.x / JSON v1 | Política |
| --- | --- |
| Campo opcional adicional, que los consumidores pueden ignorar | Compatible si no cambia campos, significados, orden significativo o defaults existentes; documentar y probar con consumidor anterior. |
| Retirada/renombrado de campo, cambio de tipo o significado, nuevo campo requerido para consumir un informe, cambio de defaults/streams/códigos/gramática documentada | Ruptura: major de CLI/contrato y guía de migración; no introducir silenciosamente en 1.x. |
| Nuevo valor de enum propio o reclasificación local | No tratar automáticamente como adición inocua: los valores conocidos son parte del contrato. Requiere major/migración si cambia decisiones de consumidores. IDs y enteros upstream ya son abiertos; `unknown` es su manejo explícito. |
| Prosa de ayuda/errores/explicaciones, espacios/orden de claves/escapes JSON equivalentes, presentación HTML/Markdown | No se congela el texto o los bytes como API; conservar significado, streams, datos y seguridad, con regresión revisada. |
| Mensajes, fingerprints, campos adicionales y resultados semánticos de oasdiff | Contenido upstream, no garantías propias de estabilidad semántica entre motores. Mantenerlo completo y revisar cambios, sin sustituir el ID por prosa ni reclasificarlo localmente. |

### Actualización de oasdiff y herramientas

La comparación exige oasdiff **1.39.0**. No se admite una versión arbitraria, descarga automática, actualización silenciosa ni cambio de motor por parte del consumidor. Para actualizar:

1. Abrir issue con necesidad, versión exacta propuesta, alcance y compatibilidad Windows x64. Revisar notas upstream y versiones/casos OpenAPI realmente verificables; no convertir sus declaraciones generales en soporte propio.
2. Revisar licencia y componentes realmente distribuidos; comprobar hash oficial del archivo y registrar procedencia. Mantener coste obligatorio 0 €; motor separado y sin nuevas APIs/servicios.
3. En una rama/PR propia, ejecutar el corpus existente con motor antiguo y propuesto; registrar diferencias de IDs, niveles/clasificación, campos/tipos, orden, fuentes/fingerprints y mensajes. Comprobar los tres formatos, fallback, gate, errores/aislamiento y smoke Windows PS5.1/7. No regenerar expectativas en bloque para aceptar divergencias.
4. Evaluar por separado contrato propio y semántica del motor. Si el cambio rompe garantías propias, major/migración antes de integrar. Una actualización semántica aceptada y compatible con el contrato requiere notas explícitas y, como mínimo, una versión minor propia; no ocultarla en un patch. El fallback debe seguir siendo neutral hasta verificar cada regla/nivel con el motor propuesto.
5. Revisar diff, versiones/hashes/documentación/avisos y evidencia real del SHA. Actualizar únicamente después de revisión y aprobación de integración; la publicación requiere autorización propia según RELEASE.md. Go y acciones CI también se actualizan por necesidad y PR explícitas, con regresión, sin alterar versiones como parte de #37.

### Cobertura contractual del issue #37

Se reutilizan tests de ayuda/versión/argumentos, matriz de políticas, snapshots, Unicode, fallback, corpus 3.0.3/3.1 y exportaciones existentes. Se añaden pruebas de consumidor sobre la serialización pública (campos y tipos requeridos, omisiones, arrays vacíos, fuentes, orden/repeticiones y datos upstream desconocidos) y del uso real de PATH/defaults, `--opción=valor` y separador `--`. Cancelación activa, límites de pipes y recuperación de escrituras parciales pertenecen a #38 y no se implementan aquí. No se publica una versión estable ni se cambia el número de versión.

## Tests y comandos

Desde la raíz, con Go 1.27.4 en PATH. Estas variables solo afectan a la sesión PowerShell; las pruebas no necesitan descargas:

```powershell
$env:GOTOOLCHAIN = 'local'
$env:GOPROXY = 'off'
$env:GOSUMDB = 'off'
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
gofmt -l cmd internal
go test ./...
go vet ./...
go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
./dist/openapi-contract-lab.exe --help
./dist/openapi-contract-lab.exe --version
./dist/openapi-contract-lab.exe diff
$LASTEXITCODE # 2: diff aún no es un comando
```

`gofmt -l` no debe listar archivos. Ayuda/version y comparaciones que superan la política terminan con código 0; argumentos inválidos con 2; fallos de escritura, entrada o motor con 1. Una comparación completada que no supera `--fail-on` devuelve 3, conservando el informe. El adaptador proporciona `*oasdiff.Error` con etapa y causa original.

Para incluir las pruebas con el motor real (no se descarga durante los tests):

```powershell
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
go test ./...
```

Sin esta variable se ejecutan los tests unitarios y golden; `TestIntegration` se omite con un mensaje explícito. Las pruebas reales usan las entradas sintéticas del spike, prueban errores de YAML/referencias externas y comprueban aislamiento de configuración. No es necesario configurar CI para ejecutarlas.

Con Go/gofmt en PATH y `OASDIFF_BINARY` definido, `./scripts/verify.ps1` ejecuta las comprobaciones usadas en CI y rechaza omitir el motor. `.gitattributes` fija LF para los archivos Go; si un árbol anterior ya tiene CRLF, normalizarlo con `gofmt -w cmd internal` antes de comprobar el formato.

El issue #8 añade [instalación reproducible](INSTALL.md), [preparación de publicación](RELEASE.md) y `./scripts/smoke.ps1`. El smoke usa el ejecutable compilado y oasdiff local, sin SDK ni descargas; se ejecuta también en CI después de compilar. No crea tags, releases ni paquetes.

Los issues #18/#33 están completados: v0.1.0 y [v0.2.0 experimentales](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.2.0) publicadas bajo Apache-2.0. PR #34 integrada en `main`, #35 sincronizada a `develop`. Sus notas/artefactos y la [preparación histórica #33](CANDIDATE-V0.2.0.md) se conservan. `develop` incorpora después contrato #37 y pruebas/corrección de diagnóstico #38; su código no es el ZIP publicado, aunque `--version` siga mostrando v0.2.0 hasta preparar otra versión. No autoriza nuevas fusiones, tags ni publicación.

Con la misma variable también se ejecutan los trece casos sintéticos de `internal/explain`, uno por combinación documentada (diez originales y tres del issue #21), y se comprueba que el motor emite los IDs/niveles esperados y las fuentes enlazadas. `internal/explain/testdata/breaking/` contiene los seis contratos nuevos, también usados por la integración CLI de JSON/Markdown/HTML y política. Las pruebas de control verifican entradas idénticas e inversas; no se infiere la regla inversa. Los snapshots de mensajes se revisan en `internal/explain/testdata/`; no hay actualización automática de expectativas.

`internal/report/testdata/` contiene snapshots revisables de HTML y Markdown para un informe mixto y otro vacío. Los tests de escape y enlaces se ejecutan sin motor. Con `OASDIFF_BINARY`, los tests de CLI también exportan JSON/Markdown/HTML reales a archivos UTF-8 en un directorio temporal. La exportación usa creación exclusiva y no sobrescribe archivos existentes.

El issue #23 añade agrupación en la vista HTML/Markdown, snapshots globales y [ejemplos antes/después](REPORTS.md#ejemplos-comparables-y-revisión-manual-del-issue-23). Los ejemplos agrupados son snapshots comprobados por tests, sin actualización automática. Las pruebas verifican orden, recuentos, conservación de ocurrencias y campos, encabezados hostiles, CSP y JSON sin mutaciones; la integración CLI compara varios endpoints/métodos con el motor real y conserva JSON y gate.

El issue #24 amplía la [regresión OpenAPI 3.1](OPENAPI31-REGRESSION.md): `TestOpenAPI31CorpusIntegration` contrasta cinco pares, reutiliza los tres casos iniciales y añade unión con null y tipo de propiedad de respuesta. `internal/cli/testdata/openapi31/cases.json` contiene expectativas observadas con oasdiff 1.39.0, nunca regeneradas automáticamente. Los tres formatos, fallback, códigos 0/3, fuentes con espacios/Unicode/# y conservación de entradas se comprueban con el motor real al ejecutar verify.

El [spike #25](CHANGE-KIND-SPIKE.md) se reproduce por separado con `./spikes/change-kind/verify.ps1` (PowerShell 7 y motor local). Sus cinco pares, quince salidas congeladas y 25 controles no cambian código ni contrato público. El runner genera nueva evidencia en un directorio nuevo bajo `dist/`, comprueba expectativas sin actualizarlas e incluye configuración ambiental inválida y refs externas bloqueadas. La CI existente no ejecuta este runner de investigación automáticamente.

### SDK portátil en Windows x64

Si no hay Go instalado, obtener el [ZIP oficial Go 1.27.4](https://go.dev/dl/go1.27.4.windows-amd64.zip), cuyo SHA-256 publicado en [go.dev](https://go.dev/dl/?mode=json) es `178ab98950ac69ab25069307dc6ade79f3816e414cce95b4da166a115cde1e7e`. Comprobarlo antes de extraer:

```powershell
New-Item -ItemType Directory -Force dist | Out-Null
Invoke-WebRequest 'https://go.dev/dl/go1.27.4.windows-amd64.zip' -OutFile dist/go1.27.4.windows-amd64.zip
$sdkHash = (Get-FileHash dist/go1.27.4.windows-amd64.zip -Algorithm SHA256).Hash
if ($sdkHash -ne '178ab98950ac69ab25069307dc6ade79f3816e414cce95b4da166a115cde1e7e') { throw 'Go checksum mismatch' }
Expand-Archive dist/go1.27.4.windows-amd64.zip -DestinationPath dist/_tools/go1.27.4
$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path
go version
```

No requiere instalación global, privilegios de administrador ni oasdiff para ayuda/versión. La descarga inicial necesita Internet; compilación y pruebas usan la biblioteca estándar y pueden hacerse offline. El SDK conserva su [licencia BSD-3-Clause](https://github.com/golang/go/blob/go1.27.4/LICENSE) en `dist/`, que permanece ignorado por Git. Coste adicional obligatorio: 0 €. No se añade framework CLI, servicio, CI ni dependencia de terceros.

El SDK se extrae bajo `_tools` porque `go test ./...` y `go vet ./...` ignoran los directorios que empiezan por `_`. `.gitignore` solo controla Git: extraer el SDK directamente en `dist/go1.27.4` haría que Go intentase recorrer también los tests del compilador.

### Validación histórica del issue #2 — 2026-11-21

Ejecutada en Windows x64, build 26300.9974, PowerShell 7.6.9, Go 1.27.4 portátil con checksum comprobado. No se modificó el PATH global. Con `GOTOOLCHAIN=local`, `GOPROXY=off` y `GOSUMDB=off`:

| Comando | Resultado real |
| --- | --- |
| `go list ./...` | Solo los paquetes `cmd/openapi-contract-lab` e `internal/cli` |
| `go test ./...` | PASS, exit 0; `internal/cli` en 0.513 s, entrada main sin tests propios |
| `go vet ./...` | Sin diagnósticos, exit 0 |
| `gofmt -l cmd internal` | Sin archivos pendientes, exit 0 |
| `go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab` | Binario Windows generado, exit 0 |
| Ejecutable con `--help` / `--version` | Ayuda / `openapi-contract-lab dev`, exit 0 |
| Ejecutable con `diff` | `unknown command "diff"; use --help for usage`, stderr, exit 2 esperado |

Los tests cubren las cuatro formas de ayuda (incluida invocación sin argumentos), dos de versión, cinco entradas inválidas y un fallo de escritura simulado. No invocan oasdiff ni servicios. El primer `go test ./...` falló al recorrer el SDK extraído directamente en `dist/`; se movió a `dist/_tools/` y se repitió la verificación completa correctamente. No se oculta ese fallo de preparación. No se probó ARM64 ni Linux/macOS.

## Mantenimiento y reporte de problemas

Mantenimiento individual, sin SLA universal, soporte continuo garantizado ni fecha comprometida para v1. Priorizar bloqueantes del [roadmap](ROADMAP.md), regresiones y seguridad dentro del ámbito Windows x64/corpus declarado. Mantener operación local/offline y coste obligatorio 0 €; no añadir servicios, telemetría, soporte de otros sistemas o reglas por defecto al corregir un fallo.

Para comunicar un problema en [GitHub Issues](https://github.com/v1kz25/openapi-contract-lab/issues), incluir versión CLI y motor, Windows/PowerShell, comando mínimo con rutas ficticias, código de salida esperado/real y dos contratos **sintéticos** que lo reproduzcan. Añadir solo el fragmento de JSON/diagnóstico imprescindible; retirar rutas personales, nombres internos, tokens y contenido sensible también de `upstream`/mensajes. No adjuntar especificaciones privadas ni informes completos por defecto. Si no puede compartirse sin datos sensibles, describir el impacto sin adjuntos y acordar un canal privado con el propietario antes de enviar material; no se presupone infraestructura de recepción privada.

Antes de preparar un candidato y cuando aparezca un aviso relevante o un fallo reproducible, revisar manualmente releases/avisos oficiales de Go y oasdiff, licencias, hashes y acciones CI realmente usadas. Son consultas de mantenimiento, no un autoupdater ni una monitorización programada. La [política de actualización](#actualización-de-oasdiff-y-herramientas) exige issue/PR y comparación del corpus: comprobar pertinencia del aviso a nuestro uso y componentes distribuidos, no actualizar a ciegas ni prometer ausencia absoluta de vulnerabilidades. Registrar aplazamientos con su motivo y limitación real.

Un fallo corregido dentro del contrato 1.x puede publicarse como patch tras revisión; una actualización semántica upstream compatible requiere al menos minor y notas, una ruptura propia exige major/migración. Mientras v0.2.0 siga experimental, documentar igualmente cambios de comportamiento/migración sin atribuirle estabilidad retrospectiva. El [proceso de publicación](RELEASE.md) y la aprobación explícita son obligatorios para cada versión; [instalación/update/rollback](INSTALL.md#actualizar-y-volver-a-la-versión-anterior) conserva carpetas anteriores. No prometer que actualizar/revertir transforma informes ya archivados.

## Gitflow

- `main`: versiones estables; rama predeterminada. No trabajar directamente en ella después de la preparación inicial.
- `develop`: integración de la próxima versión; punto de partida de tareas y spikes.
- `feature/<issue>-<nombre>`: una tarea pequeña o spike del issue activo, creada desde `develop`; PR hacia `develop`. No avanzar al siguiente issue.
- `release/<version>`: estabilización de una versión desde `develop`; solo correcciones, pruebas y documentación. PR hacia `main` y sincronización posterior hacia `develop` mediante otra PR.
- `hotfix/<issue>-<nombre>`: corrección urgente de una versión publicada, desde `main`; PR hacia `main` y otra hacia `develop` para propagar la corrección.
- Preparación de ramas y PRs autorizada dentro del issue activo. Toda fusión requiere aprobación explícita del usuario. No habilitar auto-merge.
- No hacer force push, rebase de ramas compartidas, ni eliminar ramas. Si la interfaz ofrece eliminar la rama tras fusionar, desactivar esa opción.
- Usar Git y GitHub CLI estándar; no se necesita instalar una extensión Gitflow ni añadir servicios o pipelines.
- Antes de una PR: revisar diff, pruebas relevantes, identidad Git y ausencia de secretos. Registrar resultados y bloqueos reales.
