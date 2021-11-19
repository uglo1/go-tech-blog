package repository

import (
	"database/sql"
	"time"

	"go-tech-blog/model"
)

// ArticleList ...
func ArticleList() ([]*model.Article, error) {
	query := `SELECT * FROM articles;`

	var articles []*model.Article
	if err := db.Select(&articles, query); err != nil {
		return nil, err
	}

	return articles, nil
}

// ArticleCreate ...
func ArticleCreate(article *model.Article) (sql.Result, error) {
	now := time.Now()

	article.Created = now
	article.Updated = now

	query := `INSERT INTO articles (title, body, created, updated)
						VALUES (:title, :body, :created, :updated);`

	tx := db.MustBegin()

	// :titleなどは構造体の値で置換される
	res, err := tx.NamedExec(query, article)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return res, nil
}
