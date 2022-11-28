package data

import "testing"

func RecoveryEmail_Table_Test(t testing.T) {
	rec := &RecoveryEmail{}
	if rec.Table() != "recovery_emails" {
		t.Error("incorrect table reported")
	}
}
