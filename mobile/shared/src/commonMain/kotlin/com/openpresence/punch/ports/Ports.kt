package com.openpresence.punch.ports

import com.openpresence.punch.domain.DeviceIntegrityReport
import com.openpresence.punch.domain.GpsCoordinate
import com.openpresence.punch.domain.OfflineSyncResult
import com.openpresence.punch.domain.PunchRequest
import com.openpresence.punch.domain.PunchResult

interface DeviceIntegrityChecker {
    suspend fun check(): DeviceIntegrityReport
}

interface LocationProvider {
    suspend fun currentLocation(): GpsCoordinate
}

interface GeofenceValidator {
    suspend fun validate(location: GpsCoordinate): GeofenceValidation
}

data class GeofenceValidation(
    val isInside: Boolean,
    val distanceToNearestMeters: Double,
)

interface BiometricProcessor {
    suspend fun captureFrame(): ByteArray
    suspend fun livenessScore(frame: ByteArray): Float
}

interface PunchRepository {
    suspend fun submit(request: PunchRequest): PunchResult
    suspend fun queueOffline(request: PunchRequest)
    suspend fun isOnline(): Boolean
    suspend fun syncOfflineQueue(): OfflineSyncResult
}

interface PunchApi {
    suspend fun submit(request: PunchRequest): PunchResult
}
