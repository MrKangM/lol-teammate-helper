<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { PHASE_LABELS, queueLabel } from "@/lib/format"
import { snapshot, status } from "@/lib/store"
import PlayerCard from "@/components/PlayerCard.vue"
import type { types } from "../../wailsjs/go/models"

const emit = defineEmits<{ (e: "goto", page: string): void }>()

type Side = "team" | "enemy"
const side = ref<Side>("team")
const expandedKey = ref<string | null>(null)

const members = computed(() => snapshot.value?.[side.value] ?? [])
const hasSnapshot = computed(() => (snapshot.value?.team?.length ?? 0) > 0)
const hasEnemy = computed(() => (snapshot.value?.enemy?.length ?? 0) > 0)
const modeLabel = computed(() => (snapshot.value?.queueId ? queueLabel(snapshot.value.queueId) : ""))
const phaseLabel = computed(() => PHASE_LABELS[snapshot.value?.phase ?? ""] ?? snapshot.value?.phase ?? "")

watch(hasEnemy, (has) => {
  if (!has) side.value = "team"
})
watch(hasSnapshot, (has) => {
  if (!has) expandedKey.value = null
})

const rowKey = (m: types.TeamMemberSummary) => `${side.value}:${m.puuid || m.cellId}`
const toggle = (m: types.TeamMemberSummary) => {
  const key = rowKey(m)
  expandedKey.value = expandedKey.value === key ? null : key
}
</script>

<template>
  <div class="space-y-4">
    <div v-if="hasSnapshot" class="flex flex-wrap items-center justify-between gap-3">
      <div class="flex gap-px rounded-md border border-line-gold p-0.5">
        <button
          type="button"
          class="rounded px-4 py-1 text-sm transition-colors"
          :class="side === 'team' ? 'bg-gold/20 font-semibold text-gold-bright' : 'text-muted hover:text-ink'"
          @click="side = 'team'"
        >
          我方
        </button>
        <button
          type="button"
          class="rounded px-4 py-1 text-sm transition-colors disabled:opacity-40"
          :class="side === 'enemy' ? 'bg-gold/20 font-semibold text-gold-bright' : 'text-muted hover:text-ink'"
          :disabled="!hasEnemy"
          :title="hasEnemy ? '' : '排位选人阶段客户端不会提供对手信息，进入加载界面后会自动出现'"
          @click="side = 'enemy'"
        >
          对手
        </button>
      </div>
      <div class="flex items-center gap-2 text-xs text-muted">
        <span v-if="phaseLabel" class="rounded-sm border border-teal/40 bg-teal/10 px-2 py-0.5 font-medium text-teal">{{ phaseLabel }}</span>
        <span v-if="modeLabel">{{ modeLabel }}</span>
      </div>
    </div>

    <div v-if="hasSnapshot" class="space-y-2">
      <PlayerCard v-for="m in members" :key="rowKey(m)" :member="m" :expanded="expandedKey === rowKey(m)" @toggle="toggle(m)" />
      <p class="px-1 pt-1 text-[11px] leading-5 text-muted">
        评价综合近 20 场排位的胜率、KDA、段位和当前英雄熟练度估算：大腿 ≥ 68，上等马 ≥ 55，中等马 ≥ 42，其余为下等马，仅供参考。点击一行展开战绩，再点某一场查看双方结局图。
      </p>
    </div>

    <div v-else class="panel flex flex-col items-center gap-3 px-6 py-16 text-center">
      <template v-if="status?.connected">
        <p class="title-gold text-lg">等待进入对局</p>
        <p class="max-w-md text-sm text-muted">客户端已连接。进入英雄选择后，这里会自动显示队友的段位、战绩和评价；进入加载界面后可以查看对手。</p>
      </template>
      <template v-else>
        <p class="title-gold text-lg">尚未连接到英雄联盟客户端</p>
        <p class="max-w-md text-sm text-muted">请先启动并登录客户端。软件会自动重连；如果长时间连不上，到“诊断”页查看原因。</p>
        <button type="button" class="btn" @click="emit('goto', 'diagnostics')">查看诊断</button>
      </template>
    </div>
  </div>
</template>
