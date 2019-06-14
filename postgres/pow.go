package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/raedahgroup/dcrextdata/pow"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func (pg *PgDb) LastPowEntryTime(source string) (time int64) {
	rows := pg.db.QueryRow(LastPowEntryTime, source)
	_ = rows.Scan(&time)
	return
}

//
func (pg *PgDb) AddPowData(ctx context.Context, data []pow.PowData) error {
	added := 0
	for _, d := range data {
		powModel, err := responseToPowModel(d)
		if err != nil {
			return err
		}

		err = powModel.Insert(ctx, pg.db, boil.Infer())
		if err != nil {
			if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
				return err
			}
		}
		added++
	}
	if len(data) == 1 {
		log.Infof("%s - Added %d pow entry from %s", UnixTimeToString(data[0].Time), added, data[0].Source)
	} else {
		last := data[len(data)-1]
		log.Infof("%s to %s - Added %d pow entries from %s",
			UnixTimeToString(data[0].Time), UnixTimeToString(last.Time), added, last.Source)
	}

	return nil
}

func responseToPowModel(data pow.PowData) (models.PowDatum, error) {
	return models.PowDatum{
		BTCPrice:          null.StringFrom(fmt.Sprint(data.BtcPrice)),
		CoinPrice:         null.StringFrom(fmt.Sprint(data.CoinPrice)),
		NetworkDifficulty: null.Float64From(data.NetworkDifficulty),
		NetworkHashrate:   null.StringFrom(fmt.Sprint(data.NetworkHashrate)),
		PoolHashrate:      null.StringFrom(fmt.Sprint(data.PoolHashrate)),
		Source:            data.Source,
		Time:              int(data.Time),
		Workers:           null.IntFrom(int(data.Workers)),
	}, nil
}
