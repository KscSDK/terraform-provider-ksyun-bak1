# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Create an dedicated host
resource "ksyun_dedicated_host" "default1" {
  dedicated_type ="DC1"
  name = "xuan"
  charge_type="Daily"
  purchase_time =1
  project_id=0
  dedicated_cluster_id=""
  availability_zone="cn-beijing-6a"
  total_cpu=60
}
