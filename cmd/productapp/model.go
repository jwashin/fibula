// model.go

package productapp

import (
	"context"

	"cloud.google.com/go/datastore"
)

/*id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
	CONSTRAINT products_pkey PRIMARY KEY (id)*/

// Product has the basic info about the Product
type Product struct {
	ID    *datastore.Key `json:"id" datastore:"__key__"`
	Name  string         `json:"name"`
	Price float32        `json:"price"`
}

// // Product is product information
// type Product struct {
// 	ID    int     `json:"id"`
// 	Name  string  `json:"name"`
// 	Price float64 `json:"price"`
// }

func (p *Product) getProduct(db *datastore.Client) error {
	context := context.Background()
	key := p.ID
	if key != nil {
		// productKey := datastore.NameKey("Product", p.ID, nil)
		err := db.Get(context, key, p)
		// t, err := json.Marshal(p)

		if err != nil {
			// fmt.Printf("this error is from Get: %v (%v)", err, key)
			return err
		}
	}
	return nil
}

func (p *Product) updateProduct(db *datastore.Client) error {
	context := context.Background()
	if p.ID != nil {
		productKey := p.ID
		_, err := db.Put(context, productKey, p)
		if err != nil {
			return err
		}
	}
	return nil

}

func (p *Product) deleteProduct(db *datastore.Client) error {
	context := context.Background()
	if p.ID != nil {
		// productKey := datastore.NameKey("Product", p.ID, nil)
		err := db.Delete(context, p.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Product) createProduct(db *datastore.Client) error {
	context := context.Background()
	productKey := datastore.IncompleteKey("Product", nil)
	_, err := db.Put(context, productKey, p)
	if err != nil {
		return err
	}
	return nil
}

func getProducts(db *datastore.Client, start, count int) ([]Product, error) {
	// fmt.Printf("start is %v and count is %v", start, count)
	context := context.Background()
	query := datastore.NewQuery("Product")
	var products []Product
	_, err := db.GetAll(context, query, products)

	if err != nil {
		// fmt.Printf("%v", err)
		if err == datastore.ErrInvalidEntityType {
			// we got "invalid entity type". return empty list
			return []Product{}, nil
		}
	}

	if start+count > len(products) {
		return products, err
	}

	return products[start : start+count], err
}
