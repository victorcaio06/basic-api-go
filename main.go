package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var Students = []Student{
	Student{ID: 1, Name: "Victor Caio", Age: 23}, Student{ID: 2, Name: "Ludmila Martins", Age: 21},
}

func main() {
	service := gin.Default()

	getRoutes(service)

	service.Run(":3434")

}

func getRoutes(ctx *gin.Engine) *gin.Engine {
	ctx.GET("/heart", routeHeart)

	groupStudents := ctx.Group("/students")
	groupStudents.GET("/", routeGetStudents)
	groupStudents.POST("/create-student", createStudents)
	groupStudents.PUT("/update-student/:id", updateStudent)
	groupStudents.DELETE("/delete-student/:id", deleteStudent)

	return ctx
}

func routeHeart(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "OK!"})

	ctx.Done()
}

func routeGetStudents(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Students)

	ctx.Done()
}

func createStudents(ctx *gin.Context) {
	var student Student

	err := ctx.Bind(&student)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "Invalid fields!",
		})

		return
	}

	student.ID = Students[len(Students)-1].ID + 1

	Students = append(Students, student)

	ctx.JSON(http.StatusCreated, student)

	ctx.Done()
}

func updateStudent(ctx *gin.Context) {
	var studentPayload Student
	var studentLocal Student
	var positionSlice int

	err := ctx.BindJSON(&studentPayload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "Fields invalid!",
		})

		return
	}

	id, err := strconv.Atoi(ctx.Params.ByName("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "ID invalid!",
		})

		return
	}

	for i, studentElement := range Students {
		if studentElement.ID == id {
			studentLocal = studentElement
			positionSlice = i
		}
	}

	if studentLocal.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "Student not found!",
		})

		return
	}

	studentLocal.Name = studentPayload.Name
	studentLocal.Age = studentPayload.Age

	Students[positionSlice].Name = studentLocal.Name
	Students[positionSlice].Age = studentLocal.Age

	ctx.JSON(http.StatusCreated, Students[positionSlice])

	ctx.Done()
}

func deleteStudent(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Params.ByName("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"ERROR": "ID invalid!",
		})
	}

	for index, e := range Students {
		if e.ID == id {
			Students = append(Students[0:index], Students[index+1:]...)

			ctx.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{
		"ERROR": "Student not found!",
	})
}
