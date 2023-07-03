package model

import "errors"

var (
	ErrTitleEmpty          = errors.New("title cannot be empty")
	ErrItemDeleted         = errors.New("item is deleted")
	ErrCannotCreateItem    = errors.New("cannot create new TODO item")
	ErrCannotGetItemLikes  = errors.New("cannot get item likes")
	ErrCannotGetItem       = errors.New("cannot get TODO item")
	ErrCannotGetItems      = errors.New("cannot get TODO items")
	ErrCannotUpdateItem    = errors.New("cannot update TODO item")
	ErrCannotDeleteItem    = errors.New("cannot delete TODO item")
	ErrRequesterIsNotOwner = errors.New("no permission, only TODO owner can do this")
	ErrItemIDInvalid       = errors.New("invalid TODO id")
)
