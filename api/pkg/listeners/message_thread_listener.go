package listeners

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/entities"

	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/palantir/stacktrace"
)

// MessageThreadListener handles cloud events which need to update entities.MessageThread
type MessageThreadListener struct {
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.MessageThreadService
}

// NewMessageThreadListener creates a new instance of MessageThreadListener
func NewMessageThreadListener(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.MessageThreadService,
) (l *MessageThreadListener, routes map[string]events.EventListener) {
	l = &MessageThreadListener{
		logger:  logger.WithService(fmt.Sprintf("%T", l)),
		tracer:  tracer,
		service: service,
	}

	return l, map[string]events.EventListener{
		events.EventTypeMessageAPISent:               l.OnMessageAPISent,
		events.MessageAPIDeleted:                     l.onMessageDeleted,
		events.EventTypeMessagePhoneSending:          l.OnMessagePhoneSending,
		events.EventTypeMessagePhoneSent:             l.OnMessagePhoneSent,
		events.EventTypeMessagePhoneDelivered:        l.OnMessagePhoneDelivered,
		events.EventTypeMessageSendFailed:            l.OnMessagePhoneFailed,
		events.EventTypeMessagePhoneReceived:         l.OnMessagePhoneReceived,
		events.EventTypeMessageNotificationScheduled: l.onMessageNotificationScheduled,
		events.EventTypeMessageSendExpired:           l.onMessageExpired,
		events.UserAccountDeleted:                    l.onUserAccountDeleted,
	}
}

// OnMessageAPISent handles the events.EventTypeMessageAPISent event
func (listener *MessageThreadListener) OnMessageAPISent(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessageAPISentPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		Contact:   payload.Contact,
		UserID:    payload.UserID,
		Status:    entities.MessageStatusPending,
		Timestamp: payload.RequestReceivedAt,
		Content:   payload.Content,
		MessageID: payload.MessageID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// onMessageDeleted handles the events.MessageAPIDeleted event
func (listener *MessageThreadListener) onMessageDeleted(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	payload := new(events.MessageAPIDeletedPayload)
	if err := event.DataAs(payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.service.UpdateAfterDeletedMessage(ctx, payload); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", payload.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// OnMessagePhoneSending handles the events.EventTypeMessagePhoneSending event
func (listener *MessageThreadListener) OnMessagePhoneSending(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessagePhoneSendingPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		UserID:    payload.UserID,
		Contact:   payload.Contact,
		Status:    entities.MessageStatusSending,
		Timestamp: payload.Timestamp,
		Content:   payload.Content,
		MessageID: payload.ID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// OnMessagePhoneSent handles the events.EventTypeMessagePhoneSent event
func (listener *MessageThreadListener) OnMessagePhoneSent(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessagePhoneSentPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		Contact:   payload.Contact,
		UserID:    payload.UserID,
		Status:    entities.MessageStatusSent,
		Timestamp: payload.Timestamp,
		Content:   payload.Content,
		MessageID: payload.ID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// OnMessagePhoneDelivered handles the events.EventTypeMessagePhoneDelivered event
func (listener *MessageThreadListener) OnMessagePhoneDelivered(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessagePhoneDeliveredPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		UserID:    payload.UserID,
		Contact:   payload.Contact,
		Status:    entities.MessageStatusDelivered,
		Timestamp: payload.Timestamp,
		Content:   payload.Content,
		MessageID: payload.ID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// OnMessagePhoneFailed handles the events.EventTypeMessageSendFailed event
func (listener *MessageThreadListener) OnMessagePhoneFailed(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessageSendFailedPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T] for event [%s]", event.Data(), payload, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		Contact:   payload.Contact,
		UserID:    payload.UserID,
		Status:    entities.MessageStatusFailed,
		Timestamp: payload.Timestamp,
		Content:   payload.Content,
		MessageID: payload.ID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// OnMessagePhoneReceived handles the events.EventTypeMessagePhoneReceived event
func (listener *MessageThreadListener) OnMessagePhoneReceived(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessagePhoneReceivedPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		Contact:   payload.Contact,
		Timestamp: payload.Timestamp,
		UserID:    payload.UserID,
		Status:    entities.MessageStatusReceived,
		Content:   payload.Content,
		MessageID: payload.MessageID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// onMessageNotificationScheduled handles the events.EventTypeMessageNotificationScheduled event
func (listener *MessageThreadListener) onMessageNotificationScheduled(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessageNotificationScheduledPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		Contact:   payload.Contact,
		Timestamp: payload.ScheduledAt,
		UserID:    payload.UserID,
		Content:   payload.Content,
		Status:    entities.MessageStatusScheduled,
		MessageID: payload.MessageID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// onMessageNotificationScheduled handles the events.EventTypeMessageNotificationScheduled event
func (listener *MessageThreadListener) onMessageExpired(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.MessageSendExpiredPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	updateParams := services.MessageThreadUpdateParams{
		Owner:     payload.Owner,
		Contact:   payload.Contact,
		Timestamp: payload.Timestamp,
		UserID:    payload.UserID,
		Content:   payload.Content,
		Status:    entities.MessageStatusExpired,
		MessageID: payload.MessageID,
	}

	if err := listener.service.UpdateThread(ctx, updateParams); err != nil {
		msg := fmt.Sprintf("cannot update thread for message with ID [%s] for event with ID [%s]", updateParams.MessageID, event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (listener *MessageThreadListener) onUserAccountDeleted(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.UserAccountDeletedPayload
	if err := event.DataAs(&payload); err != nil {
		msg := fmt.Sprintf("cannot decode [%s] into [%T]", event.Data(), payload)
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err := listener.service.DeleteAllForUser(ctx, payload.UserID); err != nil {
		msg := fmt.Sprintf("cannot delete [entities.MessageThread] for user [%s] on [%s] event with ID [%s]", payload.UserID, event.Type(), event.ID())
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (listener *MessageThreadListener) updateThread(ctx context.Context, params services.MessageThreadUpdateParams) error {
	return listener.service.UpdateThread(ctx, params)
}
