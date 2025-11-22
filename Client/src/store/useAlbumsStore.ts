import { create } from "zustand";
import { api } from "../lib/utils";
import type { MessageResponse } from "./userAuthStore";
import toast from "react-hot-toast";

export interface AlbumData {
  id: number;
  title: string;
  artist: string;
  year: number;
  cover_url: string;
}

export type AlbumList = AlbumData[];

interface AlbumState {
  album: AlbumData | null;
  albums: AlbumList;
  isLoading: boolean;
  getAllAlbums: () => Promise<void>;
  getAlbumById: (id: number) => Promise<AlbumData | undefined>;
  createAlbum: (data: Partial<AlbumData> | FormData) => Promise<void>;
  updateAlbum: (
    data: Partial<AlbumData> | FormData,
    id: number
  ) => Promise<void>;
  deleteAlbum: (id: number) => Promise<void>;
}

export const useAlbumStore = create<AlbumState>((set) => ({
  album: null,
  albums: [],
  isLoading: false,

  getAllAlbums: async () => {
    set({ isLoading: true });

    try {
      const res = await api.get<AlbumList>("/albums"); // <-- un tableau
      set({
        albums: res.data,
        isLoading: false,
      });
      console.log("log de getAllAlbums store", res.data);
    } catch (error) {
      console.error("Error fetching albums :", error);
      set({ isLoading: false });
    }
  },

  getAlbumById: async (id) => {
    set({ isLoading: true });
    try {
      const res = await api.get<AlbumData>(`/albums/${id}`);
      set({
        album: res.data,
        isLoading: false,
      });
      console.log("log de getAlbumById store", res.data);
      return res.data;
    } catch (error) {
      console.error("Error fetching albums :", error);
    }
  },

  createAlbum: async (data: Partial<AlbumData> | FormData) => {
    set({ isLoading: true });
    try {
      let res;
      if (data instanceof FormData) {
        // Cas multipart/form-data (avec fichier)
        res = await api.post<AlbumData>("/albums", data, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });
      }

      console.log("log de createAlbum store", res?.data);
      set({ isLoading: false });
    } catch (error) {
      console.error("Error creating album:", error);
      set({ isLoading: false });
    }
  },

  updateAlbum: async (data: Partial<AlbumData> | FormData, id) => {
    set({ isLoading: true });
    try {
      let res;
      if (data instanceof FormData) {
        // Cas multipart/form-data (avec fichier)
        res = await api.put<AlbumData>(`/albums/${id}`, data, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });
      }
      console.log("log de createAlbum store", res?.data);
      set({ isLoading: false });
    } catch (error) {
      console.log("Error dans le store updateAlbum", error);
    }
  },

  deleteAlbum: async (id) => {
    set({ isLoading: true });
    try {
      const res = await api.delete<MessageResponse>(`/albums/${id}`);
      console.log("log de deleteAlbum store", res.data.message);
      set((state) => ({
        ...state,
        albums: [...state.albums.filter((album) => album.id !== id)],
        isLoading: false,
      }));
    } catch (error) {
      console.log("Error dans le store updateAlbum", error);
    }
  },
}));
