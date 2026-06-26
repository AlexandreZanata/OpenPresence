package punch

// isValidSequence enforces BR-014 punch order for the current workday history.
func isValidSequence(history []PunchRecord, next PunchType) bool {
	if len(history) == 0 {
		return next == PunchTypeClockIn
	}
	last := history[len(history)-1]
	switch last.Type {
	case PunchTypeClockIn:
		return next == PunchTypeBreakStart || next == PunchTypeClockOut
	case PunchTypeBreakStart:
		return next == PunchTypeBreakEnd
	case PunchTypeBreakEnd:
		return next == PunchTypeClockOut
	case PunchTypeClockOut:
		return false
	default:
		return false
	}
}

// CanTransition reports whether a status change is allowed (see docs/DOMAIN-MODEL.md).
func CanTransition(from, to PunchStatus) bool {
	switch from {
	case PunchStatusPending:
		return to == PunchStatusValid || to == PunchStatusDiscarded || to == PunchStatusSuspicious
	case PunchStatusSuspicious:
		return to == PunchStatusValid || to == PunchStatusRejected
	case PunchStatusValid, PunchStatusRejected, PunchStatusDiscarded:
		return false
	default:
		return false
	}
}

// Transition moves status forward or returns ErrInvalidTransition.
func Transition(from, to PunchStatus) (PunchStatus, error) {
	if !CanTransition(from, to) {
		return from, ErrInvalidTransition
	}
	return to, nil
}
