---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "readarr_download_client_sabnzbd Resource - terraform-provider-readarr"
subcategory: "Download Clients"
description: |-
  Download Client Sabnzbd resource.
  For more information refer to Download Client https://wiki.servarr.com/readarr/settings#download-clients and Sabnzbd https://wiki.servarr.com/readarr/supported#sabnzbd.
---

# readarr_download_client_sabnzbd (Resource)

<!-- subcategory:Download Clients -->Download Client Sabnzbd resource.
For more information refer to [Download Client](https://wiki.servarr.com/readarr/settings#download-clients) and [Sabnzbd](https://wiki.servarr.com/readarr/supported#sabnzbd).

## Example Usage

```terraform
resource "readarr_download_client_sabnzbd" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "sabnzbd"
  url_base = "/sabnzbd/"
  port     = 9091
  api_key  = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `api_key` (String, Sensitive) API key.
- `book_category` (String) Book category.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `older_book_priority` (Number) Older Music priority. `-100` Default, `-2` Paused, `-1` Low, `0` Normal, `1` High, `2` Force.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `recent_book_priority` (Number) Recent Music priority. `-100` Default, `-2` Paused, `-1` Low, `0` Normal, `1` High, `2` Force.
- `tags` (Set of Number) List of associated tags.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import readarr_download_client_sabnzbd.example 1
```