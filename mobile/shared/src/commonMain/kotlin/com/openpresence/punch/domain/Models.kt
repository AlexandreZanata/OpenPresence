package com.openpresence.punch.domain

data class GpsCoordinate(
    val latitude: Double,
    val longitude: Double,
    val accuracyMeters: Double,
    val isMocked: Boolean = false,
)

data class DeviceFlag(
    val type: String,
    val severity: String,
)

data class DeviceIntegrityReport(
    val isRooted: Boolean,
    val isVpnActive: Boolean,
    val isDeveloperOptionsEnabled: Boolean,
    val flags: List<DeviceFlag> = emptyList(),
)

data class PunchResult(
    val id: String,
    val status: PunchStatus,
    val punchedAt: String,
    val type: PunchType,
)

data class PunchRequest(
    val type: PunchType,
    val location: GpsCoordinate,
    val frameBytes: ByteArray,
    val deviceReport: DeviceIntegrityReport,
    val deviceTimeIso: String,
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other == null || this::class != other::class) return false
        other as PunchRequest
        return type == other.type &&
            location == other.location &&
            frameBytes.contentEquals(other.frameBytes) &&
            deviceReport == other.deviceReport &&
            deviceTimeIso == other.deviceTimeIso
    }

    override fun hashCode(): Int {
        var result = type.hashCode()
        result = 31 * result + location.hashCode()
        result = 31 * result + frameBytes.contentHashCode()
        result = 31 * result + deviceReport.hashCode()
        result = 31 * result + deviceTimeIso.hashCode()
        return result
    }
}

/** Result of flushing the offline punch queue to the API. */
data class OfflineSyncResult(
    val synced: Int,
    val failed: Int,
    val pending: Int,
)
