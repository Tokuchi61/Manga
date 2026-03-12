package contract_test

import (
	"testing"

	socialevents "github.com/Tokuchi61/Manga/apps/api/internal/modules/social/events"
	"github.com/stretchr/testify/require"
)

func TestSocialEventConstants(t *testing.T) {
	require.Equal(t, "social.friendship.changed", socialevents.EventSocialFriendshipChanged)
	require.Equal(t, "social.follow.changed", socialevents.EventSocialFollowChanged)
	require.Equal(t, "social.message.sent", socialevents.EventSocialMessageSent)
	require.Equal(t, "social.wall.posted", socialevents.EventSocialWallPosted)
}
