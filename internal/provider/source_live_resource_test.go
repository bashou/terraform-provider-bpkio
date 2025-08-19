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

// 1. Basic creation with required fields.
func TestAccSourceLive_Basic(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resourceName := "bpkio_source_live.test"
	liveName := "tf-acc-test-live-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLiveConfig(apiKey, liveName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", liveName),
					resource.TestCheckResourceAttr(resourceName, "url", "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
		},
	})
}

func testAccSourceLiveConfig(apiKey, name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "test" {
  name = "%s"
  url  = "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"
}
`, apiKey, name)
}

// 2. Invalid URL (asset does not exist).
func TestAccSourceLive_InvalidURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-invalid-live-" + randomSuffix()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccSourceLiveInvalidURLConfig(apiKey, name),
				ExpectError: regexp.MustCompile(`(?i)400|not found|invalid|unreachable`),
			},
		},
	})
}

func testAccSourceLiveInvalidURLConfig(apiKey, name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "test" {
  name = "%s"
  url  = "https://this-does-not-exist.example.com/nonexistent.m3u8"
}
`, apiKey, name)
}

// 3. Missing required field (name).
func TestAccSourceLive_MissingName(t *testing.T) {
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
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "test" {
  url = "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"
}
`, apiKey),
				ExpectError: regexp.MustCompile(`(?s)The argument "name" is required, but no definition was found.`),
			},
		},
	})
}

// 4. Duplicate name+url (if API returns error).
func TestAccSourceLive_DuplicateNameURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	dupName := "tf-acc-dup-live-" + randomSuffix()
	config := testAccSourceLiveDuplicateConfig(apiKey, dupName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?i)403|500`),
			},
		},
	})
}

func testAccSourceLiveDuplicateConfig(apiKey, name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "first" {
  name = "%s"
  url  = "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"
}

resource "bpkio_source_live" "second" {
  name = "%s"
  url  = "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"
}
`, apiKey, name, name)
}

// 5. Check computed fields are always set.
func TestAccSourceLive_ComputedFields(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-computed-live-" + randomSuffix()
	resourceName := "bpkio_source_live.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLiveConfig(apiKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "format"),
				),
			},
		},
	})
}

// 6. Minimal config (omit optional description, multi_period, origin).
func TestAccSourceLive_MinimalConfig(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	name := "tf-acc-minimal-live-" + randomSuffix()
	resourceName := "bpkio_source_live.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLiveConfig(apiKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", ""), // should default to empty
					resource.TestCheckResourceAttr(resourceName, "multi_period", "false"),
				),
			},
		},
	})
}

// 7. Long names and special characters.
func TestAccSourceLive_LongSpecialName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resourceName := "bpkio_source_live.test"
	name := fmt.Sprintf("tf-acc-test-live-live-live-live-live-live-live-live-live-live-live-live-live-live-live-live--特殊字符-🚀-%s", randomSuffix())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "test" {
  name = "%s"
  url  = "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"
}
`, apiKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
				ExpectError: regexp.MustCompile(`(?s)Bad Request`),
			},
		},
	})
}
