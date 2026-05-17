package com.webmonitor.tv.data.repository

import com.webmonitor.tv.data.api.ApiService
import com.webmonitor.tv.data.models.*

sealed class Result<out T> {
    data class Success<T>(val data: T) : Result<T>()
    data class Error(val message: String) : Result<Nothing>()
}

class MonitorRepository(private val api: ApiService) {

    private fun bearerToken(token: String) = "Bearer $token"

    // ─── Auth ────────────────────────────────────────────────────────────────

    suspend fun login(login: String, password: String): Result<LoginResponse> {
        return try {
            val response = api.login(LoginRequest(login, password))
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.token != null) Result.Success(body)
                else Result.Error(body?.message ?: "Неверный логин или пароль")
            } else {
                Result.Error("Неверный логин или пароль")
            }
        } catch (e: Exception) {
            Result.Error("Ошибка подключения: ${e.localizedMessage}")
        }
    }

    suspend fun register(username: String, login: String, password: String): Result<GenericResponse> {
        return try {
            val response = api.register(RegistrationRequest(username, login, password))
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                val msg = if (response.code() == 409) "Пользователь уже существует" else "Ошибка регистрации"
                Result.Error(msg)
            }
        } catch (e: Exception) {
            Result.Error("Ошибка подключения: ${e.localizedMessage}")
        }
    }

    // ─── Monitors ────────────────────────────────────────────────────────────

    suspend fun getMonitors(token: String): Result<List<Monitor>> {
        return try {
            val response = api.getMonitors(bearerToken(token))
            if (response.isSuccessful) {
                Result.Success(response.body()?.monitors ?: emptyList())
            } else {
                Result.Error("Ошибка загрузки мониторов: ${response.code()}")
            }
        } catch (e: Exception) {
            Result.Error("Ошибка подключения: ${e.localizedMessage}")
        }
    }

    suspend fun addMonitor(token: String, url: String, intervalSeconds: Int): Result<GenericResponse> {
        return try {
            val response = api.addMonitor(
                bearerToken(token),
                AddMonitorRequest(url = url, timeInterval = intervalSeconds)
            )
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Не удалось добавить монитор")
            }
        } catch (e: Exception) {
            Result.Error("Ошибка подключения: ${e.localizedMessage}")
        }
    }

    suspend fun deleteMonitor(token: String, monitorId: String): Result<GenericResponse> {
        return try {
            val response = api.deleteMonitor(bearerToken(token), DeleteMonitorRequest(monitorId))
            if (response.isSuccessful) {
                Result.Success(response.body()!!)
            } else {
                Result.Error("Не удалось удалить монитор")
            }
        } catch (e: Exception) {
            Result.Error("Ошибка подключения: ${e.localizedMessage}")
        }
    }

    // ─── Checks ──────────────────────────────────────────────────────────────

    suspend fun getChecks(token: String, monitorId: String): Result<List<CheckResult>> {
        return try {
            val response = api.getChecks(bearerToken(token), monitorId)
            if (response.isSuccessful) {
                Result.Success(response.body()?.checks ?: emptyList())
            } else {
                Result.Error("Не удалось загрузить историю проверок")
            }
        } catch (e: Exception) {
            Result.Error("Ошибка подключения: ${e.localizedMessage}")
        }
    }
}
