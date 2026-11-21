# openapi-contract-lab

> Estado: CLI de desarrollo local con comparación oasdiff, diez explicaciones, informes y política de severidad. Preparación y smoke tests de la primera publicación disponibles; todavía sin tag ni release estable.

Explica cambios incompatibles entre contratos OpenAPI usando oasdiff, con informes locales y CI.

## Objetivo

Comparar dos especificaciones OpenAPI sin reimplementar el análisis semántico; enriquecer hallazgos con explicación determinista y generar informes offline.

## Stack

Go 1.27.4, sin módulos externos. La comparación ejecuta el binario local oasdiff 1.39.0 (Apache-2.0), cuya versión se comprueba antes de comparar.

## Uso disponible

Para empezar desde cero en Windows x64, ver [instalación y compilación verificables](docs/INSTALL.md). No hay instalador ni paquete publicado. La guía fija Go 1.27.4, oasdiff 1.39.0 y sus checksums; no requiere instalación global.

Desde la raíz del repositorio, con Go 1.27.4 y el motor local ya preparados:

```powershell
go run ./cmd/openapi-contract-lab --help
go run ./cmd/openapi-contract-lab --version
go build -trimpath -buildvcs=false -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
./dist/openapi-contract-lab.exe --help
./dist/openapi-contract-lab.exe compare --oasdiff ./dist/oasdiff-1.39.0/oasdiff.exe ./spikes/oasdiff/fixtures/base.yaml ./spikes/oasdiff/fixtures/required.yaml
$LASTEXITCODE # 3: informe generado; breaking change
```

Sin argumentos muestra ayuda. `--version` devuelve `openapi-contract-lab dev`: aún no hay release. Un comando desconocido devuelve un mensaje en stderr y código 2. La CLI no inicia un servidor ni una interfaz gráfica.

Para comparar, preparar el binario oficial según el [spike](docs/OASDIFF-SPIKE.md). `--oasdiff` puede omitirse si oasdiff está en PATH. Las dos entradas deben ser archivos locales; se deshabilitan las referencias externas. La salida conserva rule ID, severidad, ruta, ubicaciones y el evento original. Ver [modelo y pruebas del adaptador](docs/OASDIFF-ADAPTER.md).

Cada hallazgo incluye `explanation` con resumen, impacto, acción sugerida y ubicaciones. Las [diez reglas verificadas](docs/EXPLANATIONS.md) tienen texto determinista en español; otras reglas, niveles o versiones reciben un fallback neutral. La explicación conserva la clasificación de oasdiff y añade enlaces locales solo para fuentes que coinciden con las entradas conocidas.

La salida JSON utiliza escapes Unicode para que Windows PowerShell 5.1 conserve los acentos incluso con codificaciones de consola antiguas. `ConvertFrom-Json` recupera el texto original; no hace falta cambiar la configuración global de Windows.

Para exportar un informe local, elegir un archivo nuevo:

```powershell
./dist/openapi-contract-lab.exe compare --oasdiff ./dist/oasdiff-1.39.0/oasdiff.exe --format html --output dist/revision.html ./spikes/oasdiff/fixtures/base.yaml ./spikes/oasdiff/fixtures/required.yaml
```

`--format` admite `json` (predeterminado), `markdown` y `html`. `--output` escribe UTF-8 directamente y evita la redirección de PowerShell; rechaza destinos existentes. El HTML incluye todo lo necesario, sin servidor, scripts ni recursos de red. Ver [exportación y pruebas](docs/REPORTS.md).

Por defecto, los hallazgos `error` (breaking) devuelven **3 después de generar el informe**. `--fail-on error|warning|info|none` permite cambiar el umbral; `none` conserva el comportamiento de exportación sin gate. Severidades desconocidas bloquean con cualquier política activa. Los fallos operativos devuelven 1 y los argumentos inválidos 2. Ver [política y ejemplos PowerShell](docs/SEVERITY-POLICY.md).

Una lista vacía de hallazgos no implica que los contratos sean idénticos: los cambios solo de descripción no aparecen en el changelog.

## Funciones disponibles y límites

- Comparación de dos archivos locales mediante oasdiff 1.39.0. Regresión probada con fixtures 3.0.3 y casos mínimos 3.1.0; no se afirma cobertura completa de 3.1 ni soporte verificado de 3.2.
- Clasificación conservada del motor y diez explicaciones deterministas en español; fallback neutral para otras reglas/niveles.
- Informes JSON/Markdown/HTML offline y gate configurable por severidad. Timeout de comparación de 30 segundos; referencias externas y entradas URL/UNC deshabilitadas.
- Solo Windows x64 verificado; PowerShell 5.1 y 7. El esquema de informe y la CLI siguen en desarrollo. Sin GUI, comparación por carpetas, editor ni autofix.

### Fuera del MVP

- motor propio de diff
- LLM o API de pago
- editor de contratos completo
- autofix/Overlays automáticos generales
- servicio cloud

## Desarrollo

- Ver [roadmap](docs/ROADMAP.md), [guía de desarrollo](docs/DEVELOPMENT.md) y [AGENTS.md](AGENTS.md).
- El punto de entrada son los GitHub Issues, no implementar el MVP de golpe.
- Spike local de oasdiff (issue #1): [decisión, ejemplos y evidencia](docs/OASDIFF-SPIKE.md). No es una CLI del proyecto.
- Desarrollo principalmente desde Windows 11; sin infraestructura de pago.
- [CI y regresión](docs/CI.md): runner estándar Windows, versiones fijadas, pruebas reales con oasdiff y análisis estático. El mismo script se puede ejecutar localmente.
- [Preparación de publicación](docs/RELEASE.md): evidencias, límites y revisión pendiente antes de cualquier tag/release. `scripts/smoke.ps1` comprueba el ejecutable ya compilado sin SDK ni descargas.

## Riesgos

Mapear los hallazgos de oasdiff a un modelo estable, diferencias entre OpenAPI 3.0/3.1/3.2, y evitar explicaciones inexactas.

## Referencias

- https://github.com/oasdiff/oasdiff
- https://spec.openapis.org/overlay/

## Licencia

**Pendiente de decisión.** Este repositorio se publica sin archivo LICENSE por ahora. Ser público en GitHub no otorga automáticamente permiso para copiar, modificar o redistribuir su contenido.
