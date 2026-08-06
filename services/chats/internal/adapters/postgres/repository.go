package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	chats "wch/services/chats/internal/domain"
	"wch/services/chats/internal/domain/model"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var ErrHasntBeenDone error = errors.New("the method hasn't been done yet")

type ChatRepositoryPg struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepositoryPg {
	return &ChatRepositoryPg{db}
}

func (r *ChatRepositoryPg) CreateChat(ctx context.Context, c *model.Chat) error {
	if row := r.db.QueryRow(
		`INSERT INTO chats (id, type, name, created_at, updated_at) VALUES (
		$1, $2, $3, $4, $5
		)`,
		c.ID,
		model.ChatTypeToString(c.Type),
		c.Name,
		c.CreatedAt,
		c.UpdatedAt,
	); row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (r *ChatRepositoryPg) GetChatByID(ctx context.Context, chatID uuid.UUID) (*model.Chat, error) {
	log.Println("ChatRepositoryPg.GetChatByID: entering...")
	c := &model.Chat{}
	var chat_type string
	err := r.db.QueryRow(
		`SELECT id, type, name, created_at, updated_at FROM chats WHERE id=$1`,
		chatID.String(),
	).Scan(
		&c.ID,
		&chat_type,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, chats.ErrChatNotFound
		}
		return nil, err
	}

	c.Type, err = model.ChatTypeFromString(chat_type)
	if err != nil {
		return nil, err
	}

	log.Println("ChatRepositoryPg.GetChatByID: exiting...")
	return c, nil
}

func (r *ChatRepositoryPg) GetChatByName(ctx context.Context, name string) (*model.Chat, error) {
	c := &model.Chat{}
	var chat_type string
	err := r.db.QueryRow(
		`SELECT id, type, name, created_at, updated_at FROM chats WHERE name=$1`,
		name,
	).Scan(
		&c.ID,
		&chat_type,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.Type, err = model.ChatTypeFromString(chat_type)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatRepositoryPg) UpdateChat(ctx context.Context, c *model.Chat) error {
	if err := r.db.QueryRow(
		`UPDATE chats SET
		name = $1,
		updated_at = $2 WHERE id = $3;`,
		c.Name,
		c.UpdatedAt,
		c.ID,
	); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (r *ChatRepositoryPg) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	return ErrHasntBeenDone
}

func (r *ChatRepositoryPg) AddParticipant(
	ctx context.Context, chatID uuid.UUID, p *model.ChatParticipant,
) error {
	err := r.db.QueryRow(
		`INSERT INTO chat_participants (chat_id, user_id, role) VALUES ($1, $2, $3);`,
		p.ChatID,
		p.UserID,
		model.ChatParticipantRoleToString(p.Role),
	).Err()
	if err != nil {
		log.Printf("not ok")
		return err
	}

	return nil
}

func (r *ChatRepositoryPg) RemoveParticipant(
	ctx context.Context, chatID uuid.UUID, userID uuid.UUID,
) error {
	err := r.db.QueryRow(
		`DELETE FROM chat_participants WHERE chat_id=$1, user_id=$2;`,
		chatID.String(),
		userID.String(),
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepositoryPg) GetParticipantsIDs(
	ctx context.Context, chatID uuid.UUID,
) ([]uuid.UUID, error) {
	rows, err := r.db.Query(
		`SELECT user_id FROM chat_participants WHERE chat_id=$1;`,
		chatID.String(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids = make([]uuid.UUID, 0)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uuid.MustParse(id))
	}

	return ids, nil
}

func (r *ChatRepositoryPg) GetParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (
	*model.ChatParticipant,
	error,
) {
	row := r.db.QueryRow(
		`SELECT * FROM chat_participants WHERE chat_id=$1 AND user_id=$2;`,
		chatID.String(),
		userID,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	p := &model.ChatParticipant{}
	var role string
	if err := row.Scan(&p.ChatID, &p.UserID, &role, &p.JoinedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, chats.ErrParticipateNotFound
		} else {
			return nil, err
		}
	}
	p.Role = model.ChatParticipantRoleFromString(role)

	return p, nil
}

func (r *ChatRepositoryPg) GetChatsForUser(ctx context.Context, userID uuid.UUID) (
	[]model.Chat,
	error,
) {
	chats_ids_rows, err := r.db.QueryContext(
		ctx,
		`SELECT chat_id FROM chat_participants WHERE user_id = $1;`,
		userID.String(),
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, chats.ErrChatNotFound
		} else {
			return nil, err
		}
	}
	defer chats_ids_rows.Close()

	var ids []string
	for chats_ids_rows.Next() {
		var id string
		err = chats_ids_rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	chats_data_rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, type, name, created_at, updated_at FROM chats WHERE id = ANY($1)`,
		pq.Array(ids),
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, chats.ErrChatNotFound
		} else {
			return nil, err
		}
	}
	defer chats_data_rows.Close()

	var chats []model.Chat
	for chats_data_rows.Next() {
		var c model.Chat
		var chat_type string
		err = chats_data_rows.Scan(
			&c.ID,
			&chat_type,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		t, err := model.ChatTypeFromString(chat_type)
		if err != nil {
			return nil, err
		}

		c.Type = t
		chats = append(chats, c)
	}

	return chats, nil
}
