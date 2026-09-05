"use client";

import {
  FileText,
  Folder,
  Home,
  Image as ImageIcon,
  Music,
  Settings,
  Star,
  Trash2,
  Video,
  Signpost,
  UserRound,
  Plus,
  XIcon
} from "lucide-react";
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarFooter
} from "@/components/ui/sidebar";
import { Link } from "react-router";

import { Progress } from "@/components/ui/progress"
import { BACKEND_BASE_URL } from "@/App";
import { useEffect, useState } from "react";

import { formatFileSize } from "@/lib/files";

const mainItems = [
  { title: "Home", icon: Home, link: "/" },
  { title: "Route Rules", icon: Signpost, link: "/.@/rules" },
  { title: "Create", icon: Plus, link: "/.@/create" },

  { title: "Account", icon: UserRound, link: "/.@/account" },
];

interface Usage {
  percent: number
  free: number
  base: number
}

export function AppSidebar() {

  const [usage, setUsage] = useState<Usage>(null)

  useEffect(()=>{
    fetch(`${BACKEND_BASE_URL}/info/disk-usage`, { credentials: "include" }).then(r => r.json()).then(res => {
      let newUsage: Usage = { percent: res.percent, free: res.free, base: res.base }
      setUsage(newUsage)
    })
  },[])
  

  return (
    <Sidebar variant="inset">
      <SidebarContent>

        <SidebarGroup>
          <SidebarGroupContent>
            <div className="flex flex-row">
              <div className="flex flex-col justify-center items-center ">
                <XIcon className="w-[2rem] h-[2rem]"></XIcon>
              </div>
              <div className="text-2xl font-semibold">Dirlist</div>
            </div>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {mainItems.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton >
                    <Link to={item.link} className="flex flex-row gap-4">
                      <item.icon className="h-4 w-4" />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        
      </SidebarContent>
      {usage !== null &&
        <SidebarFooter>
          <SidebarGroupLabel>Storage space</SidebarGroupLabel>
          <SidebarGroupContent>

            <div className="p-2 pt-4 pb-6 flex flex-col gap-4 bg-popover/10 rounded-lg">

              <div className="text-ring "><span>{formatFileSize(usage.free)}</span> of <span>{formatFileSize(usage.base)}</span> used</div>

              <Progress value={usage.percent}
                max={100}
                className="mx-auto w-full max-w-xs"
              ></Progress>
            </div>
          </SidebarGroupContent>
        </SidebarFooter>
      }
    </Sidebar>
  );
}