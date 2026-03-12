package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/entity"
	moderationrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/moderation/repository"
	"github.com/google/uuid"
)

func (s *ModerationService) CreateCaseFromSupportHandoff(ctx context.Context, request dto.CreateCaseFromSupportHandoffRequest) (dto.CreateCaseFromSupportHandoffResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}
	if s.supportLookup == nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, ErrSupportReferenceRequired
	}

	supportID, err := parseID(request.SupportID, "support_id")
	if err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}
	_, err = parseID(request.ActorUserID, "actor_user_id")
	if err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}

	existing, err := s.store.GetCaseBySourceRef(ctx, entity.CaseSourceSupportReport, supportID)
	if err == nil {
		return dto.CreateCaseFromSupportHandoffResponse{Case: toCaseResponse(existing, true), Created: false}, nil
	}
	if !errors.Is(err, moderationrepository.ErrNotFound) {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}

	reference, err := s.supportLookup.GetModerationHandoffReference(ctx, supportID, strings.TrimSpace(request.RequestID), strings.TrimSpace(request.CorrelationID))
	if err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, ErrSupportReferenceInvalid
	}

	targetType, err := parseTargetType(reference.TargetType)
	if err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}
	targetID, err := parseID(reference.TargetID, "target_id")
	if err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}

	now := s.now().UTC()
	sourceRefID := supportID
	requestID := strings.TrimSpace(request.RequestID)
	if requestID == "" {
		requestID = strings.TrimSpace(reference.RequestID)
	}
	correlationID := strings.TrimSpace(request.CorrelationID)
	if correlationID == "" {
		correlationID = strings.TrimSpace(reference.CorrelationID)
	}

	moderationCase := entity.Case{
		ID:               uuid.NewString(),
		Source:           entity.CaseSourceSupportReport,
		SourceRefID:      &sourceRefID,
		RequestID:        requestID,
		CorrelationID:    correlationID,
		TargetType:       targetType,
		TargetID:         targetID,
		Status:           entity.CaseStatusQueued,
		AssignmentStatus: entity.AssignmentStatusUnassigned,
		EscalationStatus: entity.EscalationStatusNotEscalated,
		ActionResult:     entity.ActionResultNone,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.store.CreateCase(ctx, moderationCase); err != nil {
		if errors.Is(err, moderationrepository.ErrConflict) {
			persisted, getErr := s.store.GetCaseBySourceRef(ctx, entity.CaseSourceSupportReport, supportID)
			if getErr != nil {
				return dto.CreateCaseFromSupportHandoffResponse{}, getErr
			}
			return dto.CreateCaseFromSupportHandoffResponse{Case: toCaseResponse(persisted, true), Created: false}, nil
		}
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}

	if s.supportLinker != nil {
		if err := s.supportLinker.LinkModerationCase(ctx, supportID, moderationCase.ID); err != nil {
			return dto.CreateCaseFromSupportHandoffResponse{}, err
		}
	}

	persisted, err := s.store.GetCaseByID(ctx, moderationCase.ID)
	if err != nil {
		return dto.CreateCaseFromSupportHandoffResponse{}, err
	}
	return dto.CreateCaseFromSupportHandoffResponse{Case: toCaseResponse(persisted, true), Created: true}, nil
}
