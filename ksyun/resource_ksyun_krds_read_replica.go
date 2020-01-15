package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/krds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"time"
)

//var krdsRrApiField = []string{
//	"DBInstanceIdentifier",
//	"DBInstanceClass",
//	"DBInstanceName",
//	"BillType",
//	"Duration",
//	"SecurityGroupId",
//	"AvailabilityZone.1",
//	"ProjectId",
//}
//var krdsRrTfField = []string{
//	"db_instance_identifier",
//	"db_instance_class",
//	"db_instance_name",
//	"bill_type",
//	"duration",
//	"security_group_id",
//	"availability_zone.1",
//	"project_id",
//}

func resourceKsyunKrdsRr() *schema.Resource {

	return &schema.Resource{
		Create: resourceKsyunKrdsRrCreate,
		Update: resourceKsyunKrdsRrUpdate,
		Read:   resourceKsyunKrdsRrRead,
		Delete: resourceKsyunKrdsRrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_instance_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "source instance identifier",
			},
			"db_instance_class": {
				Type:     schema.TypeString,
				Required: true,
				Description: "this value regex db.ram.d{1,3}|db.disk.d{1,5} , " +
					"db.ram is rds random access memory size, db.disk is disk size",
			},
			"db_instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bill_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "duration unit is month",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "proprietary security group id for mysql",
			},
			"availability_zone_1": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set:      parameterToHash,
				Optional: true,
				Computed: true,
			},
			// 与存入数据一致datakey
			"krds": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_class": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vcpus": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"disk": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"max_conn": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"mem": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
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
						"db_instance_status": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"db_instance_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"instance_create_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"master_user_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"publicly_accessible": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"read_replica_db_instance_identifiers": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vip": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"read_replica_db_instance_identifier": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"bill_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_source": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"multi_availability_zone": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_use": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"bill_type_id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"db_parameter_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"datastore_version_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_used": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
						},
						"preferred_backup_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"product_what": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"service_start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"order_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_order_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"audit": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"db_source": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							MaxItems: 1,
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
									"db_instance_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"service_end_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceKsyunKrdsRrCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	var err error
	creates := []string{
		"DBInstanceIdentifier",
		"DBInstanceClass",
		"DBInstanceName",
		"DBInstanceType",
		"BillType",
		"Duration",
		"SecurityGroupId",
		"AvailabilityZone.1",
		"ProjectId",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(FuckHump2Downline(v)); ok {
			createReq[v] = fmt.Sprintf("%v", v1)
		}
	}

	checkBackupComplete(d, meta)

	action := "CreateDBInstanceReadReplica"
	logger.Debug(logger.RespFormat, action, createReq)
	resp, err = conn.CreateDBInstanceReadReplica(&createReq)
	logger.Debug(logger.AllFormat, action, createReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on creating Instance(krds): %s", err)
	}

	if resp != nil {
		bodyData := (*resp)["Data"].(map[string]interface{})
		krdsInstance := bodyData["DBInstance"].(map[string]interface{})
		instanceId := krdsInstance["DBInstanceIdentifier"].(string)
		logger.DebugInfo("RR DBInstanceIdentifier : %v", instanceId)
		d.SetId(instanceId)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{tCreatingStatus},
		Target:     []string{tActiveStatus, tFailedStatus, tDeletedStatus, tStopedStatus},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    mysqlInstanceStateRefresh(conn, d.Id(), []string{tCreatingStatus}),
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on creating Instance(krds): %s", err)
	}

	err = resourceKsyunMysqlRead(d, meta)
	if err != nil {
		return fmt.Errorf("error on ModifyDBParameterGroup Instance(krds): %s", err)
	}

	return modifyParameters(d, meta)
}

func checkBackupComplete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).krdsconn
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"wait"},
		Target:     []string{"complete", "err"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
		Refresh:    mysqlBackupStateRefresh(conn, d.Id()),
	}
	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on check backup Instance(krds): %s", err)
	} else {
		return nil
	}
}

func mysqlBackupStateRefresh(client *krds.Krds, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"DBInstanceIdentifier": instanceId}
		action := "DescribeDBInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeDBBackups(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil {
			return nil, "err", err
		}
		bodyData := (*resp)["Data"].(map[string]interface{})
		backups := bodyData["DBBackup"].([]interface{})
		if len(backups) > 0 {
			return resp, "complete", nil
		} else {
			return resp, "wait", nil
		}

	}
}

func resourceKsyunKrdsRrUpdate(d *schema.ResourceData, meta interface{}) error {

	return fmt.Errorf("read replica instance do not support update ...")
}
func resourceKsyunKrdsRrRead(d *schema.ResourceData, meta interface{}) error {
	return resourceKsyunMysqlRead(d, meta)
}
func resourceKsyunKrdsRrDelete(d *schema.ResourceData, meta interface{}) error {
	return resourceKsyunMysqlDelete(d, meta)
}
