package model_test

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/k0kubun/pp"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

func aTestGetdata(t *testing.T) {
	url := "http://13.114.229.73:5000/dsp-campaign/getMany"
	data := []byte(`{"cmd":"select * from BaseDspCampaign where id=2  limit 1"}`)
	res, err := util.DoBytesPost(url, data)
	if err != nil {
		log.Printf("request err:%s", err)
	}
	mcom := &model.Compaign{}
	compaign := mcom.GetData(res)
	pp.Println(compaign[0].ID)
}

func aTestJsonParse(t *testing.T) {
	count := 20
	// a, b, c, d := jsonparser.Get(data, "data", "[0]", "campaign", "id")
	// pp.Println(string(a), b, c, d)
	t1 := time.Now() // get current time
	for i := 0; i < count; i++ {
		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			// a, _, _, _ := jsonparser.Get(value, "uuid")
			jsonparser.Get(value, "uuid")
			// fmt.Println(string(a))
		}, "data", "[0]", "creatives")
	}
	elapsed := time.Since(t1)
	fmt.Println("jsonparser: ", elapsed)

	t2 := time.Now() // get current time

	for i := 0; i < count; i++ {
		var cmps model.ResponseCompaign
		json.Unmarshal(data, &cmps)
		for _, v := range cmps.Data[0].Creatives {
			_ = v.UUID
			// fmt.Println(v.UUID)
		}
	}
	elapsed2 := time.Since(t2)
	fmt.Println("json.Unmarshal: ", elapsed2)
}

var data = []byte(`{"status":1,"message":"success","data":[{"campaign":{"id":2,"userId":14,"trackingCampaignId":997,"name":"dsp campaign native","hash":"bb24e4e0-c337-4c5f-bd63-cdba943a5923","status":true,"type":"NATIVE","domain":"nike.com","revenueType":"CPAFIXED","revenueValue":0.0001,"conversionActionType":"NUN","conversionActionUrl":"http://www.cwzpvo.com/postback?cid=REPLACE&payout=OPTIONAL&txid=OPTIONAL","bidPrice":1.0003,"dailyBudget":0.00001,"spendStrategy":"SMOOTH","dailyPerPlacementBudget":0.000002,"totalBudget":0.00005,"unlimitedBudget":false,"freqCapEnabled":false,"freqCapType":"USER","freqCountLimit":10,"freqTimeWindow":12,"activityPeriodsTimezone":"+08:00","dayParting":[[true,true,true,true,true,false,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true],[true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true,true]],"countries":["AF","AR"],"countries_is_all":0,"states":["Beijing"],"cities":["Xiongan"],"adExchanges":["ba75911f-4a48-49f1-a860-57e3e4bbef72"],"adx_is_all":0,"sourceType":"APP","clientIds":["942229a9-96f0-11e7-8095-02c8570e7558","a8f89911-96f0-11e7-8095-02c8570e7558"],"clientType":"BLACK_LIST","categories":[{"id":4,"key":"IAB19","name":"Technology & Computing"},{"id":6,"key":"IAB20-17","name":"Honeymoons/Getaways"}],"categories_is_all":0,"connectionType":"WIFI","carriers":["AWCC","Roshan"],"deviceTypes":{"TABLET":false,"MOBILE":true,"DESKTOP":false},"oses":["Windows Phone 7","Windows Phone 8"],"isIdfaGaid":true,"ips":["123.123.123.123","123.123.123.123/12"],"audienceIds":["ba75911f-4a48-49f1-a860-57e3e4bbef72"],"audienceType":"BLACK_LIST","redirectType":"FLOW","flowId":"123123","redirectUrl":"http://iytg3a.nbtrk7.com/ecfe9e37-6fab-4c02-afd4-3056ff3b4dff","score":0,"fromDate":"2017-08-08","fromTime":"12:00","toDate":"2017-08-09","toTime":"12:00"},"creatives":[{"uuid":"bb24e4e0-c337-4c5f-bd63-cdba943a5923","status":true,"rejectionComment":"","delete":false,"approvalStatus":"PADDING","nativeBanner":{"brand":"testBrand","brandingText":"testbrandingText","ctaText":"testCtaTExt","headline":"testHeadline","rating":0,"iconImage":[{"id":"777f5e36-210d-4c5b-94fe-8d33b97464d0","bannerMeta":{"cdnUrl":"upload/QHL_EZVJRzbYjny_sltTcNey.png","height":100,"width":100,"size":7085,"mime":"image/png"}}],"mainImage":[{"id":"82a4b5b4-9923-11e7-8095-02c8570e7558","bannerMeta":{"cdnUrl":"mlgbee.com/images/123.png","height":200,"width":200,"size":12000,"mime":"image/png"}}]}},{"uuid":"bb24e4e0-c337-4c5f-bd63-cdba943a5924","status":true,"rejectionComment":"","delete":false,"approvalStatus":"PADDING","nativeBanner":{"brand":"testBrand","brandingText":"testbrandingText","ctaText":"testCtaTExt","headline":"testHeadline","rating":0,"iconImage":[{"id":"777f5e36-210d-4c5b-94fe-8d33b97464d0","bannerMeta":{"cdnUrl":"upload/QHL_EZVJRzbYjny_sltTcNey.png","height":100,"width":100,"size":7085,"mime":"image/png"}}],"mainImage":[{"id":"82a4b5b4-9923-11e7-8095-02c8570e7558","bannerMeta":{"cdnUrl":"mlgbee.com/images/123.png","height":200,"width":200,"size":12000,"mime":"image/png"}}]}}]}]}`)
