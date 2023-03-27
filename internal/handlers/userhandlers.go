package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

var sampleSecretKey = "hellofromback"

func (dbh DBHandler) Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var user models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return
	}
	email := user.Account.Email
	password := user.Account.Password

	if email == "" || password == "" {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		rw.Write([]byte("You need to fill all fields"))
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	dbuser, err := dbh.repo.Login(dbh.ctx, email, hashedPass)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		rw.Write([]byte("You didnt sign up"))
	}

	token, err := generateJWT(models.TokenData{Name: dbuser.Name, Surname: dbuser.Surname, Email: dbuser.Email})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		rw.Write([]byte("Error generate JWT"))
		return
	}
	var resp models.SignupResponse
	resp.User.Token = token
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(resp)

}
func (dbh DBHandler) Signup(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var user models.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Account.Email == "" || user.Account.Name == "" || user.Account.Surname == "" || user.Account.Password == "" {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		rw.Write([]byte("You need to fill all fields"))
		return
	}

	hash := md5.Sum([]byte(user.Account.Password))
	hashedPass := hex.EncodeToString(hash[:])

	err = dbh.repo.AddNewUser(dbh.ctx, user.Account.Name+user.Account.Surname, user.Account.Email, hashedPass)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		rw.Write([]byte("Error in database"))
		return
	}
	tokenData := models.TokenData{Name: user.Account.Name, Surname: user.Account.Surname, Email: user.Account.Email}
	token, err := generateJWT(tokenData)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		rw.Write([]byte("Error generate JWT"))
		return
	}
	var resp models.SignupResponse
	resp.User.Token = token
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(resp)

}
func generateJWT(user models.TokenData) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(40 * time.Minute)
	claims["authorized"] = true
	claims["first_name"] = user.Name
	claims["last_name"] = user.Surname
	claims["email"] = user.Email

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "Signing Error", err
	}

	return tokenString, nil
}

// comment these
func VerifyJWT(endpointHandler httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter,
		request *http.Request, ps httprouter.Params) {
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
				endpointHandler(writer, request, ps)
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
	}
}

func extractClaims(_ http.ResponseWriter, request *http.Request) (td models.TokenData, err error) {
	if request.Header["Token"] != nil {
		tokenString := request.Header["Token"][0]
		token, error := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return sampleSecretKey, nil
		})
		if error != nil {
			err = error
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {

			td.Name = claims["first_name"].(string)
			td.Surname = claims["last_name"].(string)
			td.Email = claims["email"].(string)
			return
		}
	}
	err = errors.New("Cannot extract claims")
	return
}

func GetDataFromToken(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	resp, err := extractClaims(rw, r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(resp)
}
