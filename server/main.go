package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Response struct {
    AccesToken string      `json:"access_token"`
    RefreshToken string    `json:"refresh_token"`
}

func main() {
    Port := "8080"
    
    
    http.HandleFunc("/", HelloServer)
    http.HandleFunc("/token", Login)
    http.ListenAndServe(":"+Port, nil)
}


func HelloServer(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func Login(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query();
    code := query["code"];
    
    if code == nil{
        println("no code")
        return
    }

    enableCors(&w)
    ClientID := getEnv("ClientID")
    ClientSecret := getEnv("ClientSecret")
    callbackUrl :="http://localhost:3000/callback"
    endpoint := "https://accounts.spotify.com/api/token"
    stringToEncode := ClientID+":"+ClientSecret
    log.Println(stringToEncode)

    base64EncodedString := b64.StdEncoding.EncodeToString([]byte(stringToEncode));

    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Set("code", code[0])
    data.Set("redirect_uri",callbackUrl)
    

    client := &http.Client{}
    r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
    if err != nil {
        log.Fatal(err)
    }

    r.Header.Add("Authorization","Basic "+base64EncodedString)
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    res, err := client.Do(r)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(res.Status)
    defer res.Body.Close()
    body, err := ioutil.ReadAll(res.Body)

    if err != nil {
        log.Fatal(err)
    }

    spotifyResponseData := Response{}

    json.Unmarshal([]byte(body), &spotifyResponseData)
    tokens, _ :=json.Marshal(spotifyResponseData)
    
    w.Header().Set("Content-type", "application/json")
    w.WriteHeader(http.StatusOK);
    w.Write(tokens)
}




func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getEnv(key string) (string){
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
    return os.Getenv(key)
}