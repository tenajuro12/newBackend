package routes

import (
	"github.com/gofiber/fiber/v2"
	controller2 "github.com/tenajuro12/newBackend/cmd/controller"
	"github.com/tenajuro12/newBackend/pkg/middleware"
)

func Setup(app *fiber.App) {
	controller2.LoadTemplates()
	app.Post("/api/register", controller2.Register)
	app.Post("/api/login", controller2.Login)

	app.Use(middleware.IsAuthenticate)

	app.Get("/login", controller2.RenderLoginPage)
	app.Get("/register", controller2.RenderRegisterPage)

	app.Get("createBlog", controller2.RenderCreateBlogPage)
	app.Post("/api/posts", controller2.CreatePost)

	app.Get("/allPost", controller2.RenderAllPostPage)
	app.Get("/api/allpost", controller2.AllPost)
	app.Get("/api/allpost/:id", controller2.DetailPost)

	app.Get("updateBlog/:id", controller2.RenderUpdateBlogPage)
	app.Put("/api/updatepost/:id", controller2.UpdatePost)

	app.Get("/api/uniquepost", controller2.UniquePost)
	app.Delete("/api/deletepost/:id", controller2.DeletePost)
	app.Post("/api/uploads", controller2.UploadImage)

	app.Get("/api/user", controller2.GetUserInfo)
	app.Delete("/api/user", controller2.DeleteUser) // Delete user account
	app.Put("/api/user", controller2.UpdateUser)

	app.Post("/api/post/:id/comment", controller2.CreateComment)           // Create a new comment for a blog post
	app.Get("/api/post/:id/comments", controller2.ReadComments)            // Retrieve all comments for a blog post
	app.Put("/api/post/:id/comment/:commentID", controller2.UpdateComment) // Update a specific comment
	app.Delete("/api/post/:id/comment/:commentID", controller2.DeleteComment)

	app.Post("/api/follow/:id", controller2.FollowUser)
	app.Delete("/api/unfollow/:id", controller2.UnfollowUser)

	app.Get("/api/posts/followed", controller2.GetPostsFromFollowedUsers)

	app.Static("/api/uploads", "./uploads")
}
