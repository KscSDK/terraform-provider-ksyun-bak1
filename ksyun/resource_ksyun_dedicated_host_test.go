package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"testing"
)

func TestAccKsyunDedicatedHost_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_dedicated_host.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedHostConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostExists("ksyun_dedicated_host.foo", &val),
					testAccCheckDedicatedHostAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunDedicatedHost_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_dedicated_host.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDedicatedHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedHostConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostExists("ksyun_dedicated_host.foo", &val),
					testAccCheckDedicatedHostAttributes(&val),
				),
			},
			{
				Config: testAccDedicatedHostUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedHostExists("ksyun_dedicated_host.foo", &val),
					testAccCheckDedicatedHostAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckDedicatedHostExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DedicatedHost id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		dedicatedHost := make(map[string]interface{})
		dedicatedHost["DedicatedHostId"] = rs.Primary.ID
		ptr, err := client.dedicatedconn.DescribeDedicatedHosts(&dedicatedHost)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["DedicatedHostSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}
		*val = *ptr
		return nil
	}
}
func testAccCheckDedicatedHostAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["DedicatedHostSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("DedicatedHost id is empty")
			}
		}
		return nil
	}
}
func testAccCheckDedicatedHostDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_dedicated_host" {
			continue
		}
		client := testAccProvider.Meta().(*KsyunClient)
		dedicatedHost := make(map[string]interface{})
		dedicatedHost["DedicatedHostId"] = rs.Primary.ID
		action := "DescribeDedicatedHosts"
		logger.Debug(logger.ReqFormat, action, dedicatedHost)
		ptr, err := client.dedicatedconn.DescribeDedicatedHosts(&dedicatedHost)
		logger.Debug(logger.AllFormat, action, *ptr, err)
		// Verify the error is what we want
		if err != nil && notFoundError(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["DedicatedHostSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("DedicatedHost still exist")
			}
		}
	}

	return nil
}

const testAccDedicatedHostConfig = `
# Create an dedicatedHost
resource "ksyun_dedicated_host" "foo" {
  dedicated_type ="DC1"
  name = "xuan"
  charge_type="Daily"
  purchase_time =1
  project_id=0
  dedicated_cluster_id=""
  availability_zone="cn-beijing-6a"
}
`

const testAccDedicatedHostUpdateConfig = `

# Create an dedicatedHost
resource "ksyun_dedicated_host" "foo" {
  dedicated_type ="DC1"
  name = "xuan-update"
  charge_type="Daily"
  purchase_time =1
  project_id=0
  dedicated_cluster_id=""
  availability_zone="cn-beijing-6a"
  total_cpu=90
}
`
