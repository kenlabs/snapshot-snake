// Code generated by github.com/londobell/tool/genapi. DO NOT EDIT.

package api

import (
	"context"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

var ErrNotSupported = xerrors.New("method not supported")

type SnapAPIStruct struct {
	Internal struct {
		ChainGetTipSet func(p0 context.Context, p1 types.TipSetKey) (*types.TipSet, error) ``

		GetDagNode func() ([]cid.Cid, error) ``

		SnapDagExport func(p0 context.Context, p1 *types.TipSet, p2 int64) (<-chan []byte, error) ``
	}
}

type SnapAPIStub struct {
}

func (s *SnapAPIStruct) ChainGetTipSet(p0 context.Context, p1 types.TipSetKey) (*types.TipSet, error) {
	if s.Internal.ChainGetTipSet == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetTipSet(p0, p1)
}

func (s *SnapAPIStub) ChainGetTipSet(p0 context.Context, p1 types.TipSetKey) (*types.TipSet, error) {
	return nil, ErrNotSupported
}

func (s *SnapAPIStruct) GetDagNode() ([]cid.Cid, error) {
	if s.Internal.GetDagNode == nil {
		return *new([]cid.Cid), ErrNotSupported
	}
	return s.Internal.GetDagNode()
}

func (s *SnapAPIStub) GetDagNode() ([]cid.Cid, error) {
	return *new([]cid.Cid), ErrNotSupported
}

func (s *SnapAPIStruct) SnapDagExport(p0 context.Context, p1 *types.TipSet, p2 int64) (<-chan []byte, error) {
	if s.Internal.SnapDagExport == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.SnapDagExport(p0, p1, p2)
}

func (s *SnapAPIStub) SnapDagExport(p0 context.Context, p1 *types.TipSet, p2 int64) (<-chan []byte, error) {
	return nil, ErrNotSupported
}

var _ SnapAPI = new(SnapAPIStruct)
