# AGENTS.md — Reglas para Codex

## Autoridad y alcance
- Lee README.md, docs/ROADMAP.md, docs/DEVELOPMENT.md y el GitHub Issue activo antes de editar.
- Trabaja solamente en el issue indicado. Si falta una decisión, documenta el bloqueo y consulta; no amplíes alcance.
- No implementes otras fases, no generes carpetas vacías, no añadas complejidad sin razón.
- El proyecto debe tener coste monetario adicional obligatorio 0 €: offline/local por defecto, sin servicios ni APIs de pago, sin credenciales ni tarjetas, sin runners facturables.
- Antes de incorporar dependencias comprueba licencia, compatibilidad Windows 11 y necesidad real. No copies código ajeno sin cumplir licencias.

## Calidad
- Código sencillo, modular, bien tipado, mantenible; SOLID cuando resulte útil.
- Valida entradas, trata fallos con mensajes útiles, cubre lógica crítica con tests y usa fixtures sintéticos.
- Documenta decisiones significativas. Ejecuta pruebas reales; nunca afirmes que pasaron sin haberlas ejecutado.
- No publiques datos personales, credenciales, archivos privados o trazas sensibles.

## Git y autoría
- Respeta `git config user.name` / `git config user.email` existentes; no los modifiques global ni localmente sin autorización.
- No añadas trailers `Co-authored-by:` atribuidos a Codex, ChatGPT, OpenAI ni otros asistentes de IA.
- No suplantes author/committer. Si una integración hace commits como bot, conserva su identidad real y advierte las atribuciones obligatorias.
- Commits descriptivos; nunca `force push`, nunca borres ramas ni sobrescribas cambios ajenos.
- Sigue Gitflow según docs/DEVELOPMENT.md: tareas y spikes en `feature/<issue>-<nombre>` desde `develop`, con PR hacia `develop`.
- No mezcles cambios fuera del issue. Propón PR para revisión; nunca fusiones una PR sin aprobación explícita.
- Antes de push/commit verifica diff, staged files, author/committer y que el mensaje no contiene trailers de IA.

## Entrega por issue
- Resume archivos modificados, decisiones, tests ejecutados con resultados, limitaciones y próximos pasos; incluye referencia al issue y PR si existen.

## Particularidades de openapi-contract-lab
- Mapear los hallazgos de oasdiff a un modelo estable, diferencias entre OpenAPI 3.0/3.1/3.2, y evitar explicaciones inexactas.
- El motor de oasdiff ofrece compatibilidad semántica; la aportación del proyecto es la revisión y explicación accionable.

- No cambies la configuración global de Git. Una tarea pequeña por rama y Pull Request; comprueba pruebas y diff antes de proponer integración.
- Cuando se necesite CI, usa exclusivamente recursos gratuitos de GitHub Actions; no habilites servicios que requieran tarjeta.

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
