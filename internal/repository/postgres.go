package repository

import (
	"context"
	"errors"

	"github.com/aakash-1857/codebin/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
type SnippetRepository struct{
	DB *pgxpool.Pool
}
func NewSnippetRepository(db *pgxpool.Pool) *SnippetRepository{
	return &SnippetRepository{
		DB: db,
	}
}

func (r *SnippetRepository) Insert(title,content string) (string,error){
	stmt:= `INSERT INTO snippets (title, content, created_at, expires_at)
    VALUES($1, $2, NOW(), NOW() + INTERVAL '365 days')
    RETURNING id`
	var id string
	err:=r.DB.QueryRow(context.Background(),stmt,title,content).Scan(&id)
	if err!=nil{
		return "",err
	}
	return id,nil
}

func (r *SnippetRepository) Get(id string) (*models.Snippet,error){
	stmt:=`SELECT id, title, content, created_at, expires_at FROM snippets
	WHERE id = $1 AND expires_at > NOW()`
	row:=r.DB.QueryRow(context.Background(),stmt,id)
	s:=&models.Snippet{}//initialize a pointer to a new snippet struct
	err:=row.Scan(&s.ID,&s.Title,&s.Content,&s.CreatedAt,&s.ExpiresAt)
	if err!=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,models.ErrNoRecord
		}else{
			return nil,err
		}
	}
	return s,nil
}

func (r *SnippetRepository) Latest() ([]*models.Snippet,error){
	// SQL
	stmt := `SELECT id, title, content, created_at, expires_at FROM snippets
	WHERE expires_at > NOW() ORDER BY created_at DESC LIMIT 10`
	//multiple rows => Query()
	rows,err:=r.DB.Query(context.Background(),stmt)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()
	snippets := []*models.Snippet{}
	for rows.Next(){
		s:=&models.Snippet{}
		err:=rows.Scan(&s.ID,&s.Title,&s.Content,&s.CreatedAt,&s.ExpiresAt)
		if err!=nil{
			return nil,err
		}
		snippets=append(snippets,s)
	}
	if err=rows.Err();err!=nil{return nil,err}
	return snippets,nil
}