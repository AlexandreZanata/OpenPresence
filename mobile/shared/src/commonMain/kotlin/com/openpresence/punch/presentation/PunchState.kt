package com.openpresence.punch.presentation

import com.openpresence.punch.domain.DeviceFlag
import com.openpresence.punch.domain.PunchErrorCode
import com.openpresence.punch.domain.PunchResult

sealed class PunchState {
    data object Idle : PunchState()

    data object CheckingDevice : PunchState()

    data class DeviceWarning(val flags: List<DeviceFlag>) : PunchState()

    data object WaitingLocation : PunchState()

    data class OutOfGeofence(val distanceMeters: Double) : PunchState()

    data object OpeningCamera : PunchState()

    data class DetectingFace(val frameCount: Int) : PunchState()

    data class CheckingLiveness(val score: Float) : PunchState()

    data object Submitting : PunchState()

    data class Success(val punch: PunchResult) : PunchState()

    data class Suspicious(val punch: PunchResult, val reason: String) : PunchState()

    data class Error(val code: PunchErrorCode, val message: String) : PunchState()
}
