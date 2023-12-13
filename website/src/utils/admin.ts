import axios from "axios";
import { clearTokenFromLocalStorage, saveTokenToLocalStorage, sha256 } from ".";

const serverUrl = import.meta.env.VITE_API_URL;

let key: string | null = null;

let client = createApi();

function createApi() {
  const client = axios.create({
    baseURL: serverUrl,
    timeout: 32000,
    headers: {
      "Content-Type": "application/json",
    },
  });

  // Request Middleware
  client.interceptors.request.use(
    function (config) {
      // Config before Request
      config.params["key"] = key;
      return config;
    },
    function (err) {
      // If Request error
      console.error(err);
      return Promise.reject(err);
    }
  );

  // Response Middleware
  client.interceptors.response.use(
    function (res) {
      if (res.data.invalidToken) {
        admin.logout();
      }
      return res;
    },
    function (error) {
      // Do something with response error
      return Promise.reject(error);
    }
  );

  return client;
}

export function setKey(token: string) {
  key = token;
}

export function clearKey() {
  key = null;
}

const admin = {
  client: client,

  async login(username: string, password: string) {
    const key = await sha256(`${username}+${password}`);

    const response = await client.post(`/admin/validate-key?key=${key}`);

    if (response.data.error || response.data.Error) {
      throw new Error(response.data.error || response.data.Error);
    }

    if (response.status === 200) {
      setKey(key);
      saveTokenToLocalStorage(key);
      return true;
    }

    return false;
  },

  isAuthed() {
    return key ? true : false;
  },

  logout() {
    clearTokenFromLocalStorage();
    clearKey();
    location.reload()
  },
};

export default admin;
