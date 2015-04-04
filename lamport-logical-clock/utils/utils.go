package utils
import (
    "strconv"
    "log"
)

/* Small helper to translate strings to integers */
func GetInt(str string) (val int) {
    val, err := strconv.Atoi(str)
    if err != nil {
        log.Fatal("Cannot convert to integer : ", str)
    }

    return
}