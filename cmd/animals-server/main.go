// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// animal represents data about an animal.
type animal struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

var animals = map[string]animal{}

func main() {
	run(8080)
}

func run(port int) {
	router := gin.Default()
	router.GET("/v1/:collection", getAnimals)
	router.GET("/v1/:collection/:id", getAnimalByID)
	router.POST("/v1/:collection", postAnimals)

	router.Run(fmt.Sprintf("0.0.0.0:%d", port))
}

// getAnimals responds with the list of all animals as JSON.
func getAnimals(c *gin.Context) {
	values := []animal{}
	for _, a := range animals {
		values = append(values, a)
	}
	c.IndentedJSON(http.StatusOK, values)
}

// postAnimals adds an animal from JSON received in the request body.
func postAnimals(c *gin.Context) {
	var newAnimal animal

	// Call BindJSON to bind the received JSON to newAnimal.
	if err := c.BindJSON(&newAnimal); err != nil {
		return
	}

	// Add the new animal.
	animals[newAnimal.ID] = newAnimal
	c.IndentedJSON(http.StatusCreated, newAnimal)
}

// getAnimalByID locates the animal whose ID value matches the id
// parameter sent by the client, then returns that animal as a response.
func getAnimalByID(c *gin.Context) {
	id := c.Param("id")

	if a, ok := animals[id]; ok {
		c.IndentedJSON(http.StatusOK, a)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"code":    http.StatusNotFound,
		"message": "animal not found",
	})
}
