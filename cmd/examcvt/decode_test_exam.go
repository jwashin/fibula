package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sort"

	"github.com/jwashin/fibula.git/pkg/exam"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data, err := ioutil.ReadFile("../../test/data/test_exam_blob.txt")
	check(err)

	sDec, _ := base64.StdEncoding.DecodeString(string(data))

	var oldexam exam.V1

	err = json.Unmarshal(sDec, &oldexam)
	if err != nil {
		log.Println("Error parsing JSON: ", err)
	}

	var newexam exam.Examination

	newexam.ID = oldexam.ID
	newexam.Title = oldexam.Title

	// We're going though a horribke struct of
	// things that are the same, but have different types,
	// named Num001, Num002, etc.
	x := reflect.TypeOf(oldexam.Items)

	// values := make([]interface{}, v.NumField())
	ids := make([]string, x.NumField())

	for i := 0; i < x.NumField(); i++ {
		ids[i] = x.Field(i).Type.Name()
	}

	fmt.Println(ids)
	m, _ := json.Marshal(oldexam.Items)
	// chucking something into JSON and back for conversion
	// is slow, but we won't have to do this much.
	var itemsmap map[string]interface{}
	_ = json.Unmarshal(m, &itemsmap)
	// fmt.Println(itemsmap)
	for _, val := range ids {
		idx := val[3:]
		jsonitem := itemsmap[idx]
		var item exam.Item
		// fmt.Println(jsonitem)
		newjson, _ := json.Marshal(jsonitem)
		json.Unmarshal(newjson, &item)
		item.ID = idx
		fmt.Println(item)
		newexam.Items = append(newexam.Items, item)
	}
	sort.Sort(exam.ItemsByID(newexam.Items))
	data2, _ := json.MarshalIndent(newexam, "", "\t")
	ioutil.WriteFile("../../test/data/test_quiz.json", data2, 0664)

}
