# Preparación de la primera publicación — issue #8

Este issue prepara instrucciones y comprobaciones. **No publica una versión estable, no crea tags, no cambia la licencia ni distribuye binarios.** La CLI conserva `openapi-contract-lab dev`. [Instalación y ejemplos](INSTALL.md) describen funciones disponibles, sin anunciar un instalador ni una GUI.

## Evidencia real — 2026-11-21

Base de código: `fbc97e61bfc1e7ceec46772a2a60a2a9f5c0dbc4` de `develop`, con los scripts/documentación de preparación del issue #8; el SHA completo de la revisión se registra en su PR. No hay cambios en el código Go de producción.

- Host Windows x64, build **26300.9974** (API Windows NT 10.0.26300); objetivo Windows 11. No se infiere otra edición a partir de la API de compatibilidad.
- Go **1.27.4**, oasdiff **1.39.0**; hashes de los archivos oficiales comprobados y coincidentes con [INSTALL.md](INSTALL.md). Herramientas portátiles existentes; no se simuló una instalación limpia del sistema operativo ni se repitió la descarga del SDK en este issue.
- `scripts/verify.ps1`: PASS con `OASDIFF_BINARY`, tests reales sin caché, formato, `go vet` y compilación. Sin descargas de módulos ni toolchains.
- `scripts/smoke.ps1` sobre el ejecutable compilado, en PowerShell **5.1.26100.9818** y **7.6.9**: ayuda/versión/compare-help 0; breaking 3; compatible/documental/none 0; exportaciones JSON/Markdown/HTML 3 con archivo completo UTF-8 y stdout vacío; reintentos de escritura 1 y SHA-256 sin cambios; política inválida 2; entrada ausente 1. Acentos JSON y UTF-8 comprobados, incluidas rutas de salida con espacios. Logs locales en `dist/smoke-*`, sin publicarlos.
- Dos compilaciones con código copiado en directorios distintos, Go/GOOS/GOARCH/GOAMD64/CGO/flags fijados según INSTALL: **bytes idénticos**, SHA-256 `a18cde9bf7a8618bea985ea6a0a2436ecca60ddb4cff45833222331e831dc72c`. Es evidencia local de esta revisión, no un checksum de release ni una garantía para cualquier entorno. `go version -m` confirmó Go, trimpath, Windows/amd64, GOAMD64=v1 y CGO deshabilitado.
- CI añade el smoke tras compilar. Su ejecución real y revisión del diff se registran en la PR; no se confunde la prueba local con el runner alojado.

Los smoke automatizan comandos del ejecutable real: no hay que pedir al usuario que los repita. La revisión visual del HTML de breaking/vacío ya fue aportada por el usuario durante el issue #5 (PR #14); este issue no cambia su presentación ni afirma nuevas comprobaciones visuales. No se verificó ARM64, Linux/macOS ni una instalación en otra máquina.

## Límites que deben acompañar una publicación

- oasdiff 1.39.0 es el único motor semántico y debe obtenerse por separado. Referencias externas, URLs/UNC y servicios remotos deshabilitados. Timeout de comparación: 30 segundos.
- Fixtures mínimos 3.0.3 y 3.1.0; cobertura limitada, sin validación completa de 3.1 ni soporte verificado de 3.2. Diez explicaciones específicas y fallback neutral para otras combinaciones.
- Un informe sin hallazgos no implica contratos idénticos. Códigos de salida propios 0/1/2/3; JSON/CLI en desarrollo. Rutas locales forman parte de los informes y deben revisarse antes de compartirlos.
- Coste obligatorio 0 €; ejecución local y CI estándar gratuita en repositorio público. Sin APIs de IA, cuentas comerciales ni secretos nuevos. Go sin módulos externos; oasdiff Apache-2.0 y Go BSD-3-Clause ya documentados.

## Revisión antes de tag/release

1. Aprobar la PR de preparación hacia `develop`. No equivale a aprobar la publicación.
2. Resolver la **licencia propia pendiente** con decisión explícita del propietario. Si se redistribuyera oasdiff u otro componente, revisar además inventario, licencias y avisos de los binarios incluidos; esta preparación no los empaqueta.
3. Acordar versión, alcance estable y notas de release después de revisar las limitaciones. No se reserva ni crea un tag en este issue. Si se solicita que `--version` identifique una versión publicada, preparar y verificar ese cambio antes de publicar; hoy devuelve `dev`.
4. Crear `release/<version>` desde el `develop` aprobado para estabilizar según Gitflow; ejecutar CI y smoke sobre su SHA exacto. PR hacia `main`, con revisión y autorización explícita para fusionar; después otra PR de sincronización hacia `develop`. No borrar ramas.
5. Solo tras esa revisión y aprobación explícita de publicación, definir el tag sobre el commit aprobado de `main` y la release correspondiente. Si se acuerdan artefactos gratuitos, fijar fuente/toolchain/flags, generar checksums y avisos, y probar exactamente los archivos que se publicarán. No usar runners facturables ni habilitar auto-merge.

La licencia y las decisiones finales de versión/publicación son pendientes de la futura release, no motivos para inventar una licencia, tag o paquete en esta tarea de preparación. El issue termina con documentación, smoke verificados y PR para revisión.
