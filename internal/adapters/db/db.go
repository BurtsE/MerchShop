package db

import (
	"MerchShop/internal/application/core/domain"
	"MerchShop/internal/ports"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

const DefaultCoinsValue = 1000

var _ ports.DBPort = (*DBAdapter)(nil)

type DBAdapter struct {
	db *pgx.Conn
}

func NewDBAdapter(source string) (*DBAdapter, error) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}
	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}
	return &DBAdapter{db: db}, nil
}

func (a DBAdapter) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	query := `
		INSERT INTO users (username, password_hash, coins)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err := a.db.QueryRow(ctx, query, user.Username, user.PasswordHash, DefaultCoinsValue).Scan(&user.ID)
	if err != nil {
		return domain.User{}, fmt.Errorf("creating user: %w", err)
	}
	return user, nil
}

func (a DBAdapter) User(ctx context.Context, userID uint) (domain.User, error) {
	user := domain.User{}
	query := `
		SELECT username, password_hash, coins
		FROM users 
			WHERE id = $1;
	`
	err := a.db.QueryRow(ctx, query, userID).Scan(&user.Username, &user.PasswordHash, &user.Coins)
	if err != nil {
		return domain.User{}, fmt.Errorf("getting user: %w", err)
	}
	user.ID = userID
	return user, nil
}

func (a DBAdapter) UserByName(ctx context.Context, username string) (domain.User, error) {
	user := domain.User{}
	query := `
		SELECT id, password_hash, coins
		FROM users 
			WHERE username = $1;
	`
	err := a.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.PasswordHash, &user.Coins)
	if err != nil {
		return domain.User{}, fmt.Errorf("getting user: %w", err)
	}
	user.Username = username
	return user, nil
}

func (a DBAdapter) UpdateUser(ctx context.Context, user domain.User) error {
	query := `
		UPDATE users username, password_hash, coins
		FROM users 
		SET username = $2, password_hash = $3, coins = $4
		    WHERE id = $1;
	`
	_, err := a.db.Exec(ctx, query, user.ID, user.Username, user.PasswordHash, user.Coins)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

func (a DBAdapter) UserWallet(ctx context.Context, user domain.User) ([]domain.WalletOperation, error) {
	var wallet = make([]domain.WalletOperation, 0)

	query := `
		SELECT wallet_operations.id, receiver_id, "value", username, password_hash, coins
		FROM wallet_operations LEFT JOIN users on(receiver_id=users.id)
			WHERE sender_id = $1;
	`
	rows, err := a.db.Query(ctx, query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting sent coins: %w", err)
	}
	for rows.Next() {
		op := domain.WalletOperation{Sender: user}
		err = rows.Scan(&op.ID, &op.Receiver.ID, &op.Value, &op.Receiver.Username, &op.Receiver.PasswordHash, &op.Receiver.Coins)
		if err != nil {
			return nil, fmt.Errorf("scanning sent coins: %w", err)
		}
		wallet = append(wallet, op)
	}

	query = `
		SELECT wallet_operations.id, sender_id, "value", username, password_hash, coins
		FROM wallet_operations LEFT JOIN users on(sender_id=users.id)
			WHERE receiver_id = $1;
	`
	rows, err = a.db.Query(ctx, query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting received coins: %w", err)
	}

	for rows.Next() {
		op := domain.WalletOperation{Sender: user}
		err = rows.Scan(&op.ID, &op.Receiver.ID, &op.Value, &op.Receiver.Username, &op.Receiver.PasswordHash, &op.Receiver.Coins)
		if err != nil {
			return nil, fmt.Errorf("scanning received coins: %w", err)
		}
		wallet = append(wallet, op)
	}
	return wallet, nil
}

func (a DBAdapter) UserInventory(ctx context.Context, user domain.User) (domain.Inventory, error) {
	var inventory = make(domain.Inventory, 0)
	query := `
		SELECT item_name, amount
		FROM inventory 
			WHERE user_id = $1;
	`
	rows, err := a.db.Query(ctx, query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting user inventory: %w", err)
	}

	for rows.Next() {
		items := domain.Items{}
		err = rows.Scan(&items.Type, &items.Quantity)
		if err != nil {
			return nil, fmt.Errorf("scanning user inventory: %w", err)
		}
		inventory = append(inventory, items)
	}
	return inventory, nil
}

func (a DBAdapter) BuyItem(ctx context.Context, user domain.User, item string) (uint, error) {
	var (
		inventoryID uint
		itemCost    int
	)
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT cost
		FROM Merch 
			WHERE type = $1;
	`
	err = tx.QueryRow(ctx, query, item).Scan(&itemCost)
	if err != nil {
		return 0, fmt.Errorf("getting item: %w", err)
	}
	if user.Coins < itemCost {
		return 0, fmt.Errorf("buying item: not enough coins")
	}
	query = `
		INSERT INTO inventory (user_id, item_name)
		VALUES ($1, $2)
		ON CONFLICT (user_id, item_name) DO UPDATE 
  			SET amount = excluded.amount + 1
		RETURNING id;
	`
	err = tx.QueryRow(ctx, query, user.ID, item).Scan(&inventoryID)
	if err != nil {
		return 0, fmt.Errorf("updating inventory: %w", err)
	}
	query = `
		UPDATE users 
			SET coins = coins - $2
			WHERE id = $1;
	`
	_, err = tx.Exec(ctx, query, user.ID, itemCost)
	if err != nil {
		return 0, fmt.Errorf("updating user coins: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}
	return inventoryID, nil
}

func (a DBAdapter) SendCoins(ctx context.Context, from domain.User, to domain.User, amount int) (uint, error) {
	var operationID uint
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	query := `
		INSERT INTO wallet_operations (sender_id, receiver_id, value)
		VALUES ($1, $2, $3)
		RETURNING id;
	`
	err = tx.QueryRow(ctx, query, from.ID, to.ID, amount).Scan(&operationID)
	if err != nil {
		return 0, fmt.Errorf("inserting operation: %w", err)
	}
	query = `
		UPDATE users 
			SET coins = coins - $2
			WHERE id = $1;
	`
	_, err = tx.Exec(ctx, query, from.ID, amount)
	if err != nil {
		return 0, fmt.Errorf("updating sender coins: %w", err)
	}
	query = `
		UPDATE users 
			SET coins = coins + $2
			WHERE id = $1;
	`
	_, err = tx.Exec(ctx, query, to.ID, amount)
	if err != nil {
		return 0, fmt.Errorf("updating receiver coins: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}
	return operationID, nil
}
