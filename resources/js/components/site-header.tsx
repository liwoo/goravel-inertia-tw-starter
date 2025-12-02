import * as React from "react"
import { Separator } from "@/components/ui/separator"
import { SidebarTrigger } from "@/components/ui/sidebar"
import { ThemeToggleIcon } from "@/components/ThemeToggleIcon"
import { NotificationDrawer } from "@/components/Notifications/NotificationDrawer"
import { Button } from "@/components/ui/button"
import { MessageCircle } from "lucide-react"
import { MessageSidebar } from "@/components/Messages/MessageSidebar"
import { MessageChat } from "@/components/Messages/MessageChat"
import { BroadcastHistory } from "@/components/Messages/BroadcastHistory"
import { useMessages, MessageUser } from "@/contexts/MessageContext"
import { usePage } from "@inertiajs/react"
import { SharedData } from "@/types/app"
import { useUI } from "@/contexts/UIContext"
import {
  Drawer,
  DrawerContent,
} from "@/components/ui/drawer"

export function SiteHeader({title}: { title: string }) {
  const { props } = usePage<SharedData>();
  const user = props.auth?.user;
  const { selectedConversation, unreadCount, setSelectedConversation } = useMessages();
  const { isMessagesOpen, openMessages, closeMessages, isNotificationsOpen, openNotifications, closeNotifications } = useUI();
  const [showNewMessage, setShowNewMessage] = React.useState(false);
  const [isBroadcastMode, setIsBroadcastMode] = React.useState(false);
  const [broadcastRefreshTrigger, setBroadcastRefreshTrigger] = React.useState(0);

  // Handle compose click - switch to "New Message" view
  const handleComposeClick = () => {
    setShowNewMessage(true);
    setSelectedConversation(null);
  };

  // Handle broadcast sent - trigger refresh of broadcast history
  const handleBroadcastSent = () => {
    setBroadcastRefreshTrigger(prev => prev + 1);
  };

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
              onClick={openMessages}
            >
              <MessageCircle className="h-5 w-5" />
              {unreadCount > 0 && (
                <span className="absolute -top-1 -right-1 h-5 min-w-5 rounded-full bg-destructive text-destructive-foreground text-xs flex items-center justify-center px-1">
                  {unreadCount > 99 ? "99+" : unreadCount}
                </span>
              )}
            </Button>

            {/* Notifications */}
            <NotificationDrawer isOpen={isNotificationsOpen} onOpenChange={(open) => open ? openNotifications() : closeNotifications()} />

            {/* Theme Toggle */}
            <ThemeToggleIcon />
          </div>
        </div>
      </header>

      {/* Messages Drawer */}
      <Drawer open={isMessagesOpen} onOpenChange={(open) => open ? openMessages() : closeMessages()}>
        <DrawerContent className="max-w-5xl mx-auto h-[85vh]">
          <div className="flex h-full overflow-hidden">
            <div className="w-96 border-r shrink-0">
              {user && (
                <MessageSidebar
                  user={user as any}
                  showNewMessage={showNewMessage}
                  onShowNewMessageChange={setShowNewMessage}
                  onBroadcastModeChange={setIsBroadcastMode}
                  onBroadcastSent={handleBroadcastSent}
                />
              )}
            </div>
            <div className="flex-1 min-w-0">
              {isBroadcastMode ? (
                <BroadcastHistory refreshTrigger={broadcastRefreshTrigger} />
              ) : (
                user && (
                  <MessageChat
                    currentUser={user as MessageUser}
                    conversation={selectedConversation}
                    onComposeClick={handleComposeClick}
                  />
                )
              )}
            </div>
          </div>
        </DrawerContent>
      </Drawer>
    </>
  )
}
