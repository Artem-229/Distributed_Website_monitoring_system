package com.webmonitor.tv.viewmodel

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.viewModelScope
import com.webmonitor.tv.data.api.RetrofitClient
import com.webmonitor.tv.data.models.CheckResult
import com.webmonitor.tv.data.models.Monitor
import com.webmonitor.tv.data.repository.MonitorRepository
import com.webmonitor.tv.data.repository.Result
import com.webmonitor.tv.data.repository.SessionManager
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch

// ─── UI States ───────────────────────────────────────────────────────────────

sealed class AuthState {
    object Idle : AuthState()
    object Loading : AuthState()
    object Success : AuthState()
    data class Error(val message: String) : AuthState()
}

data class MonitorsUiState(
    val monitors: List<Monitor> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val successMessage: String? = null
)

data class ChecksUiState(
    val checks: List<CheckResult> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null,
    val selectedMonitor: Monitor? = null
)

// ─── ViewModel ───────────────────────────────────────────────────────────────

class MainViewModel(application: Application) : AndroidViewModel(application) {

    private val repository = MonitorRepository(RetrofitClient.instance)
    private val session = SessionManager(application)

    // ── Auth ─────────────────────────────────────────────────────────────────

    val token: StateFlow<String?> = session.tokenFlow
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5000), null)

    val isLoggedIn: StateFlow<Boolean> = token
        .map { it != null && it.isNotBlank() }
        .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5000), false)

    private val _authState = MutableStateFlow<AuthState>(AuthState.Idle)
    val authState: StateFlow<AuthState> = _authState.asStateFlow()

    fun login(login: String, password: String) {
        viewModelScope.launch {
            _authState.value = AuthState.Loading
            when (val result = repository.login(login, password)) {
                is Result.Success -> {
                    val tkn = result.data.token ?: ""
                    session.saveSession(tkn, login)
                    _authState.value = AuthState.Success
                    loadMonitors()
                }
                is Result.Error -> {
                    _authState.value = AuthState.Error(result.message)
                }
            }
        }
    }

    fun register(username: String, login: String, password: String) {
        viewModelScope.launch {
            _authState.value = AuthState.Loading
            when (val result = repository.register(username, login, password)) {
                is Result.Success -> {
                    // auto-login after registration
                    loginAfterRegister(login, password)
                }
                is Result.Error -> {
                    _authState.value = AuthState.Error(result.message)
                }
            }
        }
    }

    private fun loginAfterRegister(login: String, password: String) {
        viewModelScope.launch {
            when (val result = repository.login(login, password)) {
                is Result.Success -> {
                    session.saveSession(result.data.token ?: "", login)
                    _authState.value = AuthState.Success
                    loadMonitors()
                }
                is Result.Error -> _authState.value = AuthState.Error(result.message)
            }
        }
    }

    fun logout() {
        viewModelScope.launch {
            session.clearSession()
            _authState.value = AuthState.Idle
            _monitorsState.value = MonitorsUiState()
            _checksState.value = ChecksUiState()
        }
    }

    fun resetAuthState() {
        _authState.value = AuthState.Idle
    }

    // ── Monitors ─────────────────────────────────────────────────────────────

    private val _monitorsState = MutableStateFlow(MonitorsUiState())
    val monitorsState: StateFlow<MonitorsUiState> = _monitorsState.asStateFlow()

    fun loadMonitors() {
        val tkn = token.value ?: return
        viewModelScope.launch {
            _monitorsState.update { it.copy(isLoading = true, error = null) }
            when (val result = repository.getMonitors(tkn)) {
                is Result.Success -> {
                    _monitorsState.value = MonitorsUiState(monitors = result.data)
                }
                is Result.Error -> {
                    _monitorsState.update { it.copy(isLoading = false, error = result.message) }
                }
            }
        }
    }

    fun addMonitor(url: String, intervalSeconds: Int) {
        val tkn = token.value ?: return
        viewModelScope.launch {
            _monitorsState.update { it.copy(isLoading = true, error = null) }
            when (val result = repository.addMonitor(tkn, url, intervalSeconds)) {
                is Result.Success -> {
                    _monitorsState.update { it.copy(successMessage = "Монитор добавлен") }
                    loadMonitors()
                }
                is Result.Error -> {
                    _monitorsState.update { it.copy(isLoading = false, error = result.message) }
                }
            }
        }
    }

    fun deleteMonitor(monitorId: String) {
        val tkn = token.value ?: return
        viewModelScope.launch {
            when (val result = repository.deleteMonitor(tkn, monitorId)) {
                is Result.Success -> {
                    _monitorsState.update { it.copy(successMessage = "Монитор удалён") }
                    loadMonitors()
                }
                is Result.Error -> {
                    _monitorsState.update { it.copy(error = result.message) }
                }
            }
        }
    }

    fun clearMonitorMessage() {
        _monitorsState.update { it.copy(successMessage = null, error = null) }
    }

    // ── Checks ───────────────────────────────────────────────────────────────

    private val _checksState = MutableStateFlow(ChecksUiState())
    val checksState: StateFlow<ChecksUiState> = _checksState.asStateFlow()

    fun loadChecks(monitor: Monitor) {
        val tkn = token.value ?: return
        viewModelScope.launch {
            _checksState.value = ChecksUiState(isLoading = true, selectedMonitor = monitor)
            when (val result = repository.getChecks(tkn, monitor.id)) {
                is Result.Success -> {
                    _checksState.value = ChecksUiState(
                        checks = result.data,
                        selectedMonitor = monitor
                    )
                }
                is Result.Error -> {
                    _checksState.value = ChecksUiState(
                        error = result.message,
                        selectedMonitor = monitor
                    )
                }
            }
        }
    }

    fun clearChecks() {
        _checksState.value = ChecksUiState()
    }
}
