<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { PHASE_LABELS, RATING_STYLES, queueLabel } from "@/lib/format"
import { snapshot, status } from "@/lib/store"
import PlayerCard from "@/components/PlayerCard.vue"
import PlayerDrawer from "@/components/PlayerDrawer.vue"
import type { types } from "../../wailsjs/go/models"

const emit = defineEmits<{ (e: "goto", page: string): void }>()

const selected = ref<types.TeamMemberSummary | null>(null)

const team = computed(() => snapshot.value?.team ?? [])
const enemy = computed(() => snapshot.value?.enemy ?? [])
const hasSnapshot = computed(() => team.value.length > 0)
const modeLabel = computed(() => (snapshot.value?.queueId ? queueLabel(snapshot.value.queueId) : ""))
const phaseLabel = computed(() => PHASE_LABELS[snapshot.value?.phase ?? ""] ?? snapshot.value?.phase ?? "")

const HORSES = ["大腿", "上等马", "中等马", "下等马"]

const summarise = (members: types.TeamMemberSummary[]) => {
  const rated = members.filter((m) => m.rating?.valid)
  const avg = rated.length ? rated.reduce((sum, m) => sum + m.rating.score, 0) / rated.length : null
  const counts = HORSES.map((label) => ({ label, n: members.filter((m) => m.rating?.label === label).length })).filter((c) => c.n > 0)
  return { avg, counts }
}
const teamSummary = computed(() => summarise(team.value))
const enemySummary = computed(() => summarise(enemy.value))

// Keep the drawer showing fresh data when a new snapshot arrives for the same player.
watch(snapshot, (next) => {
  if (!selected.value) return
  const all = [...(next?.team ?? []), ...(next?.enemy ?? [])]
  selected.value = all.find((m) => m.puuid === selected.value?.puuid) ?? null
})
</script>

<template>
  <div>
    <div v-if="hasSnapshot" class="grid grid-cols-1 gap-5 xl:grid-cols-2">
      <section v-for="col in [
        { key: 'team' as const, title: '我方', dot: 'bg-win', members: team, summary: teamSummary },
        { key: 'enemy' as const, title: '对手', dot: 'bg-loss', members: enemy, summary: enemySummary },
      ]" :key="col.key" class="min-w-0">
        <header class="mb-2.5 flex flex-wrap items-center gap-x-3 gap-y-1 px-1">
          <span class="size-2.5 rounded-full" :class="col.dot" />
          <h2 class="text-sm font-semibold text-ink">{{ col.title }}</h2>
          <span v-if="col.summary.avg !== null" class="num text-xs text-muted">平均 {{ col.summary.avg.toFixed(0) }} 分</span>
          <span v-for="c in col.summary.counts" :key="c.label" class="rounded px-1.5 text-[11px] font-semibold leading-5" :class="RATING_STYLES[c.label]?.badge">
            {{ c.label }} ×{{ c.n }}
          </span>
          <span v-if="col.key === 'team'" class="ml-auto flex items-center gap-2 text-xs text-muted">
            <span v-if="phaseLabel" class="rounded bg-panel-3 px-2 py-0.5 text-ink">{{ phaseLabel }}</span>
            <span v-if="modeLabel">{{ modeLabel }}</span>
          </span>
        </header>

        <div v-if="col.members.length" class="space-y-2">
          <PlayerCard v-for="m in col.members" :key="`${col.key}:${m.puuid || m.cellId}`" :member="m" :side="col.key" @open="selected = m" />
        </div>
        <div v-else class="card flex min-h-[260px] flex-col items-center justify-center gap-2 border-dashed px-8 text-center">
          <p class="text-sm font-semibold text-ink">对手信息暂不可见</p>
          <p class="max-w-xs text-xs leading-5 text-muted">排位选人阶段客户端不会提供对手的账号信息。进入加载界面后，这里会自动显示对手的段位和战绩。</p>
        </div>
      </section>
    </div>

    <p v-if="hasSnapshot" class="px-1 pt-4 text-[11px] leading-5 text-muted">
      评价综合近 20 场排位的胜率、KDA、段位和当前英雄熟练度估算：大腿 ≥ 68，上等马 ≥ 55，中等马 ≥ 42，其余为下等马，仅供参考。点击卡片查看详细战绩。
    </p>

    <div v-else class="card flex flex-col items-center gap-3 px-6 py-20 text-center">
      <template v-if="status?.connected">
        <p class="text-lg font-semibold text-ink">等待进入对局</p>
        <p class="max-w-md text-sm text-muted">客户端已连接。进入英雄选择后，这里会自动显示队友的段位、战绩和评价；进入加载界面后可以看到对手。</p>
      </template>
      <template v-else>
        <p class="text-lg font-semibold text-ink">尚未连接到英雄联盟客户端</p>
        <p class="max-w-md text-sm text-muted">请先启动并登录客户端。软件会自动重连；如果长时间连不上，到“诊断”页查看原因。</p>
        <button type="button" class="btn" @click="emit('goto', 'diagnostics')">查看诊断</button>
      </template>
    </div>

    <PlayerDrawer v-if="selected" :member="selected" @close="selected = null" />
  </div>
</template>
