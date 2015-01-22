package main

import (
  "flag"
  "os"
  "os/user"
  "fmt"
  "bufio"
  "strings"

  "github.com/garyburd/redigo/redis"
)

func main() {
  admin := flag.Bool("admin", false, "adminstrate keys instead of retrieving")

  c, _ := user.Current()
  runningUser := c.Username
  username := flag.String("username", runningUser, "the username we are administering")

  keyFile := flag.String("key_file", fmt.Sprintf("%s/.ssh/id_rsa.pub", os.Getenv("HOME")), "the public key file to add")

  flag.Parse()

  conn, err := redis.Dial("tcp", os.Getenv("REDIS_HOST"))
  if (err == nil) {
    conn.Do("AUTH", os.Getenv("REDIS_PASSWORD"))
  } else {
    conn.Close()
    fmt.Fprintf(os.Stderr, "Cannot connect to redis\n")
    os.Exit(1)
  }

  if (*admin) {
    println("Adding key")
    file, _ := os.Open(*keyFile)
    reader := bufio.NewReader(file)
    key, _ := reader.ReadString('\n')
    key = strings.TrimSuffix(key, "\n")
    addKey(conn, *username, key)
  } else {
    if (os.Getenv("REDIS_HOST") == "" || os.Getenv("REDIS_PASSWORD") == "") {
      fmt.Fprintf(os.Stderr, "Please ensure environment variables REDIS_HOST and REDIS_PASSWORD have values\n")
      conn.Close()
      os.Exit(1)
    }
    if (len(os.Args) < 2) {
      fmt.Fprintf(os.Stderr, "No username to check provided\n")
      conn.Close()
      os.Exit(1)
    }
    getKeys(conn, os.Args[1])
  }
  defer conn.Close()
}

func getKeys(conn redis.Conn, username string) {
  keys, err := redis.Strings(conn.Do("SMEMBERS", fmt.Sprintf("user:%s", username)))

  if (err == nil) {
    for _, key := range keys {
      println(key)
    }
  }
}

func addKey(conn redis.Conn, username string, key string) {
  exists, err := redis.Bool(conn.Do("SISMEMBER", fmt.Sprintf("user:%s", username), key))
  if (err == nil) {
    if (!exists) {
      conn.Do("SADD", fmt.Sprintf("user:%s", username), key)
    }
  }
}