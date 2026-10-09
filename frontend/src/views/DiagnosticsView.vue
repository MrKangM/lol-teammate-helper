<script setup lang="ts">
import { computed, ref } from "vue"
import { Copy, FolderOpen } from "lucide-vue-next"
import { OpenLogDir } from "../../wailsjs/go/main/App"
import { PHASE_LABELS } from "@/lib/format"
import { snapshot, status } from "@/lib/store"

const copied = ref(false)
const openId = ref<number | null>(null)

const events = computed(() => status.value?.events ?? [])
const shortUri = (uri: string) => uri.replace("/lol-", "").replace("/v1/", "/")
const time = (value: unknown) => new Date(value as string).toLocaleTimeString()

const pretty = (data: string) => {
  try {
    return JSON.stringify(JSON.parse(data), null, 2)
  } catch {
    return data
  }
}

const bundle = computed(() => {
  const s = status.value
  const lines = [
    `connected: ${s?.connected}`,
    `port: ${s?.port}`,
    `region: ${s?.region}`,
    `phase: ${s?.phase}`,
    `lastError: ${s?.lastError}`,
    `logPath: ${s?.logPath}`,
    `snapshot: team=${snapshot.value?.team?.length ?? 0} enemy=${snapshot.value?.enemy?.length ?? 0} phase=${snapshot.value?.phase}`,
    "",
    ...events.value.slice(0, 6).map((e) => `[${time(e.time)}] ${e.uri} ${e.type} (${e.bytes} bytes)\n${e.data.slice(0, 6000)}`),
  ]
  return lines.join("\n")
})

const copy = async () => {
  try {
    await navigator.clipboard.writeText(bundle.value)
    copied.value = true
    window.setTimeout(() => (copied.value = false), 2000)
  } catch (error) {
    console.warn("copy failed", error)
  }
}
</script>

<template>
  <div class="space-y-4">
    <section class="panel p-4">
      <p class="eyebrow mb-3">连接状态</p>
      <dl class="grid grid-cols-[110px_1fr] gap-y-2 text-sm">
        <dt class="text-muted">客户端</dt>
        <dd :class="status?.connected ? 'text-win' : 'text-loss'">{{ status?.connected ? "已连接" : "未连接" }}</dd>
        <dt class="text-muted">端口 / 大区</dt>
        <dd class="num text-ink">{{ status?.port || "-" }} / {{ status?.region || "-" }}</dd>
        <dt class="text-muted">游戏阶段</dt>
        <dd class="text-ink">{{ PHASE_LABELS[status?.phase ?? ""] ?? status?.phase ?? "-" }}</dd>
        <dt class="text-muted">最近错误</dt>
        <dd class="break-all" :class="status?.lastError ? 'text-loss' : 'text-muted'">{{ status?.lastError || "无" }}</dd>
        <dt class="text-muted">当前数据</dt>
        <dd class="text-ink">我方 {{ snapshot?.team?.length ?? 0 }} 人 · 对手 {{ snapshot?.enemy?.length ?? 0 }} 人</dd>
        <dt class="text-muted">日志文件</dt>
        <dd class="break-all text-ink">{{ status?.logPath || "-" }}</dd>
      </dl>
      <div class="mt-4 flex gap-2">
        <button type="button" class="btn" @click="OpenLogDir()"><FolderOpen class="size-3.5" />打开日志目录</button>
        <button type="button" class="btn" @click="copy"><Copy class="size-3.5" />{{ copied ? "已复制" : "复制诊断信息" }}</button>
      </div>
      <p class="mt-3 text-[11px] leading-5 text-muted">反馈问题时点“复制诊断信息”，把内容贴给开发者即可，里面包含最近收到的原始事件。内容含对局玩家的名字，请自行确认后再发送。</p>
    </section>

    <section class="panel p-4">
      <p class="eyebrow mb-3">最近收到的事件（最新在前）</p>
      <p v-if="!events.length" class="py-6 text-center text-sm text-muted">还没有收到任何事件。进入房间、选人或游戏后会出现。</p>
      <div v-for="(e, i) in events" :key="i" class="border-b border-line/60 last:border-0">
        <button type="button" class="flex w-full items-center gap-3 py-1.5 text-left text-xs hover:text-gold-bright" @click="openId = openId === i ? null : i">
          <span class="num w-[68px] shrink-0 text-muted">{{ time(e.time) }}</span>
          <span class="truncate text-ink">{{ shortUri(e.uri) }}</span>
          <span class="rounded-sm border border-line px-1 text-muted">{{ e.type }}</span>
          <span class="num ml-auto text-muted">{{ e.bytes }} B</span>
        </button>
        <pre v-if="openId === i" class="mb-2 max-h-72 overflow-auto rounded-sm bg-bg p-2 text-[11px] leading-4 text-muted">{{ pretty(e.data) }}</pre>
      </div>
    </section>
  </div>
</template>
