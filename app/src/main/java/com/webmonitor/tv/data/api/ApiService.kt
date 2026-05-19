package com.webmonitor.tv.data.api

import com.webmonitor.tv.data.models.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {

    // ─── Auth (no JWT required) ──────────────────────────────────────────────

    @POST("login")
    suspend fun login(@Body request: LoginRequest): Response<LoginResponse>

    @POST("registration")
    suspend fun register(@Body request: RegistrationRequest): Response<GenericResponse>

    // ─── Monitors (JWT required → passed via header) ─────────────────────────

    @GET("api/monitors")
    suspend fun getMonitors(
        @Header("Authorization") token: String
    ): Response<MonitorsResponse>

    @POST("api/addmonitor")
    suspend fun addMonitor(
        @Header("Authorization") token: String,
        @Body request: AddMonitorRequest
    ): Response<GenericResponse>

    @POST("api/deletemonitor")
    suspend fun deleteMonitor(
        @Header("Authorization") token: String,
        @Body request: DeleteMonitorRequest
    ): Response<GenericResponse>

    // ─── Checks ──────────────────────────────────────────────────────────────

    @GET("api/checks/{monitor_id}")
    suspend fun getChecks(
        @Header("Authorization") token: String,
        @Path("monitor_id") monitorId: String
    ): Response<ChecksResponse>

    @GET("api/checks/{monitor_id}/regions")
    suspend fun getRegions(
        @Header("Authorization") token: String,
        @Path("monitor_id") monitorId: String
    ): Response<RegionsResponse>
}
