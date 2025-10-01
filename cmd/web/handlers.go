package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aakash-1857/codebin/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
)

func (app *application) healthcheck(w http.ResponseWriter,r *http.Request){
	data:=map[string]string{
		"status":"available",
		"environment":"development",
		"version":"1.0.0",
	}
	err:=app.writeJSON(w,http.StatusOK,data,nil)
	if err!=nil{
		http.Error(w,"Internal Server Error",http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter,r *http.Request){
	id := chi.URLParam(r,"id")
	// check if snippet is in our cache
	cachedSnippet,found:=app.cache.Get(id)
	if found{
		app.logger.Info("cache hit for snippet","id",id)
		err:=app.writeJSON(w,http.StatusOK,cachedSnippet,nil)
		if err!=nil{
			app.serverErrorResponse(w,r,err)
		}
		return
	}
	app.logger.Info("cache miss for snippet","id",id)
	snippet,err:=app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// Use a specific helper for 404s (we'll add it next).
			app.notFoundResponse(w, r)
		} else {
			// Use the generic server error helper.
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.cache.Set(id,snippet,cache.DefaultExpiration)
	err = app.writeJSON(w,http.StatusOK,snippet,nil)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
	}

}
func (app *application) snippetCreate(w http.ResponseWriter,r *http.Request){
	var input struct {
		Title string `json:"title"`
		Content string `json:"content"`
	}
	err:=json.NewDecoder(r.Body).Decode(&input)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return
	}
	if input.Title == "" || input.Content == "" {
		app.serverErrorResponse(w,r,err)
		return
	}
	id,err:=app.snippets.Insert(input.Title,input.Content)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return
	}
	snippet,err:=app.snippets.Get(id)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return
	}
	err = app.writeJSON(w,http.StatusCreated,snippet,nil)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
	}
}
func (app *application) snippetLatest(w http.ResponseWriter,r *http.Request){
	snippets,err:=app.snippets.Latest()
	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return 
	}
	err=app.writeJSON(w,http.StatusOK,snippets,nil)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Create a temporary user struct just to use the password hashing method.
	user := &models.User{}
	err = user.SetPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// MODIFIED: Call the new, explicit Insert method.
	// It now returns a fully populated user object.
	createdUser, err := app.users.Insert(input.Name, input.Email, user.PasswordHash)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the returned user object back to the client.
	err = app.writeJSON(w, http.StatusCreated, createdUser, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter,r *http.Request){
	var input struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	err:=json.NewDecoder(r.Body).Decode(&input)
	if err!=nil{
		app.badRequestResponse(w,r,err)
		return
	}
	// Authentication
	// 1. Look at user by email
	user,err:=app.users.GetByEmail(input.Email)
	if err!=nil{
		app.invalidCredentialsResponse(w,r)
		return
	}
	// 2. Check if the provided password matches the stored hash
	match,err:=user.MatchesPassword(input.Password)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return
	}
	if !match{
		app.invalidCredentialsResponse(w,r)
		return
	}
	// 3. If credentials are correct, create new JWT
	token:=jwt.New(jwt.SigningMethodHS256)

	claims:=token.Claims.(jwt.MapClaims)
	claims["sub"]=user.ID// sub (subject) is standard claim for user ID
	claims["exp"]=time.Now().Add(24*time.Hour).Unix() //token valid for 24 hrs

	// 4. Sign the token with our secrety key
	signedToken,err:=token.SignedString([]byte(app.config.jwt.secret))
	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return 
	}
	// 5. send token back to cliend
	err = app.writeJSON(w, http.StatusOK,map[string]string{"token":signedToken},nil)
	if err!=nil{
		app.serverErrorResponse(w,r,err)
	}
}