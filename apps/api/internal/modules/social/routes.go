package social

import (
	"github.com/Tokuchi61/Manga/apps/api/internal/modules/social/handler"
	"github.com/Tokuchi61/Manga/apps/api/internal/shared/identity"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(router chi.Router, httpHandler *handler.HTTPHandler) {
	if router == nil {
		panic("social registerRoutes: router cannot be nil")
	}
	if httpHandler == nil {
		panic("social registerRoutes: http handler cannot be nil")
	}

	router.Route("/social", func(r chi.Router) {
		r.With(identity.RequireUser).Post("/friendships/requests", httpHandler.CreateFriendRequest)
		r.With(identity.RequireUser).Post("/friendships/requests/{requester_user_id}/respond", httpHandler.RespondFriendRequest)
		r.With(identity.RequireUser).Delete("/friendships/{target_user_id}", httpHandler.RemoveFriend)
		r.With(identity.RequireUser).Get("/friendships/requests", httpHandler.ListFriendRequests)
		r.With(identity.RequireUser).Get("/friendships", httpHandler.ListFriends)

		r.With(identity.RequireUser).Post("/follow/{target_user_id}", httpHandler.FollowUser)
		r.With(identity.RequireUser).Delete("/follow/{target_user_id}", httpHandler.UnfollowUser)
		r.With(identity.RequireUser).Get("/follow/followers", httpHandler.ListFollowers)
		r.With(identity.RequireUser).Get("/follow/following", httpHandler.ListFollowing)

		r.With(identity.RequireUser).Post("/wall/posts", httpHandler.CreateWallPost)
		r.With(identity.RequireUser).Post("/wall/posts/{post_id}/replies", httpHandler.CreateWallReply)
		r.With(identity.RequireUser).Get("/wall/posts", httpHandler.ListWallPosts)

		r.With(identity.RequireUser).Post("/messages/threads", httpHandler.OpenThread)
		r.With(identity.RequireUser).Get("/messages/threads", httpHandler.ListThreads)
		r.With(identity.RequireUser).Post("/messages/threads/{thread_id}/read", httpHandler.MarkThreadRead)
		r.With(identity.RequireUser).Post("/messages/threads/{thread_id}/messages", httpHandler.SendMessage)
		r.With(identity.RequireUser).Get("/messages/threads/{thread_id}/messages", httpHandler.ListThreadMessages)

		r.With(identity.RequireUser).Post("/relations/block/{target_user_id}", httpHandler.UpdateBlock)
		r.With(identity.RequireUser).Post("/relations/mute/{target_user_id}", httpHandler.UpdateMute)
		r.With(identity.RequireUser).Post("/relations/restrict/{target_user_id}", httpHandler.UpdateRestrict)
		r.With(identity.RequireUser).Get("/relations/{relation_type}", httpHandler.ListRelations)

		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Get("/admin/runtime", httpHandler.GetRuntimeConfig)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/friendship-state", httpHandler.UpdateFriendshipState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/follow-state", httpHandler.UpdateFollowState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/wall-state", httpHandler.UpdateWallState)
		r.With(identity.RequireUser, identity.RequireAnyRole("admin")).Post("/admin/messaging-state", httpHandler.UpdateMessagingState)
	})
}
