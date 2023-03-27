package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	repository "github.com/Nickolaygoloburdin/droneapp/internal/database"
	"github.com/Nickolaygoloburdin/droneapp/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

type DBHandler struct {
	ctx  context.Context
	repo *repository.Repository
}

func NewDBHandler(ctx context.Context, repo *repository.Repository) *DBHandler {
	return &DBHandler{ctx, repo}
}

func (dbh DBHandler) handlePage(writer http.ResponseWriter, request *http.Request) {
	_, err := generateJWT()
	if err != nil {
		log.Fatalln("Error generating JWT", err)
	}

	writer.Header().Set("Token", "%v")
	type_ := "application/json"
	writer.Header().Set("Content-Type", type_)
	var message Message
	err = json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		return
	}
	err = json.NewEncoder(writer).Encode(message)
	if err != nil {
		return
	}
}

func (dbh DBHandler) Signup(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("first_name")) + strings.TrimSpace(r.FormValue("last_name"))

	// if name == "" || surname == "" || login == "" || password == "" {
	// 	a.SignupPage(rw, "Все поля должны быть заполнены!")
	// 	return
	// }

	// if password != password2 {
	// 	a.SignupPage(rw, "Пароли не совпадают! Попробуйте еще")
	// 	return
	// }

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	err := a.repo.AddNewUser(a.ctx, username, email, hashedPass)
	if err != nil {
		//a.SignupPage(rw, fmt.Sprintf("Ошибка создания пользователя: %v", err))
		return
	}

}

var sampleSecretKey = "hellofromback"

func (dbh DBHandler) Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	var user models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	login := user.Username
	password := user.Password
	email := user.Email

	if login == "" || password == "" {
		//a.LoginPage(rw, "Необходимо указать логин и пароль!")
		return
	}

	return c.JSON(LoginResponse{AccessToken: t})

	a.cache[hashedToken] = dbuser

}
func (dbh DBHandler) Signup(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("first_name")) + strings.TrimSpace(r.FormValue("last_name"))

	// if name == "" || surname == "" || login == "" || password == "" {
	// 	a.SignupPage(rw, "Все поля должны быть заполнены!")
	// 	return
	// }

	// if password != password2 {
	// 	a.SignupPage(rw, "Пароли не совпадают! Попробуйте еще")
	// 	return
	// }

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	err := a.repo.AddNewUser(a.ctx, username, email, hashedPass)
	if err != nil {
		//a.SignupPage(rw, fmt.Sprintf("Ошибка создания пользователя: %v", err))
		return
	}
}
func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "username"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "Signing Error", err
	}

	return tokenString, nil
}

// comment these
func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized"))
					if err != nil {
						return nil, err
					}
				}
				return "", nil

			})
			// parsing errors result
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}

			}
			// if there's a token
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
		// response for if there's no token header
	})
}

func extractClaims(_ http.ResponseWriter, request *http.Request) (string, error) {
	if request.Header["Token"] != nil {
		tokenString := request.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return sampleSecretKey, nil
		})
		if err != nil {
			return "Error Parsing Token: ", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			return username, nil
		}
	}

	return "unable to extract claims", nil
}

func authPage() {
	token, _ := generateJWT()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	req.Header.Set("Token", token)
	_, _ = client.Do(req)
}
