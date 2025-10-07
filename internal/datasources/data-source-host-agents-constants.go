package datasources

// Data source name constants
const (
	// DataSourceInstanaHostAgentsFramework the name of the terraform-provider-instana data source to read host agents
	DataSourceInstanaHostAgentsFramework = "host_agents"
	// DataSourceHostAgents the name of the terraform-provider-instana data source
	DataSourceHostAgents = "instana_host_agents"
)

// Field name constants for host agents
const (
	// HostAgentFieldFilter constant value for the schema field filter
	HostAgentFieldFilter = "filter"
	// HostAgentFieldItems constant value for the schema field items
	HostAgentFieldItems = "items"
	// HostAgentFieldHost constant value for the schema field host
	HostAgentFieldHost = "host"
	// HostAgentFieldLabel constant value for the schema field label
	HostAgentFieldLabel = "label"
	// HostAgentFieldSnapshotId constant value for the schema field snapshot_id
	HostAgentFieldSnapshotId = "snapshot_id"
	// HostAgentFieldPlugin constant value for the schema field plugin
	HostAgentFieldPlugin = "plugin"
	// HostAgentFieldTags constant value for the schema field tags
	HostAgentFieldTags = "tags"
)

// Field ID constant
const (
	// HostAgentFieldID constant value for the schema field id
	HostAgentFieldID = "id"
)

// Description constants for host agents fields
const (
	// HostAgentDescDataSource description for the data source
	HostAgentDescDataSource = "Data source for Instana host agents. Host agents are the Instana agents installed on hosts."
	// HostAgentDescID description for the ID field
	HostAgentDescID = "The ID of the data source."
	// HostAgentDescFilter description for the filter field
	HostAgentDescFilter = "Dynamic Focus Query filter."
	// HostAgentDescItems description for the items field
	HostAgentDescItems = "A list of host agents."
	// HostAgentDescSnapshotID description for the snapshot_id field
	HostAgentDescSnapshotID = "The snapshot ID of the host agent."
	// HostAgentDescLabel description for the label field
	HostAgentDescLabel = "The label of the host agent."
	// HostAgentDescHost description for the host field
	HostAgentDescHost = "The host identifier of the host agent."
	// HostAgentDescPlugin description for the plugin field
	HostAgentDescPlugin = "The plugin of the host agent."
	// HostAgentDescTags description for the tags field
	HostAgentDescTags = "The tags of the host agent."
)

// Error message constants
const (
	// HostAgentErrUnexpectedConfigureType error message for unexpected configure type
	HostAgentErrUnexpectedConfigureType = "Unexpected Data Source Configure Type"
	// HostAgentErrUnexpectedConfigureTypeDetail error message detail for unexpected configure type
	HostAgentErrUnexpectedConfigureTypeDetail = "Expected *restapi.ProviderMeta, got: %T. Please report this issue to the provider developers."
	// HostAgentErrReadingAgents error message for reading host agents
	HostAgentErrReadingAgents = "Error reading host agents"
	// HostAgentErrReadingAgentsDetail error message detail for reading host agents
	HostAgentErrReadingAgentsDetail = "Could not read host agents: %s"
)

// Query parameter constants
const (
	// HostAgentQueryParamQuery query parameter name
	HostAgentQueryParamQuery = "query"
)
