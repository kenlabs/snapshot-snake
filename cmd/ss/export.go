package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/snapshot_snake/api"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"time"
)

var exportCmd = &cli.Command{
	Name: "export",
	Subcommands: []*cli.Command{
		exportSnapshotCmd,
	},
}

var exportSnapshotCmd = &cli.Command{
	Name: "snapshot",
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name:  "recent-stateroots",
			Usage: "specify the number of recent state roots to include in the export",
		},
	},
	Action: func(cctx *cli.Context) error {
		apiv0, _, err := GetAPIV0(cctx)
		if err != nil {
			return fmt.Errorf("get apiv0 err: %s", err)
		}
		ctx := context.Background()

		//CreateExportFile
		fi, err := CreateExportFile(cctx.App, cctx.Args().First())
		if err != nil {
			log.Errorf("create export file err: %s", err)
			return err
		}

		ts, err := LoadTipSet(ctx, apiv0)
		if err != nil {
			fmt.Println(err)
			return err
		}

		rs := cctx.Int64("recent-stateroots")

		begin := time.Now()
		stream, err := apiv0.SnapDagExport(ctx, ts, rs)
		if err != nil {
			return err
		}

		var last bool
		for b := range stream {
			last = len(b) == 0

			fmt.Println(len(b))

			_, err := fi.Write(b)
			if err != nil {
				return err
			}
		}

		log.Infof("done export %d tipset height elapsed %s", rs, time.Now().Sub(begin).String())

		if !last {
			return xerrors.Errorf("incomplete export (remote connection lost /  daemon process has not yet loaded the block into the cache?)")
		}

		return nil
	},
}

func LoadTipSet(ctx context.Context, api api.SnapAPI) (*types.TipSet, error) {
	nodes, _ := api.GetDagNode()
	// get from cache or build a ts
	key := types.NewTipSetKey(nodes...)
	// load tipset
	ts, err := api.ChainGetTipSet(ctx, key)
	if err != nil {
		return nil, err
	}
	return ts, nil

}
