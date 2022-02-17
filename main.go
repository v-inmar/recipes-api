package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}


// Create new recipe handler
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

// Read (ALL) recipe handler
func ListRecipesHandler(c *gin.Context){
	c.JSON(200, recipes)
}

// Read recipe by id handler
// This is Bonus
func ReadRecipeHandler(c *gin.Context){
	id := c.Param("id")
	for i := 0; i < len(recipes); i++{
		if recipes[i].ID == id{
			c.JSON(http.StatusOK, recipes[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":"Recipe not found",
	})
}

// Update a recipe by id handler
func UpdateRecipeHandler(c *gin.Context){
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	index := -1
	for i := 0; i < len(recipes); i++{
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
		})
	}

	recipes[index] = recipe

	// TODO: More input checking is needed to be production ready
	// TODO: Empty fields passed will remove values including id
	// TODO: For now, it is good
	c.JSON(http.StatusOK, recipe)
}

// Delete a recipe by id handler
func DeleteRecipeHandler(c *gin.Context){
	id := c.Param("id")
	for i := 0; i < len(recipes); i++{
		if recipes[i].ID == id{
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message":fmt.Sprintf("Recipe with id %s has been removed", id),
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found",
	})
}

// Search a recipe by tags or keywords
func SearchRecipesHandler(c *gin.Context){
	tag := c.Query("tag")
	resultedRecipes := make([]Recipe, 0)

	for i := 0; i<len(recipes); i++{
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag){
				resultedRecipes = append(resultedRecipes, recipes[i])
				break
			}
		}
	}

	c.JSON(http.StatusOK, resultedRecipes)
}

func main(){
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":"OK Response 200",
		})
	})

	r.POST("/recipes", NewRecipeHandler)
	r.GET("/recipes", ListRecipesHandler)
	r.GET("/recipes/:id", ReadRecipeHandler)
	r.PUT("/recipes/:id", UpdateRecipeHandler)
	r.DELETE("/recipes/:id", DeleteRecipeHandler)
	r.GET("/recipes/search", SearchRecipesHandler)
	r.Run()
}