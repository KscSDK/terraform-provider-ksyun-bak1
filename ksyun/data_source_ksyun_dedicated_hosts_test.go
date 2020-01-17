package ksyun

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccKsyunDedicatedHostsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDedicatedHostsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ksyun_dedicated_hosts.foo"),
				),
			},
		},
	})
}

const testAccDataDedicatedHostsConfig = `

data "ksyun_dedicated_hosts" "foo" {
  output_file="output_result"
  id=""
  project_id=["0"]
}
`
