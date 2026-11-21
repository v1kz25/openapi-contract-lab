# OpenAPI Contract Lab

## Revisión del contrato

**Motor:** oasdiff 1.39.0

**Base:** base.yaml

**Head:** head.yaml

## Resumen

- Hallazgos: 9
- Errores: 3
- Advertencias: 2
- Informativos: 2
- Sin clasificar: 2

El changelog incluye cambios que afectan a consumidores. Una lista vacía no implica contratos idénticos; los cambios solo de descripción no generan hallazgos.

## Hallazgos

### GET /items

Hallazgos: 2 · Errores: 1 · Advertencias: 1 · Informativos: 0 · Sin clasificar: 0

#### Severidad: error · Hallazgos: 1

##### 1. El parámetro pasa de opcional a obligatorio.

**Regla:** request-parameter-became-required

**Severidad:** error · **Compatibilidad:** breaking · **Nivel upstream:** 3

**Operación:** GET /items · **Sección:** paths

**Mensaje de oasdiff:** the \`query\` request parameter \`filter\` became required

**Explicación:** template

**Impacto:** Una petición que omitía el parámetro puede dejar de cumplir el contrato revisado.

**Acción sugerida:** Revise las llamadas que lo omiten y coordine su adaptación antes de exigirlo.

**Ubicaciones:**

- **base:** [base.yaml](<file:///synthetic/base.yaml#L12>) · línea 12, columna 11

- **head:** [head.yaml](<file:///synthetic/head.yaml#L12>) · línea 12, columna 11

#### Severidad: warning · Hallazgos: 1

##### 2. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** unverified-warning

**Severidad:** warning · **Compatibilidad:** potentially-breaking · **Nivel upstream:** 2

**Operación:** GET /items

**Mensaje de oasdiff:** synthetic warning

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.

### POST /items

Hallazgos: 2 · Errores: 1 · Advertencias: 0 · Informativos: 1 · Sin clasificar: 0

#### Severidad: error · Hallazgos: 1

##### 6. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** post-error

**Severidad:** error · **Compatibilidad:**  · **Nivel upstream:** 3

**Operación:** POST /items

**Mensaje de oasdiff:** Synthetic POST error

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.

#### Severidad: info · Hallazgos: 1

##### 3. El parámetro deja de ser obligatorio.

**Regla:** request-parameter-became-optional

**Severidad:** info · **Compatibilidad:** non-breaking · **Nivel upstream:** 1

**Operación:** POST /items

**Mensaje de oasdiff:** synthetic optional parameter

**Explicación:** template

**Impacto:** Omitirlo pasa a estar permitido por esta condición; enviarlo sigue siendo válido respecto a su obligatoriedad.

**Acción sugerida:** Compruebe el comportamiento del servidor cuando el parámetro no se envía.

**Ubicaciones:**

- Sin ubicación reportada.

### GET /orders

Hallazgos: 2 · Errores: 1 · Advertencias: 0 · Informativos: 1 · Sin clasificar: 0

#### Severidad: error · Hallazgos: 1

##### 7. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** orders-error

**Severidad:** error · **Compatibilidad:**  · **Nivel upstream:** 3

**Operación:** GET /orders

**Mensaje de oasdiff:** Synthetic orders error

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.

#### Severidad: info · Hallazgos: 1

##### 5. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** orders-info

**Severidad:** info · **Compatibilidad:**  · **Nivel upstream:** 1

**Operación:** GET /orders

**Mensaje de oasdiff:** Synthetic orders info

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.

### Sin ruta u operación

Hallazgos: 3 · Errores: 0 · Advertencias: 1 · Informativos: 0 · Sin clasificar: 2

#### Severidad: warning · Hallazgos: 1

##### 8. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** path-only

**Severidad:** warning · **Compatibilidad:**  · **Nivel upstream:** 2

**Operación:** /items

**Mensaje de oasdiff:** Synthetic path event

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.

#### Severidad: future · Hallazgos: 1

##### 9. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** method-only

**Severidad:** future · **Compatibilidad:**  · **Nivel upstream:** 42

**Operación:** PATCH

**Mensaje de oasdiff:** Synthetic method event

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.

#### Severidad: unknown · Hallazgos: 1

##### 4. No hay una explicación verificada para esta combinación de regla y nivel.

**Regla:** future-rule

**Severidad:** unknown · **Compatibilidad:** unknown · **Nivel upstream:** 99

**Operación:** No reportada

**Mensaje de oasdiff:** synthetic unknown event

**Explicación:** fallback

**Impacto:** Consulte el mensaje y los datos originales de oasdiff; no se ha inferido una causa adicional.

**Acción sugerida:** Revise manualmente el cambio y las ubicaciones disponibles.

**Ubicaciones:**

- Sin ubicación reportada.
