package com.openpresence.punch.di

import com.openpresence.punch.presentation.PunchViewModel
import com.openpresence.punch.ports.BiometricProcessor
import com.openpresence.punch.ports.DeviceIntegrityChecker
import com.openpresence.punch.ports.GeofenceValidator
import com.openpresence.punch.ports.LocationProvider
import com.openpresence.punch.ports.PunchRepository
import kotlinx.coroutines.CoroutineScope
import org.koin.dsl.module

fun punchModule(scope: CoroutineScope) = module {
    factory {
        PunchViewModel(
            deviceChecker = get<DeviceIntegrityChecker>(),
            locationProvider = get<LocationProvider>(),
            geofenceValidator = get<GeofenceValidator>(),
            biometricProcessor = get<BiometricProcessor>(),
            repository = get<PunchRepository>(),
            scope = scope,
        )
    }
}
