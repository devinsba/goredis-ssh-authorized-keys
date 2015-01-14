package main

import (
  "os"
  "fmt"

  "github.com/garyburd/redigo/redis"
)

func main() {
  if (os.Getenv("REDIS_HOST") == "" || os.Getenv("REDIS_PASSWORD") == "") {
    fmt.Fprintf(os.Stderr, "Please ensure environment variables REDIS_HOST and REDIS_PASSWORD have values\n")
    os.Exit(1)
  }
  if (len(os.Args) < 2) {
    fmt.Fprintf(os.Stderr, "No username to check provided\n")
    os.Exit(1)
  }
  getKeys(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), os.Args[1])
}

func getKeys(connectionString string, password string, username string) {
  conn, err := redis.Dial("tcp", connectionString)

  if (err == nil) {
    conn.Do("AUTH", password)
    
    keys, err := redis.Strings(conn.Do("SMEMBERS", fmt.Sprintf("user:%s", username)))

    if (err == nil) {
      for _, key := range keys {
        println(key)
      }
    }
  }

  defer conn.Close()
}