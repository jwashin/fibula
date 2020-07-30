// main_test.go

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud.google.com/go/datastore"
)

var a App

// var usedKeys []*datastore.Key

func TestMain(m *testing.M) {
	a.Initialize()
	code := m.Run()
	// clearTable()
	os.Exit(code)
}

func clearProducts() {
	// get all the keys for Products
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "")

	keysQuery := datastore.NewQuery("Product").KeysOnly()
	keys, _ := client.GetAll(ctx, keysQuery, nil)

	// if any, delete them
	if len(keys) > 0 {
		client.DeleteMulti(ctx, keys)
	}

}

// const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
// (
//     id SERIAL,
//     name TEXT NOT NULL,
//     price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
//     CONSTRAINT products_pkey PRIMARY KEY (id)
// )`

func TestEmptyTable(t *testing.T) {
	clearProducts()

	req, _ := http.NewRequest("GET", "/products", nil)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearProducts()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {

	clearProducts()

	// var jsonStr = []byte(`{"name":"test product", "id": "1", "price": 11.22}`)
	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}

	// if m["id"] != "1" {
	// t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	// }
}

func TestGetProduct(t *testing.T) {
	clearProducts()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// main_test.go

func addProducts(count int) {
	if count < 1 {
		count = 1
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		fmt.Printf("Could not add item to database: %v", err)
	}
	for i := 0; i < count; i++ {
		// idstr := strconv.Itoa(i + 1)
		key := datastore.Key{Kind: "Product", ID: int64(i + 1)}
		s := Product{Name: "Product", Price: 33.45}
		// k := datastore.NameKey("Product", idstr, nil)
		_, err := client.Put(ctx, &key, &s)
		// fmt.Printf("Adding %v to database", key)
		if err != nil {
			fmt.Printf("Could not add item to database: %v ", err)
		}
	}
}

func TestUpdateProduct(t *testing.T) {

	clearProducts()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	var originalProduct map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalProduct)

	var jsonStr = []byte(`{"name":"test product - updated name", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearProducts()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
