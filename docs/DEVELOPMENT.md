# Desarrollo — openapi-contract-lab

## Estado

No se ha creado código ni hay comandos específicos de compilación verificados. No inventar instrucciones de ejecución hasta completar el spike/setup.

## Entorno propuesto

- Windows 11; Git, GitHub CLI (`gh`), Codex en desktop/VS Code/CLI.
- Stack en evaluación: Go + oasdiff (dependencia Apache-2.0; confirmar versión y API antes de integrarla).
- Sin hosting, cuentas de servicios comerciales, API keys ni infraestructura de pago.

## Flujo de trabajo

1. Leer AGENTS.md, roadmap y el issue activo.
2. Crear `feature/<issue>-<nombre>` desde `develop` para una única tarea (Gitflow).
3. Implementar sólo su alcance y tests; revisar dependencias/licencias.
4. Ejecutar y registrar pruebas; comprobar Git diff y la identidad Git.
5. Proponer Pull Request; no fusionarla sin permiso.
6. Actualizar el README únicamente para funciones realmente disponibles.

## Tests y comandos

Pendiente del issue de inicialización. No se supone instalado ningún SDK ni versión concreta.

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
