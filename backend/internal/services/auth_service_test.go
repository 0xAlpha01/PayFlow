package services

import (
	"github.com/ngwebwork/payflow/internal/utils"
	"testing"
)

func TestPasswordPolicy(t *testing.T) {
	if err := utils.ValidatePassword("SecurePassword123!"); err != nil {
		t.Fatal(err)
	}
	if err := utils.ValidatePassword("weakpassword"); err == nil {
		t.Fatal("expected weak password to fail")
	}
}
