package entity

import "errors"

var (
	ErrLikedItem            = errors.New("already like this item")
	ErrUnlikedItem          = errors.New("already unlike this item")
	ErrDidNotLikeItem       = errors.New("you have not liked this item")
	ErrItemIDInvalid        = errors.New("invalid TODO id")
	ErrCannotListLikedUsers = errors.New("cannot list users liked item")
)
