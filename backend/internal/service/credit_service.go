package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/creditledger"
	entuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/redis/go-redis/v9"
)

const creditBalanceCachePrefix = "credit:balance:"
const creditBalanceCacheTTL = 5 * time.Minute

// ModelCreditRate represents a parsed model rate entry used for credit calculation.
type ModelCreditRate struct {
	ID              int64
	ModelPattern    string
	CreditsPer1kIn  int64
	CreditsPer1kOut int64
	Priority        int
}

// CreditBalance is the current credit state for a user.
type CreditBalance struct {
	Balance   int64
	ExpiresAt *time.Time
	PlanID    *int64
}

// CreditLedgerEntry is one row returned from the credit_ledger table.
type CreditLedgerEntry struct {
	ID           int64
	Delta        int64
	Reason       string
	RefID        *string
	BalanceAfter int64
	Model        *string
	CreatedAt    time.Time
}

// MatchModelCreditRate returns the highest-priority rate whose ModelPattern matches
// modelName using filepath.Match glob semantics. Returns nil if nothing matches.
func MatchModelCreditRate(modelName string, rates []*ModelCreditRate) *ModelCreditRate {
	var best *ModelCreditRate
	for _, r := range rates {
		matched, _ := filepath.Match(r.ModelPattern, modelName)
		if !matched {
			continue
		}
		if best == nil || r.Priority > best.Priority {
			best = r
		}
	}
	return best
}

// CalculateCreditDelta returns the credit cost (positive integer) for a request.
// Uses integer division so partial thousands are not charged.
func CalculateCreditDelta(inputTokens, outputTokens int64, rate *ModelCreditRate) int64 {
	return (inputTokens/1000)*rate.CreditsPer1kIn + (outputTokens/1000)*rate.CreditsPer1kOut
}

// CreditService handles all credit-related operations.
type CreditService struct {
	db    *dbent.Client
	redis *redis.Client
}

// NewCreditService constructs a CreditService.
func NewCreditService(db *dbent.Client, redis *redis.Client) *CreditService {
	return &CreditService{db: db, redis: redis}
}

func (s *CreditService) cacheKey(userID int64) string {
	return creditBalanceCachePrefix + fmt.Sprintf("%d", userID)
}

// GetBalance returns the credit balance for a user, preferring the Redis cache.
func (s *CreditService) GetBalance(ctx context.Context, userID int64) (*CreditBalance, error) {
	key := s.cacheKey(userID)
	cached, err := s.redis.Get(ctx, key).Int64()
	if err == nil {
		u, dbErr := s.db.User.Get(ctx, userID)
		if dbErr != nil {
			return nil, dbErr
		}
		return &CreditBalance{
			Balance:   cached,
			ExpiresAt: u.CreditExpiresAt,
			PlanID:    u.CreditPlanID,
		}, nil
	}

	u, err := s.db.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.redis.Set(ctx, key, u.CreditBalance, creditBalanceCacheTTL)
	return &CreditBalance{
		Balance:   u.CreditBalance,
		ExpiresAt: u.CreditExpiresAt,
		PlanID:    u.CreditPlanID,
	}, nil
}

// HasSufficientCredits checks whether a user has a positive credit balance.
// Prefers Redis cache; falls back to DB on cache miss.
func (s *CreditService) HasSufficientCredits(ctx context.Context, userID int64) (bool, error) {
	key := s.cacheKey(userID)
	bal, err := s.redis.Get(ctx, key).Int64()
	if err == redis.Nil {
		u, dbErr := s.db.User.Query().
			Where(entuser.ID(userID)).
			Select(entuser.FieldCreditBalance, entuser.FieldCreditExpiresAt).
			Only(ctx)
		if dbErr != nil {
			return false, dbErr
		}
		s.redis.Set(ctx, key, u.CreditBalance, creditBalanceCacheTTL)
		return u.CreditBalance > 0, nil
	}
	if err != nil {
		return false, err
	}
	return bal > 0, nil
}

// DeductCredits subtracts delta from the user's balance and appends a ledger row.
// delta must be positive. If the resulting balance would go negative it is floored to 0.
func (s *CreditService) DeductCredits(ctx context.Context, userID, delta int64, model, refID string) error {
	if delta <= 0 {
		return nil
	}
	u, err := s.db.User.Get(ctx, userID)
	if err != nil {
		return err
	}
	newBalance := u.CreditBalance - delta
	if newBalance < 0 {
		newBalance = 0
	}
	if err := s.db.User.UpdateOneID(userID).
		SetCreditBalance(newBalance).
		Exec(ctx); err != nil {
		return err
	}
	q := s.db.CreditLedger.Create().
		SetUserID(userID).
		SetDelta(-delta).
		SetReason("api_call").
		SetBalanceAfter(newBalance).
		SetModel(model)
	if refID != "" {
		q = q.SetRefID(refID)
	}
	if _, err := q.Save(ctx); err != nil {
		slog.Warn("credit_service: failed to write deduction ledger row", "userID", userID, "err", err)
	}
	s.redis.Del(ctx, s.cacheKey(userID))
	return nil
}

// CreditUser adds credits to a user after a successful purchase and sets expiry.
func (s *CreditService) CreditUser(ctx context.Context, userID, credits, planID int64, validityDays int, orderID string) error {
	u, err := s.db.User.Get(ctx, userID)
	if err != nil {
		return err
	}
	newBalance := u.CreditBalance + credits
	expiresAt := time.Now().AddDate(0, 0, validityDays)
	if err := s.db.User.UpdateOneID(userID).
		SetCreditBalance(newBalance).
		SetCreditExpiresAt(expiresAt).
		SetCreditPlanID(planID).
		Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.CreditLedger.Create().
		SetUserID(userID).
		SetDelta(credits).
		SetReason("purchase").
		SetRefID(orderID).
		SetBalanceAfter(newBalance).
		Save(ctx); err != nil {
		slog.Warn("credit_service: failed to write purchase ledger row", "userID", userID, "err", err)
	}
	s.redis.Del(ctx, s.cacheKey(userID))
	return nil
}

// AdminGrantCredits adds credits manually (admin action, reason = "admin_grant").
func (s *CreditService) AdminGrantCredits(ctx context.Context, userID, credits int64, notes string) error {
	u, err := s.db.User.Get(ctx, userID)
	if err != nil {
		return err
	}
	newBalance := u.CreditBalance + credits
	if err := s.db.User.UpdateOneID(userID).
		SetCreditBalance(newBalance).
		Exec(ctx); err != nil {
		return err
	}
	if _, err := s.db.CreditLedger.Create().
		SetUserID(userID).
		SetDelta(credits).
		SetReason("admin_grant").
		SetRefID(notes).
		SetBalanceAfter(newBalance).
		Save(ctx); err != nil {
		slog.Warn("credit_service: failed to write admin_grant ledger row", "err", err)
	}
	s.redis.Del(ctx, s.cacheKey(userID))
	return nil
}

// ListLedger returns paginated credit ledger entries for a user (newest first).
func (s *CreditService) ListLedger(ctx context.Context, userID int64, offset, limit int) ([]*CreditLedgerEntry, error) {
	rows, err := s.db.CreditLedger.Query().
		Where(creditledger.UserID(userID)).
		Order(dbent.Desc(creditledger.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*CreditLedgerEntry, 0, len(rows))
	for _, r := range rows {
		out = append(out, &CreditLedgerEntry{
			ID:           r.ID,
			Delta:        r.Delta,
			Reason:       r.Reason,
			RefID:        r.RefID,
			BalanceAfter: r.BalanceAfter,
			Model:        r.Model,
			CreatedAt:    r.CreatedAt,
		})
	}
	return out, nil
}

// GetAllModelRates fetches all model credit rates from the DB.
func (s *CreditService) GetAllModelRates(ctx context.Context) ([]*ModelCreditRate, error) {
	rows, err := s.db.ModelCreditRate.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*ModelCreditRate, 0, len(rows))
	for _, r := range rows {
		out = append(out, &ModelCreditRate{
			ID:              r.ID,
			ModelPattern:    r.ModelPattern,
			CreditsPer1kIn:  r.CreditsPer1kTokensInput,
			CreditsPer1kOut: r.CreditsPer1kTokensOutput,
			Priority:        r.Priority,
		})
	}
	return out, nil
}

// CheckExpiry resets credit_balance to 0 if credit_expires_at has passed.
// Call on login or lazily before balance checks.
func (s *CreditService) CheckExpiry(ctx context.Context, u *dbent.User) {
	if u.CreditExpiresAt == nil || !time.Now().After(*u.CreditExpiresAt) {
		return
	}
	if err := s.db.User.UpdateOneID(u.ID).SetCreditBalance(0).Exec(ctx); err != nil {
		slog.Warn("credit_service: expiry reset failed", "userID", u.ID, "err", err)
		return
	}
	s.db.CreditLedger.Create().
		SetUserID(int64(u.ID)).
		SetDelta(-u.CreditBalance).
		SetReason("expiry_reset").
		SetBalanceAfter(0).
		SaveX(ctx)
	s.redis.Del(ctx, s.cacheKey(int64(u.ID)))
}
