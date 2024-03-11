package reaction

type Reaction struct {
	ID int64
	UserId int64
	PostId int64
	Status int
}

type ReactionComment struct {
	ID int64
	UserId int64
	CommentId int64
	Status int
}