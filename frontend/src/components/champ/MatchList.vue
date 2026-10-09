<script setup lang="ts">
import { ref } from "vue"
import { GetGameDetail } from "../../../wailsjs/go/controller/MatchHistory"
import type { types } from "../../../wailsjs/go/models"
import { compactNumber, formatDuration, formatKda, kdaRatio, positionLabel, queueLabel, timeAgo } from "@/lib/format"
import GameDetailPanel from "./GameDetailPanel.vue"

const props = defineProps<{ matches: types.RecentMatchSummary[]; puuid: string }>()

const openGameId = ref<number | null>(null)
const loadingId = ref<number | null>(null)
const errorText = ref("")
const details = ref<Record<number, types.GameDetail>>({})

const toggle = async (gameId: number) => {
  errorText.value = ""
  if (openGameId.value === gameId) {
    openGameId.value = null
    return
  }
  openGameId.value = gameId
  if (details.value[gameId]) return

  loadingId.value = gameId
  try {
    details.value[gameId] = await GetGameDetail(gameId, props.puuid)
  } catch (error) {
    console.error("failed to load game detail", error)
    errorText.value = "对局详情加载失败，可能已超出客户端保留的历史。"
  } finally {
    loadingId.value = null
  }
}
</script>

<template>
  <div class="space-y-1">
    <p v-if="matches.length === 0" class="py-2 text-xs text-muted-foreground">暂无排位对局记录</p>

    <div v-for="m in matches" :key="m.gameId">
      <button
        type="button"
        class="grid w-full grid-cols-[28px_44px_minmax(90px,1fr)_92px_64px_120px_110px_70px] items-center gap-2 rounded px-2 py-1.5 text-left text-xs transition-colors hover:bg-muted/60"
        :class="m.win ? 'border-l-2 border-emerald-500 bg-emerald-500/5' : 'border-l-2 border-rose-500 bg-rose-500/5'"
        @click="toggle(m.gameId)"
      >
        <img v-if="m.championIcon" :src="m.championIcon" :alt="m.championName" class="size-7 rounded" />
        <div v-else class="size-7 rounded bg-muted" />

        <span class="font-semibold" :class="m.win ? 'text-emerald-500' : 'text-rose-500'">{{ m.win ? "胜利" : "失败" }}</span>
        <span class="truncate text-foreground">{{ m.championName || `英雄${m.championId}` }} · {{ positionLabel(m.position) }}</span>
        <span class="font-semibold text-foreground">{{ formatKda(m.kills, m.deaths, m.assists) }}</span>
        <span class="text-muted-foreground">{{ kdaRatio(m.kills, m.deaths, m.assists) }} KDA</span>
        <span class="whitespace-nowrap text-muted-foreground">补刀 {{ m.cs }} · {{ compactNumber(m.damage) }}</span>
        <span class="whitespace-nowrap text-muted-foreground">{{ queueLabel(m.queueId) }} {{ formatDuration(m.gameDuration) }}</span>
        <span class="text-right text-muted-foreground">{{ timeAgo(m.gameCreation) }}</span>
      </button>

      <div v-if="openGameId === m.gameId" class="px-1 py-2">
        <p v-if="loadingId === m.gameId" class="text-xs text-muted-foreground">正在加载对局详情…</p>
        <p v-else-if="errorText" class="text-xs text-rose-500">{{ errorText }}</p>
        <GameDetailPanel v-else-if="details[m.gameId]" :detail="details[m.gameId]" />
      </div>
    </div>
  </div>
</template>
