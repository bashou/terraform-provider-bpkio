// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	adServerURL = "https://vast-prep.staging.olyzon.tv/sources/1042586b/serve"
)

func TestAccSourceAdServer_Basic(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	name := "tf-acc-test-adserver-" + randomSuffix()
	resourceName := "bpkio_source_adserver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceAdServerConfig(apiKey, adServerURL, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "url", adServerURL),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			{
				// Test update: Change description and check result
				Config: testAccSourceAdServerConfigUpdate(apiKey, adServerURL, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acceptance test"),
				),
			},
		},
	})
}

func testAccSourceAdServerConfig(apiKey, url, name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  name = "%s"
  url  = "%s"
}
`, apiKey, name, url)
}

func testAccSourceAdServerConfigUpdate(apiKey, url, name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  name        = "%s"
  url         = "%s"
  description = "updated by acceptance test"
}
`, apiKey, name, url)
}

func TestAccSourceAdServer_ComputedFields(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-test-adserver-" + randomSuffix()
	resourceName := "bpkio_source_adserver.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceAdServerConfig(apiKey, adServerURL, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
		},
	})
}

func TestAccSourceAdServer_MinimalConfig(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-test-adserver-" + randomSuffix()
	resourceName := "bpkio_source_adserver.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceAdServerConfig(apiKey, adServerURL, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "url", adServerURL),
				),
			},
		},
	})
}

func TestAccSourceAdServer_MissingName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  url = "%s"
}
`, apiKey, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)The argument "name" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccSourceAdServer_LongSpecialName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	longName := strings.Repeat("x", 101) // >100 chars to trigger validation error
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  name = "%s"
  url  = "%s"
}
`, apiKey, longName, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)Bad Request`),
			},
		},
	})
}

func TestAccSourceAdServer_DuplicateNameURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	unique := "tf-acc-dupe-" + randomSuffix()
	dupeConfig := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}
resource "bpkio_source_adserver" "first" {
  name = "%s"
  url  = "%s"
}
resource "bpkio_source_adserver" "second" {
  name = "%s"
  url  = "%s"
}
`, apiKey, unique, adServerURL, unique, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      dupeConfig,
				ExpectError: regexp.MustCompile(`(?s)403|500`),
			},
		},
	})
}

func TestAccSourceAdServer_QueryParameters(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-test-adserver-" + randomSuffix()
	resourceName := "bpkio_source_adserver.test"
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}
resource "bpkio_source_adserver" "test" {
  name = "%s"
  url  = "%s"
  query_parameters = [
    {
      type  = "from-header"
      name  = "X-Test"
      value = "value1"
    },
    {
      type  = "custom"
      name  = "X-Custom"
      value = "value2"
    }
  ]
}
`, apiKey, name, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query_parameters.#", "2"),
				),
			},
		},
	})
}

func TestAccSourceAdServer_InvalidURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-test-adserver-" + randomSuffix()
	badURL := "https://this-url-will-not-exist.example.com"
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}
resource "bpkio_source_adserver" "test" {
  name = "%s"
  url  = "%s"
}
`, apiKey, name, badURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
		},
	})
}
