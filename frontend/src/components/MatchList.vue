<script setup lang="ts">
import { ref } from "vue"
import { GetGameDetail } from "../../wailsjs/go/controller/MatchHistory"
import type { types } from "../../wailsjs/go/models"
import { compactNumber, formatDuration, formatKda, kdaRatio, positionLabel, queueLabel, timeAgo } from "@/lib/format"
import GameDetailPanel from "@/components/GameDetailPanel.vue"

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
    <p v-if="matches.length === 0" class="py-3 text-center text-xs text-muted">暂无排位对局记录</p>

    <div v-for="m in matches" :key="m.gameId">
      <button
        type="button"
        class="grid w-full grid-cols-[4px_34px_minmax(0,1.3fr)_minmax(0,0.9fr)_minmax(0,1fr)_auto] items-center gap-3 overflow-hidden rounded-md py-1.5 pr-3 text-left text-xs transition-colors hover:brightness-125"
        :class="m.win ? 'bg-win-bg' : 'bg-loss-bg'"
        @click="toggle(m.gameId)"
      >
        <span class="h-full min-h-9 w-1" :class="m.win ? 'bg-win' : 'bg-loss'" />
        <img v-if="m.championIcon" :src="m.championIcon" :alt="m.championName" class="size-[34px] rounded-md" />
        <div v-else class="size-[34px] rounded-md bg-panel-3" />

        <div class="min-w-0 leading-tight">
          <p class="truncate font-semibold text-ink">{{ m.championName || `英雄${m.championId}` }}</p>
          <p class="truncate text-muted">{{ queueLabel(m.queueId) }} · {{ positionLabel(m.position) }}</p>
        </div>
        <div class="leading-tight">
          <p class="num font-semibold text-ink">{{ formatKda(m.kills, m.deaths, m.assists) }}</p>
          <p class="num text-muted">{{ kdaRatio(m.kills, m.deaths, m.assists) }} KDA</p>
        </div>
        <p class="num whitespace-nowrap leading-tight text-muted">
          补刀 {{ m.cs }}<br />
          伤害 {{ compactNumber(m.damage) }}
        </p>
        <div class="text-right leading-tight">
          <p class="font-semibold" :class="m.win ? 'text-win' : 'text-loss'">{{ m.win ? "胜利" : "失败" }}</p>
          <p class="num whitespace-nowrap text-muted">{{ formatDuration(m.gameDuration) }} · {{ timeAgo(m.gameCreation) }}</p>
        </div>
      </button>

      <div v-if="openGameId === m.gameId" class="py-2">
        <p v-if="loadingId === m.gameId" class="text-xs text-muted">正在加载对局详情…</p>
        <p v-else-if="errorText" class="text-xs text-loss">{{ errorText }}</p>
        <GameDetailPanel v-else-if="details[m.gameId]" :detail="details[m.gameId]" />
      </div>
    </div>
  </div>
</template>
