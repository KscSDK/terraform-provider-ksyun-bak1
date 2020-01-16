# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

# Get  eips
data "ksyun_dedicated_hosts" "default" {
  output_file="output_result"
  id=""
  project_id=["0"]
}

