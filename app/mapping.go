package app

import "kafka/controller"

func Routers() {
	r.POST("/users/create", controller.CreateUser)
	r.GET("/users/login", controller.Login)
	r.GET("/users/getbyemail", controller.GetUserByEmail)
	r.DELETE("/users/deletebyemail", controller.DeleteUserByEmail)
	r.PATCH("/users/updatebyemail", controller.UpdateUserByEmail)
}
