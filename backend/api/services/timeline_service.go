package services

import (
	"context"
	"sort"
	"time"

	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// TimelineEvent represents a single chronological event on a matter's timeline.
type TimelineEvent struct {
	Type string      `json:"type"` // "hearing" or "order"
	Date time.Time   `json:"date"`
	Data interface{} `json:"data"`
}

// TimelineService defines the interface for building a matter's chronological timeline.
type TimelineService interface {
	GetMatterTimeline(ctx context.Context, matterID string) ([]*TimelineEvent, error)
}

// timelineService implements TimelineService.
type timelineService struct {
	hearingRepo repositories.HearingRepository
	orderRepo   repositories.OrderRepository
}

// NewTimelineService creates a new TimelineService.
func NewTimelineService(hearingRepo repositories.HearingRepository, orderRepo repositories.OrderRepository) TimelineService {
	return &timelineService{hearingRepo: hearingRepo, orderRepo: orderRepo}
}

// GetMatterTimeline returns hearings and orders for a matter merged into a single
// chronologically sorted timeline (most recent first).
func (s *timelineService) GetMatterTimeline(ctx context.Context, matterID string) ([]*TimelineEvent, error) {
	hearings, err := s.hearingRepo.ListByMatter(ctx, matterID)
	if err != nil {
		return nil, err
	}

	orders, err := s.orderRepo.ListByMatter(ctx, matterID)
	if err != nil {
		return nil, err
	}

	events := make([]*TimelineEvent, 0, len(hearings)+len(orders))
	for _, hearing := range hearings {
		events = append(events, &TimelineEvent{Type: "hearing", Date: hearing.Date, Data: hearing})
	}
	for _, order := range orders {
		events = append(events, &TimelineEvent{Type: "order", Date: order.Date, Data: order})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.After(events[j].Date)
	})

	return events, nil
}
