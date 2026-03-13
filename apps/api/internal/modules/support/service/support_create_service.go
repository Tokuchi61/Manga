package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/dto"
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/support/entity"
	supportrepository "github.com/Tokuchi61/Manga/apps/api/internal/modules/support/repository"
	"github.com/google/uuid"
)

func (s *SupportService) CreateCommunication(ctx context.Context, request dto.CreateCommunicationRequest) (dto.CreateSupportResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateSupportResponse{}, err
	}

	requesterID, err := parseID(request.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	priority, err := parsePriority(request.Priority, entity.SupportPriorityNormal)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	requestID, err := requiredRequestID(request.RequestID)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	return s.createCase(ctx, createCaseInput{
		RequesterUserID: requesterID,
		Kind:            entity.SupportKindCommunication,
		Category:        normalizeValue(request.Category),
		Priority:        priority,
		ReasonCode:      normalizeValue(request.ReasonCode),
		ReasonText:      request.ReasonText,
		Attachments:     request.Attachments,
		RequestID:       requestID,
	})
}

func (s *SupportService) CreateTicket(ctx context.Context, request dto.CreateTicketRequest) (dto.CreateSupportResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateSupportResponse{}, err
	}

	requesterID, err := parseID(request.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	priority, err := parsePriority(request.Priority, entity.SupportPriorityNormal)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	requestID, err := requiredRequestID(request.RequestID)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	return s.createCase(ctx, createCaseInput{
		RequesterUserID: requesterID,
		Kind:            entity.SupportKindTicket,
		Category:        normalizeValue(request.Category),
		Priority:        priority,
		ReasonCode:      normalizeValue(request.ReasonCode),
		ReasonText:      request.ReasonText,
		Attachments:     request.Attachments,
		RequestID:       requestID,
	})
}

func (s *SupportService) CreateReport(ctx context.Context, request dto.CreateReportRequest) (dto.CreateSupportResponse, error) {
	if err := s.validateInput(request); err != nil {
		return dto.CreateSupportResponse{}, err
	}

	requesterID, err := parseID(request.RequesterUserID, "requester_user_id")
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	priority, err := parsePriority(request.Priority, entity.SupportPriorityNormal)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	targetType, err := parseTargetType(request.TargetType)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}
	targetID, err := parseID(request.TargetID, "target_id")
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}
	if err := s.ensureTargetExists(ctx, targetType, targetID); err != nil {
		return dto.CreateSupportResponse{}, err
	}

	requestID, err := requiredRequestID(request.RequestID)
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}

	return s.createCase(ctx, createCaseInput{
		RequesterUserID: requesterID,
		Kind:            entity.SupportKindReport,
		Category:        normalizeValue(request.Category),
		Priority:        priority,
		ReasonCode:      normalizeValue(request.ReasonCode),
		ReasonText:      request.ReasonText,
		TargetType:      &targetType,
		TargetID:        &targetID,
		Attachments:     request.Attachments,
		RequestID:       requestID,
	})
}

type createCaseInput struct {
	RequesterUserID string
	Kind            entity.SupportKind
	Category        string
	Priority        entity.SupportPriority
	ReasonCode      string
	ReasonText      string
	TargetType      *entity.SupportTargetType
	TargetID        *string
	Attachments     []string
	RequestID       string
}

func (s *SupportService) createCase(ctx context.Context, input createCaseInput) (dto.CreateSupportResponse, error) {
	if input.Kind == entity.SupportKindReport {
		if input.TargetType == nil || input.TargetID == nil {
			return dto.CreateSupportResponse{}, ErrTargetRequired
		}
	}

	reasonText, err := sanitizeText(input.ReasonText, "reason_text")
	if err != nil {
		return dto.CreateSupportResponse{}, err
	}
	attachments := normalizeAttachments(input.Attachments)

	existing, lookupErr := s.store.FindCaseByRequesterRequestID(ctx, input.RequesterUserID, input.RequestID)
	if lookupErr == nil {
		return toCreateResponse(existing), nil
	}
	if !errors.Is(lookupErr, supportrepository.ErrNotFound) {
		return dto.CreateSupportResponse{}, lookupErr
	}

	now := s.now().UTC()
	var duplicateOfSupportID *string
	similar, similarErr := s.store.FindRecentSimilarCase(
		ctx,
		input.RequesterUserID,
		input.Kind,
		input.TargetType,
		input.TargetID,
		input.ReasonCode,
		now.Add(-s.duplicateWindow),
	)
	if similarErr == nil {
		duplicateOfSupportID = &similar.ID
	} else if !errors.Is(similarErr, supportrepository.ErrNotFound) {
		return dto.CreateSupportResponse{}, similarErr
	}

	status := entity.SupportStatusOpen
	risk := spamRiskScore(reasonText, attachments)
	if risk >= 90 {
		status = entity.SupportStatusSpam
	}

	supportCase := entity.SupportCase{
		ID:                   uuid.NewString(),
		RequesterUserID:      input.RequesterUserID,
		Kind:                 input.Kind,
		Category:             input.Category,
		Priority:             input.Priority,
		ReasonCode:           input.ReasonCode,
		ReasonText:           reasonText,
		TargetType:           input.TargetType,
		TargetID:             input.TargetID,
		Status:               status,
		DuplicateOfSupportID: duplicateOfSupportID,
		RequestID:            input.RequestID,
		SpamRiskScore:        risk,
		Attachments:          attachments,
		Replies:              nil,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := s.store.CreateCase(ctx, supportCase); err != nil {
		if errors.Is(err, supportrepository.ErrConflict) {
			existing, lookupErr := s.store.FindCaseByRequesterRequestID(ctx, input.RequesterUserID, input.RequestID)
			if lookupErr == nil {
				return toCreateResponse(existing), nil
			}
			return dto.CreateSupportResponse{}, ErrSupportAlreadyExists
		}
		return dto.CreateSupportResponse{}, err
	}

	return toCreateResponse(supportCase), nil
}

func toCreateResponse(support entity.SupportCase) dto.CreateSupportResponse {
	return dto.CreateSupportResponse{
		SupportID:            support.ID,
		SupportKind:          string(support.Kind),
		Status:               string(support.Status),
		DuplicateOfSupportID: support.DuplicateOfSupportID,
		TargetType:           toStringPtr(support.TargetType),
		TargetID:             support.TargetID,
	}
}

func (s *SupportService) ensureTargetExists(ctx context.Context, targetType entity.SupportTargetType, targetID string) error {
	switch targetType {
	case entity.SupportTargetTypeManga:
		if s.mangaLookup == nil {
			return nil
		}
		exists, err := s.mangaLookup.TargetExists(ctx, targetID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidTarget
		}
	case entity.SupportTargetTypeChapter:
		if s.chapterLookup == nil {
			return nil
		}
		exists, err := s.chapterLookup.TargetExists(ctx, targetID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidTarget
		}
	case entity.SupportTargetTypeComment:
		if s.commentLookup == nil {
			return nil
		}
		exists, err := s.commentLookup.TargetExists(ctx, targetID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidTarget
		}
	}
	return nil
}

func requiredRequestID(raw string) (string, error) {
	requestID := strings.TrimSpace(raw)
	if requestID == "" {
		return "", fmt.Errorf("%w: request_id is required", ErrValidation)
	}
	return requestID, nil
}
