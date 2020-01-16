package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunEpc_basic(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_epc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEpcDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccEpcConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEpcExists("ksyun_epc.foo", &val),
					testAccCheckEpcAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunEpc_update(t *testing.T) {
	var val map[string]interface{}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ksyun_epc.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEpcDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccEpcConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEpcExists("ksyun_epc.foo", &val),
					testAccCheckEpcAttributes(&val),
				),
			},
			{
				Config: testAccEpcConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckEpcExists("ksyun_epc.foo", &val),
					testAccCheckEpcAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckEpcExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("host id is empty")
		}

		client := testAccProvider.Meta().(*KsyunClient)
		epc := make(map[string]interface{})
		epc["HostId.1"] = rs.Primary.ID
		ptr, err := client.epcconn.DescribeEpcs(&epc)

		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["HostSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}

func testAccCheckEpcAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["HostSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("host id is empty")
			}
		}
		return nil
	}
}

func testAccCheckEpcDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_epc" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		epc := make(map[string]interface{})
		epc["HostId.1"] = rs.Primary.ID
		ptr, err := client.epcconn.DescribeEpcs(&epc)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["HostSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Host still exist")
			}
		}
	}
	return nil
}

const testAccEpcConfig = `
data "ksyun_ssh_keys" "default" {
 output_file=""
 ids=[]
}
resource "ksyun_vpc" "default" {
 vpc_name   = "ksyun-vpc-tf"
 cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
 subnet_name      = "ksyun-subnet-tf"
 cidr_block = "10.7.0.0/21"
 subnet_type = "Physical"
 dhcp_ip_from = "10.7.0.2"
 dhcp_ip_to = "10.7.0.253"
 vpc_id  = "${ksyun_vpc.default.id}"
 gateway_ip = "10.7.0.1"
 dns1 = "198.18.254.41"
 dns2 = "198.18.254.40"
 availability_zone = "cn-beijing-6b"
}
resource "ksyun_security_group" "default" {
 vpc_id = "${ksyun_vpc.default.id}"
 security_group_name="ksyun-security-group"
}
resource "ksyun_epc" "foo" {
	host_name= "tf-acc-epc"
  	host_type = "CAL"
  	image_id = "2c9d8f29-6eb9-4bc7-90e5-b0bd7a9e2d3a"
  	key_id = "${data.ksyun_ssh_keys.default.keys.0.key_id}"
  	network_interface_mode = "bond4"
  	raid = "Raid5"
  	availability_zone = "cn-beijing-6b"
  	charge_type = "PostPaidByDay"
  	subnet_id = "${ksyun_subnet.default.id}"
  	security_group_id = ["${ksyun_security_group.default.id}"]
  	password = "Abc@1234"
}
`

const testAccEpcConfigUpdate = `
data "ksyun_ssh_keys" "default" {
 output_file=""
 ids=[]
}
resource "ksyun_vpc" "default" {
 vpc_name   = "ksyun-vpc-tf"
 cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "default" {
 subnet_name      = "ksyun-subnet-tf"
 cidr_block = "10.7.0.0/21"
 subnet_type = "Physical"
 dhcp_ip_from = "10.7.0.2"
 dhcp_ip_to = "10.7.0.253"
 vpc_id  = "${ksyun_vpc.default.id}"
 gateway_ip = "10.7.0.1"
 dns1 = "198.18.254.41"
 dns2 = "198.18.254.40"
 availability_zone = "cn-beijing-6b"
}
resource "ksyun_security_group" "default" {
 vpc_id = "${ksyun_vpc.default.id}"
 security_group_name="ksyun-security-group"
}
resource "ksyun_epc" "foo" {
    host_name="tf-acc-epc-1"
  	host_type = "CAL"
  	image_id = "2c9d8f29-6eb9-4bc7-90e5-b0bd7a9e2d3a"
  	key_id = "${data.ksyun_ssh_keys.default.keys.0.key_id}"
  	network_interface_mode = "bond4"
  	raid = "Raid5"
  	availability_zone = "cn-beijing-6b"
  	charge_type = "PostPaidByDay"
  	subnet_id = "${ksyun_subnet.default.id}"
  	security_group_id = ["${ksyun_security_group.default.id}"]
  	password = "Abc@1234"
}
`
