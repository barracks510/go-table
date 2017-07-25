package table

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestMakeTableEmpty(t *testing.T) {
	data := map[string]interface{}{}
	expected := []byte(`<table>
<thead>
<tr>

</tr>
</thead>

<tbody>
<tr>

</tr>
</tbody>
</table>
`)
	r, err := MakeTable(data)
	if err != nil {
		t.Error(err)
	}
	actual, _ := ioutil.ReadAll(r)
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("Got %s, expected %s", actual, expected)
	}
}

func TestMakeTableSimple(t *testing.T) {
	data := map[string]interface{}{"intValue": 1, "floatValue": 2.3}
	expected := []byte(`<table>
<thead>
<tr>
<th align="center">intValue</th><th align="center">floatValue</th>
</tr>
</thead>

<tbody>
<tr>
<td align="center">1</td><td align="center">2.3</td>
</tr>
</tbody>
</table>
`)
	r, err := MakeTable(data)
	if err != nil {
		t.Error(err)
	}
	actual, _ := ioutil.ReadAll(r)
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("Got %s, expected %s", actual, expected)
	}
}

func TestMakeTableComplex(t *testing.T) {
	data := map[string]interface{}{"tableValue": map[interface{}]interface{}{"intValue": 1}}
	expected := []byte(`<table>
<thead>
<tr>
<th align="center">tableValue</th>
</tr>
</thead>

<tbody>
<tr>
<td align="center"><table>
<thead>
<tr>
<th align="center">intValue</th>
</tr>
</thead>

<tbody>
<tr>
<td align="center">1</td>
</tr>
</tbody>
</table>
</td>
</tr>
</tbody>
</table>
`)
	r, err := MakeTable(data)
	if err != nil {
		t.Error(err)
	}
	actual, _ := ioutil.ReadAll(r)
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("Got %s, expected %s", actual, expected)
	}
}

func TestMakeTableArray(t *testing.T) {
	data := map[string]interface{}{"array": []interface{}{1, 2, 3}}
	expected := []byte(`<table>
<thead>
<tr>
<th align="center">array</th>
</tr>
</thead>

<tbody>
<tr>
<td align="center"><table>

<tbody>
<tr>
<td align="center">1</td><td align="center">2</td><td align="center">3</td>
</tr>
</tbody>
</table>
</td>
</tr>
</tbody>
</table>
`)
	r, err := MakeTable(data)
	if err != nil {
		t.Error(err)
	}
	actual, _ := ioutil.ReadAll(r)
	if bytes.Compare(expected, actual) != 0 {
		t.Errorf("Got %s, expected %s", actual, expected)
	}

}

func TestMakeTableBadType(t *testing.T) {
	data := []interface{}{}
	_, err := MakeTable(data)
	if err != ErrBadType {
		t.Errorf("Expected ErrBadType")
	}
}
