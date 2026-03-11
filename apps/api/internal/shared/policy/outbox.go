package policy

// OutboxComponent defines required transactional outbox components.
type OutboxComponent string

const (
	OutboxComponentTransactionalRecord OutboxComponent = "transactional_record"
	OutboxComponentBackgroundPublisher OutboxComponent = "background_publisher"
	OutboxComponentRetryBackoff        OutboxComponent = "retry_backoff"
	OutboxComponentDeadLetter          OutboxComponent = "dead_letter"
	OutboxComponentObservability       OutboxComponent = "observability_metrics"
)

var RequiredOutboxComponents = []OutboxComponent{
	OutboxComponentTransactionalRecord,
	OutboxComponentBackgroundPublisher,
	OutboxComponentRetryBackoff,
	OutboxComponentDeadLetter,
	OutboxComponentObservability,
}

// OutboxMessageField defines mandatory message metadata.
type OutboxMessageField string

const (
	OutboxFieldEventID       OutboxMessageField = "event_id"
	OutboxFieldSchemaVersion OutboxMessageField = "schema_version"
	OutboxFieldRequestID     OutboxMessageField = "request_id"
	OutboxFieldCorrelationID OutboxMessageField = "correlation_id"
	OutboxFieldCausationID   OutboxMessageField = "causation_id"
)

var RequiredOutboxMessageFields = []OutboxMessageField{
	OutboxFieldEventID,
	OutboxFieldSchemaVersion,
	OutboxFieldRequestID,
	OutboxFieldCorrelationID,
	OutboxFieldCausationID,
}

var OutboxPriorityModules = []string{
	"payment",
	"inventory",
	"mission",
	"royalpass",
	"notification",
	"support",
	"moderation",
	"history",
}
