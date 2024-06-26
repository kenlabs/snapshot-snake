package api

import (
	"bufio"
	"context"
	"fmt"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/snapshot_snake/common"
	"github.com/snapshot_snake/snapshot/saaf"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
	"io"
)

var _ SnapAPI = (*SnapNodeAPI)(nil)
var log = logging.Logger("rpc")

type SnapNodeAPI struct {
	fx.In

	Ds common.DagStore

	Src *saaf.SnapSource
}

func (f *SnapNodeAPI) ChainGetTipSet(ctx context.Context, tsk types.TipSetKey) (*types.TipSet, error) {
	// Fetch tipset block headers from blockstore in parallel
	var eg errgroup.Group
	cids := tsk.Cids()
	blks := make([]*types.BlockHeader, len(cids))
	for i, c := range cids {
		i, c := i, c
		eg.Go(func() error {
			b, err := f.Ds.Get(ctx, c)
			if err != nil {
				return xerrors.Errorf("get block %s: %w", c, err)
			}

			blk, err := types.DecodeBlock(b.RawData())
			if err != nil {
				return xerrors.Errorf("decode block err: %s", err)
			}
			blks[i] = blk
			return nil
		})
	}
	err := eg.Wait()
	if blks[0].Cid() == blks[1].Cid() {
		return nil, fmt.Errorf("common...")
	}
	if err != nil {
		return nil, err
	}

	ts, err := types.NewTipSet(blks)
	if err != nil {
		return nil, err
	}

	return ts, nil

}

func (f *SnapNodeAPI) SnapDagExport(ctx context.Context, ts *types.TipSet, n int64) (<-chan []byte, error) {
	r, w := io.Pipe()
	out := make(chan []byte)
	go func() {
		bw := bufio.NewWriterSize(w, 1<<20)

		err := f.Ds.Export(ctx, ts, bw, n)
		bw.Flush()
		w.CloseWithError(err)
	}()

	go func() {
		defer close(out)
		for {
			buf := make([]byte, 1<<20)
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				log.Errorf("chain export pipe read failed: %s", err)
				return
			}
			if n > 0 {
				select {
				case out <- buf[:n]:
				case <-ctx.Done():
					log.Warnf("export writer failed: %s", ctx.Err())
					return
				}
			}
			if err == io.EOF {
				// send empty slice to indicate correct eof
				select {
				case out <- []byte{}:
				case <-ctx.Done():
					log.Warnf("export writer failed: %s", ctx.Err())
					return
				}

				return
			}
		}
	}()

	return out, nil
}

func (f *SnapNodeAPI) GetDagNode() ([]cid.Cid, error) {
	latest := f.Src.Latest()
	return latest, nil
}

func (f *SnapNodeAPI) GetCacheRange() (int, error) {
	return f.Src.HpRange(), nil
}
