// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceSlate_Basic(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resourceName := "bpkio_source_slate.test"
	name := "tf-acc-test-slate-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceSlateConfig(apiKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "url", "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
		},
	})
}

func testAccSourceSlateConfig(apiKey, name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "test" {
  name = "%s"
  url  = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
}
`, apiKey, name)
}

func TestAccSourceSlate_InvalidURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	badURL := "https://this-url-does-not-exist.broadpeak.io/foo.jpg"
	name := "tf-acc-test-invalid-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "test" {
  name = "%s"
  url  = "%s"
}
`, apiKey, name, badURL),
				ExpectError: regexp.MustCompile(`(?i)400|not found|invalid|unreachable`),
			},
		},
	})
}

func TestAccSourceSlate_Update(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	resourceName := "bpkio_source_slate.test"
	initialName := "tf-acc-test-slate-update-" + randomSuffix()
	updatedName := "tf-acc-test-slate-updated-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "test" {
  name = "%s"
  url  = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
  description = "first description"
}
`, apiKey, initialName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", initialName),
					resource.TestCheckResourceAttr(resourceName, "description", "first description"),
				),
			},
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "test" {
  name = "%s"
  url  = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
  description = "updated description"
}
`, apiKey, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
				),
			},
		},
	})
}

func TestAccSourceSlate_Import(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	resourceName := "bpkio_source_slate.test"
	name := "tf-acc-test-slate-import-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "test" {
  name = "%s"
  url  = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
}
`, apiKey, name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSourceSlate_MissingName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "test" {
  url = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
}
`, apiKey),
				ExpectError: regexp.MustCompile(`(?i)The argument\s+"name"\s+is required`),
			},
		},
	})
}

func TestAccSourceSlate_MissingURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-test-missing-url-" + randomSuffix()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "test" {
  name = "%s"
}
`, apiKey, name),
				ExpectError: regexp.MustCompile(`(?i)The argument\s+"url"\s+is required`),
			},
		},
	})
}

func TestAccSourceSlate_DuplicateNameURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-test-duplicate-" + randomSuffix()
	url := "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "first" {
  name = "%s"
  url  = "%s"
}
resource "bpkio_source_slate" "second" {
  name = "%s"
  url  = "%s"
}
`, apiKey, name, url, name, url),
				// Should fail to create the second one
				ExpectError: regexp.MustCompile(`(?s)(500|403)`),
			},
		},
	})
}

func TestAccSourceSlate_ComputedFields(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resourceName := "bpkio_source_slate.test"
	name := "tf-acc-test-computed-" + randomSuffix()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceSlateConfig(apiKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "format"),
				),
			},
		},
	})
}

func TestAccSourceSlate_MinimalConfig(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resourceName := "bpkio_source_slate.minimal"
	name := "tf-acc-test-minimal-" + randomSuffix()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "minimal" {
  name = "%s"
  url  = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
}
`, apiKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "format"),
				),
			},
		},
	})
}

func TestAccSourceSlate_LongNameAndSpecialChars(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	longName := fmt.Sprintf("tf-acc-test-超级长的名字-🚀-%s", randomSuffix())
	resourceName := "bpkio_source_slate.special"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" { api_key = "%s" }
resource "bpkio_source_slate" "special" {
  name = "%s"
  url  = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
}
`, apiKey, longName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", longName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}
