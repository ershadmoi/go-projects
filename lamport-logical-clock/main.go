package main
import(
    "github.com/ershadmoi/go-projects/lamport-logical-clock/config"
    "fmt"
)

func main() {
    fmt.Println(config.ReadConfig("config/config.txt"))
}
