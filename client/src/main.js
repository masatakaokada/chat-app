import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import { initializeApp } from "firebase/app";
import { getAuth, onAuthStateChanged } from "firebase/auth";
import uuid from "vue-uuid";
import vuetify from "./plugins/vuetify";
import { loadFonts } from "./plugins/webfontloader";

const firebaseConfig = {
  apiKey: process.env.VUE_APP_API_KEY,
  authDomain: process.env.VUE_APP_AUTH_DOMAIN,
  projectId: process.env.VUE_APP_PROJECT_ID,
  storageBucket: process.env.VUE_APP_STORAGE_BUCKET,
  messagingSenderId: process.env.VUE_APP_MESSAGING_SENDER_ID,
  appId: process.env.VUE_APP_APP_ID,
};

initializeApp(firebaseConfig);

loadFonts();

let app;
const auth = getAuth();
onAuthStateChanged(auth, () => {
  if (!app) {
    createApp(App).use(router).use(uuid).use(vuetify).mount("#app");
  }
});
