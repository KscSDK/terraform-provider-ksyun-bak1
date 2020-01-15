package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunKrdsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKrdsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_krds.search-krds"),
				),
			},
		},
	})
}

const testAccDataKrdsConfig = `


variable "available_zone" {
  default = "cn-guangzhou-1a"
}
resource "ksyun_vpc" "default" {
  vpc_name   = "ksyun-vpc-tf"
  cidr_block = "10.7.0.0/21"
}
resource "ksyun_subnet" "foo" {
  subnet_name      = "ksyun-subnet-tf"
  cidr_block = "10.7.0.0/21"
  subnet_type = "Reserve"
  dhcp_ip_from = "10.7.0.2"
  dhcp_ip_to = "10.7.0.253"
  vpc_id  = "${ksyun_vpc.default.id}"
  gateway_ip = "10.7.0.1"
  dns1 = "198.18.254.41"
  dns2 = "198.18.254.40"
  availability_zone = "${var.available_zone}"
}

resource "ksyun_krds" "rds_terraform_3"{
  output_file = "output_file"
  db_instance_class= "db.ram.2|db.disk.50"
  db_instance_name = "houbin_terraform_1-n"
  db_instance_type = "HRDS"
  engine = "mysql"
  engine_version = "5.5"
  master_user_name = "admin"
  master_user_password = "123qweASD123"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"
  preferred_backup_time = "01:00-02:00"
  parameters {
    name = "auto_increment_increment"
    value = "8"
  }

  parameters {
    name = "binlog_format"
    value = "ROW"
  }

}

data "ksyun_krds" "search-krds"{
  output_file = "output_file"
  db_instance_identifier = "${ksyun_krds.rds_terraform_3.id}"
}
`
