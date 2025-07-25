---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bpkio_source_live Resource - bpkio"
subcategory: ""
description: |-
  
---

# bpkio_source_live (Resource)



## Example Usage

```terraform
terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
}

resource "bpkio_source_live" "this" {
  name        = "foobar-test-tf"
  description = "test"
  url         = "https://live.stream/master.m3u8"

  //TODO: Find way to handle when origin is empty
  origin = {}
}

resource "bpkio_source_slate" "this" {
  name        = "foobar-test-tf"
  description = "test slate"
  url         = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the source live.
- `url` (String) The URL of the source live.

### Optional

- `description` (String) The description of the source live.
- `multi_period` (Boolean) Whether the source live supports multiple periods.(Default: `false`)
- `origin` (Attributes) The origin configuration for the source live. (see [below for nested schema](#nestedatt--origin))

### Read-Only

- `format` (String) The format of the source live.
- `id` (Number) The ID of the source live.
- `type` (String) The type of the source live.

<a id="nestedatt--origin"></a>
### Nested Schema for `origin`

Optional:

- `custom_headers` (Attributes List) (see [below for nested schema](#nestedatt--origin--custom_headers))

<a id="nestedatt--origin--custom_headers"></a>
### Nested Schema for `origin.custom_headers`

Required:

- `name` (String) The name of the custom header.
- `value` (String) The value of the custom header.

## Import

Import is supported using the following syntax:

```shell
# Live Source can be imported by specifying the numeric identifier.
terraform import bpkio_source_live.example 123
```
