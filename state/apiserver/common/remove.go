// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"fmt"

	"launchpad.net/juju-core/state"
	"launchpad.net/juju-core/state/api/params"
)

// Remover implements a common Remove method for use by various facades.
type Remover struct {
	st           state.EntityFinder
	getCanModify GetAuthFunc
}

// NewRemover returns a new Remover. The GetAuthFunc will be used on
// each invocation of Remove to determine current permissions.
func NewRemover(st state.EntityFinder, getCanModify GetAuthFunc) *Remover {
	return &Remover{
		st:           st,
		getCanModify: getCanModify,
	}
}

func (r *Remover) removeEntity(entityp params.Entity) error {
	entity, err := r.st.FindEntity(entityp.Tag)
	if err != nil {
		return err
	}
	remover, ok := entity.(interface {
		state.Lifer
		state.Remover
	})
	if !ok {
		return NotSupportedError(entityp.Tag, "removal")
	}
	// Only remove entites that are not Alive.
	if life := remover.Life(); life == state.Alive {
		return fmt.Errorf("cannot remove entity %q: still alive", entityp.Tag)
	}
	if err = remover.EnsureDead(); err != nil {
		return err
	}
	return remover.Remove()
}

// Remove removes every given entity from state, calling EnsureDead
// first, then Remove. It will fail if the entity is not present.
func (r *Remover) Remove(args params.Entities) (params.ErrorResults, error) {
	result := params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.Entities)),
	}
	if len(args.Entities) == 0 {
		return result, nil
	}
	canModify, err := r.getCanModify()
	if err != nil {
		return params.ErrorResults{}, err
	}
	for i, entity := range args.Entities {
		err := ErrPerm
		if canModify(entity.Tag) {
			err = r.removeEntity(entity)
		}
		result.Results[i].Error = ServerError(err)
	}
	return result, nil
}
