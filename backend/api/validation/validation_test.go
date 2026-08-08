package validation

import "testing"

type sampleRequest struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidate_Passes(t *testing.T) {
	req := sampleRequest{Name: "Jane", Email: "jane@example.com"}
	if err := Validate(&req); err != nil {
		t.Errorf("expected valid struct to pass, got error: %v", err)
	}
}

func TestValidate_RejectsMissingRequired(t *testing.T) {
	req := sampleRequest{Email: "jane@example.com"}
	if err := Validate(&req); err == nil {
		t.Error("expected missing required field to fail validation")
	}
}

func TestValidate_RejectsInvalidEmail(t *testing.T) {
	req := sampleRequest{Name: "Jane", Email: "not-an-email"}
	if err := Validate(&req); err == nil {
		t.Error("expected invalid email to fail validation")
	}
}
