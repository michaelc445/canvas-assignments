# Instructions



1. Go to your canvas profile -> settings -> new access token
2. Replace the token in the getAssignmentForCourse function with your new token
3. Update the url in the main function. Currently it points to UCC's canvas api. You'll probably need to update this to your own university's canvas api 
4. Update your modules in the main function 
    ```go
    courses := map[string]int{
		"Team Software Project":           48420,
		"Workplace Technology and Skills": 48425,
	}
    ```
    The key can be anything you like but the name makes the most sense. The integer value is your module code on canvas, you can get this by visiting the module page and it will be in the url.
    Here is the url the home page of the Team Software Project module
    ```
    https://ucc.instructure.com/courses/48420
    ```
    Just create an entry in the courses map for each of your modules.

5.  ```sh
    $ go run assignment.go
    ```