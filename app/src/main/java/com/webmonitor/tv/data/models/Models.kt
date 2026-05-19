package com.webmonitor.tv.data.models

import com.google.gson.annotations.SerializedName

// ─── Auth ───────────────────────────────────────────────────────────────────

data class LoginRequest(
    @SerializedName("Login") val login: String,
    @SerializedName("Password") val password: String
)

data class RegistrationRequest(
    @SerializedName("Username") val username: String,
    @SerializedName("Login") val login: String,
    @SerializedName("Password") val password: String,
    @SerializedName("Telegram_id") val telegramId: Long = 0
)

data class LoginResponse(
    @SerializedName("message") val message: String,
    @SerializedName("token") val token: String?
)

data class GenericResponse(
    @SerializedName("message") val message: String
)

// ─── Monitor ─────────────────────────────────────────────────────────────────

data class Monitor(
    @SerializedName("Id") val id: String = "",
    @SerializedName("Users_id") val usersId: String = "",
    @SerializedName("Url") val url: String,
    @SerializedName("Time_interval") val timeInterval: Int,
    @SerializedName("Is_active") val isActive: Boolean = true,
    @SerializedName("Created_at") val createdAt: String = ""
)

data class MonitorsResponse(
    @SerializedName("message") val message: String,
    @SerializedName("monitors") val monitors: List<Monitor>?
)

data class AddMonitorRequest(
    @SerializedName("Url") val url: String,
    @SerializedName("Time_interval") val timeInterval: Int,
    @SerializedName("Is_active") val isActive: Boolean = true
)

data class DeleteMonitorRequest(
    @SerializedName("Id") val id: String
)

// ─── Checks ──────────────────────────────────────────────────────────────────

data class CheckResult(
    @SerializedName("Id") val id: String = "",
    @SerializedName("Monitor_id") val monitorId: String = "",
    @SerializedName("Time_Interval") val timeInterval: Int = 0,
    @SerializedName("Checked_at") val checkedAt: String = "",
    @SerializedName("Responce_time") val responseTime: Double = 0.0,   // typo in backend kept intentionally
    @SerializedName("Status_ok") val statusOk: Boolean = false
)

data class ChecksResponse(
    @SerializedName("message") val message: String,
    @SerializedName("monitors") val checks: List<CheckResult>?          // backend returns key "monitors" for checks too
)

// ─── Regions ──────────────────────────────────────────────────────────────────

data class RegionCheck(
    @SerializedName("Id") val id: String = "",
    @SerializedName("Monitor_id") val monitorId: String = "",
    @SerializedName("Responce_time") val responseTime: Double = 0.0,   // typo in backend kept intentionally
    @SerializedName("Status_ok") val statusOk: Boolean = false,
    @SerializedName("Checked_at") val checkedAt: String = "",
    @SerializedName("Region") val region: String = ""
)

data class RegionsResponse(
    @SerializedName("message") val message: String,
    @SerializedName("regions") val regions: Map<String, List<RegionCheck>>
)
