package controller

type AppController struct {
	Project  interface{ Project }
	Secret   interface{ Secret }
	Resource interface{ Resource }
}
