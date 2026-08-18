## [v7.6.0](https://github.com/instana/terraform-provider-instana/tree/v7.6.0) (2026-08-18)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.5.0...v7.6.0)

**Implemented enhancements:**

- Added new `instana_synthetic_credential` resource to manage synthetic credentials in Instana, including `credential_name`, write-only sensitive `credential_value`, optional application/mobile/website scoping, and RBAC tag associations
- Added resource support for write-only synthetic credential secrets by preserving `credential_value` in state and allowing import by `credential_name`
- Bumped `instana-go-client` dependency from `v1.2.0` to `v1.3.0` to support synthetic credential APIs
- Added resource documentation for `instana_synthetic_credential` with usage examples and Terraform import instructions

**Merged pull requests:**

- added new resource : synthetic credentials [\#113](https://github.com/instana/terraform-provider-instana/pull/113) ([BlessyElza](https://github.com/BlessyElza))

## [v7.5.0](https://github.com/instana/terraform-provider-instana/tree/v7.5.0) (2026-08-17)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.4.2...v7.5.0)

**Implemented enhancements:**

- Added new `instana_session_settings` resource to manage session settings in Instana, supporting configuration of session timeout and inactivity timeout durations

**Merged pull requests:**

- add instana_session_settings resource [\#112](https://github.com/instana/terraform-provider-instana/pull/112) ([georgekutty-1](https://github.com/georgekutty-1))

## [v7.4.2](https://github.com/instana/terraform-provider-instana/tree/v7.4.2) (2026-08-06)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.4.1...v7.4.2)

**Fixed bugs:**

- Fixed decimal precision loss in `instana_slo_config` and `instana_slo_alert_config` threshold state mapping by removing `FormatFloat`/`ParseFloat` round-trips that truncated values like `0.9995` to `1.0`, resolving Terraform's "inconsistent result after apply" error

**Merged pull requests:**

- fix(slo): fix decimal precision loss in threshold state mapping [\#110](https://github.com/instana/terraform-provider-instana/pull/110) ([BlessyElza](https://github.com/BlessyElza))

## [v7.4.1](https://github.com/instana/terraform-provider-instana/tree/v7.4.1) (2026-08-04)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.4.0...v7.4.1)

**Fixed bugs:**

- Fixed perpetual diff after `terraform apply` in `instana_rbac_role` by changing the `member` field from `[]RoleMemberModel` to `types.Set` and adding a `UseStateForUnknown` plan modifier
- Fixed perpetual diff after `terraform apply` in `instana_rbac_team` by changing all set-typed scope fields from `[]string` to `types.Set` and adding `UseStateForUnknown` plan modifiers
- Fixed panic `fatal error: concurrent map writes` caused by unsynchronised access to a shared map during parallel Terraform operations

**Merged pull requests:**

- fix(rbac): resolve apply consistency errors for instana_rbac_team scope and instana_rbac_role member [\#108](https://github.com/instana/terraform-provider-instana/pull/108) ([BlessyElza](https://github.com/BlessyElza))
- fix: resolve concurrent map writes panic [\#109](https://github.com/instana/terraform-provider-instana/pull/109) ([BlessyElza](https://github.com/BlessyElza))

## [v7.4.0](https://github.com/instana/terraform-provider-instana/tree/v7.4.0) (2026-07-28)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.3.4...v7.4.0)

**Implemented enhancements:**

- Added new `instana_rbac_group_mapping` resource to manage RBAC group mappings in Instana, mapping IdP (LDAP, OIDC, SAML) attribute key/value pairs to Instana roles for automatic user role assignment on login
- New resource supports `key`, `value`, `group_id` (required) and optional `team_id` fields
- Bumped `instana-go-client` dependency from `v1.1.3` to `v1.1.4` to include the new `GroupMapping` API struct and `GroupMappings()` endpoint
- Added resource documentation for `instana_rbac_group_mapping` with usage examples and Terraform import instructions

**Merged pull requests:**

- RBAC group mapping resource [\#106](https://github.com/instana/terraform-provider-instana/pull/106) ([blessyelzabyju](https://github.com/BlessyElza))

## [v7.3.4](https://github.com/instana/terraform-provider-instana/tree/v7.3.4) (2026-07-14)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.3.3...v7.3.4)

**Implemented enhancements:**

- Added `slack_app` and `ms_teams_app` channel types to the `instana_alerting_channel` data source, bringing it to full parity with the resource
- Updated `instana_alerting_channel` data source documentation to include `slack_app` and `ms_teams_app` attribute reference and field descriptions

**Merged pull requests:**

- Alert channel datasource - slack and teams app [\#103](https://github.com/instana/terraform-provider-instana/pull/103) ([georgekutty-1](https://github.com/georgekutty-1))

## [v7.3.3](https://github.com/instana/terraform-provider-instana/tree/v7.3.3) (2026-07-14)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.3.2...v7.3.3)

**Implemented enhancements:**

- Added Mobile App SLO support to `instana_slo_config` resource with `mobile` entity block, supporting `mobile_ids` and `filter_expression`
- Bumped `instana-go-client` dependency from `v1.1.2` to `v1.1.3` to include new `MobileIds`, `SloAdvancedFilter`, and `SloEntityMetric` API model fields
- Added application config mandatory validation for tag filter
- Updated resource documentation for `instana_alerting_channel`, `instana_custom_dashboard`, and `instana_synthetic_test` to correctly describe RBAC tag usage

**Fixed bugs:**

- Fixed custom dashboard creation with `rbac_tags` by adding a generic post-create update hook that issues a follow-up Update after Create to preserve the rbac fields

**Merged pull requests:**

- Mobile slo suport [\#102](https://github.com/instana/terraform-provider-instana/pull/102) ([georgekutty-1](https://github.com/georgekutty-1))
- Rbac tags documentation and custom dashboard fix [\#101](https://github.com/instana/terraform-provider-instana/pull/101) ([georgekutty-1](https://github.com/georgekutty-1))

## [v7.3.2](https://github.com/instana/terraform-provider-instana/tree/v7.3.2) (2026-07-09)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.3.1...v7.3.2)

**Implemented enhancements:**

- Added `rbac_tags` support to `instana_alerting_channel`, `instana_custom_dashboard` and `instana_synthetic_test` resources

**Fixed bugs:**

- Provide default value (`Any`) for website alert rule `value` field when `operator = "NOT_EMPTY"` to prevent broken GUI configurations
- Fixed `CustomPayloadFieldsToTerraform` returning `null` instead of empty list

**Merged pull requests:**

- Rbac tags addition and custom payload fix [\#99](https://github.com/instana/terraform-provider-instana/pull/99) ([georgekutty-1](https://github.com/georgekutty-1))
- Provide default value for website alert message for Not Empty Operator [\#96](https://github.com/instana/terraform-provider-instana/pull/96) ([rmainwork](https://github.com/rmainwork))

## [v7.3.1](https://github.com/instana/terraform-provider-instana/tree/v7.3.1) (2026-07-01)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.3.0...v7.3.1)

**Implemented enhancements:**

- Added `grace_period` field support to Infrastructure and Website Alert Configuration resources
- Added Apdex smart alert support to SLO Alert Configuration resource

**Merged pull requests:**

- Add grace_period field support to InfraAlertConfig and WebsiteAlertConfig [\#98](https://github.com/instana/terraform-provider-instana/pull/98) ([georgekutty-1](https://github.com/georgekutty-1))
- Add Apdex smart alert support to instana_slo_alert_config [\#97](https://github.com/instana/terraform-provider-instana/pull/97) ([dhinesh-sr](https://github.com/dhinesh-sr))

## [v7.3.0](https://github.com/instana/terraform-provider-instana/tree/v7.3.0) (2026-06-29)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.2.0...v7.3.0)

**Implemented enhancements:**

- Added support for Apdex Configuration resource to manage application and website Apdex thresholds

**Merged pull requests:**

- Add Terraform support for Apdex Configuration resource [\#92](https://github.com/instana/terraform-provider-instana/pull/92) ([dhinesh-sr](https://github.com/dhinesh-sr))

## [v7.2.0](https://github.com/instana/terraform-provider-instana/tree/v7.2.0) (2026-06-17)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.1.2...v7.2.0)

**Implemented enhancements:**

- Added new data source for retrieving Instana RBAC teams

**Merged pull requests:**

- Add RBAC team data source and disabled legacy GitHub release workflow [\#90](https://github.com/instana/terraform-provider-instana/pull/90) ([georgekutty-1](https://github.com/georgekutty-1))

## [v7.1.2](https://github.com/instana/terraform-provider-instana/tree/v7.1.2) (2026-06-12)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.1.1...v7.1.2)

**Implemented enhancements:**

- Added 7 API token permission fields to achieve complete API parity
  - `can_collect_net_trace_logs` - Permission to collect network trace logs
  - `can_configure_custom_entities` - Permission to configure custom entities
  - `can_configure_website_conversions` - Permission to configure website conversions
  - `can_configure_ip_filtering` - Permission to configure IP filtering
  - `can_configure_llm_model_price` - Permission to configure LLM model pricing
  - `can_configure_personally_identifiable_information_masking` - Permission to configure PII masking
  - `can_download_agent_configuration` - Permission to download agent configuration

**Merged pull requests:**

- API token resource enhancement [\#89](https://github.com/instana/terraform-provider-instana/pull/89) ([georgekutty-1](https://github.com/georgekutty-1))

## [v7.1.1](https://github.com/instana/terraform-provider-instana/tree/v7.1.1) (2026-06-10)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.1.0...v7.1.1)

**Implemented enhancements:**

- Updated RBAC permissions to align with the latest Instana API and added new roles
- Added support for `IncludeUnscheduledTestResults` field in SLO synthetic entity configurations

**Fixed bugs:**

- Added provider validation for API token permissions that are normalized by the backend to prevent inconsistent-result-after-apply errors
- Improved validation to reject invalid synthetic permission combinations during planning phase

**Merged pull requests:**

- RBAC permissions update and SLO configuration update [\#83](https://github.com/instana/terraform-provider-instana/pull/83) ([georgekutty-1](https://github.com/georgekutty-1))
- Validate backend-normalized API token permissions in provider [\#88](https://github.com/instana/terraform-provider-instana/pull/88) ([georgekutty-1](https://github.com/georgekutty-1))

# Changelog
## [v7.1.0](https://github.com/instana/terraform-provider-instana/tree/v7.1.0) (2026-06-03)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.0.1...v7.1.0)

**Implemented enhancements:**

- Added new data source for retrieving Instana RBAC roles

**Merged pull requests:**

- Add RBAC roles data source [\#85](https://github.com/instana/terraform-provider-instana/pull/85) ([nicoleyson](https://github.com/nicoleyson))


## [v7.0.1](https://github.com/instana/terraform-provider-instana/tree/v7.0.1) (2026-05-21)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v7.0.0...v7.0.1)

**Implemented enhancements:**

- Enhanced provider logging for better troubleshooting capabilities

**Merged pull requests:**

- Add debug logging for Terraform provider [\#82](https://github.com/instana/terraform-provider-instana/pull/82) ([BlessyElza](https://github.com/BlessyElza))

## [v7.0.0](https://github.com/instana/terraform-provider-instana/tree/v7.0.0) (2026-04-20)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.3.0...v7.0.0)

**changes:**

- Extracted Instana API Client into a standalone Go library for better modularity and reusability

**Merged pull requests:**

- Extract Instana API Client into Standalone Go Library [\#78](https://github.com/instana/terraform-provider-instana/pull/78) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.3.0](https://github.com/instana/terraform-provider-instana/tree/v6.3.0) (2026-03-30)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.2.0...v6.3.0)

**Implemented enhancements:**

- Added support for Mobile App Configuration resource for Mobile App Monitoring

**Merged pull requests:**

- Mobile app config resource [\#77](https://github.com/instana/terraform-provider-instana/pull/77) ([BlessyElza](https://github.com/BlessyElza))

## [v6.2.0](https://github.com/instana/terraform-provider-instana/tree/v6.2.0) (2026-03-09)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.1.2...v6.2.0)

**Implemented enhancements:**

- Added support for Maintenance Window resource configuration
- Added support for Mobile Alert configuration

**Merged pull requests:**

- Maintenance window resource [\#76](https://github.com/instana/terraform-provider-instana/pull/76) ([georgekutty-1](https://github.com/georgekutty-1))
- Mobile alert config [\#75](https://github.com/instana/terraform-provider-instana/pull/75) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.1.2](https://github.com/instana/terraform-provider-instana/tree/v6.1.2) (2026-02-19)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.1.1...v6.1.2)

**Implemented enhancements:**

- Enable/Disable support for application and website alerts

**Merged pull requests:**

- Enable/Disable support for application and website alerts [\#74](https://github.com/instana/terraform-provider-instana/pull/74) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.1.1](https://github.com/instana/terraform-provider-instana/tree/v6.1.1) (2026-02-05)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.1.0...v6.1.1)

**Fixed bugs:**

- Fix state management issue in API token resource

**Merged pull requests:**

- api-token-resource fix [\#73](https://github.com/instana/terraform-provider-instana/pull/73) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.1.0](https://github.com/instana/terraform-provider-instana/tree/v6.1.0) (2026-01-15)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.0.4...v6.1.0)

**Implemented enhancements:**

- Added new instana_user data source that allows users to retrieve Instana user details by email address.

**Merged pull requests:**

- Datasource user [\#71](https://github.com/instana/terraform-provider-instana/pull/71) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.0.4](https://github.com/instana/terraform-provider-instana/tree/v6.0.4) (2026-01-08)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.0.3...v6.0.4)

**Fixed bugs:**
- Eliminate false positive changes from list ordering by converting to Sets

**Implemented enhancements:**

- Added support for creating and managing Infrastructure SLOs using the newly introduced Saturation blueprint

**Merged pull requests:**

- List ordering fix [\#70](https://github.com/instana/terraform-provider-instana/pull/70) ([georgekutty-1](https://github.com/georgekutty-1))
- Add Saturation Blueprint support to SLO configs [\#69](https://github.com/instana/terraform-provider-instana/pull/69) ([nikhilgowda123](https://github.com/nikhilgowda123))


## [v6.0.3](https://github.com/instana/terraform-provider-instana/tree/v6.0.3) (2025-12-17)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.0.2...v6.0.3)

**Merged pull requests:**

- Resource documentation update and bug fixes [\#68](https://github.com/instana/terraform-provider-instana/pull/68) ([georgekutty-1](https://github.com/georgekutty-1))


## [v6.0.2](https://github.com/instana/terraform-provider-instana/tree/v6.0.2) (2025-12-10)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.0.1...v6.0.2)

**Merged pull requests:**

- Resource documentation update [\#66](https://github.com/instana/terraform-provider-instana/pull/66) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.0.1](https://github.com/instana/terraform-provider-instana/tree/v6.0.1) (2025-12-05)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v6.0.0...v6.0.1)

**Merged pull requests:**

- Provider migration [\#65](https://github.com/instana/terraform-provider-instana/pull/65) ([georgekutty-1](https://github.com/georgekutty-1))

## [v6.0.0](https://github.com/instana/terraform-provider-instana/tree/v6.0.0) (2025-12-04)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.4.1...v6.0.0)

**Merged pull requests:**

- Provider migration [\#63](https://github.com/instana/terraform-provider-instana/pull/63) ([georgekutty-1](https://github.com/georgekutty-1))


## [v5.4.1](https://github.com/instana/terraform-provider-instana/tree/v5.4.1) (2025-12-03)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.4.0...v5.4.1)

**Merged pull requests:**

- Add support for relative_diff & absolute_diff aggregation in custom event specification [\#62](https://github.com/instana/terraform-provider-instana/pull/62) ([parekh-raj](https://github.com/parekh-raj))

## [v5.4.0](https://github.com/instana/terraform-provider-instana/tree/v5.4.0) (2025-10-09)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.3.1...v5.4.0)

**Implemented enhancements:**

- Log smart alert resource [\#57](https://github.com/instana/terraform-provider-instana/pull/57) ([georgekutty-1](https://github.com/georgekutty-1))

## [v5.3.1](https://github.com/instana/terraform-provider-instana/tree/v5.3.1) (2025-09-30)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.3.0...v5.3.1)

**Merged pull requests:**

- Update index.md to highlight Synthetic Alert Config support [\#56](https://github.com/instana/terraform-provider-instana/pull/56) ([parekh-raj](https://github.com/parekh-raj))

## [v5.3.0](https://github.com/instana/terraform-provider-instana/tree/v5.3.0) (2025-09-24)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.2.0...v5.3.0)

**Merged pull requests:**

- Synthetic monitoring alert resource handle [\#55](https://github.com/instana/terraform-provider-instana/pull/55) ([georgekutty-1](https://github.com/georgekutty-1))

## [v5.2.0](https://github.com/instana/terraform-provider-instana/tree/v5.2.0) (2025-08-21)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.1.1...v5.2.0)

**Implemented enhancements:**

- Extend instana\_infra\_alert\_config resource doc for evaluation\_type field [\#50](https://github.com/instana/terraform-provider-instana/pull/50) ([parekh-raj](https://github.com/parekh-raj))

**Merged pull requests:**

- Group permissions update [\#53](https://github.com/instana/terraform-provider-instana/pull/53) ([rorywelch](https://github.com/rorywelch))
- Adding Timezone to the SLO payload [\#52](https://github.com/instana/terraform-provider-instana/pull/52) ([dhinesh-sr](https://github.com/dhinesh-sr))

## [v5.1.1](https://github.com/instana/terraform-provider-instana/tree/v5.1.1) (2025-07-28)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.1.0...v5.1.1)

## [v5.1.0](https://github.com/instana/terraform-provider-instana/tree/v5.1.0) (2025-07-24)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v5.0.0...v5.1.0)

**Implemented enhancements:**

- Extend instana\_infra\_alert\_config schema with evaluation\_type [\#46](https://github.com/instana/terraform-provider-instana/pull/46) ([parekh-raj](https://github.com/parekh-raj))

**Merged pull requests:**

- Importing an Application config fix [\#48](https://github.com/instana/terraform-provider-instana/pull/48) ([rorywelch](https://github.com/rorywelch))

## [v5.0.0](https://github.com/instana/terraform-provider-instana/tree/v5.0.0) (2025-07-03)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v4.0.5...v5.0.0)

**Implemented enhancements:**

- API Token Permission Support Update [\#45](https://github.com/instana/terraform-provider-instana/pull/45) ([ChinmayGitHub](https://github.com/ChinmayGitHub))
- Update group permissions  [\#44](https://github.com/instana/terraform-provider-instana/pull/44) ([ChinmayGitHub](https://github.com/ChinmayGitHub))

**Closed issues:**

- Newer permissions should be added to the list of allowed permissions. [\#20](https://github.com/instana/terraform-provider-instana/issues/20)

**Merged pull requests:**

- Created SLO correction configuration resources [\#42](https://github.com/instana/terraform-provider-instana/pull/42) ([dhinesh-sr](https://github.com/dhinesh-sr))

## [v4.0.5](https://github.com/instana/terraform-provider-instana/tree/v4.0.5) (2025-06-30)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v4.0.4...v4.0.5)

**Implemented enhancements:**

- Use correct tag filter as part of examples in the application\_config.md [\#43](https://github.com/instana/terraform-provider-instana/pull/43) ([parekh-raj](https://github.com/parekh-raj))

## [v4.0.4](https://github.com/instana/terraform-provider-instana/tree/v4.0.4) (2025-06-17)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v4.0.3...v4.0.4)

**Merged pull requests:**

- Fix mapAlertChannelsToSchema [\#41](https://github.com/instana/terraform-provider-instana/pull/41) ([parekh-raj](https://github.com/parekh-raj))

## [v4.0.3](https://github.com/instana/terraform-provider-instana/tree/v4.0.3) (2025-06-16)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v4.0.2...v4.0.3)

**Merged pull requests:**

- Fix alert\_channels issue in Infra Smart Alert [\#40](https://github.com/instana/terraform-provider-instana/pull/40) ([parekh-raj](https://github.com/parekh-raj))
- Include the automation framework resources to the doc [\#39](https://github.com/instana/terraform-provider-instana/pull/39) ([epostea](https://github.com/epostea))

## [v4.0.2](https://github.com/instana/terraform-provider-instana/tree/v4.0.2) (2025-06-02)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v4.0.1...v4.0.2)

**Merged pull requests:**

- Implement the automation framework resources [\#38](https://github.com/instana/terraform-provider-instana/pull/38) ([epostea](https://github.com/epostea))
- Update SLO Burn Rate Smart Alert with v2 [\#37](https://github.com/instana/terraform-provider-instana/pull/37) ([nikhilgowda123](https://github.com/nikhilgowda123))

## [v4.0.1](https://github.com/instana/terraform-provider-instana/tree/v4.0.1) (2025-05-26)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v4.0.0...v4.0.1)

**Merged pull requests:**

- Fix call to Instana API always recieving a http 500 status response code for application configs [\#36](https://github.com/instana/terraform-provider-instana/pull/36) ([rorywelch](https://github.com/rorywelch))

## [v4.0.0](https://github.com/instana/terraform-provider-instana/tree/v4.0.0) (2025-04-10)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.3.0...v4.0.0)

**Merged pull requests:**

- Update alerting permissions [\#30](https://github.com/instana/terraform-provider-instana/pull/30) ([parekh-raj](https://github.com/parekh-raj))

## [v3.3.0](https://github.com/instana/terraform-provider-instana/tree/v3.3.0) (2025-04-04)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.2.1...v3.3.0)

**Implemented enhancements:**

- Make tag\_filter optional in host availability rule [\#31](https://github.com/instana/terraform-provider-instana/pull/31) ([parekh-raj](https://github.com/parekh-raj))

## [v3.2.1](https://github.com/instana/terraform-provider-instana/tree/v3.2.1) (2025-03-26)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.2.0...v3.2.1)

## [v3.2.0](https://github.com/instana/terraform-provider-instana/tree/v3.2.0) (2025-03-21)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.1.2...v3.2.0)

**Merged pull requests:**

- Created SLO configuration and SLO smart alert resources [\#29](https://github.com/instana/terraform-provider-instana/pull/29) ([dhinesh-sr](https://github.com/dhinesh-sr))

## [v3.1.2](https://github.com/instana/terraform-provider-instana/tree/v3.1.2) (2025-03-03)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.1.1...v3.1.2)

**Merged pull requests:**

- Allow 1m granularity for smart alert resources [\#28](https://github.com/instana/terraform-provider-instana/pull/28) ([parekh-raj](https://github.com/parekh-raj))
- Set up GPG in .github/workflows/release.yml [\#27](https://github.com/instana/terraform-provider-instana/pull/27) ([parekh-raj](https://github.com/parekh-raj))
- Change release --rm-dist with --clean in release.yml [\#26](https://github.com/instana/terraform-provider-instana/pull/26) ([parekh-raj](https://github.com/parekh-raj))

## [v3.1.1](https://github.com/instana/terraform-provider-instana/tree/v3.1.1) (2025-02-10)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.1.0...v3.1.1)

**Merged pull requests:**

- Extend index.md with Infrastructure Alert Config [\#22](https://github.com/instana/terraform-provider-instana/pull/22) ([parekh-raj](https://github.com/parekh-raj))

## [v3.1.0](https://github.com/instana/terraform-provider-instana/tree/v3.1.0) (2025-01-16)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.0.1...v3.1.0)

**Implemented enhancements:**

- Add infra alert config resource [\#19](https://github.com/instana/terraform-provider-instana/pull/19) ([parekh-raj](https://github.com/parekh-raj))

## [v3.0.1](https://github.com/instana/terraform-provider-instana/tree/v3.0.1) (2024-12-20)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v3.0.0...v3.0.1)

**Implemented enhancements:**

- Add API implementation for Infra Alert Config [\#18](https://github.com/instana/terraform-provider-instana/pull/18) ([parekh-raj](https://github.com/parekh-raj))

**Merged pull requests:**

- Bug fix [\#16](https://github.com/instana/terraform-provider-instana/pull/16) ([rorywelch](https://github.com/rorywelch))
- Added support to send a User-Agent header with current Terraform Provider Version to Instana  [\#15](https://github.com/instana/terraform-provider-instana/pull/15) ([rorywelch](https://github.com/rorywelch))

## [v3.0.0](https://github.com/instana/terraform-provider-instana/tree/v3.0.0) (2024-05-30)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/v1.0.0...v3.0.0)

## [v1.0.0](https://github.com/instana/terraform-provider-instana/tree/v1.0.0) (2024-05-27)

[Full Changelog](https://github.com/instana/terraform-provider-instana/compare/627e6874cfda8cf8e5d5793f016aaf60b5285e6f...v1.0.0)

**Merged pull requests:**

- add Terraform Registry Manifest File [\#10](https://github.com/instana/terraform-provider-instana/pull/10) ([ChinmayGitHub](https://github.com/ChinmayGitHub))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
