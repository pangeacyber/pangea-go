package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	// Pangea
	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/pangea"
	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/service/audit"
	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/service/embargo"
	"github.com/pangeacyber/pangea-go/packages/pangea-sdk/service/redact"
)

type App struct {
	Router       *mux.Router
	pangea_token string

	store *DB

	// Pangea
	embargo *embargo.Embargo
	audit   *audit.Audit
	redact  *redact.Redact
}

type resp struct {
	Message string `json:"message"`
}

func (a *App) Initialize(pangea_token string) {
	var err error
	a.pangea_token = pangea_token

	a.Router = mux.NewRouter()

	a.initializeRoutes()

	a.store = NewDB()

	a.embargo = embargo.New(&pangea.Config{
		Token:    a.pangea_token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		Insecure: false,
	})

	a.audit, err = audit.New(&pangea.Config{
		Token:    a.pangea_token,
		Domain:   os.Getenv("PANGEA_DOMAIN"),
		Insecure: false,
	})
	if err != nil {
		panic(err)
	}

	a.redact = redact.New(&pangea.Config{
		Token:  a.pangea_token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

}

func (a *App) Run(addr string) {
	log.Println("[App.Run] start...")

	http.ListenAndServe(addr, a.Router)

	log.Println("[App.Run] exit")
	a.shutdown()
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) setup(w http.ResponseWriter, r *http.Request) {
	log.Println("[App.setup handling request]")

	err := a.store.setupEmployeeTable()

	resp := resp{
		Message: "App setup complete",
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (a *App) uploadResume(w http.ResponseWriter, r *http.Request) {
	// TODO: Bearer Token
	user, _, ok := r.BasicAuth()

	if !ok {
		log.Println("[App.uploadResume] Failed to parse Basic Auth")
		respondWithError(w, http.StatusUnauthorized, "Missing auth info")
		return
	}

	var emp employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&emp); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// obtain "ClientAPAddress" from header
	client_ip := r.Header.Get("ClientIPAddress")

	log.Println("[App.uploadResume] processing request from IP: ", client_ip)

	// Check Embargo
	ctx := context.Background()
	eminput := &embargo.IPCheckInput{
		IP: pangea.String(client_ip),
	}

	checkOutput, err := a.embargo.IPCheck(ctx, eminput)
	if err != nil {
		log.Println("[App.uploadResume] embargo check error: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if checkOutput.Result.Count > 0 {
		log.Println("[App.uploadResume] embargo check result positive: ", pangea.Stringify(checkOutput))

		// Audit log
		event := audit.Event{
			Action:  "add_employee",
			Actor:   user,
			Target:  emp.PersonalEmail,
			Status:  "error",
			Message: "Resume denied - sanctioned country from " + client_ip,
			Source:  "web",
		}

		logResponse, err := a.audit.Log(ctx, event, true, true)
		if err != nil {
			log.Println("[App.uploadResume] audit log error: ", err.Error())
		} else {
			log.Println(pangea.Stringify(logResponse.Result))
		}

		respondWithError(w, http.StatusForbidden, "Submissions from sanctioned country not allowed")
		return
	}

	// FIXME: What is StatusCandidate? Should we remove it?
	// set status to CANDIDATE
	emp.Status = StatusCandidate

	// Redact
	redinput := &redact.StructuredInput{
		Format: pangea.String("json"),
		Debug:  pangea.Bool(false),
	}
	redinput.SetData(emp)

	redactResponse, err := a.redact.RedactStructured(ctx, redinput)
	redacted := ""
	if err != nil {
		log.Println("[App.uploadResume] redact error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(redactResponse.Result.RedactedData))
		// FIXME: What about employee?
		var redactedEmployee employee
		redactResponse.Result.GetRedactedData(&redactedEmployee)
		log.Printf("%+v", redactedEmployee)
	}

	// Finally store to DB
	if err := a.store.addEmployee(emp); err != nil {
		// Audit log
		event := audit.Event{
			Action:  "add_employee",
			Actor:   user,
			Target:  emp.PersonalEmail,
			Status:  "error",
			Message: "Resume denied",
			Source:  "web",
		}

		logResponse, err1 := a.audit.Log(ctx, event, true, true)
		if err1 != nil {
			log.Println("[App.uploadResume] audit log error: ", err1.Error())
		} else {
			log.Println(pangea.Stringify(logResponse.Result))
		}

		log.Println("[App.uploadResume] datastore error: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Datastore error")
		return
	}

	resp := resp{
		Message: "Resume accepted",
	}

	// Audit log
	event := audit.Event{
		Action:  "add_employee",
		Actor:   user,
		Target:  emp.PersonalEmail,
		Status:  "success",
		Message: "Resume accepted",
		New:     redacted,
		Source:  "web",
	}

	logResponse, err := a.audit.Log(ctx, event, true, true)
	if err != nil {
		log.Println("[App.uploadResume] audit log error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(logResponse.Result))
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func (a *App) fetchEmployeeRecord(w http.ResponseWriter, r *http.Request) {
	// TODO: Bearer Token
	user, _, ok := r.BasicAuth()

	ctx := context.Background()

	if !ok {
		log.Println("[App.uploadResume] Failed to parse Basic Auth")
		respondWithError(w, http.StatusUnauthorized, "Missing auth info")
		return
	}

	// get the email
	vars := mux.Vars(r)
	email := vars["email"]

	log.Println("[App.fetchEmployeeRecord] Processing input from user ", user, " ", email)

	emp, err := a.store.lookupEmployee(email)
	if err != nil {
		// Audit log
		event := audit.Event{
			Action:  "lookup_employee",
			Actor:   user,
			Target:  email,
			Status:  "error",
			Message: "Requested employee record",
			Source:  "web",
		}

		logResponse, err1 := a.audit.Log(ctx, event, true, true)
		if err1 != nil {
			log.Println("[App.fetchEmployeeRecord] audit log error: ", err1.Error())
		} else {
			log.Println(pangea.Stringify(logResponse.Result))
		}

		log.Println("[App.fetchEmployeeRecord] datastore error: ", err.Error())
		respondWithError(w, http.StatusNotFound, "Employee not found")
		return
	}

	// Audit log
	event := audit.Event{
		Action:  "lookup_employee",
		Actor:   user,
		Target:  emp.PersonalEmail,
		Status:  "success",
		Message: "Requested employee record",
		Source:  "web",
	}

	logResponse, err := a.audit.Log(ctx, event, true, true)
	if err != nil {
		log.Println("[App.fetchEmployeeRecord] audit log error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(logResponse.Result))
	}

	respondWithJSON(w, http.StatusOK, emp)
}

func (a *App) updateEmployee(w http.ResponseWriter, r *http.Request) {
	// TODO: Bearer Token
	user, _, ok := r.BasicAuth()

	ctx := context.Background()

	if !ok {
		log.Println("[App.updateEmployee] Failed to parse Basic Auth")
		respondWithError(w, http.StatusUnauthorized, "Missing auth info")
		return
	}

	var input employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		log.Println("[App.updateEmployee] failed to decode: ", err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// fetch employee record
	oldemp, err := a.store.lookupEmployee(input.PersonalEmail)
	if err != nil {
		log.Println("[App.updateEmployee] Employee not found ", input.PersonalEmail, ", ", err.Error())
		respondWithError(w, http.StatusNotFound, "Employee not found")
		return
	}

	// These are provided in input: StartDate, TermDate, ManagerId, Department, Salary, Status, CompanyEmail
	// copy the rest from oldemp - for Audit logging purposes
	input.ID = oldemp.ID
	input.FirstName = oldemp.FirstName
	input.LastName = oldemp.LastName
	input.Phone = oldemp.Phone
	input.DateOfBirth = oldemp.DateOfBirth
	input.Medical = oldemp.Medical
	input.ProfilePicture = oldemp.ProfilePicture
	input.DLPicture = oldemp.DLPicture
	input.SSN = oldemp.SSN

	// update the record
	err = a.store.updateEmployee(input)

	if err != nil {
		log.Println("[App.updateEmployee] Database update error: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Database update error")
		return
	}

	resp := resp{
		Message: "Successfully updated employee record",
	}

	outold, err := json.Marshal(oldemp)
	if err != nil {
		log.Println("[App.updateEmployee] failed to marshal oldemp to JSON")
	}

	outnew, err := json.Marshal(input)
	if err != nil {
		log.Println("[App.updateEmployee] failed to marshal input to JSON")
	}

	// Audit log
	event := audit.Event{
		Action:  "update_employee",
		Actor:   user,
		Target:  input.PersonalEmail,
		Status:  "success",
		Message: "Updated employee record",
		Old:     string(outold),
		New:     string(outnew),
		Source:  "web",
	}

	logResponse, err := a.audit.Log(ctx, event, true, true)
	if err != nil {
		log.Println("[App.fetchEmployeeRecord] audit log error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(logResponse.Result))
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (a *App) shutdown() {
	a.store.Close()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/setup", a.setup).Methods("POST")
	a.Router.HandleFunc("/upload_resume", a.uploadResume).Methods("POST")
	a.Router.HandleFunc("/employee/{email}", a.fetchEmployeeRecord).Methods("GET")
	a.Router.HandleFunc("/update_employee", a.updateEmployee).Methods("POST")
}
