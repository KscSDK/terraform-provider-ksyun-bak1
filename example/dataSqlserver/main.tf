provider "ksyun"{
  region = "cn-shanghai-3"
  access_key = "666"
  secret_key = "666"
}


variable "available_zone" {
  default = "cn-shanghai-2b"
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

resource "ksyun_sqlserver" "ks-ss-233"{
  output_file = "output_file"
  db_instance_class= "db.ram.2|db.disk.100"
  db_instance_name = "ksyun_sqlserver_1"
  db_instance_type = "HRDS_SS"
  engine = "SQLServer"
  engine_version = "2008r2"
  master_user_name = "admin"
  master_user_password = "123qweASD"
  vpc_id = "${ksyun_vpc.default.id}"
  subnet_id = "${ksyun_subnet.foo.id}"
  bill_type = "DAY"

}

data "ksyun_sqlservers" "search-sqlserver"{
  output_file = "output_file"
  db_instance_identifier = "${ksyun_sqlserver.ks-ss-233.id}"
}