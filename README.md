# jsonschematics

`jsonschematics` is a Go package designed to validate and manipulate JSON data structures using schematics.

## Features

- Validate JSON objects against defined schematics
- Convert schematics to JSON schemas
- Handle complex data validation scenarios

## Installation

To install the package, use the following command:

```sh
go get github.com/ashbeelghouri/jsonschematics
```

## Usage

### Validating JSON Data
You can validate JSON data against a defined schematic using the Validate function. Here's an example:

```sh
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schema := jsonschematics.Schematics{
        // Define your schema here
    }

    data := map[string]interface{}{
        "Name": "John",
        "Age":  30,
    }

    err := schema.Validate(data)
    if err != nil {
        fmt.Println("Validation errors:", err)
    } else {
        fmt.Println("Validation successful")
    }
}
```

### Converting Schematics to JSON Schemas
Convert your schematics to JSON schemas easily:

```
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schema := jsonschematics.Schematics{
        // Define your schema here
    }

    jsonSchema, err := schema.ToJSONSchema()
    if err != nil {
        fmt.Println("Error converting to JSON schema:", err)
    } else {
        fmt.Println("JSON Schema:", jsonSchema)
    }
}
```

## API Reference
### Structs

#### Schematics
- Validate(data map[string]interface{}) error
- ToJSONSchema() (string, error)

#### Errors
- ArrayOfErrors
- ErrorMessages
- ErrorMessage

## Contributing
1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix.
3. Write tests to cover your changes.
4. Send a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/ashbeelghouri/jsonschematics/blob/master/LICENSE) file for details.
