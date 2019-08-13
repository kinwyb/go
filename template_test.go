package heldiamgo

import (
	"testing"
)

func Test_Template(t *testing.T) {
	var template = `{{if d}}22{{else if e}}333{{else}}11{{/if}}{{#a}}你好`
	tp := NewJSTemplate("", "")
	data, err := tp.Template([]byte(template), map[string]interface{}{"a": "haha", "d": false, "e": false})
	if err != nil {
		println(err.Error())
		t.Error(err)
	} else {
		println(string(data))
	}
}

func Test_e(t *testing.T) {

}
