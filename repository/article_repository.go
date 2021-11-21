package repository

import (
	"database/sql"
	"math"
	"time"

	"go-tech-blog/model"
)

// // ArticleList ...
// func ArticleList() ([]*model.Article, error) {
// 	query := `SELECT * FROM articles;`

// 	var articles []*model.Article
// 	if err := db.Select(&articles, query); err != nil {
// 		return nil, err
// 	}

// 	return articles, nil
// }

// ArticleListByCursor ...
func ArticleListByCursor(cursor int) ([]*model.Article, error) {
	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	query := `SELECT *
	FROM articles
	WHERE id < ?
	ORDER BY id desc
	LIMIT 10`

	articles := make([]*model.Article, 0, 10)

	if err := db.Select(&articles, query, cursor); err != nil {
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

// ArticleDelete ...
func ArticleDelete(id int) error {
	query := `DELETE FROM articles WHERE id = ?`

	tx := db.MustBegin()

	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
