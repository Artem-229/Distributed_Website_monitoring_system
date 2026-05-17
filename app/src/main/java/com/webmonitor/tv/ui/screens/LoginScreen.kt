package com.webmonitor.tv.ui.screens

import androidx.compose.animation.*
import androidx.compose.foundation.*
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.*
import com.webmonitor.tv.ui.components.*
import com.webmonitor.tv.ui.theme.*
import com.webmonitor.tv.viewmodel.AuthState
import com.webmonitor.tv.viewmodel.MainViewModel
import androidx.compose.ui.geometry.Offset

@Composable
fun LoginScreen(viewModel: MainViewModel) {
    val authState by viewModel.authState.collectAsState()

    var isRegisterMode by remember { mutableStateOf(false) }
    var username by remember { mutableStateOf("") }
    var login by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }

    // Reset fields when switching modes
    LaunchedEffect(isRegisterMode) {
        username = ""; login = ""; password = ""
        viewModel.resetAuthState()
    }

    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.radialGradient(
                    colors = listOf(NavyMid, NavyDeep),
                    radius = 1200f
                )
            )
    ) {
        // Decorative grid lines
        Canvas(modifier = Modifier.fillMaxSize()) {
            val step = 80f
            for (x in 0..size.width.toInt() step step.toInt()) {
                drawLine(
                    color = NavyBorder.copy(alpha = 0.15f),
                    start = Offset(x.toFloat(), 0f),
                    end = Offset(x.toFloat(), size.height),
                    strokeWidth = 1f
                )
            }
            for (y in 0..size.height.toInt() step step.toInt()) {
                drawLine(
                    color = NavyBorder.copy(alpha = 0.15f),
                    start = Offset(0f, y.toFloat()),
                    end = Offset(size.width, y.toFloat()),
                    strokeWidth = 1f
                )
            }
        }

        Row(
            modifier = Modifier.fillMaxSize(),
            verticalAlignment = Alignment.CenterVertically
        ) {
            // ── Left panel – branding ─────────────────────────────────────────
            Column(
                modifier = Modifier
                    .weight(1f)
                    .padding(48.dp),
                verticalArrangement = Arrangement.Center,
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Icon(
                    Icons.Default.MonitorHeart,
                    contentDescription = null,
                    tint = AccentBlue,
                    modifier = Modifier.size(80.dp)
                )
                Spacer(Modifier.height(20.dp))
                Text(
                    "WebMonitor",
                    fontSize = 40.sp,
                    fontWeight = FontWeight.ExtraBold,
                    color = TextPrimary,
                    letterSpacing = 2.sp
                )
                Text(
                    "TV",
                    fontSize = 40.sp,
                    fontWeight = FontWeight.ExtraBold,
                    color = AccentCyan,
                    letterSpacing = 2.sp
                )
                Spacer(Modifier.height(16.dp))
                Text(
                    "Мониторинг доступности\nвеб-сайтов в реальном времени",
                    fontSize = 15.sp,
                    color = TextSecondary,
                    textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                    lineHeight = 22.sp
                )

                Spacer(Modifier.height(40.dp))

                // Stats decorative chips
                Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                    StatChip("Серверов", "6+", AccentCyan)
                    StatChip("Регионов", "EU/AS/US", AccentBlue)
                }
            }

            // ── Right panel – form ────────────────────────────────────────────
            Box(
                modifier = Modifier
                    .weight(1f)
                    .fillMaxHeight()
                    .padding(32.dp),
                contentAlignment = Alignment.Center
            ) {
                GlassCard(modifier = Modifier.fillMaxWidth()) {
                    AnimatedContent(
                        targetState = isRegisterMode,
                        transitionSpec = { fadeIn() togetherWith fadeOut() },
                        label = "form"
                    ) { registering ->
                        Column(verticalArrangement = Arrangement.spacedBy(16.dp)) {
                            Text(
                                if (registering) "Регистрация" else "Вход",
                                fontSize = 24.sp,
                                fontWeight = FontWeight.Bold,
                                color = TextPrimary
                            )
                            Text(
                                if (registering) "Создайте аккаунт для начала работы"
                                else "Войдите в систему мониторинга",
                                fontSize = 13.sp,
                                color = TextSecondary
                            )

                            HorizontalDivider(color = NavyBorder)

                            if (registering) {
                                TvTextField(
                                    value = username,
                                    onValueChange = { username = it },
                                    label = "Имя пользователя",
                                    leadingIcon = Icons.Default.Person
                                )
                            }

                            TvTextField(
                                value = login,
                                onValueChange = { login = it },
                                label = "Логин",
                                leadingIcon = Icons.Default.AccountCircle,
                                keyboardType = KeyboardType.Email
                            )

                            TvTextField(
                                value = password,
                                onValueChange = { password = it },
                                label = "Пароль",
                                isPassword = true,
                                leadingIcon = Icons.Default.Lock
                            )

                            // Error
                            AnimatedVisibility(authState is AuthState.Error) {
                                val msg = (authState as? AuthState.Error)?.message ?: ""
                                ErrorBar(msg) { viewModel.resetAuthState() }
                            }

                            Spacer(Modifier.height(4.dp))

                            TvButton(
                                text = if (registering) "Зарегистрироваться" else "Войти",
                                onClick = {
                                    if (registering) viewModel.register(username, login, password)
                                    else viewModel.login(login, password)
                                },
                                modifier = Modifier.fillMaxWidth(),
                                isLoading = authState is AuthState.Loading
                            )

                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.Center,
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Text(
                                    if (registering) "Уже есть аккаунт?" else "Нет аккаунта?",
                                    color = TextSecondary,
                                    fontSize = 13.sp
                                )
                                TextButton(onClick = { isRegisterMode = !isRegisterMode }) {
                                    Text(
                                        if (registering) "Войти" else "Зарегистрироваться",
                                        color = AccentBlue,
                                        fontSize = 13.sp,
                                        fontWeight = FontWeight.SemiBold
                                    )
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
