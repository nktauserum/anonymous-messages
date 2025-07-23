package article

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, article Article) error
	ReadArticle(ctx context.Context, articleUUID string) (Article, error)
	UpdateArticle(ctx context.Context, articleUUID string, newText string) error
	DeleteArticle(ctx context.Context, articleUUID string) error
	ListArticles(ctx context.Context) ([]Article, error)
}

type ArticleStorage struct {
	db *sql.DB
}

func NewArticleStorage(path string) (*ArticleStorage, error) {
	instance, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = instance.Ping()
	if err != nil {
		return nil, err
	}

	db := &ArticleStorage{db: instance}
	err = db.initTables(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *ArticleStorage) initTables(ctx context.Context) error {
	const sqlStmt = `
    CREATE TABLE IF NOT EXISTS articles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL UNIQUE,
        title TEXT NOT NULL,
        body TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := s.db.ExecContext(ctx, sqlStmt)
	return err
}

func (s *ArticleStorage) CreateArticle(ctx context.Context, article Article) error {
	_, err := s.db.ExecContext(ctx, `
		insert into articles (uuid, title, body, created_at) values (?, ?, ?, ?)
	`, article.UUID, article.Title, article.Text, article.DatePublished)
	return err
}

func (s *ArticleStorage) ReadArticle(ctx context.Context, articleUUID string) (Article, error) {
	var article Article

	row := s.db.QueryRowContext(ctx, `
		select title, body, created_at from articles where uuid = $1
 	`, articleUUID)
	err := row.Scan(&article.Title, &article.Text, &article.DatePublished)
	if err != nil {
		return article, err
	}

	return article, nil
}

func (s *ArticleStorage) UpdateArticle(ctx context.Context, articleUUID, newText string) error {
	_, err := s.db.ExecContext(context.Background(), `
		update articles
		set body = $1
		where uuid = $2
	`, newText, articleUUID)

	return err
}

func (s *ArticleStorage) DeleteArticle(ctx context.Context, articleUUID string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM articles WHERE uuid = ?", articleUUID)
	return err
}

func (s *ArticleStorage) ListArticles(ctx context.Context) ([]Article, error) {
	rows, err := s.db.QueryContext(ctx, `
		select uuid, title, body, created_at from articles order by created_at desc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.UUID, &a.Title, &a.Text, &a.DatePublished); err != nil {
			return nil, err
		}

		articles = append(articles, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
