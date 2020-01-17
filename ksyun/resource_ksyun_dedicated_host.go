package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func resourceKsyunDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunDedicatedHostCreate,
		Read:   resourceKsyunDedicatedHostRead,
		Update: resourceKsyunDedicatedHostUpdate,
		Delete: resourceKsyunDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"dedicated_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Monthly",
					"Daily",
				}, false),
			},
			"purchase_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				//ForceNew: true,
			},

			"dedicated_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
				//	ForceNew: true,
			},
			"total_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ori_cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_cpu": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_datadisk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_datadisk": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},

			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}
func resourceKsyunDedicatedHostCreate(d *schema.ResourceData, m interface{}) error {
	dedicatedconn := m.(*KsyunClient).dedicatedconn
	createDedicatedHost := make(map[string]interface{})
	creates := []string{
		"dedicated_type",
		"name",
		"charge_type",
		"purchase_time",
		"project_id",
		"dedicated_cluster_id",
		"availability_zone",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createDedicatedHost[vv] = fmt.Sprintf("%v", v1)
		}
	}
	createDedicatedHost["Number"] = "1"
	action := "CreateDedicatedHosts"
	logger.Debug(logger.ReqFormat, action, createDedicatedHost)
	resp, err := dedicatedconn.CreateDedicatedHosts(&createDedicatedHost)
	logger.Debug(logger.AllFormat, action, createDedicatedHost, *resp, err)
	if err != nil {
		return fmt.Errorf("createDedicatedHost Error  : %s", err)
	}
	set, ok := (*resp)["DedicatedHostSet"]
	if !ok {
		return fmt.Errorf("createDedicatedHost Error  : no hostset found")
	}
	sets, ok := set.([]interface{})
	if !ok || len(sets) == 0 {
		return fmt.Errorf("createDedicatedHost Error  : no hostset found")
	}
	dedicated, ok := sets[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("createDedicatedHost Error  : no host found")
	}
	id, ok := dedicated["DedicatedHostId"]
	if !ok {
		return fmt.Errorf("createDedicatedHost Error  : no id found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("createDedicatedHost Error : no id found")
	}
	d.SetId(idres)
	//dedicate 创建非同步，也不能查状态。
	time.Sleep(time.Second * 3)
	return resourceKsyunDedicatedHostRead(d, m)
}

func resourceKsyunDedicatedHostRead(d *schema.ResourceData, m interface{}) error {
	dedicatedconn := m.(*KsyunClient).dedicatedconn
	readDedicatedHost := make(map[string]interface{})
	readDedicatedHost["DedicatedHostId"] = d.Id()
	if pd, ok := d.GetOk("project_id"); ok {
		readDedicatedHost["project_id"] = fmt.Sprintf("%v", pd)
	}
	action := "DescribeDedicatedHosts"
	logger.Debug(logger.ReqFormat, action, readDedicatedHost)
	resp, err := dedicatedconn.DescribeDedicatedHosts(&readDedicatedHost)
	logger.Debug(logger.AllFormat, action, readDedicatedHost, *resp, err)
	if err != nil {
		return fmt.Errorf("Error  : %s", err)
	}
	itemset, ok := (*resp)["DedicatedHostSet"]
	if !ok {
		d.SetId("")
		return nil
	}
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, items[0], dedicatedHostKeys, map[string]bool{"DedicatedHostName": true, "Model": true, "ProjectId": true})
	itemMap, ok := items[0].(map[string]interface{})
	d.Set(Hump2Downline("ProjectId"), fmt.Sprintf("%v", itemMap["ProjectId"]))
	d.Set("name", itemMap["DedicatedHostName"])
	d.Set("dedicated_type", itemMap["Model"])
	return nil
}

func resourceKsyunDedicatedHostUpdate(d *schema.ResourceData, m interface{}) error {
	dedicatedconn := m.(*KsyunClient).dedicatedconn
	// Enable partial attribute modification
	d.Partial(true)
	// Whether the representative has any modifications
	nameUpdate := false
	cpuUpdate := false
	updateReq := make(map[string]interface{})
	// modify
	if d.HasChange("name") && !d.IsNewResource() {
		if v, ok := d.GetOk("name"); ok {
			updateReq["NewDedicatedHostName"] = fmt.Sprintf("%v", v)
		} else {
			return fmt.Errorf("cann't change name to empty string")
		}
		nameUpdate = true
	}
	if nameUpdate {
		updateReq["DedicatedHostId"] = d.Id()
		action := "RenameDedicatedHost"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := dedicatedconn.RenameDedicatedHost(&updateReq)
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		if err != nil {
			return fmt.Errorf("update name (%v)error:%v", updateReq, err)
		}
		d.SetPartial("name")
	}

	if d.HasChange("total_cpu") {
		if v, ok := d.GetOk("total_cpu"); ok {
			updateReq["VCPU"] = fmt.Sprintf("%v", v)
		} else {
			return fmt.Errorf("cann't change total_cpu to empty string")
		}
		cpuUpdate = true
	}
	if cpuUpdate {
		updateReq["DedicatedHostId.1"] = d.Id()
		action := "SetvCPU"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := dedicatedconn.SetvCPU(&updateReq)
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		if err != nil {
			return fmt.Errorf("update cpu (%v)error:%v", updateReq, err)
		}
		d.SetPartial("total_cpu")
	}
	d.Partial(false)
	return resourceKsyunDedicatedHostRead(d, m)
}

func resourceKsyunDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).dedicatedconn
	//delete
	deleteDedicatedHost := make(map[string]interface{})
	deleteDedicatedHost["DedicatedHostId.1"] = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteDedicatedHost"
		logger.Debug(logger.ReqFormat, action, deleteDedicatedHost)
		resp, err1 := conn.DeleteDedicatedHost(&deleteDedicatedHost)
		logger.Debug(logger.AllFormat, action, deleteDedicatedHost, *resp, err1)
		//dedicate 删除非同步
		if err1 == nil {
			time.Sleep(time.Second * 5)
		}
		if err1 != nil && notFoundError(err1) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}

		//check
		readDedicatedHost := make(map[string]interface{})
		readDedicatedHost["DedicatedHostId"] = d.Id()
		if pd, ok := d.GetOk("project_id"); ok {
			readDedicatedHost["project_id"] = fmt.Sprintf("%v", pd)
		}
		action = "DescribeDedicatedHosts"
		logger.Debug(logger.ReqFormat, action, readDedicatedHost)
		resp, err := conn.DescribeDedicatedHosts(&readDedicatedHost)
		logger.Debug(logger.AllFormat, action, readDedicatedHost, *resp, err)
		if err != nil && notFoundError(err) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading dedicated host when deleting %q, %s", d.Id(), err))
		}
		itemset, ok := (*resp)["DedicatedHostSet"]
		if !ok {
			return nil
		}
		item, ok := itemset.([]interface{})
		if !ok || len(item) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("error on  deleting dedicated host %v,%v", d.Id(), err1))
	})

}
