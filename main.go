package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Tags []string `json:"tags"`
	Ingredients []string `json:"ingredients"`
	Instruction []string `json:"instructions"`
	PublishedAt time.Time `json:"publishedAt"`
}

var recipes []Recipe
func init(){
	recipes = make([]Recipe, 0)
}

func NewRecipeHandler(c *gin.Context){
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":err.Error(),
		})
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now().UTC()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)
}

func main(){
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":"OK Response 200",
		})
	})

	r.POST("/recipes", NewRecipeHandler)

	r.Run()
}