import axios, { Axios, AxiosError } from "axios";
import { useAuthStore } from "@/modules/auth/stores/authStore";
import { useToast } from "vue-toast-notification";
import { pinia } from "@/config/setPinia";

const authStore = useAuthStore(pinia);

export const mangocatnotesApi: Axios = axios.create({
  baseURL: import.meta.env.VITE_MANGOCATAPI_URL,
});

mangocatnotesApi.interceptors.request.use(
  (config) => {
    if (authStore.jwt) {
      config.headers.Authorization = `Bearer ${authStore.jwt}`;
    }
    return config;
  },
  (err) => {
    return Promise.reject(err);
  },
);

mangocatnotesApi.interceptors.response.use(
  (response) => {
    return response;
  },
  (err: AxiosError) => {
    if (err.code == AxiosError.ERR_NETWORK)
      useToast().open({
        message: "Server is currently unavailable. Please try again later.",
        type: "error",
        duration: 3500,
      });

    return Promise.reject(err);
  },
);
