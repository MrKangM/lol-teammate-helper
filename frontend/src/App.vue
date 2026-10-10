<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { Activity, Swords, UserRound } from "lucide-vue-next"
import { PHASE_LABELS } from "@/lib/format"
import { startStore, status, summoner, summonerIcon } from "@/lib/store"
import CareerView from "@/views/CareerView.vue"
import DiagnosticsView from "@/views/DiagnosticsView.vue"
import LiveView from "@/views/LiveView.vue"

type Page = "live" | "career" | "diagnostics"

const page = ref<Page>("live")

const nav = [
  { id: "live" as Page, label: "对局分析", icon: Swords },
  { id: "career" as Page, label: "我的生涯", icon: UserRound },
  { id: "diagnostics" as Page, label: "诊断", icon: Activity },
]

const phaseLabel = computed(() => PHASE_LABELS[status.value?.phase ?? ""] ?? "")

const goto = (next: string) => {
  page.value = next as Page
}

onMounted(startStore)
</script>

<template>
  <div class="flex h-full flex-col">
    <header class="flex h-14 shrink-0 items-center gap-6 border-b border-line bg-panel px-5">
      <div class="flex items-center gap-2.5">
        <div class="grid size-8 place-items-center overflow-hidden rounded-lg bg-accent text-sm font-black text-bg">
          <img v-if="summonerIcon" :src="summonerIcon" alt="" class="size-full object-cover" />
          <span v-else>LH</span>
        </div>
        <span class="text-sm font-semibold text-ink">队友助手</span>
      </div>

      <nav class="flex h-full items-stretch gap-1">
        <button
          v-for="item in nav"
          :key="item.id"
          type="button"
          class="relative flex items-center gap-2 px-3 text-[13px] transition-colors"
          :class="page === item.id ? 'font-semibold text-ink' : 'text-muted hover:text-ink'"
          @click="page = item.id"
        >
          <component :is="item.icon" class="size-4" />
          {{ item.label }}
          <span v-if="page === item.id" class="absolute inset-x-2 bottom-0 h-[2px] rounded-full bg-accent" />
        </button>
      </nav>

      <div class="ml-auto flex items-center gap-3 text-xs">
        <span v-if="summoner?.gameName" class="text-muted">{{ summoner.gameName }}</span>
        <span v-if="status?.connected && phaseLabel" class="rounded bg-panel-3 px-2 py-0.5 text-ink">{{ phaseLabel }}</span>
        <span class="flex items-center gap-1.5 rounded-full bg-panel-2 px-2.5 py-1" :class="status?.connected ? 'text-good' : 'text-loss'">
          <span class="size-1.5 rounded-full" :class="status?.connected ? 'bg-good' : 'bg-loss'" />
          {{ status?.connected ? "客户端已连接" : "未连接" }}
        </span>
      </div>
    </header>

    <main class="min-h-0 flex-1 overflow-y-auto p-5">
      <LiveView v-show="page === 'live'" @goto="goto" />
      <CareerView v-if="page === 'career'" />
      <DiagnosticsView v-if="page === 'diagnostics'" />
    </main>
  </div>
</template>
