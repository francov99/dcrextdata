// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/types"

	"github.com/raedahgroup/dcrextdata/vsp"
)

var (
	vspTickExistsErr = fmt.Errorf("VSPTick exists")
)

// StoreVSPs attempts to store the vsp responses by calling storeVspResponseG and returning
// a slice of errors
func (pg *PgDb) StoreVSPs(ctx context.Context, data vsp.Response) []error {
	if ctx.Err() != nil {
		return []error{ctx.Err()}
	}
	errs := make([]error, 0, len(data))
	completed := 0
	for name, tick := range data {
		err := pg.storeVspResponse(ctx, name, tick)
		if err == nil {
			completed++
		} else if err != vspTickExistsErr {
			log.Trace(err)
			errs = append(errs, err)
		}
		if ctx.Err() != nil {
			return append(errs, ctx.Err())
		}
	}
	if completed == 0 {
		log.Info("Unable to store any vsp entry")
	}
	return errs
}

func (pg *PgDb) storeVspResponse(ctx context.Context, name string, resp *vsp.ResposeData) error {
	txr, err := pg.db.Begin()
	if err != nil {
		return err
	}

	pool, err := models.VSPS(models.VSPWhere.Name.EQ(name)).One(ctx, pg.db)
	if err == sql.ErrNoRows {
		pool = responseToVSP(name, resp)
		err := pg.tryInsert(ctx, txr, pool)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	vspTick := responseToVSPTick(pool.ID, resp)
	tickTime := time.Unix(int64(resp.LastUpdated), 0)

	err = vspTick.Insert(ctx, pg.db, boil.Infer())
	// if err != nil && strings.Contains(err.Error(), "unique constraint") {
	// 	log.Tracef("Tick exits for %s", name)
	// 	err = txr.Rollback()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return vspTickExistsErr
	// } else if err != nil {
	// 	txr.Rollback()
	// 	return err
	// }
	if err != nil {
		errR := txr.Rollback()
		if errR != nil {
			return err
		}
		if strings.Contains(err.Error(), "unique constraint") {
			return vspTickExistsErr
		}
		return err
	}

	// vspTickTimeExits, err := models.VSPTickTimes(
	// 	models.VSPTickTimeWhere.UpdateTime.EQ(tickTime),
	// 	models.VSPTickTimeWhere.VSPTickID.EQ(vspTick.ID)).Exists(ctx, pg.db)

	// if err != nil {
	// 	txr.Rollback()
	// 	return err
	// }

	// if !vspTickTimeExits {
	// 	vtickTime := &models.VSPTickTime{
	// 		VSPTickID:  vspTick.ID,
	// 		UpdateTime: tickTime,
	// 	}

	// 	err = pg.tryInsert(ctx, txr, vtickTime)
	// 	if err != nil {
	// 		log.Debugf("Tick time %v for %d", vtickTime.UpdateTime, vtickTime.VSPTickID)
	// 		return err
	// 	}
	// }

	err = txr.Commit()
	if err != nil {
		return txr.Rollback()
	}

	log.Tracef("Stored data for vsp %s at %v", name, tickTime.UTC())
	return nil
}

func responseToVSP(name string, resp *vsp.ResposeData) *models.VSP {
	return &models.VSP{
		Name:                 name,
		APIEnabled:           resp.APIEnabled,
		APIVersionsSupported: types.Int64Array(resp.APIVersionsSupported),
		Network:              resp.Network,
		URL:                  resp.URL,
		Launched:             time.Unix(resp.Launched, 0),
	}
}

func responseToVSPTick(poolID int, resp *vsp.ResposeData) *models.VSPTick {
	return &models.VSPTick{
		VSPID:            poolID,
		Immature:         resp.Immature,
		Live:             resp.Live,
		Voted:            resp.Voted,
		Missed:           resp.Missed,
		PoolFees:         resp.PoolFees,
		ProportionLive:   resp.ProportionLive,
		ProportionMissed: resp.ProportionMissed,
		UserCount:        resp.UserCount,
		UsersActive:      resp.UserCountActive,
		Time:             time.Unix(resp.LastUpdated, 0),
	}
}
