// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storagecommon_test

import (
	"github.com/juju/errors"
	"github.com/juju/names"
	"github.com/juju/testing"

	"github.com/juju/juju/apiserver/common/storagecommon"
	"github.com/juju/juju/state"
	"github.com/juju/juju/storage"
	"github.com/juju/juju/storage/poolmanager"
)

type fakeStorage struct {
	testing.Stub
	storagecommon.StorageInterface
	storageInstance       func(names.StorageTag) (state.StorageInstance, error)
	storageInstanceVolume func(names.StorageTag) (state.Volume, error)
	volumeAttachment      func(names.MachineTag, names.VolumeTag) (state.VolumeAttachment, error)
	blockDevices          func(names.MachineTag) ([]state.BlockDeviceInfo, error)
}

func (s *fakeStorage) StorageInstance(tag names.StorageTag) (state.StorageInstance, error) {
	s.MethodCall(s, "StorageInstance", tag)
	return s.storageInstance(tag)
}

func (s *fakeStorage) StorageInstanceVolume(tag names.StorageTag) (state.Volume, error) {
	s.MethodCall(s, "StorageInstanceVolume", tag)
	return s.storageInstanceVolume(tag)
}

func (s *fakeStorage) VolumeAttachment(m names.MachineTag, v names.VolumeTag) (state.VolumeAttachment, error) {
	s.MethodCall(s, "VolumeAttachment", m, v)
	return s.volumeAttachment(m, v)
}

func (s *fakeStorage) BlockDevices(m names.MachineTag) ([]state.BlockDeviceInfo, error) {
	s.MethodCall(s, "BlockDevices", m)
	return s.blockDevices(m)
}

type fakeStorageInstance struct {
	state.StorageInstance
	tag   names.StorageTag
	owner names.Tag
	kind  state.StorageKind
}

func (i *fakeStorageInstance) StorageTag() names.StorageTag {
	return i.tag
}

func (i *fakeStorageInstance) Tag() names.Tag {
	return i.tag
}

func (i *fakeStorageInstance) Owner() names.Tag {
	return i.owner
}

func (i *fakeStorageInstance) Kind() state.StorageKind {
	return i.kind
}

type fakeStorageAttachment struct {
	state.StorageAttachment
	storageTag names.StorageTag
}

func (a *fakeStorageAttachment) StorageInstance() names.StorageTag {
	return a.storageTag
}

type fakeVolume struct {
	state.Volume
	tag    names.VolumeTag
	params *state.VolumeParams
	info   *state.VolumeInfo
}

func (v *fakeVolume) VolumeTag() names.VolumeTag {
	return v.tag
}

func (v *fakeVolume) Tag() names.Tag {
	return v.tag
}

func (v *fakeVolume) Params() (state.VolumeParams, bool) {
	if v.params == nil {
		return state.VolumeParams{}, false
	}
	return *v.params, true
}

func (v *fakeVolume) Info() (state.VolumeInfo, error) {
	if v.info == nil {
		return state.VolumeInfo{}, errors.NotProvisionedf("volume %v", v.tag.Id())
	}
	return *v.info, nil
}

type fakeVolumeAttachment struct {
	state.VolumeAttachment
	info *state.VolumeAttachmentInfo
}

func (v *fakeVolumeAttachment) Info() (state.VolumeAttachmentInfo, error) {
	if v.info == nil {
		return state.VolumeAttachmentInfo{}, errors.NotProvisionedf("volume attachment")
	}
	return *v.info, nil
}

type fakePoolManager struct {
	poolmanager.PoolManager
}

func (pm *fakePoolManager) Get(name string) (*storage.Config, error) {
	return nil, errors.NotFoundf("pool")
}