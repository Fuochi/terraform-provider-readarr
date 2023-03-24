---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "readarr_notification_discord Resource - terraform-provider-readarr"
subcategory: "Notifications"
description: |-
  Notification Discord resource.
  For more information refer to Notification https://wiki.servarr.com/readarr/settings#connect and Discord https://wiki.servarr.com/readarr/supported#discord.
---

# readarr_notification_discord (Resource)

<!-- subcategory:Notifications -->Notification Discord resource.
For more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [Discord](https://wiki.servarr.com/readarr/supported#discord).

## Example Usage

```terraform
resource "readarr_notification_discord" "example" {
  on_grab                         = false
  on_download_failure             = false
  on_upgrade                      = false
  on_rename                       = false
  on_import_failure               = false
  on_book_delete                  = false
  on_book_file_delete             = true
  on_book_file_delete_for_upgrade = false
  on_health_issue                 = false
  on_book_retag                   = true
  on_author_delete                = false
  on_release_import               = false

  include_health_warnings = false
  name                    = "Example"

  web_hook_url  = "http://discord-web-hook.com"
  username      = "User"
  avatar        = "https://i.imgur.com/oBPXx0D.png"
  grab_fields   = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
  import_fields = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `include_health_warnings` (Boolean) Include health warnings.
- `name` (String) Notification name.
- `on_author_delete` (Boolean) On author deleted flag.
- `on_book_delete` (Boolean) On book delete flag.
- `on_book_file_delete` (Boolean) On book file delete flag.
- `on_book_file_delete_for_upgrade` (Boolean) On book file delete for upgrade flag.
- `on_book_retag` (Boolean) On book retag flag.
- `on_download_failure` (Boolean) On download failure flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_import_failure` (Boolean) On import failure flag.
- `on_release_import` (Boolean) On release import flag.
- `on_rename` (Boolean) On rename flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `web_hook_url` (String) Web hook URL.

### Optional

- `avatar` (String) Avatar.
- `tags` (Set of Number) List of associated tags.
- `username` (String) Username.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import readarr_notification_discord.example 1
```