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
- Usa ramas `feat/<issue>-<nombre>` o `chore/<issue>-<nombre>` cuando corresponda.
- No mezcles cambios fuera del issue. Propón PR para revisión; nunca fusiones una PR sin aprobación explícita.
- Antes de push/commit verifica diff, staged files, author/committer y que el mensaje no contiene trailers de IA.

## Entrega por issue
- Resume archivos modificados, decisiones, tests ejecutados con resultados, limitaciones y próximos pasos; incluye referencia al issue y PR si existen.

## Particularidades de openapi-contract-lab
- Mapear los hallazgos de oasdiff a un modelo estable, diferencias entre OpenAPI 3.0/3.1/3.2, y evitar explicaciones inexactas.
- El motor de oasdiff ofrece compatibilidad semántica; la aportación del proyecto es la revisión y explicación accionable.

- No cambies la configuración global de Git. Una tarea pequeña por rama y Pull Request; comprueba pruebas y diff antes de proponer integración.
- Cuando se necesite CI, usa exclusivamente recursos gratuitos de GitHub Actions; no habilites servicios que requieran tarjeta.
