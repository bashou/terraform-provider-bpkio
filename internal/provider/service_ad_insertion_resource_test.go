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

var (
	SlateURL     = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
	LiveURL      = "https://origin.broadpeak.io/bpk-tv/bpkiofficial/hlsv3/index.m3u8"
	LiveURLOther = "https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
	AdServerURL  = "https://bpkiovast.s3.eu-west-1.amazonaws.com/vastmultibpkio"
)

func TestAccServiceAdInsertion_Basic(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	slateName := "tf-acc-slate-" + randomSuffix()
	liveName := "tf-acc-live-" + randomSuffix()
	adServerName := "tf-acc-adserver-" + randomSuffix()
	adInsertionName := "tf-acc-adinsertion-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfig(
					apiKey, slateName, liveName, adServerName, adInsertionName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "id"),
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "source.id"),
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "live_ad_replacement.ad_server.id"),
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "live_ad_replacement.gap_filler.id"),
				),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	// Generate unique names for all resources for this test run
	slateName := "tf-acc-slate-" + randomSuffix()
	liveName := "tf-acc-live-" + randomSuffix()
	adServerName := "tf-acc-adserver-" + randomSuffix()
	initialServiceName := "tf-acc-service-ad-initial-" + randomSuffix()
	updatedServiceName := "tf-acc-service-ad-updated-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithName(
					apiKey, slateName, liveName, adServerName, initialServiceName,
				),
				Check: resource.TestCheckResourceAttr(
					"bpkio_service_ad_insertion.test", "name", initialServiceName,
				),
			},
			{
				Config: testAccServiceAdInsertionConfigWithName(
					apiKey, slateName, liveName, adServerName, updatedServiceName,
				),
				Check: resource.TestCheckResourceAttr(
					"bpkio_service_ad_insertion.test", "name", updatedServiceName,
				),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateSlate(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	// Generate unique names for every resource
	initialSlateName := "tf-acc-slate-initial-" + randomSuffix()
	updatedSlateName := "tf-acc-slate-updated-" + randomSuffix()
	liveName := "tf-acc-live-" + randomSuffix()
	adServerName := "tf-acc-adserver-" + randomSuffix()
	serviceName := "tf-acc-service-adinsertion-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithSlateName(
					apiKey, initialSlateName, liveName, adServerName, serviceName),
				Check: resource.TestCheckResourceAttr(
					"bpkio_source_slate.slate", "name", initialSlateName),
			},
			{
				Config: testAccServiceAdInsertionConfigWithSlateName(
					apiKey, updatedSlateName, liveName, adServerName, serviceName),
				Check: resource.TestCheckResourceAttr(
					"bpkio_source_slate.slate", "name", updatedSlateName),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateLiveName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	initialLiveName := "tf-acc-live-initial-" + randomSuffix()
	updatedLiveName := "tf-acc-live-updated-" + randomSuffix()
	slateName := "tf-acc-slate-" + randomSuffix()
	adServerName := "tf-acc-adserver-" + randomSuffix()
	serviceName := "tf-acc-service-adinsertion-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithLiveName(
					apiKey, slateName, initialLiveName, adServerName, serviceName),
				Check: resource.TestCheckResourceAttr("bpkio_source_live.live", "name", initialLiveName),
			},
			{
				Config: testAccServiceAdInsertionConfigWithLiveName(
					apiKey, slateName, updatedLiveName, adServerName, serviceName),
				Check: resource.TestCheckResourceAttr("bpkio_source_live.live", "name", updatedLiveName),
			},
		},
	})
}

func TestAccServiceAdInsertion_ImportStateAndDrift(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	// Generate unique names to avoid collisions across jobs
	slateName := "tf-acc-slate-" + randomSuffix()
	liveName := "tf-acc-live-" + randomSuffix()
	adServerName := "tf-acc-adserver-" + randomSuffix()
	serviceName := "tf-acc-adinsertion-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfig(apiKey, slateName, liveName, adServerName, serviceName),
			},
			{
				ResourceName:      "bpkio_service_ad_insertion.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// --- Error cases: try creating with invalid source/adserver id (should error out) ---.

func TestAccServiceAdInsertion_InvalidSource(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	badID := 999999999

	slateName := "tf-acc-slate-" + randomSuffix()
	adServerName := "tf-acc-adserver-" + randomSuffix()
	serviceName := "tf-acc-adinsertion-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithBadSource(
					apiKey, slateName, adServerName, serviceName, badID,
				),
				ExpectError: regexp.MustCompile(`(?i)403|forbidden|not allowed`),
			},
		},
	})
}

func TestAccServiceAdInsertion_InvalidAdServer(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	badID := 999999999

	liveName := "tf-acc-live-" + randomSuffix()
	slateName := "tf-acc-slate-" + randomSuffix()
	serviceName := "tf-acc-adinsertion-badadserver-" + randomSuffix()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithBadAdServer(
					apiKey, liveName, slateName, serviceName, badID,
				),
				ExpectError: regexp.MustCompile(`(?i)403|forbidden|not allowed`),
			},
		},
	})
}

// --- Config helpers ---.

func testAccServiceAdInsertionConfig(
	apiKey, slateName, liveName, adServerName, adInsertionName string,
) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "%s"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5963
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, slateName, SlateURL, liveName, LiveURL, adServerName, AdServerURL, adInsertionName)
}

// Change live source name.
func testAccServiceAdInsertionConfigWithName(
	apiKey, slateName, liveName, adServerName, serviceName string,
) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "%s"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5963
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, slateName, SlateURL, liveName, LiveURL, adServerName, AdServerURL, serviceName)
}

// Change live source name.
func testAccServiceAdInsertionConfigWithLiveName(
	apiKey, slateName, liveName, adServerName, serviceName string,
) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "%s"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5963
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, slateName, SlateURL, liveName, LiveURL, adServerName, AdServerURL, serviceName)
}

// Change slate name.
func testAccServiceAdInsertionConfigWithSlateName(
	apiKey, slateName, liveName, adServerName, serviceName string,
) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "%s"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
  id = 5963
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, slateName, SlateURL, liveName, LiveURL, adServerName, AdServerURL, serviceName)
}

// Invalid source (bad id).
func testAccServiceAdInsertionConfigWithBadSource(
	apiKey, slateName, adServerName, serviceName string, badSourceID int,
) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

data "bpkio_transcoding_profile" "test" {
  id = 5963
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_service_ad_insertion" "adservice" {
  name = "%s"

  source = {
    id = %d
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, slateName, SlateURL, adServerName, AdServerURL, serviceName, badSourceID)
}

// Invalid ad server (bad id).
func testAccServiceAdInsertionConfigWithBadAdServer(
	apiKey, liveName, slateName, serviceName string, badAdServerID int,
) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "live" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
  id = 5963
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = %d
    }

    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, liveName, LiveURL, slateName, SlateURL, serviceName, badAdServerID)
}
