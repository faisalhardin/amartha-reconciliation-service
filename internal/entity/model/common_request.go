package model

type CommonRequestParam struct {
	Limit  int `schema:"limit" validate:"omitempty,min=1,max=100"`
	Offset int `schema:"offset" validate:"omitempty,min=0"`
}
