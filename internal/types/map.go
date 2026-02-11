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
	"T00010": TABLE_PROCESSES,    // Tabela de processos
	"T00020": TABLE_STEPS,        // Tabela de passos
	"T00040": TABLE_RECORDS,      // Tabela de registos
	"T00050": TABLE_FIELDS,       // Tabela de campos
	"T00060": TABLE_HEADER_TYPES, // Tipos de cabeçalho
	"T00070": TABLE_RECORD_TYPES, // Tipos de registo

	// Tabelas de Erros
	"T05010": TABLE_SYNTAX_ERRORS, // Erros de sintaxe
	"T05020": TABLE_DATA_ERRORS,   // Erros de preenchimento

	// Tabelas Geográficas e Agentes
	"T10020": TABLE_ELECTRICAL_UNITS,   // Unidades de medida eletricas
	"T10051": TABLE_CAE_REV4,           // Tabela CAE Rev.4
	"T10110": TABLE_COUNTRIES,          // Tabela de Países
	"T10120": TABLE_DISTRICTS,          // Tabela de Distritos
	"T10130": TABLE_MUNICIPALITIES,     // Tabela de Concelhos
	"T10140": TABLE_PARISHES,           // Tabela de Freguesias
	"T10150": TABLE_INE_ZONES,          // Tabela de Zonas INE
	"T10210": TABLE_POSTAL_CODES,       // Tabela de Códigos Postais
	"T10300": TABLE_AGENT_TYPES,        // Tipos de agentes
	"T10310": TABLE_NETWORK_OPERATORS,  // Tabela de operadores de Rede
	"T10320": TABLE_RETAILERS,          // Tabela de comercializadores
	"T10380": TABLE_LOGISTICS_OPERATOR, // Operador Logístico de Mudança de Comercializador

	// Tabelas Técnicas (Ponto de Entrega / Instalação)
	"T12110": TABLE_SERVICE_QUALITY_ZONES,  // Zonas de qualidade de serviço
	"T12210": TABLE_VOLTAGE_LEVELS,         // Níveis de tensão
	"T12215": TABLE_DELIVERY_POINT_TYPES,   // Tipos de pontos de entrega
	"T12217": TABLE_INSTALLATION_CHARS,     // Características da instalação
	"T12220": TABLE_MEASUREMENT_VOLTAGE,    // Tensão de medida
	"T12225": TABLE_REPORTING_VOLTAGE,      // Tensão de medida reporting
	"T12230": TABLE_NUMBER_OF_PHASES,       // Número de fases
	"T12510": TABLE_CONTRACTED_POWER,       // Potência contratada
	"T12520": TABLE_CONSUMPTION_PROFILE,    // Perfil de consumo
	"T12610": TABLE_DELIVERY_POINT_CONTACT, // Contacto do PE

	// Tabelas de Cliente e Identificação
	"T13010": TABLE_IDENTIFICATION_TYPES,   // Tipos de identificação
	"T13020": TABLE_CUSTOMER_TYPES,         // Tipos de cliente
	"T13110": TABLE_HOLDER_ADDRESS_TYPES,   // Tipos de morada titular
	"T13210": TABLE_CNE_IDENTIFICATION,     // Identificação do CNE
	"T13220": TABLE_CNE_CONTACT,            // Contacto do CNE
	"T13230": TABLE_CNE_PREFERRED_CONTACT,  // Contacto preferencial do CNE
	"T13240": TABLE_DEFICIENCY_EQUIPMENT,   // Tipos de deficiência/Equipamento
	"T13310": TABLE_PRIORITY_CUSTOMER_LOC,  // Tipos de local Cliente Prioritário
	"T13320": TABLE_MR_CONTRACTING_REASONS, // Motivos de contratação MR
	"T13330": TABLE_RPE_ACCESS_PURPOSE,     // Finalidade do acesso ao RPE

	// Tabelas de Medição e Equipamento
	"T14050": TABLE_EQUIPMENT_BRANDS,           // Marcas de equipamentos
	"T14060": TABLE_MEASUREMENT_DEVICE_TYPE,    // Tipo de aparelho de medição
	"T14070": TABLE_OWNERSHIP,                  // Propriedade
	"T14080": TABLE_MEASUREMENT_DEVICE_FUNC,    // Função do aparelho de medição
	"T14110": TABLE_SUPPORTED_MEASUREMENT_FUNC, // Funções de medida suportadas
	"T14120": TABLE_MEASUREMENT_VARIABLE,       // Variavél de medida
	"T14210": TABLE_ALLOWED_CYCLES,             // Ciclos permitidos
	"T14220": TABLE_TIME_CYCLES,                // Ciclos horários
	"T14310": TABLE_DATA_COLLECTION_TYPE,       // Tipo de recolha de dados
	"T14320": TABLE_COLLECTED_DATA_TYPE,        // Tipo de dados recolhidos
	"T14410": TABLE_TIME_PERIODS,               // Períodos horários
	"T14420": TABLE_RECORDERS,                  // Registadores
	"T14430": TABLE_RECORDER_TYPE,              // Tipo de registador
	"T14510": TABLE_MOVEMENT_TYPE,              // Tipo de movimento
	"T14620": TABLE_READING_REASONS,            // Motivos de leitura
	"T14650": TABLE_READING_TYPE,               // Tipo de leitura
	"T14670": TABLE_READING_STATUS,             // Status da leitura
	"T14680": TABLE_READING_STATE,              // Estado da leitura
	"T14810": TABLE_ESTIMATION_METHODS,         // Metodos de estimativa

	// Tabelas de Fluxo e Motivação
	"T20100": TABLE_YES_NO_RESPONSE,            // Resposta S/N
	"T20101": TABLE_ACCEPTANCE_RESPONSE,        // Resposta S/N - Aceitação
	"T20102": TABLE_HOLDER_CHANGE_CONTEXT,      // Contexto Alteração de Titular
	"T20200": TABLE_SOCIAL_TARIFF,              // Tarifa social
	"T21151": TABLE_CANCELLATION_REASON,        // Código de motivo de anulação
	"T21510": TABLE_TERMINATION_REASON,         // Código de motivo de denúncia
	"T23100": TABLE_SUSPENSION_REACTIVATION,    // Suspensão/Reactivação
	"T23110": TABLE_PROCESS_SUSPENSION_REASONS, // Motivos de suspensão/reactivação do processo
	"T23160": TABLE_INCIDENCE_REASONS,          // Motivos de incidência
	"T23165": TABLE_INCIDENCE_ORDINAL,          // Ordinal de incidência
	"T23166": TABLE_INCIDENCE_RESPONSIBILITY,   // Responsabilidade da incidência
	"T23190": TABLE_ANOMALY_TYPE_FRAUD,         // Tipo de anomalia - Fraude
	"T23196": TABLE_ORD_INFO_TYPE,              // Tipo de informação facultada pelo ORD
	"T23200": TABLE_PROOF_DOCUMENT_TYPE,        // Tipo de documento comprovativo
	"T23210": TABLE_DP_SUPPLY_STATUS,           // Estado do fornecimento do PE
	"T23230": TABLE_PRODUCTION_TYPE,            // Tipo de produção
	"T23250": TABLE_SUBPROCESS_CODE,            // Código de subprocesso
	"T24120": TABLE_OBJECTION_REASONS,          // Motivos de objeção
	"T24150": TABLE_REFUSAL_REASONS,            // Motivos de recusa
	"T25100": TABLE_ACTIVATION_TYPE,            // Tipo de ativação
	"T25140": TABLE_SERVICES_TO_PERFORM,        // Serviços a efectuar
	"T25150": TABLE_COMMUNICATION_TYPE,         // Tipo de comunicação
	"T26100": TABLE_ACCESS_TARIFFS,             // Tarifas de acesso
	"T26110": TABLE_DP_NETWORK_AREA,            // Área de rede do ponto de entrega
	"T29000": TABLE_RECIPIENT,                  // Destinatário
	"T29100": TABLE_DEADLINE_IDENTIFIERS,       // Identificadores dos prazos
}
