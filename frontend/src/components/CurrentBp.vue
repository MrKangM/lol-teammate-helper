<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue"

import { GetCurrentChampSelectSnapshot } from "../../wailsjs/go/main/App"
import { types } from "../../wailsjs/go/models"
import { EventsOn } from "../../wailsjs/runtime"
import { PHASE_LABELS, queueLabel } from "@/lib/format"
import PlayerRow from "@/components/champ/PlayerRow.vue"

type Side = "team" | "enemy"

const snapshot = ref<types.ChampSelectSnapshot | null>(null)
const side = ref<Side>("team")
const expandedKey = ref<string | null>(null)

const members = computed(() => snapshot.value?.[side.value] ?? [])
const hasSnapshot = computed(() => (snapshot.value?.team?.length ?? 0) > 0)
const hasEnemy = computed(() => (snapshot.value?.enemy?.length ?? 0) > 0)
const phaseLabel = computed(() => PHASE_LABELS[snapshot.value?.phase ?? ""] ?? "")
const modeLabel = computed(() => (snapshot.value?.queueId ? queueLabel(snapshot.value.queueId) : ""))

const rowKey = (m: types.TeamMemberSummary) => `${side.value}:${m.puuid || m.cellId}`
const toggle = (m: types.TeamMemberSummary) => {
  const key = rowKey(m)
  expandedKey.value = expandedKey.value === key ? null : key
}

const applySnapshot = (payload: unknown) => {
  const next = types.ChampSelectSnapshot.createFrom(payload)
  snapshot.value = next
  if (side.value === "enemy" && !(next.enemy?.length ?? 0)) {
    side.value = "team"
  }
}

let stopSnapshot: (() => void) | null = null
let stopEnded: (() => void) | null = null

onMounted(async () => {
  // Events are pushed by the backend; ask for the current state once for late openers.
  stopSnapshot = EventsOn("champ-select:snapshot", applySnapshot)
  stopEnded = EventsOn("champ-select:ended", () => {
    snapshot.value = null
    expandedKey.value = null
    side.value = "team"
  })

  if (typeof (window as any)?.go?.main?.App?.GetCurrentChampSelectSnapshot === "function") {
    try {
      const current = await GetCurrentChampSelectSnapshot()
      if (current?.team?.length && !snapshot.value) {
        applySnapshot(current)
      }
    } catch (error) {
      console.warn("[CurrentBp] failed to load current snapshot", error)
    }
  }
})

onUnmounted(() => {
  stopSnapshot?.()
  stopEnded?.()
})
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-2">
      <div>
        <h2 class="text-2xl font-semibold text-foreground">对局队友分析</h2>
        <p class="text-sm text-muted-foreground">
          基于最近 20 场单双排/灵活组排战绩。点击一行查看战绩列表，再点某一场查看双方结局图。
        </p>
      </div>
      <div v-if="hasSnapshot" class="flex items-center gap-2 text-xs text-muted-foreground">
        <span v-if="phaseLabel" class="rounded bg-primary/10 px-2 py-0.5 font-medium text-primary">{{ phaseLabel }}</span>
        <span v-if="modeLabel">{{ modeLabel }}</span>
      </div>
    </div>

    <div v-if="hasSnapshot" class="flex gap-1">
      <button
        type="button"
        class="rounded-md px-3 py-1 text-sm transition-colors"
        :class="side === 'team' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground hover:text-foreground'"
        @click="side = 'team'"
      >
        我方
      </button>
      <button
        type="button"
        class="rounded-md px-3 py-1 text-sm transition-colors disabled:cursor-not-allowed disabled:opacity-50"
        :class="side === 'enemy' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground hover:text-foreground'"
        :disabled="!hasEnemy"
        :title="hasEnemy ? '' : '排位选人阶段看不到对手，进入加载界面后会自动出现'"
        @click="side = 'enemy'"
      >
        对手
      </button>
    </div>

    <div v-if="hasSnapshot" class="space-y-2">
      <div
        class="grid grid-cols-[52px_minmax(170px,1.5fr)_minmax(170px,1.2fr)_minmax(170px,1.2fr)_minmax(140px,1fr)_84px] gap-3 px-3 text-xs text-muted-foreground"
      >
        <span>位置</span><span>召唤师</span><span>段位</span><span>近期战绩</span><span>英雄熟练度</span><span class="text-right">评价</span>
      </div>
      <PlayerRow
        v-for="m in members"
        :key="rowKey(m)"
        :member="m"
        :expanded="expandedKey === rowKey(m)"
        @toggle="toggle(m)"
      />
      <p class="px-1 pt-1 text-xs text-muted-foreground">
        评价为综合近期胜率、KDA、段位和英雄熟练度的估算（大腿 ≥ 68，上等马 ≥ 55，中等马 ≥ 42，其余为下等马），仅供参考。
      </p>
    </div>

    <div v-else class="rounded-lg border border-dashed p-10 text-center text-sm text-muted-foreground">
      等待进入英雄选择阶段…进入后这里会自动显示队友的段位、战绩和评价。
    </div>
  </div>
</template>
