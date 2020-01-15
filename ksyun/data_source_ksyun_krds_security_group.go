package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

func dataSourceKsyunKrdsSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunKrdsSecurityGroupRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},

			// 与存入数据一致datakey
			"security_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"security_group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"security_group_description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_instance_identifier": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"db_instance_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"created": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"db_instance_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"security_group_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_rule_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"security_group_rule_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"security_group_rule_protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"created": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceKsyunKrdsSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	descReq := make(map[string]interface{})
	if v, ok := d.GetOk("security_group_id"); ok {
		descReq["SecurityGroupId"] = fmt.Sprintf("%v", v)
	}
	action := "DescribeSecurityGroup"
	logger.DebugInfo("+-+-+-+-+  %+v  ------  %T", descReq)
	logger.Debug(logger.ReqFormat, action, descReq)
	resp, err := conn.DescribeSecurityGroup(&descReq)
	logger.Debug(logger.AllFormat, action, descReq, *resp, err)

	if err != nil {
		return fmt.Errorf("error on request instance. security group id %q, %s", d.Id(), err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading response body, security group id %q, %+v", d.Id(), (*resp)["Error"])
	}
	instances := bodyData["SecurityGroups"].([]interface{})

	krdsIds := make([]string, len(instances))
	krdsMapList := make([]map[string]interface{}, len(instances))
	for num, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		krdsMap := make(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "Instances" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[FuckHump2Downline(j)] = q
						}
						wtf[num] = rrmap
					}
					krdsMap["instances"] = wtf
				}
			} else if k == "SecurityGroupRules" {
				rrids := v.([]interface{})
				if len(rrids) > 0 {
					wtf := make([]interface{}, len(rrids))
					for num, rrinfo := range rrids {
						rrmap := make(map[string]interface{})
						rr := rrinfo.(map[string]interface{})
						for j, q := range rr {
							rrmap[FuckHump2Downline(j)] = q
						}
						wtf[num] = rrmap
					}
					krdsMap["security_group_rules"] = wtf
				}
			} else {
				krdsMap[FuckHump2Downline(k)] = v
			}
		}
		logger.DebugInfo(" converted ---- %+v ", krdsMap)

		krdsIds[num] = krdsMap["security_group_id"].(string)
		logger.DebugInfo("krdsIds fuck : %v", krdsIds)
		krdsMapList[num] = krdsMap
	}
	//d.Set("security_group_id",krdsIds[0])
	logger.DebugInfo(" converted ---- %+v ", krdsMapList)
	dataDbSave(d, "security_groups", krdsIds, krdsMapList)

	return nil
}
