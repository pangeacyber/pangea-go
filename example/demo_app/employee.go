package main

const (
    StatusUnknown    = 1
    StatusCandidate  = 2
    StatusFullTime   = 3
    StatusContractor = 4
    StatusTerminated = 5
)

type employee struct {
    ID             int    `json:"id"`
    FirstName      string `json:"first_name"`
    LastName       string `json:"last_name"`
    CompanyEmail   string `json:"company_email"`
    PersonalEmail  string `json:"email"`
    Phone          string `json:"phone"`
    DateOfBirth    string `json:"dob"`
    StartDate      string `json:"start_date"`
    TermDate       string `json:"term_date"`
    Department     string `json:"department"`
    ManagerId      int    `json:"manager_id"`
    Salary         int    `json:"salary"`
    Medical        string `json:"medical"`
    ProfilePicture string `json:"profile_pic"`
    DLPicture      string `json:"dl_pic"`
    SSN            string `json:"ssn"`
    Status         int    `json:status"`
}
