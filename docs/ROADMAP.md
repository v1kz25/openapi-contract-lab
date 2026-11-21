# Roadmap — openapi-contract-lab

## Estado publicado

[v0.1.0 experimental](https://github.com/v1kz25/openapi-contract-lab/releases/tag/v0.1.0) ya está publicada, bajo Apache-2.0, y sincronizada a `develop` mediante la PR #20. Incluye comparación local con oasdiff 1.39.0, diez explicaciones deterministas, informes JSON/Markdown/HTML y política de severidad. CI y smoke Windows x64 verificados. No se promete estabilidad de CLI/esquema ni cobertura completa de OpenAPI. Go 1.27.4 y motor fijados; sin módulos Go externos ni coste adicional obligatorio.

La publicación cerró el [issue #18](https://github.com/v1kz25/openapi-contract-lab/issues/18). Se conservan sus [notas](releases/v0.1.0.md) como evidencia de esa versión; las futuras publicaciones requieren aprobación independiente según [RELEASE.md](RELEASE.md).

## Próxima iteración: propuesta v0.2.0

Backlog reducido, sin comprometer una fecha ni publicar una nueva versión. Prima explicar el impacto y facilitar la revisión frente a envolver más comandos de oasdiff. Solo se implementa ahora la primera funcionalidad; el mantenimiento documental se entrega en una PR independiente como preparación del backlog.

| Orden / issue | Valor y alcance | Dificultad / dependencias | Riesgo y criterio de corte |
| --- | --- | --- | --- |
| Preparación: [#22, documentación](https://github.com/v1kz25/openapi-contract-lab/issues/22) | Reflejar v0.1.0 publicada y corregir ejemplos activos de códigos de salida. | Baja; publicación y sincronización ya realizadas. | Separar estado vigente de evidencias históricas; sin cambiar código. |
| 1: [#21, tres explicaciones](https://github.com/v1kz25/openapi-contract-lab/issues/21) | Tipo incompatible de parámetro, cuerpo obligatorio y retirada de propiedad obligatoria de respuesta; impacto/acción en español. | Baja; ampliar el mecanismo existente, catálogo y fixtures reales 1.39.0. Independiente de #22. | Máximo tres pares regla/nivel comprobados; mantener fallback y todos los datos upstream. |
| 2: [#23, informes agrupados](https://github.com/v1kz25/openapi-contract-lab/issues/23) | Revisar por endpoint/operación y severidad; sección para hallazgos globales. | Media; renderer existente, no depende del spike #25. Revisar primero la PR de #21. | Orden estable, ningún hallazgo perdido, escape y CSP intactos; no cambiar JSON ni política. |
| 3: [#24, regresión 3.1](https://github.com/v1kz25/openapi-contract-lab/issues/24) | Hasta cinco pares representativos, incluidos null/uniones y respuestas; confianza verificable. | Media; motor fijado y corpus existente, después de revisar #21. | No extrapolar nullable 3.0 ni afirmar soporte total; documentar incertidumbre y mantener fallback. |
| 4: [#25, validar tipo de cambio](https://github.com/v1kz25/openapi-contract-lab/issues/25) | Investigar idéntico/documental/contractual con las capacidades de oasdiff; evitar falsa igualdad al no haber hallazgos. | Media/alta; validar antes de diseñar una ampliación del adaptador/esquema. | Solo spike con cinco pares y decisión; si no hay señal fiable o exige gran cambio público, posponer. |

El mínimo útil propuesto es #21 + #23 con regresión #24 y documentación vigente. #25 es exploración condicionada: no se promete implementar la clasificación en v0.2.0. El usuario revisará alcance y criterios de release después; no se inicia automáticamente el siguiente issue ni se cambia `--version` durante esta preparación.

## Versiones futuras, fuera de esta iteración

- **SARIF / integración CI/CD avanzada**: el gate y los códigos actuales ya funcionan en CI. Primero demostrar una mejora real de revisión frente a JSON/Markdown, definir localizaciones, severidades, privacidad de rutas y compatibilidad del consumidor. No añadir exportador ni subida automática sin ese caso de uso.
- **Distribución multiplataforma**: #24 puede identificar supuestos Windows en rutas/enlaces/scripts, pero Linux/macOS/ARM64, matrices CI y binarios requieren verificación en cada plataforma antes de anunciar soporte. No basta con que Go compile para otra arquitectura.
- **Más reglas o dialectos**: solo por casos reales y pares regla/nivel comprobables. Sin actualizar motor/toolchain, crear diff propio, añadir LLM, cloud ni dependencias de pago como parte del backlog actual.

## Fases de referencia

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
