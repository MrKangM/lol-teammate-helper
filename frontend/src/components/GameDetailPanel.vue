<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../wailsjs/go/models"
import { compactNumber, formatDuration, formatKda, queueLabel } from "@/lib/format"

const props = defineProps<{ detail: types.GameDetail }>()

const teams = computed(() => props.detail.teams ?? [])
const maxDamage = computed(() => Math.max(1, ...teams.value.flatMap((t) => (t.players ?? []).map((p) => p.damage ?? 0))))
const damageWidth = (damage: number) => `${Math.round(((damage ?? 0) / maxDamage.value) * 100)}%`
</script>

<template>
  <div class="space-y-3 rounded-md border border-line bg-bg/60 p-3">
    <p class="text-xs text-muted">{{ queueLabel(detail.queueId) }} · 用时 {{ formatDuration(detail.gameDuration) }}</p>

    <section v-for="team in teams" :key="team.teamId" class="space-y-1">
      <header class="flex items-center gap-3 text-xs">
        <span class="rounded-sm px-1.5 py-0.5 font-semibold" :class="team.win ? 'bg-win/20 text-win' : 'bg-loss/20 text-loss'">
          {{ team.win ? "胜利" : "失败" }}
        </span>
        <span class="text-muted">击杀 {{ team.kills }} · 经济 {{ compactNumber(team.gold) }}</span>
      </header>

      <div
        v-for="p in team.players"
        :key="p.puuid || p.name"
        class="grid grid-cols-[30px_44px_minmax(120px,1.2fr)_84px_minmax(80px,1fr)_62px_auto] items-center gap-2 rounded px-1.5 py-1 text-xs"
        :class="p.isTarget ? 'bg-gold/10 ring-1 ring-gold/40' : 'hover:bg-panel-2'"
      >
        <img v-if="p.championIcon" :src="p.championIcon" :alt="p.championName" class="size-[30px] rounded" />
        <div v-else class="size-[30px] rounded bg-panel-3" />

        <div class="flex gap-0.5">
          <template v-for="(s, i) in p.spells" :key="i">
            <img v-if="s?.icon" :src="s.icon" :alt="s.name" :title="s.name" class="size-[20px] rounded-sm" />
            <span v-else class="size-[20px] rounded-sm bg-panel-3 text-[9px] leading-[20px] text-muted text-center" :title="s?.name">{{ (s?.name || "").slice(0, 1) }}</span>
          </template>
        </div>

        <div class="min-w-0">
          <p class="truncate font-medium text-ink">{{ p.name || "未知玩家" }}</p>
          <p class="truncate text-muted">{{ p.championName }} · Lv{{ p.level }}</p>
        </div>

        <p class="num font-semibold text-gold-bright">{{ formatKda(p.kills, p.deaths, p.assists) }}</p>

        <div>
          <div class="h-1.5 overflow-hidden rounded bg-panel-3">
            <div class="h-full bg-gold" :style="{ width: damageWidth(p.damage) }" />
          </div>
          <p class="num mt-0.5 text-muted">伤害 {{ compactNumber(p.damage) }}</p>
        </div>

        <p class="num leading-tight text-muted">补刀 {{ p.cs }}<br />视野 {{ p.visionScore }}</p>

        <div class="flex gap-0.5">
          <template v-for="(item, i) in p.items" :key="i">
            <img v-if="item" :src="item" alt="" class="size-6 rounded-sm border border-line" />
            <div v-else class="size-6 rounded-sm border border-line/60 bg-panel-2" />
          </template>
        </div>
      </div>
    </section>
  </div>
</template>
