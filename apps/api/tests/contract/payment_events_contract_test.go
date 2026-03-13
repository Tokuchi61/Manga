package contract_test

import (
	"testing"

	paymentevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/payment/events"
	"github.com/stretchr/testify/require"
)

func TestPaymentEventConstants(t *testing.T) {
	require.Equal(t, "payment.checkout.started", paymentevents.EventPaymentCheckoutStarted)
	require.Equal(t, "payment.callback.accepted", paymentevents.EventPaymentCallbackAccepted)
	require.Equal(t, "payment.transaction.settled", paymentevents.EventPaymentTransactionSettled)
	require.Equal(t, "payment.refund.completed", paymentevents.EventPaymentRefundCompleted)
}
