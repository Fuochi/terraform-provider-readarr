---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "readarr_quality_definition Resource - terraform-provider-readarr"
subcategory: "Profiles"
description: |-
  Quality Definition resource.
  For more information refer to Quality Definition https://wiki.servarr.com/readarr/settings#quality-1 documentation.
---

# readarr_quality_definition (Resource)

<!-- subcategory:Profiles -->Quality Definition resource.
For more information refer to [Quality Definition](https://wiki.servarr.com/readarr/settings#quality-1) documentation.

## Example Usage

```terraform
resource "readarr_metadata_profile" "example" {
  title    = "PDF"
  id       = 2
  min_size = 35.0
  max_size = 400
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Quality Definition ID.
- `title` (String) Quality Definition Title.

### Optional

- `max_size` (Number) Maximum size MB/min.
- `min_size` (Number) Minimum size MB/min.

### Read-Only

- `quality_id` (Number) Quality ID.
- `quality_name` (String) Quality Name.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import readarr_quality_definition.example 10
```