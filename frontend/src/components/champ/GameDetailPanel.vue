<script setup lang="ts">
import { computed } from "vue"
import type { types } from "../../../wailsjs/go/models"
import { compactNumber, formatDuration, formatKda, queueLabel } from "@/lib/format"

const props = defineProps<{ detail: types.GameDetail }>()

const teams = computed(() => props.detail.teams ?? [])
const maxDamage = computed(() =>
  Math.max(1, ...teams.value.flatMap((t) => (t.players ?? []).map((p) => p.damage ?? 0))),
)
const damageWidth = (damage: number) => `${Math.round(((damage ?? 0) / maxDamage.value) * 100)}%`
</script>

<template>
  <div class="space-y-3 rounded-lg border bg-background/60 p-3">
    <p class="text-xs text-muted-foreground">
      {{ queueLabel(detail.queueId) }} · 用时 {{ formatDuration(detail.gameDuration) }}
    </p>

    <section v-for="team in teams" :key="team.teamId" class="space-y-1">
      <header class="flex items-center gap-3 text-xs">
        <span
          class="rounded px-1.5 py-0.5 font-semibold text-white"
          :class="team.win ? 'bg-emerald-500' : 'bg-rose-500'"
        >
          {{ team.win ? "胜利" : "失败" }}
        </span>
        <span class="text-muted-foreground">击杀 {{ team.kills }} · 经济 {{ compactNumber(team.gold) }}</span>
      </header>

      <div
        v-for="p in team.players"
        :key="p.puuid || p.name"
        class="grid grid-cols-[28px_minmax(110px,1.2fr)_84px_minmax(90px,1fr)_70px_auto] items-center gap-2 rounded px-1 py-1 text-xs"
        :class="p.isTarget ? 'bg-primary/10 ring-1 ring-primary/40' : ''"
      >
        <img v-if="p.championIcon" :src="p.championIcon" :alt="p.championName" class="size-7 rounded" />
        <div v-else class="size-7 rounded bg-muted" />

        <div class="min-w-0">
          <p class="truncate font-medium text-foreground">{{ p.name || "未知玩家" }}</p>
          <p class="truncate text-muted-foreground">{{ p.championName }} · Lv{{ p.level }} · {{ (p.spells ?? []).join(" ") }}</p>
        </div>

        <p class="font-semibold text-foreground">{{ formatKda(p.kills, p.deaths, p.assists) }}</p>

        <div>
          <div class="h-1.5 overflow-hidden rounded bg-muted">
            <div class="h-full bg-orange-400" :style="{ width: damageWidth(p.damage) }" />
          </div>
          <p class="mt-0.5 text-muted-foreground">伤害 {{ compactNumber(p.damage) }}</p>
        </div>

        <p class="text-muted-foreground">
          补刀 {{ p.cs }}<br />
          视野 {{ p.visionScore }}
        </p>

        <div class="flex gap-0.5">
          <template v-for="(item, i) in p.items" :key="i">
            <img v-if="item" :src="item" alt="" class="size-6 rounded" />
            <div v-else class="size-6 rounded bg-muted/60" />
          </template>
        </div>
      </div>
    </section>
  </div>
</template>
