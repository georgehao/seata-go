package constant

// Status of branch transaction
const (
	// BranchStatusUnknown description:BranchStatus_Unknown branch status.
	BranchStatusUnknown = iota

	// BranchStatusRegistered description:BranchStatus_Registered to TC.
	BranchStatusRegistered

	// BranchStatusPhaseOneDone description:Branch logic is successfully done at phase one.
	BranchStatusPhaseOneDone

	// BranchStatusPhaseOneFailed description:Branch logic is failed at phase one.
	BranchStatusPhaseOneFailed

	// BranchStatusPhaseOneTimeout description:Branch logic is NOT reported for a timeout.
	BranchStatusPhaseOneTimeout

	// BranchStatusPhaseTwoCommitted description:Commit logic is successfully done at phase two.
	BranchStatusPhaseTwoCommitted

	// BranchStatusPhaseTwoCommitFailedRetryable description:Commit logic is failed but retryable.
	BranchStatusPhaseTwoCommitFailedRetryable

	// BranchStatusPhaseTwoCommitFailedUnretryable description:Commit logic is failed and NOT retryable.
	BranchStatusPhaseTwoCommitFailedUnretryable

	// BranchStatusPhaseTwoRollbacked description:Rollback logic is successfully done at phase two.
	BranchStatusPhaseTwoRollbacked

	// BranchStatusPhaseTwoRollbackFailedRetryable description:Rollback logic is failed but retryable.
	BranchStatusPhaseTwoRollbackFailedRetryable

	// BranchStatusPhaseTwoRollbackFailedUnretryable description:Rollback logic is failed but NOT retryable.
	BranchStatusPhaseTwoRollbackFailedUnretryable

	// BranchStatusPhaseTwoCommitFailedXAERNOTARetryable description:Commit logic is failed because of XAException.XAER_NOTA but retryable.
	BranchStatusPhaseTwoCommitFailedXAERNOTARetryable

	// PhaseTwoRollbackFailedXAERNOTARetryable description:rollback logic is failed because of XAException.XAER_NOTA but retryable.
	PhaseTwoRollbackFailedXAERNOTARetryable
)
