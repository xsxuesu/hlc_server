// lcserver
package server

import (
	"bytes"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"log"
	"strconv"

	"model"
	"errors"
)

type Asset struct {
	Success		bool	`json:"success"`
	Payloads	[]string	`json:"payloads"`
	Timestamp	int64	`json:"timestamp"`
	Message		string	`json:"message"`
	Messages	[]string	`json:"messages"`
}
//批量入栏 和 单个 入栏相同
func (this Remote)AssetRegister(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"PCList"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"CreatePerson","username","PCNO","isType","species","TaskGps","createTime"})
	if err1 != nil {
		return nil
	}
	var str string

	for i := 0; i < len(mx["PCList"].([]interface{})); i++ {
		switch mx["PCList"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "PCList[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["PCList"].([]interface{})[i].(string) + `","batchNumber":"` + mx["PCNO"].(string) + `","kind":"` + mx["isType"].(string) + `","type":"` + mx["species"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `","operation":"Lairage","operator":"` + mx["CreatePerson"].(string) + `","createTime":"` + mx["createTime"].(string) + `"},`
	}

	batchOrSingleOperate("Register",str,args["header"]["Authorization"][0],result)
	return nil
}

func (this Remote)AssetQueryDetail(args map[string]map[string][]string, result *Asset) error {
	m := make(map[string]interface{})
	m["fcn"] = "QueryProductDetail"
	m["args"] = make([]string, 2)
	m["args"].([]string)[0] = m["fcn"].(string)
	m["args"].([]string)[1] = args["body"]["b"][0]
	m["peers"] = []string{"peer0.creator.com"}

	mJSON, err := json.Marshal(m)
	if err != nil {
		log.Println(err.Error())
	}

	reader := bytes.NewReader(mJSON)

	request, err := http.NewRequest("POST", "http://" + model.CHAIN_CODE_DOMAIN + ":" + model.CHAIN_CODE_PORT + "/query/channels/mychannel/chaincodes/hlccc", reader)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	request.Header.Set("authorization", "Bearer " + args["header"]["Authorization"][0])

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		log.Println(err.Error())
		result.Message = "500:服务器内部错误:" + err.Error() + string(body)
		return err
	}

	return nil
}
//批量喂养  喂养
func (this Remote)AssetFeed(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"PNO"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"Name","Id","SysDate","TaskGps"})
	if err1 != nil {
		return nil
	}

	var str string
	for i := 0; i < len(mx["PNO"].([]interface{})); i++ {
		switch mx["PNO"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "PNO[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["PNO"].([]interface{})[i].(string) + `","operator":"` + mx["Name"].(string) + `","feedName":"` + mx["Id"].(string) + `","feedTime":"` + mx["SysDate"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `","operation":"Feed"}`
	}
	batchOrSingleOperate("Feed",str,args["header"]["Authorization"][0],result)
	return nil
}

//防疫 批量防疫
func (this Remote)AssetMedication(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"PNO"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"OperatorName","SysDate","id","TaskGps"})
	if err1 != nil {
		return nil
	}
	var str string
	for i := 0; i < len(mx["PNO"].([]interface{})); i++ {
		switch mx["PNO"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "PNO[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["PNO"].([]interface{})[i].(string) + `","operator":"` + mx["OperatorName"].(string) + `","vaccineName":"` + mx["id"].(string) + `","VaccineTime":"` + mx["SysDate"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `","operation":"Medication","vaccineType":"vaccineType","vaccineNumber":"vaccineNumber"}`

	}
	batchOrSingleOperate("Vaccine",str,args["header"]["Authorization"][0],result)
	return nil
}
//检疫 批量检疫
func (this Remote)AssetPrevention(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"EarTag"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"OperatorName","Immunion","CheckDate","CheckResult","TaskGps"})
	if err1 != nil {
		return nil
	}

	var str string
	for i := 0; i < len(mx["EarTag"].([]interface{})); i++ {
		switch mx["EarTag"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "EarTag[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["EarTag"].([]interface{})[i].(string) + `","operator":"` + mx["OperatorName"].(string) + `","examConsequence":"` + mx["CheckResult"].(string) + `","examTime":"` + mx["CheckDate"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `","operation":"Prevention"}`
	}
	batchOrSingleOperate("Exam",str,args["header"]["Authorization"][0],result)
	log.Println(result)

	return nil
}

//救治 批量救治
func (this Remote)AssetSave(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"PNO"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"Name","InspectResult","SysDate","Treatment","TaskGps","id"})
	if err1 != nil {
		return nil
	}

	var str string
	for i := 0; i < len(mx["PNO"].([]interface{})); i++ {
		switch mx["PNO"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "PNO[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["PNO"].([]interface{})[i].(string) + `","operator":"` + mx["Name"].(string) + `","saveNumber":"` + mx["id"].(string) + `","saveName":"saveName","saveType":"` + mx["Treatment"].(string) + `","saveConsequence":"` + mx["InspectResult"].(string) + `","saveTime":"` + mx["SysDate"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `","operation":"Save"}`

	}
	batchOrSingleOperate("Save",str,args["header"]["Authorization"][0],result)

	return nil
}

//灭失 批量灭失
func (this Remote)AssetLost(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"DeathObject"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"Name","CauseDeath","SysDate","TreatMethod","TaskGps"})
	if err1 != nil {
		return nil
	}

	var str string
	for i := 0; i < len(mx["DeathObject"].([]interface{})); i++ {
		switch mx["DeathObject"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "DeathObject[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["DeathObject"].([]interface{})[i].(string) + `","operator":"` + mx["Name"].(string) + `","lostWay":"` + mx["TreatMethod"].(string) + `","lostReaso":"` + mx["CauseDeath"].(string) + `","lostTime":"` + mx["SysDate"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `","operation":"Lost"}`

	}
	batchOrSingleOperate("Lost",str,args["header"]["Authorization"][0],result)

	return nil
}

func (this Remote)AssetFattened(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"PNO"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"CreatePerson","CLPCNO","Name","SysDate","TaskGps"})
	if err1 != nil {
		return nil
	}
	var str string
	for i := 0; i < len(mx["PNO"].([]interface{})); i++ {
		switch mx["PNO"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "PNO[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["PNO"].([]interface{})[i].(string) + `","name":"` + mx["Name"].(string) + `","outputTime":"` + mx["SysDate"].(string) + `","operation":"Fattened","operator":"` + mx["CreatePerson"].(string) + `","mapPosition":"` + mx["TaskGps"].(string) + `"}`
	}
	batchOrSingleOperate("Output",str,args["header"]["Authorization"][0],result)

	return nil
}
func (this Remote)AssetButcher(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"productIds"})
	if err2 != nil {
		return nil
	}

	err1 := verifyParamString(mx,result,[] string{"operator","hookNo","operation","butcherTime","mapPosition"})
	if err1 != nil {
		return nil
	}
	var str string
	for i := 0; i < len(mx["productIds"].([]interface{})); i++ {
		switch mx["productIds"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "productIds[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["productIds"].([]interface{})[i].(string) + `","hookNo":"` + mx["hookNo"].(string) + `","butcherTime":"` + mx["butcherTime"].(string) + `","operation":"`+mx["operation"].(string)+`","operator":"` + mx["operator"].(string) + `","mapPosition":"` + mx["mapPosition"].(string) + `"}`
	}
	batchOrSingleOperate("Butcher",str,args["header"]["Authorization"][0],result)

	return nil
}


func (this Remote)AssetWaitButcher(args map[string]map[string][]string, result *Asset) error {
	var mx map[string]interface{}
	err := json.Unmarshal([]byte(args["body"]["b"][0]), &mx)
	if err != nil {
		log.Println(err)
		result.Message = err.Error()
		return nil
	}
	err2 := verifyParamArrayString(mx,result,[] string{"productIds"})
	if err2 != nil {
		return nil
	}
	err1 := verifyParamString(mx,result,[] string{"operator","operation","waitButcherTime","mapPosition"})
	if err1 != nil {
		return nil
	}

	var str string
	for i := 0; i < len(mx["productIds"].([]interface{})); i++ {
		switch mx["productIds"].([]interface{})[i].(type) {
		case string:
		default:
			result.Message = "productIds[" + strconv.Itoa(i) + "] should be string"
			return nil
		}

		str += `{"productId":"` + mx["productIds"].([]interface{})[i].(string) + `","waitButcherTime":"` + mx["waitButcherTime"].(string) + `","operation":"`+mx["operation"].(string)+`","operator":"` + mx["operator"].(string) + `","mapPosition":"` + mx["mapPosition"].(string) + `"}`
	}
	batchOrSingleOperate("Butcher",str,args["header"]["Authorization"][0],result)

	return nil
}

//批量或者单个操作
func batchOrSingleOperate(fcn string,str string,auth string ,result *Asset){
	m := make(map[string]interface{})
	m["fcn"] = fcn
	m["args"] = make([]string, 1)
	m["args"].([]string)[0] = "[" + str[0:len(str) - 1] + "]"

	log.Println(m)

	mJSON, err := json.Marshal(m)
	if err != nil {
		log.Println(err.Error())
		result.Message = "JSON Marshal error:" + err.Error()
		return
	}

	reader := bytes.NewReader(mJSON)

	request, err := http.NewRequest("POST", "http://" + model.CHAIN_CODE_DOMAIN + ":" + model.CHAIN_CODE_PORT + "/channels/mychannel/chaincodes/jiakechaincode", reader)
	if err != nil {
		log.Println(err.Error())
		result.Message = err.Error()
		return
	}

	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	request.Header.Set("authorization", "Bearer " + auth)

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		result.Message = err.Error()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		result.Message = err.Error()
		return
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		log.Println(err.Error())
		result.Message = "JSON Unmarshal error:" + err.Error()
		return
	}

	log.Println(result.Message) // = append(result.Messages, mx["PCList"].([]interface{})[i].(string) + ":" + string(body))

	return
}


func verifyParamString( mx map[string]interface{} , result *Asset , fields []string ) error {
	for i :=0;i<len(fields) ;i++  {
		_, ok := mx[fields[i]]
		if !ok {
			result.Message = fields[i] + " required"
			return errors.New(result.Message)
		}
		switch mx[fields[i]].(type) {
		case string:
		default:
			result.Message = fields[i] + " should be string"
			return errors.New(result.Message)
		}
	}
	return nil
}


func verifyParamArrayString( mx map[string]interface{} , result *Asset , fields []string ) error {
	for i :=0;i<len(fields) ;i++  {
		_, ok := mx[fields[i]]
		if !ok {
			result.Message = fields[i] + " required"
			return errors.New(result.Message)
		}
		switch mx[fields[i]].(type) {
		case []interface{}:
		default:
			result.Message = fields[i] + " should be []string"
			return errors.New(result.Message)
		}
	}
	return nil
}
