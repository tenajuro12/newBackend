package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tenajuro12/blogbackend/controller"
	"github.com/tenajuro12/blogbackend/middleware"
	"github.com/tenajuro12/blogbackend/renders"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)

	app.Post("/api/login", controller.Login)
	app.Use(middleware.IsAuthenticate)

	app.Post("/api/post", controller.CreatePost)
	app.Get("/api/allpost", controller.AllPost)
	app.Get("/api/allpost/:id", controller.DetailPost)
	app.Put("/api/updatepost/:id", controller.UpdatePost)
	app.Get("/api/uniquepost", controller.UniquePost)
	app.Delete("/api/deletepost/:id", controller.DeletePost)
	app.Post("/api/upload-image", controller.Upload)

	app.Get("/api/user", controller.GetUserInfo)
	app.Delete("/api/user", controller.DeleteUser) // Delete user account
	app.Put("/api/user", controller.UpdateUser)

	app.Post("/api/post/:id/comment", controller.CreateComment)           // Create a new comment for a blog post
	app.Get("/api/post/:id/comments", controller.ReadComments)            // Retrieve all comments for a blog post
	app.Put("/api/post/:id/comment/:commentID", controller.UpdateComment) // Update a specific comment
	app.Delete("/api/post/:id/comment/:commentID", controller.DeleteComment)

	app.Post("/api/follow/:id", controller.FollowUser)
	app.Delete("/api/unfollow/:id", controller.UnfollowUser)

	app.Get("/api/posts/followed", controller.GetPostsFromFollowedUsers)

	app.Static("/api/uploads", "./uploads")
}
