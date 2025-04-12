package product

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"mohaafaridd.dev/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)

	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) GetProductsByIds(productIds []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIds)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id In (?%s)", placeholders)

	args := make([]interface{}, len(productIds))

	for i, v := range productIds {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	products := make([]types.Product, 0)

	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) CreateProduct(p types.Product) error {
	_, err := s.db.Query("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)",
		p.Name, p.Description, p.Image, p.Price, p.Quantity,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *Store) UpdateProductQuantity(id, quantity int) error {
	_, err := s.db.Exec(
		"UPDATE products SET quantity = ? Where id = ?",
		id, quantity,
	)

	if err != nil {
		return err
	}

	return nil

}
