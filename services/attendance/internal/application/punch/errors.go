package punch

import "errors"

var (
	ErrEmployeeNotFound   = errors.New("punch: employee not found")
	ErrEmployeeInactive   = errors.New("punch: employee is not active")
	ErrNoActivePlacement  = errors.New("punch: no active primary placement")
	ErrPunchTypeNotAllowed = errors.New("punch: punch type not allowed by policy")
	ErrBiometricRequired  = errors.New("punch: biometric verification required")
	ErrDeviceLocked       = errors.New("punch: device locked due to repeated rejections")
)
