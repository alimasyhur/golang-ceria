package main

import (
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "io/ioutil"
    "fmt"
    "log"
    "os"
    "net/http"

    "<a class="vglnk" href="http://github.com/gin-gonic/contrib/sessions" rel="nofollow"><span>github</span><span>.</span><span>com</span><span>/</span><span>gin</span><span>-</span><span>gonic</span><span>/</span><span>contrib</span><span>/</span><span>sessions</span></a>"
    "<a class="vglnk" href="http://github.com/gin-gonic/gin" rel="nofollow"><span>github</span><span>.</span><span>com</span><span>/</span><span>gin</span><span>-</span><span>gonic</span><span>/</span><span>gin</span></a>"
    "<a class="vglnk" href="http://golang.org/x/oauth2" rel="nofollow"><span>golang</span><span>.</span><span>org</span><span>/</span><span>x</span><span>/</span><span>oauth2</span></a>"
    "<a class="vglnk" href="http://golang.org/x/oauth2/google" rel="nofollow"><span>golang</span><span>.</span><span>org</span><span>/</span><span>x</span><span>/</span><span>oauth2</span><span>/</span><span>google</span></a>"
)

// Credentials which stores google ids.
type Credentials struct{
    Cid     string `json:"cid"`
    Csecret string `json:"csecret"`
}

// User is a retrieved and authentiacted user.
type User struct {
    Sub string `json:"sub"`
    Name string `json:"name"`
    GivenName string `json:"given_name"`
    FamilyName string `json:"family_name"`
    Profile string `json:"profile"`
    Picture string `json:"picture"`
    Email string `json:"email"`
    EmailVerified string `json:"email_verified"`
    Gender string `json:"gender"`
}

var cred Credentials
var conf *oauth2.Config
var state string
var store = sessions.NewCookieStore([]byte("secret"))

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func init() {
    file, err := ioutil.ReadFile("./creds.json")
    if err != nil {
        log.Printf("File error: %v\n", err)
        os.Exit(1)
    }
    json.Unmarshal(file, &cred)

    conf = &oauth2.Config{
        ClientID:     cred.Cid,
        ClientSecret: cred.Csecret,
        RedirectURL:  "<a class="vglnk" href="http://127.0.0.1:9090/auth" rel="nofollow"><span>http</span><span>://</span><span>127</span><span>.</span><span>0</span><span>.</span><span>0</span><span>.</span><span>1</span><span>:</span><span>9090</span><span>/</span><span>auth</span></a>",
        Scopes: []string{
            "<a class="vglnk" href="https://www.googleapis.com/auth/userinfo.email" rel="nofollow"><span>https</span><span>://</span><span>www</span><span>.</span><span>googleapis</span><span>.</span><span>com</span><span>/</span><span>auth</span><span>/</span><span>userinfo</span><span>.</span><span>email</span></a>", // You have to select your own scope from here -> <a class="vglnk" href="https://developers.google.com/identity/protocols/googlescopes#google_sign-in" rel="nofollow"><span>https</span><span>://</span><span>developers</span><span>.</span><span>google</span><span>.</span><span>com</span><span>/</span><span>identity</span><span>/</span><span>protocols</span><span>/</span><span>googlescopes</span><span>#</span><span>google</span><span>_</span><span>sign</span><span>-</span><span>in</span></a>
        },
        Endpoint: google.Endpoint,
    }
}

func indexHandler(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func getLoginURL(state string) string {
    return conf.AuthCodeURL(state)
}

func authHandler(c *gin.Context) {
    // Handle the exchange code to initiate a transport.
    session := sessions.Default(c)
    retrievedState := session.Get("state")
    if retrievedState != c.Query("state") {
        c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
        return
    }

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
        return
	}

	client := conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("<a class="vglnk" href="https://www.googleapis.com/oauth2/v3/userinfo" rel="nofollow"><span>https</span><span>://</span><span>www</span><span>.</span><span>googleapis</span><span>.</span><span>com</span><span>/</span><span>oauth2</span><span>/</span><span>v3</span><span>/</span><span>userinfo</span></a>")
    if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
        return
	}
    defer email.Body.Close()
    data, _ := ioutil.ReadAll(email.Body)
    log.Println("Email body: ", string(data))
    c.Status(http.StatusOK)
}

func loginHandler(c *gin.Context) {
    state = randToken()
    session := sessions.Default(c)
    session.Set("state", state)
    session.Save()
    c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}

func main() {
    router := gin.Default()
    router.Use(sessions.Sessions("goquestsession", store))
    router.Static("/css", "./static/css")
    router.Static("/img", "./static/img")
    router.LoadHTMLGlob("templates/*")

    router.GET("/", indexHandler)
    router.GET("/login", loginHandler)
    router.GET("/auth", authHandler)

    router.Run("127.0.0.1:9090")
}
}