package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lampadovnikita/StorekeeperTask/pkg/data"
)

type PGStorage struct {
	pool *pgxpool.Pool
}

func NewPGStorage(pool *pgxpool.Pool) *PGStorage {
	return &PGStorage{pool: pool}
}

func (s *PGStorage) GetGatheringInfo(orderIDs []int) (gi []data.GatheringInfo, err error) {
	q := `SELECT r.name AS rack_name, p.name AS product_name, p.id AS product_id, o.id AS order_id, o.amount,
			  (SELECT array_agg(rks.name)
	 			  FROM storage
	 				  JOIN racks AS rks
	 					  ON storage.rack_id = rks.id
	 			  WHERE product_id = s.product_id AND
	 		  		    is_rack_primary = false
		      ) AS additional_racks
		      FROM orders AS o 
			  	  JOIN products AS p
				  	  ON o.product_id = p.id
			  	  JOIN storage AS s
				  	  ON s.product_id = p.id
			  	  JOIN racks AS r
				  	  ON s.rack_id = r.id
		  	  WHERE o.id = ANY($1) AND
		  	  	    is_rack_primary = true
		  	  ORDER BY r.id`

	rows, err := s.pool.Query(context.Background(), q, orderIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		gathInfo := data.GatheringInfo{}

		err := rows.Scan(&gathInfo.RackName, &gathInfo.ProductName, &gathInfo.ProductID,
			&gathInfo.OrderID, &gathInfo.Amount, &gathInfo.AdditionalRacks)
		if err != nil {
			return nil, err
		}

		gi = append(gi, gathInfo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return gi, nil
}
