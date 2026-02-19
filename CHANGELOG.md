# Changelog

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
