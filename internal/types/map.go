package types

const (
	// Tabelas Principais
	TABLE_PROCESSES    = "processes"
	TABLE_STEPS        = "steps"
	TABLE_RECORDS      = "records"
	TABLE_FIELDS       = "fields"
	TABLE_HEADERS      = "headers"
	TABLE_HEADER_TYPES = "header_types"
	TABLE_RECORD_TYPES = "record_types"

	// Tabelas de Erros
	TABLE_SYNTAX_ERRORS = "syntax_errors"
	TABLE_DATA_ERRORS   = "data_errors"

	// Tabelas Geográficas e Agentes
	TABLE_CAE_REV4           = "cae_rev4"
	TABLE_COUNTRIES          = "countries"
	TABLE_DISTRICTS          = "districts"
	TABLE_MUNICIPALITIES     = "municipalities"
	TABLE_PARISHES           = "parishes"
	TABLE_INE_ZONES          = "ine_zones"
	TABLE_POSTAL_CODES       = "postal_codes"
	TABLE_ELECTRICAL_UNITS   = "electrical_units"
	TABLE_AGENT_TYPES        = "agent_types"
	TABLE_NETWORK_OPERATORS  = "network_operators"
	TABLE_RETAILERS          = "retailers"
	TABLE_LOGISTICS_OPERATOR = "logistics_operator"

	// Tabelas Técnicas
	TABLE_SERVICE_QUALITY_ZONES  = "service_quality_zones"
	TABLE_VOLTAGE_LEVELS         = "voltage_levels"
	TABLE_DELIVERY_POINT_TYPES   = "delivery_point_types"
	TABLE_INSTALLATION_CHARS     = "installation_chars"
	TABLE_MEASUREMENT_VOLTAGE    = "measurement_voltage"
	TABLE_REPORTING_VOLTAGE      = "reporting_voltage"
	TABLE_NUMBER_OF_PHASES       = "number_of_phases"
	TABLE_CONTRACTED_POWER       = "contracted_power"
	TABLE_CONSUMPTION_PROFILE    = "consumption_profile"
	TABLE_DELIVERY_POINT_CONTACT = "delivery_point_contact"

	// Tabelas de Cliente e Identificação
	TABLE_IDENTIFICATION_TYPES   = "identification_types"
	TABLE_CUSTOMER_TYPES         = "customer_types"
	TABLE_HOLDER_ADDRESS_TYPES   = "holder_address_types"
	TABLE_CNE_IDENTIFICATION     = "cne_identification"
	TABLE_CNE_CONTACT            = "cne_contact"
	TABLE_CNE_PREFERRED_CONTACT  = "cne_preferred_contact"
	TABLE_DEFICIENCY_EQUIPMENT   = "deficiency_equipment"
	TABLE_PRIORITY_CUSTOMER_LOC  = "priority_customer_loc"
	TABLE_MR_CONTRACTING_REASONS = "mr_contracting_reasons"
	TABLE_RPE_ACCESS_PURPOSE     = "rpe_access_purpose"

	// Tabelas de Medição e Equipamento
	TABLE_EQUIPMENT_BRANDS           = "equipment_brands"
	TABLE_MEASUREMENT_DEVICE_TYPE    = "measurement_device_type"
	TABLE_OWNERSHIP                  = "ownership"
	TABLE_MEASUREMENT_DEVICE_FUNC    = "measurement_device_func"
	TABLE_SUPPORTED_MEASUREMENT_FUNC = "supported_measurement_func"
	TABLE_MEASUREMENT_VARIABLE       = "measurement_variable"
	TABLE_ALLOWED_CYCLES             = "allowed_cycles"
	TABLE_TIME_CYCLES                = "time_cycles"
	TABLE_DATA_COLLECTION_TYPE       = "data_collection_type"
	TABLE_COLLECTED_DATA_TYPE        = "collected_data_type"
	TABLE_TIME_PERIODS               = "time_periods"
	TABLE_RECORDERS                  = "recorders"
	TABLE_RECORDER_TYPE              = "recorder_type"
	TABLE_MOVEMENT_TYPE              = "movement_type"
	TABLE_READING_REASONS            = "reading_reasons"
	TABLE_READING_TYPE               = "reading_type"
	TABLE_READING_STATUS             = "reading_status"
	TABLE_READING_STATE              = "reading_state"
	TABLE_ESTIMATION_METHODS         = "estimation_methods"

	// Tabelas de Fluxo e Motivação
	TABLE_YES_NO_RESPONSE            = "yes_no_response"
	TABLE_ACCEPTANCE_RESPONSE        = "acceptance_response"
	TABLE_HOLDER_CHANGE_CONTEXT      = "holder_change_context"
	TABLE_SOCIAL_TARIFF              = "social_tariff"
	TABLE_CANCELLATION_REASON        = "cancellation_reason"
	TABLE_TERMINATION_REASON         = "termination_reason"
	TABLE_SUSPENSION_REACTIVATION    = "suspension_reactivation"
	TABLE_PROCESS_SUSPENSION_REASONS = "process_suspension_reasons"
	TABLE_INCIDENCE_REASONS          = "incidence_reasons"
	TABLE_INCIDENCE_ORDINAL          = "incidence_ordinal"
	TABLE_INCIDENCE_RESPONSIBILITY   = "incidence_responsibility"
	TABLE_ANOMALY_TYPE_FRAUD         = "anomaly_type_fraud"
	TABLE_ORD_INFO_TYPE              = "ord_info_type"
	TABLE_PROOF_DOCUMENT_TYPE        = "proof_document_type"
	TABLE_DP_SUPPLY_STATUS           = "dp_supply_status"
	TABLE_PRODUCTION_TYPE            = "production_type"
	TABLE_SUBPROCESS_CODE            = "subprocess_code"
	TABLE_OBJECTION_REASONS          = "objection_reasons"
	TABLE_REFUSAL_REASONS            = "refusal_reasons"
	TABLE_ACTIVATION_TYPE            = "activation_type"
	TABLE_SERVICES_TO_PERFORM        = "services_to_perform"
	TABLE_COMMUNICATION_TYPE         = "communication_type"
	TABLE_ACCESS_TARIFFS             = "access_tariffs"
	TABLE_DP_NETWORK_AREA            = "dp_network_area"
	TABLE_RECIPIENT                  = "recipient"
	TABLE_DEADLINE_IDENTIFIERS       = "deadline_identifiers"
)

var TableCodeMap = map[string]string{
	// Tabelas Principais (Estrutura da Mensagem)
	"T00010": TABLE_PROCESSES,
	"T00020": TABLE_STEPS,
	"T00040": TABLE_RECORDS,
	"T00050": TABLE_FIELDS,
	"T00060": TABLE_HEADER_TYPES,
	"T00070": TABLE_RECORD_TYPES,

	// Tabelas de Erros
	"T05010": TABLE_SYNTAX_ERRORS,
	"T05020": TABLE_DATA_ERRORS,

	// Tabelas Geográficas e Agentes
	"T10020": TABLE_ELECTRICAL_UNITS,
	"T10051": TABLE_CAE_REV4,
	"T10110": TABLE_COUNTRIES,
	"T10120": TABLE_DISTRICTS,
	"T10130": TABLE_MUNICIPALITIES,
	"T10140": TABLE_PARISHES,
	"T10150": TABLE_INE_ZONES,
	"T10210": TABLE_POSTAL_CODES,
	"T10300": TABLE_AGENT_TYPES,
	"T10310": TABLE_NETWORK_OPERATORS,
	"T10320": TABLE_RETAILERS,
	"T10380": TABLE_LOGISTICS_OPERATOR,

	// Tabelas Técnicas (Ponto de Entrega / Instalação)
	"T12110": TABLE_SERVICE_QUALITY_ZONES,
	"T12210": TABLE_VOLTAGE_LEVELS,
	"T12215": TABLE_DELIVERY_POINT_TYPES,
	"T12217": TABLE_INSTALLATION_CHARS,
	"T12220": TABLE_MEASUREMENT_VOLTAGE,
	"T12225": TABLE_REPORTING_VOLTAGE,
	"T12230": TABLE_NUMBER_OF_PHASES,
	"T12510": TABLE_CONTRACTED_POWER,
	"T12520": TABLE_CONSUMPTION_PROFILE,
	"T12610": TABLE_DELIVERY_POINT_CONTACT,

	// Tabelas de Cliente e Identificação
	"T13010": TABLE_IDENTIFICATION_TYPES,
	"T13020": TABLE_CUSTOMER_TYPES,
	"T13110": TABLE_HOLDER_ADDRESS_TYPES,
	"T13210": TABLE_CNE_IDENTIFICATION,
	"T13220": TABLE_CNE_CONTACT,
	"T13230": TABLE_CNE_PREFERRED_CONTACT,
	"T13240": TABLE_DEFICIENCY_EQUIPMENT,
	"T13310": TABLE_PRIORITY_CUSTOMER_LOC,
	"T13320": TABLE_MR_CONTRACTING_REASONS,
	"T13330": TABLE_RPE_ACCESS_PURPOSE,

	// Tabelas de Medição e Equipamento
	"T14050": TABLE_EQUIPMENT_BRANDS,
	"T14060": TABLE_MEASUREMENT_DEVICE_TYPE,
	"T14070": TABLE_OWNERSHIP,
	"T14080": TABLE_MEASUREMENT_DEVICE_FUNC,
	"T14110": TABLE_SUPPORTED_MEASUREMENT_FUNC,
	"T14120": TABLE_MEASUREMENT_VARIABLE,
	"T14210": TABLE_ALLOWED_CYCLES,
	"T14220": TABLE_TIME_CYCLES,
	"T14310": TABLE_DATA_COLLECTION_TYPE,
	"T14320": TABLE_COLLECTED_DATA_TYPE,
	"T14410": TABLE_TIME_PERIODS,
	"T14420": TABLE_RECORDERS,
	"T14430": TABLE_RECORDER_TYPE,
	"T14510": TABLE_MOVEMENT_TYPE,
	"T14620": TABLE_READING_REASONS,
	"T14650": TABLE_READING_TYPE,
	"T14670": TABLE_READING_STATUS,
	"T14680": TABLE_READING_STATE,
	"T14810": TABLE_ESTIMATION_METHODS,

	// Tabelas de Fluxo e Motivação
	"T20100": TABLE_YES_NO_RESPONSE,
	"T20101": TABLE_ACCEPTANCE_RESPONSE,
	"T20102": TABLE_HOLDER_CHANGE_CONTEXT,
	"T20200": TABLE_SOCIAL_TARIFF,
	"T21151": TABLE_CANCELLATION_REASON,
	"T21510": TABLE_TERMINATION_REASON,
	"T23100": TABLE_SUSPENSION_REACTIVATION,
	"T23110": TABLE_PROCESS_SUSPENSION_REASONS,
	"T23160": TABLE_INCIDENCE_REASONS,
	"T23165": TABLE_INCIDENCE_ORDINAL,
	"T23166": TABLE_INCIDENCE_RESPONSIBILITY,
	"T23190": TABLE_ANOMALY_TYPE_FRAUD,
	"T23196": TABLE_ORD_INFO_TYPE,
	"T23200": TABLE_PROOF_DOCUMENT_TYPE,
	"T23210": TABLE_DP_SUPPLY_STATUS,
	"T23230": TABLE_PRODUCTION_TYPE,
	"T23250": TABLE_SUBPROCESS_CODE,
	"T24120": TABLE_OBJECTION_REASONS,
	"T24150": TABLE_REFUSAL_REASONS,
	"T25100": TABLE_ACTIVATION_TYPE,
	"T25140": TABLE_SERVICES_TO_PERFORM,
	"T25150": TABLE_COMMUNICATION_TYPE,
	"T26100": TABLE_ACCESS_TARIFFS,
	"T26110": TABLE_DP_NETWORK_AREA,
	"T29000": TABLE_RECIPIENT,
	"T29100": TABLE_DEADLINE_IDENTIFIERS,
}
