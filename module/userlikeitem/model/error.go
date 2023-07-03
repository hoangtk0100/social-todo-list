package model

import "errors"

var (
	ErrCannotLikeItem       = errors.New("cannot like this item")
	ErrCannotUnlikeItem     = errors.New("cannot unlike this item")
	ErrDidNotLikeItem       = errors.New("you have not liked this item")
	ErrItemIDInvalid        = errors.New("invalid TODO id")
	ErrCannotListLikedUsers = errors.New("cannot list users liked item")
)
