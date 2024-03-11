package db

import (
	"database/sql"
	"forum/internal/category"
	"log"
)

type db struct {
	db *sql.DB
}

func NewCategoryStorage(dataBase *sql.DB) category.CategoryStorage {
	return &db{
		db: dataBase,
	}
}

func (d *db) GetCategory() (c []category.Category, err error) {
	query := `SELECT * FROM categories`
	rows, err := d.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var category category.Category
		if err := rows.Scan(&category.CategoryName); err != nil {
			log.Println(err)
			return nil, err
		}
		c = append(c, category)
	}
	return
}

func (d *db) GetCategoryByName(name string) (c category.Category, err error) {
	row := d.db.QueryRow("SELECT category FROM categories WHERE category = ? ", name)
	if err = row.Scan(&c.CategoryName); err != nil {
		return
	}
	return
}

func (d *db) GetCategoryByPostID(id int64) (c []category.Category, err error) {
	query := `SELECT c.category FROM categories c
    JOIN categories_posts pc ON c.category = pc.category
    WHERE pc.post_id = $1`

	rows, err := d.db.Query(query, id)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		category := category.Category{}

		err = rows.Scan(&category.CategoryName)
		if err != nil {
			log.Println(err)
			return
		}

		c = append(c, category)
	}
	return
}

func (d *db) CreateCategoryPostByPostID(name string, id int64) (err error) {
	_, err = d.db.Exec("INSERT INTO categories_posts (category, post_id) VALUES (?,?)", name, id)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

