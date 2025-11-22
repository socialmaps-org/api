package model

import (
	"context"
	"database/sql"
	"regexp"

	"codeberg.org/socialmaps/api/internal/database"
)

type Place struct {
	Created int64
	Updated int64

	ID      string
	Name    string
	Lat     float64
	Lon     float64
	OSMType string
	OSMID   uint64

	NLikes       uint64
	NDislikes    uint64
	DecNLikes    float64
	DecNDislikes float64
	DecUpdatedAt int64
}

const placeColumns = `
	  created
	, updated
	, id
	, name
	, lat
	, lon
	, osm_type
	, osm_id
	, n_likes
	, n_dislikes
	, dec_n_likes
	, dec_n_dislikes
	, dec_updated_at
`

var placeIDRegex = regexp.MustCompile(`^plc_[a-zA-Z0-9]+$`)

func scanPlace(scn database.Scanner) *Place {
	var plc Place
	err := scn.Scan(
		&plc.Created,
		&plc.Updated,
		&plc.ID,
		&plc.Name,
		&plc.Lat,
		&plc.Lon,
		&plc.OSMType,
		&plc.OSMID,
		&plc.NLikes,
		&plc.NDislikes,
		&plc.DecNLikes,
		&plc.DecNDislikes,
		&plc.DecUpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		panic(err)
	}
	return &plc
}

func CreatePlace(ctx context.Context, db *sql.DB, name string, lat float64, lon float64, osmType string, osmID uint64) *Place {
	id := NewRandomID("plc")

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO Places (
			  created
			, id
			, name
			, lat
			, lon
			, osm_type
			, osm_id
		) VALUES (
		 	  @created
		 	, @id
			, @name
			, @lat
			, @lon
			, @osm_type
			, @osm_id
		) RETURNING `+placeColumns+`;`,
		sql.Named("created", id.time.Unix()),
		sql.Named("id", id.String()),
		sql.Named("name", name),
		sql.Named("lat", lat),
		sql.Named("lon", lon),
		sql.Named("osm_type", osmType),
		sql.Named("osm_id", osmID),
	)

	plc := scanPlace(row)

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return plc
}

func ListPlacesByCoord(ctx context.Context, db *sql.DB, south, west, north, east float64) (places []*Place, err error) {
	rows, err := db.QueryContext(
		ctx,
		`
		SELECT `+placeColumns+`
		FROM
			Places
		WHERE
			lat BETWEEN @lat_min AND @lat_max
			AND lon BETWEEN @lon_min AND @lon_max
		;
		`,
		sql.Named("lat_min", south),
		sql.Named("lon_min", west),
		sql.Named("lat_max", north),
		sql.Named("lon_max", east),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		places = append(places, scanPlace(rows))
	}

	return
}

func LoadPlaceByID(ctx context.Context, db *sql.DB, id string) *Place {
	row := db.QueryRowContext(
		ctx,
		`
		SELECT `+placeColumns+`
		FROM
			Places
		WHERE
			id = @id
		;
		`,
		sql.Named("id", id),
	)
	return scanPlace(row)
}

func IsValidPlaceID(s string) bool {
	return placeIDRegex.MatchString(s)
}
