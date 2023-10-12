# XYZ Books Codebase 2

---
 
 This program will get all books from backend(codebase 1) and will convert and compute the existing ISBN 13 / ISBN 10 to the missing ISBN. It will then update via the update endpoint of the backend. And save the new ISBN in a CSV file.

---

## Instructions
1. Edit and save the **.env** file that is included in the codebase 2.
	- Input the correct API Host and API Port of the codebase 1 (default is recommended).
2. Go to the project directory via terminal.
3. Run `go get`
4. Run `go build` 
5. After a successful build, you may run the `main` executable.

## Output
The CSV file will be outputted to the **output** directory of the codebase 2 project. Named, **output.csv**

## Versions used
- Go: 1.21.1