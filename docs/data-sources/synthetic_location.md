# Synthetic Location Data Source

Data source to get the synthetic locations from Instana API. This allows you to retrieve the specification
by label and location type and reference it in other resources such as Synthetic tests.

API Documentation: <https://instana.github.io/openapi/#operation/getSyntheticLocation>

## Example Usage

```hcl
data "instana_synthetic_location" "locations" {}
```

## Argument Reference

* `label` - Optional - the label of the synthetic location
* `location_type` - Optional - indicates if the location is public or private

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - The synthetic location identifier
* `description` - The description of the synthetic location
