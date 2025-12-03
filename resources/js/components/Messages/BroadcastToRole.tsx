"use client";

import * as React from "react";
import { Megaphone, Send, Loader2, Users } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import axios from "@/lib/axios";

interface Role {
  id: number;
  name: string;
  slug: string;
  level: number;
  is_active: boolean;
}

interface BroadcastToRoleProps {
  onSuccess?: () => void;
}

export function BroadcastToRole({ onSuccess }: BroadcastToRoleProps) {
  const [roles, setRoles] = React.useState<Role[]>([]);
  const [selectedRoleId, setSelectedRoleId] = React.useState<string>("");
  const [message, setMessage] = React.useState("");
  const [loading, setLoading] = React.useState(false);
  const [loadingRoles, setLoadingRoles] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [success, setSuccess] = React.useState<string | null>(null);

  // Load roles on mount
  React.useEffect(() => {
    const loadRoles = async () => {
      try {
        setLoadingRoles(true);
        const response = await axios.get("/api/roles?pageSize=100");
        if (response.data.success && response.data.data?.data) {
          setRoles(response.data.data.data.filter((r: Role) => r.is_active));
        }
      } catch (err) {
        console.error("Failed to load roles:", err);
        setError("Failed to load roles");
      } finally {
        setLoadingRoles(false);
      }
    };

    loadRoles();
  }, []);

  const handleSend = async () => {
    if (!selectedRoleId || !message.trim()) {
      setError("Please select a role and enter a message");
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await axios.post("/api/messages/broadcast-to-role", {
        role_id: parseInt(selectedRoleId),
        content: message.trim(),
      });

      if (response.data.success) {
        const data = response.data.data;
        setSuccess(
          `Message sent to ${data.sent_count} user${data.sent_count !== 1 ? "s" : ""} in "${data.role_name}" role`
        );
        setMessage("");
        setSelectedRoleId("");
        onSuccess?.();
      } else {
        setError(response.data.message || "Failed to send broadcast");
      }
    } catch (err: any) {
      setError(err.response?.data?.message || "Failed to send broadcast");
    } finally {
      setLoading(false);
    }
  };

  const selectedRole = roles.find((r) => r.id.toString() === selectedRoleId);

  return (
    <div className="flex flex-col h-full p-4">
      <div className="flex items-center gap-2 mb-4">
        <Megaphone className="h-5 w-5 text-primary" />
        <h2 className="font-semibold">Broadcast to Role</h2>
      </div>

      <p className="text-sm text-muted-foreground mb-4">
        Send a message to all users with a specific role. This is only available
        to super administrators.
      </p>

      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {success && (
        <Alert className="mb-4 border-green-500 bg-green-50 dark:bg-green-950">
          <AlertDescription className="text-green-700 dark:text-green-300">
            {success}
          </AlertDescription>
        </Alert>
      )}

      <div className="space-y-4 flex-1">
        <div className="space-y-2">
          <Label htmlFor="role">Select Role</Label>
          {loadingRoles ? (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Loader2 className="h-4 w-4 animate-spin" />
              Loading roles...
            </div>
          ) : (
            <Select value={selectedRoleId} onValueChange={setSelectedRoleId}>
              <SelectTrigger>
                <SelectValue placeholder="Choose a role..." />
              </SelectTrigger>
              <SelectContent>
                {roles.map((role) => (
                  <SelectItem key={role.id} value={role.id.toString()}>
                    <div className="flex items-center gap-2">
                      <Users className="h-4 w-4" />
                      <span>{role.name}</span>
                      <Badge variant="secondary" className="text-xs">
                        Level {role.level}
                      </Badge>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        </div>

        {selectedRole && (
          <div className="p-3 bg-muted rounded-lg">
            <p className="text-sm">
              <strong>Selected:</strong> {selectedRole.name}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              This message will be sent to all active users with the "
              {selectedRole.name}" role.
            </p>
          </div>
        )}

        <div className="space-y-2">
          <Label htmlFor="message">Message</Label>
          <Textarea
            id="message"
            placeholder="Type your broadcast message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            rows={6}
            maxLength={5000}
            className="resize-none"
          />
          <p className="text-xs text-muted-foreground text-right">
            {message.length}/5000
          </p>
        </div>
      </div>

      <div className="mt-4 pt-4 border-t">
        <Button
          onClick={handleSend}
          disabled={loading || !selectedRoleId || !message.trim()}
          className="w-full"
        >
          {loading ? (
            <>
              <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              Sending...
            </>
          ) : (
            <>
              <Send className="h-4 w-4 mr-2" />
              Send Broadcast
            </>
          )}
        </Button>
      </div>
    </div>
  );
}
