function transform(input) {
    inputJSON = JSON.parse(input)
    outputJSON = {}

    // your code goes here
    outputJSON.statement = inputJSON.greeting + " " + inputJSON.target + inputJSON.punctuation

    return JSON.stringify(outputJSON)
}