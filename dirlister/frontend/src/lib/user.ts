import { BACKEND_BASE_URL } from "@/App";

export interface UserStruct {
  email: string;
  createdAt: string;

}


export interface RegisterData {
  email: string;
  password: string;
  code: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface ChangePasswordData {
  currentPassword: string;
  newPassword: string;
}

export interface ResetPasswordData {
  email: string;
  code: string;
  newPassword: string;
}

export interface ChangeUserStatus {
  newstatus: string,
  email: string
}
function logout(on_finish: () => void) {
  fetch(`${BACKEND_BASE_URL}/auth/logout`, {credentials: "include"}).then(r => {
    on_finish()
  })


}


export async function sendRequest<T>(url: string, data: T): Promise<any> {


  const response = await fetch(url, {
    credentials: "include",
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const errorText = await response.text().catch(() => "Ошибка сервера");
    throw new Error(errorText || `Ошибка запроса: ${response.status}`);
  }

  if (response.headers.get("content-type")?.includes("application/json")) {
    return response.json();
  }
  return response.text();
}

export const authService = {
  register: (data: RegisterData) => sendRequest(`${BACKEND_BASE_URL}/auth/register`, data),
  login: (data: LoginData) => sendRequest(`${BACKEND_BASE_URL}/auth/login`, data),
  changePassword: (data: ChangePasswordData) => sendRequest(`${BACKEND_BASE_URL}/auth/change`, data),
  resetPassword: (data: ResetPasswordData) => sendRequest(`${BACKEND_BASE_URL}/auth/reset`, data),
  deleteAccount: () => { fetch(`${BACKEND_BASE_URL}/auth/delete`) },
  sendCode: async (data: string) => { return await sendRequest(`${BACKEND_BASE_URL}/auth/code`, {email:data}) },
  logout: (on_finish: () => void) => { logout(on_finish) },
};