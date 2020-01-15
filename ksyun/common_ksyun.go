package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"log"
	"strings"
)

//Convert sdk response type (map[string]interface{}) to the type terraform can realized([]map[string]interface).
//params data limit ： [k,v]:the type of k must be string ,the type of v must be basic type.
func GetSubDByRep(data interface{}, include, exclude map[string]bool) []interface{} {
	ma, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	subD := make(map[string]interface{})
	for k, v := range ma {
		if exclude[k] || !include[k] {
			continue
		}
		subD[Hump2Downline(k)] = v
	}
	return []interface{}{subD}
}

//sdk resp []map[string]interface{}->terraform schema.ResourceData
//Convert sdk response type ([]map[string]interface{}) to the type terraform can realized([]map[string]interface).
//include ：representing the key terraform has defined.
//exclude ：representing the key which the type is not basic type.
//Suitable for the value in d.Set（ key，value）,and the type of value must be List.
func GetSubSliceDByRep(items []interface{}, include /*,exclude*/ map[string]bool) []map[string]interface{} {
	datas := []map[string]interface{}{}
	for _, v := range items {
		data := map[string]interface{}{}
		vv, _ := v.(map[string]interface{})
		for key, value := range vv {
			//ignore keys whose type is not basic type,and need to deal later.
			if /*exclude[key]||*/ !include[key] {
				continue //if not judge,sdk may set value to terraform which can identify,and will panic.
			}
			data[Hump2Downline(key)] = value
		}
		datas = append(datas, data)
	}
	return datas
}

//sdk resp map[string]interface{} inline struct ->terraform schema.ResourceData
//convert inline struct from sdk response type ([]map[string]interface{}) to the type terraform can realized([]map[string]interface).
//exclude ：representing the key which the type is not basic type.
func GetSubStructDByRep(datas interface{}, exclude map[string]bool) map[string]interface{} {

	subStruct := map[string]interface{}{}
	items, ok := datas.(map[string]interface{})
	if !ok {
		return subStruct
	}
	for k, v := range items {
		if exclude[k] {
			continue
		}
		subStruct[Hump2Downline(k)] = v
	}
	return subStruct
}

//set sdk response (map[string]interface{}) to the terr`aform ([]map[string]interface).
//params data limit ： [k,v]:the type of k must be string ,the type of v must be basic type.
//exclude ：representing the key which the type is not basic type (terraform can't identity the type which is not basic type).
//mre: the params not set to terraform .
func SetDByRespV1(d *schema.ResourceData, m interface{}, exclud map[string]bool) map[string]interface{} {
	ma, ok := m.(map[string]interface{})
	mre := make(map[string]interface{}, 0)
	if !ok {
		return mre
	}
	for k, v := range ma {
		if exclud[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				mre[k] = mm
			} else {
				mre[k] = v
			}
			continue
		}
		d.Set(Hump2Downline(k), v)
	}
	return mre
}

//set sdk response (map[string]interface{}) to the terraform ([]map[string]interface).
//params data limit ： [k,v]:the type of k must be string ,the type of v must be basic type.
//include ：representing the key terraform has defined. terraform will panic if set the key that not defined.
//exclude ：representing the key which the type is not basic type (terraform can't identity the type which is not basic type).
//mre: the params not set to terraform .
func SetDByResp(d *schema.ResourceData, m interface{}, includ, exclude map[string]bool) map[string]interface{} {
	mre := make(map[string]interface{}, 0)
	ma, ok := m.(map[string]interface{})
	if !ok {
		return mre
	}
	for k, v := range ma {
		if !includ[k] || exclude[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				mre[k] = mm
			} else {
				mre[k] = v
			}
			continue
		}

		d.Set(Hump2Downline(k), v)
	}
	return mre
}

//The hump is converted to an underline simply, and no special treatment is required for even uppercase letters.
//ex:aDDCC ->a_d_d_c_c
func Hump2Downline(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 string
	if len(s) == 1 {
		s1 = strings.ToLower(s[:1])
		return s1
	}
	for k, v := range s {
		if k == 0 {
			s1 = strings.ToLower(s[0:1])
			continue
		}
		if v >= 65 && v <= 90 {
			v1 := "_" + strings.ToLower(s[k:k+1])
			s1 = s1 + v1
		} else {
			s1 = s1 + s[k:k+1]
		}
	}
	return s1
}

//The underline is converted to an hump simply.
func Downline2Hump(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 []string
	ss := strings.Split(s, "_")
	for _, v := range ss {
		vv := strings.ToUpper(v[:1]) + v[1:]
		s1 = append(s1, vv)
	}
	return strings.Join(s1, "")
}

//flattern struct
// convert input param struct to map when create(with out prefix).
func FlatternStruct(v interface{}, req *map[string]interface{}) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				vv := Downline2Hump(k2)
				(*req)[vv] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//flattern struct Suitable for inline struct
//convert input param struct to map when create(with  prefix).
//prefix: the name of the outer structure
func FlatternStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				kk := Downline2Hump(k2)
				(*req)[fmt.Sprintf("%s.%s", prex, kk)] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//FlatternStructSlicePrefix 用于创建时，结构体切片类型的入参转换为map型 ,【
//Flattern StructSlice Suitable for the slice of inline struct
//convert input param struct to map when create(with  prefix).
//prefix: the name of the slice
func FlatternStructSlicePrefix(values interface{}, req *map[string]interface{}, prex string) {
	v, _ := values.([]interface{})
	k := 0
	for _, v1 := range v {
		vv := v1.(map[string]interface{})
		if len(vv) == 0 {
			continue
		}
		k++
		for k2, v2 := range vv {
			kk := Downline2Hump(k2)
			(*req)[fmt.Sprintf("%s.%d.%s", prex, k, kk)] = fmt.Sprintf("%v", v2)
		}
	}
}

//Suitable for filter which need conver param with "_"(terraform) to "-"(sdk) when read .
//convert input param struct to map when create(without prefix).
func ConvertFilterStruct(v interface{}, req *map[string]interface{}) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				vv := strings.ReplaceAll(k2, "_", "-")
				(*req)[vv] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//Suitable for filter which need conver param with "_"(terraform) to "-"(sdk) when read.
//convert input param struct to map when create(with prefix).
//prefix:the name of the elemet from filter
func ConvertFilterStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			if v1[0] == nil {
				return
			}
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				vv := strings.ReplaceAll(k2, "_", "-")
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = v2
			}
		}
	}
}

/*
func ConvertFilterStructStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			if v1[0] == nil {
				return
			}
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := strings.ReplaceAll(k2, "-", "_")
				v3, ok3 := v2.([]string)
				if !ok3 || len(v3) == 0 {
					(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v2)
				}
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v3[0])

			}
		}
	}
}

*/
func dataSourceKscSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {

	d.SetId(hashStringArray(ids))
	d.Set("total_count", len(datas))

	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), datas)
	}

	return nil
}
func dataSourceKscSaveSlice(d *schema.ResourceData, dataKey string, ids []string, datas []string) error {

	d.SetId(hashStringArray(ids))
	d.Set("total_count", len(datas))

	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), datas)
	}

	return nil
}

func dataSourceDbSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {
	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(strings.Join(ids, ","))
	}

	d.Set("total_count", len(datas))
	log.Printf("reset  dataKey: %v datas: %v", dataKey, datas)

	if err := d.Set(dataKey, datas); err != nil {
		logger.DebugInfo("err %+v", err)
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		logger.DebugInfo(" output file name : %+v", outputFile.(string)+"_"+d.Id())
		writeToFile(outputFile.(string)+"_"+d.Id(), datas)
	} else {
		return fmt.Errorf(" output file error,  %+v", outputFile)
	}
	return nil
}

func dataDbSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {
	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(strings.Join(ids, ","))
	}

	d.Set("total_count", len(datas))
	log.Printf("reset  dataKey: %v datas: %v", dataKey, datas)

	if err := d.Set(dataKey, datas); err != nil {
		logger.DebugInfo("$$err$$ %+v", err)
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		logger.DebugInfo(" ------------ %+v", outputFile)

		writeToFile(outputFile.(string)+"_"+"data", datas)
	} else {
		return fmt.Errorf(" !!! %+v", outputFile)
	}
	return nil
}

func FuckHump2Downline(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 string
	if len(s) == 1 {
		s1 = strings.ToLower(s[:1])
		return s1
	}
	for k, v := range s {
		if k == 0 {
			s1 = strings.ToLower(s[0:1])
			continue
		}
		if v >= 65 && v <= 90 {
			v1 := "_" + strings.ToLower(s[k:k+1])
			s1 = s1 + v1
		} else {
			s1 = s1 + s[k:k+1]
		}
	}
	if len(s1) < 3 {
		return s1
	}

	var s2 string
	if s1[:3] == "d_b" {
		s2 = "db" + s1[3:]
	} else {
		if s1 == "read_replica_d_b_instance_identifiers" {
			s2 = "read_replica_db_instance_identifiers"
		} else if s1 == "read_replica_d_b_instance_identifier" {
			s2 = "read_replica_db_instance_identifier"
		} else {
			s2 = s1
		}
	}
	var s3 string
	if s2[len(s2)-2:] == ".1" || s2[len(s2)-2:] == ".2" {
		s3 = s2[:len(s2)-2] + "_" + s2[len(s2)-1:]
	} else {
		s3 = s2
	}

	return s3
}
