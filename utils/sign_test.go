// Copyright 2011 ZHU HAIHUA <kimiazhu@gmail.com>.
// All rights reserved.
// Use of this source code is governed by MIT
// license that can be found in the LICENSE file.

// Description: util
// Author: ZHU HAIHUA
// Since: 9/19/16
package util

import (
	"fmt"
	log "github.com/kimiazhu/log4go"
	"testing"
	"time"
)

func init() {
	cfg := `
	<logging>
	    <filter enabled="true">
		    <tag>stdout</tag>
		    <type>console</type>
		    <level>DEBUG</level>
		    <exclude>github.com/xgsdk2/betatest/tako.lib/mgox</exclude>
	    </filter>
	</logging>
	`
	log.Setup([]byte(cfg))
}

func TestSignVerify2(t *testing.T) {
	//jsStr := "{\"Key1\":\"Value1\",\"Key2\":\"\",\"Key3\":[\"Value3.1\",\"Value3.2\"],\"sign\":\"aef03ab309f230d32f612daee3fb882e4f533404\"}"
	jsStrs := []string{
		`{"Key1":"Value1","Key2":"","Key3":["Value3.1","Value3.2"],"key4":true,"sign":"aa021a87221279bcf5a1252c65d167a1a8737d4b"}`,
		`{"Key1":"Value1","Key2":"","Key3":[3.1,3.2],"key4":true,"sign":"5618bf54cc5099f0d09b2e2751d96d6b4d2e947d"}`,
		`{"Key1":"Value1","Key2":"","Key3":[true, false, true],"key4":2 ,"sign":"50502e0da04a258b22c86d6625622df55e2630b8"}`,
		`{"ts":"20160920114517178","xgAppId":17007,"sign":"6909b5efc0d1f5c245d65b0a9eef1ea697f548cf"}`,
	}
	for _,s := range jsStrs {
		r := SignVerify2(s, DefaultSecretKey)
		if !r {
			t.Errorf("sign %s, expert true, but false", s)
		} else {
			t.Logf("\n")
		}
	}

	time.Sleep(time.Second)
}

func TestBuildSignedJsonStr(t *testing.T) {
	r, _ := BuildSignedJsonStr(map[string]interface{}{"Key1": "Value1", "Key2": "", "Key3": []string{"Value3.1", "Value3.2"}}, DefaultSecretKey)
	fmt.Printf("%s\n", r)
}
