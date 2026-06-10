package model

import "time"

type Inspector struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	EmployeeNo string   `json:"employee_no"`
	Role      string    `json:"role"`
	Processes []string  `json:"processes"`
	Status    string    `json:"status"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Delegation struct {
	ID            string    `json:"id"`
	DelegatorID   string    `json:"delegator_id"`
	DelegateeID   string    `json:"delegatee_id"`
	DelegatorRole string    `json:"delegator_role"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}
