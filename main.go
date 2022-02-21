// Package classification Recipes API.
//
// Documentation for Recipes API
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
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

// swagger:parameters recipes newRecipe
type Recipe struct {
	// swagger:ignore
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


// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
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

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
func ListRecipesHandler(c *gin.Context){
	c.JSON(200, recipes)
}

// swagger:operation GET /recipes/{id} recipes oneRecipe
// Get one recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
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

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: Invalid input
//     '404':
//         description: Invalid recipe ID
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

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
//     '404':
//         description: Invalid recipe ID
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

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
// responses:
//     '200':
//         description: Successful operation
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