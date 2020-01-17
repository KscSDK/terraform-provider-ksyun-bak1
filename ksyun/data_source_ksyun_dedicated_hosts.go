package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceKsyunDedicatedHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunDedicatedHostsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"project_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dedicated_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dedicated_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"dedicated_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_cpu": {
							Type:     schema.TypeInt,
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
						},

						"create_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunDedicatedHostsRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*KsyunClient).dedicatedconn
	var dedicated_hosts []string
	req := make(map[string]interface{})
	if id, ok := d.GetOk("id"); ok {
		req[fmt.Sprintf("DedicatedHostId")] = id
	}
	var projectIds []string
	if ids, ok := d.GetOk("project_id"); ok {
		projectIds = SchemaSetToStringSlice(ids)
	}
	for k, v := range projectIds {
		if v == "" {
			continue
		}
		req[fmt.Sprintf("ProjectId.%d", k+1)] = v
	}

	var allDedicatedHosts []interface{}
	resp, err := conn.DescribeDedicatedHosts(&req)
	if err != nil {
		return fmt.Errorf("error on reading  dedicated host list req(%v):%v", req, err)
	}
	itemSet, ok := (*resp)["DedicatedHostSet"]
	if !ok {
		return fmt.Errorf("error on reading dedicated host set")
	}
	allDedicatedHosts, ok = itemSet.([]interface{})
	if !ok {
		return fmt.Errorf("error on reading dedicated host set")
	}

	datas := GetSubSliceDByRep(allDedicatedHosts, dedicatedHostKeys)
	dealDedicatedHostData(datas)
	err = dataSourceKscSave(d, "dedicated_hosts", dedicated_hosts, datas)
	if err != nil {
		return fmt.Errorf("error on save dedicated host list, %s", err)
	}
	return nil
}
func dealDedicatedHostData(datas []map[string]interface{}) {
	for k, v := range datas {
		for k1, v1 := range v {
			switch k1 {
			case "dedicated_host_name":
				datas[k]["name"] = v1
				delete(datas[k], "dedicated_host_name")
			case "model":
				datas[k]["dedicated_type"] = v1
				delete(datas[k], "model")
			}
		}
	}
}
