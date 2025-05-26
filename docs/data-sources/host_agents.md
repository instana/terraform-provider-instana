# Host Agents Data Source

Data source to get the host agents from Instana API. This allows you to retrieve the host agents
using Dynamic Focus Query filter and reference them in other resources such as Automation Policy.

API Documentation: <https://instana.github.io/openapi/#tag/Host-Agent>

## Example Usage

```hcl
data "instana_host_agents" "sample" {
  filter = "entity.agent.capability:action-script"
}
```

## Argument Reference

* `filter` - Required - Dynamic Focus Query filter.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `items` - List of available host agents.
    * `snapshot_id` - The snapshot ID of the host agent.
    * `label` - The label of the host agent.
    * `host` - The host identifier of the host agent.
    * `plugin` - The plugin name.
    * `tags` - List of host agent tags.