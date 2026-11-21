# openapi-contract-lab

> Estado: comparación local con oasdiff, JSON normalizado y explicaciones deterministas para diez reglas. Informes y política CI aún pendientes; no es un MVP completo.

Explica cambios incompatibles entre contratos OpenAPI usando oasdiff, con informes locales y CI.

## Objetivo

Comparar dos especificaciones OpenAPI sin reimplementar el análisis semántico; enriquecer hallazgos con explicación determinista y generar informes offline.

## Stack

Go 1.27.4, sin módulos externos. La comparación ejecuta el binario local oasdiff 1.39.0 (Apache-2.0), cuya versión se comprueba antes de comparar.

## Uso disponible

Desde la raíz del repositorio, con Go 1.27.4:

```powershell
go run ./cmd/openapi-contract-lab --help
go run ./cmd/openapi-contract-lab --version
go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
./dist/openapi-contract-lab.exe --help
./dist/openapi-contract-lab.exe compare --oasdiff ./dist/oasdiff-1.39.0/oasdiff.exe ./spikes/oasdiff/fixtures/base.yaml ./spikes/oasdiff/fixtures/required.yaml
```

Sin argumentos muestra ayuda. `--version` devuelve `openapi-contract-lab dev`: aún no hay release. Un comando desconocido devuelve un mensaje en stderr y código 2. La CLI no inicia un servidor ni una interfaz gráfica.

Para comparar, preparar el binario oficial según el [spike](docs/OASDIFF-SPIKE.md). `--oasdiff` puede omitirse si oasdiff está en PATH. Las dos entradas deben ser archivos locales; se deshabilitan las referencias externas. La salida conserva rule ID, severidad, ruta, ubicaciones y el evento original. Ver [modelo y pruebas del adaptador](docs/OASDIFF-ADAPTER.md).

Cada hallazgo incluye `explanation` con resumen, impacto, acción sugerida y ubicaciones. Las [diez reglas verificadas](docs/EXPLANATIONS.md) tienen texto determinista en español; otras reglas, niveles o versiones reciben un fallback neutral. La explicación conserva la clasificación de oasdiff y añade enlaces locales solo para fuentes que coinciden con las entradas conocidas.

La salida JSON utiliza escapes Unicode para que Windows PowerShell 5.1 conserve los acentos incluso con codificaciones de consola antiguas. `ConvertFrom-Json` recupera el texto original; no hace falta cambiar la configuración global de Windows.

Una comparación ejecutada correctamente devuelve **0 incluso si hay breaking changes**; aún no hay política CI por severidad. Una lista vacía de hallazgos no implica que los contratos sean idénticos: los cambios solo de descripción no aparecen en el changelog.

## MVP previsto

- comparar dos OpenAPI 3.x con oasdiff como fuente de verdad
- clasificar hallazgos de ruptura y no ruptura
- explicar cambios con reglas deterministas, localización y severidad
- informes Markdown/HTML y código de salida para CI

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

## Riesgos

Mapear los hallazgos de oasdiff a un modelo estable, diferencias entre OpenAPI 3.0/3.1/3.2, y evitar explicaciones inexactas.

## Referencias

- https://github.com/oasdiff/oasdiff
- https://spec.openapis.org/overlay/

## Licencia

**Pendiente de decisión.** Este repositorio se publica sin archivo LICENSE por ahora. Ser público en GitHub no otorga automáticamente permiso para copiar, modificar o redistribuir su contenido.
