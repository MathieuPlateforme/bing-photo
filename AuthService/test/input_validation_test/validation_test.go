package input_validation_test

import (
    "testing"
    "AuthService/utils"
)

func TestValidationFunctions(t *testing.T) {
    if utils.CheckEmail("invalid") {
        t.Error("Expected invalid email")
    }
    if !utils.CheckEmail("valid@example.com") {
        t.Error("Expected valid email")
    }

    if utils.CheckPassword("short") {
        t.Error("Expected weak password")
    }
    if !utils.CheckPassword("longenoughpassword") {
        t.Error("Expected strong password")
    }
}
