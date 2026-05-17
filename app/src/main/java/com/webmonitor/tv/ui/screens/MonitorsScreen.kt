package com.webmonitor.tv.ui.screens

import androidx.compose.animation.*
import androidx.compose.foundation.*
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.*
import androidx.compose.foundation.lazy.grid.*
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
import com.webmonitor.tv.data.models.Monitor
import com.webmonitor.tv.ui.components.*
import com.webmonitor.tv.ui.theme.*
import com.webmonitor.tv.viewmodel.MainViewModel

private val INTERVAL_OPTIONS = listOf(
    "30 сек" to 30,
    "1 мин" to 60,
    "5 мин" to 300,
    "15 мин" to 900,
    "1 час" to 3600
)

@Composable
fun MonitorsScreen(
    viewModel: MainViewModel,
    onOpenChecks: (Monitor) -> Unit
) {
    val state by viewModel.monitorsState.collectAsState()

    var showAddDialog by remember { mutableStateOf(false) }
    var deleteTarget by remember { mutableStateOf<Monitor?>(null) }

    LaunchedEffect(Unit) { viewModel.loadMonitors() }

    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(
                Brush.verticalGradient(listOf(NavyDeep, NavyDark))
            )
    ) {
        Column(modifier = Modifier.fillMaxSize()) {

            // ── Top bar ───────────────────────────────────────────────────────
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .background(NavyCard)
                    .padding(horizontal = 32.dp, vertical = 16.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Icon(Icons.Default.MonitorHeart, contentDescription = null, tint = AccentBlue, modifier = Modifier.size(28.dp))
                Spacer(Modifier.width(12.dp))
                Text("WebMonitor TV", fontSize = 22.sp, fontWeight = FontWeight.ExtraBold, color = TextPrimary)

                Spacer(Modifier.weight(1f))

                // Summary stats
                val online = state.monitors.count { it.isActive }
                val total = state.monitors.size

                StatChip("Всего", total.toString(), AccentBlue)
                Spacer(Modifier.width(8.dp))
                StatChip("Активных", online.toString(), AccentGreen)
                Spacer(Modifier.width(24.dp))

                TvButton(
                    text = "Добавить",
                    onClick = { showAddDialog = true },
                    modifier = Modifier.height(44.dp),
                    color = AccentBlue
                )
                Spacer(Modifier.width(12.dp))
                IconButton(
                    onClick = { viewModel.loadMonitors() },
                    modifier = Modifier
                        .clip(RoundedCornerShape(10.dp))
                        .background(NavyMid)
                        .size(44.dp)
                ) {
                    Icon(Icons.Default.Refresh, contentDescription = "Обновить", tint = AccentCyan)
                }
                Spacer(Modifier.width(12.dp))
                IconButton(
                    onClick = { viewModel.logout() },
                    modifier = Modifier
                        .clip(RoundedCornerShape(10.dp))
                        .background(AccentRed.copy(alpha = 0.15f))
                        .size(44.dp)
                ) {
                    Icon(Icons.Default.Logout, contentDescription = "Выход", tint = AccentRed)
                }
            }

            // ── Error / success bar ────────────────────────────────────────────
            AnimatedVisibility(state.error != null) {
                Box(Modifier.fillMaxWidth().padding(horizontal = 32.dp, vertical = 8.dp)) {
                    ErrorBar(state.error ?: "") { viewModel.clearMonitorMessage() }
                }
            }
            AnimatedVisibility(state.successMessage != null) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 32.dp, vertical = 8.dp)
                        .clip(RoundedCornerShape(12.dp))
                        .background(AccentGreen.copy(alpha = 0.15f))
                        .padding(12.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(Icons.Default.CheckCircle, contentDescription = null, tint = AccentGreen, modifier = Modifier.size(18.dp))
                    Spacer(Modifier.width(8.dp))
                    Text(state.successMessage ?: "", color = AccentGreen, fontSize = 13.sp)
                }
            }

            // ── Content ───────────────────────────────────────────────────────
            if (state.isLoading && state.monitors.isEmpty()) {
                Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        CircularProgressIndicator(color = AccentBlue, strokeWidth = 3.dp)
                        Spacer(Modifier.height(16.dp))
                        Text("Загрузка мониторов...", color = TextSecondary)
                    }
                }
            } else if (state.monitors.isEmpty()) {
                Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(Icons.Default.CloudOff, contentDescription = null, tint = TextDim, modifier = Modifier.size(64.dp))
                        Spacer(Modifier.height(16.dp))
                        Text("Мониторов пока нет", fontSize = 20.sp, color = TextSecondary)
                        Spacer(Modifier.height(8.dp))
                        Text("Нажмите «Добавить» чтобы начать", color = TextDim)
                    }
                }
            } else {
                LazyVerticalGrid(
                    columns = GridCells.Adaptive(minSize = 340.dp),
                    contentPadding = PaddingValues(32.dp),
                    verticalArrangement = Arrangement.spacedBy(16.dp),
                    horizontalArrangement = Arrangement.spacedBy(16.dp),
                    modifier = Modifier.fillMaxSize()
                ) {
                    items(state.monitors.size) { idx ->
                        MonitorCard(
                            monitor = state.monitors[idx],
                            onViewChecks = { onOpenChecks(state.monitors[idx]) },
                            onDelete = { deleteTarget = state.monitors[idx] }
                        )
                    }
                }
            }
        }

        // ── Add Monitor Dialog ─────────────────────────────────────────────────
        if (showAddDialog) {
            AddMonitorDialog(
                onDismiss = { showAddDialog = false },
                onConfirm = { url, interval ->
                    viewModel.addMonitor(url, interval)
                    showAddDialog = false
                }
            )
        }

        // ── Delete Confirm Dialog ──────────────────────────────────────────────
        deleteTarget?.let { mon ->
            AlertDialog(
                onDismissRequest = { deleteTarget = null },
                title = { Text("Удалить монитор?", color = TextPrimary) },
                text = { Text("${mon.url}\nЭто действие нельзя отменить.", color = TextSecondary) },
                confirmButton = {
                    TvButton("Удалить", onClick = {
                        viewModel.deleteMonitor(mon.id)
                        deleteTarget = null
                    }, color = AccentRed)
                },
                dismissButton = {
                    TextButton(onClick = { deleteTarget = null }) {
                        Text("Отмена", color = TextSecondary)
                    }
                },
                containerColor = NavyCard,
                shape = RoundedCornerShape(16.dp)
            )
        }
    }
}

// ─── Monitor Card ─────────────────────────────────────────────────────────────

@Composable
fun MonitorCard(
    monitor: Monitor,
    onViewChecks: () -> Unit,
    onDelete: () -> Unit
) {
    GlassCard(
        glowColor = if (monitor.isActive) GlowGreen else GlowRed,
        modifier = Modifier.fillMaxWidth()
    ) {
        Row(verticalAlignment = Alignment.Top) {
            Column(modifier = Modifier.weight(1f)) {
                Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                    StatusBadge(isOnline = monitor.isActive)
                }
                Spacer(Modifier.height(8.dp))
                Text(
                    monitor.url,
                    color = TextPrimary,
                    fontSize = 15.sp,
                    fontWeight = FontWeight.SemiBold,
                    maxLines = 2,
                    overflow = androidx.compose.ui.text.style.TextOverflow.Ellipsis
                )
                Spacer(Modifier.height(4.dp))
                Row(
                    horizontalArrangement = Arrangement.spacedBy(6.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(Icons.Default.Timer, contentDescription = null, tint = TextDim, modifier = Modifier.size(13.dp))
                    Text(
                        "Интервал: ${formatInterval(monitor.timeInterval)}",
                        color = TextSecondary,
                        fontSize = 12.sp
                    )
                }
                if (monitor.createdAt.isNotBlank()) {
                    Spacer(Modifier.height(2.dp))
                    Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                        Icon(Icons.Default.CalendarToday, contentDescription = null, tint = TextDim, modifier = Modifier.size(13.dp))
                        Text(
                            monitor.createdAt.take(10),
                            color = TextDim,
                            fontSize = 11.sp
                        )
                    }
                }
            }

            Column(horizontalAlignment = Alignment.End, verticalArrangement = Arrangement.spacedBy(6.dp)) {
                IconButton(
                    onClick = onViewChecks,
                    modifier = Modifier
                        .clip(RoundedCornerShape(8.dp))
                        .background(AccentBlue.copy(alpha = 0.15f))
                        .size(36.dp)
                ) {
                    Icon(Icons.Default.Analytics, contentDescription = "История", tint = AccentBlue, modifier = Modifier.size(18.dp))
                }
                IconButton(
                    onClick = onDelete,
                    modifier = Modifier
                        .clip(RoundedCornerShape(8.dp))
                        .background(AccentRed.copy(alpha = 0.12f))
                        .size(36.dp)
                ) {
                    Icon(Icons.Default.Delete, contentDescription = "Удалить", tint = AccentRed, modifier = Modifier.size(18.dp))
                }
            }
        }
    }
}

// ─── Add Dialog ───────────────────────────────────────────────────────────────

@Composable
fun AddMonitorDialog(
    onDismiss: () -> Unit,
    onConfirm: (url: String, intervalSeconds: Int) -> Unit
) {
    var url by remember { mutableStateOf("https://") }
    var selectedInterval by remember { mutableStateOf(INTERVAL_OPTIONS[1]) }
    var expanded by remember { mutableStateOf(false) }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = {
            Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                Icon(Icons.Default.AddCircle, contentDescription = null, tint = AccentBlue, modifier = Modifier.size(22.dp))
                Text("Новый монитор", color = TextPrimary, fontWeight = FontWeight.Bold)
            }
        },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(16.dp)) {
                TvTextField(
                    value = url,
                    onValueChange = { url = it },
                    label = "URL сайта",
                    leadingIcon = Icons.Default.Link,
                    keyboardType = KeyboardType.Uri
                )

                // Interval dropdown
                Box {
                    OutlinedButton(
                        onClick = { expanded = true },
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(12.dp),
                        colors = ButtonDefaults.outlinedButtonColors(contentColor = TextPrimary),
                        border = BorderStroke(1.dp, NavyBorder)
                    ) {
                        Icon(Icons.Default.Timer, contentDescription = null, tint = AccentBlue, modifier = Modifier.size(18.dp))
                        Spacer(Modifier.width(8.dp))
                        Text("Интервал: ${selectedInterval.first}", modifier = Modifier.weight(1f))
                        Icon(Icons.Default.ArrowDropDown, contentDescription = null, tint = TextSecondary)
                    }
                    DropdownMenu(
                        expanded = expanded,
                        onDismissRequest = { expanded = false },
                        modifier = Modifier.background(NavyCard)
                    ) {
                        INTERVAL_OPTIONS.forEach { option ->
                            DropdownMenuItem(
                                text = { Text(option.first, color = if (option == selectedInterval) AccentBlue else TextPrimary) },
                                onClick = { selectedInterval = option; expanded = false },
                                leadingIcon = {
                                    if (option == selectedInterval)
                                        Icon(Icons.Default.Check, contentDescription = null, tint = AccentBlue, modifier = Modifier.size(16.dp))
                                }
                            )
                        }
                    }
                }
            }
        },
        confirmButton = {
            TvButton(
                "Добавить",
                onClick = {
                    val trimmed = url.trim()
                    if (trimmed.isNotBlank()) onConfirm(trimmed, selectedInterval.second)
                },
                color = AccentBlue,
                enabled = url.trim().length > 8
            )
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("Отмена", color = TextSecondary) }
        },
        containerColor = NavyCard,
        shape = RoundedCornerShape(20.dp)
    )
}

private fun formatInterval(seconds: Int): String = when {
    seconds < 60 -> "$seconds сек"
    seconds < 3600 -> "${seconds / 60} мин"
    else -> "${seconds / 3600} ч"
}
