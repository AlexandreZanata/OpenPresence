package com.openpresence.punch.domain

enum class PunchType {
    CLOCK_IN,
    CLOCK_OUT,
    BREAK_START,
    BREAK_END,
}

enum class PunchStatus {
    VALID,
    SUSPICIOUS,
    REJECTED,
    PENDING,
}

enum class PunchErrorCode {
    DEVICE_BLOCKED,
    LOCATION_UNAVAILABLE,
    NETWORK,
    SUBMIT_FAILED,
    UNKNOWN,
}
