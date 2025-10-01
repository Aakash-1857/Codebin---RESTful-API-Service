// cmd/web/helpers.go
package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter,status int, data interface{},headers http.Header) error{
	js,err:=json.Marshal(data)
	if err!=nil{
		return err
	}
	for key,value := range headers{
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
func (app *application) logError(r *http.Request,err error){
	app.logger.Error(err.Error(),"request_method",r.Method,"request_url",r.URL.String())
}
func (app *application) errorResponse(w http.ResponseWriter,r *http.Request,status int,message interface{}){
	env:=map[string]interface{}{"error":message}
	err:=app.writeJSON(w,status,env,nil)
	if err!=nil{
		app.logError(r,err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
// cmd/web/helpers.go

// serverErrorResponse logs the detailed error and sends a generic 500 response.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse sends a 404 Not Found response.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// badRequestResponse sends a 400 Bad Request response.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
func (app *application) invalidCredentialsResponse(w http.ResponseWriter,r *http.Request){
	message:="invalid authentication credentials"
	app.errorResponse(w,r,http.StatusUnauthorized,message)
}
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer") // Inform client how to authenticate.
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}