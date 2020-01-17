# Specify the provider and access details
provider "ksyun" {
  region = "cn-beijing-6"
}

resource "ksyun_instance" "default" {
  image_id="6e37ed46-61a2-4f0a-9f4a-dcdd0817917c"
  instance_type="N3.2B"
  system_disk{
    disk_type="SSD3.0"
    disk_size=30
  }
  data_disk_gb=0
  #only support part type
  data_disk =[
   {
      type="SSD3.0"
      size=20
      delete_with_instance=true
   }
 ]
  subnet_id="2ea4195d-8111-4cd8-91da-6a34bb06663b"
  instance_password="Xuan663222"
  keep_image_login=false
  charge_type="Daily"
  purchase_time=1
  security_group_id=["6e3dee9c-291c-4647-bfc2-4c1eaa93fb80"]
  private_ip_address=""
  instance_name="xuan-tf"
  instance_name_suffix=""
  sriov_net_support=false
  project_id=0
  data_guard_id=""
  key_id=[]
  d_n_s1 =""
  d_n_s2 =""
  force_delete =true
  user_data=""
  host_name="xuan-tf"
}



#dedicatedInstance
/*
resource "ksyun_instance" "default" {
  image_id="6e37ed46-61a2-4f0a-9f4a-dcdd0817917c"
  dedicated_host_id="bbdae3bf-d2a9-4c0e-bcb7-2263343ec810"
  instance_configure{
   memory_gb=20
   v_c_p_u=20
  }
  instance_type="DVM1.NONE"
  data_disk_gb=20
  subnet_id="f6033a90-8c87-4a72-8cca-a7f61bf38c2b"
  instance_password="Xuan663222"
  charge_type="Daily"
  purchase_time=1
  security_group_id=["e10287dd-6702-4c45-b14e-3497f38b8f58"]
  private_ip_address=""
  instance_name="xuan-tf"
  project_id=0
  force_delete=true
  data_disk =[
   {
     type=""
     size=0
     delete_with_instance=false
   }
 ]
}
*/

