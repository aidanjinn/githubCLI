package main

import (
	"fmt"
	"githubCLI/githubAPI"
	"net/http"
	"strconv"
)

func mainLoop() {
	var userName string
	fmt.Print("Input Github Username: ")
	_, err := fmt.Scan(&userName)

	if err != nil {
		fmt.Println(err)
		return
	}

	var eventString string
	fmt.Print("Input Number of Events to Display: ")
	_, err = fmt.Scan(&eventString)

	if err != nil {
		fmt.Println(err)
		return
	}

	var eventInt int
	eventInt, err = strconv.Atoi(eventString)

	if err != nil {
		fmt.Println(err)
		return
	}

	var userData githubAPI.GithubUser
	var statusCode int

	statusCode, userData, err = githubAPI.Ping(userName, eventInt)

	if statusCode != http.StatusOK {
		fmt.Printf("Expected status code %d, got %d", http.StatusOK, statusCode)
		return
	}

	githubAPI.ToString(&userData)
}
