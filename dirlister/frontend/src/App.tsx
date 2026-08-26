import { useState, useEffect } from 'react'

import './App.css'
import { createContext, useContext } from "react";

import { ThemeProvider } from "next-themes";
import { useTheme } from "next-themes";
import type { UserStruct } from './lib/user';
import { Routes, Route, useNavigate } from "react-router-dom";
import type { ReactNode } from 'react';
import IndexPage from './pages';

import {
  SidebarProvider,
} from "@/components/ui/sidebar";
import EditPage from './pages/editor';
import Auth from './pages/Auth';
import AccountPage from './pages/account';
import FileUpload from './pages/upload';
import RulesPage from './pages/rules';

export const AuthContext = createContext<{
  user: UserStruct | null;
  setUser: React.Dispatch<React.SetStateAction<UserStruct | null>>;
  updateUser: () => Promise<void>;
} | null>(null);

export const BACKEND_BASE_URL = "https://127.0.0.1:6060"
export const RESULTS_PER_PAGE = 12

function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserStruct | null>(null);
  const navigate = useNavigate()
  useEffect(() => {
    fetch(`${BACKEND_BASE_URL}/auth/me`, {
      credentials: "include"
    })
      .then(r => {

        if (r.ok) {
          if (r.headers.get("content-type")?.includes("application/json")) {
            return r.json()
          }
          else {
            return r.text()
          }
        }
        else if (r.status===401){
        }
        return null;
      })
      .then((data) => {
        if (typeof (data) === "string" && data === "no-user") {
          navigate("/.@/auth/sign-up")
        } else {
          setUser(data);
        }
      });
  }, []);


  const updateUser = async () => {
    const r = await fetch(`${BACKEND_BASE_URL}/auth/me`, {
      credentials: "include"
    });

    if (!r.ok) {
      setUser(null);
      return;
    }
    if (r.ok) {
      if (r.headers.get("content-type")?.includes("application/json")) {
        setUser(await r.json());
      }
      else {
        navigate("/.@/auth/sign-up")
        setUser(null)
      }
    };

  };
  return (
    <AuthContext.Provider value={{ user, setUser, updateUser }}>
      {children}
    </AuthContext.Provider>
  );
}


export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("No AuthProvider");
  return ctx;
}


function App() {

  //  включить loading editor

  return (<div className="flex h-screen w-full">
    <ThemeProvider attribute="class"
      defaultTheme="system"
      enableSystem>


      <AuthProvider>

        <Routes>
          <Route path='/.@/edit' element={<EditPage />}></Route>
          <Route path='/.@/create' element={<EditPage />}></Route>
          <Route path='/.@/upload' element={<FileUpload />}></Route>
          <Route path='/.@/rules' element={<RulesPage />}></Route>
          <Route path="/.@/auth/*" element={<Auth />}></Route>
          <Route path="/.@/account" element={<AccountPage />}></Route>
          <Route path='/*' element={<IndexPage />}></Route>


        </Routes>



      </AuthProvider>


    </ThemeProvider>        </div>

  )
}

export default App
