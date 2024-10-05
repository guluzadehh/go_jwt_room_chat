package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/storage"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	cli *redis.Client
}

func New(config *config.Config) (*Storage, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DefaultDB,
	})

	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Storage{
		cli: cli,
	}, nil
}

func (s *Storage) CreateRoom(name, password string, owner_id int64) (*models.Room, error) {
	const op = "storage.redis.CreateRoom"

	ctx := context.Background()

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	room := &models.Room{
		Uuid:     id.String(),
		Name:     name,
		Password: password,
		OwnerId:  owner_id,
	}

	hashKey := fmt.Sprintf("room:%s", room.Uuid)
	err = s.cli.HSet(ctx, hashKey, map[string]interface{}{
		"name":     room.Name,
		"password": room.Password,
		"owner_id": room.OwnerId,
	}).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return room, nil
}

func (s *Storage) Rooms() ([]*models.Room, error) {
	const op = "storage.redist.Rooms"

	ctx := context.Background()

	keys, err := s.cli.Keys(ctx, "room:*").Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var rooms []*models.Room
	for _, key := range keys {
		roomData, err := s.cli.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		owner_id, err := strconv.ParseInt(roomData["owner_id"], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		room := &models.Room{
			Uuid:     parseRoomUuid(key),
			Name:     roomData["name"],
			Password: roomData["password"],
			OwnerId:  owner_id,
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *Storage) RoomByUuid(uuid string) (*models.Room, error) {
	const op = "storage.redis.RoomByUuid"

	ctx := context.Background()

	roomData, err := s.cli.HGetAll(ctx, fmt.Sprintf("room:%s", uuid)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, storage.RoomNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	owner_id, err := strconv.ParseInt(roomData["owner_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Room{
		Uuid:     roomData["uuid"],
		Name:     roomData["name"],
		Password: roomData["password"],
		OwnerId:  owner_id,
	}, nil
}

func (s *Storage) DeleteRoom(uuid string) error {
	const op = "storage.redis.DeleteRoom"

	ctx := context.Background()
	res, err := s.cli.Del(ctx, fmt.Sprintf("room:%s", uuid)).Result()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res == 0 {
		return storage.RoomNotFound
	}
	return nil
}

func parseRoomUuid(key string) string {
	var uuid string
	fmt.Sscanf(key, "room:%s", &uuid)
	return uuid
}
