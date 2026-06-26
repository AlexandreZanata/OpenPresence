package com.openpresence.punch.presentation

import com.openpresence.punch.domain.DeviceIntegrityReport
import com.openpresence.punch.domain.GpsCoordinate
import com.openpresence.punch.domain.PunchErrorCode
import com.openpresence.punch.domain.PunchRequest
import com.openpresence.punch.domain.PunchResult
import com.openpresence.punch.domain.PunchStatus
import com.openpresence.punch.domain.PunchType
import com.openpresence.punch.ports.BiometricProcessor
import com.openpresence.punch.ports.DeviceIntegrityChecker
import com.openpresence.punch.ports.GeofenceValidation
import com.openpresence.punch.ports.GeofenceValidator
import com.openpresence.punch.ports.LocationProvider
import com.openpresence.punch.ports.PunchRepository
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertIs

@OptIn(ExperimentalCoroutinesApi::class)
class PunchViewModelTest {
    @Test
    fun startPunch_emitsCheckingDeviceFirst() = runTest {
        lateinit var viewModel: PunchViewModel
        val device = object : DeviceIntegrityChecker {
            override suspend fun check(): DeviceIntegrityReport {
                assertIs<PunchState.CheckingDevice>(viewModel.state.value)
                return DeviceIntegrityReport(
                    isRooted = false,
                    isVpnActive = false,
                    isDeveloperOptionsEnabled = false,
                )
            }
        }
        viewModel = viewModel(scope = this, device = device)
        viewModel.startPunch(PunchType.CLOCK_IN)
        advanceUntilIdle()
    }

    @Test
    fun startPunch_outOfGeofence_whenValidatorFails() = runTest {
        val viewModel = viewModel(
            scope = this,
            geofence = FakeGeofenceValidator(inside = false, distance = 42.5),
        )
        viewModel.startPunch(PunchType.CLOCK_IN)
        advanceUntilIdle()
        val state = viewModel.state.value
        assertIs<PunchState.OutOfGeofence>(state)
        assertEquals(42.5, state.distanceMeters)
    }

    @Test
    fun startPunch_success_whenRepositoryReturnsValid() = runTest {
        val viewModel = viewModel(scope = this)
        viewModel.startPunch(PunchType.CLOCK_IN)
        advanceUntilIdle()
        val state = viewModel.state.value
        assertIs<PunchState.Success>(state)
        assertEquals(PunchStatus.VALID, state.punch.status)
    }

    @Test
    fun startPunch_error_onNetworkFailure() = runTest {
        val viewModel = viewModel(
            scope = this,
            repository = FakePunchRepository(failSubmit = true),
        )
        viewModel.startPunch(PunchType.CLOCK_IN)
        advanceUntilIdle()
        val state = viewModel.state.value
        assertIs<PunchState.Error>(state)
        assertEquals(PunchErrorCode.NETWORK, state.code)
    }

    @Test
    fun handleOfflinePunch_queuesPendingResult() = runTest {
        val repository = FakePunchRepository()
        val viewModel = viewModel(scope = this, repository = repository)
        val request = PunchRequest(
            type = PunchType.CLOCK_IN,
            location = GpsCoordinate(0.0, 0.0, 10.0),
            frameBytes = byteArrayOf(1),
            deviceReport = DeviceIntegrityReport(false, false, false),
            deviceTimeIso = "2026-06-26T08:00:00Z",
        )
        viewModel.handleOfflinePunch(request)
        advanceUntilIdle()
        assertEquals(1, repository.offlineCount)
        assertIs<PunchState.Success>(viewModel.state.value)
        assertEquals(PunchStatus.PENDING, (viewModel.state.value as PunchState.Success).punch.status)
    }

    private fun viewModel(
        scope: kotlinx.coroutines.CoroutineScope,
        device: DeviceIntegrityChecker = FakeDeviceChecker(),
        geofence: GeofenceValidator = FakeGeofenceValidator(inside = true, distance = 0.0),
        repository: PunchRepository = FakePunchRepository(),
    ): PunchViewModel = PunchViewModel(
        deviceChecker = device,
        locationProvider = FakeLocationProvider(),
        geofenceValidator = geofence,
        biometricProcessor = FakeBiometricProcessor(),
        repository = repository,
        scope = scope,
        deviceTimeIso = { "2026-06-26T08:01:00Z" },
    )
}

private class FakeDeviceChecker : DeviceIntegrityChecker {
    override suspend fun check() = DeviceIntegrityReport(
        isRooted = false,
        isVpnActive = false,
        isDeveloperOptionsEnabled = false,
    )
}

private class FakeLocationProvider : LocationProvider {
    override suspend fun currentLocation() = GpsCoordinate(
        latitude = -12.5458,
        longitude = -55.7061,
        accuracyMeters = 12.0,
    )
}

private class FakeGeofenceValidator(
    private val inside: Boolean,
    private val distance: Double,
) : GeofenceValidator {
    override suspend fun validate(location: GpsCoordinate) = GeofenceValidation(
        isInside = inside,
        distanceToNearestMeters = distance,
    )
}

private class FakeBiometricProcessor : BiometricProcessor {
    override suspend fun captureFrame(): ByteArray = byteArrayOf(1, 2, 3)
    override suspend fun livenessScore(frame: ByteArray): Float = 0.95f
}

private class FakePunchRepository(
    private val failSubmit: Boolean = false,
) : PunchRepository {
    var offlineCount = 0

    override suspend fun submit(request: PunchRequest): PunchResult {
        if (failSubmit) error("network down")
        return PunchResult(
            id = "punch-1",
            status = PunchStatus.VALID,
            punchedAt = "2026-06-26T08:01:02Z",
            type = request.type,
        )
    }

    override suspend fun queueOffline(request: PunchRequest) {
        offlineCount++
    }

    override suspend fun isOnline(): Boolean = true
}
