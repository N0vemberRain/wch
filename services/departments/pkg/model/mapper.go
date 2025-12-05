package model

import (
	"wch/gen"
)

func DepToProto(d *Department) *gen.Department {
	return &gen.Department{
		Id:   int32(d.ID),
		Name: d.Name,
	}
}

func DepFromProto(d *gen.Department) *Department {
	return &Department{
		ID:   int(d.Id),
		Name: d.Name,
	}
}
