package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

// хранилище где данные складываются в БД PostgreSQL.
type SrvRepository struct {
	con        *sql.DB
	logger     *zap.SugaredLogger
	commonRepo commonRep.CommonRepositoryInterface
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *zap.SugaredLogger, common commonRep.CommonRepositoryInterface) (*SrvRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &SrvRepository{con: db, logger: logger, commonRepo: common}, nil
}

func (repo *SrvRepository) GetCommonRepo() commonRep.CommonRepositoryInterface {
	return repo.commonRepo
}

func (repo *SrvRepository) Source() any {
	return repo.con
}

func (repo *SrvRepository) CreateUser(ctx context.Context, user models.UserRepo) (int64, error) {
	sql := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var lastInsertId int64

	err := repo.con.QueryRowContext(ctx, sql, user.Email, user.PasswordHash).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (repo *SrvRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserRepo, error) {
	sql := "SELECT id,email,password, created_at FROM users WHERE email=$1"

	user := models.UserRepo{}

	err := repo.con.QueryRow(sql, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *SrvRepository) CreateItem(ctx context.Context, item models.ItemRepo) (int64, error) {
	sql := "INSERT INTO items (id, type ,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := repo.con.QueryRowContext(ctx, sql, item.ID, item.Type, item.UserID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (int64, error) {
	sql := "INSERT INTO items (id, type ,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := tx.QueryRowContext(ctx, sql, item.ID, item.Type, item.UserID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) CreateItemPayloadTx(ctx context.Context, tx *sql.Tx, item models.ItemPayloadRepo) error {
	sql := "INSERT INTO item_payloads (item_id,ciphertext) VALUES ($1,$2)"

	_, err := tx.ExecContext(ctx, sql, item.ItemID, item.Ciphertext)

	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) CreateSyncChangesTx(ctx context.Context, tx *sql.Tx, itemID int64, op enum.OpType, userID int64) (int64, error) {
	sql := "INSERT INTO sync_changes (item_id,operation,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := tx.QueryRowContext(ctx, sql, itemID, op, userID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) GetSyncChangesByID(ctx context.Context, ID int64) (*models.SyncChangesRepo, error) {
	sql := "SELECT id,item_id,revision,operation,created_at FROM sync_changes WHERE id=$1"

	var change models.SyncChangesRepo

	err := repo.con.QueryRowContext(ctx, sql, ID).Scan(&change.ID, &change.ID, &change.ItemID, &change.Revision, &change.Operation, &change.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &change, nil
}

func (repo *SrvRepository) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	sql := "SELECT MAX(revision) FROM sync_changes WHERE user_id=$1 GROUP BY user_id"

	var rev int64

	err := repo.con.QueryRow(sql, userID).Scan(&rev)

	if err != nil {
		return 0, err
	}

	return rev, nil
}

func (repo *SrvRepository) GetUserSyncChanges(ctx context.Context, userID int64, since int64) ([]models.SyncChangesRepo, error) {
	sql := "SELECT id, item_id, revision, operation, created_at FROM sync_changes WHERE user_id=$1 and revision > $2;"

	rows, err := repo.con.QueryContext(ctx, sql, userID, since)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
		}
	}()

	var items []models.SyncChangesRepo

	for rows.Next() {
		var item models.SyncChangesRepo

		if err := rows.Scan(
			&item.ID,
			&item.ItemID,
			&item.Revision,
			&item.Operation,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SrvRepository) GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error) {
	sql := "SELECT id, type, created_at FROM items WHERE id = ANY($1);"

	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
		}
	}()

	var items []models.ItemRepo

	for rows.Next() {
		var item models.ItemRepo

		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SrvRepository) GetPayloadByItemIDs(ctx context.Context, itemIDs []int64) ([]models.ItemPayloadRepo, error) {
	sql := "SELECT item_id, ciphertext FROM item_payloads WHERE item_id = ANY($1);"

	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
		}
	}()

	var items []models.ItemPayloadRepo

	for rows.Next() {
		var item models.ItemPayloadRepo

		if err := rows.Scan(
			&item.ItemID,
			&item.Ciphertext,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SrvRepository) DeletePayloadByItemIDTx(ctx context.Context, tx *sql.Tx, itemID int64) error {
	sql := "DELETE FROM item_payloads WHERE item_id=$1"

	_, err := tx.ExecContext(ctx, sql, itemID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) UpdateUserItemPayloadTx(ctx context.Context, tx *sql.Tx, itemID int64, userID int64, val string) error {
	sql := "UPDATE item_payloads p SET ciphertext=$1 FROM items i WHERE p.item_id=$2 AND i.user_id=$3"

	_, err := tx.ExecContext(ctx, sql, val, itemID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) ApplySync(ctx context.Context, in []shrModel.SyncPutChange, userID int64) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		// todo error
		_ = tx.Rollback()
	}()

	for _, item := range in {

		repo.logger.Debugf("Sync Operation: %s", item.Operation)

		line := models.ItemRepo{
			ID:     int64(item.Item.ID),
			UserID: userID,
			Type:   item.Item.Type,
		}

		switch item.Operation {
		case enum.OpCreate:
			repo.logger.Debug("create strategy")
			itemID, err := repo.CreateItemTx(ctx, tx, line)

			if err != nil {
				return err
			}

			pl := models.ItemPayloadRepo{
				ItemID:     itemID,
				Ciphertext: item.Item.Ciphertext,
			}

			if err := repo.CreateItemPayloadTx(ctx, tx, pl); err != nil {
				return err
			}

			for k, v := range item.Metadata {
				if _, err := repo.GetCommonRepo().CreateMetaTx(ctx, tx, itemID, k, v); err != nil {
					return err
				}
			}

			if _, err := repo.CreateSyncChangesTx(ctx, tx, itemID, item.Operation, userID); err != nil {
				return err
			}
		case enum.OpDelete:
			repo.logger.Debug("Delete strategy")
			if err := repo.GetCommonRepo().DeleteUserMetaByItemIDTx(ctx, tx, int64(item.Item.ID), userID); err != nil {
				return err
			}

			if err := repo.DeletePayloadByItemIDTx(ctx, tx, int64(item.Item.ID)); err != nil {
				return err
			}

			if err := repo.GetCommonRepo().DeleteUserItemByIDTx(ctx, tx, int64(item.Item.ID), userID); err != nil {
				return err
			}

			// todo op to const
			if _, err := repo.CreateSyncChangesTx(ctx, tx, int64(item.Item.ID), enum.OpDelete, userID); err != nil {
				return err
			}
		case enum.OpUpdate:
			repo.logger.Debug("Update strategy")
			if err := repo.UpdateUserItemPayloadTx(ctx, tx, int64(item.Item.ID), userID, item.Item.Ciphertext); err != nil {
				return err
			}

			// todo update metadata

			if _, err := repo.CreateSyncChangesTx(ctx, tx, int64(item.Item.ID), enum.OpUpdate, userID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown operation: %s", item.Operation)
		}

		// todo op co conts

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
