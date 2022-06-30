package main

import (
    "os"
)

func main() {
    app := App{}
    app.Initialize(
            os.Getenv("PANGEA_TOKEN"))
    app.Run(":8080")
}

