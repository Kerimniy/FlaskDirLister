import { useEffect, useState } from "react";

import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/password-input";
import { Loader2 } from "lucide-react";
import type { ChangePasswordData } from "@/lib/user"
import { authService } from "@/lib/user"
import { useNavigate } from "react-router-dom";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { AppHeader } from "@/components/app-header";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useAuth } from '@/App';

export default function AccountPage() {
  const { user, updateUser } = useAuth();
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (user === null) {
      navigate("/.@/auth/login")
    }
  })


  return (

    <SidebarProvider className="w-full">
      <AppSidebar />



      <SidebarInset className="flex flex-col min-w-0">
        <AppHeader
          user={user}
          onLoginClick={() => console.log("login")}
          onProfileClick={() => console.log("profile")}
        />

        <div className="flex-col flex items-center gap-4 p-6">
          <Card className="w-full p-4 items-center">
            <AccountInfo />
          </Card>

          <Card className="w-full p-4 items-center">
            <ChangePassword />
          </Card>

          <Card className="w-full p-4 ">

            <div><Button onClick={() => { setIsLoading(true); authService.logout(() => setIsLoading(false)); updateUser(); navigate("/") }}>
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}

              Log out
            </Button></div>
          </Card>

        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}

function AccountInfo() {
  const { user } = useAuth()

  return (
    <section className="max-w-2xl flex flex-col items-center">

      <h2 className="text-lg font-semibold mb-4 ">
        Current info
      </h2>

      <div className="space-y-2">



        <div>
          <b>Email:</b> {user?.email}
        </div>

        <div>
          <b>Registration date:</b> {user?.createdAt.slice(0, 10)}
        </div>


      </div>

    </section>
  );
}

export function ChangePassword() {
  const [current, setCurrent] = useState("");
  const [next, setNext] = useState("");
  const [repeat, setRepeat] = useState("");

  const [isLoading, setIsLoading] = useState(false);
  const [status, setStatus] = useState<{ text: string; isError: boolean } | null>(null);

  const isPasswordValid = next.length >= 6;
  const isMatch = next === repeat;
  const isFormValid = current.length > 0 && isPasswordValid && isMatch;
  const isDisabled = isLoading || !isFormValid;

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setStatus(null);

    const payload: ChangePasswordData = {
      currentPassword: current,
      newPassword: next,
    };

    try {
      await authService.changePassword(payload);
      setStatus({ text: "Password changed", isError: false });
      setCurrent("");
      setNext("");
      setRepeat("");
    } catch (err: any) {
      setStatus({ text: err.message || "Password change failed", isError: true });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <section className="max-w-2xl w-full">
      <h2 className="text-lg font-semibold mb-4">Change password</h2>

      {status && (
        <div className={`p-3 mb-4 rounded-md text-sm font-medium border text-center ${status.isError
          ? 'bg-destructive/10 text-destructive border-destructive/20'
          : 'bg-green-500/10 text-green-600 border-green-500/20'
          }`}>
          {status.text}
        </div>
      )}

      <form onSubmit={handleChangePassword} className="space-y-3">
        <PasswordInput
          placeholder="Current password"
          value={current}
          onChange={(e) => setCurrent(e.target.value)}
        />

        <PasswordInput
          placeholder="New password (min 6 chars)"
          value={next}
          onChange={(e) => setNext(e.target.value)}
        />

        <PasswordInput
          placeholder="Repeat password"
          value={repeat}
          onChange={(e) => setRepeat(e.target.value)}
        />

        {next && repeat && !isMatch && (
          <p className="text-xs text-destructive">Passwords don't match</p>
        )}
        <div>
          <Button type="submit" disabled={isDisabled}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Change password
          </Button>
        </div>
      </form>
    </section>
  );
}

