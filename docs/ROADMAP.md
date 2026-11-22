# Roadmap — openapi-contract-lab

## Estado publicado

[v0.2.0 experimental](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.2.0) está publicada desde el **2026-11-22**, como prerelease, no draft. PR [#34](https://github.com/v1kz25/openapi-contract-lab/pull/34) integrada en `main` y [#35](https://github.com/v1kz25/openapi-contract-lab/pull/35) sincronizada hacia `develop`. Incluye comparación local con oasdiff 1.39.0, trece explicaciones deterministas, JSON e informes HTML/Markdown agrupados, gate por severidad y cinco pares concretos OpenAPI 3.1.0. ZIP Windows x64 y checksums publicados; motor obtenido por separado. CLI/esquema siguen siendo experimentales: no se promete estabilidad ni cobertura completa de OpenAPI.

Apache-2.0, Go 1.27.4, oasdiff 1.39.0, sin módulos Go externos ni coste adicional obligatorio. Solo Windows x64 verificado, PowerShell 5.1 y 7. [v0.1.0](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.1.0), sus diez explicaciones y [notas](releases/v0.1.0.md) permanecen como historia. Las [notas v0.2.0](releases/v0.2.0.md) y [guía de candidato](CANDIDATE-V0.2.0.md) describen la preparación de #33: sus frases «no publicado» son anteriores a la publicación. #39 reconcilia guías activas de instalación, empaquetado y mantenimiento sin modificar bytes ni notas publicados. #37/#38 ya están integrados en `develop`; la siguiente preparación de candidato corresponde a #40 después de revisar #39. Cada futura integración/publicación requiere aprobación independiente según [RELEASE.md](RELEASE.md).

## Candidato experimental v0.2.0 — trabajo completado

Este backlog histórico está cerrado; no son tareas nuevas para v1. La documentación #22, las tres explicaciones #21, los informes agrupados #23 y la [regresión 3.1 con cinco pares concretos](OPENAPI31-REGRESSION.md) del #24 se integraron mediante las PRs #26, #27, #28 y #31. El [spike #25](CHANGE-KIND-SPIKE.md), integrado en #32, concluyó reducir la propuesta a una señal de diferencias detectadas/no detectadas/desconocido y posponer su integración pública; no demuestra categorías generales de igualdad o cambio solo documental. La clasificación no está implementada.

| Orden / issue | Valor y alcance | Dificultad / dependencias | Riesgo y criterio de corte |
| --- | --- | --- | --- |
| Preparación: [#22, documentación](https://github.com/v1kz25/openapi-contract-lab/issues/22) | Reflejar v0.1.0 publicada y corregir ejemplos activos de códigos de salida. | Baja; publicación y sincronización ya realizadas. | Separar estado vigente de evidencias históricas; sin cambiar código. |
| 1: [#21, tres explicaciones](https://github.com/v1kz25/openapi-contract-lab/issues/21) | Tipo incompatible de parámetro, cuerpo obligatorio y retirada de propiedad obligatoria de respuesta; impacto/acción en español. | Baja; ampliar el mecanismo existente, catálogo y fixtures reales 1.39.0. Independiente de #22. | Máximo tres pares regla/nivel comprobados; mantener fallback y todos los datos upstream. |
| 2: [#23, informes agrupados](https://github.com/v1kz25/openapi-contract-lab/issues/23) | Revisar por endpoint/operación y severidad; sección para hallazgos globales. | Media; renderer existente, no depende del spike #25. Revisar primero la PR de #21. | Orden estable, ningún hallazgo perdido, escape y CSP intactos; no cambiar JSON ni política. |
| 3: [#24, regresión 3.1](https://github.com/v1kz25/openapi-contract-lab/issues/24) | Hasta cinco pares representativos, incluidos null/uniones y respuestas; confianza verificable. | Media; motor fijado y corpus existente, después de revisar #21. | No extrapolar nullable 3.0 ni afirmar soporte total; documentar incertidumbre y mantener fallback. |
| 4: [#25, validar tipo de cambio](https://github.com/v1kz25/openapi-contract-lab/issues/25) | Investigar idéntico/documental/contractual con las capacidades de oasdiff; evitar falsa igualdad al no haber hallazgos. | Media/alta; validar antes de diseñar una ampliación del adaptador/esquema. | Solo spike con cinco pares y decisión; si no hay señal fiable o exige gran cambio público, posponer. |

El [issue #33](https://github.com/v1kz25/openapi-contract-lab/issues/33) quedó cerrado: preparó identificación `v0.2.0`, notas, ZIP y una corrección acotada del directorio de salida del runner del spike. Las fusiones #34/#35 y la publicación ya ocurrieron. Se reutilizan su evidencia de reproducibilidad y la aceptación visual de #23; no se reabre ese trabajo.

## Auditoría actual — issue #36, 2026-11-22

Alcance de esta ejecución: auditoría y planificación exclusivamente en `v1kz25/openapi-contract-lab`. La única modificación versionada es este documento. Issues futuros no implementados; sin cambios de código, tests, workflows, herramientas, versiones o licencias.

### Estado real y procedencia

- Remoto `https://github.com/v1kz25/openapi-contract-lab.git`, repositorio público, rama predeterminada `main`. `git fetch origin --tags` completado; árbol inicial limpio, sin cambios ajenos que mezclar.
- Base `develop`/HEAD: **`aa194b224683161a79ff21da37e6b2711f54a842`**. `main`/tag v0.2.0 desreferenciado: **`72ebbd87558bd73ae138696fd9dc62a58b6e2c04`**. Sus árboles son iguales (`git diff origin/main origin/develop --stat` vacío), aunque los commits difieren por la sincronización. Tag v0.1.0 desreferenciado: `5d4b87d9cee3e5c2f089b3f401b6cc0a3b0ecd9a`.
- Releases v0.1.0 y v0.2.0 publicadas como experimentales; al empezar, ningún issue ni PR abierto. Se revisaron todos los issues/PRs existentes, especialmente #7/#8/#18/#21–#25/#33 y sus PRs, antes de crear seis issues #36–#41. #36 vincula esta planificación; rama `feature/36-v1-planning` desde `develop` actualizado, PR hacia `develop`.
- Identidad author/committer existente comprobada y conservada; sin configuración Git cambiada ni trailers de asistentes. No se publican logs locales con rutas personales.
- Leídos AGENTS, README, ROADMAP, DEVELOPMENT, RELEASE, INSTALL, CI, adaptador, explicaciones, informes, política, corpus 3.1, ambos spikes, guía de candidato y notas v0.1/v0.2; contrastados con `internal/cli`, `oasdiff`, `explain`, `report`, scripts, tests y workflow.

### Diagnóstico y contrato a estabilizar

**Base funcional suficiente para una v1 acotada; aún no lista para declararse estable.** No se observó fallo funcional en las verificaciones actuales. Los bloqueantes son compromiso contractual, cobertura concreta de fallos, instrucciones vigentes y validación/aprobación de un candidato futuro. Aprobar este plan no declara estable v0.2.0.

| Aspecto | Evidencia vigente / límite | Acción hacia v1 |
| --- | --- | --- |
| CLI | `cli.go`/`policy.go`: ayuda sin argumentos; `help/-h/--help`, `version/--version`, `compare`. Opciones antes de BASE/HEAD: `--oasdiff` por defecto PATH, `--format json`, `--output` ausente = stdout, `--fail-on error`. Códigos 0/1/2/3 y tests existentes. | #37 consolida referencia normativa y pruebas de consumidor de los huecos; no añadir comandos. |
| Streams y archivos | Informe solo en stdout; con output, stdout vacío. Diagnósticos en stderr. Informe completo antes de exit 3; exit 1 no garantiza salida completa. `O_EXCL` protege destinos existentes e inputs; no crea directorios padre. JSON con escapes ASCII/Unicode; output UTF-8 sin BOM/LF. | Congelar significados/defaults, no texto literal de diagnósticos. #38 caracteriza destino nuevo parcial tras fallo de escritura/cierre. |
| JSON | `model.go` + `explain.go`: engine/engineVersion/findings; base/head omitibles en el tipo (CLI normal aporta rutas absolutas). Finding: ruleId/severity/compatibility/upstreamLevel/upstream/explanation obligatorios; message/operation/operationId/path/section/sources/fingerprint omitibles. Source: file, line/column omitidos en cero. Explanation: kind/summary/impact/action/locations; Location: side/file y line/column/url opcionales. Arrays vacíos `[]`; sin campo de versión de esquema. | #37 especifica todos los tipos/enums/ausencias, datos opacos, evolución y migración. No equiparar ausencia a null/false. |
| Orden y consumidores | Orden de findings y locations significativo; upstream conserva valores/campos desconocidos, no necesariamente espacios/escapes. HTML/Markdown reagrupan sin mutar JSON. | Conservar orden recibido para motor fijado; no prometer orden de claves JSON, mensajes/fingerprints entre versiones ni prosa como API. Consumidores tolerantes a campos opcionales desconocidos. |
| Informes | `html/template`, CSP, escape Markdown, enlaces locales restringidos; snapshots mixto/vacío/global y orden/conservación/determinismo. Mismo informe = mismos bytes; rutas distintas pueden cambiar bytes. Aceptación visual #23 ya registrada. | Reutilizar suite y aceptación; no exigir otro renderer ni más explicaciones. Revisión visual nueva solo si cambia presentación. |
| Entradas y recursos | Dos archivos regulares locales, URLs/UNC/stdin y refs externas rechazados; configuración aislada. Deadline total versión+compare 30 s, WaitDelay 2 s; stdout motor 16 MiB, stderr 64 KiB. No hay cuota explícita de tamaño de entrada/RAM del motor ni sandbox. | #38 prueba proceso iniciado y límites de pipes reales: hoy solo cancelación previa y buffer unitario. Declarar archivos locales revisados como ámbito; no prometer seguridad ante entrada arbitraria hostil. |
| OpenAPI | Trece pares regla/nivel 3.0.3; cinco pares 3.1.0: idéntico, documental, obligatoriedad, `type: [string,null]` → string y tipo de propiedad de respuesta. Fallback verificado, no extrapolación de nullable 3.0. | Conservar corpus #24. No exigir más dialectos/casos sin necesidad; 3.1 completo, 3.2, oneOf/anyOf generales y comparación entre dialectos no verificados. |
| Motor | oasdiff exactamente 1.39.0; clasificación semántica suya, explicación accionable propia. `[]` no prueba igualdad ni ausencia de todo cambio contractual. #25 ya investigó la señal adicional, no implementada. | #37 fija actualización mediante issue/PR, hashes/licencia/corpus/diferencias revisadas. Sin actualización silenciosa ni sustitución del motor. |
| Distribución y calidad | ZIP+checksums, BUILD/avisos, motor separado, Go sin módulos externos. Verify integra motor real; CI pública Windows estándar, smoke local PS5.1/7. Reproducibilidad desde dos carpetas probada en #34. | #39 documenta procedimiento replicable y actualización/reversión. #40 revalida bytes futuros exactos; historial no sustituye CI del candidato. |
| Usuario y mantenimiento | Documentación técnica extensa, pero encabezados activos obsoletos y sin política completa de contrato/mantenimiento. Apache-2.0, NOTICE, Go BSD-3-Clause incluidos; oasdiff externo Apache-2.0. Rutas/datos en informes, sin subida automática. | #39 unifica instalación, límites, privacidad, reporte sintético y mantenimiento sin SLA universal; no inventar nuevas licencias ni servicios. |

### Verificaciones ejecutadas ahora

Windows x64; Go **1.27.4**, oasdiff **1.39.0**, PowerShell **7.6.9** y Windows PowerShell **5.1.26100.10068**. Se reutilizaron herramientas locales fijadas; sin instalar/actualizar dependencias ni descargar toolchains/módulos. Comandos desde raíz, sobre la base `aa194b2` (la rama de planificación solo modifica este documento):

```powershell
git remote -v
git status --short
git config user.name
git config user.email
git var GIT_AUTHOR_IDENT
git var GIT_COMMITTER_IDENT
git fetch origin --tags
git rev-parse HEAD origin/main origin/develop 'v0.1.0^{}' 'v0.2.0^{}'
gh release list --repo v1kz25/openapi-contract-lab
gh issue list --repo v1kz25/openapi-contract-lab --state all --limit 100
gh pr list --repo v1kz25/openapi-contract-lab --state all --limit 70
gh release view v0.2.0 --repo v1kz25/openapi-contract-lab
gh run list --repo v1kz25/openapi-contract-lab --limit 5

$env:Path = (Join-Path (Get-Location) 'dist/_tools/go1.27.4/go/bin') + [IO.Path]::PathSeparator + $env:Path
$env:GOCACHE = Join-Path (Get-Location) 'dist/go-cache'
$env:GOMODCACHE = Join-Path (Get-Location) 'dist/go-mod-cache'
$env:GOENV = 'off'
$env:OASDIFF_BINARY = (Resolve-Path 'dist/oasdiff-1.39.0/oasdiff.exe').Path
./scripts/verify.ps1
./scripts/smoke.ps1
powershell.exe -NoProfile -ExecutionPolicy Bypass -File ./scripts/smoke.ps1

Get-FileHash dist/candidate-v0.2.0-issue-33/openapi-contract-lab_0.2.0_windows_amd64.zip,dist/candidate-v0.2.0-issue-33/SHA256SUMS.txt -Algorithm SHA256
$auditExtract = Join-Path (Get-Location) ('dist/v1-audit-extracted-' + [guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath dist/candidate-v0.2.0-issue-33/openapi-contract-lab_0.2.0_windows_amd64.zip -DestinationPath $auditExtract
Get-Content (Join-Path $auditExtract 'BUILD.txt')
Get-FileHash (Join-Path $auditExtract 'openapi-contract-lab.exe') -Algorithm SHA256
./scripts/smoke.ps1 -Binary (Join-Path $auditExtract 'openapi-contract-lab.exe')
powershell.exe -NoProfile -ExecutionPolicy Bypass -File ./scripts/smoke.ps1 -Binary (Join-Path $auditExtract 'openapi-contract-lab.exe')
```

| Comprobación actual | Resultado real |
| --- | --- |
| `scripts/verify.ps1` | PASS, exit 0: gofmt sin pendientes, `go test -count=1 ./...` con motor real sin omisión de integración, vet y build. CLI 6.553 s, explain 1.712 s, oasdiff 0.880 s, report 0.825 s; main sin tests propios. |
| Smoke binario compilado, PS7 y PS5.1 | PASS, exit 0 ambos: ayuda/versión 0, breaking/exportación 3 con informe completo, compatible/documental/none 0, formato JSON y acentos, rechazo de sobrescritura 1, argumentos inválidos 2, entrada ausente 1 y originales intactos. |
| ZIP v0.2.0 existente | SHA-256 recalculado `6c2bfd76f4c123399029e908bb58f73a40b073fe5f19b7eb4069e1af9aff4c81`, coincide con digest de GitHub y SHA256SUMS. Hash del archivo SHA256SUMS `b4d6266594e4acffee952fcea931ab5f2528f91b4e27a997b18517d67c1aa4b4`, coincide con digest remoto. Reutilizado, no descargado otra vez. |
| Extracción nueva del ZIP y smoke de ese ejecutable | Exe SHA-256 `a5eb2523bce835e8a7ab33796c0baa78bed0e9406796b894fff0c3658f1a15f7`, coincide con BUILD y #34; smoke PS7 y PS5.1 PASS actuales. BUILD/notas incluidos conservan texto de preparación «NOT PUBLISHED»; carencia documental de #39, no indicio de que falte la release remota. |
| Revisión de código/tests y diff inicial | Sin defecto funcional demostrado; lagunas de prueba de pipes/cancelación activa y recuperación de salida parcial identificadas para #38. No se cambiaron expectativas. |

El primer fetch/consultas gh falló por restricciones de permisos/red del entorno; el reintento autorizado funcionó. La primera ejecución de verify terminó PASS pero emitió diagnóstico al intentar escribir caché de módulos global. Se fijó también GOMODCACHE local y se repitió verify sin ese diagnóstico; no fue un fallo de tests ni motivó cambios de código. No queda bloqueo de entorno.

**Evidencia histórica reutilizada:** [CI 45003351377](https://github.com/v1kz25/openapi-contract-lab/actions/runs/45003351377) success en `72ebbd8`, consultada ahora pero ejecutada previamente; [CI 44962771145](https://github.com/v1kz25/openapi-contract-lab/actions/runs/44962771145) success en candidato `f6e7f68ac31d81b8e23cff0e18b41b0e19197c35`. #34 registra dos builds idénticos y nueve pares empaquetados × tres formatos en PS5.1/7. No se repitieron esos builds ni el runner de investigación #25 (no cambia ni es contrato público); no se afirma otra inspección visual, instalación limpia de SO, otro equipo, otra arquitectura ni auditoría completa de red/dependencias transitivas del motor. CI de la PR de planificación se registra en ella; no sustituye la del futuro candidato.

## Objetivo y fases hacia v1.0.0

**Estable significa CLI local fiable, documentada y distribuible dentro de un alcance declarado.** Propuesta: Windows x64/PowerShell 5.1 y 7, Go 1.27.4 para construir y oasdiff 1.39.0 externo para comparar, corpus actual 3.0.3 y casos enumerados 3.1.0, trece explicaciones más fallback, tres formatos y gate vigente. Sin garantía universal sobre contratos ni comportamiento real de servidores. El contrato CLI/JSON propio será estable en 1.x con política explícita de compatibilidad; la semántica upstream permanece ligada al motor comprobado.

Propuesta mínima de versionado JSON: **contrato normativo v1 documentado y ligado a CLI 1.x**, sin añadir campos durante esta planificación ni exigir `schemaVersion` para v1.0.0. #37 debe ratificarlo o justificar un identificador embebido por necesidad de consumidores, documentando migración desde v0.2.0 sin ese campo. Si se mantienen los campos, la migración puede ser solo de garantías/documentación. Campos opcionales aditivos requieren tolerancia de consumidores; rupturas de tipos, significados, campos requeridos, defaults o códigos necesitan major y guía de migración. No congelar prosa de explicaciones ni datos opacos del motor como API propia. Un JSON Schema formal adicional es recomendable si aporta valor, no un requisito automático.

| Fase / issue | Entregable y criterio de salida | Dependencias | Prioridad / dificultad / esfuerzo |
| --- | --- | --- | --- |
| Auditoría y acuerdo: [#36](https://github.com/v1kz25/openapi-contract-lab/issues/36) | Este diagnóstico, backlog y checklist revisables; solo planificación, sin iniciar #37. | v0.2.0 sincronizada; reutiliza #21–#25/#33. | P0; baja/media; 0,5–1 jornada (esta auditoría). |
| Auditoría técnica y estabilización: [#37, contrato CLI/JSON](https://github.com/v1kz25/openapi-contract-lab/issues/37) | Referencia completa, versión normativa/migración y política upstream decididas, tests contractuales que falten. | Plan #36 revisado. | P0; media; 1–2 jornadas; **bloqueante**. |
| Compatibilidad y pruebas: [#38, fronteras/fallos](https://github.com/v1kz25/openapi-contract-lab/issues/38) | Proceso iniciado cancelado, límites por pipes reales, recuperación de salida parcial caracterizada; corpus/seguridad actual intactos. Sin prometer cuota RAM ni sandbox. | #37; reutiliza #7/#23/#24, sin ampliar dialectos por defecto. | P1; media; 1–3 jornadas; **bloqueante** para las garantías declaradas. |
| Distribución y documentación: [#39](https://github.com/v1kz25/openapi-contract-lab/issues/39) | Instalación/update/rollback y empaquetado replicables desde docs versionadas, estado vigente, mantenimiento, privacidad y avisos. | #37; ejemplos finales tras #38. | P1; baja/media; 1–2 jornadas; **bloqueante**. |
| Candidato de release: [#40](https://github.com/v1kz25/openapi-contract-lab/issues/40) | `release/1.0.0`, PR hacia main, checklist y bytes exactos probados, sin publicar ni fusionar. | #37 + #38 + #39 completados y revisados. | P1; media; 1–2 jornadas; **bloqueante**. |
| Publicación estable pendiente de autorización: [#41](https://github.com/v1kz25/openapi-contract-lab/issues/41) | Solo tras aprobación explícita de PR/SHA/archivos: integrar, tag/release estable, descarga/hash/smoke y PR de sincronización propia. | #40 aprobado; fusión de sincronización exige aprobación propia. | P1; baja; 0,5–1 jornada + espera; **bloqueante de publicación**. |

Seis issues nuevos en total, contando planificación. Cada uno contiene evidencia, alcance/exclusiones, dependencias, aceptación, pruebas y estimación. No se duplica #24 para exigir «soporte 3.1» ni #25 para añadir clasificación. Las dependencias concretas son #36 → #37 → #38; #39 puede prepararse tras #37, pero #40 espera #38/#39; #41 espera candidato y aprobación. Una tarea por rama/PR; no avanzar automáticamente al siguiente issue.

### Checklist verificable de estabilidad

Estas casillas son aceptación de **v1.0.0 futura**; que una capacidad pase hoy no acredita los bytes de esa versión.

- [ ] **Alcance**: Windows x64/PS5.1/7 y corpus exacto enumerados; 3.1 parcial/3.2 no verificado, refs externas/URL/UNC/stdin y demás exclusiones explícitas. Sin afirmar igualdad por `[]`.
- [ ] **CLI**: comandos/aliases, orden de argumentos, defaults, formatos, streams, codificación y códigos 0/1/2/3 normativos, probados; exit 3 entrega informe y error operativo tiene prioridad.
- [ ] **JSON/compatibilidad**: tabla completa de campos/tipos/omisión/arrays/enums/upstream/orden, decisión de versión de esquema y migración desde experimental, política major/aditiva y actualización del motor revisadas (#37).
- [ ] **Regresión**: verify sin omisiones con motor fijado; trece pares 3.0.3, cinco pares 3.1.0, desconocidos/fallback, gate y Unicode conservados; CI gratuita del SHA candidato exacto correcta.
- [ ] **Seguridad/recursos**: entradas y configuración aisladas; timeout/pipes reales y comportamiento de fallos probados; limitaciones de tamaño/RAM/sandbox explícitas, sin promesas no verificadas (#38).
- [ ] **Informes/archivos**: escape HTML/Markdown, CSP y enlaces restringidos, orden/conservación de hallazgos y determinismo sin regresión; existentes/inputs intactos; recuperación de destino nuevo parcial documentada. Revisión visual reutilizable o nueva según cambios.
- [ ] **Instalación/artefactos**: descarga/verificación/extracción/uso y actualización/reversión en carpeta nueva documentados y recorridos; smoke PS5.1/7 del ejecutable exacto del ZIP, sin SDK para ejecutar; motor separado.
- [ ] **Versiones/avisos/checksums**: SHA fuente, versiones/flags/build reproducible, BUILD e inventario de archivos; LICENSE/NOTICE/avisos Go conservados; SHA256SUMS de bytes exactos; sin motor/SDK/cachés/logs privados empaquetados.
- [ ] **Notas/mantenimiento/privacidad**: guías activas coherentes, notas de release con límites/migración, reporte sintético de errores y proceso de parches/actualización upstream; revisar rutas/datos antes de compartir, sin subirlos automáticamente ni prometer SLA universal (#39/#40).
- [ ] **Autorización**: propietario aprueba explícitamente PR/SHA/archivos antes de merge/tag/release. Nunca auto-merge ni borrado de ramas; sincronización posterior por otra PR con aprobación propia (#41).

No se exige ausencia absoluta de errores. Son bloqueantes los fallos reproducibles de garantías declaradas (pérdida/clasificación alterada, informe activo inseguro, sobrescritura, streams/códigos incorrectos, incompatibilidad propia sin migración), la falta de evidencia del candidato y la falta de aprobación. Un defecto fuera del alcance documentado se registra y evalúa; no se amplía el soporte para poder cerrar la checklist.

### Esfuerzo y primera tarea

Restan **4,5–10 jornadas efectivas** para #37–#41 (6 h de trabajo concentrado por jornada, aproximadamente **27–60 h**), además de revisión/espera; presupuesto prudente con margen por fallos de Windows/subprocesos: **6–13 jornadas**. No es una fecha comprometida. Supuestos: un desarrollador, herramientas existentes, sin nuevas funciones/plataformas, corpus semántico vigente suficiente y revisión disponible. Incertidumbre principal: pruebas deterministas de fallos de escritura/proceso y posibles defectos descubiertos; un cambio material de contrato obliga a revisar estimación y alcance. Esta auditoría #36 supone 0,5–1 jornada adicional al cálculo restante.

**Primera tarea recomendada: #37**, después de revisar este plan. Fijar qué se promete permite que #38 pruebe garantías concretas y #39 documente la misma interfaz; evita congelar por accidente prosa/bytes upstream o ampliar OpenAPI. No se inicia aquí.

**Recomendables, no bloqueantes por sí mismas:** esquema JSON validable por máquina, helper de empaquetado automatizado, revisión en otro equipo Windows y corpus ampliado si aparece un caso real. Se pueden mantener como notas sin crear más issues. **Posteriores:** plataformas/arquitecturas nuevas, SARIF, GUI, Overlays/autofix, clasificación del tipo de cambio, más plantillas y cobertura 3.1 general/3.2. Ninguna justifica retrasar una v1 acotada sin una necesidad esencial demostrada.

## Versiones futuras, fuera de esta iteración

- **SARIF / integración CI/CD avanzada**: el gate y los códigos actuales ya funcionan en CI. Primero demostrar una mejora real de revisión frente a JSON/Markdown, definir localizaciones, severidades, privacidad de rutas y compatibilidad del consumidor. No añadir exportador ni subida automática sin ese caso de uso.
- **Distribución multiplataforma**: #24 ya documentó supuestos Windows en rutas/enlaces/scripts; Linux/macOS/ARM64, matrices CI y binarios requieren verificación en cada plataforma antes de anunciar soporte. No basta con que Go compile para otra arquitectura.
- **Más reglas o dialectos**: solo por casos reales y pares regla/nivel comprobables. Sin actualizar motor/toolchain, crear diff propio, añadir LLM, cloud ni dependencias de pago como parte de esta planificación. La futura política de mantenimiento sí deberá permitir actualizaciones explícitas, revisadas y probadas por issue/PR.

## Fases de referencia

1. **Validación técnica / PoC**: completada, #1.
2. **Configuración inicial**: completada, #2.
3. **MVP incremental**: completado, #3–#8; v0.1.0 publicada y v0.2.0 experimental publicada.
4. **Pruebas y estabilización**: siguiente objetivo, #37–#39 tras revisión de #36.
5. **Primera versión estable**: candidato #40 y publicación #41 solo con aprobación explícita.
6. **Ampliaciones**: según necesidades demostradas, separadas de la estabilidad.

## Primeras tareas

Las tareas iniciales ya están completadas; no son un backlog pendiente. Revisar primero la PR documental de #36. Después, #37 es la primera implementación recomendada; las demás siguen las dependencias anteriores. Esta ejecución termina al entregar planificación/PR y no comienza ninguna de ellas.

## Principio de alcance

No desarrollar simultáneamente los cuatro repositorios; comenzar por openapi-contract-lab y abrir los demás sólo cuando se decida.

## Etiquetas de fase

1. Validación: `phase:validation`.
2. Preparación: `phase:setup`.
3. MVP: `phase:mvp`.
4. Estabilización: `phase:quality`.
5. Versión pública: `phase:release`.
6. Mejoras posteriores: `phase:future`.
