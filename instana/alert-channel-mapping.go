package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func MapAlertChannelsToState(ctx context.Context, alertChannels map[restapi.AlertSeverity][]string) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(alertChannels) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				ResourceFieldThresholdRuleWarningSeverity:  types.ListType{ElemType: types.StringType},
				ResourceFieldThresholdRuleCriticalSeverity: types.ListType{ElemType: types.StringType},
			},
		}), diags
	}

	alertChannelsObj := map[string]attr.Value{}

	// Map warning severity
	if warningChannels, ok := alertChannels[restapi.WarningSeverity]; ok && len(warningChannels) > 0 {
		warningList, warningDiags := types.ListValueFrom(ctx, types.StringType, warningChannels)
		diags.Append(warningDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		alertChannelsObj[ResourceFieldThresholdRuleWarningSeverity] = warningList
	} else {
		alertChannelsObj[ResourceFieldThresholdRuleWarningSeverity] = types.ListNull(types.StringType)
	}

	// Map critical severity
	if criticalChannels, ok := alertChannels[restapi.CriticalSeverity]; ok && len(criticalChannels) > 0 {
		criticalList, criticalDiags := types.ListValueFrom(ctx, types.StringType, criticalChannels)
		diags.Append(criticalDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{}), diags
		}
		alertChannelsObj[ResourceFieldThresholdRuleCriticalSeverity] = criticalList
	} else {
		alertChannelsObj[ResourceFieldThresholdRuleCriticalSeverity] = types.ListNull(types.StringType)
	}

	objVal, objDiags := types.ObjectValue(
		map[string]attr.Type{
			ResourceFieldThresholdRuleWarningSeverity:  types.ListType{ElemType: types.StringType},
			ResourceFieldThresholdRuleCriticalSeverity: types.ListType{ElemType: types.StringType},
		},
		alertChannelsObj,
	)
	diags.Append(objDiags...)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}

	return types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				ResourceFieldThresholdRuleWarningSeverity:  types.ListType{ElemType: types.StringType},
				ResourceFieldThresholdRuleCriticalSeverity: types.ListType{ElemType: types.StringType},
			},
		},
		[]attr.Value{objVal},
	)
}

func MapAlertChannelsFromState(ctx context.Context, alertChannelsList types.List) (map[restapi.AlertSeverity][]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	alertChannelsMap := make(map[restapi.AlertSeverity][]string)

	if alertChannelsList.IsNull() || alertChannelsList.IsUnknown() {
		return alertChannelsMap, diags
	}

	var alertChannelsElements []types.Object
	diags.Append(alertChannelsList.ElementsAs(ctx, &alertChannelsElements, false)...)
	if diags.HasError() {
		return nil, diags
	}

	if len(alertChannelsElements) == 0 {
		return alertChannelsMap, diags
	}

	var alertChannels struct {
		Warning  types.List `tfsdk:"warning"`
		Critical types.List `tfsdk:"critical"`
	}

	diags.Append(alertChannelsElements[0].As(ctx, &alertChannels, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil, diags
	}

	// Map warning severity
	if !alertChannels.Warning.IsNull() && !alertChannels.Warning.IsUnknown() {
		var warningChannels []string
		diags.Append(alertChannels.Warning.ElementsAs(ctx, &warningChannels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(warningChannels) > 0 {
			alertChannelsMap[restapi.WarningSeverity] = warningChannels
		}
	}

	// Map critical severity
	if !alertChannels.Critical.IsNull() && !alertChannels.Critical.IsUnknown() {
		var criticalChannels []string
		diags.Append(alertChannels.Critical.ElementsAs(ctx, &criticalChannels, false)...)
		if diags.HasError() {
			return nil, diags
		}
		if len(criticalChannels) > 0 {
			alertChannelsMap[restapi.CriticalSeverity] = criticalChannels
		}
	}

	return alertChannelsMap, diags
}
