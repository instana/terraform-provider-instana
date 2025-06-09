# Terraform Provider Instana

> **Note:** This **terraform-provider-instana** project was originally maintained under the [gessnerfl](https://github.com/gessnerfl) organization and is now officially maintained by IBM under the [Instana](https://github.com/instana) organization since version [v3.0.0](https://github.com/instana/terraform-provider-instana/releases/tag/v3.0.0).

---

## Overview

This Terraform provider offers comprehensive support for managing Instana resources via the Instana REST API.

- **Terraform Registry:** [instana/instana](https://registry.terraform.io/providers/instana/instana/latest)
- **Changelog:** [CHANGELOG.md](https://github.com/instana/terraform-provider-instana/blob/master/CHANGELOG.md)

---

## Documentation

Full documentation is available on the Terraform Registry page:  
<https://registry.terraform.io/providers/instana/instana/latest>

---

## Implementation Details

### Testing

- **Mocking:**  
  Tests are colocated with the implementation packages.  
  We use [gomock](https://github.com/golang/mock) for mocking interfaces. Mocks are generated using *source mode* and placed in the `mock` package.

- **Generating Mocks:**  
  You can generate mocks using the helper script from the root directory:
  ```bash
  generate-mock-for-file <source-file>
  ```
  Alternatively you can manually execute `mockgen` as follows
  ```bash
  mockgen -source=<source_file> -destination=mocks/<source_file_name>_mocks.go -package=mocks
  ```

---

### Releasing a New Version

1. Follow [Semantic Versioning](https://semver.org/) to create a new tag.
2. Update the changelog before releasing using [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator).
3. Push the tag to the remote repository to trigger the release process.

