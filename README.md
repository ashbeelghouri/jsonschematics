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

### Loading Schematics From JSON file
Instead of defining the Schema directly, Load the schema from JSON file:

```sh
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schematics, err := jsonschematics.LoadFromJsonFile("path-to-your-schema.json")
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }else {
        fmt.Println("Schema Loaded Successfully")
    }
}
```
see the API Reference for json fields mapping.


### Loading Schematics From map[string]interface{}
If you want to load the schema from map[string]interface, you can use the below example:

```sh
package main

import (
    "fmt"
    "github.com/ashbeelghouri/jsonschematics"
)

func main() {
    schema := map[string]interface{}{
        ... define your schema
    }

    schematics, err := jsonschematics.LoadFromMap(&schema)
    if err != nil {
        fmt.Println("Unable to load the schema:", err)
    }else {
        fmt.Println("Schema Loaded Successfully")
    }
}
```

## API Reference

### Example Files
- [Schema](https://github.com/ashbeelghouri/jsonschematics/blob/master/json/schema.json)
- [Data](https://github.com/ashbeelghouri/jsonschematics/blob/master/json/data.json)

### Structs

#### Schematics
```sh
- Schema                                       Schema
- Validators                                   validators.Validators
- Prefix                                       string
- Separator                                    string
- ArrayIdKey                                   string
- LoadSchema(filePath string)                  error
- Validate(data map[string]interface{})        *ErrorMessages
- ValidateArray(data []map[string]interface{}) *[]ArrayOfErrors
- MakeFlat(data map[string]interface)          *map[string]interface{}
```

##### Schema
```sh
- Version string `json:"version"`
- Fields []Field `json:"fields"`
```

###### >Explanation
```sh
* Version is for the maintenance of the schema
* Fields contains the validation logic for all the keys
```

##### Field
```sh
- DependsOn   []string `json:"depends_on"`
- TargetKey   string `json:"target_key"`
- Description string `json:"description"`
- Validators  []string `json:"validators"`
- Constants   map[string]interface{} `json:"constants"`
```

###### >Explanation
```sh
* DependsOn will check if the keys in array exists in data
* TargetKey will target the value in the data throught the key
* Description can have anything to explain the data, this can also be empty
* Validators is an array of string "validation functions"
* Constants will have dynanmic constants for each validator
```

##### Constant
```sh
- Attributes map[string]interface{} `json:"attributes"`
- ErrMsg     string `json:"err"`
```
###### >Explanation
```sh
* Attributes are passed in to the validation function so, it can have any map string interface.
* ErrMsg is a string that is shown as an error when validation fails
```

#### Errors
- ArrayOfErrors
- ErrorMessages
- ErrorMessage

##### ArrayOfErrors
```sh
- Errors ErrorMessages
- ID     interface{}
```

##### ErrorMessages
```
- Messages                                                 []ErrorMessage
- AddError(validator string, target string, err string)
- HaveErrors()                                             bool
```

##### ErrorMessage
```sh
- Message   string
- Validator string
- Target    string
```

#### Go Version
```sh
go 1.22.1
```

## Contributing
1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix.
3. Write tests to cover your changes.
4. Send a pull request.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/ashbeelghouri/jsonschematics/blob/master/LICENSE) file for details.
