## TODO Command Line Tool in GOlang 
This tool is written in go to practice the capabilities of Go lang in command line tools. 

## Commands
### Make a Build Command 
`cd cmd`
`go build -o <binaryName>`

### Example 
`go build -o todo`

### Run the build 
`./todo --add "Get a pizza"`

### List of Commands
--add -> Add an item is todo
--list -> list all the items 
--complete -> complete given index item ( using a number )
--delete -> delete given index item
--verbose -> print details like time created, completed etc.

### Run Tests

`go test`
