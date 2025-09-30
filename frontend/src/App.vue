<script lang="ts">
export const description = "A sidebar that collapses to icons."
export const iframeHeight = "800px"
export const containerClass = "w-full h-full"
</script>

<script setup lang="ts">
import { onMounted, ref } from "vue"
import AppSidebar from "@/components/AppSidebar.vue"
import HelloWorld from "@/components/HelloWorld.vue"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar"
import { GetImgSrc, Greet } from "../wailsjs/go/main/App"
import { IplayerBaseData } from "@/interface/baseData"
import { GetPlayerRankData } from "../wailsjs/go/controller/PlayerController"
import type { IRankedStats } from "@/interface/rankData"

const playerData = ref<IplayerBaseData>({} as IplayerBaseData)
const rankData = ref<IRankedStats>()
const activePanel = ref<"dashboard" | "helloWorld">("dashboard")

onMounted(() => {
  loadPlayerData()
})

const parsePlayerData = (raw: unknown): Partial<IplayerBaseData> => {
  if (typeof raw === "string") {
    try {
      return JSON.parse(raw) as Partial<IplayerBaseData>
    } catch (error) {
      console.warn("failed to parse player data", error)
      return {}
    }
  }

  if (typeof raw === "object" && raw !== null) {
    return raw as Partial<IplayerBaseData>
  }

  return {}
}

async function loadPlayerData() {
  try {
    const [rawBaseData, iconPath] = await Promise.all([
      Greet("test"),
      GetImgSrc(),
    ])

    playerData.value = {
      ...parsePlayerData(rawBaseData),
      iconImgSrc: iconPath as string,
    }

    console.log("player data", playerData.value)

    const rankIdentifier =
      playerData.value.puuid ??
      playerData.value.accountId ??
      (playerData.value.summonerId != null
        ? String(playerData.value.summonerId)
        : undefined)

    if (rankIdentifier) {
      try {
        rankData.value = await GetPlayerRankData(rankIdentifier)
        console.log("player rank data", rankData.value)
      } catch (rankError) {
        console.error("failed to fetch player rank data", rankError)
      }
    } else {
      console.warn("missing identifier for rank data")
    }
  } catch (error) {
    console.error("failed to initialise player data", error)
  }
}

type SidebarNavigatePayload = {
  parent: {
    title: string
    url: string
  }
  item: {
    title: string
    url: string
  }
}

const handleSidebarNavigate = (payload: SidebarNavigatePayload) => {
  if (payload.item.url === "hello-world") {
    activePanel.value = "helloWorld"
  } else {
    activePanel.value = "dashboard"
  }
}
</script>

<template>
  <SidebarProvider>
    <AppSidebar
      :avatar-src="playerData.iconImgSrc"
      :player-data="playerData"
      :rank-data="rankData"
      @navigate="handleSidebarNavigate"
    />
    <SidebarInset>
      <header
        class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
      >
        <div class="flex items-center gap-2 px-4">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="mr-2 h-4" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem class="hidden md:block">
                <BreadcrumbLink href="#">
                  Building Your Application
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator class="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>Data Fetching</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4 pt-0">
        <template v-if="activePanel === 'helloWorld'">
          <HelloWorld />
        </template>
        <template v-else>
          <div class="grid auto-rows-min gap-4 md:grid-cols-3">
            <div class="aspect-video rounded-xl bg-muted/50" />
            <div class="aspect-video rounded-xl bg-muted/50" />
            <div class="aspect-video rounded-xl bg-muted/50" />
          </div>
          <div class="min-h-[100vh] flex-1 rounded-xl bg-muted/50 md:min-h-min" />
        </template>
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
