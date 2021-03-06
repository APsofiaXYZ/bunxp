package xp

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
)

type Explorer struct {
	bundle *bundle.Bundle
}

func NewExplorer(b *bundle.Bundle) *Explorer {
	return &Explorer{bundle: b}
}

func (e *Explorer) Root() (Object, error) {
	t := MustGetObjectType("cluster")
	cluster, err := t.New(e.bundle, "")
	if err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (e *Explorer) Object(n ObjectTypeName, id ObjectId) (Object, error) {
	t, err := GetObjectType(n)
	if err != nil {
		return nil, fmt.Errorf("cannot get object type: %s", err.Error())
	}
	return t.New(e.bundle, id)
}
