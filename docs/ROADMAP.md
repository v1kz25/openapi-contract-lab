# Roadmap — openapi-contract-lab

Hay un spike local de oasdiff y una CLI Go con comparación local, explicaciones para diez reglas, informes JSON/Markdown/HTML y política de severidad con códigos de salida. La CI para PR ejecuta pruebas de regresión y análisis estático. La release sigue pendiente. Este plan es orientativo; cada fase requiere validar la anterior.

La preparación de la primera publicación incluye instalación desde fuente y smoke tests del ejecutable en Windows. Tag, versión, licencia propia y publicación se revisan por separado según [RELEASE.md](RELEASE.md).

1. **Validación técnica / PoC**: validar el hueco real y dependencias
2. **Configuración inicial**: crear base mínima verificable
3. **MVP incremental**: funcionalidad por issues pequeños
4. **Pruebas y estabilización**: pruebas, documentación y límites
5. **Primera versión pública**: release revisada y reproducible
6. **Ampliaciones**: funciones adicionales tras feedback

## Primeras tareas

- Validar integración local de oasdiff y formatos de salida
- Inicializar módulo Go, CLI mínima y pruebas base
- Crear adaptador y modelo normalizado de hallazgos oasdiff
- Incorporar explicaciones deterministas para reglas prioritarias
- Exportar informes Markdown y HTML locales
- Configurar política de severidad y códigos de salida CI
- Añadir CI gratuita y conjunto de fixtures de regresión
- Preparar primera versión pública reproducible

## Principio de alcance

No desarrollar simultáneamente los cuatro repositorios; comenzar por openapi-contract-lab y abrir los demás sólo cuando se decida.

## Etiquetas de fase

1. Validación: `phase:validation`.
2. Preparación: `phase:setup`.
3. MVP: `phase:mvp`.
4. Estabilización: `phase:quality`.
5. Versión pública: `phase:release`.
6. Mejoras posteriores: `phase:future`.
