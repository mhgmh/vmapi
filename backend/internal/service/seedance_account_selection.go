package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
)

// SelectSeedanceAccount picks a schedulable Seedance API-key account for the group.
// Unlike OpenAI/Grok media, Seedance is not OpenAI-compatible, so we use a dedicated
// platform-filtered list instead of the OpenAI capability scheduler.
func (s *OpenAIGatewayService) SelectSeedanceAccount(
	ctx context.Context,
	groupID *int64,
	sessionHash string,
	requestedModel string,
	excludedIDs map[int64]struct{},
) (*AccountSelectionResult, error) {
	if s == nil || s.accountRepo == nil {
		return nil, ErrNoAvailableAccounts
	}

	// Sticky session first.
	if sessionHash != "" && s.cache != nil {
		if accountID, err := s.getStickySessionAccountID(ctx, groupID, sessionHash); err == nil && accountID > 0 {
			if excludedIDs == nil || !mapHasInt64Key(excludedIDs, accountID) {
				if account, err := s.getSchedulableAccount(ctx, accountID); err == nil && account != nil {
					if isSeedanceAccountEligible(ctx, account, requestedModel) {
						return s.trySelectSeedanceAccount(ctx, groupID, sessionHash, account)
					}
				}
			}
		}
	}

	accounts, err := s.listSchedulableAccounts(ctx, groupID, PlatformSeedance)
	if err != nil {
		return nil, err
	}
	candidates := make([]*Account, 0, len(accounts))
	for i := range accounts {
		account := &accounts[i]
		if excludedIDs != nil {
			if _, excluded := excludedIDs[account.ID]; excluded {
				continue
			}
		}
		if !isSeedanceAccountEligible(ctx, account, requestedModel) {
			continue
		}
		candidates = append(candidates, account)
	}
	if len(candidates) == 0 {
		if strings.TrimSpace(requestedModel) != "" {
			return nil, fmt.Errorf("%w supporting model: %s", ErrNoAvailableAccounts, requestedModel)
		}
		return nil, ErrNoAvailableAccounts
	}

	// Prefer sticky-compatible random shuffle for load spread.
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	var lastErr error
	for _, account := range candidates {
		selection, err := s.trySelectSeedanceAccount(ctx, groupID, sessionHash, account)
		if err == nil && selection != nil {
			return selection, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, ErrNoAvailableAccounts
}

func (s *OpenAIGatewayService) trySelectSeedanceAccount(
	ctx context.Context,
	groupID *int64,
	sessionHash string,
	account *Account,
) (*AccountSelectionResult, error) {
	if account == nil {
		return nil, ErrNoAvailableAccounts
	}
	if s.concurrencyService == nil {
		if sessionHash != "" {
			_ = s.BindStickySession(ctx, groupID, sessionHash, account.ID)
		}
		return s.newSelectionResult(ctx, account, false, nil, nil)
	}
	result, err := s.tryAcquireAccountSlot(ctx, account.ID, account.Concurrency)
	if err != nil {
		return nil, err
	}
	if result != nil && result.Acquired {
		if sessionHash != "" {
			_ = s.BindStickySession(ctx, groupID, sessionHash, account.ID)
		}
		return s.newAcquiredSelectionResult(ctx, account, result.ReleaseFunc)
	}
	// Soft-wait plan when slot is busy.
	cfg := s.schedulingConfig()
	return s.newSelectionResult(ctx, account, false, nil, &AccountWaitPlan{
		AccountID:      account.ID,
		MaxConcurrency: account.Concurrency,
		Timeout:        cfg.FallbackWaitTimeout,
		MaxWaiting:     cfg.FallbackMaxWaiting,
	})
}

func isSeedanceAccountEligible(ctx context.Context, account *Account, requestedModel string) bool {
	if account == nil || !account.IsSeedanceAPIKey() || !account.IsSchedulable() {
		return false
	}
	if strings.TrimSpace(account.GetSeedanceAPIKey()) == "" {
		return false
	}
	if strings.TrimSpace(requestedModel) == "" {
		return true
	}
	return account.IsSchedulableForModelWithContext(ctx, requestedModel)
}

func mapHasInt64Key(m map[int64]struct{}, key int64) bool {
	if m == nil {
		return false
	}
	_, ok := m[key]
	return ok
}
