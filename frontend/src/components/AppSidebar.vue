<script setup lang="ts">
import { computed } from "vue"
import type { SidebarProps } from "@/components/ui/sidebar"
import type { IPlayerBaseData } from "@/interface/baseData"
import type { IRankedStats } from "@/interface/rankData"
import type { LucideIcon } from "lucide-vue-next"

import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
} from "lucide-vue-next"
import NavMain from "@/components/NavMain.vue"
import NavProjects from "@/components/NavProjects.vue"
import NavUser from "@/components/NavUser.vue"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar"

type NavSubItem = {
  title: string
  url: string
}

type NavItem = NavSubItem & {
  icon?: LucideIcon
  isActive?: boolean
  items?: NavSubItem[]
}

type AppSidebarProps = SidebarProps & {
  avatarSrc?: string
  playerData?: IPlayerBaseData
  rankData?: IRankedStats
}

const emit = defineEmits<{
  (event: "navigate", payload: { parent: NavItem; item: NavSubItem }): void
}>()

const props = withDefaults(defineProps<AppSidebarProps>(), {
  collapsible: "icon",
  avatarSrc: "",
  playerData: undefined,
  rankData: undefined,
})

const handleNavNavigate = (payload: { parent: NavItem; item: NavSubItem }) => {
  emit("navigate", payload)
}

const sanitize = (value?: string | null) => {
  const trimmed = value?.trim()
  return trimmed && trimmed.length > 0 ? trimmed : undefined
}

const resolvePlayerName = (player?: IPlayerBaseData) => {
  if (!player) return undefined

  const displayName = sanitize(player.displayName)
  if (displayName) return displayName

  const gameName = sanitize(player.gameName)
  const tagLine = sanitize(player.tagLine)
  if (gameName && tagLine) return `${gameName}#${tagLine}`
  if (gameName) return gameName

  return sanitize(player.internalName)
}

const resolvePlayerRankLabel = (rankData?: IRankedStats): string | undefined => {
  const entry = rankData?.queueMap?.["RANKED_SOLO_5x5"]
  if (!entry) {
    console.log(`数组不存在，${JSON.stringify(rankData?.queueMap)}`)
    return undefined
  }

  const tier = sanitize(entry.tier)
  const division = sanitize(entry.division)
  console.log(`段位:${tier},阶级:${division}`)
  if (tier && division) {
    return `${tier} ${division}`
  }

  return tier ?? division ?? undefined
}

const defaultSidebarData = {
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  teams: [
    {
      name: "Acme Inc",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
    {
      name: "Acme Corp.",
      logo: AudioWaveform,
      plan: "Startup",
    },
    {
      name: "Evil Corp.",
      logo: Command,
      plan: "Free",
    },
  ],
  navMain: [
    {
      title: "Dashboard",
      url: "#",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "History",
          url: "hello-world",
        },
        {
          title: "Starred",
          url: "#",
        },
        {
          title: "Settings",
          url: "#",
        },
      ],
    },
    {
      title: "Models",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Genesis",
          url: "#",
        },
        {
          title: "Explorer",
          url: "#",
        },
        {
          title: "Quantum",
          url: "#",
        },
      ],
    },
    {
      title: "Documentation",
      url: "#",
      icon: BookOpen,
      items: [
        {
          title: "Introduction",
          url: "#",
        },
        {
          title: "Get Started",
          url: "#",
        },
        {
          title: "Tutorials",
          url: "#",
        },
        {
          title: "Changelog",
          url: "#",
        },
      ],
    },
    {
      title: "Settings",
      url: "#",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "#",
        },
        {
          title: "Team",
          url: "#",
        },
        {
          title: "Billing",
          url: "#",
        },
        {
          title: "Limits",
          url: "#",
        },
      ],
    },
  ],
  projects: [
    {
      name: "Design Engineering",
      url: "#",
      icon: Frame,
    },
    {
      name: "Sales & Marketing",
      url: "#",
      icon: PieChart,
    },
    {
      name: "Travel",
      url: "#",
      icon: Map,
    },
  ],
}

const resolvedUser = computed(() => {
  const player = props.playerData
  const fallback = defaultSidebarData.user
  const name = resolvePlayerName(player) ?? fallback.name

  const rankLabel = resolvePlayerRankLabel(props.rankData)
  console.log(`排位标签:${rankLabel}`)
  const identifier = rankLabel ?? "未定级"
  const avatar = sanitize(player?.iconImgSrc) ?? props.avatarSrc ?? fallback.avatar

  return {
    ...fallback,
    name,
    email: identifier,
    avatar,
  }
})

const sidebarData = computed(() => ({
  ...defaultSidebarData,
  user: resolvedUser.value,
}))
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <NavUser :user="sidebarData.user" />
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="sidebarData.navMain" @navigate="handleNavNavigate" />
      <NavProjects :projects="sidebarData.projects" />
    </SidebarContent>
    <SidebarFooter />
    <SidebarRail />
  </Sidebar>
</template>
