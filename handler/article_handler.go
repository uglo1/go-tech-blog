package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-tech-blog/model"
	"go-tech-blog/repository"

	"github.com/labstack/echo/v4"
)

// ArticleIndex ...
func ArticleIndex(c echo.Context) error {
	// 統一してGoogle Analyticsなどでのアクセス解析で分析しやすくなる
	if c.Request().URL.Path == "/articles" {
		c.Redirect(http.StatusPermanentRedirect, "/")
	}

	// リポジトリの処理を呼び出して記事の一覧データを取得する
	articles, err := repository.ArticleListByCursor(0)

	if err != nil {
		log.Println(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	// 取得できた最後の記事のIDをカーソルとして設定
	var cursor int
	if len(articles) != 0 {
		cursor = articles[len(articles)-1].ID
	}

	data := map[string]interface{}{
		"Articles": articles, // 記事データをテンプレートエンジンに渡す
		"Cursor":   cursor,
	}
	return render(c, "article/index.html", data)
}

// ArticleList ...
func ArticleList(c echo.Context) error {
	// クエリパラメータからカーソルの値を取得する
	cursor, _ := strconv.Atoi(c.QueryParam("cursor"))

	articles, err := repository.ArticleListByCursor(cursor)

	if err != nil {
		c.Logger().Error(err.Error())
		// JSON形式でデータのみを返却するので、c.HTMLBlob()でなく、c.JSON()を呼ぶ
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, articles)
}

// ArticleNew ...
func ArticleNew(c echo.Context) error {
	data := map[string]interface{}{
		"Message": "Article New",
		"Now":     time.Now(),
	}

	return render(c, "article/new.html", data)
}

// ArticleShow ...
func ArticleShow(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	article, err := repository.ArticleGetByID(id)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article": article,
	}

	return render(c, "article/show.html", data)
}

// ArticleEdit ...
func ArticleEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	article, err := repository.ArticleGetByID(id)

	if err != nil {
		c.Logger().Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Article": article,
	}

	return render(c, "article/edit.html", data)
}

// ArticleCreateOutput ...
type ArticleCreateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

// ArticleCreate ...
func ArticleCreate(c echo.Context) error {
	var article model.Article
	var out ArticleCreateOutput

	// フォームの内容を構造体に埋め込む
	if err := c.Bind(&article); err != nil {
		c.Logger().Error(err.Error())

		return c.JSON(http.StatusBadRequest, out)
	}

	// バリデーションチェックを実行する
	if err := c.Validate(&article); err != nil {
		c.Logger().Error(err.Error())
		out.ValidationErrors = article.ValidationErrors(err)
		// 解釈できたパラメータが許可されていない値の場合は、422エラーを返す。
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	// 保存処理を実行する
	res, err := repository.ArticleCreate(&article)
	if err != nil {
		c.Logger().Error(err.Error())

		return c.JSON(http.StatusInternalServerError, out)
	}

	id, _ := res.LastInsertId()
	article.ID = int(id)
	out.Article = &article

	return c.JSON(http.StatusOK, out)
}

// ArticleDelete ...
func ArticleDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("articleID"))

	if err := repository.ArticleDelete(id); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Article %d is deleted.", id))
}

// ArticleUpdateOutput ...
type ArticleUpdateOutput struct {
	Article          *model.Article
	Message          string
	ValidationErrors []string
}

// ArticleUpdate ...
func ArticleUpdate(c echo.Context) error {
	// リスエスト送信元のパスを取得
	ref := c.Request().Referer()
	// リスエスト送信元のパスから記事IDを抽出する
	refID := strings.Split(ref, "/")[4]

	// リクエストURLのパスパラメータから記事IDを抽出する
	reqID := c.Param("articleID")

	// 編集画面で表示している記事と更新しようとしている記事が異なる
	if reqID != refID {
		return c.JSON(http.StatusBadRequest, "")
	}

	var article model.Article
	var out ArticleUpdateOutput

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&article); err != nil {
		out.ValidationErrors = article.ValidationErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	articleID, _ := strconv.Atoi(reqID)
	article.ID = articleID

	_, err := repository.ArticleUpdate(&article)

	if err != nil {
		out.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, out)
	}

	out.Article = &article

	return c.JSON(http.StatusOK, out)
}
