# Butler Labs CLI

A CLI for interacting with [butlerlabs.ai](https://butlerlabs.ai). 

# Instructions 

1. Follow the instructions on https://docs.butlerlabs.ai/reference/overview to get the Queue Id and API Keys for your image model.
2. Create a file in the root of this repository titled secrets.yaml (see example_secrets.yaml for format) and put your keys in. 
3. Run `go build` in the project root to build the butler-cli binary
4. Run commands or type `./butler-cli` for help 

# Examples

## Upload a specific image to Butler and wait for the results
```
❯ ./butler-cli processImages /Users/maxwolffe/Desktop/ThymeChurros.png
```

## Extract a previous upload and save as CSV
```
❯ ./butler-cli -o "/Users/maxwolffe/Desktop/output.csv" getExtractionResults e506078d-c135-4d7d-9392-1767ae2ebdd7
```
