package com.webmonitor.tv

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.viewModels
import androidx.compose.animation.*
import androidx.compose.runtime.*
import androidx.compose.ui.graphics.Color
import androidx.core.view.WindowCompat
import com.webmonitor.tv.data.models.Monitor
import com.webmonitor.tv.ui.screens.ChecksScreen
import com.webmonitor.tv.ui.screens.LoginScreen
import com.webmonitor.tv.ui.screens.MonitorsScreen
import com.webmonitor.tv.viewmodel.MainViewModel

sealed class Screen {
    object Login : Screen()
    object Monitors : Screen()
    data class Checks(val monitor: Monitor) : Screen()
}

class MainActivity : ComponentActivity() {

    private val viewModel: MainViewModel by viewModels()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        // Full-screen immersive for TV
        WindowCompat.setDecorFitsSystemWindows(window, false)

        setContent {
            // Collect login state
            val isLoggedIn by viewModel.isLoggedIn.collectAsState()
            var currentScreen by remember { mutableStateOf<Screen>(Screen.Login) }

            // Auto-navigate based on auth state
            LaunchedEffect(isLoggedIn) {
                if (isLoggedIn && currentScreen is Screen.Login) {
                    currentScreen = Screen.Monitors
                } else if (!isLoggedIn) {
                    currentScreen = Screen.Login
                }
            }

            AnimatedContent(
                targetState = currentScreen,
                transitionSpec = {
                    fadeIn(animationSpec = androidx.compose.animation.core.tween(300)) togetherWith
                    fadeOut(animationSpec = androidx.compose.animation.core.tween(300))
                },
                label = "navigation"
            ) { screen ->
                when (screen) {
                    is Screen.Login -> LoginScreen(viewModel = viewModel)
                    is Screen.Monitors -> MonitorsScreen(
                        viewModel = viewModel,
                        onOpenChecks = { monitor -> currentScreen = Screen.Checks(monitor) }
                    )
                    is Screen.Checks -> ChecksScreen(
                        viewModel = viewModel,
                        monitor = screen.monitor,
                        onBack = { currentScreen = Screen.Monitors }
                    )
                }
            }
        }
    }
}
