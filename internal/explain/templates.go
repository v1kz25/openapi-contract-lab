package explain

type template struct {
	level                   int
	summary, impact, action string
}

// These rule/level pairs are verified against oasdiff 1.39.0 with synthetic specs.
// Messages describe only this change, not the compatibility of the whole API.
var templates = map[string]template{
	"request-parameter-became-required": {
		3, "El parámetro pasa de opcional a obligatorio.",
		"Una petición que omitía el parámetro puede dejar de cumplir el contrato revisado.",
		"Revise las llamadas que lo omiten y coordine su adaptación antes de exigirlo.",
	},
	"request-parameter-became-optional": {
		1, "El parámetro deja de ser obligatorio.",
		"Omitirlo pasa a estar permitido por esta condición; enviarlo sigue siendo válido respecto a su obligatoriedad.",
		"Compruebe el comportamiento del servidor cuando el parámetro no se envía.",
	},
	"request-parameter-became-enum": {
		3, "Se añade una enumeración al esquema del parámetro.",
		"Los valores fuera de la enumeración dejan de estar permitidos por esta restricción.",
		"Contraste los valores utilizados por los clientes con la enumeración revisada.",
	},
	"request-parameter-became-not-enum": {
		1, "Se elimina la enumeración del esquema del parámetro.",
		"Se amplía el conjunto admitido por esta restricción; las demás restricciones del esquema siguen aplicándose.",
		"Verifique que la validación y la documentación reflejen los valores ahora admitidos.",
	},
	"request-parameter-enum-value-added": {
		1, "Se añade un valor permitido a la enumeración del parámetro.",
		"El valor añadido pasa a estar admitido por esta enumeración.",
		"Compruebe que el servidor acepte el nuevo valor antes de que los clientes lo utilicen.",
	},
	"request-parameter-enum-value-removed": {
		3, "Se retira un valor de la enumeración del parámetro.",
		"Una petición que envíe ese valor puede dejar de cumplir la enumeración revisada.",
		"Localice los clientes que envían el valor retirado y acuerde una alternativa compatible.",
	},
	"request-parameter-became-nullable": {
		1, "El esquema del parámetro pasa a admitir null.",
		"Se amplía la aceptación de valores nulos por este esquema; null no equivale a omitir el parámetro.",
		"Revise cómo se serializan y validan los valores nulos del parámetro.",
	},
	"request-parameter-became-not-nullable": {
		3, "El esquema del parámetro deja de admitir null.",
		"Un valor nulo antes permitido puede incumplir el esquema revisado; esto no cambia por sí solo la obligatoriedad del parámetro.",
		"Revise los clientes que envían null y coordine un valor válido antes de restringirlo.",
	},
	"response-property-enum-value-added": {
		3, "Se amplía la enumeración de una propiedad de respuesta.",
		"El servidor puede devolver un valor que el contrato anterior no contemplaba; un cliente que solo maneje los valores anteriores puede fallar.",
		"Revise el manejo de valores desconocidos y coordine el cambio con los consumidores.",
	},
	"response-property-enum-value-removed": {
		1, "Se reduce la enumeración de una propiedad de respuesta.",
		"La retirada de ese valor reduce las posibilidades de respuesta contempladas por este hallazgo.",
		"Compruebe que el servidor deje de emitir el valor retirado y revise los flujos que dependían de recibirlo.",
	},
}
