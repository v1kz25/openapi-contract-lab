# openapi-contract-lab

> Estado: base Go con ayuda y versión. **Aún no hay comparación de contratos ni MVP implementado**.

Explica cambios incompatibles entre contratos OpenAPI usando oasdiff, con informes locales y CI.

## Objetivo

Comparar dos especificaciones OpenAPI sin reimplementar el análisis semántico; enriquecer hallazgos con explicación determinista y generar informes offline.

## Stack

Go 1.27.4, sin dependencias externas en el módulo inicial. El spike recomienda integrar el binario local oasdiff 1.39.0 (Apache-2.0); la CLI todavía no lo invoca.

## Uso disponible

Desde la raíz del repositorio, con Go 1.27.4:

```powershell
go run ./cmd/openapi-contract-lab --help
go run ./cmd/openapi-contract-lab --version
go build -o dist/openapi-contract-lab.exe ./cmd/openapi-contract-lab
./dist/openapi-contract-lab.exe --help
```

Sin argumentos muestra ayuda. `--version` devuelve `openapi-contract-lab dev`: aún no hay release. Un comando desconocido devuelve un mensaje en stderr y código 2. La CLI no inicia un servidor ni una interfaz gráfica.

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
