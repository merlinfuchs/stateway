package audit

import (
	"context"
	"fmt"

	"github.com/merlinfuchs/stateway/stateway-lib/service"
)

type AuditService struct {
	auditor Auditor
}

func NewAuditService(auditor Auditor) *AuditService {
	return &AuditService{auditor: auditor}
}

func (s *AuditService) ServiceType() service.ServiceType {
	return service.ServiceTypeAudit
}

func (s *AuditService) HandleRequest(ctx context.Context, method AuditMethod, request AuditRequest) (any, error) {
	switch req := request.(type) {
	case ConfigureAuditLoggingRequest:
		err := s.auditor.ConfigureAuditLogging(ctx, ConfigureAuditLoggingParams{
			GuildID: req.GuildID,
			Enabled: req.Enabled,
		}, req.Options.Destructure()...)
		if err != nil {
			return nil, err
		}
		return struct{}{}, nil
	case ListEntityChangesRequest:
		return s.auditor.GetEntityChanges(ctx, GetEntityChangesParams{
			GuildID:    req.GuildID,
			EntityType: req.EntityType,
			EntityID:   req.EntityID,
			Path:       req.Path,
		}, req.Options.Destructure()...)
	}
	return nil, fmt.Errorf("unknown request type: %T", request)
}
