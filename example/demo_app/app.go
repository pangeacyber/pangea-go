package main

import (
	"context"
	"log"
	"os"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	// Pangea
	"github.com/pangeacyber/go-pangea/pangea"
	"github.com/pangeacyber/go-pangea/service/audit"
	"github.com/pangeacyber/go-pangea/service/embargo"
	"github.com/pangeacyber/go-pangea/service/redact"
)

type App struct {
	Router       *mux.Router
	pangea_token string

	store *DB

	// Pangea
	pangea_embargo *embargo.Embargo
	pangea_audit   *audit.Audit
	pangea_redact  *redact.Redact
}

type resp struct {
	Message string `json:"message"`
}

func (a *App) Initialize(pangea_token string) {
	a.pangea_token = pangea_token

	a.Router = mux.NewRouter()

	a.initializeRoutes()

	a.store = NewDB()

	csp := os.Getenv("PANGEA_CSP")

	embargoConfigID := os.Getenv("EMBARGO_CONFIG_ID")

	a.pangea_embargo = embargo.New(&pangea.Config{
		Token: a.pangea_token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    csp,
		},
		CfgToken: embargoConfigID,
	})

	auditConfigID := os.Getenv("AUDIT_CONFIG_ID")

	a.pangea_audit = audit.New(&pangea.Config{
		Token: a.pangea_token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    csp,
		},
		CfgToken: auditConfigID,
	})

	redactConfigID := os.Getenv("REDACT_CONFIG_ID")

	a.pangea_redact = redact.New(&pangea.Config{
		Token: a.pangea_token,
		EndpointConfig: &pangea.EndpointConfig{
			Scheme: "https",
			CSP:    csp,
		},
		CfgToken: redactConfigID,
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
	return
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
	eminput := &embargo.CheckInput{
		IP: pangea.String(client_ip),
	}

	checkOutput, _, err := a.pangea_embargo.Check(ctx, eminput)
	if err != nil {
		log.Println("[App.uploadResume] embargo check error: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if *checkOutput.Count > 0 {
		log.Println("[App.uploadResume] embargo check result positive: ", pangea.Stringify(checkOutput))

		// Audit log
		audinput := &audit.LogInput{
			Event: &audit.LogEventInput{
				Action:  pangea.String("add_employee"),
				Actor:   pangea.String(user),
				Target:  pangea.String(emp.PersonalEmail),
				Status:  pangea.String("error"),
				Message: pangea.String("Resume denied - sanctioned country from " + client_ip),
				Source:  pangea.String("web")},
			ReturnHash: pangea.Bool(true),
			Verbose:    pangea.Bool(false),
		}

		logOutput, _, err := a.pangea_audit.Log(ctx, audinput)
		if err != nil {
			log.Println("[App.uploadResume] audit log error: ", err.Error())
		} else {
			log.Println(pangea.Stringify(logOutput))
		}

		respondWithError(w, http.StatusForbidden, "Submissions from sanctioned country not allowed")
		return
	}

	// set status to CANDIDATE
	emp.Status = StatusCandidate

	// Redact
	out, err := json.Marshal(emp)
	if err != nil {
		log.Println("[App.uploadResume] failed to marshal emp to JSON")
	}

	redinput := &redact.StructuredInput{
		Data:   pangea.String(string(out)),
		Format: pangea.String("json"),
		Debug:  pangea.Bool(false),
	}

	redactOutput, _, err := a.pangea_redact.RedactStructured(ctx, redinput)
	redacted := ""
	if err != nil {
		log.Println("[App.uploadResume] redact error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(redactOutput.RedactedData))
		redacted = *redactOutput.RedactedData
		log.Println(redacted)
	}

	// Finally store to DB
	if err := a.store.addEmployee(emp); err != nil {
		// Audit log
		audinput := &audit.LogInput{
			Event: &audit.LogEventInput{
				Action:  pangea.String("add_employee"),
				Actor:   pangea.String(user),
				Target:  pangea.String(emp.PersonalEmail),
				Status:  pangea.String("error"),
				Message: pangea.String("Resume denied"),
				Source:  pangea.String("web")},
			ReturnHash: pangea.Bool(true),
			Verbose:    pangea.Bool(false),
		}

		logOutput, _, err1 := a.pangea_audit.Log(ctx, audinput)
		if err1 != nil {
			log.Println("[App.uploadResume] audit log error: ", err1.Error())
		} else {
			log.Println(pangea.Stringify(logOutput))
		}

		log.Println("[App.uploadResume] datastore error: ", err.Error())
		respondWithError(w, http.StatusInternalServerError, "Datastore error")
		return
	}

	resp := resp{
		Message: "Resume accepted",
	}

	// Audit log
	audinput := &audit.LogInput{
		Event: &audit.LogEventInput{
			Action:  pangea.String("add_employee"),
			Actor:   pangea.String(user),
			Target:  pangea.String(emp.PersonalEmail),
			Status:  pangea.String("success"),
			Message: pangea.String("Resume accepted"),
			New:     pangea.String(redacted),
			Source:  pangea.String("web")},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(false),
	}

	logOutput, _, err := a.pangea_audit.Log(ctx, audinput)
	if err != nil {
		log.Println("[App.uploadResume] audit log error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(logOutput))
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
		audinput := &audit.LogInput{
			Event: &audit.LogEventInput{
				Action:  pangea.String("lookup_employee"),
				Actor:   pangea.String(user),
				Target:  pangea.String(email),
				Status:  pangea.String("error"),
				Message: pangea.String("Requested employee record"),
				Source:  pangea.String("web")},
			ReturnHash: pangea.Bool(true),
			Verbose:    pangea.Bool(false),
		}

		logOutput, _, err1 := a.pangea_audit.Log(ctx, audinput)
		if err1 != nil {
			log.Println("[App.fetchEmployeeRecord] audit log error: ", err1.Error())
		} else {
			log.Println(pangea.Stringify(logOutput))
		}

		log.Println("[App.fetchEmployeeRecord] datastore error: ", err.Error())
		respondWithError(w, http.StatusNotFound, "Employee not found")
		return
	}

	// Audit log
	audinput := &audit.LogInput{
		Event: &audit.LogEventInput{
			Action:  pangea.String("lookup_employee"),
			Actor:   pangea.String(user),
			Target:  pangea.String(emp.PersonalEmail),
			Status:  pangea.String("success"),
			Message: pangea.String("Requested employee record"),
			Source:  pangea.String("web")},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(false),
	}

	logOutput, _, err := a.pangea_audit.Log(ctx, audinput)
	if err != nil {
		log.Println("[App.fetchEmployeeRecord] audit log error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(logOutput))
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
	audinput := &audit.LogInput{
		Event: &audit.LogEventInput{
			Action:  pangea.String("update_employee"),
			Actor:   pangea.String(user),
			Target:  pangea.String(input.PersonalEmail),
			Status:  pangea.String("success"),
			Message: pangea.String("Updated employee record"),
			Old:     pangea.String(string(outold)),
			New:     pangea.String(string(outnew)),
			Source:  pangea.String("web")},
		ReturnHash: pangea.Bool(true),
		Verbose:    pangea.Bool(false),
	}

	logOutput, _, err := a.pangea_audit.Log(ctx, audinput)
	if err != nil {
		log.Println("[App.fetchEmployeeRecord] audit log error: ", err.Error())
	} else {
		log.Println(pangea.Stringify(logOutput))
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
