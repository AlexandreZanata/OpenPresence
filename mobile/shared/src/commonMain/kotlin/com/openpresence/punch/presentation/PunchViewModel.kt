package com.openpresence.punch.presentation

import com.openpresence.punch.domain.DeviceFlag
import com.openpresence.punch.domain.PunchErrorCode
import com.openpresence.punch.domain.PunchRequest
import com.openpresence.punch.domain.PunchResult
import com.openpresence.punch.domain.PunchStatus
import com.openpresence.punch.domain.PunchType
import com.openpresence.punch.ports.BiometricProcessor
import com.openpresence.punch.ports.DeviceIntegrityChecker
import com.openpresence.punch.ports.GeofenceValidator
import com.openpresence.punch.ports.LocationProvider
import com.openpresence.punch.ports.PunchRepository
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class PunchViewModel(
    private val deviceChecker: DeviceIntegrityChecker,
    private val locationProvider: LocationProvider,
    private val geofenceValidator: GeofenceValidator,
    private val biometricProcessor: BiometricProcessor,
    private val repository: PunchRepository,
    private val scope: CoroutineScope,
    private val livenessThreshold: Float = 0.80f,
    private val deviceTimeIso: () -> String = { "1970-01-01T00:00:00Z" },
) {
    private val _state = MutableStateFlow<PunchState>(PunchState.Idle)
    val state: StateFlow<PunchState> = _state.asStateFlow()

    fun startPunch(type: PunchType) {
        scope.launch { runPunchFlow(type) }
    }

    fun handleOfflinePunch(request: PunchRequest) {
        scope.launch { queueOffline(request) }
    }

    fun syncOfflineQueue() {
        scope.launch { runOfflineSync() }
    }

    private suspend fun runOfflineSync() {
        if (!repository.isOnline()) {
            return
        }
        _state.value = PunchState.Submitting
        val result = repository.syncOfflineQueue()
        _state.value = when {
            result.failed > 0 -> PunchState.Error(
                PunchErrorCode.NETWORK,
                "Offline sync failed",
            )
            else -> PunchState.Idle
        }
    }

    private suspend fun runPunchFlow(type: PunchType) {
        try {
            _state.value = PunchState.CheckingDevice
            val deviceReport = deviceChecker.check()
            if (deviceReport.isRooted) {
                _state.value = PunchState.Error(
                    PunchErrorCode.DEVICE_BLOCKED,
                    "Device integrity check failed",
                )
                return
            }
            if (deviceReport.flags.isNotEmpty()) {
                _state.value = PunchState.DeviceWarning(deviceReport.flags)
            }

            _state.value = PunchState.WaitingLocation
            val location = locationProvider.currentLocation()
            if (location.isMocked) {
                _state.value = PunchState.Error(
                    PunchErrorCode.LOCATION_UNAVAILABLE,
                    "Mock location detected",
                )
                return
            }

            val geofence = geofenceValidator.validate(location)
            if (!geofence.isInside) {
                _state.value = PunchState.OutOfGeofence(geofence.distanceToNearestMeters)
                return
            }

            _state.value = PunchState.OpeningCamera
            val frame = biometricProcessor.captureFrame()
            _state.value = PunchState.DetectingFace(frameCount = 1)

            val liveness = biometricProcessor.livenessScore(frame)
            _state.value = PunchState.CheckingLiveness(liveness)
            if (liveness < livenessThreshold) {
                _state.value = PunchState.Error(
                    PunchErrorCode.SUBMIT_FAILED,
                    "Liveness below threshold",
                )
                return
            }

            val request = PunchRequest(
                type = type,
                location = location,
                frameBytes = frame,
                deviceReport = deviceReport,
                deviceTimeIso = deviceTimeIso(),
            )

            if (!repository.isOnline()) {
                queueOffline(request)
                return
            }

            _state.value = PunchState.Submitting
            val result = repository.submit(request)
            emitResult(result)
        } catch (_: Exception) {
            _state.value = PunchState.Error(PunchErrorCode.NETWORK, "Punch submission failed")
        }
    }

    private suspend fun queueOffline(request: PunchRequest) {
        _state.value = PunchState.Submitting
        repository.queueOffline(request)
        _state.value = PunchState.Success(
            PunchResult(
                id = "offline-pending",
                status = PunchStatus.PENDING,
                punchedAt = request.deviceTimeIso,
                type = request.type,
            ),
        )
    }

    private fun emitResult(result: PunchResult) {
        _state.value = when (result.status) {
            PunchStatus.VALID -> PunchState.Success(result)
            PunchStatus.SUSPICIOUS -> PunchState.Suspicious(result, "Awaiting manager review")
            PunchStatus.REJECTED -> PunchState.Error(PunchErrorCode.SUBMIT_FAILED, "Punch rejected")
            PunchStatus.PENDING -> PunchState.Success(result)
        }
    }

    fun reset() {
        _state.value = PunchState.Idle
    }
}
