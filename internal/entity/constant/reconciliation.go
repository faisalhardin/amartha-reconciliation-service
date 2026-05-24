package constant

const (
	ReconciliationStatusMatched       = "MATCHED"
	ReconciliationStatusMismatch      = "MISMATCH"
	ReconciliationStatusOnlyInSystem  = "ONLY_IN_SYSTEM"
	ReconciliationStatusOnlyInBank    = "ONLY_IN_BANK"
	ReconciliationStatusAmbiguous     = "AMBIGUOUS"

	ReconciliationMatchMethodReference = "REFERENCE"
	ReconciliationMatchMethodComposite = "COMPOSITE"
	ReconciliationMatchMethodNone      = "NONE"
)
