package instana

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DataSourceHostAgents = "instana_host_agents"

	HostAgentFieldFilter     = "filter"
	HostAgentFieldItems      = "items"
	HostAgentFieldHost       = "host"
	HostAgentFieldLabel      = "label"
	HostAgentFieldSnapshotId = "snapshot_id"
	HostAgentFieldPlugin     = "plugin"
	HostAgentFieldTags       = "tags"
)

// NewHostAgentsDataSource creates a new DataSource for Host Agents
func NewHostAgentsDataSource() DataSource {
	return &hostAgentsDataSource{}
}

type hostAgentsDataSource struct{}

// CreateResource creates the terraform Resource for the data source for Instana host agents
func (ds *hostAgentsDataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		ReadContext: ds.read,
		Schema: map[string]*schema.Schema{
			HostAgentFieldFilter: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Dynamic Focus Query filter.",
			},
			HostAgentFieldItems: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of host agents.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						HostAgentFieldSnapshotId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spanshot ID of the host agent.",
						},
						HostAgentFieldLabel: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The label of the host agent.",
						},
						HostAgentFieldHost: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host identifier of the host agent.",
						},
						HostAgentFieldPlugin: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The plugin of the host agent.",
						},
						HostAgentFieldTags: {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "The tags of the host agent.",
						},
					},
				},
			},
		},
	}
}

func (ds *hostAgentsDataSource) read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	filter := d.Get(HostAgentFieldFilter).(string)
	queryParams := map[string]string{"query": filter}

	data, err := instanaAPI.HostAgents().GetByQuery(queryParams)
	if err != nil {
		return diag.FromErr(err)
	}

	err = ds.updateState(d, data)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func (ds *hostAgentsDataSource) updateState(d *schema.ResourceData, hostAgents *[]*restapi.HostAgent) error {
	d.SetId(time.Now().UTC().String())
	return tfutils.UpdateState(d, map[string]interface{}{
		HostAgentFieldItems: ds.mapHostAgentsToSchema(hostAgents),
	})
}

func (ds *hostAgentsDataSource) mapHostAgentsToSchema(hostAgents *[]*restapi.HostAgent) []interface{} {
	result := make([]interface{}, len(*hostAgents))
	i := 0
	for _, hostAgent := range *hostAgents {
		item := make(map[string]interface{})
		item[HostAgentFieldSnapshotId] = hostAgent.SnapshotID
		item[HostAgentFieldLabel] = hostAgent.Label
		item[HostAgentFieldHost] = hostAgent.Host
		item[HostAgentFieldPlugin] = hostAgent.Plugin
		item[HostAgentFieldTags] = hostAgent.Tags
		result[i] = item
		i++
	}
	return result
}
