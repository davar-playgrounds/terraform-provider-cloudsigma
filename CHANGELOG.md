
<a name="v1.2.1"></a>
## v1.2.1 (2021-02-18)

### Bug Fixes
* **resource/cloudsigma_server:** revert validate function with diagnostics for ssh_keys attribute

### Documentation
* **resource/cloudsigma_server:** add an example assigning private and public networks


<a name="v1.2.0"></a>
## v1.2.0 (2021-02-17)

### Documentation
* **data-source/cloudsigma_vlan:** fix hcl example
* **resource/cloudsigma_server:** document network field with example

### Features
* **resource/cloudsigma_server:** add vlan_uuid field
* **resource/cloudsigma_server:** add network field

### Maintaining
* update terraform-plugin-sdk to v2.4.3
* **github-actions:** upgrade Go version to 1.15


<a name="v1.1.0"></a>
## v1.1.0 (2020-12-24)

### Features
* **resource/cloudsigma_drive:** replace legacy SchemaValidateFunc with wrapper function
* **resource/cloudsigma_server:** replace legacy SchemaValidateFunc with wrapper function

### Maintaining
* add support for openbsd systems


<a name="v1.0.1"></a>
## v1.0.1 (2020-11-25)

### Bug Fixes
* **resource/cloudsigma_server:** check if server runtime before IP address assignment

### Maintaining
* update terraform-plugin-sdk to v2.3.0


<a name="v1.0.0"></a>
## v1.0.0 (2020-11-09)

### Documentation
* add examples with additional drive for server resource

### Features
* **resource/cloudsigma_drive:** add uuid field
* **resource/cloudsigma_server:** add ipv4_address, ssh_keys fields
* **resource/cloudsigma_server:** add mounted_on field
* **resource/cloudsigma_ssh_key:** add private_key field

### Maintaining
* upgrade terraform-plugin-sdk to v2.2.0


<a name="v0.3.0"></a>
## v0.3.0 (2020-10-08)

### Documentation
* add authentication section

### Features
* **resource/cloudsigma_ssh_key:** add fingerprint field


<a name="v0.2.0"></a>
## v0.2.0 (2020-09-27)

### Features
* **resource/cloudsigma_drive_attachment:** remove drive_attachment resource
* **resource/cloudsigma_server:** add 'drive' option to attach server drives

### Maintaining
* upload draft changelog when creating github release
* add breaking changes notes to unreleased changelog
* remove dependency between lint and test tasks
* upgrade terraform-plugin-sdk to v2.0.3

### BREAKING CHANGE

resource drive_attachment is replaced with internal map in server resource


<a name="v0.1.2"></a>
## v0.1.2 (2020-09-07)

### Bug Fixes
* add missing  prefix by binary file name

### Documentation
* add getting-started guide to documentation


<a name="v0.1.1"></a>
## v0.1.1 (2020-08-28)

### Features
* **resource/cloudsigma_remote_snapshot:** add new resource for remote snapshots

### Maintaining
* run GoReleaser with GitHub actions
* enable GitHub actions
* automate release process with GoReleaser


<a name="v0.1.0"></a>
## v0.1.0 (2020-08-23)

### Maintaining
* generate changelog with git-chglog
* disable non-compile acc tests
* support caching for golanci-lint binary tool

