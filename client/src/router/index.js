import { createRouter, createWebHistory } from "vue-router";
import { getAuth } from "firebase/auth";
import Top from "@/components/Top";
import Signup from "@/components/Signup";
import Signin from "@/components/Signin";
import Chat from "@/components/Chat";
import RoomNew from "@/components/RoomNew";
import Room from "@/components/Room";

const routes = [
  {
    path: "/:pathMatch(.*)*",
    redirect: "/signin",
  },
  {
    path: "/",
    name: "Top",
    component: Top,
    meta: { requiresAuth: true },
  },
  {
    path: "/signup",
    name: "Signup",
    component: Signup,
  },
  {
    path: "/signin",
    name: "Signin",
    component: Signin,
  },
  {
    path: "/chat",
    name: "Chat",
    component: Chat,
    meta: { requiresAuth: true },
  },
  {
    path: "/rooms/new",
    component: RoomNew,
    meta: { requiresAuth: true },
  },
  {
    path: "/rooms/:id",
    component: Room,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

// firebaseの認証
router.beforeEach((to, _, next) => {
  let currentUser = getAuth().currentUser;
  let requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  if (requiresAuth && !currentUser) next("signin");
  else if (!requiresAuth && currentUser) next();
  else next();
});

export default router;
