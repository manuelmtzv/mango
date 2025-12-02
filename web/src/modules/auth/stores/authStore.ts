import { useStorage } from "@vueuse/core";
import { defineStore } from "pinia";
import { ref } from "vue";
import { useNotePaginationStore, useNoteStore } from "@/modules/note/stores";

export const useAuthStore = defineStore("auth", () => {
  const username = ref<string | undefined>("");
  const jwt = useStorage("jwt", "");

  const noteStore = useNoteStore();
  const notePaginatioStore = useNotePaginationStore();

  const isLoggedIn = () => {
    return !!jwt.value;
  };

  const setUser = (newUsername: string) => {
    username.value = newUsername;
    localStorage.setItem("username", newUsername);
  };

  const setJwt = (newJwt: string) => {
    jwt.value = newJwt;
    localStorage.setItem("jwt", newJwt);
  };

  const logout = () => {
    username.value = undefined;
    jwt.value = "";

    noteStore.resetState();
    notePaginatioStore.resetState();
  };

  return {
    username,
    jwt,
    isLoggedIn,
    setUser,
    setJwt,
    logout,
  };
});
