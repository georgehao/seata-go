/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
