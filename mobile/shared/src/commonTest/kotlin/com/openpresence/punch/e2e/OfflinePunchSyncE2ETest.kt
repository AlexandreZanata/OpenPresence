package com.openpresence.punch.e2e

import com.openpresence.punch.data.OfflinePunchRepository
import com.openpresence.punch.domain.DeviceIntegrityReport
import com.openpresence.punch.domain.GpsCoordinate
import com.openpresence.punch.domain.PunchRequest
import com.openpresence.punch.domain.PunchResult
import com.openpresence.punch.domain.PunchStatus
import com.openpresence.punch.domain.PunchType
import com.openpresence.punch.ports.BiometricProcessor
import com.openpresence.punch.ports.DeviceIntegrityChecker
import com.openpresence.punch.ports.GeofenceValidation
import com.openpresence.punch.ports.GeofenceValidator
import com.openpresence.punch.ports.LocationProvider
import com.openpresence.punch.ports.PunchApi
import com.openpresence.punch.presentation.PunchState
import com.openpresence.punch.presentation.PunchViewModel
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertIs

@OptIn(ExperimentalCoroutinesApi::class)
class OfflinePunchSyncE2ETest {
    @Test
    fun offlineQueue_syncsToMockApiWhenOnline() = runTest {
        var online = false
        val api = RecordingPunchApi()
        val repository = OfflinePunchRepository(api) { online }
        val request = sampleRequest()

        repository.queueOffline(request)
        assertEquals(1, repository.pendingCount())

        val syncWhileOffline = repository.syncOfflineQueue()
        assertEquals(0, syncWhileOffline.synced)
        assertEquals(1, syncWhileOffline.pending)

        online = true
        val syncResult = repository.syncOfflineQueue()
        assertEquals(1, syncResult.synced)
        assertEquals(0, syncResult.pending)
        assertEquals(1, api.submitCount)
        assertEquals(PunchType.CLOCK_IN, api.lastRequest?.type)
    }

    @Test
    fun startPunch_E2E_offlineQueuesThenSyncWhenOnline() = runTest {
        var online = false
        val api = RecordingPunchApi()
        val repository = OfflinePunchRepository(api) { online }
        val viewModel = PunchViewModel(
            deviceChecker = StubDeviceChecker(),
            locationProvider = StubLocationProvider(),
            geofenceValidator = StubGeofenceValidator(),
            biometricProcessor = StubBiometricProcessor(),
            repository = repository,
            scope = this,
            deviceTimeIso = { "2026-06-26T08:01:00Z" },
        )

        online = false
        viewModel.startPunch(PunchType.CLOCK_IN)
        advanceUntilIdle()
        assertIs<PunchState.Success>(viewModel.state.value)
        assertEquals(1, repository.pendingCount())
        assertEquals(0, api.submitCount)

        online = true
        viewModel.syncOfflineQueue()
        advanceUntilIdle()
        assertEquals(PunchState.Idle, viewModel.state.value)
        assertEquals(0, repository.pendingCount())
        assertEquals(1, api.submitCount)
    }

    private fun sampleRequest() = PunchRequest(
        type = PunchType.CLOCK_IN,
        location = GpsCoordinate(-12.5458, -55.7061, 12.0),
        frameBytes = byteArrayOf(1, 2, 3),
        deviceReport = DeviceIntegrityReport(false, false, false),
        deviceTimeIso = "2026-06-26T08:00:00Z",
    )
}

private class RecordingPunchApi : PunchApi {
    var submitCount = 0
    var lastRequest: PunchRequest? = null

    override suspend fun submit(request: PunchRequest): PunchResult {
        submitCount++
        lastRequest = request
        return PunchResult(
            id = "api-punch-$submitCount",
            status = PunchStatus.VALID,
            punchedAt = "2026-06-26T08:01:02Z",
            type = request.type,
        )
    }
}

private class StubDeviceChecker : DeviceIntegrityChecker {
    override suspend fun check() = DeviceIntegrityReport(false, false, false)
}

private class StubLocationProvider : LocationProvider {
    override suspend fun currentLocation() = GpsCoordinate(-12.5458, -55.7061, 12.0)
}

private class StubGeofenceValidator : GeofenceValidator {
    override suspend fun validate(location: GpsCoordinate) =
        GeofenceValidation(isInside = true, distanceToNearestMeters = 0.0)
}

private class StubBiometricProcessor : BiometricProcessor {
    override suspend fun captureFrame(): ByteArray = byteArrayOf(1, 2, 3)
    override suspend fun livenessScore(frame: ByteArray): Float = 0.95f
}
