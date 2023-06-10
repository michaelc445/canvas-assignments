# Instructions

1. Go to your canvas profile -> settings -> new access token
2. Replace the token var with your new token
3. Update the universityUrl var. Currently it points to UCC's canvas page. You'll probably need to update this to your own university's canvas page. 
4. Update the userId var to your own canvas userId. I found this by going to Account > Folio and it was in the URL.
    ```
    https://ucc.instructure.com/users/{user_id}/external_tools/962
    ```
5. Update your modules in the main function 
    ```go
    courses = map[string]int{
		"C-Programming for Microcontrollers": 48520,
		"Networks and data communications":   48480,
		"Theory of Computation":              48496,
		"Advanced Programming with Java":     48465,
		"Software Engineering":               48475,
		"Ethical Hacking and Web Security":   48505,
	}
    ```
    The key can be anything you like but the name makes the most sense. The integer value is your module code on canvas, you can get this by visiting the module page and it will be in the url.
    Here is the url the home page of the Software Engineering module
    ```
    https://ucc.instructure.com/courses/48475
    ```
    Just create an entry in the courses map for each of your modules.

6.  Change directory to where you have the go files. Build the binary: 
    ```sh
    $ go build .
    ```

7. Run the binary with no flags to print upcoming assignments.
    ```sh
    $ ./assignments
    ```
    or 
    ```sh
    $ ./assignments --job=all
    ```
    Print out all assignment grades for all modules:
    ```sh
    $ ./assignments --job=grades
    ```