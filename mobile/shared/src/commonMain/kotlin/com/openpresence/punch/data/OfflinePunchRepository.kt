package com.openpresence.punch.data

import com.openpresence.punch.domain.OfflineSyncResult
import com.openpresence.punch.domain.PunchRequest
import com.openpresence.punch.domain.PunchResult
import com.openpresence.punch.ports.PunchApi
import com.openpresence.punch.ports.PunchRepository

/** Queues punches offline and syncs them through [PunchApi] when online. */
class OfflinePunchRepository(
    private val api: PunchApi,
    private val connectivity: () -> Boolean,
) : PunchRepository {
    private val queue = mutableListOf<PunchRequest>()

    override suspend fun isOnline(): Boolean = connectivity()

    override suspend fun submit(request: PunchRequest): PunchResult = api.submit(request)

    override suspend fun queueOffline(request: PunchRequest) {
        queue.add(request)
    }

    override suspend fun syncOfflineQueue(): OfflineSyncResult {
        if (!isOnline()) {
            return OfflineSyncResult(synced = 0, failed = 0, pending = queue.size)
        }
        var synced = 0
        var failed = 0
        for (request in queue.toList()) {
            try {
                api.submit(request)
                queue.remove(request)
                synced++
            } catch (_: Exception) {
                failed++
                break
            }
        }
        return OfflineSyncResult(synced = synced, failed = failed, pending = queue.size)
    }

    fun pendingCount(): Int = queue.size
}
