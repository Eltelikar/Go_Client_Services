package in_memory

import "client-services/internal/graph/model"

type inMemStorage struct {
	posts    []*model.Post
	comments []*model.Comment
}

//TODO: конкуретный доступ к бд

// TODO: путь для in-memory storage
func NewStorage() {

}

//TODO: Сохранение поста
//TODO: Сохранение комментария

//TODO: Получение всех постов Posts
//TODO: Получение поста по ID Post

//TODO: Сохранить комментарий
//TODO: Получить комментарий
