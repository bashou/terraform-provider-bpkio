---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bpkio_transcoding_profile Data Source - bpkio"
subcategory: ""
description: |-
  
---

# bpkio_transcoding_profile (Data Source)



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

data "bpkio_transcoding_profile" "this" {
  id = 4694
}

output "this_transcoding_profile" {
  value = data.bpkio_transcoding_profile.this
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `content` (String)
- `id` (Number) The ID of this resource.
- `internal_id` (String)
- `name` (String)
