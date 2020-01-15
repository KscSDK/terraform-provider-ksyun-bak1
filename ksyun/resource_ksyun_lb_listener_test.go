package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccKsyunListener_basic(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_listener.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerExists("ksyun_lb_listener.foo", &val),
					testAccCheckListenerAttributes(&val),
				),
			},
		},
	})
}

func TestAccKsyunListener_update(t *testing.T) {
	var val map[string]interface{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "ksyun_lb_listener.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerExists("ksyun_lb_listener.foo", &val),
					testAccCheckListenerAttributes(&val),
				),
			},
			{
				Config: testAccListenerUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerExists("ksyun_lb_listener.foo", &val),
					testAccCheckListenerAttributes(&val),
				),
			},
		},
	})
}

func testAccCheckListenerExists(n string, val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Listener id is empty")
		}
		client := testAccProvider.Meta().(*KsyunClient)
		listener := make(map[string]interface{})
		listener["ListenerId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeListeners(&listener)
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ListenerSet"].([]interface{})
			if len(l) == 0 {
				return err
			}
		}

		*val = *ptr
		return nil
	}
}
func testAccCheckListenerAttributes(val *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val != nil {
			l := (*val)["ListenerSet"].([]interface{})
			if len(l) == 0 {
				return fmt.Errorf("Listener id is empty")
			}
		}
		return nil
	}
}
func testAccCheckListenerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_lb_listener" {
			continue
		}

		client := testAccProvider.Meta().(*KsyunClient)
		listener := make(map[string]interface{})
		listener["ListenerId.1"] = rs.Primary.ID
		ptr, err := client.slbconn.DescribeListeners(&listener)

		// Verify the error is what we want
		if err != nil {
			return err
		}
		if ptr != nil {
			l := (*ptr)["ListenerSet"].([]interface{})
			if len(l) == 0 {
				continue
			} else {
				return fmt.Errorf("Listener still exist")
			}
		}
	}

	return nil
}

const testAccListenerConfig = `
resource "ksyun_vpc" "default" {
	vpc_name        = "tf-acc-vpc"
	cidr_block = "192.168.0.0/16"
}
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "tf-xun-2"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "foo" {
  listener_name = "tf-xun-2"
  listener_port = "8080"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "stop"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}`

const testAccListenerUpdateConfig = `
resource "ksyun_vpc" "default" {
	vpc_name        = "tf-acc-vpc"
	cidr_block = "192.168.0.0/16"
}
resource "ksyun_lb" "default" {
  vpc_id = "${ksyun_vpc.default.id}"
  load_balancer_name = "tf-xun-2"
  type = "public"
  subnet_id = ""
  load_balancer_state = "stop"
  private_ip_address = ""
}
# Create Load Balancer Listener with tcp protocol
resource "ksyun_lb_listener" "foo" {
  listener_name = "tf-xun-update"
  listener_port = "8080"
  listener_protocol = "HTTP"
  listener_state = "stop"
  load_balancer_id = "${ksyun_lb.default.id}"
  method = "RoundRobin"
  certificate_id = ""
  session {
    session_state = "stop"
    session_persistence_period = 100
    cookie_type = "RewriteCookie"
    cookie_name = "cookiexunqq"
  }
}
`
