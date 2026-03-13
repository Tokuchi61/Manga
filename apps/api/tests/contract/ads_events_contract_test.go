package contract_test

import (
	"testing"

	adsevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/ads/events"
	"github.com/stretchr/testify/require"
)

func TestAdsEventConstants(t *testing.T) {
	require.Equal(t, "ads.impression.accepted", adsevents.EventAdsImpressionAccepted)
	require.Equal(t, "ads.click.accepted", adsevents.EventAdsClickAccepted)
	require.Equal(t, "ads.campaign.state_changed", adsevents.EventAdsCampaignStateChanged)
}
