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
  { id: "live" as Page, label: "对局", icon: Swords },
  { id: "career" as Page, label: "生涯", icon: UserRound },
  { id: "diagnostics" as Page, label: "诊断", icon: Activity },
]
const titles: Record<Page, string> = { live: "对局分析", career: "我的生涯", diagnostics: "诊断" }

const phaseLabel = computed(() => PHASE_LABELS[status.value?.phase ?? ""] ?? "")

const goto = (next: string) => {
  page.value = next as Page
}

onMounted(startStore)
</script>

<template>
  <div class="flex h-full">
    <nav class="flex w-[76px] shrink-0 flex-col items-center gap-1 border-r border-line-gold/60 bg-panel/80 py-4">
      <div class="mb-4 grid size-11 place-items-center rounded-full border border-gold bg-bg">
        <img v-if="summonerIcon" :src="summonerIcon" alt="" class="size-full rounded-full object-cover" />
        <span v-else class="title-gold text-lg">LH</span>
      </div>
      <button
        v-for="item in nav"
        :key="item.id"
        type="button"
        class="group relative flex w-full flex-col items-center gap-1 py-2.5 text-[11px] transition-colors"
        :class="page === item.id ? 'text-gold-bright' : 'text-muted hover:text-ink'"
        @click="page = item.id"
      >
        <span v-if="page === item.id" class="absolute left-0 top-2 h-[calc(100%-1rem)] w-[3px] bg-gold" />
        <component :is="item.icon" class="size-5" />
        {{ item.label }}
      </button>
    </nav>

    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex h-14 shrink-0 items-center justify-between border-b border-line px-6">
        <div>
          <h1 class="title-gold text-lg leading-tight">{{ titles[page] }}</h1>
          <p class="eyebrow leading-none">Teammate Helper</p>
        </div>
        <div class="flex items-center gap-3 text-xs">
          <span v-if="summoner?.gameName" class="text-muted">{{ summoner.gameName }}</span>
          <span v-if="status?.connected && phaseLabel" class="rounded-sm border border-line px-2 py-0.5 text-muted">{{ phaseLabel }}</span>
          <span class="flex items-center gap-1.5" :class="status?.connected ? 'text-win' : 'text-loss'">
            <span class="size-2 rounded-full" :class="status?.connected ? 'bg-win shadow-[0_0_8px_var(--color-win)]' : 'bg-loss'" />
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
  </div>
</template>
