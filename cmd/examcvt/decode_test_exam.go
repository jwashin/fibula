package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := ioutil.ReadFile("test_exam_blob.txt")
	check(err)

	sDec, _ := base64.StdEncoding.DecodeString(string(data))

	var exam ExamV1

	err = json.Unmarshal(sDec, &exam)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
	}

	// JSON object parses into a map with string keys

	var exam2 ExamV2
	exam2.ID = exam.ID
	// exam2.Time = exam.Time
	exam2.Title = exam.Title

	// v := reflect.ValueOf(exam.Items)
	x := reflect.TypeOf(exam.Items)

	// values := make([]interface{}, v.NumField())
	ids := make([]string, x.NumField())
	// var id string
	// var item Item

	for i := 0; i < x.NumField(); i++ {
		// values[i] = v.Field(i).Interface()
		ids[i] = x.Field(i).Type.Name()

	}
	fmt.Println(ids)
	m, _ := json.Marshal(exam.Items)
	// chucking something into JSON and back for conversion
	// is slow, but we won't have to do this much.
	var itemsmap map[string]interface{}
	_ = json.Unmarshal(m, &itemsmap)
	// fmt.Println(itemsmap)
	for _, val := range ids {
		idx := val[3:]
		jsonitem := itemsmap[idx]
		var item ExamItem
		// fmt.Println(jsonitem)
		newjson, _ := json.Marshal(jsonitem)
		json.Unmarshal(newjson, &item)
		item.ID = idx
		fmt.Println(item)
		exam2.Items = append(exam2.Items, item)
	}
	sort.Sort(sortedItems(exam2.Items))
	data2, _ := json.MarshalIndent(exam2, "", "\t")
	ioutil.WriteFile("test_quiz.json", data2, 0664)

}
