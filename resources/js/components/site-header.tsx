import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { ThemeToggleIcon } from "@/components/ThemeToggleIcon"
import { NotificationDrawer } from "@/components/Notifications/NotificationDrawer"
import { Button } from "@/components/ui/button"
import { MessageCircle } from "lucide-react"
import { useState } from "react"
import { MessageSidebar } from "@/components/Messages/MessageSidebar"
import { MessageChat } from "@/components/Messages/MessageChat"
import { useMessages } from "@/contexts/MessageContext"
import { usePage } from "@inertiajs/react"
import { SharedData } from "@/types/app"
import {
  Drawer,
  DrawerContent,
  DrawerHeader,
  DrawerTitle,
} from "@/components/ui/drawer"

export function SiteHeader({title}: { title: string }) {
  const [isMessagesOpen, setIsMessagesOpen] = useState(false);
  const { props } = usePage<SharedData>();
  const user = props.auth?.user;
  const { selectedConversation, unreadCount } = useMessages();

  return (
    <>
      <header className="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear">
        <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
          <SidebarTrigger className="-ml-1" />
          <Separator
            orientation="vertical"
            className="mx-2 data-[orientation=vertical]:h-4"
          />
          <h1 className="text-base font-medium">{title}</h1>
          <div className="ml-auto flex items-center gap-2">
            {/* Messages Button */}
            <Button
              variant="ghost"
              size="icon"
              className="relative"
              onClick={() => setIsMessagesOpen(true)}
            >
              <MessageCircle className="h-5 w-5" />
              {unreadCount > 0 && (
                <span className="absolute -top-1 -right-1 h-5 min-w-5 rounded-full bg-destructive text-destructive-foreground text-xs flex items-center justify-center px-1">
                  {unreadCount > 99 ? "99+" : unreadCount}
                </span>
              )}
            </Button>
            
            {/* Notifications */}
            <NotificationDrawer />
            
            {/* Theme Toggle */}
            <ThemeToggleIcon />
          </div>
        </div>
      </header>

      {/* Messages Drawer */}
      <Drawer open={isMessagesOpen} onOpenChange={setIsMessagesOpen}>
        <DrawerContent className="max-w-6xl mx-auto h-[80vh]">
          <DrawerHeader className="pb-4">
            <DrawerTitle>Messages</DrawerTitle>
          </DrawerHeader>
          <div className="flex h-full overflow-hidden">
            <div className="w-80 border-r">
              <MessageSidebar user={user} />
            </div>
            <div className="flex-1">
              <MessageChat 
                currentUser={user} 
                conversation={selectedConversation} 
              />
            </div>
          </div>
        </DrawerContent>
      </Drawer>
    </>
  )
}
