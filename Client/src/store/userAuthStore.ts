import { create } from "zustand";
import { api } from "../lib/utils";
import toast from "react-hot-toast";

interface userData {
  email: string;
  exp: number;
}

interface AuthResponse {
  authenticated: boolean;
  message: string;
  user: userData;
}

export interface MessageResponse {
  message: string;
  success: string;
}

interface LoginData {
  email: string;
  password: string;
}

export interface RegisterData {
  fullName: string;
  email: string;
  password: string;
}

interface AuthState {
  authUser: userData | null;
  isLoading: boolean;
  checkAuth: () => Promise<void>;
  login: (data: LoginData) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  logout: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  authUser: null,
  isLoading: false,
  checkAuth: async () => {
    set({ isLoading: true });
    try {
      const res = await api.post<AuthResponse>("/auth/check");
      console.log("log de checkAuth store", res.data);
      set({ authUser: res.data.user, isLoading: false });
    } catch (error) {
      console.error("Error in checkAuth :", error);
      // Optionnel : réinitialiser l'utilisateur en cas d'erreur
      set({ authUser: null });
    }
  },

  login: async (data) => {
    set({ isLoading: true });
    try {
      const res = await api.post<AuthResponse>("/login", data);
      console.log("log de login store", res.data);
      set({ authUser: res.data.user, isLoading: true });
      toast.success("Bon retour parmi nous!🎉");
    } catch (error) {
      console.error("Error in login :", error);
      set({ authUser: null });
      toast.error("Mot de passe ou email incorrect");
    }
  },

  register: async (data) => {
    set({ isLoading: true });
    try {
      const res = await api.post<AuthResponse>("/register", data);
      console.log("log de register store", res.data);
      set({ authUser: res.data.user, isLoading: false });
      toast.success("Bienvenu parmi nous!🎉");
    } catch (error) {
      console.error("Error in register :", error);
      set({ authUser: null });
      toast.error("Mot de passe ou email incorrect");
    }
  },

  logout: async () => {
    set({ isLoading: true });
    try {
      const res = await api.post<MessageResponse>("/logout");
      console.log("lo de logout store ", res.data.message);
      toast.success(res.data.message);
      set({ authUser: null });
    } catch (error) {
      console.error("Error in login :", error);
      set({ authUser: null });
    }
  },
}));
