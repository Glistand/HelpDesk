package timers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const zsetKey = "sla:timers"

// Kind identifies a timer slot.
const (
	KindFRWarn       = "fr_warn"
	KindFRBreach     = "fr_breach"
	KindResolveWarn  = "resolve_warn"
	KindResolveBreach = "resolve_breach"
)

type Store struct {
	rdb    *redis.Client
	logger *slog.Logger
}

func New(rdb *redis.Client, logger *slog.Logger) *Store {
	if logger == nil {
		logger = slog.Default()
	}
	return &Store{rdb: rdb, logger: logger}
}

func member(ticketID, kind string) string {
	return ticketID + ":" + kind
}

func ParseMember(m string) (ticketID, kind string, ok bool) {
	i := strings.LastIndex(m, ":")
	if i <= 0 || i == len(m)-1 {
		return "", "", false
	}
	return m[:i], m[i+1:], true
}

func (s *Store) Schedule(ctx context.Context, ticketID string, due map[string]time.Time) error {
	pipe := s.rdb.Pipeline()
	for kind, at := range due {
		pipe.ZAdd(ctx, zsetKey, redis.Z{
			Score:  float64(at.UnixMilli()),
			Member: member(ticketID, kind),
		})
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Store) CancelTicket(ctx context.Context, ticketID string) error {
	kinds := []string{KindFRWarn, KindFRBreach, KindResolveWarn, KindResolveBreach}
	members := make([]any, 0, len(kinds))
	for _, k := range kinds {
		members = append(members, member(ticketID, k))
	}
	return s.rdb.ZRem(ctx, zsetKey, members...).Err()
}

type Due struct {
	TicketID string
	Kind     string
	DueAt    time.Time
}

func (s *Store) PopDue(ctx context.Context, now time.Time, limit int64) ([]Due, error) {
	vals, err := s.rdb.ZRangeByScoreWithScores(ctx, zsetKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(now.UnixMilli(), 10),
		Count: limit,
	}).Result()
	if err != nil {
		return nil, err
	}
	out := make([]Due, 0, len(vals))
	for _, z := range vals {
		m, ok := z.Member.(string)
		if !ok {
			continue
		}
		tid, kind, ok := ParseMember(m)
		if !ok {
			continue
		}
		// remove before processing to avoid stampede; handler is idempotent via fired_timers
		if err := s.rdb.ZRem(ctx, zsetKey, m).Err(); err != nil {
			return out, fmt.Errorf("zrem: %w", err)
		}
		out = append(out, Due{
			TicketID: tid,
			Kind:     kind,
			DueAt:    time.UnixMilli(int64(z.Score)).UTC(),
		})
	}
	return out, nil
}
