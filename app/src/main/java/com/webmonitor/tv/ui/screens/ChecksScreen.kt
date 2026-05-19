package com.webmonitor.tv.ui.screens

import androidx.compose.animation.*
import androidx.compose.foundation.*
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.*
import com.webmonitor.tv.data.models.CheckResult
import com.webmonitor.tv.data.models.Monitor
import com.webmonitor.tv.data.models.RegionCheck
import com.webmonitor.tv.ui.components.*
import com.webmonitor.tv.ui.theme.*
import com.webmonitor.tv.viewmodel.MainViewModel
import kotlinx.coroutines.delay
import kotlin.math.roundToInt

private val REGION_COLOR = mapOf(
    "Russia"         to AccentCyan,
    "USA"            to AccentBlue,
    "China"          to AccentRed,
    "Central Europe" to AccentGreen,
    "Asia Pacific"   to AccentAmber
)

@Composable
fun ChecksScreen(
    viewModel: MainViewModel,
    monitor: Monitor,
    onBack: () -> Unit
) {
    val checksState by viewModel.checksState.collectAsState()
    val regionsState by viewModel.regionsState.collectAsState()
    var selectedTab by remember { mutableStateOf(0) }

    LaunchedEffect(monitor.id) { viewModel.loadChecks(monitor) }

    // Auto-poll regions every 30s while that tab is open
    LaunchedEffect(selectedTab, monitor.id) {
        if (selectedTab == 1) {
            while (true) {
                viewModel.loadRegions(monitor.id)
                delay(30_000L)
            }
        }
    }

    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Brush.verticalGradient(listOf(NavyDeep, NavyDark)))
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
                IconButton(
                    onClick = { viewModel.clearChecks(); viewModel.clearRegions(); onBack() },
                    modifier = Modifier
                        .clip(RoundedCornerShape(10.dp))
                        .background(NavyMid)
                        .size(44.dp)
                ) {
                    Icon(Icons.Default.ArrowBack, contentDescription = "Назад", tint = TextPrimary)
                }
                Spacer(Modifier.width(16.dp))
                Column {
                    Text("История проверок", fontSize = 20.sp, fontWeight = FontWeight.Bold, color = TextPrimary)
                    Text(
                        monitor.url,
                        fontSize = 12.sp,
                        color = TextSecondary,
                        maxLines = 1,
                        overflow = androidx.compose.ui.text.style.TextOverflow.Ellipsis
                    )
                }
                Spacer(Modifier.weight(1f))
                StatusBadge(isOnline = monitor.isActive)
                Spacer(Modifier.width(16.dp))
                IconButton(
                    onClick = {
                        if (selectedTab == 0) viewModel.loadChecks(monitor)
                        else viewModel.loadRegions(monitor.id)
                    },
                    modifier = Modifier
                        .clip(RoundedCornerShape(10.dp))
                        .background(NavyMid)
                        .size(44.dp)
                ) {
                    Icon(Icons.Default.Refresh, contentDescription = "Обновить", tint = AccentCyan)
                }
            }

            // ── Tab selector ──────────────────────────────────────────────────
            Row(
                modifier = Modifier.padding(horizontal = 32.dp, vertical = 12.dp),
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                TabButton(text = "История", selected = selectedTab == 0) { selectedTab = 0 }
                TabButton(text = "Регионы", selected = selectedTab == 1) { selectedTab = 1 }
            }

            // ── Content ───────────────────────────────────────────────────────
            AnimatedContent(
                targetState = selectedTab,
                modifier = Modifier.weight(1f).fillMaxWidth(),
                transitionSpec = { fadeIn() togetherWith fadeOut() },
                label = "tab"
            ) { tab ->
                when (tab) {
                    0 -> ChecksHistoryContent(
                        isLoading = checksState.isLoading,
                        error = checksState.error,
                        checks = checksState.checks
                    )
                    else -> RegionsContent(
                        isLoading = regionsState.isLoading,
                        error = regionsState.error,
                        regions = regionsState.regions
                    )
                }
            }
        }
    }
}

@Composable
private fun TabButton(text: String, selected: Boolean, onClick: () -> Unit) {
    Button(
        onClick = onClick,
        shape = RoundedCornerShape(10.dp),
        colors = ButtonDefaults.buttonColors(
            containerColor = if (selected) AccentBlue else NavyMid,
            contentColor = if (selected) TextPrimary else TextSecondary
        ),
        modifier = Modifier.height(36.dp),
        contentPadding = PaddingValues(horizontal = 20.dp, vertical = 0.dp)
    ) {
        Text(text, fontSize = 14.sp, fontWeight = if (selected) FontWeight.Bold else FontWeight.Normal)
    }
}

// ─── История ─────────────────────────────────────────────────────────────────

@Composable
private fun ChecksHistoryContent(
    isLoading: Boolean,
    error: String?,
    checks: List<CheckResult>
) {
    if (isLoading) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                CircularProgressIndicator(color = AccentBlue, strokeWidth = 3.dp)
                Spacer(Modifier.height(16.dp))
                Text("Загрузка истории...", color = TextSecondary)
            }
        }
        return
    }

    if (error != null) {
        Box(Modifier.fillMaxSize().padding(32.dp), contentAlignment = Alignment.Center) {
            ErrorBar(error) {}
        }
        return
    }

    if (checks.isEmpty()) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(Icons.Default.HourglassEmpty, contentDescription = null, tint = TextDim, modifier = Modifier.size(64.dp))
                Spacer(Modifier.height(16.dp))
                Text("Проверок пока нет", fontSize = 20.sp, color = TextSecondary)
                Text("Данные появятся после первой проверки бэкендом", color = TextDim)
            }
        }
        return
    }

    val avgTime = checks.map { it.responseTime }.average()
    val maxTime = checks.maxOf { it.responseTime }
    val minTime = checks.minOf { it.responseTime }
    val okCount = checks.count { it.statusOk }
    val uptimePercent = (okCount.toDouble() / checks.size * 100).roundToInt()

    Column(modifier = Modifier.fillMaxSize()) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 32.dp, vertical = 16.dp),
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            StatChip("Uptime", "$uptimePercent%", if (uptimePercent >= 95) AccentGreen else AccentAmber)
            StatChip("Avg", "${avgTime.roundToInt()} ms", AccentCyan)
            StatChip("Min", "${minTime.roundToInt()} ms", AccentGreen)
            StatChip("Max", "${maxTime.roundToInt()} ms", AccentAmber)
            StatChip("Записей", "${checks.size}", AccentBlue)
        }

        if (checks.size >= 2) {
            ResponseTimeChart(
                checks = checks,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(100.dp)
                    .padding(horizontal = 32.dp)
            )
            Spacer(Modifier.height(12.dp))
        }

        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 32.dp)
                .clip(RoundedCornerShape(topStart = 12.dp, topEnd = 12.dp))
                .background(NavyCard)
                .padding(horizontal = 16.dp, vertical = 10.dp)
        ) {
            Text("Время", color = TextSecondary, fontSize = 12.sp, modifier = Modifier.weight(2f))
            Text("Статус", color = TextSecondary, fontSize = 12.sp, modifier = Modifier.weight(1f))
            Text("Отклик", color = TextSecondary, fontSize = 12.sp, modifier = Modifier.weight(1f))
            Text("Интервал", color = TextSecondary, fontSize = 12.sp, modifier = Modifier.weight(1f))
        }

        LazyColumn(
            modifier = Modifier
                .weight(1f)
                .fillMaxWidth()
                .padding(horizontal = 32.dp)
                .clip(RoundedCornerShape(bottomStart = 12.dp, bottomEnd = 12.dp))
        ) {
            items(checks.size) { idx ->
                CheckRow(check = checks[idx], isEven = idx % 2 == 0)
            }
        }
    }
}

// ─── Регионы ─────────────────────────────────────────────────────────────────

@Composable
private fun RegionsContent(
    isLoading: Boolean,
    error: String?,
    regions: Map<String, List<RegionCheck>>
) {
    if (isLoading) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                CircularProgressIndicator(color = AccentBlue, strokeWidth = 3.dp)
                Spacer(Modifier.height(16.dp))
                Text("Загрузка региональных данных...", color = TextSecondary)
            }
        }
        return
    }

    if (error != null) {
        Box(Modifier.fillMaxSize().padding(32.dp), contentAlignment = Alignment.Center) {
            ErrorBar(error) {}
        }
        return
    }

    if (regions.isEmpty()) {
        Box(Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
            Column(horizontalAlignment = Alignment.CenterHorizontally) {
                Icon(Icons.Default.CloudOff, contentDescription = null, tint = TextDim, modifier = Modifier.size(64.dp))
                Spacer(Modifier.height(16.dp))
                Text("Региональных данных пока нет", fontSize = 20.sp, color = TextSecondary)
            }
        }
        return
    }

    Row(
        modifier = Modifier
            .fillMaxSize()
            .padding(horizontal = 32.dp, vertical = 16.dp),
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        regions.forEach { (name, checks) ->
            RegionCard(
                modifier = Modifier.weight(1f),
                regionName = name,
                checks = checks,
                accentColor = REGION_COLOR[name] ?: AccentBlue
            )
        }
    }
}

@Composable
private fun RegionCard(
    modifier: Modifier = Modifier,
    regionName: String,
    checks: List<RegionCheck>,
    accentColor: Color
) {
    val latest = checks.firstOrNull()
    val chronological = checks.reversed()

    GlassCard(
        modifier = modifier,
        glowColor = if (latest?.statusOk == true) GlowGreen else GlowRed
    ) {
        Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
            Row(
                verticalAlignment = Alignment.CenterVertically,
                horizontalArrangement = Arrangement.spacedBy(8.dp)
            ) {
                Icon(
                    Icons.Default.Language,
                    contentDescription = null,
                    tint = accentColor,
                    modifier = Modifier.size(18.dp)
                )
                Text(
                    regionName,
                    color = TextPrimary,
                    fontSize = 14.sp,
                    fontWeight = FontWeight.Bold,
                    modifier = Modifier.weight(1f)
                )
            }

            if (latest != null) {
                StatusBadge(isOnline = latest.statusOk)
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    Icon(Icons.Default.Timer, contentDescription = null, tint = accentColor, modifier = Modifier.size(13.dp))
                    Text(
                        "${latest.responseTime.roundToInt()} ms",
                        color = accentColor,
                        fontSize = 13.sp,
                        fontWeight = FontWeight.SemiBold
                    )
                }
                Text(
                    latest.checkedAt.take(19).replace("T", " "),
                    color = TextDim,
                    fontSize = 11.sp
                )
            } else {
                Text("Нет данных", color = TextDim, fontSize = 12.sp)
            }

            if (chronological.size >= 2) {
                RegionMiniChart(
                    checks = chronological,
                    accentColor = accentColor,
                    modifier = Modifier.fillMaxWidth().height(60.dp)
                )
            }
        }
    }
}

@Composable
private fun RegionMiniChart(
    checks: List<RegionCheck>,
    accentColor: Color,
    modifier: Modifier = Modifier
) {
    val maxVal = checks.maxOf { it.responseTime }.coerceAtLeast(1.0)

    Canvas(modifier = modifier) {
        val barWidth = (size.width / checks.size - 3f).coerceAtLeast(2f)
        checks.forEachIndexed { idx, check ->
            val height = (check.responseTime / maxVal * size.height).toFloat().coerceAtLeast(3f)
            val x = idx * (size.width / checks.size)
            drawRoundRect(
                color = accentColor.copy(alpha = if (check.statusOk) 0.75f else 0.3f),
                topLeft = androidx.compose.ui.geometry.Offset(x + 1.5f, size.height - height),
                size = androidx.compose.ui.geometry.Size(barWidth, height),
                cornerRadius = androidx.compose.ui.geometry.CornerRadius(2f)
            )
        }
    }
}

// ─── Shared subcomponents ────────────────────────────────────────────────────

@Composable
private fun CheckRow(check: CheckResult, isEven: Boolean) {
    val bg = if (isEven) NavyDark else NavyMid
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .background(bg)
            .padding(horizontal = 16.dp, vertical = 10.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            check.checkedAt.take(19).replace("T", " "),
            color = TextSecondary,
            fontSize = 12.sp,
            modifier = Modifier.weight(2f)
        )
        Row(modifier = Modifier.weight(1f), verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(4.dp)) {
            Icon(
                if (check.statusOk) Icons.Default.CheckCircle else Icons.Default.Cancel,
                contentDescription = null,
                tint = if (check.statusOk) AccentGreen else AccentRed,
                modifier = Modifier.size(14.dp)
            )
            Text(
                if (check.statusOk) "OK" else "FAIL",
                color = if (check.statusOk) AccentGreen else AccentRed,
                fontSize = 12.sp,
                fontWeight = FontWeight.Bold
            )
        }
        val responseColor = when {
            check.responseTime < 200 -> AccentGreen
            check.responseTime < 1000 -> AccentAmber
            else -> AccentRed
        }
        Text(
            "${check.responseTime.roundToInt()} ms",
            color = responseColor,
            fontSize = 12.sp,
            fontWeight = FontWeight.SemiBold,
            modifier = Modifier.weight(1f)
        )
        Text(
            "${check.timeInterval}s",
            color = TextDim,
            fontSize = 12.sp,
            modifier = Modifier.weight(1f)
        )
    }
}

@Composable
private fun ResponseTimeChart(checks: List<CheckResult>, modifier: Modifier = Modifier) {
    val maxVal = checks.maxOf { it.responseTime }.coerceAtLeast(1.0)
    val recent = checks.takeLast(30)

    Canvas(modifier = modifier) {
        val barWidth = size.width / recent.size - 4f
        recent.forEachIndexed { idx, check ->
            val height = (check.responseTime / maxVal * size.height).toFloat().coerceAtLeast(4f)
            val x = idx * (size.width / recent.size)
            val color = when {
                check.responseTime < 200 -> AccentGreen
                check.responseTime < 1000 -> AccentAmber
                else -> AccentRed
            }
            drawRoundRect(
                color = color.copy(alpha = if (check.statusOk) 0.8f else 0.4f),
                topLeft = androidx.compose.ui.geometry.Offset(x + 2, size.height - height),
                size = androidx.compose.ui.geometry.Size(barWidth, height),
                cornerRadius = androidx.compose.ui.geometry.CornerRadius(3f)
            )
        }
    }
}
